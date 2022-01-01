package domain

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
	H2H     string
	Lineups *Lineups
	Date    string
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
