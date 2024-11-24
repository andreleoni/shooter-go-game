package templates

import "game/internal/plugins/weapon/entities"

const (
	BasicWeapon entities.WeaponType = iota
	DaggersWeapon
	ProtectionWeapon
)

var EnemyTemplates = map[entities.WeaponType]*entities.WeaponType{}
