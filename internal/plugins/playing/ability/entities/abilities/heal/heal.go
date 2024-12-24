package heal

// import (
// 	"game/internal/core"
// 	"game/internal/plugins/playing/ability/entities"

// 	"github.com/hajimehoshi/ebiten/v2"
// )

// type Dagger struct {
// 	plugins     *core.PluginManager
// 	Projectiles []*entities.Projectile
// 	Power       float64

// 	// Shoot cooldown
// 	ShootTimer    float64
// 	ShootCooldown float64
// }

// func ID() string {
// 	return "Heal"
// }

// func (d *Dagger) SetPluginManager(plugins *core.PluginManager) {
// 	d.plugins = plugins
// }

// 	Shoot(x, y float64)
// 	Update(wui WeaponUpdateInput)
// 	Draw(screen *ebiten.Image, wdi WeaponDrawInput)

// 	ActiveProjectiles() []*entities.Projectile

// 	GetPower() float64

// 	DamageType() string

// 	AttackSpeed() float64
