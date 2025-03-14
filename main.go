package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/krbreyn/lilrl/game"
)

func main() {
	m := initialModel()
	p := tea.NewProgram(m)

	if m, err := p.Run(); err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	} else {
		s, ok := m.(model)
		if ok && s.exitMsg != "" {
			fmt.Println(s.exitMsg)
		}
		os.Exit(0)
	}
}

func initialModel() model {
	game := game.RLGame{
		Player: game.NPC{
			Name: "player",
			Char: '@',
			Pos:  game.Vec2{X: 5, Y: 5},
		},
		M: game.GameMap{
			Tiles: [][]rune{
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
			},
			Entities: []*game.NPC{
				{Name: "bat", Char: 'b', Pos: game.Vec2{X: 4, Y: 8}, AI: game.WanderAI{ExtraReach: 1}},
			},
		},
	}
	m := model{Game: &game, exitMsg: "Goodbye!"}
	return m
}

type model struct {
	Game    *game.RLGame
	exitMsg string
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) View() string {
	return m.Game.RenderScreen()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		key := msg.String()

		switch key {

		case "ctrl+c", "ctrl+d", "ctrl+q":
			return m, tea.Quit

		default:
			m.Game.HandleUpdate(key)
		}

	case tea.WindowSizeMsg:
	}

	return m, tea.Batch(cmds...)
}
