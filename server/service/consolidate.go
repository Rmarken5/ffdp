package service

import (
	"errors"
	"fmt"
	"github.com/rmarken5/ffdp/server/web-scraper"
	"hash/fnv"
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

func ConvertADPSliceToMap(players []web_scraper.ADPPlayer) map[uint32]web_scraper.ADPPlayer {
	playerMap := make(map[uint32]web_scraper.ADPPlayer, len(players))
	for _, player := range players {
		playerMap[hash(player.Position+player.FirstName+player.LastName)] = player
	}

	return playerMap
}

func ConvertPlayerStatsSliceToMap(players []web_scraper.PlayerStats) map[uint32]web_scraper.PlayerStats {
	playerMap := make(map[uint32]web_scraper.PlayerStats, len(players))
	for _, player := range players {
		playerMap[hash(player.Position+player.FirstName+player.LastName)] = player
	}
	return playerMap
}

func BuildPlayerSlice(draftPicks map[uint32]web_scraper.ADPPlayer, playerStats map[uint32]web_scraper.PlayerStats) (players []Player, notMatched []web_scraper.ADPPlayer) {
	for key, val := range draftPicks {
		if stats, ok := playerStats[key]; ok {
			player, err := DraftPickAndStatsToPlayer(val, stats, 170)
			if err != nil {
				fmt.Println(err)
			}
			players = append(players, player)
		} else {
			notMatched = append(notMatched, val)
		}
	}
	return players, notMatched
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
