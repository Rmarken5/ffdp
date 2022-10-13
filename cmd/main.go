package main

import (
	"fmt"
	"github.com/rmarken5/ffdp/server/web-scraper"
	"net/http"
)

func main() {
	scraperSvc := web_scraper.WebScraperImpl{
		Client: http.Client{},
	}

	/*adpPlayers, err := scraperSvc.GetAverageDraftPickList(web_scraper.FantasySharksCurrentADPURL)
	if err != nil {
		panic(err)
	}*/
	playerStats, err := scraperSvc.GetTotalPlayerPointsProjected(web_scraper.FantasySharksCurrentYearProjectedPointsURL)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", playerStats)
	/*adpMap := service.ConvertADPSliceToMap(adpPlayers.Players)
	statMap := service.ConvertPlayerStatsSliceToMap(playerStats)

	players, notFound := service.BuildPlayerSlice(adpMap, statMap)
	for _, player := range players {
		fmt.Println(player)
	}

	fmt.Printf("\n\n\n\n\n")
	fmt.Println("NotFound!")

	for _, player := range notFound {
		fmt.Println(player)
	}*/
}
