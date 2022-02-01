package processor

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"log"
	"strings"
	"time"
)

func (a App) getEventDate(url string) time.Time {
	infoUrl := url + a.configuration.Scrapper.InfoParams
	selector := a.configuration.Scrapper.InfoSelector
	startTimeDetail := a.scrapper.GoPage(infoUrl).Text(selector)
	startTime := strings.Split(startTimeDetail, "\n")[0]
	day := strings.Split(startTime, " ")[0]
	month := strings.Split(startTime, " ")[1]
	year := strings.Split(startTime, " ")[2]
	normalizedStartTime := fmt.Sprintf("%s-%s-%s", year, common.NormalizeMonth(month), common.NormalizeDay(day))
	eventDate, _ := time.Parse("2006-01-02", normalizedStartTime)
	return eventDate
}

func (a App) getLineups(url string) (*domain.Team, *domain.Team) {
	lineupsUrl := url + a.configuration.Scrapper.LineupsParams
	homeTeamSelector := a.configuration.Scrapper.HomeTeamSelector
	awayTeamSelector := a.configuration.Scrapper.AwayTeamSelector
	homeFormSelector := a.configuration.Scrapper.HomeSelector
	awayFormSelector := a.configuration.Scrapper.AwaySelector
	substituteSelector := a.configuration.Scrapper.SubstituteSelector

	page := a.scrapper.GoPage(lineupsUrl)
	homeTeamName := page.Text(homeTeamSelector)
	awayTeamName := page.Text(awayTeamSelector)
	rawHomeFormation := page.Text(homeFormSelector)
	rawAwayFormation := page.Text(awayFormSelector)
	rawSubstitutes := page.Elements(substituteSelector)

	homeFormation := getFormation(rawHomeFormation)
	awayFormation := getFormation(rawAwayFormation)
	homeSubstitutes, awaySubstitutes := getSubstitutes(&rawSubstitutes)

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

func getSubstitutes(elements *rod.Elements) ([]domain.Player, []domain.Player) {
	parseSubstitute := false
	players := make([]string, 0)
	normalizedPlayers := make([]string, 0)
	homeSubstitutes := make([]domain.Player, 0)
	awaySubstitutes := make([]domain.Player, 0)

	for _, element := range *elements {
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
			homeSubstitutes = append(homeSubstitutes, *newPlayer(player))
		} else {
			awaySubstitutes = append(awaySubstitutes, *newPlayer(player))
		}
	}

	return homeSubstitutes, awaySubstitutes
}

func getFormation(raw string) []domain.Player {
	players := make([]domain.Player, 0)
	values := strings.Split(raw, "\n")
	for _, value := range values {
		if !common.IsFormationNumber(value) {
			players = append(players, *newPlayer(value))
		}
	}
	return players
}

func newPlayer(name string) *domain.Player {
	return &domain.Player{
		Name:  name,
	}
}

func (a App) produce(url string) {
	go a.matchDate(url)
	go a.lineups(url)
	go a.commentary(url)
}

func (a App) matchDate(url string) {
	a.channels.matchDate[url] <- a.getEventDate(url)
	done()
}

func (a App) lineups(url string) {
	homeTeam, awayTeam := a.getLineups(url)
	lineups := &domain.Lineups{
		HomeTeam: homeTeam,
		AwayTeam: awayTeam,
	}
	a.channels.lineups[url] <- lineups
	done()
}

// TODO: implement proper stop in loop but scan all partial events
func (a App) commentary(url string) {
	sent := 0
	countDown := 0
	endOfEvent := false
	matchInProgress := true
	eventName := strings.Split(url, "/")[7]
	stopFlag := a.configuration.Events.StopFlag
	graceEndTime := a.configuration.Events.GraceEndTime
	commentaryUrl := url + a.configuration.Scrapper.CommentaryParams

	log.Println("START event processing:", eventName)

	for ok := true; ok; ok = matchInProgress {
		if endOfEvent && countDown == 0 {
			time.Sleep(graceEndTime * time.Millisecond)
			countDown += 1
		} else if endOfEvent && countDown == a.configuration.Events.CountDown {
			matchInProgress = false
			a.channels.commentary[url] <- NewFlagCommentary("end")
			break
		}
		rawEvents := a.getEvents(commentaryUrl)
		if rawEvents != nil {
			commentaries := normalize(*rawEvents)
			if sent != len(commentaries) {
				for _, commentary := range commentaries {
					a.channels.commentary[url] <- commentary
					sent += 1
					if strings.Contains(commentary.Comment, stopFlag) {
						endOfEvent = true
					}
				}
			}
		} else {
			log.Println("Match is not started")
			a.channels.commentary[url] <- NewFlagCommentary("not-started")
			break
		}
	}

	log.Println("END event processing:", eventName)

	close(a.channels.commentary[url])
	done()
}

func NewFlagCommentary(flag string) *domain.Commentary {
	return &domain.Commentary{
		Time:    flag,
		Comment: flag,
	}
}

// GetEvents TODO: read events as unbounded streams or until conditions (e.g. 90' time, message pattern, etc)
func (a App) getEvents(url string) *[]string {
	moreCommentSelector := a.configuration.Scrapper.MoreCommentsSelector
	commentSelector := a.configuration.Scrapper.CommentarySelector
	page := a.scrapper.GoPage(url)
	page.Click(moreCommentSelector)
	rawEvents := page.Text(commentSelector)
	if rawEvents != "" {
		events := strings.Split(rawEvents, "\n")
		return &events
	}

	return nil
}

func normalize(comments []string) []*domain.Commentary {
	var time string
	var commentaries = make([]*domain.Commentary, 0)

	for _, value := range comments {
		if common.IsTimeCounter(value) {
			time = value
		} else {
			commentary := domain.Commentary{
				Time:    time,
				Comment: value,
			}
			commentaries = append(commentaries, &commentary)
			time = ""
		}
	}
	reverse(&commentaries)
	return commentaries
}

func reverse(list *[]*domain.Commentary) {
	for i := 0; i < len(*list)/2; i++ {
		j := len(*list) - i - 1
		(*list)[i], (*list)[j] = (*list)[j], (*list)[i]
	}
}
