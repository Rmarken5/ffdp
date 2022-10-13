package service

import (
	"context"
	"fmt"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	ws "github.com/rmarken5/ffdp/server/web-scraper"
	"log"
	"sort"
)

type Logic interface {
	GetPlayersByPreviousYearPoints(ctx context.Context) (player_proto.Players, error)
	GetPlayersByCurrentYearProjections(ctx context.Context) (player_proto.Players, error)
}

type LogicImpl struct {
	webScraper ws.WebScraper
	logger     *log.Logger
}

func NewLogicImpl(webScraper ws.WebScraper, logger *log.Logger) *LogicImpl {
	return &LogicImpl{
		webScraper: webScraper,
		logger:     logger,
	}
}

func (logic *LogicImpl) GetPlayersByPreviousYearPoints(ctx context.Context) (player_proto.Players, error) {
	adpList, err := logic.webScraper.GetAverageDraftPickList(ws.FantasySharksCurrentADPURL)
	if err != nil {
		logic.logger.Println("error in retrieving adp list.")
		return player_proto.Players{}, fmt.Errorf("error in retrieving adp list")
	}
	playerPoints, err := logic.webScraper.GetTotalPlayerPointsLastYear(ws.FantasySharksPreviousYearPointsURL)
	if err != nil {
		logic.logger.Println("error in retrieving point list.")
		return player_proto.Players{}, fmt.Errorf("error in retrieving point list")
	}

	adpMap := ConvertADPSliceToMap(adpList.Players)
	playerPointsMap := ConvertPlayerStatsSliceToMap(playerPoints)
	players, _ := DraftPicksAndStatsToPlayerProtos(adpMap, playerPointsMap, 170)
	sort.Sort(ByLastNameDesc(players))
	proto := player_proto.Players{Players: players}
	return proto, nil
}

func (logic *LogicImpl) GetPlayersByCurrentYearProjections(ctx context.Context) (player_proto.Players, error) {
	adpList, err := logic.webScraper.GetAverageDraftPickList(ws.FantasySharksCurrentADPURL)
	if err != nil {
		logic.logger.Println("error in retrieving adp list.")
		return player_proto.Players{}, fmt.Errorf("error in retrieving adp list")
	}
	playerPoints, err := logic.webScraper.GetTotalPlayerPointsProjected(ws.FantasySharksCurrentYearProjectedPointsURL)
	if err != nil {
		logic.logger.Println("error in retrieving point list.")
		return player_proto.Players{}, fmt.Errorf("error in retrieving point list")
	}

	adpMap := ConvertADPSliceToMap(adpList.Players)
	playerPointsMap := ConvertPlayerStatsSliceToMap(playerPoints)
	players, _ := DraftPicksAndStatsToPlayerProtos(adpMap, playerPointsMap, 170)
	sort.Sort(ByLastNameDesc(players))
	proto := player_proto.Players{Players: players}
	return proto, nil
}
