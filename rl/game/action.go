package game

type Action any

// player only
type InvalidKeyAction struct{}

type WaitAction struct{}

type MoveAction struct {
	Target Vec2
}
