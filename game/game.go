package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type RLGame struct {
	Player NPC
	M      GameMap
	UI     UI
}

//AlertMsg func - pause until user hits enter, ex. health below 50%

type GameMap struct {
	Tiles    [][]rune
	Entities []*NPC
}

// Follows the comma, ok pattern.
func (m *GameMap) EntityAtPos(x, y int) (*NPC, bool) {
	for _, e := range m.Entities {
		if x == e.Pos.X && y == e.Pos.Y {
			return e, true
		}
	}
	return &NPC{}, false
}

func (g *RLGame) HandleEntityMove(e *NPC, target Vec2) {
	target = Vec2{X: target.X + e.Pos.X, Y: target.Y + e.Pos.Y}
	if target.X < 0 || target.X > len(g.M.Tiles)-1 || target.Y < 0 || target.Y > len(g.M.Tiles[0])-1 {
		if e == &g.Player {
			g.NewStatusMsg("You bump into the edge!")
		}
		return
	}

	if other_e, ok := g.M.EntityAtPos(target.X, target.Y); !ok {
		e.Pos.X = target.X
		e.Pos.Y = target.Y
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
		g.HandleEntityMove(&g.Player, Vec2{0, -1})

	case "down", "j":
		g.HandleEntityMove(&g.Player, Vec2{0, 1})

	case "left", "h":
		g.HandleEntityMove(&g.Player, Vec2{-1, 0})

	case "right", "l":
		g.HandleEntityMove(&g.Player, Vec2{1, 0})

	default:
		return // do not process turn
	}

	for _, npc := range g.M.Entities {
		action := npc.AI.DecideAction(npc, nil)
		switch action.Type {
		case MoveAction:
			g.HandleEntityMove(npc, action.Target)
		}
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
		ui.StatusMsgs = append(ui.StatusMsgs, "")
		ui.StatusMsgs = append(ui.StatusMsgs, "")
		ui.StatusMsgs = append(ui.StatusMsgs, "Welcome to Lil' RL!")
	}
	var last3 string
	j := 3
	for range 3 {
		last3 += ui.StatusMsgs[len(ui.StatusMsgs)-j]
		if j != 1 {
			last3 += "\n"
		}
		j--
	}

	return last3
}

func (g *RLGame) RenderMap() string {
	var sb strings.Builder

	for yi, y := range g.M.Tiles {
		for xi, x := range y {
			if g.Player.Pos.X == xi && g.Player.Pos.Y == yi {
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

var statusMsgStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(50)

//var sideBarStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2).Width(10).Height(10)
