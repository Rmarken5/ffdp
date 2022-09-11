package web_scraper

import (
	"encoding/xml"
	"fmt"
	"github.com/google/brotli/go/cbrotli"
	"golang.org/x/net/html/charset"
	"net/http"
	"strings"
)

const FantasySharksURL = "https://www.fantasysharks.com/apps/bert/forecasts/adp.php?Position=97&xml=1&adpsort=99&Segment=746"

type WebScraper interface {
	GetPlayersFromSource(url string) (Adp, error)
	PrintPlayers(Adp)
}

type WebScraperImpl struct {
	http.Client
}
type Adp struct {
	XMLName xml.Name `xml:"adp"`
	Text    string   `xml:",chardata"`
	Players []Player `xml:"player"`
}
type Player struct {
	Text      string `xml:",chardata"`
	Rank      string `xml:"Rank,attr"`
	ID        string `xml:"ID,attr"`
	FullName  string `xml:"Name,attr"`
	FirstName string
	LastName  string
	Position  string `xml:"Position,attr"`
	Team      string `xml:"Team,attr"`
	Bye       string `xml:"Bye,attr"`
	ADP       string `xml:"ADP,attr"`
	StdDev    string `xml:"StdDev,attr"`
	MFL       string `xml:"MFL,attr"`
	RTS       string `xml:"RTS,attr"`
	FFCalc    string `xml:"FFCalc,attr"`
}

func (w *WebScraperImpl) GetPlayersFromSource(myUrl string) (Adp, error) {
	var adp Adp
	request, err := http.NewRequest(http.MethodGet, myUrl, nil)
	if err != nil {
		return Adp{}, err
	}
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-Encoding", "br")
	request.Header.Set("Cookie", "phpbb3_54dir_u=1; phpbb3_54dir_k=; PHPSESSID=925c3f786cb9acd0f85f16ecb21b5734; phpbb3_54dir_sid=035316db15317149f51fb1a5ac2ddef7; FFTools=1855314257")
	request.Header.Set("User-Agent", "PostmanRuntime/7.29.2")
	request.Header.Set("Postman-Token", "67ff4ac2-4ae6-4e83-8939-9926f9df7b7e")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Host", "www.fantasysharks.com")

	if err != nil {
		return Adp{}, err
	}
	request.Close = true
	resp, err := w.Client.Do(request)
	if err != nil {
		fmt.Println("error on do")
		return Adp{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Adp{}, fmt.Errorf("recieved status code: %d", resp.StatusCode)
	}
	reader := cbrotli.NewReader(resp.Body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&adp)

	if err != nil {
		return Adp{}, err
	}

	for i, player := range adp.Players {
		massagePlayerData(&player)
		adp.Players[i] = player
	}

	fmt.Println(adp.Players[0].LastName)

	return adp, nil
}

func massagePlayerData(player *Player) {
	if nameParts := strings.Split(player.FullName, ","); len(nameParts) > 1 {
		player.LastName = nameParts[0]
		player.FirstName = nameParts[1]
	}
}

func (w *WebScraperImpl) PrintPlayers(adp Adp) {
	players := adp.Players
	fmt.Printf("|%-4s|%-6s|%-35s|%-8s|%-4s|%-3s|%-5s|%-4s|%-5s|%-5s|%-6s|\n", "Rank", "ID", "Name", "Position", "Team", "Bye", "ADP", "SDEV", "MFL", "RTS", "FFCalc")
	for _, player := range players {
		fmt.Printf("|%4s|%6s|%35s|%8s|%4s|%3s|%5s|%4s|%5s|%5s|%6s|\n", player.Rank, player.ID, player.FullName, player.Position, player.Team, player.Bye, player.ADP, player.StdDev, player.MFL, player.RTS, player.FFCalc)
	}
}
