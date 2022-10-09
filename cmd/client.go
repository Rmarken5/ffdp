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

	playerMenuItem := models.MenuItem{
		Label:       "Previous Year ADP Vs. Projected Points",
		CreateModel: models.InitializePlayerModel,
	}
	wm.MenuItems = append(wm.MenuItems, playerMenuItem)

	program := tea.NewProgram(wm, tea.WithAltScreen())
	if err := program.Start(); err != nil {
		panic(err)
	}
}
