# murl

`Murl` (`My` + `URL`) is a minimalist HTTP client for the command line, inspired by the simplicity and power of `cURL` (`Client` + `URL`). It allows developers to make HTTP requests (GET, POST, PUT, DELETE) with custom headers, payloads, and verbose output options for debugging. Whether you're testing APIs, handling errors gracefully, or exploring web endpoints, Murl provides a straightforward and efficient way to interact with URLs directly from your terminal.

## Features
- Supports HTTP methods: `GET`, `POST`, `PUT`, `DELETE`.
- Allows custom headers with `-H`.
- Sends JSON payloads with `-d`.
- Verbose mode to debug requests and responses. `-v`
- Default connection to port 80 for HTTP.

---

## Installation

### Prerequisites
- [Go](https://golang.org/dl/) installed (version 1.19 or later).

### Build from Source
1. Clone or download the repository.
   ```bash
   git clone https://github.com/dhawal-pandya/Murl
   cd Murl
   ```
2. Build the executable:
   ```bash
   go build murl.go
   ```

3. Verify the build:
   ```bash
   ./murl -h
   ```

### Install Globally (Optional)
#### macOS/Linux
1. Move the compiled binary to a directory in your `PATH`:
   ```bash
   mv murl /usr/local/bin/
   ```
2. Verify installation:
   ```bash
   murl -h
   ```

#### Windows
1. Move the compiled binary to a directory in your `PATH` (e.g., `C:\Windows\System32`).
2. Verify installation:
   ```powershell
   murl -h
   ```

---

## Usage

### Basic Syntax
```bash
murl [options] <url>
```

### Options
- `-X <method>`: Specify the HTTP method (GET, POST, DELETE, PUT). Default is `GET`.
- `-d <data>`: Send data as the request body (used with POST and PUT).
- `-H <header>`: Add a custom header.
- `-v`: Enable verbose mode for debugging (shows detailed request and response headers).

### Examples

#### Simple GET Request
```bash
murl http://eu.httpbin.org/get
```

#### Verbose GET Request
```bash
murl -v http://eu.httpbin.org/get
```

#### POST Request with JSON Data
```bash
murl -X POST http://eu.httpbin.org/post -d '{"key": "value"}' -H "Content-Type: application/json"
```

#### DELETE Request
```bash
murl -X DELETE http://eu.httpbin.org/delete
```

#### PUT Request with JSON Data
```bash
murl -X PUT http://eu.httpbin.org/put -d '{"key": "value2"}' -H "Content-Type: application/json"
```

---

### Examples for Error Handling

#### GET Request with Non-Existing URL (404 Error)
```bash
murl http://httpbin.org/status/404
```

#### GET Request with Unauthorized Access (401 Error)
```bash
murl http://httpbin.org/status/401
```

#### POST Request with Invalid JSON Payload (400 Bad Request)
```bash
murl -X POST http://httpbin.org/post -d '{"key": "value"' -H "Content-Type: application/json"
```

- The server will return a 400 Bad Request status code due to invalid JSON formatting.

#### GET Request to an Invalid URL (Invalid Domain)
```bash
murl http://invalid.url
```

#### Server Error (500 Internal Server Error)
```bash
murl http://httpbin.org/status/500
```

#### POST Request with Missing Required Header
```bash
murl -X POST http://httpbin.org/post -d '{"key": "value"}'
```

- If a required header (like Content-Type: application/json) is missing, the server might return an error based on how it's configured to handle requests without proper headers.
- I added code that checks and set if the headers are empty for POST and PUT.

#### Timeout Error
```bash
murl http://httpbin.org/delay/10
```

- If the server takes too long to respond (e.g., due to delay endpoint), it may lead to a timeout error, depending on the client’s timeout settings.

---

## Notes
- Default port is 80 for HTTP. HTTPS is not currently supported.
- Ensure proper formatting for JSON payloads and headers.
- Test the program with [httpbin](http://httpbin.org/) for practice.

---

## Contributing
Feel free to fork the project, open issues, or submit pull requests.

---

## License
This project is licensed under the MIT License.

