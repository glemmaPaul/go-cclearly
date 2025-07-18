package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx          context.Context
	curlExecutor *CurlExecutor
	database     *Database
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize cURL executor
	a.curlExecutor = NewCurlExecutor()

	// Initialize database
	dbPath := filepath.Join(getDataDir(), "cclearly.db")
	database, err := NewDatabase(dbPath)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		// Don't return - continue without database for now
		// The app can still work for making requests, just without history
	} else {
		a.database = database
	}
}

// shutdown is called when the app shuts down
func (a *App) shutdown(ctx context.Context) {
	if a.database != nil {
		a.database.Close()
	}
}

// ExecuteRequest executes an HTTP request and returns the response
func (a *App) ExecuteRequest(req RequestData) (*ResponseData, error) {
	if a.curlExecutor == nil {
		return nil, fmt.Errorf("cURL executor not initialized")
	}

	response, err := a.curlExecutor.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}

	// Save to history if we have a database
	if a.database != nil {
		fullCommand := a.buildCurlCommand(req)
		a.database.SaveRequest(req, response, fullCommand)
	}

	return response, nil
}

// ParseCurlCommand parses a cURL command string into RequestData
func (a *App) ParseCurlCommand(curlCmd string) (*RequestData, error) {
	if a.curlExecutor == nil {
		return nil, fmt.Errorf("cURL executor not initialized")
	}

	return a.curlExecutor.ParseCurlCommand(curlCmd)
}

// GetHistory retrieves the request history
func (a *App) GetHistory(limit int) ([]HistoryItem, error) {
	if a.database == nil {
		return make([]HistoryItem, 0), nil
	}

	history, err := a.database.GetHistory(limit)
	if err != nil {
		return make([]HistoryItem, 0), err
	}

	// Ensure we always return a proper slice, not nil
	if history == nil {
		history = make([]HistoryItem, 0)
	}

	return history, nil
}

// TestDatabaseConnection tests if the database is working
func (a *App) TestDatabaseConnection() (bool, string) {
	if a.database == nil {
		return false, "Database not initialized"
	}

	// Try a simple query to test the connection
	_, err := a.database.GetHistory(1)
	if err != nil {
		return false, fmt.Sprintf("Database test failed: %v", err)
	}

	return true, "Database connection successful"
}

// GetRequestByID retrieves a specific request by ID
func (a *App) GetRequestByID(id int64) (*HistoryItem, error) {
	if a.database == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	return a.database.GetRequestByID(id)
}

// GetResponseByID retrieves the response data for a specific history item
func (a *App) GetResponseByID(id int64) (*ResponseData, error) {
	if a.database == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	item, err := a.database.GetRequestByID(id)
	if err != nil {
		return nil, err
	}

	// Parse headers from JSON string
	headers := make(map[string]string)
	if item.ResponseHeaders != "" {
		if err := json.Unmarshal([]byte(item.ResponseHeaders), &headers); err != nil {
			log.Printf("Failed to parse response headers: %v", err)
		}
	}

	response := &ResponseData{
		StatusCode:    item.StatusCode,
		Headers:       headers,
		Body:          item.ResponseBody,
		FormattedBody: item.ResponseBody,
		ResponseType:  item.ResponseType,
	}

	return response, nil
}

// buildCurlCommand builds a cURL command string from RequestData
func (a *App) buildCurlCommand(req RequestData) string {
	cmd := "curl"

	if req.Method != "GET" {
		cmd += fmt.Sprintf(" -X %s", req.Method)
	}

	for key, value := range req.Headers {
		cmd += fmt.Sprintf(" -H \"%s: %s\"", key, value)
	}

	if req.Body != "" && req.Method != "GET" {
		cmd += fmt.Sprintf(" -d \"%s\"", req.Body)
	}

	cmd += fmt.Sprintf(" \"%s\"", req.URL)

	return cmd
}

// getDataDir returns the application data directory
func getDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	dataDir := filepath.Join(homeDir, ".cclearly")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "."
	}

	return dataDir
}
