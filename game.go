package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type RLGame struct {
	Player Entity
	M      GameMap
	UI     UI
}

//AlertMsg func - pause until user hits enter, ex. health below 50%

type GameMap struct {
	Tiles    [][]rune
	Entities []*Entity
}

type Entity struct {
	Name string
	Char rune
	X, Y int
	//Behavior func() NPCActionMsg
}

// Follows the comma, ok pattern.
func (m *GameMap) EntityAtPos(x, y int) (*Entity, bool) {
	for _, e := range m.Entities {
		if x == e.X && y == e.Y {
			return e, true
		}
	}
	return &Entity{}, false
}

func (g *RLGame) HandleEntityMove(e *Entity, targetX, targetY int) {
	if targetX < 0 || targetX > len(g.M.Tiles)-1 || targetY < 0 || targetY > len(g.M.Tiles[0])-1 {
		if e == &g.Player {
			g.NewStatusMsg("You bump into the edge!")
		}
		return
	}
	if other_e, ok := g.M.EntityAtPos(targetX, targetY); !ok {
		e.X = targetX
		e.Y = targetY
	} else {
		if e == &g.Player {
			g.NewStatusMsg(fmt.Sprintf("You bump into the %s!", other_e.Name))
		}
		return // do nothing for now
	}
}

func (g *RLGame) HandleUpdate(key string) {
	switch key {
	case "up", "k":
		g.HandleEntityMove(&g.Player, g.Player.X, g.Player.Y-1)

	case "down", "j":
		g.HandleEntityMove(&g.Player, g.Player.X, g.Player.Y+1)

	case "left", "h":
		g.HandleEntityMove(&g.Player, g.Player.X-1, g.Player.Y)

	case "right", "l":
		g.HandleEntityMove(&g.Player, g.Player.X+1, g.Player.Y)
	}
}

type UI struct {
	StatusMsgs []string
}

func (g *RLGame) NewStatusMsg(msg string) {
	g.UI.StatusMsgs = append(g.UI.StatusMsgs, msg)
}

func (g *RLGame) RenderScreen() string {
	statusMsg := statusMsgStyle.Render(g.UI.RenderStatusBox())
	game_map := g.RenderMap()

	layout := lipgloss.JoinVertical(lipgloss.Left, game_map, statusMsg)
	return layout
}

func (ui *UI) RenderStatusBox() string {
	if len(ui.StatusMsgs) == 0 {
		ui.StatusMsgs = append(ui.StatusMsgs, "Welcome to Lil' RL!")
	}
	return ui.StatusMsgs[len(ui.StatusMsgs)-1]
}

func (g *RLGame) RenderMap() string {
	var sb strings.Builder

	for yi, y := range g.M.Tiles {
		for xi, x := range y {
			if g.Player.X == xi && g.Player.Y == yi {
				sb.WriteRune(g.Player.Char)
				continue
			}
			if e, ok := g.M.EntityAtPos(xi, yi); ok {
				sb.WriteRune(e.Char)
			} else {
				sb.WriteRune(x)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

var statusMsgStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Width(50)

//var sideBarStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2).Width(10).Height(10)
