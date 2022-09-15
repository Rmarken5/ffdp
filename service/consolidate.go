package service

import (
	"errors"
	"github.com/rmarken5/ffdp/web-scraper"
	"strconv"
)

type Player struct {
	DraftPick  web_scraper.ADPPlayer
	Key        PlayerKey
	PointTotal int16
	Value      float32
}

type PlayerKey struct {
	Position string
	Team     string
	LastName string
}

func (pk PlayerKey) IsDraftPickEqualToStat(otherKey PlayerKey) bool {
	return pk.Position == otherKey.Position && pk.Team == otherKey.Team && pk.LastName == otherKey.LastName
}

func DraftPickAndStatsToPlayer(draftPick web_scraper.ADPPlayer, stat web_scraper.PlayerStats, picksInDraft int) (Player, error) {
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
		return Player{}, errors.New("player stats and draft record do not match")
	}

	convertedADP, err := strconv.ParseFloat(draftPick.ADP, 32)
	if err != nil {
		return Player{}, err
	}

	return Player{
		DraftPick:  draftPick,
		PointTotal: stat.Points,
		Key:        draftKey,
		Value:      float32(stat.Points) / float32(picksInDraft) / float32(convertedADP),
	}, nil
}

func ConvertADPSliceToMap(players []web_scraper.ADPPlayer) map[PlayerKey]web_scraper.ADPPlayer {
	playerMap := make(map[PlayerKey]web_scraper.ADPPlayer, len(players))
	for _, player := range players {
		key := PlayerKey{
			Position: player.Position,
			Team:     player.Team,
			LastName: player.LastName,
		}
		playerMap[key] = player
	}
	return playerMap
}

func ConvertPlayerStatsSliceToMap(players []web_scraper.PlayerStats) map[PlayerKey]web_scraper.PlayerStats {
	playerMap := make(map[PlayerKey]web_scraper.PlayerStats, len(players))
	for _, player := range players {
		key := PlayerKey{
			Position: player.Position,
			Team:     player.Team,
			LastName: player.LastName,
		}
		playerMap[key] = player
	}
	return playerMap
}
