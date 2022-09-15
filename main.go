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




	playerStats, _ := scraperSvc.GetTotalPlayerPoints(web_scraper.FantasySharksPreviousYearPointsURL)
/*	scraperSvc.PrintPlayerStatLineHeader()
	for _, player := range playerStats {
		scraperSvc.PrintPlayerStatLine(player)

	}*/

	adpMap := service.ConvertADPSliceToMap(adpPlayers.Players)
	statMap := service.ConvertPlayerStatsSliceToMap(playerStats)

	for key, val := range adpMap {
		fmt.Println("_____________________________________________")
		fmt.Println(val)
		fmt.Println(statMap[key])
		fmt.Println("_____________________________________________")
	}

}
