package game

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	statusMsgs []string
	buffer     [][]rune
}

func (ui *UI) NewStatusMsg(msg string) {
	ui.statusMsgs = append(ui.statusMsgs, msg)
}

func (ui *UI) RenderScreen(m *GameMap) string {
	statusMsg := statusMsgStyle.Render(ui.RenderStatusBox())
	game_map := ui.RenderMap(m, m.Player.MapPos)

	layout := lipgloss.JoinVertical(lipgloss.Left, game_map, statusMsg)
	return layout
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

func (ui *UI) RenderMap(m *GameMap, room Vec3) string {
	var sb strings.Builder
	r, ok := m.RoomMap[room]

	if !ok {
		return "Something went horribly wrong! You are in a room that doesn't exist!"
	}

	if len(r.Tiles) == 0 {
		return "Something went horribly wrong! You are in a room with no tiles!"
	}

	for yi, y := range r.Tiles {
		for xi, x := range y {
			if m.Player.RoomPos.X == xi && m.Player.RoomPos.Y == yi {
				sb.WriteRune(m.Player.Char)
				continue
			}
			if e, ok := m.EntityAtPos(Vec2{xi, yi}, room); ok {
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
