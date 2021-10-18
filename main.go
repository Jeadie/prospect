package main

import (
	"encoding/csv"
	"fmt"
	"github.com/Jeadie/prospect/prospect"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {

	providers := []prospect.Provider{
		prospect.ReutersProvider{},
		prospect.MiningComProvider{},
		prospect.AFRMiningProvider{
			BaseUrl: "https://www.afr.com/companies/mining",
		},
	}

	// Channels & Wait groups
	providerWg := new(sync.WaitGroup)
	providerWg.Add(len(providers)) // Each provider + CSV waitGroup
	results := make(chan *prospect.ResourceReference)

	// Setup each provider
	for _, p := range providers {
		body, _ := getBodyFromUrl(p.GetBaseUrl().String())

		// response body is streamed on demand. Must close connection at end.
		defer closeHttpBody(body)

		doc, _ := goquery.NewDocumentFromReader(body)
		selectionChan := make(chan *goquery.Selection)

		go p.GetResources(doc, selectionChan)

		p := p
		go func(wg *sync.WaitGroup) {
			for s := range selectionChan {
				results <- p.ToResource(s)
			}
			wg.Done()
		}(providerWg)
	}

	// When all providers are done, results are done.
	go func(c chan *prospect.ResourceReference, wg *sync.WaitGroup) {
		wg.Wait()
		close(c)
	}(results, providerWg)

	// Create daily file
	y, m, d := time.Now().Date()
	f, _ := os.Create(fmt.Sprintf("%d-%s-%d.csv", y, m.String(), d))

	// close csv file on main() end
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("an error occurred when closing file. Result may be invalid")
		}
	}(f)

	csv_wg := new(sync.WaitGroup)
	csv_wg.Add(1)

	go func(wg *sync.WaitGroup) {
		// Process all results
		writeToCsv(f, results)
		wg.Done()
	}(csv_wg)
	csv_wg.Wait()
}

func closeHttpBody(b io.ReadCloser) {
	err := b.Close()
	if err != nil {
		fmt.Println("an error occurred when closing http GET")
	}
}

func getBodyFromUrl(url string) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func writeToCsv(f *os.File, c chan *prospect.ResourceReference) {
	writer := csv.NewWriter(f)
	for r := range c {
		writer.Write([]string{r.Href.String(), r.Title, r.Preview})
	}
	writer.Flush()
}
