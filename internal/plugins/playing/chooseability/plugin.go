package chooseability

import (
	"game/internal/core"
	"game/internal/plugins/menu/fontface"
	"game/internal/plugins/playing/weapon/entities"
	"game/internal/plugins/playing/weapon/templates"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type ChooseAbilityPlugin struct {
	kernel             *core.GameKernel
	availableAbilities []entities.WeaponType
}

func NewChooseAbilityPlugin(plugins *core.PluginManager) *ChooseAbilityPlugin {
	return &ChooseAbilityPlugin{
		availableAbilities: []entities.WeaponType{templates.BasicWeapon, templates.DaggersWeapon, templates.ProtectionWeapon},
	}
}

func (cp *ChooseAbilityPlugin) ID() string {
	return "ChooseAbilityPlugin"
}

func (cp *ChooseAbilityPlugin) Init(kernel *core.GameKernel) error {
	cp.kernel = kernel
	return nil
}

func (cp *ChooseAbilityPlugin) Update() error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		cp.kernel.EventBus.Publish("NewAbility", templates.BasicWeapon)
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		cp.kernel.EventBus.Publish("NewAbility", templates.DaggersWeapon)
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		cp.kernel.EventBus.Publish("NewAbility", templates.ProtectionWeapon)
	}

	return nil
}

func (cp *ChooseAbilityPlugin) Draw(screen *ebiten.Image) {
	text.Draw(screen, "Qual abilidade você quer?:", fontface.FontFace, 300, 150, color.White)

	for i, aa := range cp.availableAbilities {
		col := color.White

		name := ""

		switch aa {
		case templates.BasicWeapon:
			name = "1. Basic Weapon"
		case templates.DaggersWeapon:
			name = "2. Daggers Weapon"
		case templates.ProtectionWeapon:
			name = "3. Protection Weapon"
		}

		text.Draw(screen, name, fontface.FontFace, 300, 200+(i*30), col)
	}
}
