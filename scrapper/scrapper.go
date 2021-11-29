package scrapper

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
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
	// TODO: only wait once
	s.Page = s.Browser.MustPage(url).MustWaitLoad()
	return &s
}

func (s Scrapper) ElementsByPattern(selector string, pattern string) rod.Elements {
	return s.Page.MustElement(selector).MustWaitLoad().MustElements(pattern)
}

func (s Scrapper) DynamicElementsByPattern(selector string, pattern string) rod.Elements {
	return s.Page.MustElement(selector).
			MustWaitLoad().
			MustPress(input.ArrowDown).
			MustWaitLoad().
			MustElements(pattern)
}

func (s Scrapper) Click(selector string) {
	if s.exists(selector) {
		s.Page.MustElement(selector).MustElement("*").MustClick()
	}
}

func (s Scrapper) Text(selector string) string {
	return s.Page.MustElement(selector).MustText()
}

func (s Scrapper) Elements(selector string) rod.Elements {
	return s.Page.MustElements(selector)
}

// TODO: improve the way to check existence
func (s Scrapper) exists(selector string) bool {
	elements, err := s.Page.Elements(selector)
	return len(elements) > 0 || err != nil
}
