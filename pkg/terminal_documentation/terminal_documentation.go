package terminal_documentation

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/utils"
	"github.com/charmbracelet/glamour"
	"github.com/gernest/front"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type TerminalDocumentation struct {
	AFS             *afero.Afero
	IOFS            *afero.IOFS
	TermRenderer    *glamour.TermRenderer
	EmbeddedDocsFS  *embed.FS
	MarkdownHandler *front.Matter
	ParsedDocsCache []MarkdownDocument
}

type MarkdownDocument struct {
	FrontMatter DocumentFrontMatter
	Body        string
}

type Title struct {
	Short string
	Long  string
}

type DocumentFrontMatter struct {
	Title       Title
	Description string
	Category    string
	Tags        []string
}

type TerminalDocumentationer interface {
	ReadMarkdownDoc(text string) (mdc MarkdownDocument)
	InitRenderer() (err error)
	Render(body string) (output string, err error)
	ListByCategory(category string)
	ListByTag(tag string)
	List()
	FindAndParse(docsFolderPath string)
}

func (td *TerminalDocumentation) InitRenderer() (err error) {
	if td.TermRenderer == nil {
		td.TermRenderer, err = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(120),
		)
	}
	return err
}

func (td *TerminalDocumentation) Render(body string) (output string, err error) {
	err = td.InitRenderer()
	if err != nil {
		return "", err
	}
	output, err = td.TermRenderer.Render(body)
	return output, err
}

func (td *TerminalDocumentation) InitHandler() {
	if td.MarkdownHandler == nil {
		td.MarkdownHandler = front.NewMatter()
		td.MarkdownHandler.Handle("---", front.YAMLHandler)
	}
}

func (td *TerminalDocumentation) FindAndParse(docsFolderPath string) {
	// ignore errors for now
	dirEntries, _ := td.EmbeddedDocsFS.ReadDir(docsFolderPath)
	for _, entry := range dirEntries {
		entryPath := fmt.Sprintf("%s/%s", docsFolderPath, entry.Name())

		if entry.IsDir() {
			td.FindAndParse(entryPath)
		} else {
			log.Debug().Msgf("Parsing Documentation File: %s", entryPath)
			// read file, ignoring errors for now
			raw, _ := td.EmbeddedDocsFS.ReadFile(entryPath)
			// parse and append to cache
			td.InitHandler()
			fm, b, err := td.MarkdownHandler.Parse(strings.NewReader(string(raw)))
			if err != nil {
				log.Warn().Msgf("Could not parse %s", entryPath)
				// Some docs might need to be skipped in terminal
			} else if skip, _ := fm["skipTerminal"]; skip != true { //nolint
				// Turn the tags into an array of strings for further use
				tagsAsInterface := fm["tags"].([]interface{})
				tags := make([]string, len(tagsAsInterface))
				for i, v := range tagsAsInterface {
					tags[i] = v.(string)
				}
				td.ParsedDocsCache = append(td.ParsedDocsCache, MarkdownDocument{
					Body: b,
					FrontMatter: DocumentFrontMatter{
						Title: Title{
							Short: strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())),
							Long:  fm["title"].(string),
						},
						Description: fm["description"].(string),
						Category:    fm["category"].(string),
						Tags:        tags,
					},
				})
			}
		}
	}
}

func (td *TerminalDocumentation) ListTags(docs []MarkdownDocument) (tags []string) {
	for _, doc := range docs {
		for _, tag := range doc.FrontMatter.Tags {
			if !utils.Contains(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

func (td *TerminalDocumentation) ListCategories(docs []MarkdownDocument) (categories []string) {
	for _, doc := range docs {
		if !utils.Contains(categories, doc.FrontMatter.Category) {
			categories = append(categories, doc.FrontMatter.Category)
		}
	}
	return categories
}

func (td *TerminalDocumentation) ListTitles(docs []MarkdownDocument) (titles []Title) {
	for _, doc := range docs {
		titles = append(titles, doc.FrontMatter.Title)
	}
	return titles
}

func (td *TerminalDocumentation) CompleteTitle(docs []MarkdownDocument, match string) []string {
	var titles []string
	for _, title := range td.ListTitles(docs) {
		if strings.HasPrefix(title.Short, match) {
			titles = append(titles, fmt.Sprintf("%s\t%s", title.Short, title.Long))
		}
	}
	return titles
}

func (td *TerminalDocumentation) FilterByTag(tag string, docs []MarkdownDocument) (filteredDocs []MarkdownDocument) {
	for _, doc := range docs {
		if utils.Contains(doc.FrontMatter.Tags, tag) {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs
}

func (td *TerminalDocumentation) FilterByCategory(category string, docs []MarkdownDocument) (filteredDocs []MarkdownDocument) {
	for _, doc := range docs {
		if doc.FrontMatter.Category == category {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs
}

func (td *TerminalDocumentation) SelectDocument(shortTitle string, docs []MarkdownDocument) (document MarkdownDocument, err error) {
	for _, doc := range docs {
		if doc.FrontMatter.Title.Short == shortTitle {
			document = doc
		}
	}
	if document.FrontMatter.Title.Short == "" {
		err = fmt.Errorf("could not find document with short title: %s", shortTitle)
	}
	return document, err
}

func (td *TerminalDocumentation) FormatFrontMatter(format string, docs []MarkdownDocument) {
	var frontMatterList []DocumentFrontMatter
	for _, doc := range docs {
		frontMatterList = append(frontMatterList, doc.FrontMatter)
	}
	switch format {
	case "json":
		fm, _ := json.Marshal(frontMatterList)
		fmt.Print(string(fm))
	default:
		var table strings.Builder
		table.WriteString("| Name | Description | Category | Tags |\n")
		table.WriteString("| ---- | ----------- | -------- | ---- |\n")
		for _, doc := range frontMatterList {
			entry := fmt.Sprintf("| %s | %s | %s | %s |\n", doc.Title.Short, doc.Description, doc.Category, strings.Join(doc.Tags, ", "))
			table.WriteString(entry)
		}
		out, _ := td.Render(table.String())
		fmt.Print(out)
	}
}

func (td *TerminalDocumentation) RenderDocument(doc MarkdownDocument) (string, error) {
	// Add the title since it's captured in frontmatter and not raw markdown
	var bodyWithTitle strings.Builder
	bodyWithTitle.WriteString(fmt.Sprintf("# %s\n", doc.FrontMatter.Title.Long))
	bodyWithTitle.WriteString(doc.Body)
	return td.Render(bodyWithTitle.String())
}
