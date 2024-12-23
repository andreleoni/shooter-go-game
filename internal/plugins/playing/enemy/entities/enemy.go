package entities

type Enemy struct {
	X, Y                            float64
	Width                           float64
	Height                          float64
	Active                          bool
	Speed                           float64
	Type                            EnemyType
	Stats                           EnemyStats
	Health                          float64
	Power                           float64
	LastPlayerDamageTime            float64
	LastAreaDamageDeltaTimeByWeapon map[string]float64
	DamageFlashTime                 float64
}

func (e *Enemy) GetBounds() (float64, float64, float64, float64) {
	return e.X, e.Y, e.Width, e.Height
}
