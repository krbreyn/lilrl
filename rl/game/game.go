package game

func MakeNewDebugGame() *RLGame {
	game := RLGame{
		M: GameMap{
			DepthMap: make(map[int]Room),
			Player: Actor{
				Name:   "player",
				Rune:   '@',
				Pos:    Vec2{X: 5, Y: 5},
				Energy: 10,
				Speed:  10,
				Depth:  1,
			},
		},
		Depth: 1,
	}

	new_floor, _ := GenNewFloor(game.Depth)
	debug_room := Room{
		Tiles: new_floor,
		// Actors: []*Actor{
		// 	{
		// 		Name:   "bat",
		// 		Rune:   'b',
		// 		Pos:    Vec2{X: 4, Y: 8},
		// 		Room:   Vec3{0, 0, 0},
		// 		AI:     WanderAI{},
		// 		Energy: 5,
		// 		Speed:  5,
		// 	},
		// },
	}

	game.M.AddNewRoom(1, debug_room)

	return &game
}

type RLGame struct {
	M     GameMap
	Depth int
	UI    UI
}

//AlertMsg func - pause until user hits enter, ex. health below 50%

const turnsPerUpdate uint8 = 1

func (g *RLGame) Update(key string) {

	/* process player action */
	playerAction := GetPlayerAction(key, &g.M.Player)
	pa_result := playerAction.Perform(g)
	if !pa_result.Succeeded {
		return
	}
	g.M.Player.Energy = 0

	/* process other actor actions until it's the player's turn again */
	for {
		nextPlayerTurn := g.M.Player.Energy == g.M.Player.Speed
		if nextPlayerTurn {
			break
		}

		for _, a := range g.M.DepthMap[g.M.Player.Depth].Actors {
			nextActorTurn := a.Energy == a.Speed
			if !nextActorTurn {
				a.Energy += turnsPerUpdate
				continue
			}
			action := a.AI.DecideAction(a, nil)
			result := action.Perform(g)
			if result.Succeeded {
				a.Energy = 0
			}
		}
		g.M.Turn += int(turnsPerUpdate)
		g.M.Player.Energy += turnsPerUpdate
	}

}

func (g *RLGame) RenderUI() string {
	return g.UI.RenderScreen(&g.M, g.Depth)
}
