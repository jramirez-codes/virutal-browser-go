package api

// Data models
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	WsUrl   string `json:"wsUrl"`
}
