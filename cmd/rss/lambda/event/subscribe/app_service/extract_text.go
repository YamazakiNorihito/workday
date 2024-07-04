package app_service

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func extractText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	if n.Type == html.ElementNode && n.Data == "br" {
		buf.WriteString("\n")
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, buf)
	}
}

func extractTextFromHTML(htmlContent string) (string, error) {
	if htmlContent == "" {
		return "", nil
	}
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	extractText(doc, &buf)

	return buf.String(), nil
}
