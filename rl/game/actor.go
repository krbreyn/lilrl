package game

import "math/rand"

type ActionType int

type Actor struct {
	Name   string
	Char   rune
	Energy uint8
	Speed  uint8
	Pos    Vec2
	Map    Vec3
	AI     AI
}

type AI interface {
	DecideAction(n *Actor, ctx *AIContext) Action
}

type AIContext interface {
	GetPath(target Vec2) (path []Vec2)
}

type WanderAI struct {
	ExtraReach int
}

func (w WanderAI) DecideAction(n *Actor, ctx *AIContext) Action {
	dir := Vec2{
		rand.Intn(3) - 1,
		rand.Intn(3) - 1,
	}
	return MoveAction{Target: dir}
}
