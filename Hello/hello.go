package main

import (
	"fmt"
	"io"   
	"net/http"
	"golang.org/x/net/html"
	"os"
	"strings"
	)

const DOWNLOADDIR string = "/download/"
const PDF_SUFFIX string =".pdf"

// download a single file from url
func download_a_file(url string) {

	// request for the file content.
	rep, err := http.Get(url)

	if err != nil {
		fmt.Println("9008 HTTP GET Error trying downloading...")
		os.Exit(-1)
	}

	defer rep.Body.Close()

	// path of download directory
	download_dir := os.Getenv("GOPATH")+DOWNLOADDIR

	if _, err = os.Stat(download_dir); os.IsNotExist(err) {
		fmt.Println("creating download dir")
		os.Mkdir(download_dir, 0700)
	}

	// file_name := $GOPATH/desktop/download/***.pdf
	file_name := download_dir+url[strings.LastIndex(url,"/")+1:]

	// create a file with file_name
	file, err := os.Create(string(file_name))
	if err != nil {
		fmt.Println("9003 File OS Error: Create")
		os.Exit(-1)
	}

	// copy the content into the file created.
	_, err = io.Copy(file, rep.Body)
	if err != nil {
		fmt.Println("9003 File OS Error: Copy")
		os.Exit(-1)
	}

	file.Close()
	fmt.Println("Success!")
}

/*
Get one signle url from <a *href="xxx"*>
*/
func get_one_link(token html.Token) (ok bool, href string) {

	// loop through token's attributes.
	ok = false;
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

	// iterate the HTMLs to find specific tags
	// here we first look for <a>
	for {
		nextToken := page.Next()

		switch {
			case nextToken == html.ErrorToken:

				fmt.Println("End of HTML reached...")
				return links

			case nextToken == html.StartTagToken:

				token := page.Token() // get <XXX>

				// explicitly looking for <a>
				if token.Data == "a" {
					_, url := get_one_link(token)
					//fmt.Println("Link found! ", url)
					links = append(links, url)
				}
			}
	}
	
	return links
}

func main(){

	// single argument as the target site to parse the files.
	if len(os.Args) < 2 {
		fmt.Println("Usage: \n\thello http://xxxx")
		os.Exit(-1)
	}

	// target_url :: the first hop the crawling.
	// get the target url from the args.
	target_url := os.Args[1]

	// add http prefix if the url does not have it
	if strings.Index(target_url, "http://") == -1 && strings.Index(target_url, "https://") == -1 {
		fmt.Println("9001 No HTTP prefix...")
		target_url = "http://" + target_url
	}

	// inital http request.
	re, err := http.Get(target_url)

	if err != nil { 
		fmt.Println("9008 HTTP GET Error")
		os.Exit(-1)
	}

	defer re.Body.Close()

	// get the resulting links.
	links := get_all_links(re.Body)

	// for now, simply print all the links.
	// TODO: download the nested files
	for _, link := range(links) {
	 	fmt.Println(link)

	 	// download the pdf files
	 	if strings.Index(link, PDF_SUFFIX) != -1 {
	 		prefix := string(target_url[0:strings.LastIndex(target_url,"/")+1])
	 		download_a_file(prefix+link)
	 	}
	}
}
