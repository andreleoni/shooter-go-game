package chooseability

import (
	"game/internal/core"
	"game/internal/plugins/menu/fontface"
	abilitiesentities "game/internal/plugins/playing/ability/entities/abilities"
	abilitiesentitiesbasic "game/internal/plugins/playing/ability/entities/abilities/basic"
	abilitiesentitiesdagger "game/internal/plugins/playing/ability/entities/abilities/dagger"
	abilitiesentitiesprotection "game/internal/plugins/playing/ability/entities/abilities/protection"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type ChooseAbilityPlugin struct {
	kernel             *core.GameKernel
	availableAbilities map[string]abilitiesentities.Ability
	abilities          []string
}

func NewChooseAbilityPlugin(plugins *core.PluginManager) *ChooseAbilityPlugin {
	abilities := []string{
		"BasicWeapon",
		"DaggersWeapon",
		"ProtectionWeapon",
	}

	abilitiesByName := map[string]abilitiesentities.Ability{
		"BasicWeapon":      abilitiesentitiesbasic.New(),
		"DaggersWeapon":    abilitiesentitiesdagger.New(),
		"ProtectionWeapon": abilitiesentitiesprotection.New(),
		// "HealAbility":      abilitiesentitiesheal.New(),
	}

	cp := ChooseAbilityPlugin{
		availableAbilities: abilitiesByName,
		abilities:          abilities,
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

	for i, key := range cp.abilities {
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
