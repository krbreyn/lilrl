package game

import (
	"fmt"
)

func MakeNewDebugGame() *RLGame {
	game := RLGame{
		M: GameMap{
			RoomMap: make(map[Vec3]Room),
			Player: NPC{
				Name:    "player",
				Char:    '@',
				RoomPos: Vec2{X: 5, Y: 5},
				MapPos:  Vec3{0, 0, 0},
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
			{Name: "bat", Char: 'b', RoomPos: Vec2{X: 4, Y: 8}, MapPos: Vec3{0, 0, 0}, AI: WanderAI{ExtraReach: 1}},
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
		if pos == e.RoomPos {
			return e, true
		}
	}
	return &NPC{}, false
}

func (g *RLGame) HandleEntityMove(e *NPC, target Vec2, room Vec3) {
	target = Vec2{X: target.X + e.RoomPos.X, Y: target.Y + e.RoomPos.Y}

	t_room := g.M.RoomMap[room]
	if target.X < 0 || target.X > len(t_room.Tiles)-1 || target.Y < 0 || target.Y > len(t_room.Tiles[0])-1 {
		if e == &g.M.Player {
			g.UI.NewStatusMsg("You bump into the edge!")
		}
		return
	}

	if other_e, ok := g.M.EntityAtPos(target, room); !ok {
		e.RoomPos.X = target.X
		e.RoomPos.Y = target.Y
	} else {
		if e == &g.M.Player {
			g.UI.NewStatusMsg(fmt.Sprintf("You bump into the %s!", other_e.Name))
		}
		return // do nothing for now
	}
}

func (g *RLGame) HandleUpdate(key string) {
	switch key {

	/*
		movement
	*/
	case "up", "k":
		g.HandleEntityMove(&g.M.Player, Vec2{0, -1}, g.M.Player.MapPos)
	case "down", "j":
		g.HandleEntityMove(&g.M.Player, Vec2{0, 1}, g.M.Player.MapPos)
	case "left", "h":
		g.HandleEntityMove(&g.M.Player, Vec2{-1, 0}, g.M.Player.MapPos)
	case "right", "l":
		g.HandleEntityMove(&g.M.Player, Vec2{1, 0}, g.M.Player.MapPos)
	case "y":
		g.HandleEntityMove(&g.M.Player, Vec2{-1, -1}, g.M.Player.MapPos)
	case "u":
		g.HandleEntityMove(&g.M.Player, Vec2{1, -1}, g.M.Player.MapPos)
	case "b":
		g.HandleEntityMove(&g.M.Player, Vec2{-1, 1}, g.M.Player.MapPos)
	case "n":
		g.HandleEntityMove(&g.M.Player, Vec2{1, 1}, g.M.Player.MapPos)

	default:
		return // do not process turn
	}

	for _, npc := range g.M.RoomMap[g.M.Player.MapPos].Entities {
		action := npc.AI.DecideAction(npc, nil)
		switch action.Type {
		case MoveAction:
			g.HandleEntityMove(npc, action.Target, npc.MapPos)
		}
	}
}

func (g *RLGame) RenderUI() string {
	return g.UI.RenderScreen(&g.M)
}
