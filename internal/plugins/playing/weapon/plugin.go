package weapon

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/weapon/entities"
	"game/internal/plugins/playing/weapon/templates"

	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type WeaponPlugin struct {
	kernel  *core.GameKernel
	plugins *core.PluginManager

	weapons []*entities.Weapon
}

func NewWeaponPlugin(plugins *core.PluginManager) *WeaponPlugin {
	return &WeaponPlugin{
		plugins: plugins,
	}
}

func (wp *WeaponPlugin) ID() string {
	return "WeaponSystem"
}

func (wp *WeaponPlugin) Init(
	kernel *core.GameKernel,
) error {

	wp.kernel = kernel

	wp.weapons = append(wp.weapons, &entities.Weapon{
		Power: 20,
		Type:  templates.BasicWeapon,
	})

	return nil
}

func (wp *WeaponPlugin) AddWeapon(weapon *entities.Weapon) {
	wp.weapons = append(wp.weapons, weapon)
}

func (wp *WeaponPlugin) GetWeapons() []*entities.Weapon {
	return wp.weapons
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
		} else if weapon.Type == templates.BasicWeapon {
			for _, projectile := range weapon.Projectiles {
				if projectile.Active {
					dx := projectile.DirectionX
					dy := projectile.DirectionY

					// Calculate distance
					distance := math.Sqrt(dx*dx + dy*dy)

					if distance > 0 {
						// Normalize direction
						dx /= distance
						dy /= distance

						// Update position
						projectile.X += dx * projectile.Speed * wp.kernel.DeltaTime
						projectile.Y += dy * projectile.Speed * wp.kernel.DeltaTime
					}

					// Deactivate if off screen
					if projectile.X < 0 ||
						projectile.X > constants.WorldHeight ||
						projectile.Y < 0 ||
						projectile.Y > constants.WorldWidth {

						projectile.Active = false
					}
				}
			}
		}
	}

	return nil
}

func (wp *WeaponPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := wp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	playerPlugin := wp.plugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)
	playerX, playerY := playerPlugin.GetPosition()

	for _, weapon := range wp.weapons {
		if weapon.Type == templates.DaggersWeapon {
			for _, weapon := range weapon.Projectiles {
				if weapon.Active {
					screenX := weapon.X - cameraX
					screenY := weapon.Y - cameraY

					ebitenutil.DrawRect(screen, screenX, screenY, 5, 5, color.RGBA{255, 255, 0, 255})
				}
			}
		} else if weapon.Type == templates.BasicWeapon {
			for _, projectile := range weapon.Projectiles {
				if projectile.Active {
					// Draw bullet relative to camera position
					screenX := projectile.X - cameraX
					screenY := projectile.Y - cameraY

					// Only draw if on screen
					if screenX >= -5 && screenX <= constants.ScreenWidth+5 &&
						screenY >= -5 && screenY <= constants.ScreenHeight+5 {

						angle := math.Atan2(
							projectile.DirectionY, projectile.DirectionX)

						staticsprite := assets.NewStaticSprite()
						staticsprite.Load("assets/images/bullets/arrow/arrow.png")

						staticsprite.DrawAngle(screen, screenX, screenY, angle)
					}
				}
			}
		} else if weapon.Type == templates.ProtectionWeapon {
			screenX := playerX - cameraX
			screenY := playerY - cameraY

			circleX := screenX
			circleY := screenY

			ebitenutil.DrawCircle(screen, circleX, circleY, 50, color.RGBA{111, 222, 111, 255})
		}
	}
}

func (wp *WeaponPlugin) Shoot(x, y float64) {
	for _, weapon := range wp.weapons {
		if weapon.Type == templates.DaggersWeapon {
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
					Power:      weapon.Power,
				}

				weapon.Projectiles = append(weapon.Projectiles, projectile)
			}
		} else if weapon.Type == templates.BasicWeapon {
			// Get enemy plugin to find closest enemy
			enemyPlugin := wp.plugins.GetPlugin("EnemySystem").(plugins.EnemyPlugin)

			enemies := enemyPlugin.GetEnemies()

			if len(enemies) > 0 {
				// Find closest enemy
				closestEnemy := enemies[0]
				closestDist := math.MaxFloat64

				for _, enemy := range enemies {
					if !enemy.Active {
						continue
					}

					dx := (enemy.X + enemy.Width/2) - x
					dy := (enemy.Y + enemy.Height/2) - y
					dist := math.Sqrt(dx*dx + dy*dy)

					if dist < closestDist {
						closestDist = dist
						closestEnemy = enemy
					}
				}

				// Calcular direção
				dx := (closestEnemy.X + closestEnemy.Width/2) - x
				dy := (closestEnemy.Y + closestEnemy.Height/2) - y

				distance := math.Sqrt(dx*dx + dy*dy)

				dirX := dx / distance
				dirY := dy / distance

				// Create bullet targeting closest enemy
				projectile := &entities.Projectile{
					X:          x,
					Y:          y,
					Speed:      300,
					Active:     true,
					Power:      weapon.Power,
					TargetX:    closestEnemy.X + closestEnemy.Width/2,
					TargetY:    closestEnemy.Y + closestEnemy.Height/2,
					DirectionX: dirX,
					DirectionY: dirY,
				}

				weapon.Projectiles = append(weapon.Projectiles, projectile)
			}
		}
	}
}
