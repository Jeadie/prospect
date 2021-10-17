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
	result := make([]ResourceReference, 0)
	response, err := http.Get("https://www.reuters.com/markets/commodities")
	if err != nil {
		return
	}

	// response body is streamed on demand. Must close connection at end.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("an error occurred when closing http GET")
		}
	}(response.Body)

	doc, _ := goquery.NewDocumentFromReader(response.Body)
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

	// close csv file on main() end
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("an error occurred when closing file. Result may be invalid")
		}
	}(f)
	c := csv.NewWriter(f)

	c.Write([]string{"link", "title", "preview"})
	for _, r := range result {
		c.Write([]string{r.href.String(), r.title, r.preview})
	}

	// Ensure full output has flushed before deferred close
	c.Flush()
}
