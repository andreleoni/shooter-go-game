package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	gameOver  bool
	Obstacles []Obstacle
}

type Obstacle struct {
	x, y          float64
	width, height float64
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil // Não atualiza se o jogo estiver terminado
	}

	player.Update()
	UpdateEnemies(&player, g)
	UpdatePowerUps(&player)
	UpdateBullets()
	SpawnEnemies()
	AutoShoot(&player)

	// Verifica colisão com obstáculos para jogador e inimigos
	player.HandleObstacleCollision(g.Obstacles)
	for i := range enemies {
		enemies[i].HandleObstacleCollision(g.Obstacles)
	}

	// Verifica colisões
	if CheckEnemyCollisions(g) {
		g.gameOver = true
	}

	return nil
}
func (p *Player) HandleObstacleCollision(obstacles []Obstacle) {
	for _, obstacle := range obstacles {
		if CheckCollision(p.X, p.Y, p.Width, p.Height, obstacle.x, obstacle.y, obstacle.width, obstacle.height) {
			// Ajuste horizontal
			if p.X < obstacle.x {
				p.X = obstacle.x - p.Width // Mover para a esquerda do obstáculo
			} else if p.X > obstacle.x+obstacle.width {
				p.X = obstacle.x + obstacle.width // Mover para a direita do obstáculo
			}

			// Ajuste vertical
			if p.Y < obstacle.y {
				p.Y = obstacle.y - p.Height // Mover acima do obstáculo
			} else if p.Y > obstacle.y+obstacle.height {
				p.Y = obstacle.y + obstacle.height // Mover abaixo do obstáculo
			}
		}
	}
}

func (e *Enemy) HandleObstacleCollision(obstacles []Obstacle) {
	for _, obstacle := range obstacles {
		if CheckCollision(e.X, e.Y, e.Width, e.Height, obstacle.x, obstacle.y, obstacle.width, obstacle.height) {
			// Calcula as profundidades de colisão nos eixos X e Y
			xOverlap := min(e.X+e.Width, obstacle.x+obstacle.width) - max(e.X, obstacle.x)
			yOverlap := min(e.Y+e.Height, obstacle.y+obstacle.height) - max(e.Y, obstacle.y)

			// Ajusta o eixo com menor sobreposição para evitar travamento
			if xOverlap < yOverlap {
				// Ajuste no eixo X
				if e.X < obstacle.x {
					e.X = obstacle.x - e.Width
				} else {
					e.X = obstacle.x + obstacle.width
				}
			} else {
				// Ajuste no eixo Y
				if e.Y < obstacle.y {
					e.Y = obstacle.y - e.Height
				} else {
					e.Y = obstacle.y + obstacle.height
				}
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

	if g.gameOver {
		DrawGameOver(screen)
		return
	}

	// Desenha o jogador, inimigos e power-ups
	player.Draw(screen)

	DrawEnemies(screen)
	DrawPowerUps(screen)
	DrawBullets(screen)

	// Draw each obstacle as a filled rectangle
	for _, obstacle := range g.Obstacles {
		ebitenutil.DrawRect(screen, obstacle.x, obstacle.y, float64(obstacle.width), float64(obstacle.height), color.RGBA{33, 0, 0, 33})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

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
	currentGame := Game{}
	currentGame.generateObstacles(10)
	fmt.Println(currentGame)

	player.Avatar = img

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Jogo Básico em Ebiten")
	if err := ebiten.RunGame(&currentGame); err != nil {
		panic(err)
	}
}

var currentGame Game

func (g *Game) generateObstacles(count int) {
	for i := 0; i < count; i++ {
		obstacle := Obstacle{
			x:      float64(rand.Intn(640)), // Random x position within window width
			y:      float64(rand.Intn(480)), // Random y position within window height
			width:  32,                      // Set obstacle width
			height: 32,                      // Set obstacle height
		}
		g.Obstacles = append(g.Obstacles, obstacle)
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
	X, Y           float64
	Avatar         *ebiten.Image
	Speed          float64
	Width          float64
	Height         float64
	BulletSpeed    float64              // Velocidade das balas
	ActivePowerUps map[string]time.Time // Mapeia o tipo de power-up para seu tempo de expiração

}

var player Player

func (p *Player) Init() {
	p.X = screenWidth / 2
	p.Y = screenHeight / 2
	p.Speed = 2.0
	p.Width = 16
	p.Height = 16
	p.ActivePowerUps = make(map[string]time.Time)
}

func (p *Player) Update() {
	if (ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft)) && p.X > 0 {
		p.X -= p.Speed
	}
	if ((ebiten.IsKeyPressed(ebiten.KeyD)) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)) && p.X < screenWidth-p.Width {
		p.X += p.Speed
	}
	if ((ebiten.IsKeyPressed(ebiten.KeyW)) || ebiten.IsKeyPressed(ebiten.KeyArrowUp)) && p.Y > 0 {
		p.Y -= p.Speed
	}
	if ((ebiten.IsKeyPressed(ebiten.KeyS)) || ebiten.IsKeyPressed(ebiten.KeyArrowDown)) && p.Y < screenHeight-p.Height {
		p.Y += p.Speed
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		FireBullet(p)
	}

	for powerUpType, expiration := range p.ActivePowerUps {
		if time.Now().After(expiration) {
			delete(p.ActivePowerUps, powerUpType) // Remove power-up expirado

			if powerUpType == "speed" {
				p.Speed -= 1.0 // Reverte o aumento de velocidade
			} else if powerUpType == "power" {
				p.BulletSpeed += 2.0 // Reverte o aumento da força do tiro
			}
		}
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	// ebitenutil.DrawRect(screen, p.X, p.Y, p.Width, p.Height, color.RGBA{0, 255, 0, 255}) // Verde

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X, p.Y)

	// Apply a color key filter to remove the pink background
	op.ColorM.Scale(1, 1, 1, 1)
	op.ColorM.Translate(-1, -0, -1, 0) // Adjust for transparency in pink areas

	// Desenha a área do sprite selecionada (exemplo: x: 0, y: 0, largura: 32, altura: 32)
	subImage := p.Avatar.SubImage(image.Rect(64, 64, 32, 32)).(*ebiten.Image)

	screen.DrawImage(subImage, op)
}

// powerup

type PowerUp struct {
	X, Y          float64
	Width, Height float64
	Type          string // Tipo do power-up (velocidade ou força do tiro)
	Active        bool
}

var powerUps []PowerUp

func InitPowerUps() {
	for i := 0; i < 3; i++ {
		powerUp := PowerUp{
			X:      rand.Float64() * screenWidth,
			Y:      rand.Float64() * screenHeight,
			Width:  16,
			Height: 16,
			Type:   "speed",
			Active: true,
		}
		powerUps = append(powerUps, powerUp)
	}

	for i := 0; i < 3; i++ {
		powerUp := PowerUp{
			X:      rand.Float64() * screenWidth,
			Y:      rand.Float64() * screenHeight,
			Width:  16,
			Height: 16,
			Type:   "power",
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
			ebitenutil.DrawRect(screen, powerUp.X, powerUp.Y, powerUp.Width, powerUp.Height, color.RGBA{0, 0, 255, 255})

			if powerUp.Type == "power" {
				ebitenutil.DrawRect(screen, powerUp.X, powerUp.Y, powerUp.Width, powerUp.Height, color.RGBA{255, 255, 0, 255})
			}

		}
	}
}

func ApplyPowerUpEffect(player *Player, powerUpType string) {
	if powerUpType == "speed" {
		player.Speed += 1.0                                                    // Aumenta a velocidade do jogador
		player.ActivePowerUps[powerUpType] = time.Now().Add(100 * time.Second) // Duração de 5 segundos
	} else if powerUpType == "power" {
		player.BulletSpeed += 2.0                                              // Aumenta a força do tiro
		player.ActivePowerUps[powerUpType] = time.Now().Add(100 * time.Second) // Duração de 5 segundos
	}

	fmt.Println("New player attributes: ", player)
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

	// Limite máximo de inimigos ativos
	if len(enemies) < 9000 { // Por exemplo, 10 inimigos ativos
		enemies = append(enemies, NewEnemy())
	}

	// remover, só ta aqui para testar o obstaculo
	lastSpawnTime = time.Now().Add(30000 * time.Minute)
}

func NewEnemy() Enemy {
	// Escolhe uma posição aleatória fora da tela
	var x, y float64
	spawnSide := rand.Intn(4) // 0: cima, 1: baixo, 2: esquerda, 3: direita
	switch spawnSide {
	case 0: // Cima
		x = rand.Float64() * screenWidth
		y = -20 // Acima da tela
	case 1: // Baixo
		x = rand.Float64() * screenWidth
		y = screenHeight + 20 // Abaixo da tela
	case 2: // Esquerda
		x = -20 // À esquerda da tela
		y = rand.Float64() * screenHeight
	case 3: // Direita
		x = screenWidth + 20 // À direita da tela
		y = rand.Float64() * screenHeight
	}

	return Enemy{
		X:      x,
		Y:      y,
		Width:  16,
		Height: 16,
		Active: true,
		Speed:  1.0, // Velocidade do inimigo
	}
}

func UpdateEnemies(player *Player, g *Game) {
	for i := range enemies {
		enemies[i].Update(player, g)
	}
}

func (e *Enemy) Update(p *Player, g *Game) {
	if !e.Active {
		return
	}

	e.MoveTowardsPlayer(p, g)
}

func (e *Enemy) MoveTowardsPlayer(p *Player, g *Game) {
	playerX := p.X
	playerY := p.Y

	// Calcule a direção para o jogador
	dx := playerX - e.X
	dy := playerY - e.Y

	// Normalize o vetor para obter direção unitária
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance > 0 {
		dx /= distance
		dy /= distance
	}

	// Aplique um pequeno movimento na direção do jogador
	newX := e.X + dx*e.Speed
	newY := e.Y + dy*e.Speed

	// Verifique colisão e ajuste o movimento, se necessário
	for _, obstacle := range g.Obstacles {
		if CheckCollision(newX, newY, e.Width, e.Height, obstacle.x, obstacle.y, obstacle.width, obstacle.height) {
			// Se houver colisão, mova-se perpendicularmente
			if math.Abs(dx) > math.Abs(dy) {
				// Tente mover-se na direção vertical
				newY = e.Y + e.Speed*math.Copysign(1, dy)
			} else {
				// Tente mover-se na direção horizontal
				newX = e.X + e.Speed*math.Copysign(1, dx)
			}
		}
	}

	// Atualize a posição do inimigo
	e.X = newX
	e.Y = newY
}

func DrawEnemies(screen *ebiten.Image) {
	for _, enemy := range enemies {
		if enemy.Active {
			ebitenutil.DrawRect(screen, enemy.X, enemy.Y, enemy.Width, enemy.Height, color.RGBA{255, 0, 0, 255}) // Vermelho
		}
	}
}

// collision

var lastShotTime time.Time
var shotInterval = 10 * time.Second // Tempo entre disparos

func AutoShoot(player *Player) {
	if time.Since(lastShotTime) >= shotInterval-time.Duration(player.BulletSpeed)*time.Second {
		fmt.Println("shot in", time.Now())
		// Removi o tiro pra testar obstaculos
		// FireBullet(player)        // Adiciona uma nova bala
		lastShotTime = time.Now() // Atualiza o tempo do último disparo
	}
}

func CheckCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	collision := x1 < x2+w2 &&
		x1+w1 > x2 &&
		y1 < y2+h2 &&
		y1+h1 > y2

	if collision {
		fmt.Printf("Colisão detectada: Player (x: %.2f, y: %.2f, w: %.2f, h: %.2f) com PowerUp (x: %.2f, y: %.2f, w: %.2f, h: %.2f)\n", x1, y1, w1, h1, x2, y2, w2, h2)
	}
	return collision
}

func CheckEnemyCollisions(g *Game) bool {
	// Verifica se um inimigo colidiu com o jogador
	for _, enemy := range enemies {
		if enemy.Active && player.X < enemy.X+enemy.Width && player.X+player.Width > enemy.X &&
			player.Y < enemy.Y+enemy.Height && player.Y+player.Height > enemy.Y {
			return true // Colisão entre jogador e inimigo
		}
	}

	// Verifica se uma bala colidiu com um inimigo
	for i := len(bullets) - 1; i >= 0; i-- {
		for j := len(enemies) - 1; j >= 0; j-- {
			if bullets[i].Active && enemies[j].Active &&
				bullets[i].X < enemies[j].X+enemies[j].Width && bullets[i].X+4 > enemies[j].X &&
				bullets[i].Y < enemies[j].Y+enemies[j].Height && bullets[i].Y+10 > enemies[j].Y {
				bullets[i].Active = false // Desativa a bala
				enemies[j].Active = false // Desativa o inimigo
				// Aqui você pode incrementar a pontuação ou fazer outra lógica
				return false
			}
		}
	}

	return false
}

// Função para verificar se o inimigo está próximo de um obstáculo
func isNearObstacle(e *Enemy, o Obstacle) bool {
	// Calcula a distância do centro do inimigo ao centro do obstáculo
	distX := (e.X + e.Width/2) - (o.x + o.width/2)
	distY := (e.Y + e.Height/2) - (o.y + o.height/2)
	distance := math.Sqrt(distX*distX + distY*distY)

	// Define uma distância de proximidade (ajuste conforme necessário)
	return distance < 50
}

// Função para desviar do obstáculo ajustando a direção
func avoidObstacle(e *Enemy, o Obstacle) (float64, float64) {
	// Calcula uma direção perpendicular ao obstáculo
	if e.X < o.x {
		return -0.5, 0 // Move para a esquerda
	} else if e.X > o.x+o.width {
		return 0.5, 0 // Move para a direita
	} else if e.Y < o.y {
		return 0, -0.5 // Move para cima
	} else {
		return 0, 0.5 // Move para baixo
	}
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
			ebitenutil.DrawRect(screen, bullet.X, bullet.Y, 4, 10, color.RGBA{255, 255, 0, 255}) // Amarelo
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
		X:          player.X + player.Width/2 - 2, // Centraliza a bala em relação ao jogador
		Y:          player.Y,
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
			if bullets[i].Y < 0 || bullets[i].Y > screenHeight || bullets[i].X < 0 || bullets[i].X > screenWidth {
				bullets[i].Active = false
			}
		}
	}
}
