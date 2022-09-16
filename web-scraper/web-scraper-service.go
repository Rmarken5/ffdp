package web_scraper

import (
	"encoding/xml"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/google/brotli/go/cbrotli"
	"golang.org/x/net/html/charset"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

const (
	FantasySharksCurrentADPURL         = "https://www.fantasysharks.com/apps/bert/forecasts/adp.php?Position=97&xml=1&adpsort=99&Segment=746"
	FantasySharksPreviousYearPointsURL = "https://www.fantasysharks.com/apps/bert/stats/points.php?League=-1&Position=99&scoring=3&Segment=717"
	ADPPlayerStringFormat              = "|%-4s|%-6s|%-35s|%-8s|%-4s|%-3s|%-5d|%-4s|%-5s|%-5s|%-6s|\n"
	PlayerStatFormat                   = "|%-4s|%-35s|%-4s|%-8s|%-10d|%-8d|%-9d|%-8d|%-4d|%-9d|%-8d|%-6d|\n"
	PlayerStatFormatHeader             = "|%-4s|%-35s|%-4s|%-8s|%-10s|%-8s|%-9s|%-8s|%-4s|%-9s|%-8s|%-6s|\n"
)

type WebScraper interface {
	GetAverageDraftPickList(url string) (Adp, error)
	GetTotalPlayerPoints(url string) ([]PlayerStats, error)
	PrintPlayers(Adp)
}

type WebScraperImpl struct {
	http.Client
}
type Adp struct {
	XMLName xml.Name    `xml:"adp"`
	Text    string      `xml:",chardata"`
	Players []ADPPlayer `xml:"player"`
}
type ADPPlayer struct {
	Text      string `xml:",chardata"`
	Rank      string `xml:"Rank,attr"`
	ID        string `xml:"ID,attr"`
	FullName  string `xml:"Name,attr"`
	FirstName string
	LastName  string
	Position  string `xml:"Position,attr"`
	Team      string `xml:"Team,attr"`
	Bye       string `xml:"Bye,attr"`
	ADP       string  `xml:"ADP,attr"`
	StdDev    string `xml:"StdDev,attr"`
	MFL       string `xml:"MFL,attr"`
	RTS       string `xml:"RTS,attr"`
	FFCalc    string `xml:"FFCalc,attr"`
}

type PlayerStats struct {
	Rank           string
	Name           string
	LastName       string
	FirstName      string
	Team           string
	Position       string
	PassYrds       int16
	PassTDs        int8
	RushYrds       int16
	RushTDs        int8
	Recs           int8
	RecYrds        int16
	RecTDs         int8
	FieldGoalsMade int8
	PointsAgainst  int16
	Tackles        int16
	Points         int16
}

func (w *WebScraperImpl) GetAverageDraftPickList(myUrl string) (Adp, error) {
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

func (w *WebScraperImpl) GetTotalPlayerPoints(url string) ([]PlayerStats, error) {
	var playerStats []PlayerStats
	collector := colly.NewCollector()

	collector.OnRequest(func(req *colly.Request) {
		req.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_"+strconv.Itoa(rand.Intn(15-9)+9)+"_1) AppleWebKit/531.36 (KHTML, like Gecko) Chrome/"+strconv.Itoa(rand.Intn(79-70)+70)+".0.3945.130 Safari/531.36")
		req.Headers.Set("Accept-Encoding", "gzip")
		req.Headers.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
		req.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	})

	collector.OnHTML("#toolData > tbody > tr", func(h *colly.HTMLElement) {
		player := new(PlayerStats)

		h.ForEachWithBreak("td", func(i int, element *colly.HTMLElement) bool {
			switch i {
			case 0:
				if element.Text == "" { // Header row.
					return false
				}
				player.Rank = element.Text
				break
			case 1:
				player.Name = element.Text
				break
			case 2:
				player.Team = element.Text
				break
			case 3:
				player.Position = element.Text
				break
			case 4:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.PassYrds = int16(number)
				break
			case 5:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.PassTDs = int8(number)
				break
			case 6:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.RushYrds = int16(number)
				break
			case 7:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.RushTDs = int8(number)
				break
			case 8:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.Recs = int8(number)
				break
			case 9:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.RecYrds = int16(number)
				break
			case 10:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.RecTDs = int8(number)
				break
			case 11:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.FieldGoalsMade = int8(number)
				break
			case 12:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.PointsAgainst = int16(number)
				break
			case 13:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.Tackles = int16(number)
				break
			case 14:
				number, err := strconv.Atoi(element.Text)
				if err != nil {
					fmt.Println("cannot convert text")
				}
				player.Points = int16(number)
				massagePlayerStats(player)
				playerStats = append(playerStats, *player)
				break
			}
			return true
		})
	})

	collector.Visit(url)

	return playerStats, nil
}

func (w *WebScraperImpl) PrintAverageDraftPickList(adp Adp) {
	players := adp.Players
	w.PrintHeader()
	for _, player := range players {
		w.PrintPlayer(player)
	}
}

func (w *WebScraperImpl) PrintHeader() {
	fmt.Printf(ADPPlayerStringFormat, "Rank", "ID", "Name", "Position", "Team", "Bye", "ADP", "SDEV", "MFL", "RTS", "FFCalc")
}

func (w *WebScraperImpl) PrintPlayer(player ADPPlayer) {
	fmt.Printf(ADPPlayerStringFormat, player.Rank, player.ID, player.FullName, player.Position, player.Team, player.Bye, player.ADP, player.StdDev, player.MFL, player.RTS, player.FFCalc)
}

func (w *WebScraperImpl) PrintPlayerStatLineHeader() {
	fmt.Printf(PlayerStatFormatHeader, "Rank", "Name", "Team", "Position", "Pass Yards", "Pass TD", "Rush Yads", "Rush TD", "Recs", "Rec Yards", "Rec TD", "Points")
}

func (w *WebScraperImpl) PrintPlayerStatLine(player PlayerStats) {
	fmt.Printf(PlayerStatFormat, player.Rank, player.Name, player.Team, player.Position, player.PassYrds, player.PassTDs, player.RushYrds, player.RushTDs, player.Recs, player.RecYrds, player.RecTDs, player.Points)
}

func massagePlayerStats(player *PlayerStats) {
	if nameParts := strings.Split(player.Name, ","); len(nameParts) > 1 {
		player.LastName = nameParts[0]
		player.FirstName = nameParts[1]
	}
}
func massagePlayerData(player *ADPPlayer) {
	if nameParts := strings.Split(player.FullName, ","); len(nameParts) > 1 {
		player.LastName = nameParts[0]
		player.FirstName = nameParts[1]
	}
}

