package main

import (
	"log"
	"net/http"
	"sync"

	"go/browser/lib"
)

var (
	instanceCloseMap = make(map[string]func() error)
	wsUrlChannels    = make(chan string, 1000)
	mu               sync.RWMutex
)

func CreateInstance() (*lib.ChromeInstance, error) {
	// Get Available Port
	startPort, err := lib.GetPort()
	if err != nil {
		return nil, err
	}
	headless := true

	// Get Chrome Instance
	instance, err := lib.LaunchChrome(startPort, headless)
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

func main() {
	// Create Inital Instance N+1
	_, err := CreateInstance()
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}

	// API Server - Register routes
	apiPort := ":8080"
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		lib.GetBrowserInstanceUrl(wsUrlChannels, &mu, w, r)

		// Create New Instance N+1 (Preload)
		_, err := CreateInstance()
		if err != nil {
			log.Fatalf("Failed to create instance: %v", err)
		}
	})

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		lib.KillBrowserInstance(&mu, instanceCloseMap, w, r)
	})

	log.Printf("Server starting on http://localhost%s", apiPort)
	log.Println("\nAvailable endpoints:")
	log.Println(" GET - Get WebSocket URL")

	if err := http.ListenAndServe(apiPort, nil); err != nil {
		log.Fatal(err)
	}

	// Keep running until interrupted
	select {}
}
