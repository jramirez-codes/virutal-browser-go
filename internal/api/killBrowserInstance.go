package api

import (
	"encoding/json"
	"net/http"
	"virtual-browser/internal/types"
)

func KillBrowserInstance(instanceCloseMap *types.ServerInstanceClose, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	url := r.URL.Query().Get("url")

	// Kill instance
	instanceCloseFunc, ok := (*instanceCloseMap).InstanceCloseMapFunc[url]
	if ok {
		// Remove from map
		(*instanceCloseMap).Mu.Lock()
		delete((*instanceCloseMap).InstanceCloseMapFunc, url)
		(*instanceCloseMap).Mu.Unlock()

		// Delete Instance and Return Response if Error
		if err := instanceCloseFunc(); err != nil {
			response := types.WsApiResponse{
				Success: false,
				Message: "Failed to kill browser instance",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		// Not Found
		response := types.WsApiResponse{
			Success: false,
			Message: "Instance not found",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.WsApiResponse{
		Success: true,
		Message: "Browser Instance URL retrieved killed successfully",
	}

	json.NewEncoder(w).Encode(response)
}
