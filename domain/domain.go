package domain

import "time"

type Player struct {
	Name  string
	Score float64
}

type Team struct {
	Name              string
	Form              []Player
	SubstitutePlayers []Player
}

type Metadata struct {
	Url     string
	Lineups *Lineups
	Date    time.Time
}

type Lineups struct {
	HomeTeam *Team
	AwayTeam *Team
}

type Commentary struct {
	Time    string
	Comment string
}

type Event struct {
	Metadata *Metadata
	Data     []*Commentary
}

type Message struct {
	Metadata *Metadata
	Data     *Commentary
}
