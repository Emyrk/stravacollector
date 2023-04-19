package cli

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/pflag"

	"github.com/spf13/viper"

	"github.com/Emyrk/strava/api"

	"github.com/Emyrk/strava/database"

	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	var (
		dbURL       string
		secret      string
		clientID    string
		port        int
		accessURL   string
		config      string
		writeConfig bool
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

			srv, err := api.New(api.Options{
				OAuth: api.OAuthOptions{
					ClientID: clientID,
					Secret:   secret,
				},
				DB:        db,
				Logger:    logger.With().Str("component", "api").Logger(),
				AccessURL: u,
			})
			if err != nil {
				return fmt.Errorf("create server: %w", err)
			}

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

			// TODO: Check for server up

			time.Sleep(time.Second)
			err = srv.StartWebhook(ctx)
			if err != nil {
				return fmt.Errorf("start webhook: %w", err)
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

			return nil
		},
	}

	cmd.Flags().BoolVar(&writeConfig, "write-config", false, "Write config file and exit")
	cmd.Flags().StringVar(&config, "config", "", "Config file")
	cmd.Flags().StringVar(&accessURL, "access-url", "", "External url to talk with")
	cmd.Flags().IntVar(&port, "port", 9090, "Port to listen on")
	cmd.Flags().StringVar(&secret, "oauth-secret", "", "Strava oauth app secret")
	cmd.Flags().StringVar(&clientID, "oauth-client-id", "", "Strava oauth app client ID")
	//cmd.Flags().StringVar(&token, "access-token", "", "Strava access token")
	cmd.Flags().StringVar(&dbURL, "db-url", "postgres://postgres:postgres@localhost:5432/strava?sslmode=disable", "Database URL")

	return cmd
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
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
