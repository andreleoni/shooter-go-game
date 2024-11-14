package bullet

import "math"

type Bullet struct {
	X, Y   float64
	Width  float64
	Height float64
	Speed  float64
	Active bool

	TargetX float64
	TargetY float64
}

func (b *Bullet) GetBounds() (float64, float64, float64, float64) {
	return b.X, b.Y, b.Width, b.Height
}

// Add method to update target position
func (b *Bullet) UpdateTarget(targetX, targetY float64) {
	b.TargetX = targetX
	b.TargetY = targetY
}

func (b *Bullet) MoveTowardsTarget(deltaTime float64) {
	dx := b.TargetX - b.X
	dy := b.TargetY - b.Y

	// Calculate distance
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		// Normalize direction
		dx /= distance
		dy /= distance

		// Update position
		b.X += dx * b.Speed * deltaTime
		b.Y += dy * b.Speed * deltaTime
	}
}
