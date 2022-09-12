package main

import (
	"github.com/rmarken5/ffdp/web-scraper"
	"net/http"
)

func main() {
	scraperSvc := web_scraper.WebScraperImpl{
		Client: http.Client{},
	}

	/*players, err := scraperSvc.GetPlayersFromSource(web_scraper.FantasySharksCurrentADPURL)
	if err != nil {
		panic(err)
	}

	scraperSvc.PrintAverageDraftPickList(players)*/

	playerStats, _ := scraperSvc.GetTotalPlayerPoints(web_scraper.FantasySharksPreviousYearPointsURL)
	scraperSvc.PrintPlayerStatLineHeader()
	for _, player := range playerStats {
		scraperSvc.PrintPlayerStatLine(player)

	}
}
