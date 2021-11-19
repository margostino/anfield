package scorer

import (
	"errors"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Scoring struct {
	Team  string
	Score float64
}

type Stats struct {
	Players map[string]*Scoring
}

const BALL_POSSESSION_RULE = "ball possession"

var stats *Stats
var entityRegex *regexp.Regexp

func Initialize() {
	stats = &Stats{
		Players: make(map[string]*Scoring),
	}
}

func CalculateScoring(homeTeam *domain.Team, awayTeam *domain.Team, commentary *domain.Commentary) {
	comment := strings.ToLower(commentary.Comment)
	rules, err := getRules(comment)

	if err == nil {
		mergePlayers(homeTeam, awayTeam)
		entities := entityRegex.FindAllString(comment, -1)
		applyScoring(comment, entities, rules)
	}

	// TODO: validate VAR
	// TODO: validate 2 players in one comment
	// TODO: validate team scoring
	// TODO: ball possession
}

func applyScoring(comment string, entities []string, rules []configuration.Rule) {
	ratios := getRatios(comment)

	if entities != nil {
		for _, rule := range rules {
			applyRule(entities, ratios, &rule)
		}
	}
}

func getRatios(comment string) map[string]float64 {
	var ratios = make(map[string]float64)
	if isBallPossession(comment) {
		updateTeamsPossession(comment, &ratios)
	}
	return ratios
}

func updateTeamsPossession(comment string, ratios *map[string]float64) {
	teamsPossessionRaw := strings.ReplaceAll(comment, BALL_POSSESSION_RULE+":", "")
	teamsPossessionParts := strings.Split(teamsPossessionRaw, ",")

	for _, possessionRaw := range teamsPossessionParts {
		name := strings.TrimSpace(strings.Split(possessionRaw, ":")[0])
		posessionPercentage := strings.TrimSpace(strings.Split(possessionRaw, ":")[1])
		posessionString := strings.ReplaceAll(posessionPercentage, "%", "")
		posessionNumber, _ := strconv.ParseFloat(posessionString, 64)
		(*ratios)[name] = posessionNumber
	}
}

func applyRule(entities []string, ratios map[string]float64, rule *configuration.Rule) {
	for _, entity := range entities {
		for key, value := range stats.Players {
			match := strings.Contains(key, entity) || strings.Contains(value.Team, entity)
			ratio, hasRatio := ratios[entity]
			if match && hasRatio {
				value.Score += rule.Score * ratio / 100
			} else {
				value.Score += rule.Score
			}
		}
	}
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
		newEntityRegex(homeTeam.Name, awayTeam.Name)
	}
}

func newEntityRegex(homeTeam string, awayTeam string) {
	values := joinPlayers()
	appendTeamNames(homeTeam, &values)
	appendTeamNames(awayTeam, &values)
	pattern := "(" + strings.ToLower(values) + "){1}"
	entityRegex = regexp.MustCompile(pattern)
}

func appendTeamNames(teamName string, values *string) {
	*values += "|" + teamName
	teamParts := strings.Split(teamName, " ")
	if len(teamParts) == 2 {
		*values += "|" + teamParts[0]
	} else if len(teamParts) == 3 {
		*values += "|" + teamParts[0] + " " + teamParts[1]
		*values += "|" + teamParts[0]
	}
}

func appendPlayers(team *domain.Team) {
	teamName := strings.ToLower(team.Name)
	allPlayers := append(team.Form, team.SubstitutePlayers...)
	for _, player := range allPlayers {
		playerName := strings.ToLower(player.Name)
		scoring := &Scoring{
			Team:  teamName,
			Score: player.Score,
		}
		stats.Players[playerName] = scoring
	}
}

func ruleLookup(comment string, rules []configuration.Rule) []configuration.Rule {
	var result = make([]configuration.Rule, 0)
	for _, rule := range rules {
		if matchRule(&rule, comment) {
			result = append(result, rule)
		}
	}
	return result
}

func getRules(comment string) ([]configuration.Rule, error) {
	var rules = make([]configuration.Rule, 0)
	rules = append(rules, ruleLookup(comment, configuration.ScoringRules())...)

	if len(rules) == 0 {
		log.Printf("MISSING RULE %s", comment)
		return rules, errors.New("missing rule")
	}

	return rules, nil
}

func isBallPossession(comment string) bool {
	return strings.Contains(comment, BALL_POSSESSION_RULE)
}

func matchRule(rule *configuration.Rule, comment string) bool {
	if rule.Type == configuration.STATIC_RULE {
		return strings.Contains(comment, rule.Pattern)
	} else {
		match, _ := regexp.MatchString(rule.Pattern, comment)
		return match
	}
}

func Scorings() *Stats {
	return stats
}
