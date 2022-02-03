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

type User struct {
	Username  string
	FirstName string
	LastName  string
	Id        int
}

type Wallet struct {
	Budget      float64
	LastUpdated time.Time `bson:"last_updated"`
}

type Transaction struct {
	UserId    string
	AssetId   string
	Units     int
	Value     float64
	Timestamp time.Time
}

// MongoDB Collections
// TODO: isolate data model domain
// TODO: separate Metadata/Data from App/DB domain

type MatchDocument struct {
	Metadata *Metadata
	Data     *Data
}

type AssetDocument struct {
	Id          string `bson:"_id"`
	Name        string
	Score       float64
	LastUpdated time.Time
}

type UserDocument struct {
	Id        string `bson:"_id"`
	SocialId  int64  `bson:"social_id"`
	Username  string
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Wallet    *Wallet
}

type TransactionDocument struct {
	UserId    string
	AssetId   string
	Units     int
	Value     float64
	Timestamp time.Time
}
