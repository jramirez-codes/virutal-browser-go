# Google Browser Go

A high-performance, scalable solution for managing headless Chrome instances using Go and Nginx. This project provides a robust infrastructure for browser automation, scraping, and testing workflows by exposing Chrome DevTools Protocol (CDP) via WebSocket proxying.

## 🚀 Features

-   **Automated Instance Management**: Seamlessly launch and terminate headless Chrome instances on demand.
-   **WebSocket Proxying**: Secure and efficient WebSocket proxying via Nginx, enabling stable connections to Chrome DevTools.
-   **REST API**: Simple HTTP endpoints to request new browser sessions and manage lifecycles.
-   **Concurrency**: Designed to handle multiple concurrent browser sessions with dynamic port allocation.
-   **Health Monitoring**: Built-in health check endpoints for infrastructure reliability.

## 🏗️ Architecture

The system consists of three main components:

1.  **Nginx (Reverse Proxy)**: Acts as the entry point, handling incoming HTTP and WebSocket connections. It routes API requests to the Go server and proxies WebSocket connections directly to the specific Chrome instance.
2.  **Go Server (Manager)**: Orchestrates Chrome processes. It assigns unique ports, launches headless instances, and maintains a mapping of active sessions.
3.  **Chrome Instances**: Headless Chromium/Chrome processes running with remote debugging enabled, ready to accept CDP commands.

## 🛠️ Prerequisites

-   **Go**: Version 1.22 or higher.
-   **Nginx**: For reverse proxying and WebSocket support.
-   **Google Chrome / Chromium**: The browser executable must be installed and accessible in the system path.

## 📦 Installation

1.  **Clone the Repository**
    ```bash
    git clone https://github.com/yourusername/google-browser-go.git
    cd google-browser-go
    ```

2.  **Install Go Dependencies**
    ```bash
    go mod download
    ```

3.  **Configure Nginx**
    Copy the provided `nginx.conf` to your Nginx configuration directory or include it in your main `nginx.conf`.
    ```bash
    sudo cp nginx.conf /etc/nginx/sites-available/google-browser-go
    sudo ln -s /etc/nginx/sites-available/google-browser-go /etc/nginx/sites-enabled/
    sudo nginx -t
    sudo systemctl reload nginx
    ```

## 🚀 Usage

1.  **Start the Go Server**
    ```bash
    go run main.go
    ```
    The server will start on port `8080` (by default).

2.  **Request a Browser Instance**
    Make a GET request to the `/get` endpoint to launch a new Chrome instance and receive its WebSocket URL.
    ```bash
    curl http://localhost/get
    ```
    **Response:**
    ```json
    {
      "success": true,
      "message": "Browser Instance URL retrieved successfully",
      "wsUrl": "ws://localhost/devtools/browser/3500/uuid-string"
    }
    ```

3.  **Connect via WebSocket**
    Use the returned `wsUrl` to connect your automation tool (e.g., Puppeteer, Playwright, or direct CDP client) to the browser instance.

4.  **Terminate an Instance**
    When finished, kill the specific instance to free up resources.
    ```bash
    curl "http://localhost/kill?url=ws://localhost/devtools/browser/3500/uuid-string"
    ```

## ⚙️ Configuration

-   **`nginx.conf`**: Defines the proxy rules. Ensure the `proxy_pass` upstream matches your Go server's address.
-   **`main.go`**:
    -   `apiPort`: Port for the Go API server (default: `:8080`).
    -   `wsUrlChannels`: Buffered channel size for handling concurrent requests.
-   **`lib/GetPort.go`**: Defines the port range (default: `3000-3999`) for Chrome instances.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
