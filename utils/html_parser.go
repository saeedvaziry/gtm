package utils

import (
	"strings"

	"golang.org/x/net/html"
)

func ExtractHTMLAttributes(c string) map[string]string {
	var attributes = make(map[string]string)
	doc, err := html.Parse(strings.NewReader(c))
	if err != nil {
		panic(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			for _, attr := range n.FirstChild.Attr {
				attributes[attr.Key] = attr.Val
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return attributes
}
