package main

import (
	"context"
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
		return
	}
	a.database = database
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
		return []HistoryItem{}, nil
	}

	return a.database.GetHistory(limit)
}

// GetRequestByID retrieves a specific request by ID
func (a *App) GetRequestByID(id int64) (*HistoryItem, error) {
	if a.database == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	return a.database.GetRequestByID(id)
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
