package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	mapWidth     = 1080
	mapHeight    = 720
	screenWidth  = 1080
	screenHeight = 720
)

const (
	enemySeparationDistance = 30.0 // Minimum distance between enemies
	separationForce         = 0.5  // Force of separation
	MaxWeapons              = 5
)

type Game struct {
	gameOver        bool
	Obstacles       []Obstacle
	powerUpOptions  []string
	choosingUpgrade bool
	upgradeOptions  []UpgradeOption
}

type UpgradeOption struct {
	Type        string // "skill" or "attribute"
	SkillType   SkillType
	AttrType    AttributeType
	Name        string
	Description string
}

func (g *Game) ShowUpgradeOptions() {
	// Generate 3 random options
	g.upgradeOptions = make([]UpgradeOption, 3)
	g.choosingUpgrade = true

	// Mix skills and attributes
	availableSkills := []SkillType{FireballSkill, IceBlastSkill, ThunderStrike, ShieldSkill, HealingSkill}
	availableAttrs := []AttributeType{Strength, Agility, Vitality, Intelligence, Fortune}

	for i := 0; i < 3; i++ {
		// Randomly choose between skill or attribute
		if rand.Float64() < 0.5 {
			// Skill option
			skill := availableSkills[rand.Intn(len(availableSkills))]
			g.upgradeOptions[i] = UpgradeOption{
				Type:        "skill",
				SkillType:   skill,
				Name:        string(skill),
				Description: getSkillDescription(skill),
			}
		} else {
			// Attribute option
			attr := availableAttrs[rand.Intn(len(availableAttrs))]
			g.upgradeOptions[i] = UpgradeOption{
				Type:        "attribute",
				AttrType:    attr,
				Name:        string(attr),
				Description: getAttributeDescription(attr),
			}
		}
	}
}

func DrawUpgradeChoice(screen *ebiten.Image, options []UpgradeOption) {
	cardWidth := 200.0
	cardHeight := 100.0
	spacing := 20.0

	for i, option := range options {
		cardX := float64(screenWidth)/2 - ((cardWidth*3 + spacing*2) / 2) + float64(i)*(cardWidth+spacing)
		cardY := float64(screenHeight)/2 - cardHeight/2

		// Draw card background
		ebitenutil.DrawRect(screen, cardX, cardY, cardWidth, cardHeight, color.RGBA{50, 50, 50, 255})

		// Draw option info
		text := fmt.Sprintf("%d: %s\n%s", i+1, option.Name, option.Description)
		ebitenutil.DebugPrintAt(screen, text, int(cardX)+10, int(cardY)+10)
	}
}

func getSkillDescription(skillType SkillType) string {
	switch skillType {
	case FireballSkill:
		return "Launch a powerful fireball"
	case IceBlastSkill:
		return "Freeze enemies in place"
	case ThunderStrike:
		return "Call down lightning"
	case ShieldSkill:
		return "Create protective barrier"
	case HealingSkill:
		return "Restore health over time"
	default:
		return ""
	}
}

func (g *Game) applyUpgrade(choice int) {
	option := g.upgradeOptions[choice]
	if option.Type == "skill" {
		player.AddSkill(option.SkillType)
	} else {
		player.AddAttribute(option.AttrType)
	}
}

const PowerUpDuration = 15 * time.Second

// Update ApplyPowerUpEffect
func ApplyPowerUpEffect(player *Player, powerUpType string) {
	expiryTime := time.Now().Add(PowerUpDuration)
	player.ActivePowerUps[powerUpType] = expiryTime

	switch powerUpType {
	case "speed":
		player.Speed += 2.0
	case "power":
		player.WeaponStrength += 3.0
	case "radius":
		player.CollectionRadius += 50.0
	}
}

// Add RemoveExpiredPowerUps function
func (p *Player) RemoveExpiredPowerUps() {
	now := time.Now()
	for powerType, expiryTime := range p.ActivePowerUps {
		if now.After(expiryTime) {
			// Remove effects
			switch powerType {
			case "speed":
				p.Speed -= 2.0
			case "power":
				p.WeaponStrength -= 3.0
			case "radius":
				p.CollectionRadius -= 50.0
			}
			delete(p.ActivePowerUps, powerType)
		}
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

	if g.choosingUpgrade {
		if ebiten.IsKeyPressed(ebiten.Key1) {
			g.applyUpgrade(0)
			g.choosingUpgrade = false
		} else if ebiten.IsKeyPressed(ebiten.Key2) {
			g.applyUpgrade(1)
			g.choosingUpgrade = false
		} else if ebiten.IsKeyPressed(ebiten.Key3) {
			g.applyUpgrade(2)
			g.choosingUpgrade = false
		}
		return nil
	}

	player.Update()
	camera.Update(player.X, player.Y)

	UpdateEnemies(&player, g)

	SpawnPowerUp()
	UpdatePowerUps(&player)

	UpdateBullets()
	SpawnEnemies()
	AutoShoot(&player)
	AutoCastSkills(&player)

	UpdateSkillProjectiles()

	// Verifica colisão com obstáculos para jogador e inimigos
	player.HandleObstacleCollision(g.Obstacles)

	// Verifica colisões
	CheckEnemyCollisions(g)

	UpdateXPItems(&player)

	return nil
}

func DrawStatusUI(screen *ebiten.Image, player *Player) {
	powerupY := 10.0
	skillsY := 120.0
	attributeY := 230.0

	// Update PowerUpInfo struct
	type PowerUpInfo struct {
		Type        string
		ExpireAt    time.Time
		CollectedAt time.Time
	}

	// Convert map to slice with collection times
	var orderedPowerUps []PowerUpInfo
	for t, exp := range player.ActivePowerUps {
		orderedPowerUps = append(orderedPowerUps, PowerUpInfo{
			Type:        t,
			ExpireAt:    exp,
			CollectedAt: exp.Add(-PowerUpDuration), // Calculate collection time
		})
	}

	// Sort by collection time
	sort.Slice(orderedPowerUps, func(i, j int) bool {
		return orderedPowerUps[i].CollectedAt.Before(orderedPowerUps[j].CollectedAt)
	})

	// Draw powerups in collection order
	for _, powerUp := range orderedPowerUps {
		remaining := powerUp.ExpireAt.Sub(time.Now()).Seconds()
		if remaining > 0 {
			var text string
			switch powerUp.Type {
			case "speed":
				text = fmt.Sprintf("Speed Boost: %.1fs", remaining)
			case "power":
				text = fmt.Sprintf("Power Boost: %.1fs", remaining)
			case "radius":
				text = fmt.Sprintf("Radius Boost: %.1fs", remaining)
			}
			ebitenutil.DebugPrintAt(screen, text, 10, int(powerupY))
			powerupY += 20
		}
	}

	// Get ordered skills
	type SkillInfo struct {
		Type  SkillType
		Skill *Skill
	}
	var orderedSkills []SkillInfo
	for t, s := range player.Skills {
		if s.IsActive {
			orderedSkills = append(orderedSkills, SkillInfo{t, s})
		}
	}
	sort.Slice(orderedSkills, func(i, j int) bool {
		return string(orderedSkills[i].Type) < string(orderedSkills[j].Type)
	})

	// Draw skills in order
	ebitenutil.DebugPrintAt(screen, "Skills:", 10, int(skillsY))
	skillsY += 20
	for _, si := range orderedSkills {
		cooldownRemaining := ""
		if time.Since(si.Skill.LastUsed) < time.Duration(si.Skill.Cooldown)*time.Second {
			remaining := si.Skill.Cooldown - time.Since(si.Skill.LastUsed).Seconds()
			cooldownRemaining = fmt.Sprintf(" (CD: %.1fs)", remaining)
		}
		text := fmt.Sprintf("%s Lv.%d%s", si.Type, si.Skill.Level, cooldownRemaining)
		ebitenutil.DebugPrintAt(screen, text, 20, int(skillsY))
		skillsY += 20
	}

	// Get ordered attributes
	type AttrInfo struct {
		Type      AttributeType
		Attribute *Attribute
	}
	var orderedAttrs []AttrInfo
	for t, a := range player.Attributes {
		orderedAttrs = append(orderedAttrs, AttrInfo{t, a})
	}
	sort.Slice(orderedAttrs, func(i, j int) bool {
		return string(orderedAttrs[i].Type) < string(orderedAttrs[j].Type)
	})

	// Draw attributes in order
	ebitenutil.DebugPrintAt(screen, "Attributes:", 10, int(attributeY))
	attributeY += 20
	for _, ai := range orderedAttrs {
		text := fmt.Sprintf("%s Lv.%d (%.1f)", ai.Type, ai.Attribute.Level, ai.Attribute.Value)
		ebitenutil.DebugPrintAt(screen, text, 20, int(attributeY))
		attributeY += 20
	}
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

	if g.choosingUpgrade {
		DrawUpgradeChoice(screen, g.upgradeOptions)
		return
	}

	// Desenha o jogador, inimigos e power-ups
	player.Draw(screen)

	DrawEnemies(screen)
	DrawPowerUps(screen)
	DrawBullets(screen)
	DrawStatusUI(screen, &player)

	DrawSkillProjectiles(screen)

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
	CurrentWeapon    WeaponType
	ActivePowerUps   map[string]time.Time
	Skills           map[SkillType]*Skill
	Attributes       map[AttributeType]*Attribute

	Barrier *BarrierEffect
	Aura    *AuraEffect
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
	p.ApplyAttributeEffects()

	currentGame.choosingUpgrade = true
	currentGame.ShowUpgradeOptions()
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
	p.XP = 0.0             // Adicionar valor inicial de XP
	p.XPToNextLevel = 10.0 // Adicionar valor necessário para o próximo nível
	p.Level = 1            // Adicionar nível inicial
	p.CollectionRadius = 50.0
	p.ActivePowerUps = make(map[string]time.Time)
	p.CurrentWeapon = BasicGun

	p.Attributes = make(map[AttributeType]*Attribute)
	p.ApplyAttributeEffects()

	p.InitSkillSystem()
}

// Add attribute effects
func (p *Player) ApplyAttributeEffects() {
	// Reset to base values first
	p.WeaponStrength = 5.0 // Base damage
	p.Speed = 2.0          // Base speed
	p.Health = 100.0       // Base health
	p.BulletSpeed = 1.0    // Base projectile speed

	// Apply attribute bonuses
	for attrType, attr := range p.Attributes {
		switch attrType {
		case Strength:
			// Increases weapon damage
			p.WeaponStrength += float64(attr.Level) * 2.0

		case Agility:
			// Increases movement and attack speed
			p.Speed += float64(attr.Level) * 0.5
			p.BulletSpeed += float64(attr.Level) * 0.2

		case Vitality:
			// Increases max health
			healthBonus := float64(attr.Level) * 20.0
			p.Health += healthBonus

		case Intelligence:
			// Reduces skill cooldowns
			for _, skill := range p.Skills {
				skill.Cooldown = getBaseCooldown(skill.Type) * (1.0 - float64(attr.Level)*0.1)
			}

		case Fortune:
			// Increases XP gain (implemented in XP collection)
			p.CollectionRadius += float64(attr.Level) * 10.0
		}
	}
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

// Add constants
const (
	PowerUpSpawnInterval = 10 * time.Second
)

// Add last spawn tracking
var (
	lastPowerUpSpawn time.Time
)

func InitPowerUps() {
	powerUps = make([]PowerUp, 0)
	lastPowerUpSpawn = time.Now()
}

// Add SpawnPowerUp function
func SpawnPowerUp() {
	if time.Since(lastPowerUpSpawn) < PowerUpSpawnInterval {
		return
	}

	powerUp := PowerUp{
		X:      rand.Float64() * mapWidth,
		Y:      rand.Float64() * mapHeight,
		Width:  16,
		Height: 16,
		Active: true,
	}

	// Randomly choose power-up type
	switch rand.Intn(3) {
	case 0:
		powerUp.Type = "speed"
	case 1:
		powerUp.Type = "power"
	case 2:
		powerUp.Type = "radius"
	}

	powerUps = append(powerUps, powerUp)
	lastPowerUpSpawn = time.Now()
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
	if time.Since(lastShotTime) >= shotInterval-time.Duration(player.BulletSpeed)*time.Second {
		FireBullet(player)
		lastShotTime = time.Now()
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
	Damage     float64
	Color      color.RGBA
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
			ebitenutil.DrawRect(screen, bullet.X-camera.X, bullet.Y-camera.Y, 4, 10, color.RGBA{255, 255, 0, 255})
		}
	}
}

func getBulletColor(weaponType WeaponType) color.RGBA {
	switch weaponType {
	case Shotgun:
		return color.RGBA{255, 0, 0, 255} // Red
	case RapidFire:
		return color.RGBA{0, 255, 255, 255} // Cyan
	case SpreadShot:
		return color.RGBA{0, 255, 0, 255} // Green
	default:
		return color.RGBA{255, 255, 0, 255} // Yellow (BasicGun)
	}
}

// Add explosion effect struct
type AttributeExplosion struct {
	X, Y      float64
	Type      AttributeType
	Radius    float64
	Duration  float64
	StartTime time.Time
	Active    bool
}

var attributeExplosions []AttributeExplosion

// Add explosion creation function
func CreateAttributeExplosion(x, y float64, attrType AttributeType) {
	explosion := AttributeExplosion{
		X:         x,
		Y:         y,
		Type:      attrType,
		Radius:    50.0,
		Duration:  0.5,
		StartTime: time.Now(),
		Active:    true,
	}
	attributeExplosions = append(attributeExplosions, explosion)
}

// Update explosions and apply effects
func UpdateAttributeExplosions(player *Player) {
	for i := range attributeExplosions {
		if !attributeExplosions[i].Active {
			continue
		}

		// Check duration
		if time.Since(attributeExplosions[i].StartTime).Seconds() > attributeExplosions[i].Duration {
			attributeExplosions[i].Active = false
			continue
		}

		// Apply effects to enemies in range
		for j := range enemies {
			if !enemies[j].Active {
				continue
			}

			// Check if enemy is in explosion radius
			dx := enemies[j].X - attributeExplosions[i].X
			dy := enemies[j].Y - attributeExplosions[i].Y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist <= attributeExplosions[i].Radius {
				// Apply different effects based on attribute type
				switch attributeExplosions[i].Type {
				case Strength:
					// High damage explosion
					enemies[j].Health -= 50.0 * float64(player.Attributes[Strength].Level)
				case Intelligence:
					// Stun effect
					enemies[j].Speed *= 0.5
				case Agility:
					// Knockback effect
					force := 10.0 * float64(player.Attributes[Agility].Level)
					if dist > 0 {
						enemies[j].X += (dx / dist) * force
						enemies[j].Y += (dy / dist) * force
					}
				case Vitality:
					// Area denial effect
					enemies[j].Health -= 20.0 * float64(player.Attributes[Vitality].Level)
				case Fortune:
					// Increase XP drop chance
					if enemies[j].Health <= 0 {
						for k := 0; k < player.Attributes[Fortune].Level; k++ {
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
				}
			}
		}
	}
}

func DrawAttributeExplosions(screen *ebiten.Image) {
	for _, explosion := range attributeExplosions {
		if !explosion.Active {
			continue
		}

		color := color.RGBA{255, 0, 255, 128} // Purple
		// switch explosion.Type {
		// case Strength:
		// 	color = color.RGBA{255, 0, 0, 128} // Red
		// case Intelligence:
		// 	color = color.RGBA{0, 0, 255, 128} // Blue
		// case Agility:
		// 	color = color.RGBA{0, 255, 0, 128} // Green
		// case Vitality:
		// 	color = color.RGBA{255, 255, 0, 128} // Yellow
		// case Fortune:

		// }

		// Draw explosion circle
		ebitenutil.DrawCircle(screen,
			explosion.X-camera.X,
			explosion.Y-camera.Y,
			explosion.Radius,
			color)
	}
}

func FireBullet(player *Player) {
	nearestEnemy := GetNearestEnemy(player)
	if nearestEnemy == nil {
		return
	}

	// Calculate direction
	directionX := nearestEnemy.X - player.X
	directionY := nearestEnemy.Y - player.Y
	length := math.Sqrt(directionX*directionX + directionY*directionY)

	// Normalize direction
	if length > 0 {
		directionX /= length
		directionY /= length
	}

	bullet := Bullet{
		X:          player.X + player.Width/2,
		Y:          player.Y + player.Height/2,
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

// Add weapon types
type WeaponType string

const (
	BasicGun   WeaponType = "basic"
	Shotgun    WeaponType = "shotgun"
	RapidFire  WeaponType = "rapid"
	SpreadShot WeaponType = "spread"
)

// Weapon struct
type Weapon struct {
	Type            WeaponType
	FireRate        float64
	ProjectileCount int
	Spread          float64
	X, Y            float64
	Width           float64
	Height          float64
	Active          bool
}

// Add weapon spawning
var weapons []Weapon

func getWeaponStats(weaponType WeaponType) Weapon {
	switch weaponType {
	case Shotgun:
		return Weapon{
			Type:            Shotgun,
			FireRate:        0.8,
			ProjectileCount: 5,
			Spread:          0.3,
		}
	case RapidFire:
		return Weapon{Type: RapidFire, FireRate: 0.2, ProjectileCount: 1, Spread: 0.1}
	case SpreadShot:
		return Weapon{Type: SpreadShot, FireRate: 0.5, ProjectileCount: 3, Spread: 0.2}
	default:
		return Weapon{Type: BasicGun, FireRate: 1.0, ProjectileCount: 1, Spread: 0}
	}
}

// Constants
const (
	MaxSkills     = 5
	MaxAttributes = 5
	MaxLevel      = 5
)

// Skill types
type SkillType string

const (
	FireballSkill SkillType = "fireball"
	IceBlastSkill SkillType = "iceblast"
	ThunderStrike SkillType = "thunder"
	ShieldSkill   SkillType = "shield"
	HealingSkill  SkillType = "healing"
)

// Attribute types
type AttributeType string

const (
	Strength     AttributeType = "strength"     // Increases damage
	Agility      AttributeType = "agility"      // Increases speed
	Vitality     AttributeType = "vitality"     // Increases health
	Intelligence AttributeType = "intelligence" // Increases skill power
	Fortune      AttributeType = "fortune"      // Increases XP gain
)

// Skill definition
type Skill struct {
	Type     SkillType
	Level    int
	Cooldown float64
	Power    float64
	LastUsed time.Time
	IsActive bool
}

// Attribute definition
type Attribute struct {
	Type  AttributeType
	Level int
	Value float64
}

// Initialize player skills/attributes
func (p *Player) InitSkillSystem() {
	p.Skills = make(map[SkillType]*Skill)
	p.Attributes = make(map[AttributeType]*Attribute)
}

// Add or upgrade skill
func (p *Player) AddSkill(skillType SkillType) bool {
	// Check max skills limit
	if len(p.Skills) >= MaxSkills && p.Skills[skillType] == nil {
		return false
	}

	// Get or create skill
	skill, exists := p.Skills[skillType]
	if !exists {
		skill = &Skill{
			Type:     skillType,
			Level:    1,
			Cooldown: getBaseCooldown(skillType),
			Power:    getBasePower(skillType),
			IsActive: true,
		}
		p.Skills[skillType] = skill
	} else if skill.Level < MaxLevel {
		// Level up existing skill
		skill.Level++
		skill.Power *= 1.5
		skill.Cooldown *= 0.8
	}

	return true
}

// Add or upgrade attribute
func (p *Player) AddAttribute(attrType AttributeType) bool {
	// Check max attributes limit
	if len(p.Attributes) >= MaxAttributes && p.Attributes[attrType] == nil {
		return false
	}

	// Get or create attribute
	attr, exists := p.Attributes[attrType]
	if !exists {
		attr = &Attribute{
			Type:  attrType,
			Level: 1,
			Value: getBaseValue(attrType),
		}
		p.Attributes[attrType] = attr
	} else if attr.Level < MaxLevel {
		// Level up existing attribute
		attr.Level++
		attr.Value *= 1.5
	}

	p.ApplyAttributeEffects()
	return true
}

// Helper functions
func getBaseCooldown(skillType SkillType) float64 {
	switch skillType {
	case FireballSkill:
		return 5.0
	case IceBlastSkill:
		return 8.0
	case ThunderStrike:
		return 10.0
	case ShieldSkill:
		return 15.0
	case HealingSkill:
		return 20.0
	default:
		return 1.0
	}
}

func getBasePower(skillType SkillType) float64 {
	switch skillType {
	case FireballSkill:
		return 50.0
	case IceBlastSkill:
		return 40.0
	case ThunderStrike:
		return 70.0
	case ShieldSkill:
		return 100.0
	case HealingSkill:
		return 30.0
	default:
		return 10.0
	}
}

func getBaseValue(attrType AttributeType) float64 {
	switch attrType {
	case Strength:
		return 10.0
	case Agility:
		return 10.0
	case Vitality:
		return 100.0
	case Intelligence:
		return 10.0
	case Fortune:
		return 1.0
	default:
		return 1.0
	}
}

func getAttributeDescription(attrType AttributeType) string {
	switch attrType {
	case Strength:
		return "Increase damage dealt"
	case Agility:
		return "Move faster"
	case Vitality:
		return "Increase max health"
	case Intelligence:
		return "Boost skill power"
	case Fortune:
		return "Get more XP"
	default:
		return ""
	}
}

// Add skill effect types and properties
type SkillEffect struct {
	Type     SkillType
	Duration float64
	Range    float64
	Speed    float64
	Pattern  string // "circle", "line", "cone", etc.
	Color    color.RGBA
	Damage   float64
}

func getSkillEffect(skillType SkillType) SkillEffect {
	switch skillType {
	case FireballSkill:
		return SkillEffect{
			Type:     FireballSkill,
			Duration: 2.0,
			Range:    200.0,
			Speed:    8.0,
			Pattern:  "projectile",
			Color:    color.RGBA{255, 100, 0, 255}, // Orange fire
			Damage:   50.0,
		}
	case IceBlastSkill:
		return SkillEffect{
			Type:     IceBlastSkill,
			Duration: 3.0,
			Range:    150.0,
			Speed:    4.0,
			Pattern:  "cone",
			Color:    color.RGBA{100, 200, 255, 255}, // Light blue
			Damage:   30.0,
		}
	case ThunderStrike:
		return SkillEffect{
			Type:     ThunderStrike,
			Duration: 0.5,
			Range:    300.0,
			Speed:    15.0,
			Pattern:  "chain",
			Color:    color.RGBA{255, 255, 0, 255}, // Yellow
			Damage:   70.0,
		}
	case ShieldSkill:
		return SkillEffect{
			Type:     ShieldSkill,
			Duration: 5.0,
			Range:    50.0,
			Speed:    0.0,
			Pattern:  "circle",
			Color:    color.RGBA{100, 100, 255, 255}, // Blue shield
			Damage:   0.0,
		}
	case HealingSkill:
		return SkillEffect{
			Type:     HealingSkill,
			Duration: 4.0,
			Range:    0.0,
			Speed:    0.0,
			Pattern:  "aura",
			Color:    color.RGBA{0, 255, 0, 255}, // Green
			Damage:   -30.0,                      // Negative damage = healing
		}
	default:
		return SkillEffect{}
	}
}

// Add skill execution function
func ExecuteSkill(player *Player, skillType SkillType) {
	skill := player.Skills[skillType]
	if skill == nil || !skill.IsActive {
		return
	}

	if time.Since(skill.LastUsed) < time.Duration(skill.Cooldown)*time.Second {
		return
	}

	effect := getSkillEffect(skillType)

	switch effect.Pattern {
	case "projectile":
		// Single direction projectile (Fireball)
		nearestEnemy := GetNearestEnemy(player)
		if nearestEnemy != nil {
			FireProjectile(player, nearestEnemy, effect)
		}

	case "cone":
		// Ice Blast: Multiple projectiles in cone shape
		FireConeProjectiles(player, effect, 45.0, 5) // 45 degree spread, 5 projectiles

	case "chain":
		// Thunder: Chain lightning between enemies
		ChainLightning(player, effect, 3) // Chain to 3 enemies

	case "circle":
		// Shield: Circular barrier around player
		CreateBarrier(player, effect)

	case "aura":
		// Healing: Area effect centered on player
		CreateAuraEffect(player, effect)
	}

	skill.LastUsed = time.Now()
}

// Example implementation of attack patterns
func FireProjectile(player *Player, target *Enemy, effect SkillEffect) {
	dx := target.X - player.X
	dy := target.Y - player.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > 0 {
		dx = dx / dist * effect.Speed
		dy = dy / dist * effect.Speed
	}

	// Create projectile
	bullet := Bullet{
		X:          player.X,
		Y:          player.Y,
		Speed:      effect.Speed,
		DirectionX: dx,
		DirectionY: dy,
		Active:     true,
		Damage:     effect.Damage,
		Color:      effect.Color,
	}
	bullets = append(bullets, bullet)
}

func FireConeProjectiles(player *Player, effect SkillEffect, spreadAngle float64, count int) {
	baseAngle := math.Atan2(player.Y, player.X)
	angleStep := spreadAngle / float64(count-1)

	for i := 0; i < count; i++ {
		angle := baseAngle - (spreadAngle / 2) + (float64(i) * angleStep)
		dx := math.Cos(angle) * effect.Speed
		dy := math.Sin(angle) * effect.Speed

		bullet := Bullet{
			X:          player.X,
			Y:          player.Y,
			Speed:      effect.Speed,
			DirectionX: dx,
			DirectionY: dy,
			Active:     true,
			Damage:     effect.Damage,
			Color:      effect.Color,
		}
		bullets = append(bullets, bullet)
	}
}

func ChainLightning(player *Player, effect SkillEffect, chainCount int) {
	// Find nearest enemy
	current := GetNearestEnemy(player)
	if current == nil {
		return
	}

	hit := make(map[*Enemy]bool)
	hit[current] = true

	for i := 0; i < chainCount && current != nil; i++ {
		// Apply damage to current target
		current.Health -= effect.Damage

		// Find next closest enemy that hasn't been hit
		var next *Enemy
		minDist := math.MaxFloat64

		for _, enemy := range enemies {
			if enemy.Active && !hit[&enemy] {
				dist := math.Sqrt(math.Pow(enemy.X-current.X, 2) + math.Pow(enemy.Y-current.Y, 2))
				if dist < minDist {
					minDist = dist
					next = &enemy
				}
			}
		}

		current = next
		if current != nil {
			hit[current] = true
		}
	}
}

// Add to existing structs
type BarrierEffect struct {
	Active    bool
	Duration  float64
	Radius    float64
	StartTime time.Time
	Color     color.RGBA
}

type AuraEffect struct {
	Active     bool
	Duration   float64
	Radius     float64
	StartTime  time.Time
	Color      color.RGBA
	HealAmount float64
}

func CreateBarrier(player *Player, effect SkillEffect) {
	player.Barrier = &BarrierEffect{
		Active:    true,
		Duration:  effect.Duration,
		Radius:    effect.Range,
		StartTime: time.Now(),
		Color:     effect.Color,
	}
}

func CreateAuraEffect(player *Player, effect SkillEffect) {
	player.Aura = &AuraEffect{
		Active:     true,
		Duration:   effect.Duration,
		Radius:     effect.Range,
		StartTime:  time.Now(),
		Color:      effect.Color,
		HealAmount: -effect.Damage, // Convert damage to healing
	}
}

// Add skill projectile types
type SkillProjectile struct {
	X, Y       float64
	DirectionX float64
	DirectionY float64
	Speed      float64
	Active     bool
	SkillType  SkillType
	Power      float64
	Radius     float64
}

var skillProjectiles []SkillProjectile

// Cast skill function
func CastSkill(player *Player, skillType SkillType) {
	skill := player.Skills[skillType]
	if skill == nil || !skill.IsActive {
		return
	}

	if time.Since(skill.LastUsed) < time.Duration(skill.Cooldown)*time.Second {
		return
	}

	nearestEnemy := GetNearestEnemy(player)
	if nearestEnemy == nil {
		return
	}

	switch skillType {
	case FireballSkill:
		// Single powerful projectile
		dx := nearestEnemy.X - player.X
		dy := nearestEnemy.Y - player.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > 0 {
			dx = dx / dist
			dy = dy / dist
		}

		projectile := SkillProjectile{
			X:          player.X + player.Width/2,
			Y:          player.Y + player.Height/2,
			DirectionX: dx,
			DirectionY: dy,
			Speed:      8.0,
			Active:     true,
			SkillType:  FireballSkill,
			Power:      skill.Power,
			Radius:     20.0,
		}
		skillProjectiles = append(skillProjectiles, projectile)

	case ThunderStrike:
		// Chain lightning effect
		projectile := SkillProjectile{
			X:         nearestEnemy.X,
			Y:         nearestEnemy.Y,
			Active:    true,
			SkillType: ThunderStrike,
			Power:     skill.Power,
			Radius:    50.0,
		}
		skillProjectiles = append(skillProjectiles, projectile)
	}

	skill.LastUsed = time.Now()
}

// Update skill projectiles
func UpdateSkillProjectiles() {
	for i := range skillProjectiles {
		if !skillProjectiles[i].Active {
			continue
		}

		switch skillProjectiles[i].SkillType {
		case FireballSkill:
			// Move fireball
			skillProjectiles[i].X += skillProjectiles[i].DirectionX * skillProjectiles[i].Speed
			skillProjectiles[i].Y += skillProjectiles[i].DirectionY * skillProjectiles[i].Speed

			// Check if out of bounds
			if skillProjectiles[i].X < 0 || skillProjectiles[i].X > mapWidth ||
				skillProjectiles[i].Y < 0 || skillProjectiles[i].Y > mapHeight {
				skillProjectiles[i].Active = false
			}

		case ThunderStrike:
			// Thunder strike is instant, deactivate after one frame
			skillProjectiles[i].Active = false
		}

		// Check collision with enemies
		for _, enemy := range enemies {
			if !enemy.Active {
				continue
			}

			dx := enemy.X - skillProjectiles[i].X
			dy := enemy.Y - skillProjectiles[i].Y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist <= skillProjectiles[i].Radius {
				switch skillProjectiles[i].SkillType {
				case FireballSkill:
					enemy.Health -= skillProjectiles[i].Power
					skillProjectiles[i].Active = false
				case ThunderStrike:
					enemy.Health -= skillProjectiles[i].Power
				}

				if enemy.Health <= 0 {
					enemy.Active = false
					// Drop XP
					xpItem := XPItem{
						X:      enemy.X,
						Y:      enemy.Y,
						Width:  10,
						Height: 10,
						Active: true,
					}
					xpItems = append(xpItems, xpItem)
				}
			}
		}
	}
}

// Draw skill projectiles
func DrawSkillProjectiles(screen *ebiten.Image) {
	for _, proj := range skillProjectiles {
		if !proj.Active {
			continue
		}

		switch proj.SkillType {
		case FireballSkill:
			ebitenutil.DrawCircle(screen,
				proj.X-camera.X,
				proj.Y-camera.Y,
				proj.Radius,
				color.RGBA{255, 100, 0, 255})
		case ThunderStrike:
			ebitenutil.DrawCircle(screen,
				proj.X-camera.X,
				proj.Y-camera.Y,
				proj.Radius,
				color.RGBA{255, 255, 0, 255})
		}
	}
}

var (
	skillLastCastTime = make(map[SkillType]time.Time)
)

func AutoCastSkills(player *Player) {
	// Check each skill
	for skillType, skill := range player.Skills {
		if !skill.IsActive {
			continue
		}

		lastCast, exists := skillLastCastTime[skillType]
		if !exists {
			skillLastCastTime[skillType] = time.Now()
			continue
		}

		// Check if cooldown has passed
		if time.Since(lastCast) >= time.Duration(skill.Cooldown)*time.Second {
			nearestEnemy := GetNearestEnemy(player)
			if nearestEnemy != nil {
				CastSkill(player, skillType)
				skillLastCastTime[skillType] = time.Now()
			}
		}
	}
}
