package cmd

import (
	"fmt"
	"os"

	"github.com/natori-hrj/sre-metrics-cli/internal/slo"
	"github.com/natori-hrj/sre-metrics-cli/internal/storage"
	"github.com/spf13/cobra"
)

var budgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Show error budget remaining",
	Example: `  slo budget --service natorium-dev --target 99.9
  slo budget --service api-gateway --target 99.9 --window 30d`,
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
		fmt.Println(renderBudgetDetail(service, status, windowStr))
		return nil
	},
}

func renderBudgetDetail(service string, status slo.Status, window string) string {
	allowedErrors := int64(float64(status.TotalRequests) * (100 - status.Target) / 100)
	actualErrors := status.TotalRequests - status.SuccessRequests
	remaining := allowedErrors - actualErrors
	if remaining < 0 {
		remaining = 0
	}

	content := fmt.Sprintf(
		"Service:         %s\n"+
			"Window:          %s\n"+
			"SLO Target:      %.2f%%\n"+
			"Current SLI:     %.2f%%\n"+
			"Allowed Errors:  %d\n"+
			"Actual Errors:   %d\n"+
			"Budget Left:     %d requests\n"+
			"Budget Percent:  %s",
		service,
		window,
		status.Target,
		status.SLI,
		allowedErrors,
		actualErrors,
		remaining,
		renderBudget(status.ErrorBudgetPercent),
	)

	title := titleStyle.Render("Error Budget")
	return borderStyle.Render(centerText(title, 38) + "\n\n" + content)
}

func init() {
	budgetCmd.Flags().String("service", "", "Service name (e.g. natorium-dev)")
	budgetCmd.Flags().Float64("target", 99.9, "SLO target percentage")
	budgetCmd.Flags().String("window", "30d", "Time window (7d, 30d, 90d)")
	_ = budgetCmd.MarkFlagRequired("service")
	rootCmd.AddCommand(budgetCmd)
}
