package api

import (
	"encoding/json"
	"net/http"
	"virtual-browser/internal/types"
)

func GetServerStats(serverStates *types.ServerStatsResponse, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverStates)
}
