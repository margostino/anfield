package domain

import "time"

type Player struct {
	Name string
}

type Team struct {
	Name              string
	Form              []Player
	SubstitutePlayers []Player
}

type Metadata struct {
	Url      string
	Date     time.Time
	Finished bool
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
	Lineups  *Lineups
	Data     []*Commentary
}

type Message struct {
	Metadata *Metadata
	Lineups  *Lineups
	Data     *Commentary
}

type Data struct {
	Comments []Commentary
}

type Document struct {
	Metadata *Metadata
	Data     *Data
}
