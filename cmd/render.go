package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/natori-hrj/sre-metrics-cli/internal/slo"
)

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63"))

	labelStyle = lipgloss.NewStyle().
			Width(16).
			Foreground(lipgloss.Color("245"))

	goodStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42"))

	badStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))

	warnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))
)

func renderStatus(service string, status slo.Status) string {
	var lines []string

	title := titleStyle.Render("SLO Status")
	lines = append(lines, centerText(title, 34))
	lines = append(lines, "")

	// Service
	lines = append(lines, labelStyle.Render("Service:")+service)

	// SLI
	sliStr := fmt.Sprintf("%.2f%%", status.SLI)
	lines = append(lines, labelStyle.Render("SLI:")+sliStr)

	// Target
	targetStr := fmt.Sprintf("%.2f%%", status.Target)
	lines = append(lines, labelStyle.Render("SLO Target:")+targetStr)

	// Status
	var statusStr string
	if status.Met {
		statusStr = goodStyle.Render("MEETING SLO")
	} else {
		statusStr = badStyle.Render("BREACHING SLO")
	}
	lines = append(lines, labelStyle.Render("Status:")+statusStr)

	// Error Budget
	budgetStr := renderBudget(status.ErrorBudgetPercent)
	lines = append(lines, labelStyle.Render("Error Budget:")+budgetStr)

	// Requests
	lines = append(lines, "")
	reqStr := fmt.Sprintf("%d / %d", status.SuccessRequests, status.TotalRequests)
	lines = append(lines, labelStyle.Render("Requests:")+reqStr)

	content := strings.Join(lines, "\n")
	return borderStyle.Render(content)
}

func renderBudget(percent float64) string {
	str := fmt.Sprintf("%.1f%% remaining", percent)
	switch {
	case percent > 50:
		return goodStyle.Render(str)
	case percent > 20:
		return warnStyle.Render(str)
	default:
		return badStyle.Render(str)
	}
}

func renderRecordSuccess(service string, total, success int64) string {
	sli := float64(success) / float64(total) * 100
	content := fmt.Sprintf("[%s] Recorded: %d/%d requests (SLI: %.2f%%)", service, success, total, sli)
	return goodStyle.Render(content)
}

func centerText(s string, width int) string {
	pad := (width - lipgloss.Width(s)) / 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + s
}
