package api

import (
	"encoding/json"
	"net/http"
	"virtual-browser/internal/browser"
	"virtual-browser/internal/types"
)

func GetBrowserInstanceUrl(ch chan *browser.ChromeInstance, usedPool *types.InstancePoolUsed, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get Chrome Instance
	currInstance := <-ch

	// Add to used pool
	usedPool.Mu.Lock()
	usedPool.InstanceMap[currInstance.WsURL] = currInstance
	usedPool.Mu.Unlock()

	response := types.WsApiResponse{
		Success: true,
		Message: "Browser Instance URL retrieved successfully",
		WsUrl:   currInstance.WsURL,
	}

	json.NewEncoder(w).Encode(response)
}
