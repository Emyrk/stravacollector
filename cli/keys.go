package cli

import (
	"encoding/base64"
	"fmt"

	"github.com/Emyrk/strava/api/auth/authkeys"
	"github.com/spf13/cobra"
)

func generateKey() *cobra.Command {
	var (
		asBase64 bool
	)
	cmd := &cobra.Command{
		Use: "gen-key",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := authkeys.GenerateKey()
			if err != nil {
				return fmt.Errorf("error generating key: %w", err)
			}

			data := authkeys.MarshalPrivateKey(key)
			str := string(data)
			if asBase64 {
				str = base64.StdEncoding.EncodeToString(data)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), str)
			return nil
		},
	}
	cmd.Flags().BoolVar(&asBase64, "base64", false, "Output the key as a base64 encoded string")

	return cmd
}
