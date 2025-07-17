package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// CurlExecutor handles cURL command execution
type CurlExecutor struct{}

// NewCurlExecutor creates a new cURL executor
func NewCurlExecutor() *CurlExecutor {
	return &CurlExecutor{}
}

// ExecuteRequest executes an HTTP request using cURL
func (ce *CurlExecutor) ExecuteRequest(req RequestData) (*ResponseData, error) {
	start := time.Now()

	// Build cURL command
	args := ce.buildCurlArgs(req)

	// Execute cURL command
	cmd := exec.Command("curl", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	totalTime := time.Since(start).Seconds()

	// Parse response
	response := &ResponseData{
		StatusCode: 0,
		Headers:    make(map[string]string),
		Body:       stdout.String(),
		Timing: ResponseTiming{
			TotalTime: totalTime,
		},
	}

	if err != nil {
		response.Error = stderr.String()
		if response.Error == "" {
			response.Error = err.Error()
		}
		// Even if there's an error, try to extract status code
		ce.parseCurlOutput(stderr.String(), response)
		return response, nil // Return response with error, don't fail
	}

	// Parse the output to extract status code and headers
	output := stdout.String()
	ce.parseCurlOutput(output, response)

	return response, nil
}

// ParseCurlCommand parses a cURL command string into RequestData
func (ce *CurlExecutor) ParseCurlCommand(curlCmd string) (*RequestData, error) {
	// Remove "curl" prefix and trim whitespace
	curlCmd = strings.TrimSpace(strings.TrimPrefix(curlCmd, "curl"))

	req := &RequestData{
		Method:  "GET", // Default method
		Headers: make(map[string]string),
	}

	// Parse arguments with proper quote handling
	args := ce.tokenizeCurlCommand(curlCmd)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case strings.HasPrefix(arg, "-X") || strings.HasPrefix(arg, "--request"):
			if i+1 < len(args) {
				req.Method = args[i+1]
				i++
			}
		case strings.HasPrefix(arg, "-H") || strings.HasPrefix(arg, "--header"):
			if i+1 < len(args) {
				header := args[i+1]
				// Remove quotes if present
				header = strings.Trim(header, `"'`)
				if colonIndex := strings.Index(header, ":"); colonIndex > 0 {
					key := strings.TrimSpace(header[:colonIndex])
					value := strings.TrimSpace(header[colonIndex+1:])
					req.Headers[key] = value
				}
				i++
			}
		case strings.HasPrefix(arg, "-d") || strings.HasPrefix(arg, "--data"):
			if i+1 < len(args) {
				req.Body = strings.Trim(args[i+1], `"'`)
				if req.Method == "GET" {
					req.Method = "POST" // cURL defaults to POST with -d
				}
				i++
			}
		case strings.HasPrefix(arg, "--data-binary"):
			if i+1 < len(args) {
				req.Body = strings.Trim(args[i+1], `"'`)
				if req.Method == "GET" {
					req.Method = "POST"
				}
				i++
			}
		case !strings.HasPrefix(arg, "-") && req.URL == "":
			// First non-flag argument is the URL
			req.URL = strings.Trim(arg, `"'`)
		}
	}

	return req, nil
}

// buildCurlArgs builds the cURL command arguments from RequestData
func (ce *CurlExecutor) buildCurlArgs(req RequestData) []string {
	args := []string{
		"-s",                                // Silent mode
		"-w", "\nHTTPSTATUS:%{http_code}\n", // Write status code with prefix
		"-D", "-", // Dump headers to stdout
		"--max-time", "30", // 30 second timeout
		"--compressed",                  // Handle gzip/deflate compression automatically
		"-H", "Cache-Control: no-cache", // Prevent caching
		"-H", "Pragma: no-cache", // Prevent caching (HTTP/1.0)
	}

	// Add method
	if req.Method != "GET" {
		args = append(args, "-X", req.Method)
	}

	// Add headers
	for key, value := range req.Headers {
		args = append(args, "-H", fmt.Sprintf("%s: %s", key, value))
	}

	// Add body
	if req.Body != "" && req.Method != "GET" {
		args = append(args, "-d", req.Body)
	}

	// Add URL
	args = append(args, req.URL)

	return args
}

// tokenizeCurlCommand splits a cURL command into arguments with proper quote handling
func (ce *CurlExecutor) tokenizeCurlCommand(cmd string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(cmd); i++ {
		char := cmd[i]

		if !inQuotes && (char == '"' || char == '\'') {
			inQuotes = true
			quoteChar = char
			continue
		}

		if inQuotes && char == quoteChar {
			inQuotes = false
			quoteChar = 0
			continue
		}

		if !inQuotes && unicode.IsSpace(rune(char)) {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(char)
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

// parseCurlOutput extracts information from cURL's output
func (ce *CurlExecutor) parseCurlOutput(output string, response *ResponseData) {
	// Extract status code from our custom format
	statusRegex := regexp.MustCompile(`HTTPSTATUS:(\d+)`)
	if matches := statusRegex.FindStringSubmatch(output); len(matches) > 1 {
		if status, err := strconv.Atoi(matches[1]); err == nil {
			response.StatusCode = status
		}
	}

	// If we didn't find our custom format, try the standard HTTP response format
	if response.StatusCode == 0 {
		statusRegex := regexp.MustCompile(`HTTP/\d\.\d\s+(\d+)`)
		if matches := statusRegex.FindStringSubmatch(output); len(matches) > 1 {
			if status, err := strconv.Atoi(matches[1]); err == nil {
				response.StatusCode = status
			}
		}
	}

	// Extract response headers and body
	// Normalize line endings and split
	output = strings.ReplaceAll(output, "\r\n", "\n")
	output = strings.ReplaceAll(output, "\r", "\n")
	lines := strings.Split(output, "\n")
	var bodyLines []string
	inBody := false
	headerSectionEnded := false

	for _, line := range lines {
		if strings.HasPrefix(line, "HTTPSTATUS:") {
			continue // Skip our status line
		}

		// Skip HTTP status line (e.g., "HTTP/2 200")
		if strings.HasPrefix(strings.TrimSpace(line), "HTTP/") {
			continue
		}

		// Check if this line looks like a header (contains colon and doesn't start with { or [)
		isHeader := strings.Contains(line, ":") &&
			!strings.HasPrefix(strings.TrimSpace(line), "{") &&
			!strings.HasPrefix(strings.TrimSpace(line), "[")

		if !headerSectionEnded && isHeader {
			// This is a header line
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				response.Headers[key] = value
			}
		} else if !headerSectionEnded && !isHeader && strings.TrimSpace(line) != "" {
			// This is the start of the body (first non-header, non-empty line)
			headerSectionEnded = true
			inBody = true
			bodyLines = append(bodyLines, line)
		} else if inBody {
			// We're in the body section
			bodyLines = append(bodyLines, line)
		}
	}

	// Update body with only the body content
	if len(bodyLines) > 0 {
		response.Body = strings.Join(bodyLines, "\n")
	}

	// Detect response type and format if needed
	response.ResponseType = ce.detectResponseType(response.Headers, response.Body)

	// Format JSON responses
	if response.ResponseType == "json" {
		response.FormattedBody = ce.formatJSON(response.Body)
	} else {
		response.FormattedBody = response.Body
	}

	// Extract timing information if available
	timingRegex := regexp.MustCompile(`time_total:\s*([\d.]+)`)
	if matches := timingRegex.FindStringSubmatch(output); len(matches) > 1 {
		if timing, err := strconv.ParseFloat(matches[1], 64); err == nil {
			response.Timing.TotalTime = timing
		}
	}

	connectRegex := regexp.MustCompile(`time_connect:\s*([\d.]+)`)
	if matches := connectRegex.FindStringSubmatch(output); len(matches) > 1 {
		if timing, err := strconv.ParseFloat(matches[1], 64); err == nil {
			response.Timing.ConnectTime = timing
		}
	}
}

// detectResponseType determines the type of response based on content-type header and body content
func (ce *CurlExecutor) detectResponseType(headers map[string]string, body string) string {
	// Check content-type header first
	if contentType, exists := headers["Content-Type"]; exists {
		contentType = strings.ToLower(contentType)
		if strings.Contains(contentType, "application/json") {
			return "json"
		}
		if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/xhtml") {
			return "html"
		}
		if strings.Contains(contentType, "text/plain") {
			return "raw"
		}
		if strings.Contains(contentType, "text/xml") || strings.Contains(contentType, "application/xml") {
			return "raw"
		}
	}

	// Fallback: analyze body content
	body = strings.TrimSpace(body)
	if body == "" {
		return "raw"
	}

	// Check if it looks like JSON
	if (strings.HasPrefix(body, "{") && strings.HasSuffix(body, "}")) ||
		(strings.HasPrefix(body, "[") && strings.HasSuffix(body, "]")) {
		// Try to parse as JSON to confirm
		var js json.RawMessage
		if json.Unmarshal([]byte(body), &js) == nil {
			return "json"
		}
	}

	// Check if it looks like HTML
	if strings.HasPrefix(body, "<!DOCTYPE") ||
		strings.HasPrefix(body, "<html") ||
		strings.HasPrefix(body, "<HTML") {
		return "html"
	}

	// Default to raw
	return "raw"
}

// formatJSON pretty-prints JSON if the response is valid JSON
func (ce *CurlExecutor) formatJSON(body string) string {
	if body == "" {
		return ""
	}

	// Try to parse and format the JSON
	var js interface{}
	if err := json.Unmarshal([]byte(body), &js); err != nil {
		// If it's not valid JSON, return the original body
		return body
	}

	// Format the JSON with proper indentation
	formatted, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return body
	}

	return string(formatted)
}
