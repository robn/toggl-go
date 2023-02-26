package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(timerCmd)
}

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "what are you doing right now?",
	Run:   runTimer,
}

func runTimer(cmd *cobra.Command, args []string) {
	fmt.Println("timer")
}
