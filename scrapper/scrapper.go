package scrapper

import (
	"github.com/go-rod/rod"
)

type Scrapper struct {
	Browser *rod.Browser
	Page    *rod.Page
}

func New() *Scrapper {
	return &Scrapper{
		Browser: rod.New().MustConnect(),
	}
}

func (s Scrapper) GoPage(url string) *Scrapper {
	s.Page = s.Browser.MustPage(url)
	return &s
}

func (s Scrapper) ElementsByPattern(selector string, pattern string) rod.Elements {
	return s.Page.MustElement(selector).MustElements(pattern)
}

func (s Scrapper) Click(selector string) {
	s.Page.MustElement(selector).MustElement("*").MustClick()
}

func (s Scrapper) Text(selector string) string {
	return s.Page.MustElement(selector).MustText()
}

func (s Scrapper) Elements(selector string) rod.Elements {
	return s.Page.MustElements(selector)
}
