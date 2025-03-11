package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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
	game := RLGame{
		Player: Entity{
			Name: "player",
			Char: '@',
			X:    5,
			Y:    5,
		},
		Map: GameMap{
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
			Entities: []Entity{
				{"Bat", 'b', 3, 2},
			},
		},
	}
	m := model{Game: &game, exitMsg: "Goodbye!"}
	return m
}

type model struct {
	Game    *RLGame
	exitMsg string
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) View() string {
	return m.Game.RenderMap()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		key := msg.String()

		switch key {

		case "ctrl+c", "ctrl+d", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
	}

	return m, tea.Batch(cmds...)
}
