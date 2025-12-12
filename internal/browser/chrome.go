// Package browser
package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type ChromeInstance struct {
	cmd         *exec.Cmd
	userDataDir string
	port        int
	WsURL       string
}

// WaitForChrome waits for Chrome to be ready
func (c *ChromeInstance) WaitForChrome(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for Chrome to start")
		case <-ticker.C:
			// Check if process is still running
			if c.cmd.ProcessState != nil && c.cmd.ProcessState.Exited() {
				return fmt.Errorf("chrome process exited unexpectedly")
			}

			// Try to connect to CDP endpoint
			resp, err := exec.Command("curl", "-s", fmt.Sprintf("http://localhost:%d/json/version", c.port)).Output()
			if err == nil && len(resp) > 0 {
				return nil
			}
		}
	}
}

// GetWebSocketURL returns the WebSocket debugger URL
func (c *ChromeInstance) GetWebSocketURL() (string, error) {
	if c.WsURL != "" {
		return c.WsURL, nil
	}

	output, err := exec.Command("curl", "-s", fmt.Sprintf("http://localhost:%d/json/version", c.port)).Output()
	if err != nil {
		return "", fmt.Errorf("failed to get WebSocket URL: %w", err)
	}

	// Parse JSON response to extract webSocketDebuggerUrl
	var versionInfo struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}

	if err := json.Unmarshal(output, &versionInfo); err != nil {
		return "", fmt.Errorf("failed to parse version info: %w", err)
	}

	if versionInfo.WebSocketDebuggerURL == "" {
		return "", fmt.Errorf("webSocketDebuggerUrl not found in response")
	}

	c.WsURL = versionInfo.WebSocketDebuggerURL
	return c.WsURL, nil
}

// Close terminates the Chrome instance and cleans up
func (c *ChromeInstance) Close() error {
	if c.cmd != nil && c.cmd.Process != nil {
		log.Printf("Shutting down Chrome (PID: %d)...\n", c.cmd.Process.Pid)

		// Try graceful shutdown first
		c.cmd.Process.Signal(os.Kill)

		// Wait a bit for graceful shutdown
		done := make(chan error, 1)
		go func() {
			done <- c.cmd.Wait()
		}()

		// WARN: Might Delete
		// select {
		// case <-time.After(5 * time.Second):
		// 	// Force kill if graceful shutdown fails
		// 	fmt.Printf("Force killing Chrome...")
		// 	c.cmd.Process.Kill()
		// 	<-done
		// case <-done:
		// 	// Graceful shutdown succeeded
		// }
	} else {
		log.Println("PID NOT FOUND")
	}

	// Clean up temp directory
	if c.userDataDir != "" {
		os.RemoveAll(c.userDataDir)
		log.Printf("Cleaned up user data directory\n")
	}

	return nil
}
