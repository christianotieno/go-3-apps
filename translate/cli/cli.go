package cli

import (
	"github.com/Jeffail/gabs"
	"io"
	"log"
	"net/http"
	"sync"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateURL = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", translateURL, nil)

	query := req.URL.Query()

	query.Add("client", "gtx")

	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	if err != nil {
		log.Fatalf("1. There was an error while creating the request: %s", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("2. There was an error while sending the request: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("There was an error while closing the response body: %s", err)
		}
	}(res.Body)

	if res.StatusCode == http.StatusTooManyRequests {
		str <- "Too many requests, try again later"
		wg.Done()
		return
	}

	parsedJson, err := gabs.ParseJSONBuffer(res.Body)
	if err != nil {
		log.Fatalf("3. There was an error while parsing the response: %s", err)
	}

	nestOne, err := parsedJson.ArrayElement(0)
	if err != nil {
		log.Fatalf("4. There was an error while parsing the response: %s", err)
	}

	nestTwo, err := nestOne.ArrayElement(0)
	if err != nil {
		log.Fatalf("5. There was an error while parsing the response: %s", err)
	}

	translatedStr, err := nestTwo.ArrayElement(0)
	if err != nil {
		log.Fatalf("6. There was an error while parsing the response: %s", err)
	}

	str <- translatedStr.Data().(string)

	wg.Done()
}
