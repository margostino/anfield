package common

type Results struct {
	Url      string `yaml:"url"`
	Selector string `yaml:"selector"`
	Pattern  string `yaml:"pattern"`
}

type Matches struct {
	MoreCommentsSelector string `yaml:"moreCommentsSelector"`
	CommentsSelector     string `yaml:"commentsSelector"`
	CommentUrlParam      string `yaml:"commentUrlParam"`
	InfoUrlParam         string `yaml:"infoUrlParam"`
	InfoSelector         string `yaml:"infoSelector"`
}

type Configuration struct {
	Results *Results `yaml:"results"`
	Matches *Matches `yaml:"matches"`
}
