package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/natori-hrj/sre-metrics-cli/internal/slo"
)

type vercelDeployment struct {
	UID       string `json:"uid"`
	State     string `json:"state"`
	CreatedAt int64  `json:"createdAt"`
}

type vercelResponse struct {
	Deployments []vercelDeployment `json:"deployments"`
}

// FetchVercelDeployments fetches deployment records from Vercel API.
// Returns one Record per day with total=deployments, success=successful deployments.
func FetchVercelDeployments(token, projectID, teamID string, since time.Time) ([]slo.Record, error) {
	params := url.Values{}
	params.Set("projectId", projectID)
	params.Set("since", fmt.Sprintf("%d", since.UnixMilli()))
	params.Set("limit", "100")
	if teamID != "" {
		params.Set("teamId", teamID)
	}

	reqURL := "https://api.vercel.com/v6/deployments?" + params.Encode()
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch deployments: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Vercel API error (status %d): %s", resp.StatusCode, string(body))
	}

	var vResp vercelResponse
	if err := json.NewDecoder(resp.Body).Decode(&vResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// Aggregate by day
	type dayStats struct {
		total   int64
		success int64
	}
	days := make(map[string]*dayStats)

	for _, d := range vResp.Deployments {
		t := time.UnixMilli(d.CreatedAt)
		key := t.Format("2006-01-02")
		if _, ok := days[key]; !ok {
			days[key] = &dayStats{}
		}
		days[key].total++
		if d.State == "READY" {
			days[key].success++
		}
	}

	var records []slo.Record
	for key, stats := range days {
		t, _ := time.Parse("2006-01-02", key)
		records = append(records, slo.Record{
			Timestamp: t,
			Total:     stats.total,
			Success:   stats.success,
		})
	}

	return records, nil
}
