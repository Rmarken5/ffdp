package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
)

type WelcomeModel struct {
	client     player_proto.DraftPickServiceClient
	cursor     int
	MenuItems  []MenuItem
	windowSize int
}
type MenuItem struct {
	Label       string
	CreateModel func(client player_proto.DraftPickServiceClient) tea.Model
}

func InitializeWelcomeModel(client player_proto.DraftPickServiceClient) WelcomeModel {
	return WelcomeModel{
		cursor: 0,
		client: client,
	}
}

func (w WelcomeModel) Init() tea.Cmd {
	return nil
}

func (w WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return w, tea.Quit
		case "up", "k":
			if w.cursor > 0 {
				w.cursor--
			}
		case "down", "j":
			if w.cursor < len(w.MenuItems)-1 {
				w.cursor++
			}
		case "enter", " ":
			m := w.MenuItems[w.cursor].CreateModel(w.client)
			return m, tea.Println("loading")
		}
	}

	return w, nil
}

func (w WelcomeModel) View() string {
	s := "Select an option.\n\n"

	for i, item := range w.MenuItems {
		if w.cursor == i {
			s += "> "
		}
		s += item.Label + "\n"
	}
	return s
}
