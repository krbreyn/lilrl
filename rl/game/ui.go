package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	statusMsgs []string
}

func (ui *UI) NewStatusMsg(msg string) {
	ui.statusMsgs = append(ui.statusMsgs, msg)
}

func (ui *UI) RenderScreen(m *GameMap, depth int) string {
	statusMsg := ui.RenderStatusBox()
	game_map := ui.RenderMap(m, depth)
	sidebar := ui.RenderSideScreen(m)

	layout := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Left, game_map, statusMsg),
		sidebar,
	)
	return layout
}

func (ui *UI) RenderSideScreen(m *GameMap) string {
	return fmt.Sprintf("turn: %d", m.Turn)
}

func (ui *UI) RenderStatusBox() string {
	if len(ui.statusMsgs) == 0 {
		ui.statusMsgs = append(ui.statusMsgs, "")
		ui.statusMsgs = append(ui.statusMsgs, "")
		ui.statusMsgs = append(ui.statusMsgs, "Welcome to Lil' RL!")
	}
	var last3 string
	j := 3
	for range 3 {
		last3 += ui.statusMsgs[len(ui.statusMsgs)-j]
		if j != 1 {
			last3 += "\n"
		}
		j--
	}

	return last3
}

func (ui *UI) RenderMap(m *GameMap, depth int) string {
	var sb strings.Builder
	r, ok := m.DepthMap[depth]

	if !ok {
		return "Something went horribly wrong! You are in a room that doesn't exist!"
	}

	if len(r.Tiles) == 0 {
		return "Something went horribly wrong! You are in a room with no tiles!"
	}

	for yi, y := range r.Tiles {
		for xi, x := range y {
			if m.Player.Pos.X == xi && m.Player.Pos.Y == yi {
				sb.WriteRune(m.Player.Rune)
				continue
			}
			if e, ok := m.ActorAtPos(Vec2{xi, yi}, depth); ok {
				sb.WriteRune(e.Rune)
			} else {
				sb.WriteRune(x.Rune)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
