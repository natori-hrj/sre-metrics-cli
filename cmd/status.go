package cmd

import (
	"fmt"
	"os"

	"github.com/natori-hrj/sre-metrics-cli/internal/slo"
	"github.com/natori-hrj/sre-metrics-cli/internal/storage"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current SLO compliance status",
	Example: `  slo status --service natorium-dev --target 99.9
  slo status --service api-gateway --target 99.95 --window 7d`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service, _ := cmd.Flags().GetString("service")
		target, _ := cmd.Flags().GetFloat64("target")
		windowStr, _ := cmd.Flags().GetString("window")

		if target <= 0 || target > 100 {
			return fmt.Errorf("--target must be between 0 and 100")
		}

		window, err := slo.ParseWindow(windowStr)
		if err != nil {
			return err
		}

		records, err := storage.LoadAll(service)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(records) == 0 {
			return fmt.Errorf("no records found for service %q. Use 'slo record' to add data", service)
		}

		status := slo.Calculate(records, target, window)
		fmt.Println(renderStatus(service, status))
		return nil
	},
}

func init() {
	statusCmd.Flags().String("service", "", "Service name (e.g. natorium-dev)")
	statusCmd.Flags().Float64("target", 99.9, "SLO target percentage")
	statusCmd.Flags().String("window", "30d", "Time window (7d, 30d, 90d)")
	_ = statusCmd.MarkFlagRequired("service")
	rootCmd.AddCommand(statusCmd)
}
