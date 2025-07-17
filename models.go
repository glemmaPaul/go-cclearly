package main

import "time"

// RequestData represents the structure of an HTTP request
type RequestData struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// ResponseData represents the structure of an HTTP response
type ResponseData struct {
	StatusCode    int               `json:"statusCode"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"body"`
	FormattedBody string            `json:"formattedBody,omitempty"` // Pre-formatted body (e.g., pretty JSON)
	ResponseType  string            `json:"responseType"`            // "json", "html", "raw"
	Timing        ResponseTiming    `json:"timing"`
	Error         string            `json:"error,omitempty"`
}

// ResponseTiming contains timing information for the request
type ResponseTiming struct {
	TotalTime    float64 `json:"totalTime"`    // Total time in seconds
	ConnectTime  float64 `json:"connectTime"`  // Time to establish connection
	TransferTime float64 `json:"transferTime"` // Time to transfer data
}

// HistoryItem represents a saved request in the database
type HistoryItem struct {
	ID              int64     `json:"id"`
	Method          string    `json:"method"`
	URL             string    `json:"url"`
	FullCommand     string    `json:"fullCommand"`
	StatusCode      int       `json:"statusCode"`
	ResponseBody    string    `json:"responseBody,omitempty"`
	ResponseHeaders string    `json:"responseHeaders,omitempty"`
	ResponseType    string    `json:"responseType,omitempty"`
	RequestSize     int64     `json:"requestSize"`  // Size of request payload in bytes
	ResponseSize    int64     `json:"responseSize"` // Size of response body in bytes
	CreatedAt       time.Time `json:"createdAt"`
}

// CurlCommand represents a parsed cURL command
type CurlCommand struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Raw     string            `json:"raw"` // Original cURL command
}
