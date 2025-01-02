package player

import "github.com/hajimehoshi/ebiten/v2"

func InputHandler(p *PlayerPlugin, newX, newY, currentspeed float64) (float64, float64) {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		newY -= currentspeed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		newY += currentspeed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		newX -= currentspeed * p.kernel.DeltaTime
		p.facingRight = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		newX += currentspeed * p.kernel.DeltaTime
		p.facingRight = true
	}

	return newX, newY
}
