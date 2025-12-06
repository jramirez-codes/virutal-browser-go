package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

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
