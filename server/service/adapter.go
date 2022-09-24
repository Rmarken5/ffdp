package service

import (
	"errors"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
	web_scraper "github.com/rmarken5/ffdp/server/web-scraper"
	"strconv"
)

func DraftPicksAndStatsToPlayerProtos(draftPicks map[uint32]web_scraper.ADPPlayer,
	playerStats map[uint32]web_scraper.PlayerStats,
	picksInDraft int) (players []*player_proto.Player, unmatched []web_scraper.ADPPlayer) {
	for key, pick := range draftPicks {
		if stat, ok := playerStats[key]; ok {
			player, err := DraftPickAndStatsToPlayerProto(pick, stat, picksInDraft)
			if err != nil {
				// TODO: Make a server instance so logger can be used.
			}
			players = append(players, &player)
		} else {
			unmatched = append(unmatched, pick)
		}
	}
	return players, unmatched

}

func DraftPickAndStatsToPlayerProto(draftPick web_scraper.ADPPlayer, stat web_scraper.PlayerStats, picksInDraft int) (player_proto.Player, error) {
	draftKey := PlayerKey{
		Position: draftPick.Position,
		Team:     draftPick.Team,
		LastName: draftPick.LastName,
	}
	statKey := PlayerKey{
		Position: stat.Position,
		Team:     stat.Team,
		LastName: stat.LastName,
	}

	if isEqual := draftKey.IsDraftPickEqualToStat(statKey); !isEqual {
		return player_proto.Player{}, errors.New("player stats and draft record do not match")
	}

	convertedADP, err := strconv.ParseFloat(draftPick.ADP, 32)
	if err != nil {
		return player_proto.Player{}, err
	}

	return player_proto.Player{
		Rank:       draftPick.Rank,
		ID:         draftPick.ID,
		FullName:   draftPick.FullName,
		FirstName:  draftPick.FirstName,
		LastName:   draftPick.LastName,
		Position:   draftPick.Position,
		Team:       draftPick.Team,
		Bye:        draftPick.Bye,
		ADP:        draftPick.ADP,
		PointTotal: int32(stat.Points),
		Value:      float32(stat.Points) / float32(picksInDraft) / float32(convertedADP),
	}, nil
}
