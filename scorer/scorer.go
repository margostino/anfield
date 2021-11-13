package scorer

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"log"
	"strings"
)

type Player struct {
	Team  string
	Name  string
	Score float64
}

type Stats struct {
	Players []*Player
}

var stats *Stats

func Initialize() {
	stats = &Stats{
		Players: make([]*Player, 0),
	}
}

func CalculateScoring(homeTeam *domain.Team, awayTeam *domain.Team, commentary *domain.Commentary) {
	comment := commentary.Comment
	score := scoreByRule(comment)

	if score < 0 {
		log.Printf("MISSING RULE %s", comment)
	}

	mergePlayers(homeTeam, awayTeam)

	for _, player := range stats.Players {
		if strings.Contains(comment, player.Name) {
			player.Score += score
		}
	}
	// TODO: validate VAR
	// TODO: validate 2 players in one comment
	// TODO: validate team scoring
}

func mergePlayers(homeTeam *domain.Team, awayTeam *domain.Team) {
	if len(stats.Players) == 0 {
		appendPlayers(homeTeam)
		appendPlayers(awayTeam)
	}
}

func appendPlayers(team *domain.Team) {
	allPlayers := append(team.Form, team.SubstitutePlayers...)
	for _, player := range allPlayers {
		statsPlayer := &Player{
			Name:  player.Name,
			Score: player.Score,
			Team:  team.Name,
		}
		stats.Players = append(stats.Players, statsPlayer)
	}
}

func scoreByRule(comment string) float64 {
	for _, rule := range configuration.PlayerRules() {
		if matchRule(&rule, comment) {
			return rule.Score
		}
	}
	return -1
}

func calculateTeamScoring(homeTeam *domain.Team, awayTeam *domain.Team, commentary *domain.Commentary) {

}

func matchRule(rule *configuration.Rule, commentary string) bool {
	if strings.Contains(commentary, rule.Pattern) {
		return true
	}
	return false
}
