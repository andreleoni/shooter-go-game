package collision

import (
	"math"
)

func Check(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func CheckCircle(x1, y1, r, x2, y2, w2, h2 float64) bool {
	dx := x1 - math.Max(x2, math.Min(x1, x2+w2))
	dy := y1 - math.Max(y2, math.Min(y1, y2+h2))

	return dx*dx+dy*dy < r*r
}

func CheckCircleCollision(
	x1, y1, r1,
	x2, y2, r2 float64) bool {
	// Calculate distance between centers
	dx := x2 - x1
	dy := y2 - y1
	distanceSquared := dx*dx + dy*dy

	radiusSum := r1 + r2

	return distanceSquared <= radiusSum*radiusSum
}

type CheckSpriteCollisionInput struct {
	X1, Y1, Width1, Height1 float64
	X2, Y2, Width2, Height2 float64
}

func CheckSpriteCollision(csci CheckSpriteCollisionInput) bool {
	// Calculate centers
	c1x := csci.X1 + csci.Width1/2
	c1y := csci.Y1 + csci.Height1/2
	c2x := csci.X2 + csci.Width2/2
	c2y := csci.Y2 + csci.Height2/2

	// Calculate semi-axes
	a1 := csci.Width1 / 2  // Semi-major axis of first ellipse
	b1 := csci.Height1 / 2 // Semi-minor axis of first ellipse
	a2 := csci.Width2 / 2  // Semi-major axis of second ellipse
	b2 := csci.Height2 / 2 // Semi-minor axis of second ellipse

	// Calculate distance between centers
	dx := c2x - c1x
	dy := c2y - c1y

	// Normalize coordinates by dividing by respective semi-axes
	nx := dx / (a1 + a2)
	ny := dy / (b1 + b2)

	// Check if normalized distance is less than or equal to 1
	return (nx*nx + ny*ny) <= 1.0
}
