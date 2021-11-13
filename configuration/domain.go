package configuration

import "time"

type ResultsConfig struct {
	Url      string `yaml:"url"`
	Selector string `yaml:"selector"`
	Pattern  string `yaml:"pattern"`
}

type FixturesConfig struct {
	Url      string `yaml:"url"`
	Selector string `yaml:"selector"`
	Pattern  string `yaml:"pattern"`
}

type DataConfig struct {
	Update      bool   `yaml:"update"`
	MatchesPath string `yaml:"matchesPath"`
}

type SourceConfig struct {
	Url string `yaml:"url"`
}

type KafkaConfig struct {
	Topic           string `yaml:"topic"`
	Protocol        string `yaml:"protocol"`
	Address         string `yaml:"address"`
	ConsumerGroupId string `yaml:"consumerGroupId"`
}

type BotConfig struct {
	Token   string  `yaml:"token"`
	ChatIds []int64 `yaml:"chatIds"`
}

type RealtimeConfig struct {
	Matches      []string      `yaml:"matches"`
	StopFlag     string        `yaml:"stopFlag"`
	GraceEndTime time.Duration `yaml:"graceEndTime"`
	CountDown    int           `yaml:"countDown"`
}

type CommentaryConfig struct {
	MoreCommentsSelector string `yaml:"moreCommentsSelector"`
	Selector             string `yaml:"selector"`
	Params               string `yaml:"params"`
}

type InfoConfig struct {
	Selector string `yaml:"selector"`
	Params   string `yaml:"params"`
}

type LineupsConfig struct {
	HomeSelector       string `yaml:"homeSelector"`
	HomeTeamSelector   string `yaml:"homeTeamSelector"`
	AwaySelector       string `yaml:"awaySelector"`
	AwayTeamSelector   string `yaml:"awayTeamSelector"`
	SubstituteSelector string `yaml:"substituteSelector"`
	Params             string `yaml:"params"`
}

type Configuration struct {
	AppPath    string            `yaml:"appPath"`
	Data       *DataConfig       `yaml:"data"`
	Source     *SourceConfig     `yaml:"source"`
	Results    *ResultsConfig    `yaml:"results"`
	Fixtures   *FixturesConfig   `yaml:"fixtures"`
	Commentary *CommentaryConfig `yaml:"commentary"`
	Info       *InfoConfig       `yaml:"info"`
	Lineups    *LineupsConfig    `yaml:"lineups"`
	Bot        *BotConfig        `yaml:"bot"`
	Realtime   *RealtimeConfig   `yaml:"realtime"`
	Kafka      *KafkaConfig      `yaml:"kafka"`
	Rules      *Rules
}

type Rule struct {
	Pattern string  `yaml:"pattern"`
	Score   float64 `yaml:"score"`
}

type Rules struct {
	Team   []Rule `yaml:"team"`
	Player []Rule `yaml:"player"`
}
