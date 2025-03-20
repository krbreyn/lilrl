package game

import (
	"fmt"
	"math/rand"
)

type Action interface {
	Perform(g *RLGame) ActionResult
}

type ActionResult struct {
	Succeeded bool
	Alternate Action
}

type DebugAction struct {
	Cmd string
}

func (a DebugAction) Perform(g *RLGame) ActionResult {
	switch a.Cmd {
	case "newmap":
		new_floor, _ := GenDngRogue(90, 50)
		g.M.DepthMap[g.Depth] = Room{Tiles: new_floor, Actors: g.M.DepthMap[g.M.Player.Depth].Actors}
	}
	return ActionResult{true, nil}
}

// player only
type InvalidKeyAction struct{}

func (a InvalidKeyAction) Perform(g *RLGame) ActionResult {
	return ActionResult{false, nil}
}

type WaitAction struct{}

func (a WaitAction) Perform(g *RLGame) ActionResult {
	return ActionResult{true, nil}
}

type MoveAction struct {
	Actor  *Actor
	Target Vec2
}

func (a MoveAction) Perform(g *RLGame) ActionResult {
	target := Vec2{X: a.Target.X + a.Actor.Pos.X, Y: a.Target.Y + a.Actor.Pos.Y}

	room := g.M.DepthMap[a.Actor.Depth]

	/* bounds/wall checking */
	if target.X < 0 || target.X > len(room.Tiles[0])-1 || target.Y < 0 || target.Y > len(room.Tiles)-1 {
		if a.Actor == &g.M.Player {
			g.UI.NewStatusMsg("You bump into the edge!")
		}
		return ActionResult{true, nil}
	}

	t_tile := g.M.TileAtPos(target, a.Actor.Depth)
	switch t_tile.Type {
	case TileWall:
		if a.Actor == &g.M.Player {
			g.UI.NewStatusMsg("You bump into the wall!")
		}
		return ActionResult{true, nil}
	}

	/* actor/combat checking */
	if other_a, ok := g.M.ActorAtPos(target, a.Actor.Depth); !ok {
		a.Actor.Pos.X = target.X
		a.Actor.Pos.Y = target.Y
	} else {
		if a.Actor == &g.M.Player {
			g.UI.NewStatusMsg(fmt.Sprintf("You bump into the %s!", other_a.Name))
		} else if other_a == &g.M.Player {
			g.UI.NewStatusMsg(fmt.Sprintf("The %s bumps into you!", a.Actor.Name))
		}
	}

	return ActionResult{true, nil}
}

type StairMoveAction struct {
	Actor *Actor
	Dir   string // "up" or "down"
}

func (a StairMoveAction) Perform(g *RLGame) ActionResult {
	switch a.Dir {
	case "down":
		if g.M.TileAtPos(a.Actor.Pos, a.Actor.Depth).Type != TileStairDown {
			if a.Actor == &g.M.Player {
				g.UI.NewStatusMsg("You cannot go down here!")
			}
			return ActionResult{false, nil}
		} else {
			g.Depth++
			g.M.Player.Depth = g.Depth
			new_map, spots := GenNewFloor(g.Depth)
			newSpot := spots[rand.Intn(len(spots))]
			g.M.Player.Pos.X = newSpot.X
			g.M.Player.Pos.Y = newSpot.Y
			new_room := Room{Tiles: new_map}
			g.M.AddNewRoom(g.Depth, new_room)
		}
	}

	return ActionResult{false, nil}
}

type AttackMeleeAction struct {
	Actor  *Actor
	Target *Actor
}
