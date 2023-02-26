package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start doing a new thing",
	Run:   runStart,
}

func runStart(cmd *cobra.Command, args []string) {
	fmt.Println("start")
}
