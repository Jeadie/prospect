package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ResourceReference struct {
	href    *url.URL
	title   string
	preview string
}

func main() {
	reuters := "https://www.reuters.com/markets/commodities"
	response, err := http.Get(reuters)
	if err != nil {
		return
	}

	// Close response object once finished.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	doc, err := goquery.NewDocumentFromReader(response.Body)

	result := make([]ResourceReference, 0)
	doc.Find("[class^=\"FeedScroll-feed-container-\"]").Find(".item").Each(func(_ int, s *goquery.Selection) {
		a := s.Find("a")
		p := s.Find("p")
		href, _ := a.Attr("href")
		link, err := url.Parse(href)
		if err == nil {
			result = append(result, ResourceReference{
				href:    link,
				title:   a.Text(),
				preview: p.Text(),
			})
		}
	})
	y, m, d := time.Now().Date()
	f, _ := os.Create(fmt.Sprintf("%d-%s-%d-reuters.csv", y, m.String(), d))
	defer f.Close()
	c := csv.NewWriter(f)
	c.Write([]string{"link", "title", "preview"})

	for _, r := range result {
		c.Write([]string{r.href.String(), r.title, r.preview})
	}
	c.Flush()
}
