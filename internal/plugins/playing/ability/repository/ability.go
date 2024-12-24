package repository

import (
	abilitiesentities "game/internal/plugins/playing/ability/entities/abilities"
)

type Ability struct {
	abilities []abilitiesentities.Ability
}

func NewAbility() *Ability {
	return &Ability{}
}

func (a *Ability) Add(ability abilitiesentities.Ability) {
	a.abilities = append(a.abilities, ability)
}

func (a *Ability) Get() []abilitiesentities.Ability {
	return a.abilities
}
