package bot

import (
	"strings"
)

// TODO: enrich stats with more information (last update, highest/lowest, etc...)
// TODO: add command to explain stats
// TODO: add command to alert trends and changes and automate buy/sell operation given conditions (e.g. threshold)
// TODO: persist stats history and matches reply (realtime + batch contribution to avoid duplications)

func shouldShowStats(message string) bool {
	return message == "/stats"
}

func (a App) showStats(userId int64) (interface{}, string) {
	var reply string
	//players := scorer.Scorings().Players
	//for key, value := range players {
	//	if isFollowing(key, userId) {
	//		reply += fmt.Sprintf("Player %s, Score: %.2f\n", value.Player, value.Score)
	//	}
	//}
	//
	//if reply == "" {
	//	reply = "No stats yet!"
	//}

	return nil, reply
}

// TODO: normalize player names (and keys) everywhere to avoid FOR lookup
func isFollowing(player string, userId int64) bool {
	lowerPlayer := strings.ToLower(player)
	for _, value := range following[userId] {
		lowerValue := strings.ToLower(value)
		if strings.Contains(lowerPlayer, lowerValue) {
			return true
		}
	}
	return false
}
