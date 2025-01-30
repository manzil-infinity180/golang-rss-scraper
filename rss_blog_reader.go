package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
)

func GetPost(url string) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	//feed, err := parser.ParseURLWithContext(url)
	if err != nil {
		log.Println(fmt.Errorf("error parsing %s: %w", url, err))
		return
	}
	log.Println(feed.Items)
	return
}
