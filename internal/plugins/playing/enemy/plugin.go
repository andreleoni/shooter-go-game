package enemy

import (
	"fmt"
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/helpers/collision"

	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/enemy/entities"
	entity "game/internal/plugins/playing/enemy/entities"
	"game/internal/plugins/playing/enemy/factory"
	"game/internal/plugins/playing/enemy/templates"
	"game/internal/plugins/playing/player"

	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EnemyPlugin struct {
	kernel  *core.GameKernel
	plugins *core.PluginManager

	enemies      []*entity.Enemy
	spawnTimer   float64
	playerPlugin *player.PlayerPlugin
	StaticAsset  *assets.StaticSprite
}

func NewEnemyPlugin(playerPlugin *player.PlayerPlugin, plugins *core.PluginManager) *EnemyPlugin {
	return &EnemyPlugin{
		enemies:      []*entity.Enemy{},
		spawnTimer:   0,
		playerPlugin: playerPlugin,
		plugins:      plugins,
	}
}

func (ep *EnemyPlugin) ID() string {
	return "EnemySystem"
}

func (ep *EnemyPlugin) Init(kernel *core.GameKernel) error {
	ep.kernel = kernel

	// Assign sprite to enemy
	for enemytype, template := range templates.EnemyTemplates {
		enemypath := "assets/images/enemies/" + fmt.Sprint(enemytype) + ".png"

		ep.StaticAsset = assets.NewStaticSprite()
		err := ep.StaticAsset.Load(enemypath)
		if err != nil {
			log.Fatal("Failed to load enemy asset:", err)
		}

		template.StaticSprite = ep.StaticAsset
	}

	return nil
}

func (ep *EnemyPlugin) Update() error {
	ep.spawnTimer += ep.kernel.DeltaTime

	if ep.spawnTimer >= 0.5 {
		ep.Spawn()
		ep.spawnTimer = 0
	}

	playerX, playerY := ep.playerPlugin.GetPosition()
	playerWidth, playerHeight := ep.playerPlugin.GetSize()

	for _, enemy := range ep.enemies {
		if enemy.Active {
			ep.moveTowardsPlayer(enemy, playerX, playerY)

			playerCollision := collision.Check(
				enemy.X, enemy.Y,
				enemy.Width, enemy.Height,
				(playerX - playerWidth/2), (playerY - playerHeight/2),
				playerWidth, playerHeight)

			if playerCollision {
				if enemy.LastDamageTime == 0 || enemy.LastDamageTime >= 0.5 {
					ep.playerPlugin.DecreaseHealth(enemy.Power)
					ep.playerPlugin.DamageFlashTime = 0.3

					if enemy.LastDamageTime == 0 {
						enemy.LastDamageTime += ep.kernel.DeltaTime
					} else {
						enemy.LastDamageTime = 0
					}
				} else {
					enemy.LastDamageTime += ep.kernel.DeltaTime
					ep.playerPlugin.DamageFlashTime += ep.kernel.DeltaTime
				}
			}
		}

		// Atualizar o temporizador de flash de dano
		if enemy.DamageFlashTime > 0 {
			enemy.DamageFlashTime -= ep.kernel.DeltaTime
		}
	}

	return nil
}

func (ep *EnemyPlugin) checkEnemyCollision(x, y float64, currentEnemy *entity.Enemy) bool {
	for _, enemy := range ep.enemies {
		if enemy != currentEnemy && enemy.Active {
			if math.Abs(enemy.X-x) < currentEnemy.Width && math.Abs(enemy.Y-y) < currentEnemy.Height {
				return true
			}
		}
	}

	return false
}

func (ep *EnemyPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := ep.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, enemy := range ep.enemies {
		if enemy.Active {
			// Draw enemy relative to camera position
			screenX := enemy.X - cameraX
			screenY := enemy.Y - cameraY

			// Only draw if on screen
			if screenX >= -enemy.Width && screenX <= constants.ScreenWidth+enemy.Width &&
				screenY >= -enemy.Height && screenY <= constants.ScreenHeight+enemy.Height {

				if enemy.DamageFlashTime > 0 {
					ebitenutil.DrawRect(screen, screenX, screenY, enemy.Width, enemy.Height, color.RGBA{255, 255, 255, 255})
				} else {
					ebitenutil.DrawRect(screen, screenX, screenY, enemy.Width, enemy.Height, color.RGBA{255, 0, 255, 255})

					input := assets.DrawInput{
						Width:  enemy.Width,
						Height: enemy.Height,
						X:      screenX,
						Y:      screenY,
					}

					enemy.Stats.StaticSprite.Draw(screen, input)
				}
			}
		}
	}
}

func (ep *EnemyPlugin) Spawn() {
	playerX, playerY := ep.playerPlugin.GetPosition()

	// Escolher uma borda aleatória (0: superior, 1: inferior, 2: esquerda, 3: direita)
	border := rand.Intn(4)
	var x, y float64

	switch border {
	case 0: // Superior
		x = playerX + rand.Float64()*constants.ScreenWidth - constants.ScreenWidth/2
		y = playerY - constants.ScreenHeight/2
	case 1: // Inferior
		x = playerX + rand.Float64()*constants.ScreenWidth - constants.ScreenWidth/2
		y = playerY + constants.ScreenHeight/2
	case 2: // Esquerda
		x = playerX - constants.ScreenWidth/2
		y = playerY + rand.Float64()*constants.ScreenHeight - constants.ScreenHeight/2
	case 3: // Direita
		x = playerX + constants.ScreenWidth/2
		y = playerY + rand.Float64()*constants.ScreenHeight - constants.ScreenHeight/2
	}

	// Garantir que a posição de spawn esteja dentro dos limites do mundo
	x = math.Max(0, math.Min(x, constants.WorldWidth))
	y = math.Max(0, math.Min(y, constants.WorldHeight))

	// Escolher um tipo aleatório de inimigo
	enemyType := entities.EnemyType(rand.Intn(len(templates.EnemyTemplates)))

	ep.enemies = append(ep.enemies, factory.CreateEnemy(enemyType, x, y))
}

func (ep *EnemyPlugin) GetEnemies() []*entity.Enemy {
	return ep.enemies
}

func (ep *EnemyPlugin) moveTowardsPlayer(enemy *entity.Enemy, playerX, playerY float64) {
	dx := playerX - enemy.X
	dy := playerY - enemy.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Normalizar a direção
	dx /= distance
	dy /= distance

	// Mover o inimigo em direção ao jogador
	enemy.X += dx * enemy.Speed * ep.kernel.DeltaTime
	enemy.Y += dy * enemy.Speed * ep.kernel.DeltaTime

	// Verificar colisão com outros inimigos
	if ep.checkEnemyCollision(enemy.X, enemy.Y, enemy) {
		enemy.X -= dx * enemy.Speed * ep.kernel.DeltaTime
		enemy.Y -= dy * enemy.Speed * ep.kernel.DeltaTime
	}
}
