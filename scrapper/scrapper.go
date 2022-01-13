package scrapper

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type Scrapper struct {
	browser *rod.Browser
	page    *rod.Page
}

func New() *Scrapper {
	return &Scrapper{
		browser: rod.New().MustConnect(),
	}
}

func (s Scrapper) GoPage(url string) *Scrapper {
	// TODO: only wait once
	s.page = s.browser.MustPage(url).MustWaitLoad()
	return &s
}

func (s Scrapper) ElementsByPattern(selector string, pattern string) rod.Elements {
	return s.page.MustElement(selector).MustWaitLoad().MustElements(pattern)
}

func (s Scrapper) DynamicElementsByPattern(selector string, pattern string) rod.Elements {
	return s.page.MustElement(selector).
		MustWaitLoad().
		MustPress(input.ArrowDown).
		MustWaitLoad().
		MustElements(pattern)
}

func (s Scrapper) Click(selector string) {
	if s.exists(selector) {
		s.page.MustElement(selector).MustElement("*").MustClick()
	}
}

func (s Scrapper) Text(selector string) string {
	if s.exists(selector) {
		return s.page.MustElement(selector).MustText()
	}
	return ""
}

func (s Scrapper) Elements(selector string) rod.Elements {
	return s.page.MustElements(selector)
}

// TODO: improve the way to check existence
func (s Scrapper) exists(selector string) bool {
	elements, err := s.page.Elements(selector)
	return len(elements) > 0 || err != nil
}

func (s Scrapper) Close() {
	s.browser.MustClose()
}
