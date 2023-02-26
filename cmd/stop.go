package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop doing the thing you're doing",
	Run:   runStop,
}

func runStop(cmd *cobra.Command, args []string) {
	fmt.Println("stop")
}
