package cli

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Emyrk/strava/strava"
	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	var (
		token string
	)

	cmd := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			if token == "" {
				log.Fatal("--access-token is not set")
			}
			client := strava.New(token)
			segment, err := client.AthleteSegmentEfforts(ctx, 16659489, 2)
			fmt.Println(err)
			d, _ := json.Marshal(segment)
			fmt.Println(string(d))
		},
	}

	cmd.Flags().StringVar(&token, "access-token", "", "Strava access token")

	return cmd
}
