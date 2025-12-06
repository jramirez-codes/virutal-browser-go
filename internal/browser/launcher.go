package browser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// GetChromePath returns the Chrome executable path based on OS
func GetChromePath() string {
	var paths []string

	switch runtime.GOOS {
	case "darwin": // macOS
		paths = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}
	case "linux":
		paths = []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
		}
	case "windows":
		paths = []string{
			"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
			"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Google\\Chrome\\Application\\chrome.exe"),
		}
	}

	// Check each path
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Try to find in PATH
	if path, err := exec.LookPath("google-chrome"); err == nil {
		return path
	}
	if path, err := exec.LookPath("chromium"); err == nil {
		return path
	}
	if path, err := exec.LookPath("chromium-browser"); err == nil {
		return path
	}

	return ""
}

// LaunchChrome starts a Chrome instance with CDP enabled
func LaunchChrome(port int, headless bool) (*ChromeInstance, error) {
	chromePath := GetChromePath()
	if chromePath == "" {
		return nil, fmt.Errorf("Chrome/Chromium not found. Please install Chrome or set CHROME_PATH environment variable")
	}

	log.Printf("Found Chrome at: %s", chromePath)

	// Create temporary user data directory
	userDataDir, err := os.MkdirTemp("", "chrome-cdp-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

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
	}

	if headless {
		args = append(args, "--headless=new")
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
	if err := instance.waitForChrome(30 * time.Second); err != nil {
		instance.Close()
		return nil, err
	}

	log.Printf("✓ Chrome started successfully")
	log.Printf("CDP endpoint: http://localhost:%d", port)

	return instance, nil
}
