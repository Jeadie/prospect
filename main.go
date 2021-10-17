package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Jeadie/prospect/prospect"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	reutersUrl := "https://www.reuters.com/markets/commodities"
	miningUrl := "https://www.mining.com/#latest-section"

	reutersBody, _ := getBodyFromUrl(reutersUrl)
	miningBody, _ := getBodyFromUrl(miningUrl)

	// response body is streamed on demand. Must close connection at end.
	defer closeHttpBody(reutersBody)
	defer closeHttpBody(miningBody)

	// Channels & Wait groups
	wg := new(sync.WaitGroup)
	wg.Add(1)
	results := make(chan *prospect.ResourceReference)

	reuterSelections := make(chan *goquery.Selection)
	miningSelections := make(chan *goquery.Selection)

	reuters := prospect.ReutersProvider{}
	mining := prospect.MiningComProvider{}

	reuterDoc, _ := goquery.NewDocumentFromReader(reutersBody)
	go reuters.GetResources(reuterDoc, reuterSelections)

	miningDoc, _ := goquery.NewDocumentFromReader(miningBody)
	go mining.GetResources(miningDoc, miningSelections)

	fmt.Println("Before resources")
	fmt.Println("After resources")

	// Transform all into ResourceReferences
	go func() {
		for s := range reuterSelections {
			results <- reuters.ToResource(s)
		}
		for s := range miningSelections {
			results <- mining.ToResource(s)
		}
		close(results)
	}()

	// Create daily file
	y, m, d := time.Now().Date()
	f, _ := os.Create(fmt.Sprintf("%d-%s-%d-reuters.csv", y, m.String(), d))

	// close csv file on main() end
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("an error occurred when closing file. Result may be invalid")
		}
	}(f)

	go func(wg *sync.WaitGroup) {
		// Process all results
		writeToCsv(f, results)
		wg.Done()
	}(wg)
	wg.Wait()
}

func closeHttpBody(Body io.ReadCloser) {
	err := Body.Close()
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
	fmt.Println("Writing CSVs")
	writer := csv.NewWriter(f)
	for r := range c {
		writer.Write([]string{r.Href.String(), r.Title, r.Preview})
	}
	fmt.Println("Done writing CSVs")
	writer.Flush()
}
