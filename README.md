# ETag Caching Demo

This project demonstrates how HTTP ETag caching works between a Go server and a browser client. The server implements ETag-based conditional requests, and the browser automatically handles cache validation using the If-None-Match header.

## Project Structure

```
etag/
├── main.go      # Go HTTP server with ETag implementation
├── index.html   # Browser client that polls the server
├── test.json    # Sample data file that can be modified
└── go.mod       # Go module definition
```

## How It Works

### Server Side

The Go server implements an HTTP endpoint that:

1. Reads data from test.json
2. Generates a SHA256 checksum (ETag) for the content
3. Returns the ETag header with the response
4. Handles conditional requests using the If-None-Match header
5. Returns 304 Not Modified when the cached version is still valid
6. Returns 200 OK with full content when data has changed

Key headers used:
- `Cache-Control: public, max-age=0, must-revalidate` - Forces browser to revalidate with server
- `ETag` - Unique identifier for the current content version
- `If-None-Match` - Browser sends this to check if cached version is still valid

### Client Side

The HTML page polls the server every 5 seconds and demonstrates:

- First request: Server returns 200 OK with ETag
- Subsequent requests: Browser automatically adds If-None-Match header
- If content unchanged: Server returns 304 Not Modified (no body transferred)
- If content changed: Server returns 200 OK with new data and new ETag

## Running the Project

### Prerequisites

- Go 1.25 or later
- Air (optional, for live reload)

### Starting the Server

```bash
# Using Go directly
go run main.go

# Using Air for live reload
air
```

The server starts on http://localhost:8080

### Testing with the Browser

1. Open http://localhost:8080/ in a browser
2. Open browser DevTools (F12) and go to Network tab
3. Ensure "Disable cache" is unchecked
4. Watch the requests in the Network tab

Expected behavior:
- First request: 200 OK, full response body transferred
- Subsequent requests: 304 Not Modified, no body transferred
- After modifying test.json: 200 OK, new ETag and response body

### Verifying Cache Behavior

```bash
# First request - full response
curl -v http://localhost:8080/

# Subsequent request with ETag - should return 304
curl -v -H "If-None-Match: <etag-from-previous>" http://localhost:8080/
```

## Modifying the Data

Edit test.json to see the cache invalidation:

```json
{
    "data": "your new content here"
}
```

Save the file and the next poll will detect the change and return new data with a new ETag.

## Key Concepts

### Why Use ETags

- Reduces bandwidth by allowing 304 responses when content has not changed
- Ensures clients always have fresh content through conditional requests
- Works alongside Cache-Control for comprehensive caching strategy

### Cache Flow

```
Initial Request:
  Client -> GET / -> Server (200 OK + ETag) -> Client (caches response)

Subsequent Request (same content):
  Client -> GET / + If-None-Match: <etag> -> Server
  Server -> 304 Not Modified -> Client (uses cached response)

After Content Change:
  Client -> GET / + If-None-Match: <old-etag> -> Server
  Server -> 200 OK + New-ETag + New Content -> Client (updates cache)
```

## Dependencies

No external Go dependencies required. The project uses only the Go standard library.

## License

This project is for educational purposes to demonstrate HTTP caching mechanisms.
