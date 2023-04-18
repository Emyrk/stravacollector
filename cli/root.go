package cli

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "strava",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	cmd.AddCommand(serverCmd())

	return cmd
}

func getLogger(cmd *cobra.Command) zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	return logger
}
