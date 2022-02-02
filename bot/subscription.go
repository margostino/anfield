package bot

//func getSubscriptionOptions() interface{} {
//	buttons := make([]tgbotapi.KeyboardButton, 0)
//	for _, match := range Matches() {
//		button := tgbotapi.KeyboardButton{
//			Text:            match,
//			RequestContact:  false,
//			RequestLocation: false,
//		}
//		buttons = append(buttons, button)
//	}
//	return &tgbotapi.ReplyKeyboardMarkup{
//		Keyboard: [][]tgbotapi.KeyboardButton{
//			buttons,
//		},
//	}
//}
//
//func isSubscription(message string) bool {
//	return message == "/subscribe"
//}
//
//func shouldSubscribeToMatch(previousMessage string, message string) bool {
//	return isSubscription(previousMessage) && common.InSlice(message, Matches())
//}
//
//func subscriptionReply() (interface{}, string) {
//	markup := getSubscriptionOptions()
//	reply := "select a match to follow"
//	return markup, reply
//}
//
//func matchSubscriptionReply(message string, userId int64) (interface{}, string) {
//	//subscribe(userId, message)
//	return nil, "Done!"
//}


