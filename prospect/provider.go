package prospect

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)


type ResourceReference struct {
	Href    *url.URL
	Title   string
	Preview string
}

// Provider is responsible for providing results from a specific Document type.
type Provider interface {
	GetResources(d *goquery.Document, output chan *goquery.Selection)
	ToResource(s *goquery.Selection) *ResourceReference
	//Parse(document goquery.Document) (Resource, error)
}


type ReutersProvider struct {
	// Currently, nothing is needed.
}

func (p ReutersProvider) GetResources(d *goquery.Document, output chan *goquery.Selection)  {
	d.Find("[class^=\"FeedScroll-feed-container-\"]").Find(".item").Each(func(_ int, s *goquery.Selection) {
		output <- s
	})
	close(output)
}

func (p ReutersProvider) ToResource(s *goquery.Selection) *ResourceReference {
	a := s.Find("a")
	text := s.Find("p")
	href, _ := a.Attr("href")
	link, err := url.Parse(href)
	if err != nil {
		return &ResourceReference{}
	} else {
		return &ResourceReference{
			Href:    link,
			Title:   a.Text(),
			Preview: text.Text(),
		}
	}
}



type MiningComProvider struct {
	// Currently, nothing is needed.
}

func (m MiningComProvider) GetResources(d *goquery.Document, output chan *goquery.Selection)  {
	d.Find("[data-post-id]").Each(func(_ int, s *goquery.Selection) {
		output <- s
	})
	close(output)
}

func (m MiningComProvider) ToResource(s *goquery.Selection) *ResourceReference {
	a := s.Find("a").Last()
	text := s.Find("p")
	href, _ := a.Attr("href")
	link, err := url.Parse(href)
	if err != nil {
		return &ResourceReference{}
	} else {
		return &ResourceReference{
			Href:    link,
			Title:   cleanString(a.Text()),
			Preview: cleanString(text.Text()),
		}
	}
}

func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "…")
	return s
}
