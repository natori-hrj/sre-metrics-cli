package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/natori-hrj/sre-metrics-cli/internal/slo"
)

const (
	defaultDir = ".slo"
	timeFormat = time.RFC3339
)

func dataPath(service string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	dir := filepath.Join(home, defaultDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create data directory: %w", err)
	}
	filename := service + ".csv"
	return filepath.Join(dir, filename), nil
}

func Append(service string, record slo.Record) error {
	path, err := dataPath(service)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open data file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	return w.Write([]string{
		record.Timestamp.Format(timeFormat),
		strconv.FormatInt(record.Total, 10),
		strconv.FormatInt(record.Success, 10),
	})
}

func LoadAll(service string) ([]slo.Record, error) {
	path, err := dataPath(service)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open data file: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read CSV: %w", err)
	}

	records := make([]slo.Record, 0, len(rows))
	for i, row := range rows {
		if len(row) != 3 {
			return nil, fmt.Errorf("row %d: expected 3 fields, got %d", i+1, len(row))
		}
		ts, err := time.Parse(timeFormat, row[0])
		if err != nil {
			return nil, fmt.Errorf("row %d: parse timestamp: %w", i+1, err)
		}
		total, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: parse total: %w", i+1, err)
		}
		success, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: parse success: %w", i+1, err)
		}
		records = append(records, slo.Record{
			Timestamp: ts,
			Total:     total,
			Success:   success,
		})
	}
	return records, nil
}
