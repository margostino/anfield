package scrapper

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/domain"
	"strings"
	"time"
)

var browser *rod.Browser

func Initialize() {
	browser = rod.New().MustConnect()
}

func Browser() *rod.Browser {
	return browser
}

func GetFinishedResults() []string {
	return GetUrlsResult(context.BATCH)
}

func GetInProgressResults() []string {
	return GetUrlsResult(context.REALTIME)
}

func GetUrlsResult(mode string) []string {
	var urls []string
	var url, selector, pattern string

	if mode == context.REALTIME {
		selector = context.Config().Fixtures.Selector
		pattern = context.Config().Fixtures.Pattern
		url = context.Config().Source.Url + context.Config().Fixtures.Url
	} else {
		selector = context.Config().Results.Selector
		pattern = context.Config().Results.Pattern
		url = context.Config().Source.Url + context.Config().Results.Url
	}

	resultsPage := browser.MustPage(url)
	elements := resultsPage.MustElement(selector).MustElements(pattern)

	for _, element := range elements {
		status := element.MustText()
		if mode == context.BATCH || (mode == context.REALTIME && inProgress(status)) {
			urls = append(urls, element.MustProperty("href").String())
		}
	}

	return urls
}

func inProgress(status string) bool {
	prefix := strings.Split(status, "\n")[0]
	return common.IsTimeCounter(prefix)
}

func GetEventDate(url string) string {
	infoUrl := url + context.Config().Info.Params
	selector := context.Config().Info.Selector
	infoPage := browser.MustPage(infoUrl)
	startTimeDetail := infoPage.MustElement(selector).MustText()
	startTime := strings.Split(startTimeDetail, "\n")[0]
	day := strings.Split(startTime, " ")[0]
	month := strings.Split(startTime, " ")[1]
	year := strings.Split(startTime, " ")[2]
	normalizedStartTime := fmt.Sprintf("%s-%s-%s", year, common.NormalizeMonth(month), common.NormalizeDay(day))
	eventDate, _ := time.Parse("2006-01-02", normalizedStartTime)
	return eventDate.String()
}

func GetFormation(raw string) []string {
	players := make([]string, 0)
	values := strings.Split(raw, "\n")
	for _, value := range values {
		if !common.IsFormationNumber(value) {
			players = append(players, value)
		}
	}
	return players
}

func GetSubstitutes(elements rod.Elements) ([]string, []string) {
	parseSubstitute := false
	players := make([]string, 0)
	normalizedPlayers := make([]string, 0)
	homeSubstitutes := make([]string, 0)
	awaySubstitutes := make([]string, 0)

	for _, element := range elements {
		value := element.MustText()
		if value == "SUBSTITUTE PLAYERS" {
			parseSubstitute = true
		} else if parseSubstitute && !common.IsTimeCounter(value) && !common.InSlice(value, players) {
			players = strings.Split(value, "\n")
			break
		} else if value == "COACHES" {
			break
		}
	}

	for _, value := range players {
		if !common.IsTimeCounter(value) {
			normalizedPlayers = append(normalizedPlayers, value)
		}
	}

	for i, player := range normalizedPlayers {
		if common.Even(i) {
			homeSubstitutes = append(homeSubstitutes, player)
		} else {
			awaySubstitutes = append(awaySubstitutes, player)
		}
	}

	return homeSubstitutes, awaySubstitutes
}

func GetLineups(url string) (*domain.Team, *domain.Team) {
	lineupsUrl := url + context.Config().Lineups.Params
	page := browser.MustPage(lineupsUrl)
	homeTeamSelector := context.Config().Lineups.HomeTeamSelector
	awayTeamSelector := context.Config().Lineups.AwayTeamSelector
	homeTeamName := page.MustElement(homeTeamSelector).MustText()
	awayTeamName := page.MustElement(awayTeamSelector).MustText()
	homeFormSelector := context.Config().Lineups.HomeSelector
	rawHomeFormation := page.MustElement(homeFormSelector).MustText()
	homeFormation := GetFormation(rawHomeFormation)
	awayFormSelector := context.Config().Lineups.HomeSelector
	rawAwayFormation := page.MustElement(awayFormSelector).MustText()
	awayFormation := GetFormation(rawAwayFormation)
	substituteSelector := context.Config().Lineups.SubstituteSelector
	rawSubstitutes := page.MustElements(substituteSelector)
	homeSubstitutes, awaySubstitutes := GetSubstitutes(rawSubstitutes)
	homeTeam := domain.Team{
		Name:              homeTeamName,
		Form:              homeFormation,
		SubstitutePlayers: homeSubstitutes,
	}
	awayTeam := domain.Team{
		Name:              awayTeamName,
		Form:              awayFormation,
		SubstitutePlayers: awaySubstitutes,
	}
	return &homeTeam, &awayTeam
}

func GetMetadata(url string) *domain.Metadata {
	eventDate := GetEventDate(url)
	homeTeam, awayTeam := GetLineups(url)
	return &domain.Metadata{
		Date:     eventDate,
		HomeTeam: homeTeam,
		AwayTeam: awayTeam,
	}
}

// GetEvents TODO: read events as unbounded streams or until conditions (e.g. 90' time, message pattern, etc)
func GetEvents(url string) *[]string {
	moreCommentSelector := context.Config().Commentary.MoreCommentsSelector
	commentSelector := context.Config().Commentary.Selector
	summaryPage := browser.MustPage(url)
	summaryPage.MustElement(moreCommentSelector).MustElement("*").MustClick()
	rawEvents := summaryPage.MustElement(commentSelector).MustText()
	events := strings.Split(rawEvents, "\n")
	return &events
}
