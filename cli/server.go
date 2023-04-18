package cli

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Emyrk/strava/api"
	"golang.org/x/oauth2"

	"github.com/Emyrk/strava/database"

	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	var (
		//token string
		dbURL    string
		secret   string
		clientID string
	)

	cmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := getLogger(cmd)
			ctx := cmd.Context()
			if secret == "" || clientID == "" {
				return fmt.Errorf("missing client id or secret")
			}

			db, err := database.NewPostgresDB(ctx, logger, dbURL)
			if err != nil {
				return fmt.Errorf("connect to postgres: %w", err)
			}
			var _ = db

			accessURL := "http://localhost:8000"

			srv, err := api.New(api.Options{
				OAuthCfg: &oauth2.Config{
					ClientID:     clientID,
					ClientSecret: secret,
					Endpoint: oauth2.Endpoint{
						AuthURL:   "https://www.strava.com/oauth/authorize",
						TokenURL:  "https://www.strava.com/oauth/token",
						AuthStyle: 0,
					},
					RedirectURL: fmt.Sprintf("%s/oauth2/callback", accessURL),
					// Must be comma joined
					Scopes: []string{strings.Join([]string{"read", "read_all", "profile:read_all", "activity:read"}, ",")},
				},
			})
			if err != nil {
				return fmt.Errorf("create server: %w", err)
			}

			url := srv.OAuthCfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
			logger.Info().Msg(fmt.Sprintf("Visit the URL for the auth dialog: %s", url))

			hsrv := &http.Server{
				Addr:    "0.0.0.0:8000",
				Handler: srv.Handler,
				BaseContext: func(listener net.Listener) context.Context {
					return ctx
				},
			}
			return hsrv.ListenAndServe()

			//client := strava.New(token)
			//segment, err := client.AthleteSegmentEfforts(ctx, 16659489, 2)
			//fmt.Println(err)
			//d, _ := json.Marshal(segment)
			//fmt.Println(string(d))
		},
	}

	cmd.Flags().StringVar(&secret, "oauth-secret", "", "Strava oauth app secret")
	cmd.Flags().StringVar(&clientID, "oauth-client-id", "", "Strava oauth app client ID")
	//cmd.Flags().StringVar(&token, "access-token", "", "Strava access token")
	cmd.Flags().StringVar(&dbURL, "db-url", "postgres://postgres:postgres@localhost:5432/strava?sslmode=disable", "Database URL")

	return cmd
}
