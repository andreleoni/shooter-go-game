package player

import "github.com/hajimehoshi/ebiten/v2"

func InputHandler(p *PlayerPlugin, newX, newY float64) (float64, float64) {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		newY -= p.speed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		newY += p.speed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		newX -= p.speed * p.kernel.DeltaTime
		p.facingRight = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		newX += p.speed * p.kernel.DeltaTime
		p.facingRight = true
	}

	return newX, newY
}
