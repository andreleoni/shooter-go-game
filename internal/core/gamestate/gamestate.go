package gamestate

type State int

const (
	MenuState State = iota
	CharacterSelectState
	PlayingState
	GameOverState
)
