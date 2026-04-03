package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/natori-hrj/sre-metrics-cli/internal/slo"
)

type uptimeMonitor struct {
	ID                int    `json:"id"`
	FriendlyName      string `json:"friendly_name"`
	CustomUptimeRatio string `json:"custom_uptime_ratio"`
}

type uptimeResponse struct {
	Stat     string          `json:"stat"`
	Error    json.RawMessage `json:"error,omitempty"`
	Monitors []uptimeMonitor `json:"monitors"`
}

// FetchUptimeRobot fetches uptime ratio from UptimeRobot API.
// The free API returns uptime ratio for custom ranges.
// We convert the ratio into a synthetic Record for SLO calculation.
func FetchUptimeRobot(apiKey, monitorID string, days int) ([]slo.Record, error) {
	payload := fmt.Sprintf(
		"api_key=%s&monitors=%s&custom_uptime_ratios=%d&format=json",
		apiKey, monitorID, days,
	)

	req, err := http.NewRequest("POST", "https://api.uptimerobot.com/v2/getMonitors", strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch uptime data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("UptimeRobot API error (status %d): %s", resp.StatusCode, string(body))
	}

	var uResp uptimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&uResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if uResp.Stat != "ok" {
		return nil, fmt.Errorf("UptimeRobot API error: %s", string(uResp.Error))
	}

	if len(uResp.Monitors) == 0 {
		return nil, fmt.Errorf("monitor %s not found", monitorID)
	}

	mon := uResp.Monitors[0]
	ratio, err := strconv.ParseFloat(mon.CustomUptimeRatio, 64)
	if err != nil {
		return nil, fmt.Errorf("parse uptime ratio: %w", err)
	}

	// Convert uptime percentage to a synthetic record.
	// Assume 1 check per 5 minutes = 288 checks/day.
	checksPerDay := int64(288)
	totalChecks := checksPerDay * int64(days)
	successChecks := int64(float64(totalChecks) * ratio / 100)

	record := slo.Record{
		Timestamp: time.Now(),
		Total:     totalChecks,
		Success:   successChecks,
	}

	return []slo.Record{record}, nil
}
