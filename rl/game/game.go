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
				Map:    Vec3{0, 0, 0},
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
		Entities: []*Actor{
			{
				Name:   "bat",
				Char:   'b',
				Pos:    Vec2{X: 4, Y: 8},
				Map:    Vec3{0, 0, 0},
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

type GameMap struct {
	Turn    int
	Player  Actor
	Rooms   []Room
	RoomMap map[Vec3]Room
}

func (m *GameMap) AddNewRoom(pos Vec3, room Room) {
	m.RoomMap[pos] = room
}

type Room struct {
	Pos      Vec3
	Tiles    [][]rune
	Entities []*Actor
}

// Follows the comma, ok pattern.
func (m *GameMap) EntityAtPos(pos Vec2, room Vec3) (*Actor, bool) {
	if m.Player.Pos == pos {
		return &m.Player, true
	}
	for _, e := range m.RoomMap[room].Entities {
		if pos == e.Pos {
			return e, true
		}
	}
	return &Actor{}, false
}

func (g *RLGame) HandleEntityMove(e *Actor, target Vec2) {
	target = Vec2{X: target.X + e.Pos.X, Y: target.Y + e.Pos.Y}

	room := g.M.RoomMap[e.Map]
	if target.X < 0 || target.X > len(room.Tiles)-1 || target.Y < 0 || target.Y > len(room.Tiles[0])-1 {
		if e == &g.M.Player {
			g.UI.NewStatusMsg("You bump into the edge!")
		}
		return
	}

	if other_e, ok := g.M.EntityAtPos(target, e.Map); !ok {
		e.Pos.X = target.X
		e.Pos.Y = target.Y
	} else {
		if e == &g.M.Player {
			g.UI.NewStatusMsg(fmt.Sprintf("You bump into the %s!", other_e.Name))
		} else {
			g.UI.NewStatusMsg(fmt.Sprintf("The %s bumps into you!", e.Name))
		}
		return // do nothing for now
	}
}

func (g *RLGame) HandleAction(e *Actor, action Action) {
	e.Energy = 0

	switch action := action.(type) {
	case WaitAction:
		return
	case MoveAction:
		g.HandleEntityMove(e, action.Target)
	}
}

var turnsPerUpdate uint8 = 1

func (g *RLGame) Update(PlayerAction Action) {
	g.HandleAction(&g.M.Player, PlayerAction)

	for g.M.Player.Energy != g.M.Player.Speed {
		for _, e := range g.M.RoomMap[g.M.Player.Map].Entities {
			if e.Energy != e.Speed {
				e.Energy += turnsPerUpdate
				continue
			}
			g.HandleAction(e, e.AI.DecideAction(e, nil))
		}
		g.M.Turn += int(turnsPerUpdate)
		g.M.Player.Energy += turnsPerUpdate
	}

}

func (g *RLGame) RenderUI() string {
	return g.UI.RenderScreen(&g.M)
}
