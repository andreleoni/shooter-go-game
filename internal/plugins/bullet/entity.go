package bullet

type Bullet struct {
	X, Y   float64
	Width  float64
	Height float64
	Speed  float64
	Active bool
}

func (b *Bullet) GetBounds() (float64, float64, float64, float64) {
	return b.X, b.Y, b.Width, b.Height
}
