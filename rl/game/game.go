package game

import (
	"fmt"
)

func MakeNewDebugGame() *RLGame {
	game := RLGame{
		M: GameMap{
			RoomMap: make(map[Vec3]Room),
			Player: Actor{
				Name:   "player",
				Char:   '@',
				Pos:    Vec2{X: 5, Y: 5},
				Room:   Vec3{0, 0, 0},
				Energy: 10,
				Speed:  10,
			},
		},
	}

	debug_room := Room{
		Tiles: [][]rune{
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
		},
		Actors: []*Actor{
			{
				Name:   "bat",
				Char:   'b',
				Pos:    Vec2{X: 4, Y: 8},
				Room:   Vec3{0, 0, 0},
				AI:     WanderAI{},
				Energy: 5,
				Speed:  5,
			},
		},
	}

	game.M.AddNewRoom(Vec3{0, 0, 0}, debug_room)

	return &game
}

type RLGame struct {
	M  GameMap
	UI UI
}

//AlertMsg func - pause until user hits enter, ex. health below 50%

const turnsPerUpdate uint8 = 1

func (g *RLGame) Update(PlayerAction Action) {
	shouldReturn := g.HandleAction(&g.M.Player, PlayerAction)
	if shouldReturn {
		return
	}

	for {
		nextPlayerTurn := g.M.Player.Energy == g.M.Player.Speed
		if nextPlayerTurn {
			break
		}

		for _, a := range g.M.RoomMap[g.M.Player.Room].Actors {
			nextActorTurn := a.Energy == a.Speed
			if !nextActorTurn {
				a.Energy += turnsPerUpdate
				continue
			}
			g.HandleAction(a, a.AI.DecideAction(a, nil))
		}
		g.M.Turn += int(turnsPerUpdate)
		g.M.Player.Energy += turnsPerUpdate
	}

}

func (g *RLGame) RenderUI() string {
	return g.UI.RenderScreen(&g.M)
}

/* actions */
func (g *RLGame) HandleAction(e *Actor, action Action) (shouldReturn bool) {
	e.Energy = 0

	switch action := action.(type) {
	case InvalidKeyAction:
		return true
	case WaitAction:
		return false
	case MoveAction:
		g.HandleMoveAction(e, action.Target)
	}
	return false
}

func (g *RLGame) HandleMoveAction(a *Actor, target Vec2) {
	target = Vec2{X: target.X + a.Pos.X, Y: target.Y + a.Pos.Y}

	room := g.M.RoomMap[a.Room]
	if target.X < 0 || target.X > len(room.Tiles)-1 || target.Y < 0 || target.Y > len(room.Tiles[0])-1 {
		if a == &g.M.Player {
			g.UI.NewStatusMsg("You bump into the edge!")
		}
		return
	}

	if other_e, ok := g.M.ActorAtPos(target, a.Room); !ok {
		a.Pos.X = target.X
		a.Pos.Y = target.Y
	} else {
		if a == &g.M.Player {
			g.UI.NewStatusMsg(fmt.Sprintf("You bump into the %s!", other_e.Name))
		} else {
			g.UI.NewStatusMsg(fmt.Sprintf("The %s bumps into you!", a.Name))
		}
		return // do nothing for now
	}
}
