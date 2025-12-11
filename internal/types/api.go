// types/api.go
package types

import "sync"

// Data models
type WsApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	WsUrl   string `json:"wsUrl"`
}

type ServerStatsResponse struct {
	StartTime                 int64   `json:"startTime"`
	CPUUsage                  float64 `json:"cpuUsage"`
	MemoryUsage               int64   `json:"memoryUsage"`
	LiveChromeInstanceCount   int     `json:"liveChromeInstanceCount"`
	ServedChromeInstanceCount int     `json:"servedChromeInstanceCount"`
	Mu                        sync.RWMutex
}

type ServerInstanceClose struct {
	InstanceCloseMapFunc map[string]func() error
	Mu                   sync.RWMutex
}
