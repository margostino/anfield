package bot

import (
	"github.com/margostino/anfield/common"
)

func shouldFollow(message string) bool {
	return message == "/follow"
}

func followReply() (interface{}, string) {
	reply := "Which player would do like to follow?"
	return nil, reply
}

func shouldFollowPlayer(previousMessage string) bool {
	return shouldFollow(previousMessage)
}

func playerFollowerReply(message string, userId int64) (interface{}, string) {
	Follow(userId, message)
	return nil, "Done!"
}

func IsFollowing(message string, chatId int64) bool {
	return common.InSlice(message, Following()[chatId])
}
