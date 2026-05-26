package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type ServiceStatusRecord struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir db dir: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Initialize(ctx context.Context, db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS service_status (
			name TEXT PRIMARY KEY,
			status TEXT NOT NULL,
			message TEXT NOT NULL,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS audit_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			action TEXT NOT NULL,
			details TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}
	for _, stmt := range stmts {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func UpsertSetting(ctx context.Context, db *sql.DB, key, value string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO settings(key, value, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET
			value=excluded.value,
			updated_at=CURRENT_TIMESTAMP
	`, key, value)
	return err
}

func InsertSettingIfMissing(ctx context.Context, db *sql.DB, key, value string) error {
	_, err := db.ExecContext(ctx, `
		INSERT OR IGNORE INTO settings(key, value, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`, key, value)
	return err
}

func GetSetting(ctx context.Context, db *sql.DB, key string) (string, error) {
	var value string
	err := db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func UpsertServiceStatus(ctx context.Context, db *sql.DB, status ServiceStatusRecord) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO service_status(name, status, message, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(name) DO UPDATE SET
			status=excluded.status,
			message=excluded.message,
			updated_at=CURRENT_TIMESTAMP
	`, status.Name, status.Status, status.Message)
	return err
}

func ListServiceStatuses(ctx context.Context, db *sql.DB) ([]ServiceStatusRecord, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT name, status, message
		FROM service_status
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]ServiceStatusRecord, 0)
	for rows.Next() {
		var status ServiceStatusRecord
		if err := rows.Scan(&status.Name, &status.Status, &status.Message); err != nil {
			return nil, err
		}
		result = append(result, status)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
