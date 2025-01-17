package enemy

import (
	"fmt"
	"game/internal/assets"
	"game/internal/config"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/helpers/collision"

	"game/internal/plugins/menu/fontface"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/enemy/entities"
	entity "game/internal/plugins/playing/enemy/entities"
	"game/internal/plugins/playing/enemy/factory"
	"game/internal/plugins/playing/enemy/templates"
	"game/internal/plugins/playing/player"

	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type DamageInfo struct {
	X, Y  float64
	Value float64
	Color color.Color
	Timer float64
}

type EnemyPlugin struct {
	kernel  *core.GameKernel
	plugins *core.PluginManager

	enemies         []*entity.Enemy
	inactiveEnemies []*entity.Enemy

	deathEnemies []*entity.Enemy

	spawnTimer   float64
	playerPlugin *player.PlayerPlugin
	StaticAsset  *assets.StaticSprite

	damages []DamageInfo // Lista de danos causados pelos inimigos

	maxEnemies int // Número máximo de inimigos por nível

	gameTimer  float64
	spawnRates map[int]float64 // Maps minutes to spawn delay

	baseStats map[int]EnemyBaseStats

	globalProjectiles []*entity.Projectile
}

type EnemyBaseStats struct {
	Health float64
	Damage float64
}

func NewEnemyPlugin(playerPlugin *player.PlayerPlugin, plugins *core.PluginManager) *EnemyPlugin {
	return &EnemyPlugin{
		enemies:      []*entity.Enemy{},
		spawnTimer:   0,
		playerPlugin: playerPlugin,
		plugins:      plugins,
		maxEnemies:   10,
		spawnRates: map[int]float64{
			0: 2.0,
			1: 1.0,
			2: 0.7,
			3: 0.5,
			4: 0.3,
			5: 0.2,
			6: 0.1,
		},
		baseStats: map[int]EnemyBaseStats{
			0: {Health: 0.1, Damage: 0.1},
			1: {Health: 0.2, Damage: 0.2},
			2: {Health: 0.4, Damage: 0.4},
			3: {Health: 0.6, Damage: 0.6},
			4: {Health: 0.8, Damage: 0.8},
			5: {Health: 0.9, Damage: 0.9},
			6: {Health: 1.0, Damage: 1.0},
		},
	}
}

func (ep *EnemyPlugin) ID() string {
	return "EnemySystem"
}

func (ep *EnemyPlugin) Init(kernel *core.GameKernel) error {
	ep.kernel = kernel
	ep.globalProjectiles = []*entity.Projectile{}

	return nil
}

func (ep *EnemyPlugin) Update() error {
	// Update game timer
	ep.gameTimer += ep.kernel.DeltaTime
	ep.spawnTimer += ep.kernel.DeltaTime

	// Get current minute
	currentMinute := int(ep.gameTimer / 60)

	// Get spawn delay based on current minute
	spawnDelay := ep.spawnRates[6] // Default to highest difficulty
	if delay, exists := ep.spawnRates[currentMinute]; exists {
		spawnDelay = delay
	}

	// Spawn enemy based on current spawn rate
	if ep.spawnTimer >= spawnDelay {
		ep.Spawn()
		ep.spawnTimer = 0
	}

	playerX, playerY := ep.playerPlugin.GetPosition()
	playerWidth, playerHeight := ep.playerPlugin.GetSize()

	cameraPlugin := ep.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for i, enemy := range ep.enemies {
		if enemy.Active {
			ep.moveTowardsPlayer(enemy, playerX, playerY)

			if enemy.IsEnemyMovingRight(playerX) {
				enemy.CurrentAnimation = enemy.RunningRightAnimationSprite
			} else {
				enemy.CurrentAnimation = enemy.RunningLeftAnimationSprite
			}

			enemy.CurrentAnimation.Update(ep.kernel.DeltaTime)

			playerCollision := collision.Check(
				enemy.X, enemy.Y,
				enemy.Width, enemy.Height,
				(playerX - playerWidth/2), (playerY - playerHeight/2),
				playerWidth, playerHeight)

			if playerCollision {
				if enemy.LastPlayerDamageTime == 0 || enemy.LastPlayerDamageTime >= 0.5 {
					ep.playerPlugin.ApplyDamage(enemy.Power)
					ep.playerPlugin.DamageFlashTime = 0.3

					if enemy.LastPlayerDamageTime == 0 {
						enemy.LastPlayerDamageTime += ep.kernel.DeltaTime
					} else {
						enemy.LastPlayerDamageTime = 0
					}
				} else {
					enemy.LastPlayerDamageTime += ep.kernel.DeltaTime
					ep.playerPlugin.DamageFlashTime += ep.kernel.DeltaTime
				}
			}

			// Margem extra de 200 pixels além da tela
			const marginBeyondScreen = 200

			// Verificar se o inimigo está muito além dos limites da tela
			if enemy.X < cameraX-marginBeyondScreen ||
				enemy.X > cameraX+constants.ScreenWidth+marginBeyondScreen ||
				enemy.Y < cameraY-marginBeyondScreen ||
				enemy.Y > cameraY+constants.ScreenHeight+marginBeyondScreen {

				enemy.Active = false
				ep.inactiveEnemies = append(ep.inactiveEnemies, enemy)
				ep.enemies = append(ep.enemies[:i], ep.enemies[i+1:]...)
			}

			if enemy.Name == "ranged" {
				// Calculate distance to player
				dx := playerX - enemy.X
				dy := playerY - enemy.Y
				distance := math.Sqrt(dx*dx + dy*dy)

				// Attack if in range
				if distance <= enemy.AttackRange {
					enemy.AttackCooldown -= ep.kernel.DeltaTime
					if enemy.AttackCooldown <= 0 {
						ep.globalProjectiles = append(ep.globalProjectiles, enemy.Shoot(playerX, playerY))

						enemy.AttackCooldown = 2.0 // Reset cooldown
					}
				}
			}
		}

		// Atualizar o temporizador de flash de dano
		if enemy.DamageFlashTime > 0 {
			enemy.DamageFlashTime -= ep.kernel.DeltaTime
		}
	}

	for _, enemy := range ep.deathEnemies {
		enemy.DeathAnimation.Update(ep.kernel.DeltaTime)

		if enemy.DeathAnimation.CurrentFrame+1 >= len(enemy.DeathAnimation.Frames) {
			ep.deathEnemies = ep.deathEnemies[:0]
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
					vector.DrawFilledRect(screen,
						float32(screenX),
						float32(screenY),
						float32(enemy.Width),
						float32(enemy.Height),
						color.RGBA{255, 255, 255, 255},
						true)

				} else {
					if config.IsDebugEnv() {
						vector.DrawFilledRect(screen,
							float32(screenX),
							float32(screenY),
							float32(enemy.Width),
							float32(enemy.Height),
							color.RGBA{255, 0, 255, 255},
							true)
					}

					input := assets.DrawInput{
						Width:  enemy.Width,
						Height: enemy.Height,
						X:      screenX,
						Y:      screenY,
					}

					if enemy.CurrentAnimation != nil {
						enemy.CurrentAnimation.Draw(screen, input)
					}
				}
			}
		}
	}

	// Draw global projectiles
	for _, p := range ep.globalProjectiles {
		if p.Active {
			screenX := p.X - cameraX
			screenY := p.Y - cameraY

			vector.DrawFilledRect(
				screen,
				float32(screenX),
				float32(screenY),
				float32(p.Width),
				float32(p.Height),
				color.RGBA{255, 0, 0, 255},
				true,
			)
		}
	}

	// Desenhar os inimigos mortos
	for _, enemy := range ep.deathEnemies {
		screenX := enemy.X - cameraX
		screenY := enemy.Y - cameraY

		if enemy.DeathAnimation.CurrentFrame <= len(enemy.DeathAnimation.Frames) {
			input := assets.DrawInput{
				Width:  enemy.Width,
				Height: enemy.Height,
				X:      screenX,
				Y:      screenY,
			}

			enemy.DeathAnimation.Draw(screen, input)
		}
	}

	for i := len(ep.damages) - 1; i >= 0; i-- {
		damage := &ep.damages[i]
		damage.Timer -= ep.kernel.DeltaTime

		if damage.Timer <= 0 {
			ep.damages = append(ep.damages[:i], ep.damages[i+1:]...)
		} else {
			screenX := damage.X - cameraX
			screenY := damage.Y - cameraY - (1.0-damage.Timer)*20

			drawValue := int(damage.Value)
			text.Draw(screen, fmt.Sprintf("%d", drawValue), basicfont.Face7x13, int(screenX), int(screenY), damage.Color)
		}
	}

	minutes := int(ep.gameTimer / 60)
	seconds := int(ep.gameTimer) % 60
	timerText := fmt.Sprintf("%02d:%02d", minutes, seconds)

	text.Draw(
		screen,
		timerText,
		fontface.FontFace,
		500, // X position
		40,  // Y position
		color.White)
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

	var enemy *entity.Enemy

	if len(ep.inactiveEnemies) > 0 {
		// Reutilizar um inimigo inativo
		enemy = ep.inactiveEnemies[len(ep.inactiveEnemies)-1]
		ep.inactiveEnemies = ep.inactiveEnemies[:len(ep.inactiveEnemies)-1]
		enemy.X = x
		enemy.Y = y
		enemy.Health = enemy.MaxHealth
		enemy.Active = true

	} else {
		// Criar um novo inimigo
		enemyType := entities.EnemyType(rand.Intn(len(templates.EnemyTemplates)))
		enemy = factory.CreateEnemy(enemyType, x, y)
	}

	getEnemyStats := ep.getEnemyStats()

	enemy.Health = enemy.MaxHealth + (enemy.MaxHealth * getEnemyStats.Health)
	enemy.Power = enemy.Power + (enemy.Power * getEnemyStats.Damage)

	ep.enemies = append(ep.enemies, enemy)
}

func (ep *EnemyPlugin) GetEnemies() []*entity.Enemy {
	return ep.enemies
}

func (ep *EnemyPlugin) SetEnemies(e []*entity.Enemy) {
	ep.enemies = e
}

func (ep *EnemyPlugin) moveTowardsPlayer(enemy *entity.Enemy, playerX, playerY float64) {
	// Direction to player
	dx := playerX - enemy.X
	dy := playerY - enemy.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Define minimum range based on enemy type
	minRange := 0.0
	if enemy.Name == "ranged" {
		minRange = 200.0 // Ranged enemies try to maintain this distance
	} else {
		minRange = 10.0 // Melee enemies get closer
	}

	// Adjust movement based on range
	if distance > 0 {
		dx /= distance
		dy /= distance

		// If ranged and too close, move away from player
		if enemy.Name == "ranged" && distance < minRange {
			dx = -dx
			dy = -dy
		}
	}

	// Simple separation to avoid stacking
	separationX, separationY := 0.0, 0.0
	for _, other := range ep.enemies {
		if other != enemy && other.Active {
			diffX := enemy.X - other.X
			diffY := enemy.Y - other.Y
			dist := math.Sqrt(diffX*diffX + diffY*diffY)

			if dist < 40 { // Minimal separation distance
				separationX += diffX * 0.3
				separationY += diffY * 0.3
			}
		}
	}

	// Combine movement forces
	moveX := dx*0.7 + separationX*0.3
	moveY := dy*0.7 + separationY*0.3

	// Update position
	enemy.X += moveX * enemy.Speed * ep.kernel.DeltaTime
	enemy.Y += moveY * enemy.Speed * ep.kernel.DeltaTime
}

func (ep *EnemyPlugin) ApplyDamage(enemy *entities.Enemy, damage float64, isCriticalDamage bool) {
	// Aplicar a armadura para reduzir o dano
	effectiveDamage := damage
	enemy.Health -= effectiveDamage

	if enemy.Health <= 0 {
		enemy.Health = 0
	}

	colorByType := map[bool]color.Color{
		true:  color.RGBA{255, 0, 0, 200},
		false: color.RGBA{255, 255, 255, 200},
	}

	damageText := fmt.Sprintf("%d", int(effectiveDamage))
	bounds := text.BoundString(basicfont.Face7x13, damageText)
	textWidth := bounds.Dx()
	textHeight := bounds.Dy()

	textX := enemy.X + enemy.Width/2 - float64(textWidth)/2
	textY := enemy.Y - float64(textHeight)/2

	ep.damages = append(ep.damages, DamageInfo{
		X:     textX,
		Y:     textY,
		Value: effectiveDamage,
		Color: colorByType[isCriticalDamage],
		Timer: 0.4,
	})

}

func (ep *EnemyPlugin) AddDeathEnemies(e *entity.Enemy) {
	ep.deathEnemies = append(ep.deathEnemies, e)
}

func (ep *EnemyPlugin) getEnemyStats() EnemyBaseStats {
	currentMinute := int(ep.gameTimer / 60)

	stats := ep.baseStats[6]

	if baseStats, exists := ep.baseStats[currentMinute]; exists {
		stats = baseStats
	}

	return stats
}

func (ep *EnemyPlugin) GetGlobalProjectiles() []*entity.Projectile {
	return ep.globalProjectiles
}
