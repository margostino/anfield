package context

import "time"

type Results struct {
	Url      string `yaml:"url"`
	Selector string `yaml:"selector"`
	Pattern  string `yaml:"pattern"`
}

type Fixtures struct {
	Url      string `yaml:"url"`
	Selector string `yaml:"selector"`
	Pattern  string `yaml:"pattern"`
}

type Data struct {
	Update  bool   `yaml:"update"`
	Matches string `yaml:"matches"`
}

type Source struct {
	Url string `yaml:"url"`
}

type Bot struct {
	Token string `yaml:"token"`
}

type Realtime struct {
	Matches      []string      `yaml:"matches"`
	StopFlag     string        `yaml:"stopFlag"`
	GraceEndTime time.Duration `yaml:"graceEndTime"`
	CountDown    int           `yaml:"countDown"`
}

type Commentary struct {
	MoreCommentsSelector string `yaml:"moreCommentsSelector"`
	Selector             string `yaml:"selector"`
	Params               string `yaml:"params"`
}

type Info struct {
	Selector string `yaml:"selector"`
	Params   string `yaml:"params"`
}

type Configuration struct {
	AppPath    string      `yaml:"appPath"`
	Data       *Data       `yaml:"data"`
	Source     *Source     `yaml:"source"`
	Results    *Results    `yaml:"results"`
	Fixtures   *Fixtures   `yaml:"fixtures"`
	Commentary *Commentary `yaml:"commentary"`
	Info       *Info       `yaml:"info"`
	Bot        *Bot        `yaml:"bot"`
	Realtime   *Realtime   `yaml:"realtime"`
}
