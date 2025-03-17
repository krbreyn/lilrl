package game

import (
	"fmt"
)

func MakeNewDebugGame() *RLGame {
	game := RLGame{
		M: GameMap{
			RoomMap: make(map[Vec3]Room),
			Player: NPC{
				Name: "player",
				Char: '@',
				Pos:  Vec2{X: 5, Y: 5},
				Map:  Vec3{0, 0, 0},
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
		Entities: []*NPC{
			{Name: "bat", Char: 'b', Pos: Vec2{X: 4, Y: 8}, Map: Vec3{0, 0, 0}, AI: WanderAI{ExtraReach: 1}},
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
	Player  NPC
	Rooms   []Room
	RoomMap map[Vec3]Room
}

func (m *GameMap) AddNewRoom(pos Vec3, room Room) {
	m.RoomMap[pos] = room
}

type Room struct {
	Pos      Vec3
	Tiles    [][]rune
	Entities []*NPC
}

// Follows the comma, ok pattern.
func (m *GameMap) EntityAtPos(pos Vec2, room Vec3) (*NPC, bool) {
	for _, e := range m.RoomMap[room].Entities {
		if pos == e.Pos {
			return e, true
		}
	}
	return &NPC{}, false
}

func (g *RLGame) HandleEntityMove(e *NPC, target Vec2) {
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
		}
		return // do nothing for now
	}
}

func (g *RLGame) HandleUpdate(key string) {

	switch key {
	/* movement */
	case "up", "k", "down", "j", "left", "h", "right", "l", "y", "u", "b", "n":
		var pos Vec2
		switch key {
		case "up", "k":
			pos = Vec2{0, -1}
		case "down", "j":
			pos = Vec2{0, 1}
		case "left", "h":
			pos = Vec2{-1, 0}
		case "right", "l":
			pos = Vec2{1, 0}
		case "y":
			pos = Vec2{-1, -1}
		case "u":
			pos = Vec2{1, -1}
		case "b":
			pos = Vec2{-1, 1}
		case "n":
			pos = Vec2{1, 1}
		}
		g.HandleEntityMove(&g.M.Player, pos)

	default:
		return // do not process turn
	}

	for _, npc := range g.M.RoomMap[g.M.Player.Map].Entities {
		action := npc.AI.DecideAction(npc, nil)
		switch action.Type {
		case MoveAction:
			g.HandleEntityMove(npc, action.Target)
		}
	}
}

func (g *RLGame) RenderUI() string {
	return g.UI.RenderScreen(&g.M)
}
