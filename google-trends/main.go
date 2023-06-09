package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	Title    string `xml:"title"`
	ItemList []Item `xml:"item"`
}

type Item struct {
	Title     string `xml:"title"`
	Link      string `xml:"link"`
	Traffic   string `xml:"ht:approx_traffic"`
	NewsItems []News `xml:"ht:news_item"`
}

type News struct {
	Headline     string `xml:"news_item_title"`
	HeadlineLink string `xml:"news_item_url"`
}

func main() {
	var r RSS
	data := readGoogleTrends()

	err := xml.Unmarshal(data, &r)

	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Println("\n Below are the top 20 Google Trends for the day !")
	fmt.Println("-------------------------------------------------")
	for i, item := range r.Channel.ItemList {
		rank := i + 1
		fmt.Println("Rank: ", rank)
		fmt.Println("Search Term: ", item.Title)
		fmt.Println("Link to the trend: ", item.Link)
		fmt.Println("News: ")
		for _, news := range item.NewsItems {
			fmt.Println("Headline: ", news.Headline)
			fmt.Println("Headline Link: ", news.HeadlineLink)
		}
		fmt.Println("-------------------------------------------------")
	}
}

func getGoogleTrends() *http.Response {
	resp, err := http.Get("https://trends.google.com/trends/trendingsearches/daily/rss?geo=PT")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return resp
}

func readGoogleTrends() []byte {
	resp := getGoogleTrends()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return data
}
