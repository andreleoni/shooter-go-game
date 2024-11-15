package core

import (
	"sync"
	"time"
)

const (
	FixedTimeStep = 1.0 / 60.0 // 60 FPS
	MaxSteps      = 5          // Prevent spiral of death
)

type GameKernel struct {
	PluginManager *PluginManager
	TimeScale     float64
	DeltaTime     float64
	accumulator   float64
	lastUpdate    time.Time
	mu            sync.Mutex
}

func NewGameKernel() *GameKernel {
	return &GameKernel{
		PluginManager: NewPluginManager(),
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
