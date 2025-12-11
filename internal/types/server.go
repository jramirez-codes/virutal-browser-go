package types

import (
	"sync"
)

type ServerInstanceClose struct {
	InstanceCloseMapFunc map[string]func() error
	Mu                   sync.RWMutex
}

type IsCreatingInstance struct {
	Status bool
	Mu     sync.RWMutex
}
