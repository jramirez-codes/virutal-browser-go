package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"virtual-browser/internal/api"
	"virtual-browser/internal/browser"
	"virtual-browser/internal/types"
	"virtual-browser/internal/util"
)

// Global Variables
var instanceCloseMap = types.ServerInstanceClose{
	InstanceCloseMapFunc: make(map[string]func() error),
	Mu:                   sync.RWMutex{},
}
var wsURLChannels = make(chan string, 500)
var serverStats = types.ServerStatsResponse{
	StartTime:                 0,
	CPUUsage:                  0.0,
	MemoryUsage:               0,
	LiveChromeInstanceCount:   0,
	ServedChromeInstanceCount: 0,
	Mu:                        sync.RWMutex{},
}

func CreateInstance() (*browser.ChromeInstance, error) {
	// Get Available Port
	startPort, err := util.GetPort()
	if err != nil {
		return nil, err
	}

	// Get Chrome Instance
	instance, err := browser.LaunchChrome(startPort)
	if err != nil {
		return nil, err
	}

	// Get WebSocket URL
	wsURL, err := instance.GetWebSocketURL()
	if err != nil {
		return nil, err
	}

	// Add WebSocket URL to map
	instanceCloseMap.Mu.Lock()
	wsURLChannels <- wsURL
	instanceCloseMap.InstanceCloseMapFunc[wsURL] = instance.Close
	instanceCloseMap.Mu.Unlock()

	return instance, nil
}

func StartAPIServer() {
	// API Server - Register routes
	apiPort := ":8080"

	// Get WebSocket URL
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		// Create New Instance N+1 (Preload)
		go func() {
			_, err := CreateInstance()
			if err != nil {
				log.Fatalf("Failed to create instance: %v", err)
			}
		}()

		api.GetBrowserInstanceUrl(wsURLChannels, w, r)
	})

	// Kill WebSocket URL
	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		api.KillBrowserInstance(&instanceCloseMap, w, r)
		serverStats.Mu.Lock()
		serverStats.ServedChromeInstanceCount++
		serverStats.Mu.Unlock()
	})

	// Get Server Stats
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		// Get Live Chrome Instance Count
		instanceCloseMap.Mu.RLock()
		serverStats.Mu.Lock()
		serverStats.LiveChromeInstanceCount = len(instanceCloseMap.InstanceCloseMapFunc)
		serverStats.Mu.Unlock()
		instanceCloseMap.Mu.RUnlock()

		// Return Status
		api.GetServerStats(&serverStats, w, r)
	})

	log.Printf("Server starting on http://localhost%s", apiPort)
	log.Println("\nAvailable endpoints:")
	log.Println(" GET - Get WebSocket URL")

	if err := http.ListenAndServe(apiPort, nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Add WaitGroup
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Record Start Time
	serverStats.StartTime = time.Now().Unix()

	// Create Inital Instance N+1
	go func() {
		_, err := CreateInstance()
		if err != nil {
			log.Fatalf("Failed to create instance: %v", err)
		}
	}()

	// Start API Server
	go StartAPIServer()

	wg.Wait()

	// Keep running until interrupted
	select {}
}
