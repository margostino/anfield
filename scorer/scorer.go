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

type Scorer struct {
	Rules []configuration.Rule
}

func NewScorer(configuration *configuration.Configuration) *Scorer {
	return &Scorer{
		Rules: configuration.Rules,
	}
}

type Score struct {
	Time  string
	Score float64
}

type Scoring struct {
	Players map[string][]*Score
}

const BALL_POSSESSION_RULE = "ball possession"

var entityRegex *regexp.Regexp

func (s Scorer) CalculateScoring(lineups *domain.Lineups, commentary *domain.Commentary) map[string]float64 {
	var scores = make(map[string]float64)
	homeTeamName := lineups.HomeTeam.Name
	awayTeamName := lineups.AwayTeam.Name
	homeTeam := append(lineups.HomeTeam.Form, lineups.HomeTeam.SubstitutePlayers...)
	awayTeam := append(lineups.AwayTeam.Form, lineups.AwayTeam.SubstitutePlayers...)
	allPlayers := append(homeTeam, awayTeam...)
	comment := strings.ToLower(commentary.Comment)
	rules, err := s.getRules(comment)

	if err == nil {

		//for _, player := range allPlayers {
		//	playerName := strings.ToLower(player.Name)
		//	scoring := &Scoring{
		//		Player: player.Name,
		//		Team:   team.Name,
		//		Score:  player.Score,
		//	}
		//	stats.Players[playerName] = scoring
		//}

		playerKeysList := make([]string, 0, len(allPlayers))
		for _, player := range allPlayers {
			playerKeysList = append(playerKeysList, player.Name)
		}

		playerKeys := strings.Join(playerKeysList, "|")
		homeTeamNameKeys := appendTeamName(homeTeamName)
		awayTeamNameKeys := appendTeamName(awayTeamName)
		keys := playerKeys + "|" + homeTeamNameKeys + "|" + awayTeamNameKeys
		pattern := "(" + strings.ToLower(keys) + "){1}"
		entityRegex = regexp.MustCompile(pattern)
		entities := entityRegex.FindAllString(comment, -1)
		ratios := getSharedRatios(comment)

		if entities != nil {
			for _, rule := range rules {
				for _, entity := range entities {
					//match := matchTeamOrPlayer(entity, scoring)

					if isPlayerEntity(entity, allPlayers) {
						ratio, hasRatio := ratios[entity]
						if hasRatio {
							scores[entity] = rule.Score * ratio / 100
						} else {
							scores[entity] = rule.Score
						}
					}
				}
			}
		}
	}

	// TODO: validate VAR
	// TODO: validate 2 players in one comment
	// TODO: validate team scoring
	// TODO: ball possession
	return scores
}

func isPlayerEntity(entity string, players []domain.Player) bool {
	for _, player := range players {
		if strings.ToLower(entity) == strings.ToLower(player.Name) {
			return true
		}
	}
	return false
}

func appendTeamName(teamName string) string {
	var pattern string
	pattern = teamName
	teamParts := strings.Split(teamName, " ")
	if len(teamParts) == 2 {
		pattern += "|" + teamParts[0]
	} else if len(teamParts) == 3 {
		pattern += "|" + teamParts[0] + " " + teamParts[1] + "|" + teamParts[0]
	}
	return pattern
}

func getSharedRatios(comment string) map[string]float64 {
	var ratios = make(map[string]float64)
	if isBallPossession(comment) {
		ratios = updateTeamsPossession(comment)
	}
	return ratios
}

func updateTeamsPossession(comment string) map[string]float64 {
	var possession = make(map[string]float64)
	teamsPossessionRaw := strings.ReplaceAll(comment, BALL_POSSESSION_RULE+":", "")
	teamsPossessionParts := strings.Split(teamsPossessionRaw, ",")

	for _, possessionRaw := range teamsPossessionParts {
		name := strings.TrimSpace(strings.Split(possessionRaw, ":")[0])
		posessionPercentage := strings.TrimSpace(strings.Split(possessionRaw, ":")[1])
		posessionString := strings.ReplaceAll(posessionPercentage, "%", "")
		posessionNumber, _ := strconv.ParseFloat(posessionString, 64)
		possession[name] = posessionNumber
	}
	return possession
}

func getScoring(comment string, rules []configuration.Rule) map[string]float64 {
	var scores = make(map[string]float64)
	entities := entityRegex.FindAllString(comment, -1)
	//	ratios := getRatios(comment)

	if entities != nil {
		//for _, rule := range rules {
		//	for _, entity := range entities {
		//		for _, scoring := range stats.Players {
		//			match := matchTeamOrPlayer(entity, scoring)
		//			ratio, hasRatio := ratios[entity]
		//			if match && hasRatio {
		//				scores[scoring.Player] = rule.Score * ratio / 100
		//			} else {
		//				scores[scoring.Player] = rule.Score
		//			}
		//		}
		//	}
		//}
	}
	return scores
}

//func joinPlayers() string {
//	keys := make([]string, 0, len(stats.Players))
//	for key, _ := range stats.Players {
//		//keys = append(keys, strings.Split(key, "-")[0])
//		keys = append(keys, key)
//	}
//	return strings.Join(keys, "|")
//}

func mergePlayers(homeTeam *domain.Team, awayTeam *domain.Team) {
	// Only first time
	//if len(stats.Players) == 0 {
	//	appendPlayers(homeTeam)
	//	appendPlayers(awayTeam)
	//	newEntityRegex(homeTeam.Name, awayTeam.Name)
	//}
}

//func newEntityRegex(homeTeam string, awayTeam string) {
//	values := joinPlayers()
//	appendTeamNames(homeTeam, &values)
//	appendTeamNames(awayTeam, &values)
//	pattern := "(" + strings.ToLower(values) + "){1}"
//	entityRegex = regexp.MustCompile(pattern)
//}
//
//func appendTeamNames(teamName string, values *string) {
//	*values += "|" + teamName
//	teamParts := strings.Split(teamName, " ")
//	if len(teamParts) == 2 {
//		*values += "|" + teamParts[0]
//	} else if len(teamParts) == 3 {
//		*values += "|" + teamParts[0] + " " + teamParts[1]
//		*values += "|" + teamParts[0]
//	}
//}
//
//func appendPlayers(team *domain.Team) {
//	allPlayers := append(team.Form, team.SubstitutePlayers...)
//	for _, player := range allPlayers {
//		playerName := strings.ToLower(player.Name)
//		scoring := &Scoring{
//			Player: player.Name,
//			Team:   team.Name,
//			Score:  player.Score,
//		}
//		stats.Players[playerName] = scoring
//	}
//}

func ruleLookup(comment string, rules []configuration.Rule) []configuration.Rule {
	var result = make([]configuration.Rule, 0)
	for _, rule := range rules {
		if matchRule(&rule, comment) {
			result = append(result, rule)
		}
	}
	return result
}

func (s Scorer) getRules(comment string) ([]configuration.Rule, error) {
	var rules = make([]configuration.Rule, 0)
	rules = append(rules, ruleLookup(comment, s.Rules)...)

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

//func matchTeamOrPlayer(entity string, scoring *Scoring) bool {
//	team := strings.ToLower(scoring.Team)
//	player := strings.ToLower(scoring.Player)
//	return strings.Contains(player, entity) || strings.Contains(team, entity)
//}
