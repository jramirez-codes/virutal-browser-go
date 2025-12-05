package lib

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Data models
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	WsUrl   string `json:"wsUrl"`
}

var (
	instanceClose = make(map[string]func() error)
	mu            sync.RWMutex
)

// Get Browser Instance URL
func GetBrowserInstanceUrl(ch chan string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(ch) == 0 {
		response := Response{
			Success: false,
			Message: "Failed to retrieve WebSocket URL",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get WebSocket URL
	mu.RLock()
	WsUrl := <-ch
	mu.RUnlock()

	response := Response{
		Success: true,
		Message: "Browser Instance URL retrieved successfully",
		WsUrl:   WsUrl,
	}

	json.NewEncoder(w).Encode(response)
}

func KillBrowserInstance(ch chan string, instanceCloseMap map[string]func() error, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	url := r.URL.Query().Get("url")

	// Kill instance
	mu.Lock()
	if instanceCloseFunc, ok := instanceCloseMap[url]; ok {
		if err := instanceCloseFunc(); err != nil {
			mu.Unlock()
			response := Response{
				Success: false,
				Message: "Failed to kill browser instance",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		delete(instanceCloseMap, url)
	}
	mu.Unlock()

	response := Response{
		Success: true,
		Message: "Browser Instance URL retrieved killed successfully",
	}

	json.NewEncoder(w).Encode(response)
}
