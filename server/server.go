package server

import (
	"context"
	"fmt"
	"github.com/rmarken5/ffdp/protobuf/proto_files/adp_server"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	"github.com/rmarken5/ffdp/server/service"
	web_scraper "github.com/rmarken5/ffdp/server/web-scraper"
	"google.golang.org/appengine/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ADPServer struct {
	scraper web_scraper.WebScraper
	adp_server.UnimplementedDraftPickServiceServer
}

func (s ADPServer) GetPlayers(ctx context.Context, _ *emptypb.Empty) (*player_proto.Players, error) {
	//TODO: refactor all this into a service method.

	adpList, err := s.scraper.GetAverageDraftPickList(web_scraper.FantasySharksCurrentADPURL)
	if err != nil {
		log.Errorf(ctx, "error in retrieving adp list.")
		return &player_proto.Players{}, fmt.Errorf("error in retrieving adp list")
	}
	playerPoints, err := s.scraper.GetTotalPlayerPoints(web_scraper.FantasySharksPreviousYearPointsURL)
	if err != nil {
		log.Errorf(ctx, "error in retrieving point list.")
		return &player_proto.Players{}, fmt.Errorf("error in retrieving point list")
	}

	adpMap := service.ConvertADPSliceToMap(adpList.Players)
	playerPointsMap := service.ConvertPlayerStatsSliceToMap(playerPoints)

	players, _ := service.DraftPicksAndStatsToPlayerProtos(adpMap, playerPointsMap, 170)

	proto := player_proto.Players{Players: players}

	return &proto, nil

}
