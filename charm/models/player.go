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

func InitializePlayerModel(allPlayers []*pp.Player) PlayerModel {

	return PlayerModel{
		allPlayers:  allPlayers,
		currentPage: 1,
		cursor:      0,
	}
}

func (p PlayerModel) Init() tea.Cmd {
	return nil
}

func (p PlayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return p, tea.Quit
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			} else if p.currentPage-1 > 0 {
				p.cursor = p.pageSize - 1
				p.currentPage--
				p.playerPage = p.allPlayers[(p.currentPage-1)*p.pageSize : p.currentPage*p.pageSize]
			}
		case "down", "j":
			if p.cursor < len(p.playerPage)-1 {
				p.cursor++
			} else if p.currentPage < p.numberOfPages {
				p.cursor = 0
				p.currentPage++
				if p.currentPage == p.numberOfPages {
					p.playerPage = p.allPlayers[(p.currentPage-1)*p.pageSize:]
				} else {
					p.playerPage = p.allPlayers[(p.currentPage-1)*p.pageSize : p.currentPage*p.pageSize]
				}
			}
		case "enter", " ":
			p.selected = p.playerPage[p.cursor]
		}
	case tea.WindowSizeMsg:
		p.updatePlayerModel(msg.Height - 6)
	}

	return p, nil
}

func (p PlayerModel) View() string {
	s := "Scroll through players\n\n"

	for i, player := range p.playerPage {
		if i == p.cursor {
			s += " > "
		}
		s += printer.Print(player)
	}
	s += fmt.Sprintf("\nPage %d of %d\n\n", p.currentPage, p.numberOfPages)

	if p.selected != nil {
		s += fmt.Sprintf("Currently Selected Player:\n%s\n", printer.Print(p.selected))
	}

	return s
}

func (p *PlayerModel) updatePlayerModel(height int) {
	numberOfPages := int(math.Ceil(float64(len(p.allPlayers)) / float64(height)))
	p.numberOfPages = numberOfPages
	p.pageSize = height
	p.playerPage = p.allPlayers[(p.currentPage-1)*p.pageSize : p.currentPage*p.pageSize]
}

func InitPlayerModel(client pp.DraftPickServiceClient) tea.Model {
	players, _ := client.GetPlayers(context.Background(), &emptypb.Empty{})
	pm := InitializePlayerModel(players.Players)
	return pm
}
