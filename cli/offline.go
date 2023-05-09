package cli

import (
	"log"
	"net/http"

	"github.com/Emyrk/strava/internal/offlineserver"
	"github.com/spf13/cobra"
)

func offlineServer() *cobra.Command {
	var (
		port int
	)
	cmd := &cobra.Command{
		Use: "offline",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := offlineserver.New()
			if err != nil {
				return err
			}
			log.Printf("Starting server on http://localhost:%d\n", port)
			return http.ListenAndServe(":7000", srv.Router)
		},
	}
	cmd.Flags().IntVar(&port, "port", 7000, "port to run on")

	return cmd
}
