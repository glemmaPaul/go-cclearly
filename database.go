package main

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database handles SQLite operations
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	database := &Database{db: db}
	if err := database.init(); err != nil {
		return nil, err
	}

	return database, nil
}

// init initializes the database schema
func (d *Database) init() error {
	query := `
	CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		method TEXT NOT NULL,
		url TEXT NOT NULL,
		full_command TEXT NOT NULL,
		status_code INTEGER,
		response_body TEXT,
		response_headers TEXT,
		response_type TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := d.db.Exec(query)
	return err
}

// SaveRequest saves a request to the history
func (d *Database) SaveRequest(req RequestData, response *ResponseData, fullCommand string) error {
	// Convert headers map to JSON string
	headersJSON := "{}"
	if response.Headers != nil {
		if headersBytes, err := json.Marshal(response.Headers); err == nil {
			headersJSON = string(headersBytes)
		}
	}

	query := `
	INSERT INTO history (method, url, full_command, status_code, response_body, response_headers, response_type, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, req.Method, req.URL, fullCommand, response.StatusCode, response.FormattedBody, headersJSON, response.ResponseType, time.Now())
	return err
}

// GetHistory retrieves the request history
func (d *Database) GetHistory(limit int) ([]HistoryItem, error) {
	query := `
	SELECT id, method, url, full_command, status_code, response_body, response_headers, response_type, created_at
	FROM history
	ORDER BY created_at DESC
	LIMIT ?
	`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return make([]HistoryItem, 0), err
	}
	defer rows.Close()

	var history []HistoryItem
	for rows.Next() {
		var item HistoryItem
		err := rows.Scan(&item.ID, &item.Method, &item.URL, &item.FullCommand, &item.StatusCode, &item.ResponseBody, &item.ResponseHeaders, &item.ResponseType, &item.CreatedAt)
		if err != nil {
			return make([]HistoryItem, 0), err
		}
		history = append(history, item)
	}

	// Ensure we never return nil
	if history == nil {
		history = make([]HistoryItem, 0)
	}

	return history, nil
}

// GetRequestByID retrieves a specific request by ID
func (d *Database) GetRequestByID(id int64) (*HistoryItem, error) {
	query := `
	SELECT id, method, url, full_command, status_code, response_body, response_headers, response_type, created_at
	FROM history
	WHERE id = ?
	`

	var item HistoryItem
	err := d.db.QueryRow(query, id).Scan(&item.ID, &item.Method, &item.URL, &item.FullCommand, &item.StatusCode, &item.ResponseBody, &item.ResponseHeaders, &item.ResponseType, &item.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}
