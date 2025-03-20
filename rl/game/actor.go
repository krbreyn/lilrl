package game

import "math/rand"

type Actor struct {
	Name   string
	Rune   rune
	Energy uint8
	Speed  uint8
	Pos    Vec2
	Depth  int
	AI     AI
}

type AI interface {
	DecideAction(a *Actor, ctx *AIContext) Action
}

type AIContext interface {
	GetPath(target Vec2) (path []Vec2)
}

type WanderAI struct{}

func (w WanderAI) DecideAction(a *Actor, ctx *AIContext) Action {
	dir := Vec2{
		rand.Intn(3) - 1,
		rand.Intn(3) - 1,
	}
	return MoveAction{Actor: a, Target: dir}
}
