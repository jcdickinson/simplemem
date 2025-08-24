package memory

import (
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

var (
	wikiLinkRegex = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
)

type Link struct {
	Text   string
	Target string
	Type   string
}

func ExtractLinks(content string) []Link {
	var links []Link

	// First extract wiki-style links
	wikiMatches := wikiLinkRegex.FindAllStringSubmatch(content, -1)
	for _, match := range wikiMatches {
		if len(match) > 1 {
			target := match[1]
			if !strings.HasSuffix(target, ".md") {
				target = target + ".md"
			}
			links = append(links, Link{
				Text:   match[1],
				Target: target,
				Type:   "wiki",
			})
		}
	}

	// Parse markdown to extract regular links
	p := parser.New()
	doc := p.Parse([]byte(content))

	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if link, ok := node.(*ast.Link); ok {
				dest := string(link.Destination)
				// Only include .md file links
				if strings.HasSuffix(dest, ".md") {
					text := extractTextFromNodes(link.Children)
					links = append(links, Link{
						Text:   text,
						Target: dest,
						Type:   "markdown",
					})
				}
			}
		}
		return ast.GoToNext
	})

	return links
}

func extractTextFromNodes(nodes []ast.Node) string {
	var text strings.Builder
	for _, node := range nodes {
		if textNode, ok := node.(*ast.Text); ok {
			text.Write(textNode.Literal)
		} else if codeNode, ok := node.(*ast.Code); ok {
			text.Write(codeNode.Literal)
		}
	}
	return text.String()
}

func ParseMarkdown(content string) []byte {
	p := parser.New()
	doc := p.Parse([]byte(content))
	return markdown.Render(doc, nil)
}

func ResolveLinks(content string, basePath string) string {
	// Convert wiki-style links to markdown links
	result := wikiLinkRegex.ReplaceAllStringFunc(content, func(match string) string {
		innerMatch := wikiLinkRegex.FindStringSubmatch(match)
		if len(innerMatch) > 1 {
			target := innerMatch[1]
			if !strings.HasSuffix(target, ".md") {
				target = target + ".md"
			}
			return "[" + innerMatch[1] + "](" + target + ")"
		}
		return match
	})

	return result
}

func GetBacklinks(store *Store, targetName string) ([]string, error) {
	memories, err := store.List()
	if err != nil {
		return nil, err
	}

	var backlinks []string
	targetWithExt := targetName
	if !strings.HasSuffix(targetWithExt, ".md") {
		targetWithExt = targetWithExt + ".md"
	}

	for _, memory := range memories {
		content, err := store.Read(memory)
		if err != nil {
			continue
		}

		links := ExtractLinks(content)
		for _, link := range links {
			if link.Target == targetWithExt || link.Target == targetName {
				backlinks = append(backlinks, memory)
				break
			}
		}
	}

	return backlinks, nil
}