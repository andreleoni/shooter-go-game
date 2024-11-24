// internal/core/kernel.go
package core

import (
	"game/internal/core/eventbus"
	"sync"
	"time"
)

const (
	FixedTimeStep = 1.0 / 60.0
	MaxSteps      = 5
)

type GameKernel struct {
	EventBus    *eventbus.EventBus
	TimeScale   float64
	DeltaTime   float64
	accumulator float64
	lastUpdate  time.Time
	mu          sync.Mutex
}

func NewGameKernel() *GameKernel {
	return &GameKernel{
		EventBus:   eventbus.NewEventBus(),
		TimeScale:  1.0,
		lastUpdate: time.Now(),
	}
}

func (k *GameKernel) Update(pm *PluginManager) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	currentTime := time.Now()
	frameTime := currentTime.Sub(k.lastUpdate).Seconds()
	k.lastUpdate = currentTime

	if frameTime > 0.25 {
		frameTime = 0.25
	}

	k.accumulator += frameTime

	steps := 0
	for k.accumulator >= FixedTimeStep && steps < MaxSteps {
		k.DeltaTime = FixedTimeStep

		if err := pm.UpdateAll(); err != nil {
			return err
		}

		k.accumulator -= FixedTimeStep
		steps++
	}

	// Store leftover time
	k.DeltaTime = k.accumulator

	return nil
}
