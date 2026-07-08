package httpx

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	stringsutil "github.com/projectdiscovery/utils/strings"
	"golang.org/x/net/html"
	"slices"
)

var (
	cutset                  = "\n\t\v\f\r"
	reTitle                 = regexp.MustCompile(`(?im)<\s*title.*>(.*?)<\s*/\s*title>`)
	reContentType           = regexp.MustCompile(`(?im)\s*charset="(.*?)"|charset=(.*?)"\s*`)
	supportedTitleMimeTypes = []string{
		"text/html",
		"application/xhtml+xml",
		"application/xml",
		"application/rss+xml",
		"application/atom+xml",
		"application/xhtml+xml",
		"application/vnd.wap.xhtml+xml",
	}
)

// ExtractTitle from a response
func ExtractTitle(r *Response) (title string) {
	// Try to parse the DOM
	titleDom, err := getTitleWithDom(r)
	// In case of error fallback to regex
	if err != nil {
		for _, match := range reTitle.FindAllString(r.Raw, -1) {
			title = match
			break
		}
	} else {
		title = renderNode(titleDom)
	}

	title = html.UnescapeString(trimTitleTags(title))

	// remove unwanted chars
	title = strings.TrimSpace(strings.Trim(title, cutset))
	title = stringsutil.ReplaceAll(title, "", "\n", "\t", "\v", "\f", "\r")

	return title
}

func CanHaveTitleTag(mimeType string) bool {  
    return slices.Contains(supportedTitleMimeTypes, mimeType)  
}  

func getTitleWithDom(r *Response) (*html.Node, error) {
	var title *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" {
			title = node
			return
		}
		for child := node.FirstChild; child != nil && title == nil; child = child.NextSibling {
			crawler(child)
		}
	}
	htmlDoc, err := html.Parse(bytes.NewReader(r.Data))
	if err != nil {
		return nil, err
	}
	crawler(htmlDoc)
	if title != nil {
		return title, nil
	}
	return nil, fmt.Errorf("title not found")
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n) //nolint
	return buf.String()
}

func trimTitleTags(title string) string {
    titleBegin := strings.Index(title, ">")
    titleEnd := strings.Index(title, "</")
    if titleEnd < 0 || titleBegin < 0 || titleEnd <= titleBegin {
        return title
    }
    return title[titleBegin+1 : titleEnd]
}

var (
	reMetaCopyright = regexp.MustCompile(`(?im)<\s*meta\s[^>]*name\s*=\s*["']?(?:copyright|rights)["'\s][^>]*content\s*=\s*["']([^"']*)["']`)
	reCopySymbol    = regexp.MustCompile(`(?im)(©|&copy;)([^<]{0,300})|(?:\bcopyright\b\s*(?:[^<]{0,200}))`)
)

// ExtractCopyright from a response body
func ExtractCopyright(r *Response) string {
	// 1. Try <meta name="copyright" or "rights">
	metaMatch := reMetaCopyright.FindStringSubmatch(r.Raw)
	if len(metaMatch) > 1 && strings.TrimSpace(metaMatch[1]) != "" {
		return strings.TrimSpace(metaMatch[1])
	}

	// 2. Try DOM parsing for meta tags
	doc, err := html.Parse(bytes.NewReader(r.Data))
	if err == nil {
		if cp := extractCopyrightMetaFromDom(doc); cp != "" {
			return cp
		}
		if cp := extractCopyrightFooterFromDom(doc); cp != "" {
			return cp
		}
	}

	// 3. Fallback: regex in raw HTML for © or "copyright"
	match := reCopySymbol.FindString(r.Raw)
	if match != "" {
		match = html.UnescapeString(match)
		match = strings.TrimSpace(match)
		match = strings.Trim(match, cutset)
		if len(match) > 300 {
			match = match[:300]
		}
		return match
	}

	return ""
}

func extractCopyrightMetaFromDom(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var name, content string
		for _, attr := range n.Attr {
			switch attr.Key {
			case "name":
				name = strings.ToLower(attr.Val)
			case "content":
				content = attr.Val
			}
		}
		if (name == "copyright" || name == "rights") && content != "" {
			return strings.TrimSpace(content)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if cp := extractCopyrightMetaFromDom(c); cp != "" {
			return cp
		}
	}
	return ""
}

func extractCopyrightFooterFromDom(n *html.Node) string {
	var search func(*html.Node) string
	search = func(node *html.Node) string {
		if node.Type == html.ElementNode {
			isFooter := node.Data == "footer" || node.Data == "small"
			if isFooter {
				if cp := extractCopyrightLine(node); cp != "" {
					return cp
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if cp := search(c); cp != "" {
				return cp
			}
		}
		return ""
	}
	return search(n)
}

func extractCopyrightLine(n *html.Node) string {
	// Walk all text nodes under this element and return the line matching © or copyright
	var result string
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			if reCopySymbol.MatchString(text) {
				// Get some context from parent
				text = html.UnescapeString(text)
				text = strings.TrimSpace(text)
				text = strings.Trim(text, cutset)
				if len(text) > 300 {
					text = text[:300]
				}
				result = text
				return
			}
		}
		if result == "" {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				walk(c)
			}
		}
	}
	walk(n)
	return result
}
