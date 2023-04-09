package citations

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcuswhybrow/ray-peat-rodeo/lib/utils"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func fromTemplate(b util.BufWriter, templateText string, data interface{}) {
	t := utils.ReturnOrPanic(template.New("template").Parse(templateText))
	utils.PanicOnErr(t.Execute(b, data))
}

type BaseCitation struct {
	Position     int
	CitationKey  string
	ExternalLink string
}

type BookPrimaryAuthor struct {
	Name        string
	FirstNames  string
	LastName    string
	CitationKey string
}

type Book struct {
	Citation      BaseCitation
	Title         string
	PrimaryAuthor BookPrimaryAuthor
}

type Person struct {
	Citation   BaseCitation
	Name       string
	LastName   string
	FirstNames string
}

type SciencePaper struct {
	Citation      BaseCitation
	Title         string
	Year          int
	PrimaryAuthor string
	Doi           string
}

type ExternalLink struct {
	Citation BaseCitation
	Title    string
}

var NewLibGenSearchUrl = func(query string) string {
	return "https://libgen.is/search.php?req=" + query
}

var NewSciHubUrl = func(query string) string {
	return "https://sci-hub.ru/" + query
}

var peopleContextKey = parser.NewContextKey()
var booksContextKey = parser.NewContextKey()
var externalLinksContextKey = parser.NewContextKey()
var sciencePapersContextKey = parser.NewContextKey()

type citationsParser struct {
}

func NewParser() parser.InlineParser {
	return &citationsParser{}
}

func (w *citationsParser) Trigger() []byte {
	return []byte{'['}
}

func ensureContextMapExists[A any](pc parser.Context, contextKey parser.ContextKey) map[string]A {
	value := pc.Get(contextKey)
	if value != nil {
		return value.(map[string]A)
	} else {
		value := map[string]A{}
		pc.Set(contextKey, value)
		return value
	}
}

type CitationKind int64

const (
	PersonCitation CitationKind = iota
	BookCitation
	SciencePaperCitation
	ExternalLinkCitation
)

func (w *citationsParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	people := ensureContextMapExists[Person](pc, peopleContextKey)
	books := ensureContextMapExists[Book](pc, booksContextKey)
	sciencePapers := ensureContextMapExists[SciencePaper](pc, sciencePapersContextKey)
	externalLinks := ensureContextMapExists[ExternalLink](pc, externalLinksContextKey)

	line, _ := block.PeekLine()

	if line[1] != '[' {
		return nil
	}

	inside, post, foundEnd := strings.Cut(string(line[2:]), "]]")
	if !foundEnd || len(inside) == 0 {
		return nil
	}

	rawCitationKey, rawModifier := func() (string, *string) {
		if rawCitaionKey, rawModifier, found := strings.Cut(inside, "|"); found {
			return rawCitaionKey, &rawModifier
		}
		return inside, nil
	}()

	rawPostText := func() string {
		if i := strings.IndexFunc(post, func(r rune) bool {
			return unicode.IsSpace(r) || unicode.IsPunct(r)
		}); i >= 0 {
			return post[:i]
		}
		return post
	}()
	bytesConsumed := 2 + len(inside) + 2 + len(rawPostText)

	if rawCitationKey == "" {
		if rawModifier == nil {
			// [[]]
			return nil
		} else if *rawModifier == "" {
			// [[|]]
			return nil
		} else {
			// [[|text]] placeholder
			block.Advance(bytesConsumed)
			return ast.NewString([]byte(*rawModifier))
		}
	} else {
		if rawModifier != nil && *rawModifier == "" {
			// [[key|]] silent citation
			block.Advance(bytesConsumed)
			return ast.NewString([]byte{})
		}
	}

	kind := func() CitationKind {
		if strings.HasPrefix(rawCitationKey, "doi:") {
			return SciencePaperCitation
		} else if strings.HasPrefix(rawCitationKey, "http://") || strings.HasPrefix(rawCitationKey, "https://") {
			return ExternalLinkCitation
		} else if strings.Contains(rawCitationKey, "-by-") {
			return BookCitation
		} else {
			return PersonCitation
		}
	}()

	displayText := func(citationKeyOverride *string) string {
		citationKey := func() string {
			if citationKeyOverride != nil {
				return *citationKeyOverride
			}
			switch kind {
			case PersonCitation:
				return rawCitationKey
			case ExternalLinkCitation:
				return rawCitationKey
			case BookCitation:
				return strings.TrimSpace(strings.Split(rawCitationKey, "-by-")[0])
			case SciencePaperCitation:
				if len(rawCitationKey) > 4 {
					return rawCitationKey[4:]
				} else {
					panic("Science Paper Citation has no DOI")
				}
			}
			panic("Unknown Citation Kind")
		}()

		if rawModifier != nil {
			if *rawModifier != "" {
				return *rawModifier + rawPostText
			} else {
				return ""
			}
		} else {
			return citationKey + rawPostText
		}
	}

	// [[key]] or [[key|]] or [[key|text]]
	switch kind {
	case SciencePaperCitation:
		doi := rawCitationKey[4:]
		externalLink := NewSciHubUrl(doi)
		title := make(chan string, 1)
		go func() {
			doiUrl := "https://doi.org/" + doi
			req := utils.ReturnOrPanic(http.NewRequest("GET", doiUrl, nil))
			req.Header.Set("Accept", "application/json; charset=utf-8")
			res, err := http.DefaultClient.Do(req)

			if err != nil {
				fmt.Println("Could not get title for url:", doiUrl)
				fmt.Println(err)
				title <- doi
			} else {
				data := utils.ReturnOrPanic(ioutil.ReadAll(res.Body))
				res.Body.Close()

				doiJson := map[string]interface{}{}
				json.Unmarshal(data, &doiJson)
				title <- doiJson["title"].(string)
			}
		}()
		block.Advance(bytesConsumed)
		return NewCitationNode(withLinkRenderer(func() map[string]string {
			t := <-title
			sciencePapers[doi] = SciencePaper{
				Citation: BaseCitation{
					Position:     0,
					ExternalLink: externalLink,
					CitationKey:  doi,
				},
				Title:         t,
				Year:          0,
				PrimaryAuthor: "",
				Doi:           doi,
			}

			return map[string]string{
				"class": "citation science-paper",
				"href":  externalLink,
				"text":  displayText(&t),
			}
		}))
	case BookCitation:
		bookTitle, primaryAuthor := func() (string, string) {
			parts := strings.Split(rawCitationKey, "-by-")
			if len(parts) < 2 {
				panic("Book Citation has no primary author")
			}
			return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		}()
		firstNames, lastName := func() (string, string) {
			names := strings.Split(primaryAuthor, " ")
			lastName := names[len(names)-1]
			firstNames := strings.Join(names[:len(names)-1], " ")
			return firstNames, lastName
		}()
		authorKey := func() string {
			if len(strings.TrimSpace(firstNames)) == 0 {
				return lastName
			}
			return lastName + ", " + firstNames
		}()
		citationKey := lastName + ", " + firstNames + ". " + bookTitle
		externalLink := NewLibGenSearchUrl(bookTitle + " " + primaryAuthor)
		block.Advance(bytesConsumed)
		return NewCitationNode(withLinkRenderer(func() map[string]string {
			books[citationKey] = Book{
				Citation: BaseCitation{
					Position:     0,
					ExternalLink: externalLink,
					CitationKey:  citationKey,
				},
				Title: bookTitle,
				PrimaryAuthor: BookPrimaryAuthor{
					Name:        primaryAuthor,
					FirstNames:  firstNames,
					LastName:    lastName,
					CitationKey: authorKey,
				},
			}

			return map[string]string{
				"class": "citation book",
				"href":  externalLink,
				"text":  displayText(nil),
			}
		}))
	case ExternalLinkCitation:
		externalLink := rawCitationKey
		block.Advance(bytesConsumed)
		title := make(chan string, 1)
		go func() {
			doc, err := goquery.NewDocument(externalLink)
			if err != nil {
				fmt.Println("Could not extract title from url:", externalLink)
				fmt.Println(err)
				title <- externalLink
			} else {
				title <- strings.Trim(doc.Find("title").Text(), " \n")
			}
		}()
		return NewCitationNode(withLinkRenderer(func() map[string]string {
			t := <-title
			externalLinks[rawCitationKey] = ExternalLink{
				Citation: BaseCitation{
					Position:     0,
					ExternalLink: externalLink,
					CitationKey:  externalLink,
				},
				Title: t,
			}

			return map[string]string{
				"class": "citation external-link",
				"href":  externalLink,
				"text":  displayText(&t),
			}
		}))
	case PersonCitation:
		externalLink := NewLibGenSearchUrl(rawCitationKey)
		firstNames, lastName := func() (string, string) {
			names := strings.Split(rawCitationKey, " ")
			lastName := names[len(names)-1]
			firstNames := strings.Join(names[:len(names)-1], " ")
			return firstNames, lastName
		}()
		citationKey := func() string {
			if len(strings.TrimSpace(firstNames)) == 0 {
				return lastName
			}
			return lastName + ", " + firstNames
		}()
		block.Advance(bytesConsumed)
		return NewCitationNode(withLinkRenderer(func() map[string]string {
			people[citationKey] = Person{
				Citation: BaseCitation{
					Position:     0,
					ExternalLink: externalLink,
					CitationKey:  citationKey,
				},
				Name:       rawCitationKey,
				LastName:   lastName,
				FirstNames: firstNames,
			}

			return map[string]string{
				"class": "citation person",
				"href":  externalLink,
				"text":  displayText(nil),
			}
		}))
	}
	panic("Unknown Citation Kind")
}
