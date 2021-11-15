package scorer

import (
	"errors"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"log"
	"regexp"
	"strings"
)

type Scoring struct {
	Team  string
	Score float64
}

type Stats struct {
	Players map[string]*Scoring
}

var stats *Stats
var playersRegex *regexp.Regexp

func Initialize() {
	stats = &Stats{
		Players: make(map[string]*Scoring),
	}
}

func CalculateScoring(homeTeam *domain.Team, awayTeam *domain.Team, commentary *domain.Commentary) {
	comment := "Jarrod Bowen is leaving the field to be replaced by Vladimir Coufal in a tactical substitution." //commentary.Comment
	rules, err := getRule(comment)

	if err == nil {
		mergePlayers(homeTeam, awayTeam)
		matchedPlayers := playersRegex.FindAllString(comment, -1)

		if matchedPlayers != nil {
			for _, rule := range rules {
				if len(matchedPlayers) == 1 {
					stats.Players[matchedPlayers[0]].Score += rule.Score
				} else if len(matchedPlayers) >= rule.Pos {
					stats.Players[matchedPlayers[rule.Pos-1]].Score += rule.Score
				}
			}
		}
	}

	// TODO: validate VAR
	// TODO: validate 2 players in one comment
	// TODO: validate team scoring
	// TODO: ball possession
}

func joinPlayers() string {
	keys := make([]string, 0, len(stats.Players))
	for key, _ := range stats.Players {
		//keys = append(keys, strings.Split(key, "-")[0])
		keys = append(keys, key)
	}
	return strings.Join(keys, "|")
}

func mergePlayers(homeTeam *domain.Team, awayTeam *domain.Team) {
	// Only first time
	if len(stats.Players) == 0 {
		appendPlayers(homeTeam)
		appendPlayers(awayTeam)
		newPlayersRegex()
	}
}

func newPlayersRegex() {
	values := joinPlayers()
	pattern := "(" + values + ")+"
	playersRegex = regexp.MustCompile(pattern)
}

func appendPlayers(team *domain.Team) {
	allPlayers := append(team.Form, team.SubstitutePlayers...)
	for _, player := range allPlayers {
		scoring := &Scoring{
			Team:  team.Name,
			Score: player.Score,
		}
		stats.Players[player.Name] = scoring
	}
}

func getRule(comment string) ([]configuration.Rule, error) {
	var rules = make([]configuration.Rule, 0)
	lowerComment := strings.ToLower(comment)

	for _, rule := range configuration.PlayerRules() {
		if matchRule(&rule, lowerComment) {
			rules = append(rules, rule)
		}
	}
	for _, rule := range configuration.TeamRules() {
		if matchRule(&rule, lowerComment) {
			rules = append(rules, rule)
		}
	}

	if len(rules) == 0 {
		log.Printf("MISSING RULE %s", lowerComment)
		return rules, errors.New("missing rule")
	}

	return rules, nil
}

func matchRule(rule *configuration.Rule, commentary string) bool {
	if strings.Contains(commentary, rule.Pattern) {
		return true
	}
	return false
}
