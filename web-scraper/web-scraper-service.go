package web_scraper

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
)

type WebScraper interface {
	GetPlayersFromSource(url string) (ADP, error)
	PrintPlayers(ADP)
}

type WebScraperImpl struct {
	http.Client
}
type ADP struct {
	adp     string `xml:"adp"`
	Players []Player
}

type Player struct {
	Rank     string `xml:"Rank"`
	ID       string `xml:"ID"`
	Name     string `xml:"Name"`
	Position string `xml:"Position"`
	Team     string `xml:"Team"`
	Bye      string `xml:"Bye"`
	ADP      string `xml:"ADP"`
	StdDev   string `xml:"StdDev"`
	MFL      string `xml:"MFL"`
	RTS      string `xml:"RTS"`
	FFCalc   string `xml:"FFCalc"`
}

func (w *WebScraperImpl) GetPlayersFromSource(myUrl string) (ADP, error) {
	var players ADP
	request, err := http.NewRequest(http.MethodGet, myUrl, nil)
	if err != nil {
		return ADP{}, err
	}
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-Encoding", "br")
	request.Header.Set("Cookie", "phpbb3_54dir_u=1; phpbb3_54dir_k=; PHPSESSID=925c3f786cb9acd0f85f16ecb21b5734; phpbb3_54dir_sid=035316db15317149f51fb1a5ac2ddef7; FFTools=1855314257")
	request.Header.Set("User-Agent", "PostmanRuntime/7.29.2")
	request.Header.Set("Postman-Token", "67ff4ac2-4ae6-4e83-8939-9926f9df7b7e")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Host", "www.fantasysharks.com")

	if err != nil {
		return ADP{}, err
	}
	request.Close = true
	resp, err := w.Client.Do(request)
	fmt.Println(resp.Header)

	if err != nil {
		fmt.Println("error on do")
		return ADP{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return ADP{}, fmt.Errorf("recieved status code: %d", resp.StatusCode)
	}
	content, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(content))
	if err != nil {
		fmt.Println("error on readall")
		return ADP{}, err
	}

	reader := bytes.NewReader(content)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&players)

	if err != nil {
		return ADP{}, err
	}
	return players, nil
}

func (w *WebScraperImpl) PrintPlayers(adp ADP) {
	players := adp.Players
	fmt.Printf("|%4s|%6s|%35|%8s|%4s|%3s|%5s|%4s|%5s|%5s|%6s|", "Rank", "ID", "Name", "Position", "Team", "Bye", "ADP", "SDEV", "MFL", "RTS", "FFCalc")
	for _, player := range players {
		fmt.Printf("|%4s|%6s|%35|%8s|%4s|%3s|%5s|%4s|%5s|%5s|%6s|", player.Rank, player.ID, player.Name, player.Position, player.Team, player.Bye, player.ADP, player.StdDev, player.MFL, player.RTS, player.FFCalc)
	}
}

const FantasySharksURL = "https://www.fantasysharks.com/apps/bert/forecasts/adp.php?Position=97&xml=1&adpsort=99&Segment=746"
