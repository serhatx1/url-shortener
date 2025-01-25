# URL Shortener
<img width="810" alt="image" src="https://github.com/user-attachments/assets/adc3aee7-78ef-41ad-b7bb-bfeafc7a61b2" />

This project is a simple URL shortener service built in Go. It allows users to shorten long URLs and redirect to the original URLs.

## Project Structure

```
url-shortener
├── cmd
│   └── main.go          # Main entry point for the application
├── internal
│   ├── handler
│   │   └── handler.go   # HTTP handler functions for API endpoints
│   ├── repository
│   │   └── repository.go # Data access layer for interacting with the database
│   ├── service
│       └── service.go   # Business logic for processing URL shortening requests
├── pkg
│   └── shortener
│       └── shortener.go  # Implementation of the URL shortener logic
├── go.mod                # Go module file defining dependencies
└── README.md             # Documentation for the project
```

## Setup Instructions

1. Clone the repository:
   ```
   git clone <repository-url>
   cd url-shortener
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run the application:
   ```
   go run cmd/main.go
   ```

## Usage

- To shorten a URL, send a POST request to `/shorten` with the original URL in the request body.
- To redirect to the original URL, send a GET request to `/{shortened-path}`.

## Contributing

Feel free to submit issues or pull requests for improvements or bug fixes.
