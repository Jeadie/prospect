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
	GetBaseUrl() *url.URL
	//Parse(document goquery.Document) (Resource, error)
}

type ReutersProvider struct {
	// Currently, nothing is needed.
}

func (p ReutersProvider) GetResources(d *goquery.Document, output chan *goquery.Selection) {
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

func (p ReutersProvider) GetBaseUrl() *url.URL {
	v, err := url.Parse("https://www.reuters.com/markets/commodities")
	if err != nil {
		return &url.URL{}
	}
	return v
}

type MiningComProvider struct {
	// Currently, nothing is needed.
}

func (m MiningComProvider) GetResources(d *goquery.Document, output chan *goquery.Selection) {
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

func (m MiningComProvider) GetBaseUrl() *url.URL {
	v, err := url.Parse("https://www.mining.com/#latest-section")
	if err != nil {
		return &url.URL{}
	}
	return v
}

type AFRMiningProvider struct {
	BaseUrl string
}

func (m AFRMiningProvider) GetResources(d *goquery.Document, output chan *goquery.Selection) {
	d.Find("[data-pb-type=\"st\"]").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return s.Find("[data-pb-type=\"ab\"]").Length() > 0
	}).Each(func(_ int, s *goquery.Selection) { output <- s })
	close(output)
}

func (m AFRMiningProvider) ToResource(s *goquery.Selection) *ResourceReference {
	title := s.Find("[data-pb-type=\"hl\"]").Find("a")
	href, _ := title.Attr("href")
	link, err := url.Parse(m.BaseUrl + href)

	body := s.Find("[data-pb-type=\"ab\"]")

	if err != nil {
		return &ResourceReference{}
	} else {
		return &ResourceReference{
			Href:    link,
			Title:   cleanString(title.Text()),
			Preview: cleanString(body.Text()),
		}
	}
}

func (m AFRMiningProvider) GetBaseUrl() *url.URL {
	v, err := url.Parse(m.BaseUrl)
	if err != nil {
		return &url.URL{}
	}
	return v
}

func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "…")
	return s
}
