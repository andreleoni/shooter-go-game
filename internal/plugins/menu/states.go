package menu

type State int

const (
	MenuState State = iota
	CharacterSelectState
	PlayingState
	ChoosingAbilityState
	GameOverState
)
