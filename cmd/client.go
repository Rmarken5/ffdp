package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rmarken5/ffdp/charm/models"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure())
	if err != nil {
		fmt.Println("unable to dial.")
		panic(err)
	}
	defer conn.Close()

	pbClient := player_proto.NewDraftPickServiceClient(conn)
	wm := models.InitializeWelcomeModel(pbClient)

	previousYearPoints := models.MenuItem{
		Label:       "Previous Year ADP Vs. Previous Year Projected Points",
		CreateModel: models.InitializePreviousVsProjected,
	}

	currentYearProjections := models.MenuItem{
		Label:       "Previous Year ADP Vs. Current Year Projected Points",
		CreateModel: models.InitializePreviousVsProjected,
	}
	wm.MenuItems = []models.MenuItem{previousYearPoints, currentYearProjections}

	program := tea.NewProgram(wm, tea.WithAltScreen())
	if err := program.Start(); err != nil {
		panic(err)
	}
}
