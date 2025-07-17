# Changelog

All notable changes to the cclearly project will be documented in this file.

## [Unreleased] - 2024-07-17

### 🎉 Initial Development - cclearly cURL Inspector

**Project Overview**: Built a fast, local-first cURL inspector using Wails (Go + Svelte) with a focus on simplicity and performance. The application provides a graphical interface for the battle-tested `curl` command-line tool.

### ✨ Features Added

#### Core Application Structure
- **Three-pane layout**: History (24%), Request (38%), Response (38%) with responsive design
- **Component-based architecture**: Modular Svelte components for maintainability
  - `Header.svelte`: Application header with send request button
  - `HistoryPane.svelte`: Left sidebar showing request history
  - `RequestPane.svelte`: Middle section for configuring HTTP requests
  - `ResponsePane.svelte`: Right sidebar displaying response details

#### Backend Infrastructure (Go)
- **cURL Integration**: Direct integration with system `curl` binary for perfect compatibility
- **Request Parsing**: Intelligent cURL command parser supporting:
  - HTTP methods (GET, POST, PUT, DELETE, PATCH)
  - Headers with proper quote handling
  - Request bodies and data
  - URL parsing and validation
- **Response Processing**: Comprehensive response handling including:
  - Status code extraction with custom format (`HTTPSTATUS:200`)
  - Response headers parsing
  - Body content separation
  - Timing information capture
- **Database Layer**: SQLite integration for local request history
  - Automatic request/response storage
  - History retrieval with pagination
  - Full cURL command reconstruction

#### Frontend Features (Svelte + Tailwind CSS)
- **Modern UI**: Clean, Postman/Insomnia-inspired interface
- **Intelligent Paste Detection**: Automatic cURL command parsing on paste
- **Request Builder**: Comprehensive form for HTTP requests
  - Method selection dropdown
  - URL input with validation
  - Dynamic header management (add/remove)
  - Conditional body textarea for non-GET requests
- **Response Viewer**: Detailed response display
  - Status code with color coding
  - Response body with syntax highlighting
  - Response headers in organized format
- **History Management**: Persistent request history
  - Click-to-load previous requests
  - URL truncation with intelligent ellipsis
  - Metadata display (header count, payload size)
  - Timestamp information

### 🎨 UI/UX Improvements

#### Layout & Design
- **Responsive Grid**: Fixed percentage-based layout (24% history, 76% request/response)
- **Clean Headers**: Removed divider lines, integrated titles with content
- **Status Indicators**: Color-coded status codes throughout the application
  - 🟢 Green (200-299): Success responses
  - 🟡 Yellow (400-499): Client errors
  - 🔴 Red (500+): Server errors
  - ⚪ Gray: Unknown/other status codes

#### History Pane Enhancements
- **Smart URL Truncation**: 
  - Left-truncation for simple URLs to show endpoints
  - Middle truncation for complex URLs (domain + ... + endpoint)
  - Responsive to window size
- **Metadata Display**:
  - 📄 Header count with emoji
  - 📀 Payload size with diskette emoji
  - Formatted byte sizes (B, KB, MB, GB)
- **Visual Improvements**:
  - Larger, more readable text
  - Better spacing and padding
  - Hover states and selection indicators

#### Request Pane Features
- **Form Validation**: URL required before sending
- **Dynamic Headers**: Add/remove headers with key-value pairs
- **Method-Aware UI**: Body field only shows for non-GET requests
- **Placeholder Text**: Helpful examples and guidance

#### Response Pane Enhancements
- **Status in Header**: Status code moved to top-right of response pane
- **Color Coding**: Consistent status color scheme
- **Loading States**: Spinner during request execution
- **Error Handling**: Graceful error display with helpful tips

### 🔧 Technical Implementation

#### Backend Architecture
- **Modular Design**: Separated concerns into focused files
  - `models.go`: Data structures and types
  - `curl_executor.go`: cURL command execution and parsing
  - `database.go`: SQLite operations
  - `app.go`: Main application logic and Wails integration
- **Error Handling**: Comprehensive error management with graceful degradation
- **Data Persistence**: Local SQLite database in `~/.cclearly/cclearly.db`

#### Frontend Architecture
- **State Management**: Svelte reactive statements for clean data flow
- **Component Communication**: Props and event handlers for inter-component communication
- **Type Safety**: TypeScript definitions for Wails bindings
- **Styling**: Tailwind CSS for utility-first styling

#### Development Setup
- **Wails v2**: Modern desktop app framework
- **Svelte 3**: Reactive frontend framework
- **Tailwind CSS**: Utility-first CSS framework
- **SQLite**: Local database for persistence
- **Go Modules**: Modern Go dependency management

### 🐛 Bug Fixes & Improvements

#### cURL Integration
- **Status Code Parsing**: Fixed status code extraction from cURL output
- **Header Parsing**: Improved quote handling for headers with spaces
- **Command Tokenization**: Enhanced argument parsing with proper quote support
- **Response Separation**: Better separation of headers and body content

#### Frontend Issues
- **TypeScript Errors**: Fixed event handler type issues
- **Layout Responsiveness**: Improved pane sizing and overflow handling
- **Data Conversion**: Proper conversion between backend and frontend data formats
- **History Rehydration**: Fixed header population when loading from history

### 📝 Development Notes

#### Key Design Decisions
- **Local-First**: All data stored locally, no cloud dependencies
- **cURL Trust**: Leverages existing, battle-tested `curl` binary
- **Simplicity**: Focus on core functionality over feature bloat
- **Performance**: Lightweight Go backend with efficient SQLite storage

#### Future Considerations
- **Syntax Highlighting**: JSON/XML response formatting
- **Request Collections**: Grouping related requests
- **Environment Variables**: Support for different environments
- **Export/Import**: Request history backup and sharing
- **Keyboard Shortcuts**: Power user productivity features

### 🚀 Getting Started

The application is now functional with core cURL inspection capabilities:
1. Paste cURL commands to auto-populate request form
2. Modify requests as needed
3. Send requests and view responses
4. Browse request history with metadata
5. Click history items to reload previous requests

All data is stored locally in SQLite, ensuring privacy and offline functionality.
