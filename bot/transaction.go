package bot

func isSelling(message string) bool {
	return message == "/sell"
}

// TODO: support all commands in one shot (e.g. /buy salah 10)
func isBuying(message string) bool {
	return message == "/buy"
}

func shouldBuyAsset(previousMessages []string) bool {
	return len(previousMessages) == 1 && isBuying(previousMessages[0])
}

func shouldBuyAssetUnits(previousMessages []string) bool {
	return len(previousMessages) == 2 && isBuying(previousMessages[0]) && previousMessages[1] != ""
}

func buyAssetQuestion() (interface{}, string) {
	reply := "Which asset would do like to buy?"
	return nil, reply
}

func buyUnitsQuestion() (interface{}, string) {
	reply := "How much would do like to buy?"
	return nil, reply
}

func buyInvalidUnits() (interface{}, string) {
	reply := "Units invalid. It should be a number."
	return nil, reply
}
