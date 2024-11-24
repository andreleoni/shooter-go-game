package entities

type Weapon struct {
	Power       float64
	Type        WeaponType
	Projectiles []*Projectile
}
