package templates

import (
	"game/internal/assets"
	"game/internal/plugins/playing/enemy/entities"
	"log"
)

const (
	BasicEnemy entities.EnemyType = iota
	FastEnemy
	TankEnemy
	RangedEnemy
)

var EnemyTemplates = map[entities.EnemyType]*entities.EnemyTemplate{
	BasicEnemy: {
		Name:      "basic",
		MaxHealth: 20,
		Speed:     100,
		Damage:    10,
		Size:      40,
		Power:     10,
	},
	FastEnemy: {
		Name:      "fast",
		MaxHealth: 10,
		Speed:     200,
		Damage:    5,
		Size:      100,
		Power:     5,
	},
	TankEnemy: {
		Name:      "tank",
		MaxHealth: 50,
		Speed:     50,
		Damage:    20,
		Size:      60,
		Power:     20,
	},
	RangedEnemy: {
		Name:      "ranged",
		MaxHealth: 30,
		Speed:     75,
		Damage:    15,
		Size:      18,
		Power:     15,
	},
}

func init() {
	for _, t := range EnemyTemplates {
		t.RunningAnimationSprite = assets.NewAnimation(0.1)
		err := t.RunningAnimationSprite.LoadFromJSON(
			"assets/images/enemies/"+t.Name+"/run/asset.json",
			"assets/images/enemies/"+t.Name+"/run/asset.png")
		if err != nil {
			log.Fatal("Failed to load enemy run asset:", err)
		}

		t.DeathAnimation = assets.NewAnimation(0.1)
		err = t.DeathAnimation.LoadFromJSON(
			"assets/images/enemies/"+t.Name+"/death/asset.json",
			"assets/images/enemies/"+t.Name+"/death/asset.png")
		if err != nil {
			log.Fatal("Failed to load enemy death asset:", err)
		}
	}
}
