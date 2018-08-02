package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/net/html"
)

func travelNodes(node *html.Node, handler func(*html.Node)) {
	// fmt.Printf("next sibling %v\n", node.NextSibling)
	for ; node != nil; node = node.NextSibling {
		// fmt.Printf("node %v\n", 1)
		handler(node)
		travelNodes(node.FirstChild, handler)
	}
}

func generateMinnorLocation(origin string) string {
	// fmt.Printf("ext %v\n", path.Ext(origin))
	dir := path.Dir(origin)
	base := path.Base(origin)
	if base == "" || base == "." {
		base = "index"
	}

	ext := path.Ext(origin)

	const DATA = "./minnor"
	if ext != "" {
		return DATA + dir + "/" + base
	}
	return DATA + dir + "/" + base + ".html"
}

// create a file at location.
// return file hanlder
func openFile(location string) (*os.File, error) {
	minnorLocation := generateMinnorLocation(location)
	dir := path.Dir(minnorLocation)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("creating folder %v: %v", dir, err)
	}

	file, err := os.Create(minnorLocation)
	if err != nil {
		return nil, fmt.Errorf("creating file %s: %v", minnorLocation, err)
	}
	return file, nil
}

func crawl(link string) (links []string) {
	fmt.Printf("crawling %v\n", link)
	res, err := http.Get(link)
	if err != nil {
		fmt.Printf("make get request: %v", err)
		return nil
	}
	defer res.Body.Close()

	file, err := openFile(res.Request.URL.Path)
	if err != nil {
		fmt.Printf("can not open file %v", err)
		return
	}
	defer file.Sync()
	defer file.Close()

	node, err := html.Parse(res.Body)
	if err != nil {
		fmt.Printf("parse response: %v", err)
		return nil
	}

	saveLink := func(node *html.Node) {
		// fmt.Printf("node %v\n", node)
		if node.Type == html.ElementNode && node.Data == "a" {
			for i := range node.Attr {
				attr := node.Attr[i]
				if attr.Key == "href" {
					link, err := url.Parse(attr.Val)
					if err != nil {
						fmt.Printf("getting url %v\n", err)
						return
					}
					reference := res.Request.URL.ResolveReference(link)
					if reference.Host != res.Request.Host {
						return
					}

					rel, err := filepath.Rel(path.Dir(generateMinnorLocation(res.Request.URL.Path)), generateMinnorLocation(reference.Path))
					if err != nil {
						fmt.Printf("generating rel path: %v", err)
					}
					// fmt.Printf("updating href. Current %s Old: %s New %s\n", generateMinnorLocation(res.Request.URL.Path), generateMinnorLocation(reference.Path), rel)

					node.Attr[i].Val = rel
					links = append(links, reference.String())
					// fmt.Printf("attr after updating %#v\n", node.Attr)
					break
				}
			}
			// fmt.Printf("after updating %#v\n", *node)
		}
	}

	travelNodes(node, saveLink)
	html.Render(file, node)
	return links
}

func main() {
	workList := make(chan []string)
	tokens := make(chan struct{}, 20)
	seen := make(map[string]bool)

	n := 1
	go func() {
		fmt.Printf("args %v\n", os.Args[1:])
		workList <- os.Args[1:]
	}()

	for ; n > 0; n-- {
		list := <-workList
		for _, link := range list {
			if seen[link] {
				continue
			}

			seen[link] = true
			tokens <- struct{}{}
			n++
			go func(link string) {
				list := crawl(link)
				<-tokens
				workList <- list
			}(link)
		}
	}
}
