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
	Turn     int
	Player   Actor
	DepthMap map[int]Room
}

func (m *GameMap) AddNewRoom(depth int, room Room) {
	m.DepthMap[depth] = room
}

type Room struct {
	Tiles  [][]Tile
	Actors []*Actor
}

// Follows the comma, ok pattern.
func (m *GameMap) ActorAtPos(pos Vec2, depth int) (*Actor, bool) {
	if m.Player.Pos == pos && m.Player.Depth == depth {
		return &m.Player, true
	}
	for _, a := range m.DepthMap[depth].Actors {
		if pos == a.Pos {
			return a, true
		}
	}
	return &Actor{}, false
}

func (m *GameMap) TileAtPos(pos Vec2, depth int) Tile {
	r := m.DepthMap[depth]
	return r.Tiles[pos.Y][pos.X]
}
