package chooseability

import (
	"game/internal/core"
	"game/internal/plugins/menu/fontface"
	weaponentities "game/internal/plugins/playing/weapon/entities/weapons"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type ChooseAbilityPlugin struct {
	kernel             *core.GameKernel
	availableAbilities map[string]weaponentities.Weapon
}

func NewChooseAbilityPlugin(plugins *core.PluginManager) *ChooseAbilityPlugin {
	weaponsByName := map[string]weaponentities.Weapon{
		"BasicWeapon":   weaponentities.NewBasic(plugins),
		"DaggersWeapon": weaponentities.NewDagger(),
		// "ProtectionWeapon": weaponentities.NewProtection(),
	}

	cp := ChooseAbilityPlugin{
		availableAbilities: weaponsByName,
	}

	return &cp
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
		cp.kernel.EventBus.Publish("NewAbility", cp.availableAbilities["BasicWeapon"])
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		cp.kernel.EventBus.Publish("NewAbility", cp.availableAbilities["DaggersWeapon"])
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		cp.kernel.EventBus.Publish("NewAbility", cp.availableAbilities["ProtectionWeapon"])
	}

	return nil
}

func (cp *ChooseAbilityPlugin) Draw(screen *ebiten.Image) {
	text.Draw(screen, "Qual abilidade você quer?:", fontface.FontFace, 300, 150, color.White)

	i := 0

	for key, _ := range cp.availableAbilities {
		i += 1

		col := color.White

		name := ""

		switch key {
		case "BasicWeapon":
			name = "1. Basic Weapon"
		case "DaggersWeapon":
			name = "2. Daggers Weapon"
		case "ProtectionWeapon":
			name = "3. Protection Weapon"
		}

		text.Draw(screen, name, fontface.FontFace, 300, 200+(i*30), col)
	}
}
