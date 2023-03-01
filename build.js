import { fileURLToPath } from 'url'
import { dirname } from 'path'
import Metalsmith from 'metalsmith'
import layouts from '@metalsmith/layouts'
import permalinks from '@metalsmith/permalinks'
import collections from '@metalsmith/collections'
import frontMatter from 'front-matter'
import { Parser } from 'simple-text-parser'
import yaml from 'js-yaml'
import fs from 'fs'
import url from 'url'
import markdownIt from 'markdown-it'
import markdownItFootnote from 'markdown-it-footnote'
import markdownItContainer from 'markdown-it-container'
import replaceExtension from 'replace-ext'
import sass from 'metalsmith-sass'
import { DateTime } from 'luxon'
import inPlace from '@metalsmith/in-place'
import path from 'node:path'
import slugify from 'slugify'

export default async function build() {
  try {
    const __dirname = dirname(fileURLToPath(import.meta.url))
    const files = await Metalsmith(__dirname)
      .source('src')
      .destination('build')
      .clean(tue)
      .metadata({
        site: {
          domain: "raypeat.rodeo",
          github: "https://github.com/marcuswhybrow/ray-peat-rodeo",
          githubEdit: "https://github.com/marcuswhybrow/ray-peat-rodeo/edit/main"
        }
      })

      .use(function addDocumentMetadata(files, metalsmith) {
        metalsmith.match('documents/**').forEach(filepath => {
          files[filepath].githubEdit = [
            metalsmith.metadata().site.githubEdit,
            path.join(metalsmith._source, filepath)
          ].join('/')
        })
      })

      // Determine which lines are spoken by whom
      .use(function speakerSyntax(files, metalsmith) {
        metalsmith.match('documents/**/*.md').forEach(filepath => {
          const file = files[filepath]

          const lines = function() {
            const otherSpeakers = []
            let prevSpeakerKey = null
            const results = []
            const lines = file.contents.toString().split(/\r?\n/)
            const speakers = file.speakers || {}
            lines.forEach(line => {
              const speakerParagraphRegex = /(?<speakerDefinition>!(?<speakerKey>\S*))?\s*(?<paragraph>.*)$/
              const { speakerDefinition, speakerKey, paragraph } = speakerParagraphRegex.exec(line).groups
              if (speakerDefinition) {
                let computedSpeakerKey = speakerKey || prevSpeakerKey
                
                if (!computedSpeakerKey) {
                  computedSpeakerKey = (function defaultSpeakerKey() {
                    const speakerEntries = Object.entries(speakers)
                    switch (speakerEntries.length) {
                      case 0:
                        speakers.H = 'Host'
                        return 'H'
                      case 1:
                        return speakerEntries[0][0]
                      default:
                        return null
                    }
                  })()
                }

                if (computedSpeakerKey) {
                  if (!otherSpeakers.includes(computedSpeakerKey))
                    otherSpeakers.push(computedSpeakerKey)
                  const speakerName = speakers[computedSpeakerKey]
                  if (!speakerName)
                    throw new Error(`Unknown speaker key "${computedSpeakerKey}" for line:\n${line}`)
                  results.push({
                    type: 'Other Speaker',
                    speakerKey: computedSpeakerKey,
                    speakerName,
                    speakerNumber: otherSpeakers.indexOf(computedSpeakerKey),
                    text: paragraph
                  })
                } else {
                  throw new Error(`Multiple speakers defined, but line is missing speaker key:\n${line}`)
                }
                prevSpeakerKey = computedSpeakerKey
              } else {
                if (paragraph) {
                  results.push({
                    type: 'Ray Peat',
                    text: paragraph
                  })
                } else {
                  results.push({
                    type: 'Empty',
                    text: line
                  })
                }
              }
            })
            return results
          }()

          const groupedLines = (function() {
            if (lines.length === 0) return lines
          
            const results = [lines.shift()]
            lines.forEach(line => {
              const tailIndex = results.length - 1
              const tail = results[tailIndex]
              if (
                line.type === 'Empty' ||
                (line.type === 'Ray Peat' && tail.type === 'Ray Peat') ||
                (line.type === 'Other Speaker' && tail.type === "Other Speaker" && line.speakerKey === tail.speakerKey)
              ) {
                results[tailIndex] = {...tail, text: `${tail.text}\n${line.text}`}
              } else {
                results.push(line)
              }
            })
            return results
          })()

          const sections = groupedLines.map(line => {
            const sectionName = line.type === 'Ray Peat'
              ? 'ray'
              : line.type === 'Other Speaker'
                ? `speaker ${line.speakerNumber} ${line.speakerName}`
                : null
            if (sectionName) {
              return  `::: ${sectionName}\n${line.text}\n:::`
            } else {
              return line.text
            }
          })

          file.contents = Buffer.from(sections.join('\n'), 'utf8')
        })
      })

      // Process square bracket tags
      .use(function interviewSyntax(files, metalsmith) {
        const references = metalsmith.match('documents/**/*.md').map(filepath => {
          const file = files[filepath]
          const fileMarkdown = frontMatter(file.contents.toString())

          const interviewTagsParser = (() => {
            function getLibGenSearchUrl(query) {
              return 'https://libgen.is/search.php?req=' + encodeURIComponent(query)
            }
            function getSciHubSearchUrl(query) {
              return 'https://sci-hub.ru/' + encodeURIComponent(query)
            }
            function getGoogleSearchUrl(query) {
              return 'https://www.google.com/search?q=' + encodeURIComponent(query)
            }
            
            function computeNode (pipeAndDisplayText, fallback, f) {
              if (pipeAndDisplayText) {
                if (pipeAndDisplayText.length >= 1) return f(pipeAndDisplayText.substring(1) || fallback)
                return { type: "text", text: "" }
              }
              return f(fallback)
            }
            
            function getPerson(name) {
              const people = yaml.load(fs.readFileSync('./data/people.yml', 'utf8'))
              const searchExternallyAs = people[name] ? people[name].searchExternallyAs || name : name
              return {
                name,
                searchExternallyAs,
                libGenURL: getLibGenSearchUrl(searchExternallyAs),
                googleURL: getGoogleSearchUrl(searchExternallyAs)
              }
            }
          
            const parserRules = {
              internalLinks: [
                /\[\[(\|.*?)\]\]/gi,
                (tag, pipeAndDisplayText) => computeNode(pipeAndDisplayText, "", (displayText) => ({
                  type: "internal-link-broken",
                  text: displayText
                }))
              ],
              books: [
                /\[\[([^\]\|]*?)\s*-by-\s*([^\]\|]*?)(\|[^\]]*?)?\]\]/gi,
                (tag, bookTitle, primaryAuthorFullName, pipeAndDisplayText) => {
                  const book = (() => {
                    const authorLastName = primaryAuthorFullName.split(' ').pop()
                    const authorFirstNames = primaryAuthorFullName.substring(0, primaryAuthorFullName.lastIndexOf(' '))
                    const person = getPerson(primaryAuthorFullName)
                    const query = `${bookTitle} ${person.searchExternallyAs}`
                    const libGenURL = getLibGenSearchUrl(query)
                    const key = `${bookTitle} -by- ${primaryAuthorFullName}`
                    const books = yaml.load(fs.readFileSync('./data/books.yml', 'utf8'))
                    const bookData = books[key]
                    return {
                      title: bookTitle,
                      author: primaryAuthorFullName,
                      authorLastName,
                      authorFirstNames,
                      googleUrl: getGoogleSearchUrl(query),
                      url: bookData ? bookData.openAsURL || libGenURL : libGenURL,
                      linkTitle: bookData ? bookData.openAsURLMessage : null
                    }
                  })()

                  const indexKey = (() => {
                    const bibliographKeyFirstName = book.authorFirstNames ? `, ${book.authorFirstNames}` : ''
                    return `${book.authorLastName}${bibliographKeyFirstName}. ${bookTitle}`
                  })()
              
                  return computeNode(pipeAndDisplayText, bookTitle, (displayText) => ({
                    type: "book",
                    text: `<a href="${book.url}" target="_blank" class="book" title="${book.linkTitle}">${displayText}</a>`,
                    value: book,
                    indexKey,
                    friendlyName: `${book.title} by ${book.author}`
                  }))
                }
              ],
              sciencePapers: [
                /\[\[doi\:([^\]\|]*?)(\|[^\]]*?)?\]\]/gi,
                (tag, doi, pipeAndDisplayText) => {
                  const sciencePaper = (() => {
                    const dois = yaml.load(fs.readFileSync('./data/doi.yml', 'utf8'))
                    const sciHubUrl = getSciHubSearchUrl(doi)
                    return {
                      doi,
                      sciHubUrl,
                      url: dois[doi] ? dois[doi].url || sciHubUrl : sciHubUrl
                    }
                  })()
              
                  return computeNode(pipeAndDisplayText, sciencePaper.url, (displayText) => ({
                    type: "science-paper",
                    text: `<a href="${sciencePaper.url}" target="_blank" class="science-paper">${displayText}</a>`,
                    value: sciencePaper,
                    indexKey: sciencePaper.doi
                  }))
                }
              ],
              externalLinks: [
                /\[\[(https?\:\/\/[^\]\|]*?)(\|[^\]]*?)?\]\]/gi,
                (tag, url, pipeAndDisplayText) => {
                  return computeNode(pipeAndDisplayText, url, (displayText) => ({
                    type: "external-link",
                    text: `<a href="${url}" target="_blank" class="external">${displayText}</a>`,
                    value: url,
                    indexKey: url,
                    friendlyName: url
                  }))
                }
              ],
              people: [
                /\[\[(([^\|\]]*?)(?:['’]s?)?)(\|([^\]]*?))?\]\]/gi,
                (tag, nameAsWritten, nameWithoutPluralisation, pipeAndDisplayText) => {
                  const person = getPerson(nameWithoutPluralisation)
                  const indexKey = (() => {
                    const names = person.name.split(' ')
                    return `${names.pop()}${names.length ? `, ${names.join(' ')}` : ''}`
                  })()
                  return computeNode(pipeAndDisplayText, nameAsWritten, (displayText) => ({
                    type: "person",
                    text: `<a href="${person.libGenURL}" target="_blank" class="person">${displayText}</a>`,
                    value: person,
                    indexKey,
                    friendlyName: person.name
                  }))
                }
              ],
              timecode: [
                /\[(\d+)\:(\d+)(?:\:(\d+))?\]/gi,
                (tag, firstDigits, secondDigits, thirdDigits) => {
                  const { youTubeFormat, localFormat } = (() => {
                    let { h, m, s } = (() => {
                      if (thirdDigits)
                        return { h: firstDigits, m: secondDigits, s: thirdDigits }
                      return { h: '', m: firstDigits, s: secondDigits }
                    })()
                    
                    const hasHours = Boolean(h)
                    h = h.padStart(2, '0')
                    m = m.padStart(2, '0')
                    s = s.padStart(2, '0')
                
                    return {
                      youTubeFormat: hasHours ? `${h}h${m}m${s}s` : `${m}m${s}s`,
                      localFormat: hasHours ? `${h}:${m}:${s}` : `${m}:${s}`
                    }
                  })()
              
                  const timecode = {
                    youTubeFormat,
                    localFormat,
                    originalURL: file.source,
                    url: (() => {
                      // TODO: Don't assume source URL has no hash
                      const parsedSourceUrl = url.parse(file.source, true)
                      if (parsedSourceUrl.host.endsWith("youtube.com") || parsedSourceUrl.host.endsWith("youtu.be")) {
                        return `${file.source}#t=${youTubeFormat}`
                      } else if (parsedSourceUrl.pathname.endsWith(".mp3")) {
                        return `${file.source}#t=${localFormat}`
                      } else {
                        return `${file.source}#t=${localFormat}`
                      }
                    })()
                  }
              
                  return ({
                    type: "timecode",
                    text: `<a href="${timecode.url}" target="_blank" class="timecode">${timecode.localFormat}</a> `,
                    value: timecode
                  })
                }
              ]
            }
          
            const parser = new Parser()
            Object.values(parserRules).map(([regex, handler]) => parser.addRule(regex, handler))
            return parser
          })()

          // replace file contents with the tags rendered out
          file.contents = Buffer.from(interviewTagsParser.render(fileMarkdown.body), 'utf8')

          // extract and count references (four kinds)
          const fileReferences = Object.fromEntries(
            [
              ['people', 'person'],
              ['books', 'book'],
              ['externalLinks', 'external-link'],
              ['sciencePapers', 'science-paper']
            ].map(([category, tagName]) => [
              category, 
              interviewTagsParser
                .toTree(fileMarkdown.body)
                .filter(node => node.type === tagName)
                .reduce((tagCount, tag) => {
                  const extantKey = tagCount[tag.indexKey]
                  tagCount[tag.indexKey] = {
                    count: extantKey ? extantKey.count + 1 : 1,
                    friendlyName: tag.friendlyName,
                    value: tag.value
                  }
                  return tagCount
                }, {})
            ])
          )

          return [filepath, fileReferences]
        });

        // Add accumulated references to global metadata
        metalsmith.metadata({
          references: references.reduce((results, [filepath, fileReferences]) => {
            Object.entries(fileReferences).forEach(([category, categoryData]) => {
              results[category] = results[category] || {}
              Object.entries(categoryData).forEach(([indexKey, entryData]) => {
                const extantEntry = results[category][indexKey]
                const extantFiles = extantEntry ? extantEntry.files || {} : {}
                const extantFileCount = extantEntry ? extantEntry.fileCount || 0 : 0
                results[category][indexKey] = {
                  indexKey,
                  value: entryData.value,
                  friendlyName: entryData.friendlyName,
                  count: extantEntry ? extantEntry.count + entryData.count : entryData.count,
                  files: {...extantFiles, [filepath]: files[filepath]},
                  fileCount: extantFileCount + 1
                }
              })
            })
            return results
          }, {})
        })
      })

      // Render Markdown to HTML
      .use(function renderMarkdownToHtml(files, metalsmith) {
        const md = markdownIt({
          html: true,
          linkify: true,
          typographer: true,
        })
          .use(markdownItFootnote)
          .use(markdownItContainer, 'speaker', {
            render: function (tokens, idx) {
              var args = tokens[idx].info.trim().match(/^speaker\s+(\d*?)\s+(.*)$/);
          
              if (tokens[idx].nesting === 1) {
                return `<div class="speaker speaker-other speaker-other-${args[1]}">\n<span class="speaker-name">${args[2]}:</span>`;
              } else {
                return `</div>\n`;
              }
            }
          })
          .use(markdownItContainer, 'ray', {
            render: function (tokens, idx) {
              if (tokens[idx].nesting === 1) {
                return `<div class="speaker speaker-ray">\n<span class="speaker-name">Ray Peat:</span>`;
              } else {
                return `</div>\n`;
              }
            }
          });

        metalsmith.match('**/*.md').forEach(filepath => {
          const file = files[filepath]
          const contents = file.contents.toString()
          file.contents = Buffer.from(md.render(contents), 'utf8')
          const htmlFilepath = replaceExtension(filepath, '.html')
          files[htmlFilepath] = file
          delete files[filepath]
        })
      })

      .use(function useLuxonDateTime(files, metalsmith) {
        metalsmith.match('documents/**').forEach(filepath => {
          const file = files[filepath]
          file.date = DateTime.fromJSDate(file.date)
          file.transcription.date = DateTime.fromJSDate(file.transcription.date)
        })
      })
      .use(collections({
        transcripts: {
          pattern: 'documents/**',
          reverse: true,
          sortBy: 'date'
        }
      }))
      .use(function addGlobalMetadata(files, metalsmith) {
        Object.values(files).forEach(file => {
          const metadata = metalsmith.metadata()
          file.global = {
            site: metadata.site,
            references: metadata.references,
            collections: metadata.collections
          }
        })
      })

      .use(function computeDocumentPermalinks(files, metalsmith) {
        metalsmith.match('documents/**').forEach(filepath => {
          const file = files[filepath]
          file.filename = (() => {
            const match = slugify(path.basename(filepath)).match(/^(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})-(?<slug>[^\.]*)\.html$/i)
            if (!match) throw new Error(`${filepath} does not conform to filename format: <4 digit year>-<2 digit month>-<2 digit day>-<slug>.md`)
            return match.groups
          })()
          file.permalink = file.filename.slug
        })
      })

      .use(layouts({
        pattern: 'documents/**',
        directory: 'layouts',
        default: 'document.pug',
      }))
      
      .use(permalinks({
        relative: false,
      }))

      .use(sass({
        pattern: 'site/**',
        outputStyle: "expanded"
      }))

      .use(inPlace({
        pattern: 'site/**',
        setFilename: true
      }))

      .use(function moveSiteFilesToRoot(files, metalsmith) {
        metalsmith.match('site/**').forEach(filepath => {
          const newFilepath = filepath.substring(5, filepath.length)
          files[newFilepath] = files[filepath]
          delete files[filepath]
        })
      })

      .build(err => {
        if (err) throw err;
      })
    console.log('Build success')
    return files
  } catch (err) {
    console.error(err)
    return err
  }
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  // is main script
  build()
}
