package citations

import (
	"html/template"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

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

var citationsContextKey = parser.NewContextKey()

type CitationsContext struct {
	Count         int
	People        map[string]Person
	Books         map[string]Book
	SciencePapers map[string]SciencePaper
	ExternalLinks map[string]ExternalLink
}

func Get(pc parser.Context) CitationsContext {
	if meta := pc.Get(citationsContextKey); meta != nil {
		return meta.(CitationsContext)
	}
	return CitationsContext{
		Count:         0,
		People:        map[string]Person{},
		Books:         map[string]Book{},
		SciencePapers: map[string]SciencePaper{},
		ExternalLinks: map[string]ExternalLink{},
	}
}

type citationsParser struct {
}

func NewParser() parser.InlineParser {
	return &citationsParser{}
}

func (w *citationsParser) Trigger() []byte {
	return []byte{'['}
}

func (w *citationsParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	incrementCitationsCount := func() int {
		citations := Get(pc)
		citations.Count = citations.Count + 1
		pc.Set(citationsContextKey, citations)
		return citations.Count
	}
	updateContext := func(f func(citations *CitationsContext)) {
		extantCitations := Get(pc)
		f(&extantCitations)
		pc.Set(citationsContextKey, extantCitations)
	}

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
	if rawCitationKey == "" && rawModifier != nil && *rawModifier == "" {
		return nil
	}

	rawPostText := func() string {
		if i := strings.IndexFunc(post, func(r rune) bool {
			return unicode.IsSpace(r) || unicode.IsPunct(r)
		}); i >= 0 {
			return post[:i]
		}
		return post
	}()
	bytesConsumed := 2 + len(inside) + 2 + len(rawPostText)

	// The rawCitationKey determines the logic and rendering
	node := func() ast.Node {
		if rawCitationKey == "" {
			if rawModifier == nil {
				// [[]]
				return nil
			} else if *rawModifier == "" {
				// [[|]]
				return nil
			} else {
				// [[|text]] placeholder
				return ast.NewString([]byte(*rawModifier))
			}
		} else {
			citation := func() *Citation {
				fromTemplate := func(b util.BufWriter, templateText string, data interface{}) {
					t, err := template.New("template").Parse(templateText)
					if err != nil {
						panic(err)
					}
					err = t.Execute(b, data)
					if err != nil {
						panic(err)
					}
				}
				computedDisplayText := func(defaultDisplayText string) string {
					if rawModifier != nil {
						if *rawModifier != "" {
							return *rawModifier + rawPostText
						} else {
							return ""
						}
					} else {
						return defaultDisplayText + rawPostText
					}
				}

				type CitationFromTemplateConfig struct {
					displayTextFromCitationKey string
					entering                   string
					exiting                    string
					data                       interface{}
				}

				newCitationFromTemplates := func(config CitationFromTemplateConfig) *Citation {
					return NewCitation(CitationConfig{
						position: incrementCitationsCount(),
						text:     computedDisplayText(config.displayTextFromCitationKey),
						entering: func(b util.BufWriter) {
							fromTemplate(
								b,
								config.entering,
								config.data,
							)
						},
						exiting: func(b util.BufWriter) {
							fromTemplate(
								b,
								config.exiting,
								config.data,
							)
						},
					})
				}

				// [[key]] or [[key|]] or [[key|text]]
				if strings.HasPrefix(rawCitationKey, "doi:") {
					doi := rawCitationKey[4:]
					externalLink := NewSciHubUrl(doi)
					updateContext(func(citations *CitationsContext) {
						citations.SciencePapers[doi] = SciencePaper{
							Citation: BaseCitation{
								Position:     0,
								ExternalLink: externalLink,
								CitationKey:  doi,
							},
							Title:         "",
							Year:          0,
							PrimaryAuthor: "",
							Doi:           doi,
						}
					})
					return newCitationFromTemplates(CitationFromTemplateConfig{
						displayTextFromCitationKey: doi,
						entering: `<a
							class="citation science-paper"
							data-doi="{{.doi}}"
							href="{{.href}}"
							target="_blank"
						>`,
						exiting: `</a>`,
						data: map[string]interface{}{
							"doi":  doi,
							"href": externalLink,
						},
					})
				} else if i := strings.LastIndex(rawCitationKey, "-by-"); i >= 0 {
					bookTitle := strings.TrimSpace(rawCitationKey[:i])
					primaryAuthor := strings.TrimSpace(rawCitationKey[i+4:])
					names := strings.Split(primaryAuthor, " ")
					lastName := names[len(names)-1]
					firstNames := strings.Join(names[:len(names)-1], " ")
					authorKey := func() string {
						if len(strings.TrimSpace(firstNames)) == 0 {
							return lastName
						}
						return lastName + ", " + firstNames
					}()
					citationKey := lastName + ", " + firstNames + ". " + bookTitle
					externalLink := NewLibGenSearchUrl(bookTitle + " " + primaryAuthor)
					updateContext(func(citations *CitationsContext) {
						citations.Books[citationKey] = Book{
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
					})
					return newCitationFromTemplates(CitationFromTemplateConfig{
						displayTextFromCitationKey: bookTitle,
						entering: `<a
							class="citation book"
							data-title="{{.title}}"
							data-primary-author="{{.primaryAuthor}}"
							href="{{.href}}"
							target="_blank"
						>`,
						exiting: `</a>`,
						data: map[string]interface{}{
							"title":         bookTitle,
							"primaryAuthor": primaryAuthor,
							"href":          externalLink,
						},
					})
				} else if strings.HasPrefix(rawCitationKey, "http://") || strings.HasPrefix(rawCitationKey, "https://") {
					externalLink := rawCitationKey
					updateContext(func(citations *CitationsContext) {
						citations.ExternalLinks[externalLink] = ExternalLink{
							Citation: BaseCitation{
								Position:     0,
								ExternalLink: externalLink,
								CitationKey:  externalLink,
							},
							Title: "",
						}
					})
					return newCitationFromTemplates(CitationFromTemplateConfig{
						displayTextFromCitationKey: externalLink,
						entering: `<a
							class="citation external-link"
							href="{{.href}}"
							target="_blank"
						>`,
						exiting: `</a>`,
						data: map[string]interface{}{
							"href": externalLink,
						},
					})
				} else {
					externalLink := NewLibGenSearchUrl(rawCitationKey)
					names := strings.Split(rawCitationKey, " ")
					lastName := names[len(names)-1]
					firstNames := strings.Join(names[:len(names)-1], " ")
					citationKey := func() string {
						if len(strings.TrimSpace(firstNames)) == 0 {
							return lastName
						}
						return lastName + ", " + firstNames
					}()
					updateContext(func(citations *CitationsContext) {
						citations.People[rawCitationKey] = Person{
							Citation: BaseCitation{
								Position:     0,
								ExternalLink: externalLink,
								CitationKey:  citationKey,
							},
							Name:       rawCitationKey,
							LastName:   lastName,
							FirstNames: firstNames,
						}
					})
					return newCitationFromTemplates(CitationFromTemplateConfig{
						displayTextFromCitationKey: rawCitationKey,
						entering: `<a
							class="citation person"
							data-name="{{.name}}"
							href="{{.href}}"
							target="_blank"
						>`,
						exiting: `</a>`,
						data: map[string]interface{}{
							"name": rawCitationKey,
							"href": externalLink,
						},
					})
				}
			}()
			if rawModifier != nil && *rawModifier == "" {
				// [[key|]] silent citation
				return ast.NewString([]byte{})
			} else {
				return citation
			}
		}
	}()

	block.Advance(bytesConsumed)
	return node
}

type Citation struct {
	ast.BaseInline
	position       int
	text           string
	renderEntering func(b util.BufWriter)
	renderExiting  func(b util.BufWriter)
}

func (t *Citation) Dump(source []byte, level int) {
	ast.DumpHelper(t, source, level, nil, nil)
}

var KindCitation = ast.NewNodeKind("Citation")

func (n *Citation) Kind() ast.NodeKind {
	return KindCitation
}

type CitationConfig struct {
	position int
	entering func(b util.BufWriter)
	exiting  func(b util.BufWriter)
	text     string
}

func NewCitation(config CitationConfig) *Citation {
	return &Citation{
		BaseInline:     ast.BaseInline{},
		position:       config.position,
		renderEntering: config.entering,
		renderExiting:  config.exiting,
		text:           config.text,
	}
}

type CitationRenderer struct {
}

func NewCitationHTMLRenderer() renderer.NodeRenderer {
	return &CitationRenderer{}
}

func (r *CitationRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindCitation, r.renderPersonCitation)
}

func (t *CitationRenderer) renderPersonCitation(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	citation := node.(*Citation)
	if entering {
		citation.renderEntering(w)
		w.WriteString(citation.text)
	} else {
		citation.renderExiting(w)
	}
	return ast.WalkContinue, nil
}

type citations struct {
}

// Citations is an extension for Goldmark that replaces [[tags]] with links
var Citations = &citations{}

func New() goldmark.Extender {
	return &citations{}
}

func (e *citations) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewParser(), 100),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewCitationHTMLRenderer(), 100),
	))
}
