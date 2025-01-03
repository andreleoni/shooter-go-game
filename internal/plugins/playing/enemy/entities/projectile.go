package entities

type Projectile struct {
	X, Y          float64
	Width, Height float64
	Speed         float64
	DirectionX    float64
	DirectionY    float64
	Active        bool
	Power         float64
}
