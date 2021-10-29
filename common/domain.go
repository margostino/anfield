package common

type Results struct {
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

type Matches struct {
	MoreCommentsSelector string `yaml:"moreCommentsSelector"`
	CommentsSelector     string `yaml:"commentsSelector"`
	CommentUrlParam      string `yaml:"commentUrlParam"`
	InfoUrlParam         string `yaml:"infoUrlParam"`
	InfoSelector         string `yaml:"infoSelector"`
}

type Configuration struct {
	AppPath string   `yaml:"appPath"`
	Data    *Data    `yaml:"data"`
	Source  *Source  `yaml:"source"`
	Results *Results `yaml:"results"`
	Matches *Matches `yaml:"matches"`
}
