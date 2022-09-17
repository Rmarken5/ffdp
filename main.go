package main

import (
	"fmt"
	"github.com/rmarken5/ffdp/service"
	"github.com/rmarken5/ffdp/web-scraper"
	"net/http"
)

func main() {
	scraperSvc := web_scraper.WebScraperImpl{
		Client: http.Client{},
	}

	adpPlayers, err := scraperSvc.GetAverageDraftPickList(web_scraper.FantasySharksCurrentADPURL)
	if err != nil {
		panic(err)
	}
	playerStats, err := scraperSvc.GetTotalPlayerPoints(web_scraper.FantasySharksPreviousYearPointsURL)
	if err != nil {
		panic(err)
	}
	adpMap := service.ConvertADPSliceToMap(adpPlayers.Players)
	statMap := service.ConvertPlayerStatsSliceToMap(playerStats)

	players, notFound := service.BuildPlayerSlice(adpMap, statMap)
	for _, player := range players {
		fmt.Println(player)
	}

	fmt.Printf("\n\n\n\n\n")
	fmt.Println("NotFound!")

	for _, player := range notFound {
		fmt.Println(player)
	}
}
