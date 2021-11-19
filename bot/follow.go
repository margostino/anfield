package bot

import (
	"strings"
)

func shouldFollow(message string) bool {
	return message == "/follow"
}

func shouldUnfollow(message string) bool {
	return message == "/unfollow"
}

func followQuestion() (interface{}, string) {
	reply := "Which player would do like to follow?"
	return nil, reply
}

func unfollowQuestion() (interface{}, string) {
	reply := "Which player would do like to unfollow?"
	return nil, reply
}

func shouldFollowPlayer(previousMessage string) bool {
	return shouldFollow(previousMessage)
}

func shouldUnfollowPlayer(previousMessage string) bool {
	return shouldUnfollow(previousMessage)
}

func followReply(message string, userId int64) (interface{}, string) {
	Follow(userId, message)
	return nil, "Done!"
}

func unfollowReply(message string, userId int64) (interface{}, string) {
	Unfollow(userId, message)
	return nil, "Done!"
}

func IsFollowing(message string, chatId int64) bool {
	lowerMessage := strings.ToLower(message)
	for _, value := range following[chatId] {
		lowerValue := strings.ToLower(value)
		if strings.Contains(lowerMessage, lowerValue) {
			return true
		}
	}
	return false
}
