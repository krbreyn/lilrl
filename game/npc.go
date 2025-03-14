package game

import "math/rand"

type ActionType int

type NPC struct {
	Name string
	Char rune
	Pos  Vec2
	AI   AI
}

const (
	MoveAction ActionType = iota + 1
)

type Action struct {
	Type   ActionType
	Target Vec2
}

type AI interface {
	DecideAction(n *NPC, ctx *AIContext) Action
}

type AIContext interface {
	GetPath(target Vec2) (path []Vec2)
}

type WanderAI struct {
	ExtraReach int
}

func (w WanderAI) DecideAction(n *NPC, ctx *AIContext) Action {
	dirs := []Vec2{{0, 1 + w.ExtraReach}, {0, -1 - w.ExtraReach}, {1 + w.ExtraReach, 0}, {-1 - w.ExtraReach, 0}}
	move := dirs[rand.Intn(len(dirs))]
	return Action{Type: MoveAction, Target: move}
}
