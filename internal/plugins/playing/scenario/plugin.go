package scenario

import (
	"game/internal/core"
	"image/color"
	"math/rand"
	"sync"
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

type MapTile struct {
	Type     TileType
	Walkable bool
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

	s.generateMap()

	return s
}

func (sp *ScenarioPlugin) Init(kernel *core.GameKernel) error {
	sp.kernel = kernel

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
			if r < 0.15 {
				sp.tiles[x][y].Type = TileTree
				sp.tiles[x][y].Walkable = false
			} else if r < 0.25 {
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
	tileSize := 32
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	// Number of goroutines to use
	numWorkers := 1
	rowsPerWorker := len(sp.tiles) / numWorkers

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if w == numWorkers-1 {
			endRow = len(sp.tiles) // Handle remaining rows
		}

		go func(start, end int) {
			defer wg.Done()

			// Pre-create images for this worker
			tileImages := make(map[TileType]*ebiten.Image)
			for tileType := range tileColors {
				img := ebiten.NewImage(tileSize, tileSize)
				img.Fill(tileColors[tileType])
				tileImages[tileType] = img
			}

			for x := start; x < end; x++ {
				for y := range sp.tiles[x] {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))

					tileType := sp.tiles[x][y].Type

					mutex.Lock()
					screen.DrawImage(tileImages[tileType], op)
					mutex.Unlock()
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()
}

func (sp *ScenarioPlugin) Update() error {
	return nil
}
