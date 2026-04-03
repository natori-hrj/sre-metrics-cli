package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slo",
	Short: "SLI/SLO calculator and tracker",
	Long:  "A CLI tool for recording SLI metrics, tracking SLO compliance, and monitoring error budgets.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
