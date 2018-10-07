package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		fmt.Printf("Parsing %s\n", url)
		links, err := getLinks(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing links in %s: %s\n", url, err.Error())
			continue
		}

		fmt.Printf("%s Links:\n", url)
		for _, link := range links {
			fmt.Printf("\t%s\n", link)
		}
	}
}

func getLinks(url string) (links []string, err error) {
	links = make([]string, 0)

	// Get the html from the url.
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body := resp.Body
	defer body.Close()

	// Parse the html.
	count := 0
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		switch {
		case tokenType == html.ErrorToken:
			// End of html
			fmt.Printf("Counted %d anchors\n", count)
			return
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()
			if token.Data != "a" {
				continue
			}
			for i := 0; i < len(token.Attr); i++ {
				if token.Attr[i].Key == "href" && strings.Index(token.Attr[i].Val, "http") == 0 {
					count++
					links = append(links, token.Attr[i].Val)
					break
				}
			}

		}
	}
}
