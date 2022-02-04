package bot

import (
	"errors"
	"strconv"
	"strings"
)

func isSelling(message string) bool {
	return message == "/sell"
}

// TODO: support all commands in one shot (e.g. /buy salah 10)
func isBuying(message string) bool {
	return message == "/buy"
}

func shouldBuyAsset(messages []string) bool {
	return len(messages) == 1 && isBuying(messages[0])
}

func buyAssetValueInstruction() (interface{}, string) {
	reply := "Please send Asset Name and Value separated by space.\nExample:\nsalah 2"
	return nil, reply
}

func extractTransactionFrom(message string) (string, int, error) {
	values := strings.Split(message, " ")

	if len(values) == 2 {
		assetName := values[0]
		units, err := strconv.Atoi(values[1])
		if err != nil {
			return "", -1, errors.New("input is invalid. Unit input should be a number")
		}
		return assetName, units, nil
	}

	return "", -1, errors.New("input is invalid. Input should be 2 values separated by space (Example: Salah 10)")
}
