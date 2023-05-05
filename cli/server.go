package cli

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/Emyrk/strava/api"
	"github.com/Emyrk/strava/api/queue"
	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/strava/stravalimit"
)

func serverCmd() *cobra.Command {
	var (
		dbURL             string
		secret            string
		clientID          string
		port              int
		accessURL         string
		config            string
		writeConfig       bool
		stackDriver       bool
		verifyToken       string
		disableWebhooks   bool
		signingSecret     string
		prometheusEnabled bool
		promtheusAddress  string
	)

	v := viper.New()
	cmd := &cobra.Command{
		Use: "server",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			v.SetConfigType("yaml")
			v.SetConfigName("strava.yaml")
			v.AddConfigPath(".")

			if err := v.ReadInConfig(); err != nil {
				if config != "" {
					return err
				}
				// It's okay if there isn't a config file
				if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
					return err
				}
			}

			// When we bind flags to environment variables expect that the
			// environment variables are prefixed, e.g. a flag like --number
			// binds to an environment variable STING_NUMBER. This helps
			// avoid conflicts.
			v.SetEnvPrefix("STRAVA")

			// Environment variables can't have dashes in them, so bind them to their equivalent
			// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

			// Bind to environment variables
			// Works great for simple config names, but needs help for names
			// like --favorite-color which we fix in the bindFlags function
			v.AutomaticEnv()

			// Bind the current command's flags to viper
			bindFlags(cmd, v, writeConfig)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()
			registry := prometheus.NewRegistry()
			stravalimit.SetRegistry(registry)

			if writeConfig {
				return v.SafeWriteConfigAs(config)
			}

			logger := getLogger(cmd)
			if secret == "" || clientID == "" {
				return fmt.Errorf("missing client id or secret")
			}

			db, err := database.NewPostgresDB(ctx, logger, dbURL)
			if err != nil {
				return fmt.Errorf("connect to postgres: %w", err)
			}

			if accessURL == "" {
				accessURL = fmt.Sprintf("http://localhost:%d", port)
			}

			u, err := url.Parse(accessURL)
			if err != nil {
				return fmt.Errorf("parse access url: %w", err)
			}
			if !(u.Scheme == "http" || u.Scheme == "https") {
				return fmt.Errorf("access url scheme must be http or https")
			}

			secPem, err := base64.StdEncoding.DecodeString(strings.TrimSpace(signingSecret))
			if err != nil {
				return fmt.Errorf("decode signing key: %w", err)
			}
			srv, err := api.New(api.Options{
				OAuth: api.OAuthOptions{
					ClientID: clientID,
					Secret:   secret,
				},
				DB:            db,
				Logger:        logger.With().Str("component", "api").Logger(),
				AccessURL:     u,
				VerifyToken:   verifyToken,
				SigningKeyPEM: secPem,
				Registry:      registry,
			})
			if err != nil {
				return fmt.Errorf("create server: %w", err)
			}

			manager, err := queue.New(ctx, queue.Options{
				DBURL:    dbURL,
				Logger:   logger.With().Str("component", "queue").Logger(),
				DB:       db,
				OAuthCfg: srv.OAuthConfig,
				Registry: registry,
			})
			if err != nil {
				return fmt.Errorf("create queue manager: %w", err)
			}

			err = manager.Run(ctx)
			if err != nil {
				return fmt.Errorf("run queue: %w", err)
			}
			defer manager.Close()

			hsrv := &http.Server{
				Addr:    fmt.Sprintf("0.0.0.0:%d", port),
				Handler: srv.Handler,
				BaseContext: func(listener net.Listener) context.Context {
					return ctx
				},
			}

			go func() {
				err := hsrv.ListenAndServe()
				if err != nil {
					logger.Error().Err(err).Msg("http server error")
				}
			}()
			logger.Info().
				Int("port", port).
				Str("access_url", accessURL).Msg("Server running")

			stravaRateLimitLog(ctx, logger)
			if prometheusEnabled {
				launchPrometheus(ctx, logger, promtheusAddress, registry)
			}

			if !disableWebhooks {
				lastPrint := time.Time{}
				for {
					health := fmt.Sprintf("%s/myhealthz", strings.TrimSuffix(accessURL, "/"))
					select {
					case <-ctx.Done():
						return fmt.Errorf("server did not start in time: %s", health)
					default:
					}

					resp, err := http.Get(health)
					if err == nil && resp.StatusCode == http.StatusOK {
						break
					}
					if time.Since(lastPrint) > time.Second*10 {
						logger.Info().
							Str("url", health).
							Msg("Server not responding, cannot start webhook")
						lastPrint = time.Now()
					}
					time.Sleep(time.Second * 1)
				}

				logger.Info().Msg("Server is up, starting webhook")
				eq, err := srv.StartWebhook(ctx)
				if err == nil {
					logger.Info().Msgf("Webhook started to %s", srv.Events.Callback.String())
					go func() {
						manager.HandleWebhookEvents(ctx, eq)
					}()
				}
				if err != nil {
					now := time.Now()
					// This sucks but prevents endless loop that uses all our api limits.
					go func() {
						for {
							logger.Error().
								Str("callback", srv.Events.Callback.String()).
								Str("time", now.Format(time.RFC3339)).
								Err(err).
								Msg("Webhook failed to start, restart the server to try again")
							time.Sleep(time.Second * 10)
						}
					}()
				}
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			// CLOSE
			<-c
			logger.Info().Msg("Gracefully shutting down...")
			cancel()

			tmp, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			err = hsrv.Shutdown(tmp)
			if err != nil {
				logger.Error().Err(err).Msg("http server shutdown error")
			}
			manager.Close()

			return nil
		},
	}

	cmd.Flags().BoolVar(&writeConfig, "write-config", false, "Write config file and exit")
	cmd.Flags().StringVar(&config, "config", "", "Config file")
	cmd.Flags().StringVar(&accessURL, "access-url", "", "External url to talk with")
	cmd.Flags().IntVar(&port, "port", 9090, "Port to listen on")
	cmd.Flags().StringVar(&secret, "oauth-secret", "", "Strava oauth app secret")
	cmd.Flags().StringVar(&clientID, "oauth-client-id", "", "Strava oauth app client ID")
	cmd.Flags().StringVar(&dbURL, "db-url", "postgres://postgres:postgres@localhost:5432/strava?sslmode=disable", "Database URL")
	cmd.Flags().BoolVar(&stackDriver, "stack-driver", false, "Export stack driver logs")
	cmd.Flags().StringVar(&verifyToken, "verify-token", "", "Strava webhook verify token")
	cmd.Flags().BoolVar(&disableWebhooks, "disable-webhooks", false, "Useful for running a server without a public url")
	cmd.Flags().StringVar(&signingSecret, "signing-secret", "", "RSA signing key base64 encoded")
	cmd.Flags().BoolVar(&prometheusEnabled, "enable-prometheus", false, "Enable prometheus metrics")
	cmd.Flags().StringVar(&promtheusAddress, "prometheus-address", "0.0.0.0:9091", "Prometheus address to listen on")

	return cmd
}

func stravaRateLimitLog(ctx context.Context, logger zerolog.Logger) {
	go func() {
		logger.Debug().Msg("Will watch strava rate limits every minute.")
		ticker := time.NewTicker(time.Minute * 10)
		for {
			select {
			case <-ctx.Done():
				logger.Debug().Msg("Stopping strava rate limit watcher.")
				return
			case <-ticker.C:
				i, d := stravalimit.Remaining()
				logger.Debug().
					Int64("IntervalLeft", i).
					Int64("DailyLeft", d).
					Msg("Strava Rate Limits")
			}
		}
	}()
}

func launchPrometheus(ctx context.Context, logger zerolog.Logger, address string, registry *prometheus.Registry) {
	srv := http.Server{
		Addr:    address,
		Handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}
	go func() {
		logger.Info().Str("address", address).Msg("Starting prometheus server")
		err := srv.ListenAndServe()
		if err != nil {
			logger.Error().Str("service", "prometheus").Err(err).Msg("prometheus server error")
		}
	}()
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper, always bool) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name
		// If using camelCase in the config file, replace hyphens with a camelCased string.
		// Since viper does case-insensitive comparisons, we don't need to bother fixing the case, and only need to remove the hyphens.
		//if replaceHyphenWithCamelCase {
		//	configName = strings.ReplaceAll(f.Name, "-", "")
		//}

		if always {
			v.Set(configName, f.Value)
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
