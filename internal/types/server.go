package types

import (
	"sync"
	"virtual-browser/internal/browser"
)

type InstancePoolUsed struct {
	InstanceMap map[string]*browser.ChromeInstance
	Mu          sync.RWMutex
}

type IsCreatingInstance struct {
	Status bool
	Mu     sync.RWMutex
}
