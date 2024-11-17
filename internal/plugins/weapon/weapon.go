package weapon

import (
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/weapon/entities"
	"game/internal/plugins/weapon/templates"

	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type WeaponPlugin struct {
	kernel *core.GameKernel

	weapons []*entities.Weapon
}

func NewWeaponPlugin() *WeaponPlugin {
	return &WeaponPlugin{}
}

func (wp *WeaponPlugin) ID() string {
	return "WeaponSystem"
}

func (wp *WeaponPlugin) Init(kernel *core.GameKernel) error {
	wp.kernel = kernel

	wp.weapons = append(wp.weapons, &entities.Weapon{
		Power: 10,
		Type:  templates.DaggersWeapon,
	})

	return nil
}

func (wp *WeaponPlugin) Update() error {
	for _, weapon := range wp.weapons {
		if weapon.Type == templates.DaggersWeapon {
			for _, weapon := range weapon.Projectiles {
				weapon.X += weapon.DirectionX * weapon.Speed * wp.kernel.DeltaTime
				weapon.Y += weapon.DirectionY * weapon.Speed * wp.kernel.DeltaTime

				// Deactivate if off screen
				if weapon.X < 0 || weapon.X > constants.WorldWidth ||
					weapon.Y < 0 || weapon.Y > constants.WorldHeight {
					weapon.Active = false
				}
			}
		}
	}

	return nil
}

func (wp *WeaponPlugin) Draw(screen *ebiten.Image) {
	for _, weapon := range wp.weapons {
		if weapon.Type == templates.DaggersWeapon {
			for _, weapon := range weapon.Projectiles {
				if weapon.Active {
					screenX := weapon.X
					screenY := weapon.Y

					ebitenutil.DrawRect(screen, screenX, screenY, 5, 5, color.RGBA{255, 255, 255, 255})
				}
			}
		}
	}
}

func (wp *WeaponPlugin) Shoot(x, y float64) {
	for _, weapon := range wp.weapons {
		for i := 0; i < 10; i++ {
			angle := rand.Float64() * 2 * math.Pi
			directionX := math.Cos(angle)
			directionY := math.Sin(angle)

			projectile := &entities.Projectile{
				X:          x,
				Y:          y,
				Speed:      300,
				DirectionX: directionX,
				DirectionY: directionY,
				Active:     true,
			}

			weapon.Projectiles = append(weapon.Projectiles, projectile)
		}
	}
}
