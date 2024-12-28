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

	spawnTimer   float64
	playerPlugin *player.PlayerPlugin
	StaticAsset  *assets.StaticSprite

	damages []DamageInfo // Lista de danos causados pelos inimigos

	maxEnemies int // Número máximo de inimigos por nível
}

func NewEnemyPlugin(playerPlugin *player.PlayerPlugin, plugins *core.PluginManager) *EnemyPlugin {
	return &EnemyPlugin{
		enemies:      []*entity.Enemy{},
		spawnTimer:   0,
		playerPlugin: playerPlugin,
		plugins:      plugins,
		maxEnemies:   10,
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

	cameraPlugin := ep.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for i, enemy := range ep.enemies {
		if enemy.Active {
			ep.moveTowardsPlayer(enemy, playerX, playerY)

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

			// Verificar se o inimigo está fora dos limites da tela
			if enemy.X < cameraX || enemy.X > cameraX+constants.ScreenWidth ||
				enemy.Y < cameraY || enemy.Y > cameraY+constants.ScreenHeight {
				enemy.Active = false
				ep.inactiveEnemies = append(ep.inactiveEnemies, enemy)
				ep.enemies = append(ep.enemies[:i], ep.enemies[i+1:]...)
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
					vector.DrawFilledRect(screen,
						float32(screenX),
						float32(screenY),
						float32(enemy.Width),
						float32(enemy.Height),
						color.RGBA{255, 255, 255, 255},
						true)
				} else {
					vector.DrawFilledRect(screen,
						float32(screenX),
						float32(screenY),
						float32(enemy.Width),
						float32(enemy.Height),
						color.RGBA{255, 0, 255, 255},
						true)

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

	// Desenhar os danos
	for i := len(ep.damages) - 1; i >= 0; i-- {
		damage := &ep.damages[i]
		damage.Timer -= ep.kernel.DeltaTime

		if damage.Timer <= 0 {
			// Remover o dano da lista
			ep.damages = append(ep.damages[:i], ep.damages[i+1:]...)
		} else {
			// Desenhar o dano
			screenX := damage.X - cameraX
			screenY := damage.Y - cameraY - (1.0-damage.Timer)*20 // Mover o texto para cima ao longo do tempo

			drawValue := int(damage.Value)
			text.Draw(screen, fmt.Sprintf("%d", drawValue), basicfont.Face7x13, int(screenX), int(screenY), damage.Color)
		}
	}
}

func (ep *EnemyPlugin) Spawn() {
	if len(ep.enemies) >= ep.maxEnemies {
		return
	}

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

	var enemy *entity.Enemy

	fmt.Println("Enemies:", ep.enemies, ep.inactiveEnemies)

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

	ep.enemies = append(ep.enemies, enemy)
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

func (ep *EnemyPlugin) ApplyDamage(enemy *entities.Enemy, damage float64, isCriticalDamage bool) {
	// Aplicar a armadura para reduzir o dano
	effectiveDamage := damage
	enemy.Health -= effectiveDamage

	if enemy.Health < 0 {
		enemy.Health = 0
	}
	colorByType := map[bool]color.Color{
		true:  color.RGBA{255, 0, 0, 255},
		false: color.RGBA{255, 255, 255, 255},
	}

	// Adicionar uma nova entrada de dano
	ep.damages = append(ep.damages, DamageInfo{
		X:     enemy.X,
		Y:     enemy.Y,
		Value: effectiveDamage,
		Color: colorByType[isCriticalDamage],
		Timer: 0.4,
	})
}
