package game

import "github.com/charmbracelet/lipgloss"

type Vec2 struct {
	X, Y int
}

type Vec3 struct {
	X, Y, Z int
}

type Tile struct {
	Type  TileType
	Char  []rune
	Name  string
	Color lipgloss.Color
}

type TileType int

const (
	TileEmpty TileType = iota
	TileWall
	TileFloor
	TileStair
)

type GameMap struct {
	Turn    int
	Player  Actor
	RoomMap map[Vec3]Room
}

func (m *GameMap) AddNewRoom(pos Vec3, room Room) {
	m.RoomMap[pos] = room
}

type Room struct {
	Pos    Vec3
	Tiles  [][]rune
	Actors []*Actor
}

// Follows the comma, ok pattern.
func (m *GameMap) ActorAtPos(pos Vec2, room Vec3) (*Actor, bool) {
	if m.Player.Pos == pos && m.Player.Room == room {
		return &m.Player, true
	}
	for _, a := range m.RoomMap[room].Actors {
		if pos == a.Pos {
			return a, true
		}
	}
	return &Actor{}, false
}
