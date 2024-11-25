package stats

import (
	"fmt"
	"game/internal/core"
	"game/internal/plugins/playing/player"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type StatsPlugin struct {
	kernel       *core.GameKernel
	playerPlugin *player.PlayerPlugin
	gameFont     font.Face
}

func NewStatsPlugin(playerPlugin *player.PlayerPlugin) *StatsPlugin {
	return &StatsPlugin{
		playerPlugin: playerPlugin,
	}
}

func (sp *StatsPlugin) ID() string {
	return "StatsSystem"
}

func (sp *StatsPlugin) Init(kernel *core.GameKernel) error {
	sp.kernel = kernel

	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	sp.gameFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (sp *StatsPlugin) Update() error {
	// No update logic needed for stats
	return nil
}

func (sp *StatsPlugin) Draw(screen *ebiten.Image) {
	playerPower := sp.playerPlugin.GetHealth()
	text.Draw(screen, fmt.Sprintf("Life: %.0f", playerPower), sp.gameFont, 10, 30, color.White)
}
