package cli

import (
	"encoding/base64"

	"github.com/spf13/cobra"
)

func parseDump() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dump",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			cmd.Println(string(out))
			return nil
		},
	}

	return cmd
}
