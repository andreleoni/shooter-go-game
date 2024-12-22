package entities

type Projectile struct {
	Active bool
	Power  float64

	X, Y       float64
	Speed      float64
	DirectionX float64
	DirectionY float64

	TargetX float64
	TargetY float64

	Width  float64
	Height float64
}
