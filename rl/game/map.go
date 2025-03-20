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
	Name  string
	Rune  rune
	Color lipgloss.Color
}

type TileType int

const (
	TileEmpty TileType = iota
	TileWall
	TileFloor
	TileDoor
	TileStairDown
	TileStairUp
)

func NewTile(t TileType) Tile {
	name, char := t.Repr()
	return Tile{
		Type: t,
		Name: name,
		Rune: char,
	}
}

func (t TileType) Repr() (Name string, Rune rune) {
	switch t {
	case TileEmpty:
		return "Empty", ' '
	case TileWall:
		return "Wall", '#'
	case TileFloor:
		return "Floor", '.'
	case TileDoor:
		return "Door", '+'
	case TileStairDown:
		return "Stair Down", '>'
	case TileStairUp:
		return "Stair Up", '<'

	default:
		return "Unknown", '?'
	}
}

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
	Tiles  [][]Tile
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

func (m *GameMap) TileAtPos(pos Vec2, room Vec3) Tile {
	r := m.RoomMap[room]
	return r.Tiles[pos.Y][pos.X]
}
