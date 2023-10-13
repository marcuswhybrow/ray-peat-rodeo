use std::{path::{PathBuf, Path}, fs::DirEntry, collections::HashMap};

use chrono::NaiveDate;
use extract_frontmatter::{Extractor, config::Splitter::EnclosingLines};
use markdown_it::{parser::extset::MarkdownItExt, MarkdownIt};

use crate::{GITHUB_LINK, BasePage, markdown::{Frontmatter, mention::{Mention, MentionDeclaration}, github::{GitHubIssue, GitHubIssueDeclaration}}, stash::Stash, write};


#[derive(Debug, Hash, PartialEq, Eq, Clone)]
pub struct InputFileDate {
    pub year: String,
    pub month: String,
    pub day: String,
    pub as_date: NaiveDate,
}

impl From<InputFileDate> for String {
    fn from(value: InputFileDate) -> Self {
        value.to_string()
    }
}

impl ToString for InputFileDate {
    fn to_string(&self) -> String {
        format!("{}-{}-{}", self.year, self.month, self.day)
    }
}

impl InputFileDate {
    fn new(path: &PathBuf) -> Self {
        let date = &path.file_stem().unwrap().to_str().unwrap()[..=9];
        let year = &date[..=3];
        let month = &date[5..=6];
        let day = &date[8..=9];

        InputFileDate {
            year: year.to_string(),
            month: month.to_string(),
            day: day.to_string(),
            as_date: {
                let year: i32 = {
                    if year == "????" {
                        0
                    } else {
                        year.parse()
                            .unwrap_or_else(|_| panic!("Invalid year {}", year))
                    }
                };

                let month: u32 = {
                    if month == "??" {
                        1
                    } else {
                        month.parse()
                            .unwrap_or_else(|_| panic!("Invalid month {}", month))
                    }
                };

                let day: u32 = {
                    if day == "??" {
                        1
                    } else {
                        day.parse()
                            .unwrap_or_else(|_| panic!("Invalud day {}", day))
                    }
                };

                NaiveDate::from_ymd_opt(year, month, day)
                    .unwrap_or_else(|| panic!("Invalid date for {}", date))
            }
        }
    }
}

#[derive(Debug, Hash, PartialEq, Eq, Clone)]
pub struct ContentFile {
    pub path: PathBuf, 
    pub filename_date: InputFileDate,
    pub title_slug: String,
    pub frontmatter: Frontmatter,
    pub markdown: String,
}

impl TryFrom<DirEntry> for ContentFile {
    type Error = &'static str;

    fn try_from(entry: DirEntry) -> Result<Self, Self::Error> {
        let path = entry.path();
        let extension = path.extension();

        if path.is_dir() {
            return Err("Path was a directory");
        }

        if extension.is_none() || extension.is_some_and(|ext| {
            let lc_ext = ext.to_ascii_lowercase();
            !(lc_ext == "md" || lc_ext == "markdown")
        }) {
            return Err("Path does not have a markdown extension");
        }

        if path.file_stem().unwrap().to_ascii_uppercase() == "README" {
            return Err("Path was a README markdown file");
        }

        let text = std::fs::read_to_string(&path)
            .expect("Path could not be read to string");

        let (raw_frontmatter, markdown) = Extractor::new(EnclosingLines("---"))
            .extract(text.as_str());

        Ok(ContentFile {
            filename_date: InputFileDate::new(&path),
            title_slug: path.file_stem().unwrap().to_str().unwrap()
                .split_at(11).1.to_string(),
            markdown: markdown.to_string(),
            frontmatter: serde_yaml::from_str(&raw_frontmatter)
                .unwrap_or_else(|e| panic!("Invalid YAML frontmatter in {:?}\n{e}", path)),
            path,
        })
    }

}

impl ContentFile {
    pub fn as_content<'a>(&self, parser: &mut MarkdownIt, stash: &'a mut Stash) -> Content {
        parser.ext.insert(ContentFileBeingParsed(self.clone()));
        let mut ast = parser.parse(self.markdown.as_str());
        let mut mention_count: HashMap<MentionDeclaration, u32> = HashMap::new();
        let mut mentions = vec![];
        let mut issues = vec![];

        ast.walk_mut(|node, _depth| {
            if let Some(mention_declaration) = node.cast::<MentionDeclaration>() {
                *mention_count.entry((*mention_declaration).clone()).or_insert(0) += 1;
                let count = mention_count.get(mention_declaration).unwrap();
                let mention = mention_declaration.as_mention(self.clone(), count.clone(), stash);
                node.replace(mention.clone());
                mentions.push(mention);

            } else if let Some(github_issue_declaration) = node.cast::<GitHubIssueDeclaration>() {
                let issue = github_issue_declaration.as_github_issue(stash);
                node.replace(issue.clone());
                issues.push(issue);
            }
        });

        Content { file: self.clone(), ast, mentions, issues }
    }

    pub fn is_todo(&self) -> bool {
        self.path.parent().unwrap().file_name().unwrap() == "todo"
    }

    pub fn permalink(&self) -> String {
        format!("/{}", self.title_slug)
    }
}

/// An input markdownfile from ./content 
#[derive(Debug)]
pub struct Content {
    pub file: ContentFile,
    pub ast: markdown_it::Node,
    pub mentions: Vec<Mention>,
    #[allow(dead_code)]
    pub issues: Vec<GitHubIssue>,
}

impl Content {
    pub fn write(&self, output: &Path) {
        let path = output.join(format!("{}/index.html", self.file.title_slug));
        write(&path,  markup::new! { @ContentPage { content: self } });
    }
}

#[derive(Debug)]
pub struct ContentFileBeingParsed(pub ContentFile);
impl MarkdownItExt for ContentFileBeingParsed {}

markup::define! {
    ContentPage<'a>(content: &'a Content) {
        @BasePage {
            title: Some(content.file.frontmatter.source.title.as_str()),
            content: markup::new! {
                article .interview ["data-pagefind-body"] {
                    header {
                        div.hud ["data-pagefind-ignore"] {
                            span.date { @content.file.filename_date.to_string() }
                            a.series [
                                href = format!("/series/{}", content.file.frontmatter.source.series.to_lowercase().replace(" ", "-"))
                            ] {
                                @content.file.frontmatter.source.series
                            }
                        }

                        h1.title { @content.file.frontmatter.source.title }

                        @if let Some(transcription) = content.file.frontmatter.transcription.clone()  {
                            @if let Some(author) = transcription.author.clone() {
                                div."transcription-attribution" ["data-pagefind-ignore"] {
                                    @if let Some(url) = transcription.url.clone() {
                                        "Thank you to " @author " who "
                                        a[href = url] {
                                            "published"
                                        }
                                        " an earlier version of this transcript "
                                    } else {
                                        "Transcribed by " @author ", "
                                    }

                                    @ if let Some(date) = transcription.date.clone() {
                                        @date "."
                                    }
                                }
                            }
                        }

                        div.actions ["data-pagefind-ignore"] {
                            a."view-source"[
                                href = content.file.frontmatter.source.url.clone()
                            ] { "View Source" }
                            a."edit"[
                                href = format!("{GITHUB_LINK}/edit/main/{}", content.file.path.to_str().unwrap())
                            ] { "Edit on GitHub" }
                        }

                    }

                    main {
                        @if !content.file.is_todo() {
                            @markup::raw(content.ast.render())
                        } else {
                            section {
                                p {
                                    br {}
                                    "A transcript for this interview doesn't yet exist." 

                                    @if let Some(t) = &content.file.frontmatter.transcription {
                                        @if let Some (url) = &t.url {
                                            " However an external transcript is available here:\n"
                                            a [href = url] { @url }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            },
        }
    }
}
