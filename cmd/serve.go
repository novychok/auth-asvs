package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		// app, cleanup, err := internal.Init()
		// if err != nil {
		// 	return err
		// }
		// defer cleanup()

		go func() {
			// err := app.StartAdminAPIV1(cmd.Context())
			// if err != nil {
			// 	panic(err)
			// }
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		return nil
	},
}

func InitServeCommands() {
	rootCmd.AddCommand(serveCmd)
}
