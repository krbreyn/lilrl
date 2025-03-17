package game

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	StatusMsgs []string
}

func (ui *UI) NewStatusMsg(msg string) {
	ui.StatusMsgs = append(ui.StatusMsgs, msg)
}

func (ui *UI) RenderScreen(m *GameMap) string {
	statusMsg := statusMsgStyle.Render(ui.RenderStatusBox())
	game_map := ui.RenderMap(m, m.Player.MapPos)

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

func (ui *UI) RenderMap(m *GameMap, room Vec3) string {
	var sb strings.Builder
	r := m.RoomMap[room]

	if len(r.Tiles) == 0 {
		return "Something went horribly wrong! You are in a room that doesn't exist!"
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
