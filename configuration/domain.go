package configuration

import "time"

const (
	DYNAMIC_RULE = "dynamic"
	STATIC_RULE  = "static"
)

type Configuration struct {
	App        *AppConfig        `yaml:"app"`
	Source     *SourceConfig     `yaml:"source"`
	Scrapper   *ScrapperConfig   `yaml:"scrapper"`
	Bot        *BotConfig        `yaml:"bot"`
	Logger     *LoggerConfig     `yaml:"logger"`
	DataLoader *DataLoaderConfig `yaml:"dataLoader"`
	Events     *EventsConfig     `yaml:"events"`
	Kafka      *KafkaConfig      `yaml:"kafka"`
	Mongo      *MongoConfig      `yaml:"mongo"`
	Rules      []Rule
}

type LoggerConfig struct {
	CompletionStep int `yaml:"completionStep"`
}

type ScrapperConfig struct {
	Url                  string `yaml:"url"`
	ResultsPath          string `yaml:"resultsPath"`
	FixturesPath         string `yaml:"fixturesPath"`
	MatchRowsSelector    string `yaml:"matchRowsSelector"`
	HrefPattern          string `yaml:"hrefPattern"`
	MoreCommentsSelector string `yaml:"moreCommentsSelector"`
	CommentarySelector   string `yaml:"commentarySelector"`
	CommentaryParams     string `yaml:"commentaryParams"`
	InfoSelector         string `yaml:"infoSelector"`
	InfoParams           string `yaml:"infoParams"`
	HomeSelector         string `yaml:"homeSelector"`
	HomeTeamSelector     string `yaml:"homeTeamSelector"`
	AwaySelector         string `yaml:"awaySelector"`
	AwayTeamSelector     string `yaml:"awayTeamSelector"`
	SubstituteSelector   string `yaml:"substituteSelector"`
	LineupsParams        string `yaml:"lineupsParams"`
	UrlProperty          string `yaml:"urlProperty"`
}

type SourceConfig struct {
	Update      bool   `yaml:"update"`
	MatchesPath string `yaml:"matchesPath"`
}

type KafkaConfig struct {
	Topic           string `yaml:"topic"`
	Protocol        string `yaml:"protocol"`
	Address         string `yaml:"address"`
	ConsumerGroupId string `yaml:"consumerGroupId"`
}

type MongoConfig struct {
	Hostname          string `yaml:"hostname"`
	Port              int    `yaml:"port"`
	Database          string `yaml:"database"`
	MatchesCollection string `yaml:"matchesCollection"`
	AssetsCollection  string `yaml:"assetsCollection"`
	UsersCollection   string `yaml:"usersCollection"`
}
type BotConfig struct {
	Token                string  `yaml:"token"`
	ChatIds              []int64 `yaml:"chatIds"`
	KafkaConsumerGroupId string  `yaml:"kafkaConsumerGroupId"`
}

type DataLoaderConfig struct {
	KafkaConsumerGroupId string `yaml:"kafkaConsumerGroupId"`
}

type EventsConfig struct {
	Matches      []string      `yaml:"matches"`
	StopFlag     string        `yaml:"stopFlag"`
	GraceEndTime time.Duration `yaml:"graceEndTime"`
	CountDown    int           `yaml:"countDown"`
}

type AppConfig struct {
	Path           string        `yaml:"path"`
	ChannelTimeout time.Duration `yaml:"channelTimeout"`
}

type Rule struct {
	Pattern string  `yaml:"pattern"`
	Score   float64 `yaml:"score"`
	Pos     int     `yaml:"pos"`
	Type    string  `yaml:"type"`
}

type Rules struct {
	ScoringRules []Rule `yaml:"scoringRules"`
}
