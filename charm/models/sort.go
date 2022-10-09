package models

import (
	tea "github.com/charmbracelet/bubbletea"
	pp "github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
)

type SortModel struct {
	cursor        int
	MenuItems     []SortMenuItem
	modelToReturn *PlayerModel
}

type SortMenuItem struct {
	Label  string
	SortBy func([]*pp.Player)
}

func InitializeSortModel(menuItems []SortMenuItem, modelToReturn *PlayerModel) SortModel {
	return SortModel{
		cursor:        0,
		MenuItems:     menuItems,
		modelToReturn: modelToReturn,
	}
}

func NewSortMenuItem(label string, sortImpl func([]*pp.Player)) SortMenuItem {
	return SortMenuItem{
		Label:  label,
		SortBy: sortImpl,
	}
}

func (s SortModel) Init() tea.Cmd {
	return nil
}

func (s SortModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return nil, tea.Quit
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			if s.cursor < len(s.MenuItems)-1 {
				s.cursor++
			}
		case "enter", " ":
			m := s.MenuItems[s.cursor]

			m.SortBy(s.modelToReturn.allPlayers)

			return *s.modelToReturn, tea.Println("loading")
		}
	}

	return s, nil
}

func (s SortModel) View() string {
	st := "Select an option.\n\n"

	for i, item := range s.MenuItems {
		if s.cursor == i {
			st += "> "
		}
		st += item.Label + "\n"

	}

	return st
}
