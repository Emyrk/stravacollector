package cli

import (
	"io"
	"os"
	"strconv"

	"github.com/google/uuid"

	"github.com/hirosassa/zerodriver"
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

	cmd.AddCommand(
		serverCmd(),
		generateKey(),
	)

	return cmd
}

func getLogger(cmd *cobra.Command) zerolog.Logger {
	useStackDriver, _ := cmd.Flags().GetBool("stack-driver")

	var out io.Writer = zerolog.ConsoleWriter{Out: os.Stderr}
	if ok, _ := strconv.ParseBool(os.Getenv("STRAVA_JSON_LOGS")); ok {
		out = os.Stderr
	}

	var logger zerolog.Logger
	if useStackDriver {
		logger = *(zerodriver.NewDevelopmentLogger().Logger)
		logger.Output(out)
	} else {
		logger = zerolog.New(out).With().Timestamp().Logger()
	}
	logger = logger.With().Str("deployment_id", uuid.NewString()).Logger()
	return logger
}
