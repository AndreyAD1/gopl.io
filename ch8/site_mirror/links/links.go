package links

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"os"
	// "strings"

	"golang.org/x/net/html"
)

func Extract(initalURL string) ([]string, error) {
	parsedURL, err := url.Parse(initalURL)
	if err != nil {
		return nil, fmt.Errorf("invalid link: %s", initalURL)
	}
	resp, err := http.Get(initalURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", initalURL, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", initalURL, err)
	}


	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}
				links = append(links, link.String())
				localLink := getLocalLink(link, parsedURL)
				a.Val = localLink
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	err = writePageToFile(doc, parsedURL)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func writePageToFile(page *html.Node, link *url.URL) error {
	file, err := getFile(link)
	if err != nil {
		return err
	}
	defer file.Close()
	err = html.Render(file, page)
	if err != nil {
		return err
	}
	fmt.Println("successfully make the file")
	return nil
}

func getFile(parsedURL *url.URL) (*os.File, error) {
	dirPath := filepath.Join("mirrors", parsedURL.Host, path.Dir(parsedURL.Path))
	filePath := filepath.Join(dirPath, path.Base(parsedURL.Path)) + ".html"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return nil, fmt.Errorf("can not create the directory %s: %v", dirPath, err)
		}
		file, err := os.Create(filePath)
		if err != nil {
			err := fmt.Errorf("can not create the file %s: %v", filePath, err)
			return nil, err
		}
		return file, err
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		err = fmt.Errorf("can not open the file: %s", filePath)
		return nil, err
	}
	return file, err
}

func getLocalLink(link, currentURL *url.URL) string {
	return link.String()
	// if link.Hostname() != currentURL.Hostname() {
	// 	return link.String()
	// }
	// linkPath := strings.Split(link.Path)
	// linkPath = append(linkPath, linkFile)
}
