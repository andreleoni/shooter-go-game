package core

import (
	"sync"
	"time"
)

type GameKernel struct {
	PluginManager *PluginManager
	EventBus      *EventBus
	TimeScale     float64
	DeltaTime     float64
	lastUpdate    time.Time
	mu            sync.Mutex
}

func NewGameKernel() *GameKernel {
	return &GameKernel{
		PluginManager: NewPluginManager(),
		EventBus:      NewEventBus(),
		TimeScale:     1.0,
		lastUpdate:    time.Now(),
	}
}

func (k *GameKernel) Update() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	currentTime := time.Now()
	k.DeltaTime = currentTime.Sub(k.lastUpdate).Seconds() * k.TimeScale
	k.lastUpdate = currentTime

	return k.PluginManager.UpdateAll()
}
