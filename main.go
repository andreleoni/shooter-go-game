package main

import (
	"fmt"
	"image/color"
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
	gameOver bool
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil // Não atualiza se o jogo estiver terminado
	}

	player.Update()
	UpdateEnemies(&player)
	UpdatePowerUps(&player)
	UpdateBullets()
	SpawnEnemies()
	AutoShoot(&player)

	// Verifica colisões
	if CheckCollisions() {
		g.gameOver = true
	}

	return nil
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	player.Init()
	InitEnemies()
	InitPowerUps()

	CheckCollisions()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Jogo Básico em Ebiten")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
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
	X, Y   float64
	Speed  float64
	Width  float64
	Height float64
}

var player Player

func (p *Player) Init() {
	p.X = screenWidth / 2
	p.Y = screenHeight / 2
	p.Speed = 2.0
	p.Width = 16
	p.Height = 16
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && p.X > 0 {
		p.X -= p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && p.X < screenWidth-p.Width {
		p.X += p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && p.Y > 0 {
		p.Y -= p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && p.Y < screenHeight-p.Height {
		p.Y += p.Speed
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		FireBullet(p)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.X, p.Y, p.Width, p.Height, color.RGBA{0, 255, 0, 255}) // Verde
}

// powerup

type PowerUp struct {
	X, Y   float64
	Width  float64
	Height float64
	Active bool
}

var powerUps []PowerUp

func InitPowerUps() {
	for i := 0; i < 3; i++ {
		powerUp := PowerUp{
			X:      rand.Float64() * screenWidth,
			Y:      rand.Float64() * screenHeight,
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
			ApplyPowerUpEffect(player)
		}
	}
}

func DrawPowerUps(screen *ebiten.Image) {
	for _, powerUp := range powerUps {
		if powerUp.Active {
			ebitenutil.DrawRect(screen, powerUp.X, powerUp.Y, powerUp.Width, powerUp.Height, color.RGBA{0, 0, 255, 255}) // Azul
		}
	}
}

func ApplyPowerUpEffect(player *Player) {
	player.Speed += 0.5 // Exemplo: aumenta a velocidade do jogador temporariamente
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
var spawnInterval = 200 * time.Millisecond // Tempo entre spawns

func SpawnEnemies() {
	if time.Since(lastSpawnTime) < spawnInterval {
		return
	}

	// Limite máximo de inimigos ativos
	if len(enemies) < 10 { // Por exemplo, 10 inimigos ativos
		enemies = append(enemies, NewEnemy())
	}
	lastSpawnTime = time.Now()
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

func UpdateEnemies(player *Player) {
	for i := range enemies {
		enemies[i].Update(player)
	}
}

func (e *Enemy) Update(player *Player) {
	if !e.Active {
		return
	}

	// Calcula a direção em que o inimigo deve se mover
	directionX := player.X - e.X
	directionY := player.Y - e.Y
	length := math.Sqrt(directionX*directionX + directionY*directionY)

	// Normaliza a direção
	if length > 0 {
		directionX /= length
		directionY /= length
	}

	// Atualiza a posição do inimigo
	e.X += directionX * e.Speed
	e.Y += directionY * e.Speed
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
var shotInterval = 2 * time.Second // Tempo entre disparos

func AutoShoot(player *Player) {
	if time.Since(lastShotTime) >= shotInterval {
		FireBullet(player)        // Adiciona uma nova bala
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

func CheckCollisions() bool {
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
