package browser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"virtual-browser/internal/util/browserUtil"
)

// LaunchChrome starts a Chrome instance with CDP enabled
func LaunchChrome(port int) (*ChromeInstance, error) {
	chromePath := browserUtil.GetChromePath()
	if chromePath == "" {
		return nil, fmt.Errorf("Chrome/Chromium not found. Please install Chrome or set CHROME_PATH environment variable")
	}

	log.Printf("Found Chrome at: %s", chromePath)

	// Create temporary user data directory
	userDataDir, err := os.MkdirTemp("", "chrome-cdp-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Select random user agent
	userAgent := browserUtil.GenerateRandomUserAgent()
	log.Printf("Using User-Agent: %s", userAgent)

	// Chrome arguments
	args := []string{
		fmt.Sprintf("--remote-debugging-port=%d", port),
		"--no-sandbox",
		"--log-level=3",
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-background-networking",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-breakpad",
		"--disable-client-side-phishing-detection",
		"--disable-component-extensions-with-background-pages",
		"--disable-default-apps",
		"--disable-dev-shm-usage",
		"--disable-extensions",
		"--disable-features=TranslateUI,GCM",
		"--disable-hang-monitor",
		"--disable-ipc-flooding-protection",
		"--disable-popup-blocking",
		"--disable-prompt-on-repost",
		"--disable-renderer-backgrounding",
		"--disable-sync",
		"--force-color-profile=srgb",
		"--metrics-recording-only",
		"--no-first-run",
		"--enable-automation",
		"--password-store=basic",
		"--use-mock-keychain",
		fmt.Sprintf("--user-data-dir=%s", userDataDir),
		"--remote-debugging-address=0.0.0.0",
		"--headless=new",
		fmt.Sprintf("--user-agent=%s", userAgent),
	}

	// Add about:blank as initial page
	args = append(args, "about:blank")

	// Create command
	cmd := exec.Command(chromePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start Chrome
	if err := cmd.Start(); err != nil {
		os.RemoveAll(userDataDir)
		return nil, fmt.Errorf("failed to start Chrome: %w", err)
	}

	instance := &ChromeInstance{
		cmd:         cmd,
		userDataDir: userDataDir,
		port:        port,
	}

	// Wait for Chrome to be ready
	log.Printf("Waiting for Chrome to start on port %d...", port)
	if err := instance.WaitForChrome(30 * time.Second); err != nil {
		instance.Close()
		return nil, err
	}

	// Get WebSocket URL
	instance.WsURL, _ = instance.GetWebSocketURL()

	log.Printf("✓ Chrome started successfully")
	log.Printf("CDP endpoint: http://localhost:%d", port)

	return instance, nil
}
