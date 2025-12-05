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

// Get Browser Instance URL
func GetBrowserInstanceUrl(ch chan string, mu *sync.RWMutex, w http.ResponseWriter, r *http.Request) {
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
	mu.Lock()
	WsUrl := <-ch
	mu.Unlock()

	response := Response{
		Success: true,
		Message: "Browser Instance URL retrieved successfully",
		WsUrl:   WsUrl,
	}

	json.NewEncoder(w).Encode(response)
}

func KillBrowserInstance(mu *sync.RWMutex, instanceCloseMap map[string]func() error, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	url := r.URL.Query().Get("url")

	// Kill instance
	instanceCloseFunc, ok := instanceCloseMap[url]
	if ok {
		// Found
		if err := instanceCloseFunc(); err != nil {
			mu.Lock()
			delete(instanceCloseMap, url)
			mu.Unlock()
			response := Response{
				Success: false,
				Message: "Failed to kill browser instance",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		// Not Found
		response := Response{
			Success: false,
			Message: "Instance not found",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Browser Instance URL retrieved killed successfully",
	}

	json.NewEncoder(w).Encode(response)
}
