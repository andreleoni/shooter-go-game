package entities

// Area of effect
type AOF interface {
	Target() string
	GetCollision(x, y float64) bool
}

type Square struct {
	X, Y          float64
	Width, Height float64
}

type Radius struct {
	X, Y   float64
	Radius float64
}
