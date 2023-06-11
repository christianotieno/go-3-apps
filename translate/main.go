package main

import (
	"flag"
	"fmt"
	"github.com/christianotieno/translate/cli"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

var sourceLang string
var targetLang string
var sourceText string

func init() {
	flag.StringVar(&sourceLang, "s", "en", "source language[en]")
	flag.StringVar(&targetLang, "t", "fr", "target language[fr]")
	flag.StringVar(&sourceText, "text", "", "Text to translate")
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	strChan := make(chan string)

	wg.Add(1)

	reqBody := &cli.RequestBody{
		SourceLang: sourceText,
		TargetLang: targetLang,
		SourceText: sourceText,
	}

	go cli.RequestTranslate(reqBody, strChan, &wg)

	processedStr := strings.ReplaceAll(<-strChan, "+", " ")

	fmt.Printf("Translation: %s\n", processedStr)

	close(strChan)

	wg.Wait()
}
