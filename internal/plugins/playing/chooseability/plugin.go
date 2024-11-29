package chooseability

import (
	"game/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ChooseAbilityPlugin struct {
	kernel *core.GameKernel
}

func NewChooseAbilityPlugin(plugins *core.PluginManager) *ChooseAbilityPlugin {
	return &ChooseAbilityPlugin{}
}

func (cp *ChooseAbilityPlugin) ID() string {
	return "ChooseAbilityPlugin"
}

func (cp *ChooseAbilityPlugin) Init(kernel *core.GameKernel) error {
	cp.kernel = kernel
	return nil
}

func (cp *ChooseAbilityPlugin) Update() error {
	return nil
}

func (cp *ChooseAbilityPlugin) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 320, 0, 320, 240, color.RGBA{0, 230, 230, 0xff})
}
