package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/natori-hrj/sre-metrics-cli/internal/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a service configuration",
	Example: `  slo init --service natorium`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service, _ := cmd.Flags().GetString("service")
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Initializing service: %s\n\n", service)

		cfg := storage.ServiceConfig{}

		cfg.Name = prompt(reader, "Display name (e.g. natorium.dev)")
		cfg.URL = prompt(reader, "URL (e.g. https://natorium.dev)")

		fmt.Println("\n--- Vercel (leave blank to skip) ---")
		cfg.VercelToken = prompt(reader, "Vercel API token")
		if cfg.VercelToken != "" {
			cfg.VercelProjectID = prompt(reader, "Vercel Project ID")
			cfg.VercelTeamID = prompt(reader, "Vercel Team ID (blank for hobby)")
		}

		fmt.Println("\n--- UptimeRobot (leave blank to skip) ---")
		cfg.UptimeRobotKey = prompt(reader, "UptimeRobot API key")
		if cfg.UptimeRobotKey != "" {
			cfg.UptimeRobotMonID = prompt(reader, "UptimeRobot Monitor ID")
		}

		if err := storage.SaveConfig(service, cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}

		fmt.Println()
		fmt.Println(goodStyle.Render(fmt.Sprintf("Service %q initialized!", service)))

		if cfg.VercelToken != "" {
			fmt.Println("  Vercel:      configured")
		}
		if cfg.UptimeRobotKey != "" {
			fmt.Println("  UptimeRobot: configured")
		}
		fmt.Println("\nRun 'slo fetch --service " + service + "' to pull data.")
		return nil
	},
}

func prompt(reader *bufio.Reader, label string) string {
	fmt.Printf("  %s: ", label)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func init() {
	initCmd.Flags().String("service", "", "Service identifier (e.g. natorium)")
	_ = initCmd.MarkFlagRequired("service")
	rootCmd.AddCommand(initCmd)
}
