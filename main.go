package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	mapWidth     = 2000
	mapHeight    = 2000
	screenWidth  = 1080
	screenHeight = 720
)

const (
	enemySeparationDistance = 30.0 // Minimum distance between enemies
	separationForce         = 0.5  // Force of separation
)

type Game struct {
	gameOver        bool
	choosingPowerUp bool
	Obstacles       []Obstacle
	powerUpOptions  []string
}

func DrawPowerUpChoice(screen *ebiten.Image) {
	cardWidth := 200.0
	cardHeight := 100.0
	cardX := (screenWidth - cardWidth) / 2
	cardY := (screenHeight - cardHeight) / 2

	// Draw the card background
	ebitenutil.DrawRect(screen, cardX, cardY, cardWidth, cardHeight, color.RGBA{255, 255, 255, 255})

	// Draw the text
	message := "Choose your power-up:\n1. Speed\n2. Power\n3. Radius"
	ebitenutil.DebugPrintAt(screen, message, int(cardX)+10, int(cardY)+10)
}

func ApplyPowerUpEffect(player *Player, powerUpType string) {
	// Ativa o modo de escolha de power-up e atualiza as opções disponíveis
	currentGame.choosingPowerUp = true

	// Atualiza o tipo de power-up que será aplicado após a escolha
	if powerUpType == "speed" {
		currentGame.powerUpOptions = []string{"speed"}

	} else if powerUpType == "radius" {
		currentGame.powerUpOptions = []string{"radius"}

	} else if powerUpType == "power" {
		currentGame.powerUpOptions = []string{"power"}
	}
}

type Obstacle struct {
	x, y          float64
	width, height float64
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil // Não atualiza se o jogo estiver terminado
	}

	if g.choosingPowerUp {
		if ebiten.IsKeyPressed(ebiten.Key1) {
			player.Speed += 1.0
			player.ActivePowerUps["speed"] = time.Now().Add(15 * time.Minute)
			g.choosingPowerUp = false
		} else if ebiten.IsKeyPressed(ebiten.Key2) {
			player.BulletSpeed += 0.2
			player.ActivePowerUps["power"] = time.Now().Add(15 * time.Minute)
			g.choosingPowerUp = false
		} else if ebiten.IsKeyPressed(ebiten.Key3) {
			player.CollectionRadius += 100.0
			player.ActivePowerUps["radius"] = time.Now().Add(15 * time.Minute)
			g.choosingPowerUp = false
		}

		return nil
	}

	player.Update()
	camera.Update(player.X, player.Y)

	UpdateEnemies(&player, g)
	UpdatePowerUps(&player)
	UpdateBullets()
	SpawnEnemies()
	AutoShoot(&player)

	// Verifica colisão com obstáculos para jogador e inimigos
	player.HandleObstacleCollision(g.Obstacles)

	// Verifica colisões
	CheckEnemyCollisions(g)

	UpdateXPItems(&player)

	return nil
}

func (p *Player) HandleObstacleCollision(obstacles []Obstacle) {
	for _, obstacle := range obstacles {
		if CheckCollision(p.X, p.Y, p.Width, p.Height, obstacle.x, obstacle.y, obstacle.width, obstacle.height) {
			// Resolve a colisão ajustando a posição do jogador
			if p.X < obstacle.x && p.X+p.Width > obstacle.x && p.X+p.Width < obstacle.x+obstacle.width {
				p.X = obstacle.x - p.Width
			} else if p.X+p.Width > obstacle.x+obstacle.width && p.X < obstacle.x+obstacle.width {
				p.X = obstacle.x + obstacle.width
			}
			if p.Y < obstacle.y && p.Y+p.Height > obstacle.y && p.Y+p.Height < obstacle.y+obstacle.height {
				p.Y = obstacle.y - p.Height
			} else if p.Y+p.Height > obstacle.y+obstacle.height && p.Y < obstacle.y+obstacle.height {
				p.Y = obstacle.y + obstacle.height
			}
		}
	}
}

// Funções auxiliares para calcular a sobreposição
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fundo preto

	// Desenha o fundo do mapa
	mapBackgroundColor := color.RGBA{50, 50, 50, 255} // Cor de fundo do mapa (cinza escuro)
	ebitenutil.DrawRect(screen, -camera.X, -camera.Y, mapWidth, mapHeight, mapBackgroundColor)

	if g.gameOver {
		DrawGameOver(screen)
		return
	}

	if g.choosingPowerUp {
		DrawPowerUpChoice(screen)
		return
	}

	// Desenha o jogador, inimigos e power-ups
	player.Draw(screen)

	DrawEnemies(screen)
	DrawPowerUps(screen)
	DrawBullets(screen)

	// Desenha a barra de vida do jogador
	DrawHealthBar(screen, &player)

	// Desenha a barra de experiência do jogador
	DrawXPBar(screen, &player)

	// Desenha os itens de experiência
	DrawXPItems(screen)

	// Desenha o mapa e os objetos levando em conta a posição da câmera
	for _, obstacle := range g.Obstacles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(obstacle.x-camera.X, obstacle.y-camera.Y)

		ebitenutil.DrawRect(screen, obstacle.x-camera.X, obstacle.y-camera.Y, float64(obstacle.width), float64(obstacle.height), color.RGBA{33, 0, 0, 33})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var camera Camera

func main() {
	rand.Seed(time.Now().UnixNano())
	player.Init()
	InitEnemies()
	InitPowerUps()

	img, _, err := ebitenutil.NewImageFromFile("otsp_creatures_01.png")
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	currentGame = Game{}
	currentGame.generateObstacles(20, 16, 16, 64, 64)
	fmt.Println(currentGame)

	player.Avatar = img

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("leonigod")
	if err := ebiten.RunGame(&currentGame); err != nil {
		panic(err)
	}
}

var currentGame Game

func (g *Game) generateObstacles(count int, minWidth, minHeight, maxWidth, maxHeight float64) {

	for i := 0; i < count; i++ {
		var obstacle Obstacle

		validPosition := false

		for !validPosition {
			// Gera posições e tamanhos aleatórios para o obstáculo
			obstacle = Obstacle{
				x:      float64(rand.Intn(mapWidth)),
				y:      float64(rand.Intn(mapHeight)),
				width:  minWidth + float64(rand.Intn(int(maxWidth-minWidth))),
				height: minHeight + float64(rand.Intn(int(maxHeight-minHeight))),
			}

			// Verifica se o obstáculo não se sobrepõe a outros obstáculos
			validPosition = true
			for _, existingObstacle := range g.Obstacles {
				if checkOverlap(obstacle, existingObstacle, 0, 0) {
					validPosition = false
					break
				}
			}
		}

		g.Obstacles = append(g.Obstacles, obstacle)
	}
}

func checkOverlap(o1, o2 Obstacle, minDistanceX, minDistanceY float64) bool {
	return !(o1.x+o1.width+minDistanceX < o2.x || o1.x > o2.x+o2.width+minDistanceX ||
		o1.y+o1.height+minDistanceY < o2.y || o1.y > o2.y+o2.height+minDistanceY)
}

func DrawXPItems(screen *ebiten.Image) {
	for _, xpItem := range xpItems {
		if xpItem.Active {
			// Draw XP item
			ebitenutil.DrawCircle(screen, xpItem.X-camera.X, xpItem.Y-camera.Y, 5, color.RGBA{255, 255, 0, 255})
		}
	}

}

// game over

func DrawGameOver(screen *ebiten.Image) {
	// Desenha uma mensagem de Game Over
	message := "GAME OVER"
	ebitenutil.DebugPrint(screen, message)
}

// player

type Player struct {
	X, Y             float64
	Avatar           *ebiten.Image
	Speed            float64
	Width            float64
	Height           float64
	BulletSpeed      float64
	WeaponStrength   float64
	Health           float64
	Level            int
	XP               float64
	XPToNextLevel    float64
	CollectionRadius float64
	ActivePowerUps   map[string]time.Time // Mapeia o tipo de power-up para seu tempo de expiração
}

func DrawXPBar(screen *ebiten.Image, player *Player) {
	barWidth := 100.0
	barHeight := 10.0
	barX := (screenWidth - barWidth) / 2    // Centraliza a barra de experiência horizontalmente
	barY := screenHeight - barHeight - 30.0 // Posiciona a barra de experiência na parte inferior da tela, acima da barra de vida

	// Calcula a largura da barra de experiência com base na experiência atual do jogador
	xpPercentage := player.XP / player.XPToNextLevel
	currentBarWidth := barWidth * xpPercentage

	// Desenha o fundo da barra de experiência (cinza)
	ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{128, 128, 128, 255})

	// Desenha a barra de experiência atual (azul)
	ebitenutil.DrawRect(screen, barX, barY, currentBarWidth, barHeight, color.RGBA{0, 0, 255, 255})
}

type XPItem struct {
	X, Y     float64
	Width    float64
	Height   float64
	Active   bool
	VelX     float64
	VelY     float64
	Moving   bool
	Progress float64 // For easing
}

func UpdateXPItems(player *Player) {
	playerCenterX := player.X + player.Width/2
	playerCenterY := player.Y + player.Height/2

	for i := range xpItems {
		if !xpItems[i].Active {
			continue
		}

		// Calculate distance to player
		dx := playerCenterX - xpItems[i].X
		dy := playerCenterY - xpItems[i].Y
		dist := math.Sqrt(dx*dx + dy*dy)

		// Check if within collection radius
		if dist < player.CollectionRadius {
			xpItems[i].Moving = true
		}

		// Update movement with easing
		if xpItems[i].Moving {
			// Ease progress
			xpItems[i].Progress += 0.05
			if xpItems[i].Progress > 1 {
				xpItems[i].Progress = 1
			}

			// Apply easing function (cubic)
			t := xpItems[i].Progress
			ease := t * t * (3 - 2*t)

			// Calculate target position
			targetX := playerCenterX
			targetY := playerCenterY

			// Interpolate position
			xpItems[i].X = xpItems[i].X + (targetX-xpItems[i].X)*ease
			xpItems[i].Y = xpItems[i].Y + (targetY-xpItems[i].Y)*ease

			// Collect when very close
			if dist < 5.0 {
				xpItems[i].Active = false
				player.XP += 10
				if player.XP >= player.XPToNextLevel {
					player.LevelUp()
				}
			}
		}
	}
}

func (p *Player) LevelUp() {
	p.Level++
	p.XP = 0
	p.XPToNextLevel *= 1.5  // Aumenta a experiência necessária para o próximo nível
	p.WeaponStrength += 1.0 // Aumenta a força da arma
	p.Health = 100.0        // Restaura a vida do jogador
}

var xpItems []XPItem

var player Player

func (p *Player) Init() {
	p.X = screenWidth / 2
	p.Y = screenHeight / 2
	p.Speed = 2.0
	p.Width = 36
	p.Height = 36
	p.BulletSpeed = 0.0
	p.WeaponStrength = 5.0
	p.Health = 100.0
	p.XP = 0.0              // Adicionar valor inicial de XP
	p.XPToNextLevel = 100.0 // Adicionar valor necessário para o próximo nível
	p.Level = 1             // Adicionar nível inicial
	p.CollectionRadius = 50.0
	p.ActivePowerUps = make(map[string]time.Time)
}

type Camera struct {
	X, Y float64
}

func (c *Camera) Update(playerX, playerY float64) {
	c.X = playerX - screenWidth/2
	c.Y = playerY - screenHeight/2

	// Limita a câmera para não sair dos limites do mapa
	if c.X < 0 {
		c.X = 0
	}
	if c.Y < 0 {
		c.Y = 0
	}
	if c.X > mapWidth-screenWidth {
		c.X = mapWidth - screenWidth
	}
	if c.Y > mapHeight-screenHeight {
		c.Y = mapHeight - screenHeight
	}
}

func (p *Player) Update() {
	// originalX, originalY := p.X, p.Y

	if (ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft)) && p.X > 0 {
		p.X -= p.Speed
	}
	if (ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)) && p.X < mapWidth-p.Width {
		p.X += p.Speed
	}
	if (ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp)) && p.Y > 0 {
		p.Y -= p.Speed
	}
	if (ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown)) && p.Y < mapHeight-p.Height {
		p.Y += p.Speed
	}

	// Verifica colisão com obstáculos
	p.HandleObstacleCollision(currentGame.Obstacles)
}

func (p *Player) Draw(screen *ebiten.Image) {
	// Calcula a posição do jogador relativa à câmera
	px := p.X - camera.X
	py := p.Y - camera.Y

	// Desenha o jogador como um retângulo
	ebitenutil.DrawRect(screen, px, py, float64(p.Width), float64(p.Height), color.RGBA{0, 255, 2, 255})
}

// powerup

type PowerUp struct {
	X, Y          float64
	Width, Height float64
	Type          string // Tipo do power-up (velocidade ou força do tiro)
	Active        bool
}

var powerUps []PowerUp

func DrawHealthBar(screen *ebiten.Image, player *Player) {
	barWidth := 100.0
	barHeight := 10.0
	barX := (screenWidth - barWidth) - 500.0 // Centraliza a barra de vida horizontalmente
	barY := screenHeight - barHeight - 10.0  // Posiciona a barra de vida na parte inferior da tela

	// Calcula a largura da barra de vida com base na vida do jogador
	healthPercentage := player.Health / 100.0
	currentBarWidth := barWidth * healthPercentage

	// Desenha o fundo da barra de vida (vermelho)
	ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{255, 0, 0, 255})

	// Desenha a barra de vida atual (verde)
	ebitenutil.DrawRect(screen, barX, barY, currentBarWidth, barHeight, color.RGBA{0, 255, 0, 255})
}

func InitPowerUps() {
	for i := 0; i < 6; i++ {
		powerUp := PowerUp{
			X:      rand.Float64() * mapWidth,
			Y:      rand.Float64() * mapHeight,
			Width:  16,
			Height: 16,
			Active: true,
		}
		powerUps = append(powerUps, powerUp)
	}
}

func UpdatePowerUps(player *Player) {
	for i := range powerUps {
		if powerUps[i].Active && CheckCollision(player.X, player.Y, player.Width, player.Height, powerUps[i].X, powerUps[i].Y, powerUps[i].Width, powerUps[i].Height) {
			powerUps[i].Active = false
			ApplyPowerUpEffect(player, powerUps[i].Type)
		}
	}
}

func DrawPowerUps(screen *ebiten.Image) {
	for _, powerUp := range powerUps {
		if powerUp.Active {
			ebitenutil.DrawRect(screen, powerUp.X-camera.X, powerUp.Y-camera.Y, powerUp.Width, powerUp.Height, color.RGBA{0, 0, 255, 255})

			if powerUp.Type == "power" {
				ebitenutil.DrawRect(screen, powerUp.X-camera.X, powerUp.Y-camera.Y, powerUp.Width, powerUp.Height, color.RGBA{255, 255, 0, 255})
			}
		}
	}
}

// enemy

type Enemy struct {
	X, Y   float64
	Width  float64
	Height float64
	SpeedX float64
	SpeedY float64
	Speed  float64
	Active bool
	Health float64
	Attack float64
}

func UpdateEnemies(player *Player, g *Game) {
	for i := range enemies {
		enemies[i].Update(player, g)
	}
}

var enemies []Enemy

func InitEnemies() {
}

var lastSpawnTime time.Time
var spawnInterval = 3000 * time.Millisecond // Tempo entre spawns

func SpawnEnemies() {
	if time.Since(lastSpawnTime) < spawnInterval {
		return
	}

	if spawnInterval > 500*time.Millisecond {
		spawnInterval -= 200 * time.Millisecond // Aumenta a frequência de spawns
	}

	// Limite máximo de inimigos ativos
	if len(enemies) < 9000 { // Por exemplo, 10 inimigos ativos
		enemies = append(enemies, NewEnemy())
	}

	// remover, só ta aqui para testar o obstaculo
	lastSpawnTime = time.Now() //.Add(30000 * time.Minute)
}

func NewEnemy() Enemy {
	// Calcula as distâncias entre o jogador e cada borda da tela, considerando a posição da câmera
	distances := map[string]float64{
		"top":    player.Y - camera.Y,
		"bottom": (camera.Y + screenHeight) - player.Y,
		"left":   player.X - camera.X,
		"right":  (camera.X + screenWidth) - player.X,
	}

	// Encontra a borda mais próxima
	var closestBorder string
	minDistance := math.MaxFloat64
	for border, distance := range distances {
		if distance < minDistance {
			minDistance = distance
			closestBorder = border
		}
	}

	// Define a posição inicial do inimigo com base na borda mais próxima
	var x, y float64
	switch closestBorder {
	case "top":
		x = camera.X + rand.Float64()*screenWidth
		y = camera.Y - 20 // Acima da tela
	case "bottom":
		x = camera.X + rand.Float64()*screenWidth
		y = camera.Y + screenHeight + 20 // Abaixo da tela
	case "left":
		x = camera.X - 20 // À esquerda da tela
		y = camera.Y + rand.Float64()*screenHeight
	case "right":
		x = camera.X + screenWidth + 20 // À direita da tela
		y = camera.Y + rand.Float64()*screenHeight
	}

	return Enemy{
		X:      x,
		Y:      y,
		Width:  16,
		Height: 16,
		Active: true,
		Speed:  1.0,  // Velocidade do inimigo
		Health: 10.0, // Vida do inimigo
		Attack: 5.0,  // Pontos de ataque do inimigo
	}
}

func (e *Enemy) Update(p *Player, g *Game) {
	if !e.Active {
		return
	}

	e.MoveForwardPlayer(g, p.X, p.Y)
}

type RelativePosition struct {
	X string
	Y string
}

func (rp RelativePosition) Top() bool {
	return rp.Y == "top"
}

func (rp RelativePosition) Down() bool {
	return rp.Y == "down"
}

func (rp RelativePosition) Left() bool {
	return rp.X == "left"
}

func (rp RelativePosition) Right() bool {
	return rp.X == "right"
}

func (e *Enemy) GetRelativePosition(playerX, playerY float64) *RelativePosition {
	rp := RelativePosition{}

	if e.Y < playerY {
		rp.Y = "down"
	} else if e.Y > playerY {
		rp.Y = "top"
	}

	if e.X < playerX {
		rp.X = "right"
	} else if e.X > playerX {
		rp.X = "left"
	}

	return &rp
}

func (e *Enemy) GetClosestCorner(obstacleX, obstacleY, obstacleWidth, obstacleHeight float64) RelativePosition {
	// Calcula as distâncias entre o inimigo e cada canto do obstáculo
	topLeftDistance := math.Sqrt(math.Pow(e.X-obstacleX, 2) + math.Pow(e.Y-obstacleY, 2))
	topRightDistance := math.Sqrt(math.Pow((e.X+e.Width)-(obstacleX+obstacleWidth), 2) + math.Pow(e.Y-obstacleY, 2))
	bottomLeftDistance := math.Sqrt(math.Pow(e.X-obstacleX, 2) + math.Pow((e.Y+e.Height)-(obstacleY+obstacleHeight), 2))
	bottomRightDistance := math.Sqrt(math.Pow((e.X+e.Width)-(obstacleX+obstacleWidth), 2) + math.Pow((e.Y+e.Height)-(obstacleY+obstacleHeight), 2))

	// Determina o canto mais próximo
	minDistance := topLeftDistance
	closestCorner := RelativePosition{X: "left", Y: "top"}

	if topRightDistance < minDistance {
		minDistance = topRightDistance
		closestCorner = RelativePosition{X: "right", Y: "top"}
	}

	if bottomLeftDistance < minDistance {
		minDistance = bottomLeftDistance
		closestCorner = RelativePosition{X: "left", Y: "down"}
	}

	if bottomRightDistance < minDistance {
		closestCorner = RelativePosition{X: "right", Y: "down"}
	}

	return closestCorner
}

type ObstacleCollision struct {
	X, Y string
}

func (e *Enemy) MoveForwardPlayer(g *Game, playerX, playerY float64) {
	incollision := false

	// Calculate centers
	playerCenterX := playerX + player.Width/2
	playerCenterY := playerY + player.Height/2
	enemyCenterX := e.X + e.Width/2
	enemyCenterY := e.Y + e.Height/2

	// Calculate separation vector
	separationX := 0.0
	separationY := 0.0

	// Check collision with other enemies
	for _, other := range enemies {
		if other.Active && other != *e {
			// Calculate distance between enemies
			dx := e.X - other.X
			dy := e.Y - other.Y
			distance := math.Sqrt(dx*dx + dy*dy)

			// Apply separation if too close
			if distance < enemySeparationDistance && distance > 0 {
				// Calculate separation force
				strength := (enemySeparationDistance - distance) / enemySeparationDistance
				separationX += (dx / distance) * strength * separationForce
				separationY += (dy / distance) * strength * separationForce
			}
		}
	}

	// Handle obstacle collisions
	for _, obstacle := range g.Obstacles {
		relativePosition := e.GetRelativePosition(playerCenterX, playerCenterY)

		oc := ObstacleCollision{}

		if CheckCollision(e.X+e.Speed, e.Y, e.Width, e.Height, obstacle.x, obstacle.y, obstacle.width, obstacle.height) {
			oc.X = "right"

			if relativePosition.Top() {
				e.Y -= e.Speed
			} else if relativePosition.Down() {
				e.Y += e.Speed
			}

			incollision = true
		}

		if CheckCollision(e.X, e.Y+e.Speed, e.Width, e.Height, obstacle.x, obstacle.y, obstacle.width, obstacle.height) {
			oc.Y = "bottom"

			if relativePosition.Left() {
				e.X -= e.Speed
			} else if relativePosition.Right() {
				e.X += e.Speed
			}

			incollision = true
		}

		if CheckCollision(e.X-e.Speed, e.Y, e.Width, e.Height, obstacle.x, obstacle.y, obstacle.width, obstacle.height) {
			oc.X = "top"

			if relativePosition.Top() {
				e.Y -= e.Speed
			} else if relativePosition.Down() {
				e.Y += e.Speed
			}

			incollision = true
		}

		if CheckCollision(e.X, e.Y-e.Speed, e.Width, e.Height, obstacle.x, obstacle.y, obstacle.width, obstacle.height) {
			oc.Y = "left"

			if relativePosition.Left() {
				e.X -= e.Speed
			} else if relativePosition.Right() {
				e.X += e.Speed
			}

			incollision = true
		}

		if incollision {
			break
		}
	}

	// If no obstacle collision, move enemy
	if !incollision {
		// Calculate direction to player
		dx := playerCenterX - enemyCenterX
		dy := playerCenterY - enemyCenterY
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance > 0 {
			// Normalize direction
			dx = dx / distance
			dy = dy / distance

			// Add separation force to movement
			dx += separationX
			dy += separationY

			// Normalize combined vector
			totalForce := math.Sqrt(dx*dx + dy*dy)
			if totalForce > 0 {
				dx = dx / totalForce * e.Speed
				dy = dy / totalForce * e.Speed
			}

			// Move enemy
			e.X += dx
			e.Y += dy
		}
	}
}

func DrawEnemies(screen *ebiten.Image) {
	for _, enemy := range enemies {
		if enemy.Active {
			ebitenutil.DrawRect(screen, enemy.X-camera.X, enemy.Y-camera.Y, enemy.Width, enemy.Height, color.RGBA{255, 0, 0, 255}) // Vermelho
		}
	}
}

// collision

var lastShotTime time.Time
var shotInterval = 1 * time.Second // Tempo entre disparos

func AutoShoot(player *Player) {
	if time.Since(lastShotTime) >=
		shotInterval-time.Duration(player.BulletSpeed)*time.Second {

		FireBullet(player)        // Adiciona uma nova bala
		lastShotTime = time.Now() // Atualiza o tempo do último disparo
	}
}

func CheckCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 &&
		x1+w1 > x2 &&
		y1 < y2+h2 &&
		y1+h1 > y2
}

var lastEnemyCollisionAt time.Time

func CheckEnemyCollisions(g *Game) bool {
	// Verifica se um inimigo colidiu com o jogador
	for _, enemy := range enemies {
		if enemy.Active && player.X < enemy.X+enemy.Width && player.X+player.Width > enemy.X &&
			player.Y < enemy.Y+enemy.Height && player.Y+player.Height > enemy.Y {

			if time.Since(lastEnemyCollisionAt) > 1*time.Second {
				player.Health -= enemy.Attack
				lastEnemyCollisionAt = time.Now()
			}

			if player.Health <= 0 {
				g.gameOver = true // Termina o jogo se a vida do jogador chegar a zero
			}

			return true // Colisão entre jogador e inimigo
		}
	}

	// Verifica se uma bala colidiu com um inimigo
	for i := len(bullets) - 1; i >= 0; i-- {
		for j := len(enemies) - 1; j >= 0; j-- {
			if bullets[i].Active && enemies[j].Active &&
				bullets[i].X < enemies[j].X+enemies[j].Width && bullets[i].X+4 > enemies[j].X &&
				bullets[i].Y < enemies[j].Y+enemies[j].Height && bullets[i].Y+10 > enemies[j].Y {
				bullets[i].Active = false                  // Desativa a bala
				enemies[j].Health -= player.WeaponStrength // Reduz a vida do inimigo

				if enemies[j].Health <= 0 {
					enemies[j].Active = false // Desativa o inimigo se a vida for menor ou igual a 0

					// Dropar item de experiência quando o inimigo é derrotado
					xpItem := XPItem{
						X:      enemies[j].X,
						Y:      enemies[j].Y,
						Width:  10,
						Height: 10,
						Active: true,
					}

					xpItems = append(xpItems, xpItem)
				}

				return false
			}
		}
	}

	// Verifica se uma bala colidiu com um inimigo
	for i := len(bullets) - 1; i >= 0; i-- {
		for j := len(enemies) - 1; j >= 0; j-- {
			if bullets[i].Active && enemies[j].Active &&
				bullets[i].X < enemies[j].X+enemies[j].Width && bullets[i].X+4 > enemies[j].X &&
				bullets[i].Y < enemies[j].Y+enemies[j].Height && bullets[i].Y+10 > enemies[j].Y {
				bullets[i].Active = false                  // Desativa a bala
				enemies[j].Health -= player.WeaponStrength // Reduz a vida do inimigo

				if enemies[j].Health <= 0 {
					enemies[j].Active = false // Desativa o inimigo se a vida for menor ou igual a 0

					// Dropar item de experiência quando o inimigo é derrotado
					if enemies[j].Health <= 0 {
						xpItem := XPItem{
							X:      enemies[j].X,
							Y:      enemies[j].Y,
							Width:  10,
							Height: 10,
							Active: true,
						}
						xpItems = append(xpItems, xpItem)
					}
				}
				return false
			}
		}
	}

	return false
}

// bullet

type Bullet struct {
	X, Y       float64
	Speed      float64
	Active     bool
	DirectionX float64
	DirectionY float64
}

var bullets []Bullet

func GetNearestEnemy(player *Player) *Enemy {
	var nearestEnemy *Enemy
	minDistance := 1e9 // Um valor muito grande

	for i := range enemies {
		if enemies[i].Active {
			distance := math.Sqrt(math.Pow(enemies[i].X-player.X, 2) + math.Pow(enemies[i].Y-player.Y, 2))
			if distance < minDistance {
				minDistance = distance
				nearestEnemy = &enemies[i]
			}
		}
	}

	return nearestEnemy
}

func DrawBullets(screen *ebiten.Image) {
	for _, bullet := range bullets {
		if bullet.Active {
			ebitenutil.DrawRect(screen, bullet.X-camera.X, bullet.Y-camera.Y, 4, 10, color.RGBA{255, 255, 0, 255}) // Amarelo
		}
	}
}

func FireBullet(player *Player) {
	nearestEnemy := GetNearestEnemy(player)
	if nearestEnemy == nil {
		return // Não há inimigos ativos
	}

	// Calcula a direção
	directionX := nearestEnemy.X - player.X
	directionY := nearestEnemy.Y - player.Y
	length := math.Sqrt(directionX*directionX + directionY*directionY)

	// Normaliza a direção
	if length > 0 {
		directionX /= length
		directionY /= length
	}

	bullet := Bullet{
		X:          player.X + player.Width/2 - 2,  // Centraliza a bala em relação ao jogador
		Y:          player.Y + player.Height/2 - 5, // Centraliza a bala em relação ao jogador
		Speed:      5.0,
		Active:     true,
		DirectionX: directionX,
		DirectionY: directionY,
	}
	bullets = append(bullets, bullet)
}

func UpdateBullets() {
	for i := range bullets {
		if bullets[i].Active {
			bullets[i].X += bullets[i].DirectionX * bullets[i].Speed
			bullets[i].Y += bullets[i].DirectionY * bullets[i].Speed

			// Desativa a bala se sair da tela
			if bullets[i].Y < 0 || bullets[i].Y > mapHeight || bullets[i].X < 0 || bullets[i].X > mapWidth {
				bullets[i].Active = false
			}
		}
	}
}
