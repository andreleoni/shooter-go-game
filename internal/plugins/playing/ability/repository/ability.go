package repository

import (
	abilitiesentities "game/internal/plugins/playing/ability/entities/abilities"
)

type AbilityRecord struct {
	ability       abilitiesentities.Ability
	acquiredOrder int
}

type Ability struct {
	records map[string]*AbilityRecord
}

func NewAbility() *Ability {
	abilityRecords := make(map[string]*AbilityRecord)

	return &Ability{
		records: abilityRecords,
	}
}

func (a *Ability) Add(ability abilitiesentities.Ability) {
	currentAbility, exists := a.records[ability.ID()]

	if !exists {
		a.records[ability.ID()] =
			&AbilityRecord{
				ability:       ability,
				acquiredOrder: len(a.records),
			}

		return
	}

	currentAbility.ability.IncreaseLevel()
}

func (a *Ability) Get() []abilitiesentities.Ability {
	sortedAbilities := make([]abilitiesentities.Ability, len(a.records))

	for _, r := range a.records {
		sortedAbilities[r.acquiredOrder] = r.ability
	}

	return sortedAbilities
}
