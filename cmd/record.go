package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/natori-hrj/sre-metrics-cli/internal/slo"
	"github.com/natori-hrj/sre-metrics-cli/internal/storage"
	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record SLI data point",
	Example: `  slo record --service natorium-dev --total 10000 --success 9990
  slo record --service api-gateway --total 5000 --success 4998`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service, _ := cmd.Flags().GetString("service")
		total, _ := cmd.Flags().GetInt64("total")
		success, _ := cmd.Flags().GetInt64("success")

		if total <= 0 {
			return fmt.Errorf("--total must be a positive integer")
		}
		if success < 0 || success > total {
			return fmt.Errorf("--success must be between 0 and --total (%d)", total)
		}

		record := slo.Record{
			Timestamp: time.Now(),
			Total:     total,
			Success:   success,
		}

		if err := storage.Append(service, record); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(renderRecordSuccess(service, total, success))
		return nil
	},
}

func init() {
	recordCmd.Flags().String("service", "", "Service name (e.g. natorium-dev)")
	recordCmd.Flags().Int64("total", 0, "Total number of requests")
	recordCmd.Flags().Int64("success", 0, "Number of successful requests")
	_ = recordCmd.MarkFlagRequired("service")
	_ = recordCmd.MarkFlagRequired("total")
	_ = recordCmd.MarkFlagRequired("success")
	rootCmd.AddCommand(recordCmd)
}
