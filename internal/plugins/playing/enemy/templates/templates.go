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
		Name:                 "basic",
		MaxHealth:            20,
		Speed:                80,
		Damage:               10,
		Size:                 40,
		Power:                10,
		RunningAnimationTime: 0.3,
	},
	FastEnemy: {
		Name:                 "fast",
		MaxHealth:            10,
		Speed:                100,
		Damage:               5,
		Size:                 30,
		Power:                5,
		RunningAnimationTime: 0.1,
	},
	TankEnemy: {
		Name:                 "tank",
		MaxHealth:            50,
		Speed:                50,
		Damage:               20,
		Size:                 70,
		Power:                20,
		RunningAnimationTime: 0.1,
	},
	RangedEnemy: {
		Name:                 "ranged",
		MaxHealth:            30,
		Speed:                75,
		Damage:               15,
		Size:                 25,
		Power:                15,
		RunningAnimationTime: 0.2,
	},
}

func init() {
	for _, t := range EnemyTemplates {
		t.RunningLeftAnimationSprite = assets.NewAnimation(t.RunningAnimationTime)
		err := t.RunningLeftAnimationSprite.LoadFromJSON(
			"assets/images/enemies/"+t.Name+"/run/left/asset.json",
			"assets/images/enemies/"+t.Name+"/run/left/asset.png")
		if err != nil {
			log.Fatal("Failed to load enemy left run asset:", err)
		}

		t.RunningRightAnimationSprite = assets.NewAnimation(t.RunningAnimationTime)
		err = t.RunningRightAnimationSprite.LoadFromJSON(
			"assets/images/enemies/"+t.Name+"/run/right/asset.json",
			"assets/images/enemies/"+t.Name+"/run/right/asset.png")
		if err != nil {
			log.Fatal("Failed to load enemy run run asset:", err)
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
