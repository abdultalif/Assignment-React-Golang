package cmd

import (
	"backend-service/internal/app"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start",
	Long:  "Start",
	Run: func(cmd *cobra.Command, args []string) {
		app.RunServer()
	},
}
