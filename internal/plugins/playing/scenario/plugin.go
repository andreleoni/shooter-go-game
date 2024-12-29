package scenario

import (
	"fmt"
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/playing/camera"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type TileType int

const (
	TileGround TileType = iota
	TileTree
	TileRock
	TilePortal
)

var tileColors = map[TileType]color.Color{
	TileGround: color.RGBA{34, 139, 34, 255},   // Green
	TileTree:   color.RGBA{0, 100, 0, 255},     // Dark Green
	TileRock:   color.RGBA{128, 128, 128, 255}, // Gray
	TilePortal: color.RGBA{138, 43, 226, 255},  // Purple
}

var tileSize = 32

type MapTile struct {
	Type     TileType
	Walkable bool
	Animated *assets.Animation
}

type ScenarioPlugin struct {
	kernel    *core.GameKernel
	plugins   *core.PluginManager
	tiles     [][]*MapTile
	mapWidth  int
	mapHeight int
}

func New(plugins *core.PluginManager) *ScenarioPlugin {
	s := &ScenarioPlugin{
		plugins:   plugins,
		mapWidth:  50, // Reduced map size
		mapHeight: 50,
	}

	return s
}

func (sp *ScenarioPlugin) Init(kernel *core.GameKernel) error {
	sp.kernel = kernel

	// Initialize tiles array
	sp.tiles = make([][]*MapTile, sp.mapWidth)
	for i := range sp.tiles {
		sp.tiles[i] = make([]*MapTile, sp.mapHeight)
		for j := range sp.tiles[i] {
			// Initialize with ground tile first
			sp.tiles[i][j] = &MapTile{
				Type:     TileGround,
				Walkable: true,
			}
		}
	}

	// Generate map features before loading animations
	sp.generateMap()

	// Load animations for tiles after map generation
	for i := range sp.tiles {
		for j := range sp.tiles[i] {
			tile := sp.tiles[i][j]

			// Select appropriate asset based on tile type
			var imagePath string
			switch tile.Type {
			case TileGround:
				imagePath = "assets/images/maps/grass/1/asset"
			case TileTree:
				imagePath = "assets/images/maps/grass/2/asset"
			case TileRock:
				imagePath = "assets/images/maps/grass/3/asset"
			case TilePortal:
				imagePath = "assets/images/maps/grass/4/asset"
			}

			if imagePath != "" {
				animation := assets.NewAnimation(0.1)
				err := animation.LoadFromJSON(
					imagePath+".json",
					imagePath+".png")

				if err != nil {
					// Log error but continue loading other tiles
					fmt.Printf("Warning: failed to load asset at %d,%d: %v\n", i, j, err)
					continue
				}

				tile.Animated = animation
			}
		}
	}

	return nil
}

func (sp *ScenarioPlugin) ID() string {
	return "ScenarioPlugin"
}

func (sp *ScenarioPlugin) generateMap() {
	rand.Seed(time.Now().UnixNano())

	// Initialize map with ground tiles
	sp.tiles = make([][]*MapTile, sp.mapWidth)

	for x := range sp.tiles {
		sp.tiles[x] = make([]*MapTile, sp.mapHeight)
		for y := range sp.tiles[x] {
			sp.tiles[x][y] = &MapTile{
				Type:     TileGround,
				Walkable: true,
			}
		}
	}

	// Add random obstacles (15% chance for trees, 10% for rocks)
	for x := 0; x < sp.mapWidth; x++ {
		for y := 0; y < sp.mapHeight; y++ {
			r := rand.Float64()
			if r < 0.10 {
				sp.tiles[x][y].Type = TileTree
				sp.tiles[x][y].Walkable = false
			} else if r < 0.15 {
				sp.tiles[x][y].Type = TileRock
				sp.tiles[x][y].Walkable = false
			}
		}
	}

	// Add random portal
	portalX := rand.Intn(sp.mapWidth)
	portalY := rand.Intn(sp.mapHeight)
	sp.tiles[portalX][portalY].Type = TilePortal
	sp.tiles[portalX][portalY].Walkable = true

	// Clear area around portal
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			newX := portalX + x
			newY := portalY + y
			if newX >= 0 && newX < sp.mapWidth && newY >= 0 && newY < sp.mapHeight {
				sp.tiles[newX][newY].Type = TileGround
				sp.tiles[newX][newY].Walkable = true
			}
		}
	}
}

func (sp *ScenarioPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := sp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	if cameraPlugin == nil {
		return
	}
	cameraX, cameraY := cameraPlugin.GetPosition()

	// Calculate visible tile range
	startTileX := int(cameraX) / tileSize
	startTileY := int(cameraY) / tileSize

	tilesX := constants.ScreenWidth/tileSize + 2
	tilesY := constants.ScreenHeight/tileSize + 2

	endTileX := startTileX + tilesX
	endTileY := startTileY + tilesY

	startTileX = clamp(startTileX, 0, sp.mapWidth-1)
	startTileY = clamp(startTileY, 0, sp.mapHeight-1)
	endTileX = clamp(endTileX, 0, sp.mapWidth)
	endTileY = clamp(endTileY, 0, sp.mapHeight)

	for x := startTileX; x < endTileX; x++ {
		for y := startTileY; y < endTileY; y++ {
			tile := sp.tiles[x][y]
			if tile == nil || tile.Animated == nil {
				// Fallback to color if animation is not available
				tileImage := ebiten.NewImage(tileSize, tileSize)
				tileImage.Fill(tileColors[tile.Type])

				op := &ebiten.DrawImageOptions{}
				screenX := float64(x*tileSize) - cameraX
				screenY := float64(y*tileSize) - cameraY
				op.GeoM.Translate(screenX, screenY)
				screen.DrawImage(tileImage, op)
				continue
			}

			// Draw animated tile
			screenX := float64(x*tileSize) - cameraX
			screenY := float64(y*tileSize) - cameraY

			tile.Animated.Draw(screen, assets.DrawInput{
				Width:  float64(tileSize),
				Height: float64(tileSize),
				X:      screenX,
				Y:      screenY,
			})
		}
	}
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (sp *ScenarioPlugin) Update() error {
	// Update animations
	if sp.tiles != nil {
		for x := range sp.tiles {
			for y := range sp.tiles[x] {
				if sp.tiles[x][y] != nil && sp.tiles[x][y].Animated != nil {
					sp.tiles[x][y].Animated.Update(sp.kernel.DeltaTime)
				}
			}
		}
	}
	return nil
}
