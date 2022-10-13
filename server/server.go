package server

import (
	"context"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	"github.com/rmarken5/ffdp/server/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ADPServer struct {
	logic service.Logic
	player_proto.UnimplementedDraftPickServiceServer
}

func NewADPServer(logic service.Logic) *ADPServer {
	return &ADPServer{
		logic: logic,
	}
}

func (s ADPServer) GetPlayersByPreviousYearPoints(ctx context.Context, _ *emptypb.Empty) (*player_proto.Players, error) {
	players, err := s.logic.GetPlayersByPreviousYearPoints(ctx)
	if err != nil {
		return nil, err
	}

	return &players, err
}

func (s ADPServer) GetPlayersByCurrentYearProjections(ctx context.Context, _ *emptypb.Empty) (*player_proto.Players, error) {
	players, err := s.logic.GetPlayersByPreviousYearPoints(ctx)
	if err != nil {
		return nil, err
	}

	return &players, err
}
