package main

import (
	"log"
	"net/http"
	"sync"

	"virtual-browser/internal/api"
	"virtual-browser/internal/browser"
	"virtual-browser/internal/util"
)

var (
	instanceCloseMap = make(map[string]func() error)
	mu               sync.RWMutex
)
var wsUrlChannels = make(chan string, 500)

func CreateInstance() (*browser.ChromeInstance, error) {
	// Get Available Port
	startPort, err := util.GetPort()
	if err != nil {
		return nil, err
	}
	headless := true

	// Get Chrome Instance
	instance, err := browser.LaunchChrome(startPort, headless)
	if err != nil {
		return nil, err
	}

	// Get WebSocket URL
	wsUrl, err := instance.GetWebSocketURL()
	if err != nil {
		return nil, err
	}

	// Add WebSocket URL to map
	mu.Lock()
	wsUrlChannels <- wsUrl
	instanceCloseMap[wsUrl] = instance.Close
	mu.Unlock()

	return instance, nil
}

func StartApiServer() {
	// API Server - Register routes
	apiPort := ":8080"
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		api.GetBrowserInstanceUrl(wsUrlChannels, w, r)

		// Create New Instance N+1 (Preload)
		_, err := CreateInstance()
		if err != nil {
			log.Fatalf("Failed to create instance: %v", err)
		}
	})

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		api.KillBrowserInstance(&mu, instanceCloseMap, w, r)
	})

	log.Printf("Server starting on http://localhost%s", apiPort)
	log.Println("\nAvailable endpoints:")
	log.Println(" GET - Get WebSocket URL")

	if err := http.ListenAndServe(apiPort, nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Create Inital Instance N+1
	go func() {
		_, err := CreateInstance()
		if err != nil {
			log.Fatalf("Failed to create instance: %v", err)
		}
	}()

	// Start API Server
	go StartApiServer()

	// Keep running until interrupted
	select {}
}
