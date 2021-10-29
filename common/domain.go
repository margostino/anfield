package common

type Results struct {
	Url      string `yaml:"url"`
	Selector string `yaml:"selector"`
	Pattern  string `yaml:"pattern"`
}

type Data struct {
	Matches string `yaml:"matches"`
}

type Batch struct {
	Results *Results `yaml:"results"`
	Matches *Matches `yaml:"matches"`
}

type Matches struct {
	MoreCommentsSelector string `yaml:"moreCommentsSelector"`
	CommentsSelector     string `yaml:"commentsSelector"`
	CommentUrlParam      string `yaml:"commentUrlParam"`
	InfoUrlParam         string `yaml:"infoUrlParam"`
	InfoSelector         string `yaml:"infoSelector"`
}

type Configuration struct {
	AppPath string `yaml:"appPath"`
	Data    *Data  `yaml:"data"`
	Batch   *Batch `yaml:"data"`
}
