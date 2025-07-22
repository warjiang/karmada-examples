package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

type Metric struct {
	ID        int64
	Name      string
	Type      string
	Pod       string
	Labels    []Label
	Value     float64
	Timestamp int64
}

type Label struct {
	Name  string
	Value string
}

func New() (*DB, error) {
	db, err := sql.Open("sqlite3", "./metrics.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &DB{db}, nil
}

func createTables(db *sql.DB) error {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	schemaPath := filepath.Join(basepath, "..", "..", "db", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	_, err = db.Exec(string(schema))
	return err
}

func (db *DB) InsertMetric(metric *Metric) error {
	labelsJSON, err := json.Marshal(metric.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	_, err = db.Exec(
		"INSERT INTO metrics (name, type, pod, labels, value, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
		metric.Name,
		metric.Type,
		metric.Pod,
		string(labelsJSON),
		metric.Value,
		metric.Timestamp,
	)
	return err
}

func (db *DB) GetMetrics(name string, labels map[string]string, pod string) ([]*Metric, error) {
	query := "SELECT id, name, type, pod, labels, value, timestamp FROM metrics WHERE name = ?"
	args := []interface{}{name}

	if pod != "" {
		query += " AND pod = ?"
		args = append(args, pod)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var metrics []*Metric
	for rows.Next() {
		metric := &Metric{}
		var labelsJSON string
		if err := rows.Scan(&metric.ID, &metric.Name, &metric.Type, &metric.Pod, &labelsJSON, &metric.Value, &metric.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan metric: %w", err)
		}
		if err := json.Unmarshal([]byte(labelsJSON), &metric.Labels); err != nil {
			return nil, fmt.Errorf("failed to unmarshal labels: %w", err)
		}

		if len(labels) > 0 {
			match := true
			for k, v := range labels {
				found := false
				for _, l := range metric.Labels {
					if l.Name == k && l.Value == v {
						found = true
						break
					}
				}
				if !found {
					match = false
					break
				}
			}
			if match {
				metrics = append(metrics, metric)
			}
		} else {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

func (db *DB) GetMetricNames() ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT name FROM metrics")
	if err != nil {
		return nil, fmt.Errorf("failed to query metric names: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan metric name: %w", err)
		}
		names = append(names, name)
	}

	return names, nil
}
