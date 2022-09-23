package main

import (
	"fmt"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	"github.com/rmarken5/ffdp/server"
	"github.com/rmarken5/ffdp/server/service"
	web_scraper "github.com/rmarken5/ffdp/server/web-scraper"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "ffdp", 0)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8001))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	client := http.Client{}
	webScraper := web_scraper.WebScraperImpl{
		Client: client,
	}
	logic := service.NewLogicImpl(&webScraper, logger)
	adpServer := server.NewADPServer(logic)

	player_proto.RegisterDraftPickServiceServer(grpcServer, adpServer)
	fmt.Println("Server up")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal(err)
	}

}
