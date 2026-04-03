package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/natori-hrj/sre-metrics-cli/internal/fetcher"
	"github.com/natori-hrj/sre-metrics-cli/internal/slo"
	"github.com/natori-hrj/sre-metrics-cli/internal/storage"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch SLI data from configured sources",
	Example: `  slo fetch --service natorium
  slo fetch --service natorium --window 7d`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service, _ := cmd.Flags().GetString("service")
		windowStr, _ := cmd.Flags().GetString("window")

		window, err := slo.ParseWindow(windowStr)
		if err != nil {
			return err
		}

		cfg, err := storage.LoadConfig(service)
		if err != nil {
			return err
		}

		since := time.Now().Add(-window)
		days := int(window.Hours() / 24)
		fetched := false

		// Fetch from Vercel
		if cfg.VercelToken != "" && cfg.VercelProjectID != "" {
			fmt.Print("Fetching Vercel deployments... ")
			records, err := fetcher.FetchVercelDeployments(cfg.VercelToken, cfg.VercelProjectID, cfg.VercelTeamID, since)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nVercel error: %v\n", err)
			} else {
				for _, r := range records {
					if err := storage.Append(service+"-deploy", r); err != nil {
						return fmt.Errorf("save deploy record: %w", err)
					}
				}
				fmt.Println(goodStyle.Render(fmt.Sprintf("OK (%d days)", len(records))))
				fetched = true
			}
		}

		// Fetch from UptimeRobot
		if cfg.UptimeRobotKey != "" && cfg.UptimeRobotMonID != "" {
			fmt.Print("Fetching UptimeRobot data... ")
			records, err := fetcher.FetchUptimeRobot(cfg.UptimeRobotKey, cfg.UptimeRobotMonID, days)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nUptimeRobot error: %v\n", err)
			} else {
				for _, r := range records {
					if err := storage.Append(service+"-uptime", r); err != nil {
						return fmt.Errorf("save uptime record: %w", err)
					}
				}
				fmt.Println(goodStyle.Render("OK"))
				fetched = true
			}
		}

		if !fetched {
			return fmt.Errorf("no data sources configured. Run 'slo init --service %s' to set up Vercel or UptimeRobot", service)
		}

		fmt.Println("\nRun these to see results:")
		fmt.Printf("  slo status --service %s-deploy --target 99.9\n", service)
		fmt.Printf("  slo status --service %s-uptime --target 99.9\n", service)
		return nil
	},
}

func init() {
	fetchCmd.Flags().String("service", "", "Service identifier")
	fetchCmd.Flags().String("window", "30d", "Time window (7d, 30d, 90d)")
	_ = fetchCmd.MarkFlagRequired("service")
	rootCmd.AddCommand(fetchCmd)
}
