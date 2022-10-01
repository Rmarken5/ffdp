package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rmarken5/ffdp/charm/models"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure())
	if err != nil {
		fmt.Println("unable to dial.")
		panic(err)
	}
	defer conn.Close()
	pbClient := player_proto.NewDraftPickServiceClient(conn)

	players, _ := pbClient.GetPlayers(context.Background(), &emptypb.Empty{})

	program := tea.NewProgram(models.InitialWelcomeModel(players, 10), tea.WithAltScreen())
	if err := program.Start(); err != nil {
		panic(err)
	}
}
