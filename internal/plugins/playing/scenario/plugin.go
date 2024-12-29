package scenario

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/playing/camera"
	"image/color"
	"math/rand"

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

type MapTile struct {
	Type     TileType
	Walkable bool
	Animated *assets.Animation
}

type Chunk struct {
	Tiles     [][]*MapTile
	Generated bool
}

type ScenarioPlugin struct {
	kernel    *core.GameKernel
	plugins   *core.PluginManager
	chunks    map[int]map[int]*Chunk
	chunkSize int
	tileSize  int
}

func New(plugins *core.PluginManager) *ScenarioPlugin {
	return &ScenarioPlugin{
		plugins:   plugins,
		chunks:    make(map[int]map[int]*Chunk),
		chunkSize: 16,
		tileSize:  32,
	}
}

func (sp *ScenarioPlugin) Init(kernel *core.GameKernel) error {
	sp.kernel = kernel
	return nil
}

func (sp *ScenarioPlugin) ID() string {
	return "ScenarioSystem"
}

func (sp *ScenarioPlugin) generateChunk(chunkX, chunkY int) *Chunk {
	chunk := &Chunk{
		Tiles: make([][]*MapTile, sp.chunkSize),
	}

	for x := range chunk.Tiles {
		chunk.Tiles[x] = make([]*MapTile, sp.chunkSize)
		for y := range chunk.Tiles[x] {
			tile := &MapTile{
				Type:     TileGround,
				Walkable: true,
			}

			// Generate random obstacles
			r := rand.Float64()
			if r < 0.05 {
				tile.Type = TileTree
				tile.Walkable = false
			} else if r < 0.10 {
				tile.Type = TileRock
				tile.Walkable = false
			}

			// Load animation
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
				if err == nil {
					tile.Animated = animation
				}
			}

			chunk.Tiles[x][y] = tile
		}
	}

	chunk.Generated = true
	return chunk
}

func (sp *ScenarioPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := sp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	if cameraPlugin == nil {
		return
	}
	cameraX, cameraY := cameraPlugin.GetPosition()

	// Calculate visible chunks
	startChunkX := int(cameraX) / (sp.chunkSize * sp.tileSize)
	startChunkY := int(cameraY) / (sp.chunkSize * sp.tileSize)

	chunksX := (constants.ScreenWidth / (sp.chunkSize * sp.tileSize)) + 2
	chunksY := (constants.ScreenHeight / (sp.chunkSize * sp.tileSize)) + 2

	// Generate and draw visible chunks
	for cx := startChunkX - 1; cx <= startChunkX+chunksX; cx++ {
		if sp.chunks[cx] == nil {
			sp.chunks[cx] = make(map[int]*Chunk)
		}

		for cy := startChunkY - 1; cy <= startChunkY+chunksY; cy++ {
			if sp.chunks[cx][cy] == nil {
				sp.chunks[cx][cy] = sp.generateChunk(cx, cy)
			}

			chunk := sp.chunks[cx][cy]
			for x := 0; x < sp.chunkSize; x++ {
				for y := 0; y < sp.chunkSize; y++ {
					worldX := cx*sp.chunkSize*sp.tileSize + x*sp.tileSize
					worldY := cy*sp.chunkSize*sp.tileSize + y*sp.tileSize
					screenX := float64(worldX) - cameraX
					screenY := float64(worldY) - cameraY

					tile := chunk.Tiles[x][y]
					if tile.Animated != nil {
						tile.Animated.Draw(screen, assets.DrawInput{
							Width:  float64(sp.tileSize),
							Height: float64(sp.tileSize),
							X:      screenX,
							Y:      screenY,
						})
					} else {
						// Fallback to color rendering
						tileImage := ebiten.NewImage(sp.tileSize, sp.tileSize)
						tileImage.Fill(tileColors[tile.Type])
						op := &ebiten.DrawImageOptions{}
						op.GeoM.Translate(screenX, screenY)
						screen.DrawImage(tileImage, op)
					}
				}
			}
		}
	}

	go sp.cleanupFarChunks(startChunkX, startChunkY)
}

func (sp *ScenarioPlugin) cleanupFarChunks(centerX, centerY int) {
	const cleanupDistance = 4

	for cx := range sp.chunks {
		if abs(cx-centerX) > cleanupDistance {
			delete(sp.chunks, cx)
			continue
		}
		for cy := range sp.chunks[cx] {
			if abs(cy-centerY) > cleanupDistance {
				delete(sp.chunks[cx], cy)
			}
		}
	}
}

func (sp *ScenarioPlugin) Update() error {
	// Update animations for visible chunks
	cameraPlugin := sp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	if cameraPlugin == nil {
		return nil
	}

	cameraX, cameraY := cameraPlugin.GetPosition()
	startChunkX := int(cameraX) / (sp.chunkSize * sp.tileSize)
	startChunkY := int(cameraY) / (sp.chunkSize * sp.tileSize)

	for cx := range sp.chunks {
		for cy, chunk := range sp.chunks[cx] {
			if abs(cx-startChunkX) <= 2 && abs(cy-startChunkY) <= 2 {
				for x := range chunk.Tiles {
					for y := range chunk.Tiles[x] {
						if chunk.Tiles[x][y].Animated != nil {
							chunk.Tiles[x][y].Animated.Update(sp.kernel.DeltaTime)
						}
					}
				}
			}
		}
	}

	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (sp *ScenarioPlugin) IsTileWalkable(worldX, worldY float64) bool {
	chunkX := int(worldX) / (sp.chunkSize * sp.tileSize)
	chunkY := int(worldY) / (sp.chunkSize * sp.tileSize)

	if sp.chunks[chunkX] == nil || sp.chunks[chunkX][chunkY] == nil {
		return false
	}

	tileX := (int(worldX) / sp.tileSize) % sp.chunkSize
	tileY := (int(worldY) / sp.tileSize) % sp.chunkSize

	return sp.chunks[chunkX][chunkY].Tiles[tileX][tileY].Walkable
}
