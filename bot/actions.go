package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Action interface {
	shouldReply(input string) bool
	reply(update *tgbotapi.Update) (interface{}, string)
}

type Start struct {
	Command string
}

// NewActions TODO: tbd by config
func NewActions() []Action {
	var actions = make([]Action, 0)

	start := Start{
		Command: "/start",
	}

	actions = append(actions, start)

	return actions
}

func (s Start) shouldReply(input string) bool {
	return s.Command == input
}

func (s Start) reply(update *tgbotapi.Update) (interface{}, string) {
	return nil, "nil"
}
