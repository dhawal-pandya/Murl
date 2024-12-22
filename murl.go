package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	method := flag.String("X", "GET", "HTTP method (GET, POST, DELETE, PUT)")
	data := flag.String("d", "", "Data payload for POST or PUT requests")
	headers := flag.String("H", "", "Additional headers")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: murl <url>")
		return
	}
	rawURL := flag.Args()[0]

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port == "" {
		port = "80"
	}
	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	// request for direct TCP connection
	requestLine := fmt.Sprintf("%s %s HTTP/1.1\r\n", *method, path)
	hostHeader := fmt.Sprintf("Host: %s\r\n", host)
	defaultHeaders := "Accept: */*\r\nConnection: close\r\n"
	customHeaders := ""
	if *headers != "" {
		customHeaders = fmt.Sprintf("%s\r\n", *headers)
	}
	content := ""
	if *data != "" {
		content = fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(*data), *data)
	} else {
		content = "\r\n"
	}
	tcpRequest := requestLine + hostHeader + defaultHeaders + customHeaders + content

	req, err := http.NewRequest(*method, rawURL, bytes.NewBufferString(*data))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		os.Exit(1)
	}

	// headers
	if *headers != "" {
		parts := strings.SplitN(*headers, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		} else {
			fmt.Println("Error: Invalid header format. Use 'Key: Value'.")
			os.Exit(1)
		}
	}
	if *method == "POST" || *method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// make the request here
	resp, err := client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			fmt.Println("Error: Request timed out")
		} else if strings.Contains(err.Error(), "no such host") {
			fmt.Println("Error: Invalid URL or domain")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		os.Exit(1)
	}
	defer resp.Body.Close()

	if *verbose {
		fmt.Printf("> %s %s HTTP/1.1\n", *method, req.URL.Path)
		for k, v := range req.Header {
			fmt.Printf("> %s: %s\n", k, strings.Join(v, ", ")) // > for the headers we send
		}
		fmt.Println()
		fmt.Printf("< HTTP/%d.%d %s\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status)
		for k, v := range resp.Header {
			fmt.Printf("< %s: %s\n", k, strings.Join(v, ", ")) // < for the headers we get
		}
		fmt.Println()
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 { // hack for not having to handle all requests
		fmt.Printf("Error: HTTP %d\n", resp.StatusCode)
		fmt.Println(string(body))
		os.Exit(1)
	}

	fmt.Println(string(body))

	// for debugging
	fmt.Printf("Connecting to %s:%s via TCP...\n", host, port)
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	if *verbose {
		fmt.Println("> " + strings.ReplaceAll(tcpRequest, "\r\n", "\n> "))
	}

	_, err = conn.Write([]byte(tcpRequest))
	if err != nil {
		fmt.Println("Error sending TCP request:", err)
		return
	}

	scanner := bufio.NewScanner(conn)
	if *verbose {
		for scanner.Scan() {
			fmt.Println("< " + scanner.Text())
		}
	} else { // just prints the body, skips the headers
		bodyStarted := false
		for scanner.Scan() {
			line := scanner.Text()
			if bodyStarted {
				fmt.Println(line)
			}
			if line == "" {
				bodyStarted = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading TCP response:", err)
	}
}
