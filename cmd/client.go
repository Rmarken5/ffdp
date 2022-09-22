package main

import (
	"context"
	"fmt"
	"github.com/rmarken5/ffdp/protobuf/proto_files/adp_server"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:8001")
	if err != nil {
		fmt.Println("unable to dial.")
		panic(err)
	}
	defer conn.Close()
	pbClient := adp_server.NewDraftPickServiceClient(conn)

	players, _ := pbClient.GetPlayers(context.Background(), &emptypb.Empty{})

	fmt.Print(players)

}
