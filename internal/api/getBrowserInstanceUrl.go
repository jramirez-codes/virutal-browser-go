package api

import (
	"encoding/json"
	"net/http"
	"virtual-browser/internal/browser"
	"virtual-browser/internal/types"
)

func GetBrowserInstanceUrl(ch *browser.ChromeInstance, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := types.WsApiResponse{
		Success: true,
		Message: "Browser Instance URL retrieved successfully",
		WsUrl:   ch.WsURL,
	}

	json.NewEncoder(w).Encode(response)
}
