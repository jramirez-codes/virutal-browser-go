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

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Global Variables
var instanceCloseMap = types.ServerInstanceClose{
	InstanceCloseMapFunc: make(map[string]func() error),
	Mu:                   sync.RWMutex{},
}
var isCreatingInstance = types.IsCreatingInstance{
	Status: false,
	Mu:     sync.RWMutex{},
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
	// Set Creating Instance
	isCreatingInstance.Mu.Lock()
	isCreatingInstance.Status = true
	isCreatingInstance.Mu.Unlock()

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

	// Set Creating Instance
	isCreatingInstance.Mu.Lock()
	isCreatingInstance.Status = false
	isCreatingInstance.Mu.Unlock()

	return instance, nil
}

func StartAPIServer() {
	// API Server - Register routes
	apiPort := ":8080"

	// Get WebSocket URL
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			// Check the CPU Usage
			isCpuHigh := false
			serverStats.Mu.RLock()
			if serverStats.CPUUsage > 80 {
				isCpuHigh = true
			}
			serverStats.Mu.RUnlock()

			// Check if instance is already being created
			var localIsCreatingInstance bool
			isCreatingInstance.Mu.RLock()
			localIsCreatingInstance = isCreatingInstance.Status
			isCreatingInstance.Mu.RUnlock()

			// Wait for CPU to be low and instance to not be being created
			for isCpuHigh && !localIsCreatingInstance {
				time.Sleep(time.Second * 5)
				serverStats.Mu.RLock()
				if serverStats.CPUUsage < 80 {
					isCpuHigh = false
				}
				serverStats.Mu.RUnlock()

				// Check if instance is already being created
				isCreatingInstance.Mu.RLock()
				localIsCreatingInstance = isCreatingInstance.Status
				isCreatingInstance.Mu.RUnlock()
			}

			// Create New Instance N+1 (Preload)
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
		// Get Memory Usage
		memoryInfo, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("Failed to get memory info: %v", err)
		}
		// Get Live Chrome Instance Count
		serverStats.Mu.Lock()
		instanceCloseMap.Mu.RLock()
		serverStats.LiveChromeInstanceCount = len(instanceCloseMap.InstanceCloseMapFunc)
		instanceCloseMap.Mu.RUnlock()

		serverStats.MemoryUsage = int64(memoryInfo.UsedPercent)
		serverStats.Mu.Unlock()

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
	wg.Add(3)

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

	// Server Stats
	go func() {
		for {
			percentCPU, _ := cpu.Percent(time.Second, false)
			serverStats.Mu.Lock()
			serverStats.CPUUsage = percentCPU[0]
			serverStats.Mu.Unlock()
			time.Sleep(time.Second * 2)
		}
	}()

	wg.Wait()

	// Keep running until interrupted
	select {}
}
