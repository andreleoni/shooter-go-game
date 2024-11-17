package entities

type Projectile struct {
	Active bool
	Power  float64

	Type WeaponType

	X, Y       float64
	Speed      float64
	DirectionX float64
	DirectionY float64
}
