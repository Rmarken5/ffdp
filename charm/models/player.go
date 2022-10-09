package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rmarken5/ffdp/client/printer"
	pp "github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	"github.com/rmarken5/ffdp/server/service"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"sort"
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

func InitializePlayerModel(client pp.DraftPickServiceClient, initHeight int) tea.Model {
	players, _ := client.GetPlayers(context.Background(), &emptypb.Empty{})
	pageSize := initHeight - 9

	return PlayerModel{
		allPlayers:    players.Players,
		currentPage:   1,
		cursor:        0,
		pageSize:      pageSize,
		numberOfPages: calculateNumberOfPages(pageSize, players.Players),
		playerPage:    calculatePlayerPage(1, pageSize, players.Players),
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
				p.playerPage = calculatePlayerPage(p.currentPage, p.pageSize, p.allPlayers)
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
					p.playerPage = calculatePlayerPage(p.currentPage, p.pageSize, p.allPlayers)
				}
			}
		case "enter", " ":
			p.selected = p.playerPage[p.cursor]
		case "s":
			return InitializeSortModel(createSortMenuItems(), &p), tea.EnterAltScreen

		}
	case tea.WindowSizeMsg:
		p.updatePlayerModel(msg.Height - 10)
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
	p.pageSize = height
	p.numberOfPages = calculateNumberOfPages(p.pageSize, p.allPlayers)
	p.playerPage = calculatePlayerPage(p.currentPage, p.pageSize, p.allPlayers)
}

func calculateNumberOfPages(height int, v []*pp.Player) int {
	return int(math.Ceil(float64(len(v)) / float64(height)))
}

func calculatePlayerPage(currentPage, pageSize int, players []*pp.Player) []*pp.Player {
	return players[(currentPage-1)*pageSize : currentPage*pageSize]
}

func createSortMenuItems() []SortMenuItem {
	byLastNameDesc := NewSortMenuItem("By Last Name Desc", func(players []*pp.Player) {
		sort.Sort(service.ByLastNameDesc(players))
	})
	byLastNameAsc := NewSortMenuItem("By Last Name Asc", func(players []*pp.Player) {
		sort.Sort(service.ByLastNameAsc(players))
	})
	byValueDesc := NewSortMenuItem("By Value Desc", func(players []*pp.Player) {
		sort.Sort(service.ByValueDesc(players))
	})
	byValueAsc := NewSortMenuItem("By Value Asc", func(players []*pp.Player) {
		sort.Sort(service.ByValueAsc(players))
	})
	sortMenuItems := []SortMenuItem{byLastNameAsc, byLastNameDesc, byValueAsc, byValueDesc}
	return sortMenuItems

}
