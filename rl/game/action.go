package game

import "fmt"

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
		g.M.RoomMap[g.M.Player.Room] = Room{Tiles: GenDngRogue(100, 50), Actors: g.M.RoomMap[g.M.Player.Room].Actors}
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

	room := g.M.RoomMap[a.Actor.Room]

	/* bounds/wall checking */
	if target.X < 0 || target.X > len(room.Tiles[0])-1 || target.Y < 0 || target.Y > len(room.Tiles)-1 {
		if a.Actor == &g.M.Player {
			g.UI.NewStatusMsg("You bump into the edge!")
		}
		return ActionResult{true, nil}
	}

	t_tile := g.M.TileAtPos(target, a.Actor.Room)
	switch t_tile.Type {
	case TileWall:
		if a.Actor == &g.M.Player {
			g.UI.NewStatusMsg("You bump into the wall!")
		}
		return ActionResult{true, nil}
	}

	/* actor/combat checking */
	if other_a, ok := g.M.ActorAtPos(target, a.Actor.Room); !ok {
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

type AttackMeleeAction struct {
	Actor  *Actor
	Target *Actor
}
