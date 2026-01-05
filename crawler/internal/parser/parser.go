package parser

import (
	"crawler/pkg/models"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HTMLParser struct{}

// constructor for Parser
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// Parse method to extract title and content from HTML
func (p *HTMLParser) Parse(pageURL string, html string) (*models.Document, error) {
	// Parse the HTML content using goquery package
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	title := doc.Find("title").First().Text()

	// Extract text content
	content := extractText(doc)

	//extract domain from URL
	parsedURL, _ := url.Parse(pageURL)
	domain := parsedURL.Host

	//count words in content
	wordCount := len(strings.Fields(content))

	return &models.Document{
		URL:       pageURL,
		Title:     title,
		Content:   content,
		Domain:    domain,
		WordCount: wordCount,
	}, nil
}

func extractText(doc *goquery.Document) string {
	//remove script and style tags
	doc.Find("script, style").Remove()

	//get the text from body
	text := doc.Find("body").Text()

	//clean up whitespace
	text = strings.Join(strings.Fields(text), " ")

	return text
}

func (p *HTMLParser) ExtractLinks(html, baseURL string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return []string{}
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return []string{}
	}

	links := []string{}
	seen := make(map[string]bool)

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if href == "" {
			return
		}

		linkURL, err := url.Parse(href)
		if err != nil {
			return
		}

		if linkURL.Host != "" && linkURL.Host != base.Host {
			return
		}

		if !seen[href] {
			seen[href] = true
			if linkURL.Scheme == "" {
				linkURL.Scheme = base.Scheme
			}
			if linkURL.Host == "" {
				linkURL.Host = base.Host
			}
			if linkURL.Path == "" {
				linkURL.Path = "/"
			}
			links = append(links, linkURL.String())
		}
	})

	return links
}
