Of course. This is the perfect kind of document to keep you grounded and focused during development. It serves as your project's constitution.

Here is the markdown file excerpt. You can copy this directly into a `README.md` or a `PROJECT.md` file in your new repository.

---

# Project: cclearly

A fast, local-first cURL inspector. No accounts, no cloud, no bloat. Just a clean interface between you and `curl`.

## Guiding Principles (The Vibe)

Every decision should be measured against these principles.

*   **Local First, Always:** All data (history, collections) lives in a local SQLite file. The app is 100% functional offline.
*   **Performance is a Feature:** The app must launch instantly and feel snappy. This is achieved by using a lightweight Go binary and a native webview (Wails).
*   **Simplicity is King:** We aggressively favor simplicity over a cluttered feature set. The core loop is sacred: **Paste -> Inspect -> Send -> View**.
*   **Trust in `curl`:** We do not re-implement an HTTP client. We are a graphical user interface for the battle-tested `curl` command-line tool. This guarantees perfect compatibility and leverages an open-source powerhouse.

## Core Features

*   **Intelligent Paste:** The app listens for a paste event. If the clipboard contains text starting with `curl `, it automatically parses it and populates the request UI.
*   **Three-Pane Layout:**
    *   **Left Pane:** A searchable list of past requests (history), with the most recent at the top.
    *   **Middle Pane:** The request editor. Clean, editable fields for the URL, method, headers, query params, and body.
    *   **Right Pane:** The response viewer. A simple view with tabs for Body, Headers, and Timings. Includes syntax highlighting for common content types (JSON, HTML, etc.).
*   **Local History:** All requests are automatically saved to a local SQLite database for later reference and re-use.
*   **Request Editor:** Modify any part of the parsed request before sending it.
*   **Cross-Platform:** Works on macOS, Windows, and Linux thanks to Wails.

## Technology Stack

| Layer     | Technology                                      | Purpose                                                                    |
| :-------- | :---------------------------------------------- | :------------------------------------------------------------------------- |
| **Core**  | **Go (Golang)**                                 | The application backend. Handles logic, file system access, and commands.  |
| **Bridge**  | **Wails v2**                                    | Binds the Go backend to a web frontend, creating a native app.             |
| **Frontend**| **Svelte**                                      | For building a reactive, simple, and highly performant user interface.     |
| **State**   | **Svelte Stores (`svelte/store`)**              | Simple, built-in state management for a clean and reactive data flow.      |
| **Styling** | **Tailwind CSS**                                | A utility-first CSS framework for rapidly building the custom UI.          |
| **Database**| **SQLite**                                      | The local, file-based database for storing request history.                |

## Software Architecture & Dependencies

The application is cleanly divided into a Go backend and a Svelte frontend, which communicate via the Wails JS/Go bridge.

```plaintext
+-----------------------------------------------------------------+
|                         Svelte Frontend                         |
|  [History Pane]     [Request Pane]       [Response Pane]        |
|  (Loads from store) (Populates store)    (Reads from store)     |
+-----------------|----------------|------------------------------+
                  ^                |
                  | JS Calls       | Svelte Store (UI State)
                  v                v
+-----------------------------------------------------------------+
|                           Wails Bridge                          |
+-----------------|----------------|------------------------------+
                  ^                |
                  | Go Functions   | Go Calls
                  v                v
+-----------------------------------------------------------------+
|                            Go Backend                           |
|                                                                 |
|  [cURL Parser] <-> [cURL Executor] <-> [SQLite Handler]         |
|  (string -> struct)   (os/exec "curl")  (mattn/go-sqlite3)       |
+-----------------------------------------------------------------+
```

### Key Dependencies:

*   **Go Backend:**
    *   `os/exec`: Standard library package used to call the system's `curl` binary.
    *   `github.com/mattn/go-sqlite3`: The Go driver for our SQLite database.
    *   *Optional:* A third-party cURL parsing library, or a simple custom one using `strings` and `regexp`.

*   **Svelte Frontend:**
    *   `tailwindcss`, `postcss`, `autoprefixer`: For styling.
    *   *Optional:* `daisyUI` (Tailwind component library) for quick scaffolding of buttons, inputs, etc.
    *   *Optional:* `highlight.js` for syntax highlighting in the response pane.

## Data Model (SQLite)

A single table, `history`, is sufficient for the MVP.

**Table: `history`**

| Column         | Type      | Description                                          |
| :------------- | :-------- | :--------------------------------------------------- |
| `id`           | `INTEGER` | Primary Key, Auto-incrementing.                      |
| `method`       | `TEXT`    | The HTTP method (e.g., "GET", "POST").               |
| `url`          | `TEXT`    | The request URL, without query params.               |
| `full_command` | `TEXT`    | The complete, raw `curl` command for perfect recall. |
| `status_code`  | `INTEGER` | The HTTP status code of the response (e.g., 200, 404). |
| `created_at`   | `INTEGER` | Unix timestamp of when the request was made.         |

## Kickstart Commands

A quick reference for getting the project started.

```bash
# 1. Initialize a new Wails project with the Svelte template
wails init -n cclearly -t svelte

# 2. Navigate into the project directory
cd cclearly

# 3. Install frontend dev dependencies for Tailwind CSS
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p

# 4. Install the Go dependency for SQLite
go get github.com/mattn/go-sqlite3

# 5. Run the app in development mode
wails dev
```