package cli

import (
	"io"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "strava",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(serverCmd())

	return cmd
}

func getLogger(cmd *cobra.Command) zerolog.Logger {
	var out io.Writer = zerolog.ConsoleWriter{Out: os.Stderr}
	if ok, _ := strconv.ParseBool(os.Getenv("STRAVA_JSON_LOGS")); ok {
		out = os.Stderr
	}

	logger := zerolog.New(out).With().Timestamp().Logger()
	return logger
}
