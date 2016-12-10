package main

import (
	"fmt"
	"net/http"
	"time"

	"strings"

	"golang.org/x/net/html"
)

func main() {
	urlProcessor := make(chan string)
	done := make(chan bool)
	timeout := make(chan bool, 1)

	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	go processURL(urlProcessor, timeout, done)

	urlProcessor <- "https://jeremywho.com"

	<-done
	fmt.Println("Done")
}

func processURL(urlProcessor chan string, timeout chan bool, done chan bool) {
	visited := make(map[string]bool)
	for {
		select {
		case url := <-urlProcessor:
			if _, ok := visited[url]; ok {
				continue
			} else {
				visited[url] = true
				go exploreURL(url, urlProcessor)
			}
		case <-timeout:
			fmt.Printf("Explored %d pages\n", len(visited))
			done <- true
		}
	}
}

func exploreURL(url string, urlProcessor chan string) {
	fmt.Printf("Visiting %s.\n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return
		}

		if tt == html.StartTagToken {
			t := z.Token()

			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {

						// if link is within jeremywho.com
						if strings.HasPrefix(a.Val, "https://jeremywho.com") {
							urlProcessor <- a.Val
						}
					}
				}
			}
		}
	}
}
