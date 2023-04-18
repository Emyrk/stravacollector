package cli

import (
	"encoding/json"
	"fmt"

	"github.com/Emyrk/strava/database"

	"github.com/Emyrk/strava/strava"
	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	var (
		token string
		dbURL string
	)

	cmd := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			logger := getLogger(cmd)
			ctx := cmd.Context()
			if token == "" {
				logger.Fatal().Msg("--access-token is not set")
			}

			db, err := database.NewPostgresDB(ctx, logger, dbURL)
			if err != nil {
				logger.Fatal().Err(err).Msg("connect to postgres")
			}
			var _ = db

			client := strava.New(token)
			segment, err := client.AthleteSegmentEfforts(ctx, 16659489, 2)
			fmt.Println(err)
			d, _ := json.Marshal(segment)
			fmt.Println(string(d))
		},
	}

	cmd.Flags().StringVar(&token, "access-token", "", "Strava access token")
	cmd.Flags().StringVar(&dbURL, "db-url", "postgres://postgres:postgres@localhost:5432/strava?sslmode=disable", "Database URL")

	return cmd
}
