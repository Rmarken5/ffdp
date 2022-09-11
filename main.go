package main

import (
	"github.com/rmarken5/ffdp/web-scraper"
	"net/http"
)

func main() {
	scraperSvc := web_scraper.WebScraperImpl{
		Client: http.Client{},
	}

	players, err := scraperSvc.GetPlayersFromSource(web_scraper.FantasySharksURL)
	if err != nil {
		panic(err)
	}

	scraperSvc.PrintPlayers(players)
}
