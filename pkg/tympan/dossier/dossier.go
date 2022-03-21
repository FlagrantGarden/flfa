package dossier

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/charmbracelet/glamour"
	"github.com/gernest/front"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

// A Dossier contains everything needed to display documentation in the terminal
type Dossier struct {
	// The handler for interacting with a filesystem - local, remote, or otherwise.
	AFS *afero.Afero
	// The handler for interacting with documents embedded in the binary itself
	EFS *embed.FS
	// The handler for writing documents to the terminal with style
	TerminalRenderer *glamour.TermRenderer
	// The handler for parsing a markdown document with frontmatter
	MarkdownHandler *front.Matter
	// The cache of parsed documents in the Dossier
	Documents []Document
}

// A Document is sourced from markdown either embedded in the binary or available via another filesystem.
type Document struct {
	// Metadata about the Document
	FrontMatter FrontMatter
	// The text of the Document in markdown
	Body string
}

// The Title of a Document has two formats for convenience.
type Title struct {
	// A short downcased string without special characters; used for topic selection
	Short string
	// The full string for the Document's title; human readable
	Long string
}

// All Documents in a Dossier have required frontmatter; the source markdown document may have other metadata but it
// *must* have these fields.
type FrontMatter struct {
	// The human-readable title for the document and the short name suitable for shell completion
	Title Title
	// A short synopsis of the document's contents
	Description string
	// Whether this document is Conceptual, Narrative, or something else
	Category string
	//  Freeform list of strings to help users filter and find documents
	Tags []string
}

// To implement your own Dossier, you must be able to read markdown documents, render them, list them all, list them
// filtered by category, list them filtered by tag, and cache them.
// TODO: Use actually needed methods
type DossierI interface {
	// CacheDocuments searches the Dossier's embedded file system at the specified folder path, parsing all discovered
	// markdown documents and caching them in the Dossier.
	CacheDocuments(docsFolderPath string)
	// RenderDocument returns the terminal-rendered markdown of the specified document but does not print it itself.
	RenderDocument(doc Document) (string, error)
	// FormatFrontMatter returns either the string representation of the JSON object containing the specified frontmatter
	// or the rendered markdown table representation of the frontmatter to the caller but does not print either itself.
	FormatFrontMatter(format string, docs []Document)
	// SelectDocument returns the Document which has the specified short title. If no Document matches, it errors.
	SelectDocument(shortTitle string) (document Document, err error)
	// CompleteShortTitle returns the list of short titles which are a valid match for the specified string as a prefix.
	CompleteShortTitle(match string) (titles []string)
	// ListCategories returns a slice of strings containing all of the unique categories used by the Documents in the
	// Dossier's cache.
	ListCategories(docs []Document) (categories []string)
	// ListTags returns a slice of strings containing all of the unique tags used by the Documents in the Dossier's cache.
	ListTags(docs []Document) (tags []string)
}

// InitializeMarkdownHandler creates the handler for markdown documents if it does not already exist and defines how it
// should behave.
func (dossier *Dossier) InitializeMarkdownHandler() {
	if dossier.MarkdownHandler == nil {
		dossier.MarkdownHandler = front.NewMatter()
		dossier.MarkdownHandler.Handle("---", front.YAMLHandler)
	}
}

// CacheEmbeddedDocuments searches the Dossier's embedded file system at the specified folder path, parsing all
// discovered markdown documents and caching them in the Dossier.
func (dossier *Dossier) CacheDocuments(docsFolderPath string) {
	// ignore errors for now
	dirEntries, _ := dossier.EFS.ReadDir(docsFolderPath)
	for _, entry := range dirEntries {
		entryPath := fmt.Sprintf("%s/%s", docsFolderPath, entry.Name())

		if entry.IsDir() {
			dossier.CacheDocuments(entryPath)
		} else {
			log.Debug().Msgf("Parsing Documentation File: %s", entryPath)
			// read file, ignoring errors for now
			raw, _ := dossier.EFS.ReadFile(entryPath)
			// parse and append to cache
			dossier.InitializeMarkdownHandler()
			fm, b, err := dossier.MarkdownHandler.Parse(strings.NewReader(string(raw)))
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
				dossier.Documents = append(dossier.Documents, Document{
					Body: b,
					FrontMatter: FrontMatter{
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

// RenderDocument returns the terminal-rendered markdown of the specified document but does not print it itself.
func (dossier *Dossier) RenderDocument(doc Document) (string, error) {
	// Add the title since it's captured in frontmatter and not raw markdown
	var bodyWithTitle strings.Builder
	bodyWithTitle.WriteString(fmt.Sprintf("# %s\n", doc.FrontMatter.Title.Long))
	bodyWithTitle.WriteString(doc.Body)
	return dossier.Render(bodyWithTitle.String())
}

// FormatFrontMatter returns either the string representation of the JSON object containing the specified frontmatter or
// the rendered markdown table representation of the frontmatter to the caller but does not print either itself.
func (dossier *Dossier) FormatFrontMatter(format string, docs []Document) {
	var frontMatterList []FrontMatter
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
		out, _ := dossier.Render(table.String())
		fmt.Print(out)
	}
}

// SelectDocument returns the Document which has the specified short title. If no Document matches, it errors.
func (dossier *Dossier) SelectDocument(shortTitle string) (document Document, err error) {
	for _, doc := range dossier.Documents {
		if doc.FrontMatter.Title.Short == shortTitle {
			document = doc
		}
	}
	if document.FrontMatter.Title.Short == "" {
		err = fmt.Errorf("could not find document with short title: %s", shortTitle)
	}
	return document, err
}

// CompleteShortTitle returns the list of short titles which are a valid match for the specified string as a prefix.
func (dossier *Dossier) CompleteShortTitle(match string) (titles []string) {
	for _, title := range dossier.ListTitles(dossier.Documents) {
		if strings.HasPrefix(title.Short, match) {
			titles = append(titles, fmt.Sprintf("%s\t%s", title.Short, title.Long))
		}
	}
	return
}

// ListCategories returns a slice of strings containing all of the unique categories used by the Documents in the
// Dossier's cache.
func (dossier *Dossier) ListCategories(docs []Document) (categories []string) {
	for _, doc := range docs {
		if !utils.Contains(categories, doc.FrontMatter.Category) {
			categories = append(categories, doc.FrontMatter.Category)
		}
	}
	return categories
}

// ListTags returns a slice of strings containing all of the unique tags used by the Documents in the Dossier's cache.
func (dossier *Dossier) ListTags(docs []Document) (tags []string) {
	for _, doc := range docs {
		for _, tag := range doc.FrontMatter.Tags {
			if !utils.Contains(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

// InitializeTerminalRenderer creates the terminal renderer if it does not already exist and defines how it should
//behave.
func (dossier *Dossier) InitalizeTerminalRenderer() (err error) {
	if dossier.TerminalRenderer == nil {
		dossier.TerminalRenderer, err = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(120),
		)
	}
	return err
}

// Render initializes the terminal renderer if needed and then uses it to render the body of a markdown document,
// returning the rendered string for printing to the screen but not printing it itself.
func (dossier *Dossier) Render(body string) (output string, err error) {
	err = dossier.InitalizeTerminalRenderer()
	if err != nil {
		return "", err
	}
	output, err = dossier.TerminalRenderer.Render(body)
	return output, err
}

// ListTitles returns a slice of Titles (the short and long names) from every Document in the Dossier's cache.
func (dossier *Dossier) ListTitles(docs []Document) (titles []Title) {
	for _, doc := range docs {
		titles = append(titles, doc.FrontMatter.Title)
	}
	return titles
}

// FilterByTag returns a slice of Documents which have the specified tag.
func FilterByTag(tag string, documents []Document) (filteredDocs []Document) {
	for _, doc := range documents {
		if utils.Contains(doc.FrontMatter.Tags, tag) {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs
}

// FilterByCategory returns a slice of Documents which have the specified category.
func FilterByCategory(category string, documents []Document) (filteredDocs []Document) {
	for _, doc := range documents {
		if doc.FrontMatter.Category == category {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs
}
