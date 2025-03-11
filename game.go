package main

import "strings"

type RLGame struct {
	Player Entity
	Map    GameMap
}

type GameMap struct {
	Tiles    [][]rune
	Entities []Entity
}

type Entity struct {
	Name string
	Char rune
	X, Y int
}

// Follows the comma, ok pattern.
func (m *GameMap) EntityAtPos(x, y int) (Entity, bool) {
	for _, e := range m.Entities {
		if x == e.X && y == e.Y {
			return e, true
		}
	}
	return Entity{}, false
}

func (g *RLGame) RenderMap() string {
	var sb strings.Builder

	for yi, y := range g.Map.Tiles {
		for xi, x := range y {
			if g.Player.X == xi && g.Player.Y == yi {
				sb.WriteRune(g.Player.Char)
				continue
			}
			if e, ok := g.Map.EntityAtPos(xi, yi); ok {
				sb.WriteRune(e.Char)
			} else {
				sb.WriteRune(x)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
