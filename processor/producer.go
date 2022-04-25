package processor

import (
	"fmt"
	"github.com/go-rod/rod/lib/proto"
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
	lineupsSelector := a.configuration.Scrapper.LineupsSelector
	page := a.scrapper.GoPage(lineupsUrl)
	homeTeamName := page.Text(homeTeamSelector)
	awayTeamName := page.Text(awayTeamSelector)
	rawElements := page.Text(lineupsSelector)
	lineupsStartFlag := a.configuration.Scrapper.LineupStartFlag
	substitutesStartFlag := a.configuration.Scrapper.SubstitutesStartFlag
	rawLineupElements := strings.Split(rawElements, "\n")
	homeFormation, awayFormation := extractFormationDataElement(rawLineupElements, lineupsStartFlag)
	homeSubstitutes, awaySubstitutes := extractSubstituteDataElement(rawLineupElements, substitutesStartFlag)

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

func newPlayer(name string) *domain.Player {
	return &domain.Player{
		Name: name,
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
			startFlag := a.configuration.Scrapper.CommentaryStartFlag
			endFlag := a.configuration.Scrapper.CommentaryEndFlag
			rawCommentaries := extractCommentaryDataElement(*rawEvents, startFlag, endFlag)
			commentaries := normalizeCommentary(rawCommentaries)
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
	moreCommentTextSelector := a.configuration.Scrapper.MoreCommentsTextSelector
	commentSelector := a.configuration.Scrapper.CommentarySelector
	page := a.scrapper.GoPage(url)
	elements := page.Elements(moreCommentSelector)

	for _, value := range elements {
		text, _ := value.Text()
		if text == moreCommentTextSelector {
			btn := value.MustElement("button")
			btn.Click(proto.InputMouseButtonLeft)
			break
		}
	}
	rawEvents := page.Text(commentSelector)
	if rawEvents != "" {
		events := strings.Split(rawEvents, "\n")
		return &events
	}

	return nil
}

func extractCommentaryDataElement(elements []string, startFlag string, endFlag string) []string {
	var index string
	var shouldStartNormalizing, shouldContinueNormalizing bool
	var results = make([]string, 0)

	for _, value := range elements {
		shouldStartNormalizing = value == startFlag
		if common.IsTimeCounter(value) {
			index = value
		} else if shouldStartNormalizing || shouldContinueNormalizing {
			shouldContinueNormalizing = !(value == endFlag)
			if !shouldStartNormalizing && shouldContinueNormalizing && index != "" {
				results = append(results, index, value)
				index = ""
			}
		}
	}
	reverse(&results)
	return results
}

func extractFormationDataElement(elements []string, startFlag string) ([]domain.Player, []domain.Player) {
	homeFormation := make([]domain.Player, 0)
	awayFormation := make([]domain.Player, 0)
	var moreFlagsBeforeExtracting bool
	startIndexLookup := 0
	startFlagValues := strings.Split(startFlag, ",")

	for _, value := range elements {
		if startIndexLookup < len(startFlagValues) {
			moreFlagsBeforeExtracting = value == startFlagValues[startIndexLookup] || startFlagValues[startIndexLookup] == "%"
			if moreFlagsBeforeExtracting && len(startFlagValues) > 1 {
				startIndexLookup += 1
			}
		} else {
			moreFlagsBeforeExtracting = false
		}
		if !common.IsFormationNumber(value) && !moreFlagsBeforeExtracting && startIndexLookup == len(startFlagValues) {
			if len(homeFormation) == 11 && len(awayFormation) == 11 {
				break
			}
			if len(homeFormation) < 11 {
				homeFormation = append(homeFormation, *newPlayer(value))
			} else {
				awayFormation = append(awayFormation, *newPlayer(value))
			}
		}
	}
	return homeFormation, awayFormation
}

func extractSubstituteDataElement(elements []string, startFlag string) ([]domain.Player, []domain.Player) {
	var index = -1
	homeSubstituteFormation := make([]domain.Player, 0)
	awaySubstituteFormation := make([]domain.Player, 0)

	for _, value := range elements {
		if index == 19 {
			break
		}

		if value == startFlag || index == 0 {
			index += 1
		}

		if index > 0 {
			index += 1
			if index%2 != 0 {
				awaySubstituteFormation = append(awaySubstituteFormation, *newPlayer(value))
			} else {
				homeSubstituteFormation = append(homeSubstituteFormation, *newPlayer(value))
			}
		}
	}
	return homeSubstituteFormation, awaySubstituteFormation
}

func normalizeCommentary(rawCommentary []string) []*domain.Commentary {
	var comment string
	var commentaries = make([]*domain.Commentary, 0)

	for _, value := range rawCommentary {
		if comment != "" {
			commentary := domain.Commentary{
				Time:    value,
				Comment: comment,
			}
			commentaries = append(commentaries, &commentary)
			comment = ""
		} else {
			comment = value
		}
	}
	return commentaries
}

//func reverse(list *[]*domain.Commentary) {
func reverse(list *[]string) {
	for i := 0; i < len(*list)/2; i++ {
		j := len(*list) - i - 1
		(*list)[i], (*list)[j] = (*list)[j], (*list)[i]
	}
}
