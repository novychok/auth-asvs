package cmd

import (
	"context"

	"github.com/novychok/authasvs/internal"
	"github.com/spf13/cobra"
)

var authApiCmd = &cobra.Command{
	Use: "authapi",
	RunE: func(cmd *cobra.Command, _ []string) error {
		app, cleanup, err := internal.Init()
		if err != nil {
			return err
		}
		defer cleanup()

		return app.StartAuthApiV1(context.TODO())
	},
}

func InitAuthApiCommands() {
	rootCmd.AddCommand(authApiCmd)
}
