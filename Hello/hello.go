package main

import (
	"fmt"
	"io"   
	"net/http"
	"golang.org/x/net/html"
	"os"
	"strings"
	)

func get_one_link(token html.Token) (ok bool, href string) {
	for _, a := range token.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func get_all_links(body io.Reader) []string {

	links := []string{}

	// ready to parse HTMLs
	page := html.NewTokenizer(body)

	// 
	for {
		nextToken := page.Next()

		switch {
			case nextToken == html.ErrorToken:
				fmt.Println("End of HTML reached...")
				return links
			case nextToken == html.StartTagToken:

				token := page.Token() // get <XXX>

				if token.Data == "a" {
					_, url := get_one_link(token)
					fmt.Println("Link found! ", url)
					links = append(links, url)
				}
			}
	}
	
	return links
}

func main(){

	// single argument as the target site to parse the files.
	if len(os.Args) < 2 {
		os.Exit(-1)
	}

	// target_url :: the first hop the crawling.
	// get the target url from the args.
	target_url := os.Args[1]

	// add http prefix if the url does not have it
	if strings.Index("http://", target_url) == -1 && strings.Index("https://", target_url) == -1 {
		fmt.Println("9001 No HTTP prefix...")
		target_url = "http://" + target_url
	}

	re, err := http.Get(target_url)

	if err != nil { 
		fmt.Println("9008 HTTP GET Error")
		os.Exit(-1)
	}

	defer re.Body.Close()

	// get the resulting links.
	links := get_all_links(re.Body)

	// for now, simply print all the links.
	// TODO: download the targeting files.
	for _, link := range(links) {
	 	fmt.Println(link)
	}
}
