package api

import (
	"encoding/json"
	"net/http"
	"virtual-browser/internal/types"
)

func GetBrowserInstanceUrl(ch chan string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(ch) == 0 {
		response := types.ApiResponse{
			Success: false,
			Message: "Failed to retrieve WebSocket URL",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get WebSocket URL
	WsUrl := <-ch

	response := types.ApiResponse{
		Success: true,
		Message: "Browser Instance URL retrieved successfully",
		WsUrl:   WsUrl,
	}

	json.NewEncoder(w).Encode(response)
}
