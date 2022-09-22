package main

import (
	"fmt"
	"github.com/rmarken5/ffdp/protobuf/proto_files/adp_server"
	"github.com/rmarken5/ffdp/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", "8001"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	adpServer := server.ADPServer{}
	adp_server.RegisterDraftPickServiceServer(grpcServer, adpServer)
	grpcServer.Serve(lis)
}
