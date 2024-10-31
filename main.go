package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/html"
)

const (
	URL = "https://www.politifact.com/"
)

func main() {
	fmt.Println("web scrapper up and running")
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("Error occured while fetching the url:", err)
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("Error occured while parsing the html:", err)
	}
	var w sync.WaitGroup
	w.Add(1)
	go traverse(doc, &w)
	w.Wait()
}

func traverse(node *html.Node, w *sync.WaitGroup) {
	if node == nil {
		return
	}
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				if isValidURL(attr.Val) {
					resp, err := http.Get(attr.Val)
					if err != nil {
						fmt.Println("Error occured:", attr.Val)
						return
					}
					if resp.StatusCode == http.StatusOK {
						fmt.Println("dead link url:", attr.Val)
					}
				}
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		w.Add(1)
		go traverse(c, w)
	}
	defer w.Done()
}
func isValidURL(link string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}
