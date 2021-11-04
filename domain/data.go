package domain

type Team struct {
	Name              string
	Form              []string
	SubstitutePlayers []string
}

type Metadata struct {
	Url      string
	Date     string
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
