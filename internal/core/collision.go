package core

type Collidable interface {
	GetBounds() (x, y, width, height float64)
	IsActive() bool
}

type CollisionPlugin struct{}

func (cp *CollisionPlugin) CheckCollision(a, b Collidable) bool {
	ax, ay, aw, ah := a.GetBounds()
	bx, by, bw, bh := b.GetBounds()

	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}
