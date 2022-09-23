package client

import (
	"context"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ADPClient struct {
	draftPickClient player_proto.DraftPickServiceClient
}

func (c *ADPClient) GetPlayers(ctx context.Context, in *emptypb.Empty) (*player_proto.Players, error) {

	return c.draftPickClient.GetPlayers(ctx, in)
}
