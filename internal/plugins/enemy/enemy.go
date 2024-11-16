package enemy

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/camera"
	"game/internal/plugins/enemy/entities"
	entity "game/internal/plugins/enemy/entities"
	"game/internal/plugins/enemy/factory"
	"game/internal/plugins/enemy/templates"
	"game/internal/plugins/obstacle"
	"game/internal/plugins/player"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyPlugin struct {
	kernel       *core.GameKernel
	enemies      []*entity.Enemy
	spawnTimer   float64
	playerPlugin *player.PlayerPlugin
	animation    *assets.Animation
	staticasset  *assets.StaticSprite
}

func NewEnemyPlugin(playerPlugin *player.PlayerPlugin) *EnemyPlugin {
	return &EnemyPlugin{
		enemies:      []*entity.Enemy{},
		spawnTimer:   0,
		playerPlugin: playerPlugin,
	}
}

func (ep *EnemyPlugin) ID() string {
	return "EnemySystem"
}

func (ep *EnemyPlugin) Init(kernel *core.GameKernel) error {
	ep.kernel = kernel

	ep.animation = assets.NewAnimation(0.01)
	err := ep.animation.LoadFromJSON(
		"assets/images/player/gunner/run/tileset.json",
		"assets/images/player/gunner/run/tileset.png")

	if err != nil {
		log.Fatal("Failed to load enemy animation:", err)
	}

	ep.staticasset = assets.NewStaticSprite()
	err = ep.staticasset.Load("assets/images/enemies/troll/image.png")

	if err != nil {
		log.Fatal("Failed to load enemy asset:", err)
	}

	return nil
}

func (ep *EnemyPlugin) Update() error {
	ep.spawnTimer += ep.kernel.DeltaTime

	if ep.spawnTimer >= 1.0 {
		ep.Spawn()
		ep.spawnTimer = 0
	}

	playerX, playerY := ep.playerPlugin.GetPosition()
	obstaclePlugin := ep.kernel.PluginManager.GetPlugin("ObstacleSystem").(*obstacle.ObstaclePlugin)

	for _, enemy := range ep.enemies {
		if enemy.Active {
			ep.moveTowardsPlayer(enemy, playerX, playerY, obstaclePlugin)

			// Verificar colisão com o jogador e cooldown de dano
			if ep.checkCollisionWithPlayer(enemy) {
				if enemy.LastDamageTime >= 0.5 {
					ep.playerPlugin.DecreaseHealth(enemy.Power)
					enemy.LastDamageTime = 0
				} else {
					enemy.LastDamageTime += ep.kernel.DeltaTime
				}
			}
		}

	}

	ep.animation.Update(ep.kernel.DeltaTime)

	return nil
}

func (ep *EnemyPlugin) checkCollisionWithPlayer(enemy *entities.Enemy) bool {
	playerX, playerY := ep.playerPlugin.GetPosition()
	playerWidth, playerHeight := ep.playerPlugin.GetSize()

	return math.Abs(enemy.X-playerX) < (enemy.Width+playerWidth)/2 &&
		math.Abs(enemy.Y-playerY) < (enemy.Height+playerHeight)/2
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
	cameraPlugin := ep.kernel.PluginManager.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, enemy := range ep.enemies {
		if enemy.Active {
			// Draw enemy relative to camera position
			screenX := enemy.X - cameraX
			screenY := enemy.Y - cameraY

			// Only draw if on screen
			if screenX >= -enemy.Width && screenX <= constants.ScreenWidth+enemy.Width &&
				screenY >= -enemy.Height && screenY <= constants.ScreenHeight+enemy.Height {

				// ep.animation.Draw(screen, screenX, screenY, enemy.Width, enemy.Height, true)
				ep.staticasset.Draw(screen, screenX, screenY, false)
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

	ep.enemies = append(ep.enemies, factory.CreateEnemy(templates.TankEnemy, x, y))
}

func (ep *EnemyPlugin) GetEnemies() []*entity.Enemy {
	return ep.enemies
}

func (ep *EnemyPlugin) moveTowardsPlayer(enemy *entity.Enemy, playerX, playerY float64, obstaclePlugin *obstacle.ObstaclePlugin) {
	// Calculate desired velocity towards player
	dx := playerX - enemy.X
	dy := playerY - enemy.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		dx /= distance
		dy /= distance

		desiredX := dx * enemy.Speed * ep.kernel.DeltaTime
		desiredY := dy * enemy.Speed * ep.kernel.DeltaTime

		// Check for obstacles and steer around them
		if !obstaclePlugin.CheckCollisionRect(enemy.X+desiredX, enemy.Y+desiredY, enemy.Width, enemy.Height) &&
			!ep.checkEnemyCollision(enemy.X+desiredX, enemy.Y+desiredY, enemy) {
			enemy.X += desiredX
			enemy.Y += desiredY
		} else {
			// Move perpendicularly to avoid collision
			perpendicularDirections := [][2]float64{
				{0, -enemy.Speed * ep.kernel.DeltaTime}, // Up
				{0, enemy.Speed * ep.kernel.DeltaTime},  // Down
			}

			// Determine which perpendicular direction is closer to the player
			if math.Abs(playerY-(enemy.Y-perpendicularDirections[0][1])) < math.Abs(playerY-(enemy.Y+perpendicularDirections[1][1])) {
				perpendicularDirections = [][2]float64{
					{0, -enemy.Speed * ep.kernel.DeltaTime}, // Up
					{0, enemy.Speed * ep.kernel.DeltaTime},  // Down
				}
			} else {
				perpendicularDirections = [][2]float64{
					{0, enemy.Speed * ep.kernel.DeltaTime},  // Down
					{0, -enemy.Speed * ep.kernel.DeltaTime}, // Up
				}
			}

			moved := false
			for _, dir := range perpendicularDirections {
				if !obstaclePlugin.CheckCollisionRect(enemy.X+dir[0], enemy.Y+dir[1], enemy.Width, enemy.Height) &&
					!ep.checkEnemyCollision(enemy.X+dir[0], enemy.Y+dir[1], enemy) {
					enemy.X += dir[0]
					enemy.Y += dir[1]
					moved = true
					break
				}
			}

			// If no direction is found, revert to old position and move randomly left or right
			if !moved {
				directions := [][2]float64{
					{-enemy.Speed * ep.kernel.DeltaTime, 0}, // Left
					{enemy.Speed * ep.kernel.DeltaTime, 0},  // Right
					{0, -enemy.Speed * ep.kernel.DeltaTime}, // Up
					{0, enemy.Speed * ep.kernel.DeltaTime},  // Down
				}

				for {
					randomDirection := rand.Intn(4)
					dir := directions[randomDirection]
					if !obstaclePlugin.CheckCollisionRect(enemy.X+dir[0], enemy.Y+dir[1], enemy.Width, enemy.Height) {
						enemy.X += dir[0]
						enemy.Y += dir[1]
						break
					}
				}
			}
		}
	}
}
