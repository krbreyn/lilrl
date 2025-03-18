package game

import "github.com/charmbracelet/lipgloss"

type Vec2 struct {
	X, Y int
}

type Vec3 struct {
	X, Y, Z int
}

type Tile struct {
	Char []rune
	//type
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
	Rooms   []Room
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
	if m.Player.Pos == pos {
		return &m.Player, true
	}
	for _, e := range m.RoomMap[room].Actors {
		if pos == e.Pos {
			return e, true
		}
	}
	return &Actor{}, false
}
