package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rmarken5/ffdp/client/printer"
	pp "github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
)

type WelcomeModel struct {
	choices  []*pp.Player
	cursor   int
	selected *pp.Player
}

func InitialWelcomeModel(players *pp.Players) WelcomeModel {
	return WelcomeModel{
		choices: players.Players,
		cursor:  0,
	}
}

func (wm WelcomeModel) Init() tea.Cmd {
	return nil
}

func (wm WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return wm, tea.Quit
		case "up", "k":
			if wm.cursor > 0 {
				wm.cursor--
			}
		case "down", "j":
			if wm.cursor < len(wm.choices)-1 {
				wm.cursor++
			}
		case "enter", " ":
			wm.selected = wm.choices[wm.cursor]
		}
	}
	return wm, nil
}

func (wm WelcomeModel) View() string {
	s := "Scroll through players\n\n"

	for i, player := range wm.choices {
		if i == wm.cursor {
			s += " > "
		}
		s += printer.Print(player)
	}

	return s
}
