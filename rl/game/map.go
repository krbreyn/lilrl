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
