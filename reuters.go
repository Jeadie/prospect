package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"github.com/jeadie/prospect/prospect"
	"sync"
	"time"
)

func main() {
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

	// Channels & Wait groups
	wg := new(sync.WaitGroup)
	wg.Add(1)
	results := make(chan *prospect.ResourceReference)
	selections := make(chan *goquery.Selection)

	reuters := prospect.ReutersProvider{}
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	fmt.Println("Before resources")
	go reuters.GetResources(doc, selections)
	fmt.Println("After resources")

	go func() {
		for s := range selections {
			results <- reuters.ToResource(s)
		}
		close(results)
	}()

	fmt.Println("After results")

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

func writeToCsv(f *os.File, c chan *prospect.ResourceReference) {
	fmt.Println("Writing CSVs")
	writer := csv.NewWriter(f)
	for r := range c {
		writer.Write([]string{r.Href.String(), r.Title, r.Preview})
	}
	fmt.Println("Done writing CSVs")
	writer.Flush()
}
