package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rmarken5/ffdp/client/printer"
	pp "github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
)

type PlayerModel struct {
	allPlayers    []*pp.Player
	playerPage    []*pp.Player
	pageSize      int
	currentPage   int
	cursor        int
	numberOfPages int
	selected      *pp.Player
}

func InitializePlayerModel(allPlayers []*pp.Player, pageSize int) PlayerModel {

	numberOfPages := int(math.Ceil(float64(len(allPlayers)) / float64(pageSize)))
	playerPage := allPlayers[:pageSize]

	return PlayerModel{
		allPlayers:    allPlayers,
		playerPage:    playerPage,
		pageSize:      pageSize,
		currentPage:   1,
		cursor:        0,
		numberOfPages: numberOfPages,
	}
}

func (wm PlayerModel) Init() tea.Cmd {
	return nil
}

func (wm PlayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return wm, tea.Quit
		case "up", "k":
			if wm.cursor > 0 {
				wm.cursor--
			} else if wm.currentPage-1 > 0 {
				wm.cursor = wm.pageSize - 1
				wm.currentPage--
				wm.playerPage = wm.allPlayers[(wm.currentPage-1)*wm.pageSize : wm.currentPage*wm.pageSize]
			}
		case "down", "j":
			if wm.cursor < len(wm.playerPage)-1 {
				wm.cursor++
			} else if wm.currentPage < wm.numberOfPages {
				wm.cursor = 0
				wm.currentPage++
				if wm.currentPage == wm.numberOfPages {
					wm.playerPage = wm.allPlayers[(wm.currentPage-1)*wm.pageSize:]
				} else {
					wm.playerPage = wm.allPlayers[(wm.currentPage-1)*wm.pageSize : wm.currentPage*wm.pageSize]
				}
			}
		case "enter", " ":
			wm.selected = wm.playerPage[wm.cursor]
		}
	}
	return wm, nil
}

func (wm PlayerModel) View() string {
	s := "Scroll through players\n\n"

	for i, player := range wm.playerPage {
		if i == wm.cursor {
			s += " > "
		}
		s += printer.Print(player)
	}
	s += fmt.Sprintf("\nPage %d of %d\n\n", wm.currentPage, wm.numberOfPages)

	if wm.selected != nil {
		s += fmt.Sprintf("Currently Selected Player:\n%s\n", printer.Print(wm.selected))
	}

	return s
}

func InitPlayerModel(client pp.DraftPickServiceClient, pageSize int) tea.Model {
	players, _ := client.GetPlayers(context.Background(), &emptypb.Empty{})
	pm := InitializePlayerModel(players.Players, pageSize)
	return pm
}
