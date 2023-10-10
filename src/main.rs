mod markdown;
mod scraper;

use std::{
    fs,
    path::{Path, PathBuf}, collections::HashMap,
};
use chrono::NaiveDate;
use clap::Parser;
use fs_extra::dir::CopyOptions;
use markdown::{
    Frontmatter,
    mention::{MentionDeclaration, Mention, Author}, github::{GitHubIssueDeclaration, GitHubIssue}
};
use markdown_it::{MarkdownIt, parser::extset::MarkdownItExt, Node};
use markup::DynRender;
use extract_frontmatter::{Extractor, config::Splitter::EnclosingLines};
use scraper::{Scraper, ScraperKind};

extern crate fs_extra;

pub const GITHUB_LINK: &str = "https://github.com/marcuswhybrow/ray-peat-rodeo";
pub const MENTION_SLUG: &str = "mentions";

#[derive(Parser, Debug)]
#[command(name = "Ray Peat Rodeo Engine")]
#[command(author = "Marcus Whybrow <marcus@whybrow.uk>")]
#[command(about = "Builds Ray Peat Rodeo into HTML from source")]
#[command(long_about = None)]
struct Args {
    /// The input path containing markdown content
    #[arg(short, long, default_value_t = String::from("./content"))]
    input: String,

    /// The output path in which Ray Peat Rodeo should build HTML
    #[arg(short, long, default_value_t = String::from("./build"))]
    output: String,

    #[arg(long, default_value_t = String::from("./src/cache.yml"))]
    cache_path: String,

    #[arg(long, default_value_t = false)]
    build_cache: bool,

    /// Whether files and directories inside of OUT_PATH should be deleted
    /// before building
    #[arg(long, default_value_t = false)]
    clean: bool,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();

    let input = Path::new(&args.input);
    let output = Path::new(&args.output);

    let cache_path = Path::new(&args.cache_path);

    println!("Ray Peat Rodeo");
    println!("  Input:  {:?}", input);
    println!("  Output: {:?}", output);
    println!("  Cache:  {:?}", cache_path);

    if args.build_cache {
        println!("  Building cache of web scraped data...");

        // In "Scraper" mode, mock data is returned, but a record is kept of 
        // all URLs requested.
        let mut scraper = Scraper::new(cache_path, ScraperKind::Scraper);

        // So we construct all the pages as normal, but do nothing with them.
        let _ = OutputPages::new(input, &mut scraper);
        let _ = OutputPages::new(input.join("todo").as_path(), &mut scraper);

        // The we write the cache to disk, which is used, when building the 
        // pages from real, using a scraper in Fulfiller mode. Probably could 
        // make this simplier, and easier to understand, but it works for now.
        scraper.into_cache().await;

        println!("  Built cache at {:?}", cache_path);
        return Ok(())
    }

    if !output.exists() {
        println!("  Creating output directory.");
        fs::create_dir(output.clone()).unwrap();
    }

    if args.clean {
        println!(
            "  Clean option enabled. Deleting files and directories inside {:?}",
            output
        );

        for entry in fs::read_dir(output.clone()).unwrap() {
            let entry = entry.unwrap();
            let path = entry.path();

            if entry.file_type().unwrap().is_dir() {
                fs::remove_dir_all(path.clone())
                    .expect(format!("Could not remove directory {:?}", path).as_str());
            } else {
                fs::remove_file(path.clone())
                    .expect(format!("Could note remove file {:?}", path).as_str());
            }
        }
    }

    println!("\nWriting Files");

    let mut scraper = Scraper::new(&cache_path, ScraperKind::Fulfiller);

    let mut output_pages = OutputPages::new(input, &mut scraper).0;
    output_pages.sort_by_key(|p| p.input_file.date.as_date);
    output_pages.reverse();
    
    let mut todo_output_pages = OutputPages::new(input.join("todo").as_path(), &mut scraper).0;
    todo_output_pages.sort_by_key(|p| p.input_file.date.as_date);
    todo_output_pages.reverse();

    output_pages.extend(todo_output_pages);

    for output_page in output_pages.iter() {
        render(&output.join(format!("{}/index.html", output_page.input_file.slug)), markup::new! {
            @Page { output_page }
        });
    }

    {
        let mut authors: HashMap<Author, Vec<Mention>> = HashMap::new();
        for output_page in output_pages.iter() {
            if let Some(mentions) = &output_page.mentions {
                for mention in mentions {
                    authors.entry(mention.author.clone()).or_insert(vec![]).push(mention.clone());
                }
            }
        }

        for (author, author_mentions) in authors {
            render(&output.join(format!("{MENTION_SLUG}/{}/index.html", author.id())), markup::new! {
                @AuthorPage {
                    author: author.clone(),
                    mentions: author_mentions.clone(),
                }
            });
        }
    }

    {
        let mut series_map: HashMap<&String, Vec<&InputFile>> = HashMap::new();
        for output_page in output_pages.iter() {
            series_map
                .entry(&output_page.input_file.frontmatter.source.series)
                .or_insert(vec![])
                .push(&output_page.input_file);
        }
        for (series, input_files) in series_map.iter() {
            render(
                &output.join(
                    format!("series/{}/index.html", series.to_lowercase().replace(" ", "-"))
                ),
                markup::new! {
                    @SeriesPage { series, input_files }
                }
            );
        }
    }


    render(&output.join("index.html"), markup::new! {
        @Homepage { output_pages: &output_pages }
    });

    copy_dir(&PathBuf::from("./src/assets"), &output.join("assets"));

    println!("\nDone");

    Ok(())
}


#[derive(Debug, Hash, PartialEq, Eq, Clone)]
pub struct InputFileDate {
    year: String,
    month: String,
    day: String,
    as_date: NaiveDate,
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

/// An input markdownfile from ./content 
#[derive(Debug, Hash, PartialEq, Eq, Clone)]
pub struct InputFile {
    todo: bool,
    date: InputFileDate,
    path: String,
    slug: String,
    frontmatter: Frontmatter,
    markdown: String,
}

impl InputFile {
    fn new(path: PathBuf) -> Result<Self, &'static str> {
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

        Ok(Self {
            todo: path.parent().unwrap().file_name().unwrap() == "todo",
            date: InputFileDate::new(&path),
            path: path.to_str().unwrap().to_string(),
            markdown: markdown.into(),
            slug: path.file_stem().unwrap().to_str().unwrap().split_at(11).1.to_string(),
            frontmatter: match serde_yaml::from_str(&raw_frontmatter) {
                Ok(f) => f,
                Err(e) => panic!("Invalid YAML frontmatter in {:?}\n{e}", path),
            },
        })
    }
}

pub struct OutputPage {
    input_file: InputFile,
    html_ast: Option<Node>,
    mentions: Option<Vec<Mention>>,

    #[allow(dead_code)]
    issues: Option<Vec<GitHubIssue>>,
}

impl OutputPage {
    fn new(input_file: InputFile, parser: &mut MarkdownIt, scraper: &mut &mut Scraper) -> Self {
        parser.ext.insert(InputFileBeingParsed(input_file.clone()));

        let mut ast = parser.parse(input_file.markdown.as_str());
        let mut mention_count: HashMap<MentionDeclaration, u32> = HashMap::new();
        let mut mentions = vec![];
        let mut issues = vec![];

        ast.walk_mut(|node, _depth| {
            if let Some(mention_declaration) = node.cast::<MentionDeclaration>() {
                *mention_count.entry((*mention_declaration).clone()).or_insert(0) += 1;
                let count = mention_count.get(mention_declaration).unwrap();
                let mention = mention_declaration.as_mention(input_file.clone(), count.clone(), scraper);
                node.replace(mention.clone());
                mentions.push(mention);

            } else if let Some(github_issue_declaration) = node.cast::<GitHubIssueDeclaration>() {
                let issue = github_issue_declaration.as_github_issue(scraper);
                node.replace(issue.clone());
                issues.push(issue);
            }
        });

        OutputPage {
            input_file, 
            html_ast: Some(ast), 
            mentions: Some(mentions),
            issues: Some(issues),
        }
    }
}

pub struct OutputPages(Vec<OutputPage>);

impl OutputPages {
    fn new(path: &Path, mut scraper: &mut Scraper) -> Self {
        let mut parser = MarkdownIt::new();

        // Standard markdown parsing rules
        markdown_it::plugins::cmark::add(&mut parser);

        // Custom markdown parsing rules
        markdown::timecode::add(&mut parser);
        markdown::speaker::add(&mut parser);
        markdown::github::add(&mut parser); // must apply before sidenote rules
        markdown::sidenote::add(&mut parser);
        markdown::mention::add(&mut parser);

        let mut input_file_results = vec![];
        for entry in fs::read_dir(path.clone()).unwrap() {

            let Ok(input_file) = InputFile::new(entry.unwrap().path()) else {
                continue;
            };

            input_file_results.push(
                OutputPage::new(input_file, &mut parser, &mut scraper)
            );
        }

        OutputPages(input_file_results.into())
    }
}


#[derive(Debug)]
pub struct InputFileBeingParsed(InputFile);
impl MarkdownItExt for InputFileBeingParsed {}


fn render(path: &PathBuf, content: DynRender) {
    let dir = path.parent()
        .expect(format!("Failed to determine directory for {:?}", path).as_str());
    std::fs::create_dir_all(dir)
        .expect(format!("Failed to create directories for {:?}", dir).as_str());
    std::fs::write(&path, content.to_string())
        .expect(format!("Failed to write {:?}", path).as_str());
    println!("  Wrote {:?}", path);
}

#[allow(dead_code)]
fn copy_file(input_path: &PathBuf, output_path: &PathBuf) {
    std::fs::copy(input_path, output_path)
        .expect(format!("Failed to copy {:?} to {:?}", input_path, output_path).as_str());
    println!("  Copied file {:?} to {:?}", input_path, output_path); 
}

fn copy_dir(input_path: &PathBuf, output_path: &PathBuf) {
    fs_extra::dir::copy(
        input_path, 
        output_path, 
        &CopyOptions::new()
            .overwrite(true)
            .content_only(true)
    ).expect(format!("Failed to copy directory {:?} to {:?}", input_path, output_path).as_str());
    println!("  Copied directory {:?} to {:?}", input_path, output_path); 
}


#[derive(Debug)]
struct GlobalScraper<'a>(Scraper<'a>);
impl MarkdownItExt for GlobalScraper<'static> {}


markup::define! {
    Base<'a>(title: Option<&'a str>, content: DynRender<'a>) {
        @markup::doctype()
        html {
            head {
                title { 
                    @if let Some(title) = title {
                        @title " - Ray Peat Rodeo"
                    } else {
                        "Ray Peat Rodeo"
                    }
                }
                meta[charset = "UTF-8"] {}
                meta[name = "viewport", 
                    content = "width=device-width, initial-scale=1.0"] {}
                link[rel = "stylesheet", href = "/assets/style.css"] {}
                // script[src="https://unpkg.com/htmx.org@1.9.6"] {}
            }
            body {
                div #"top-bar" {
                    a [href="/"] { "Ray Peat Rodeo" }

                    @PageFind {}
                }

                @content

                footer {
                    p {
                        "Ray Peat Rodeo is an "
                        a[href = GITHUB_LINK] { "open source" }
                        " website written in Rust. Last updated "
                        @chrono::offset::Local::now().format("%Y-%m-%d").to_string()
                        "."
                    }
                }
            }
        }
    }

    PageFind() {
        link[href = "/pagefind/pagefind-ui.css", rel = "stylesheet"] {}
        script[src = "/pagefind/pagefind-ui.js", type = "text/javascript"] {}
        div #search {}
        script { @markup::raw(r#"
            window.addEventListener('DOMContentLoaded', (event) => {
                new PagefindUI({
                    element: '#search',
                    showSubResults: true,
                    showImages: false,
                    translations: {
                        placeholder: 'Search Ray Peat Rodeo',
                    }
                });
            });
        "#) }
    }

    Sidenote<'a>(content: DynRender<'a>) {
        span .sidenote."sidenote-standalone" {
            @content
        }
    }

    Homepage<'a>(output_pages: &'a Vec<OutputPage>) {
        @Base {
            title: Some("Ray Peat Rodeo"),
            content: markup::new! {
                article {
                    section {
                        @Sidenote {
                            content: markup::new! {
                                div.links {
                                    a."github-link" [href = GITHUB_LINK] { 
                                        img [
                                            src = "/assets/images/github-mark.svg",
                                            title = "Visit project on GitHub",
                                        ] {}
                                    } 
                                }

                                r#"Ray Peat Rodeo offers accurate, referenced 
                                transcripts of Ray Peat interviews that can be 
                                easily searched or surveyed."#
                                br {}
                                r#"Transcripts are accessibly written in markdown, 
                                and leverage a pleasant custom syntax to describe 
                                who's speaking, mark referenced works and authors, 
                                insert sidenotes, and even to add callouts to 
                                GitHub issues discussing textual improvements."#
                                br {}
                                r#"Project longevity, flexibility and simplicity 
                                is undergirded by a beskpoke engine written in 
                                Rust. Ease of development and deployment are 
                                guaranteed by the excellent nix package manager. 
                                The project is maintained, discussed, and deployed
                                via GitHub."#
                            }
                        }

                        div #documents {
                            @for output_page in output_pages.iter() {
                                div .document {
                                    div ."document-header" {
                                        div ."document-hud" {
                                            span ."document-date" { @output_page.input_file.date.to_string() }
                                            a ."document-series" [
                                                href = format!("/series/{}", output_page.input_file.frontmatter.source.series.to_lowercase().replace(" ", "-"))
                                            ] {
                                                @output_page.input_file.frontmatter.source.series
                                            }
                                        }

                                        @if output_page.input_file.todo {
                                            a .todo."document-title" [href = output_page.input_file.slug.clone()] {
                                                @output_page.input_file.frontmatter.source.title
                                            }
                                        } else {
                                            a ."document-title" [href = output_page.input_file.slug.clone()] {
                                                @output_page.input_file.frontmatter.source.title
                                            }
                                        }

                                    }

                                    @if let Some(mentions) = &output_page.mentions {
                                        div ."document-body" {
                                            @for mention in mentions.iter().filter(|m| m.position == 1) {
                                                a .mention.{ mention.kind() } [
                                                    href = mention.slug(), 
                                                    title = mention.display_text(),
                                                ] {
                                                    @mention.display_text()
                                                }

                                                " "
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    AuthorPage(author: Author, mentions: Vec<Mention>) {
        @Base {
            title: Some(author.display_text().as_str()),
            content: markup::new! {
                article {
                    section {
                        h1 { @author.display_text() }

                        p {
                            "This author is mentioned on the following pages."
                        }

                        ul {
                            @for mention in mentions {
                                li {
                                    a[href = format!("/{}", mention.slug())] {
                                        @mention.input_file.frontmatter.source.title
                                    }
                                }
                            }
                        }
                    }
                }
            },
        }
    }

    SeriesPage<'a>(series: &'a String, input_files: &'a Vec<&'a InputFile>) {
        @Base {
            title: Some(series),
            content: markup::new! {
                article {
                    section {
                        h1 { @series }

                        p {
                            @input_files.len() " entries for " @series "."
                            br {}
                        }

                        ul {
                            @for input_file in input_files.iter() {
                                li {
                                    @if input_file.todo {
                                        a.todo[href = format!("/{}", input_file.slug)] {
                                            @input_file.frontmatter.source.title.clone()
                                        }
                                    } else {
                                        a[href = format!("/{}", input_file.slug)] {
                                            @input_file.frontmatter.source.title.clone()
                                        }

                                    }

                                    " ("
                                    a[href = input_file.frontmatter.source.url.clone()] {
                                        "source"
                                    }
                                    ")"
                                }
                            }
                        }
                    }
                }
            },
        }
    }
    Page<'a>(output_page: &'a OutputPage) {
        @Base {
            title: Some(output_page.input_file.frontmatter.source.title.as_str()),
            content: markup::new! {
                article .interview ["data-pagefind-body"] {
                    header {
                        div.hud ["data-pagefind-ignore"] {
                            span.date { @output_page.input_file.date.to_string() }
                            a.series [
                                href = format!("/series/{}", output_page.input_file.frontmatter.source.series.to_lowercase().replace(" ", "-"))
                            ] {
                                @output_page.input_file.frontmatter.source.series
                            }
                        }

                        h1.title { @output_page.input_file.frontmatter.source.title }

                        @if let Some(transcription) = output_page.input_file.frontmatter.transcription.clone()  {
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
                            a."orig-content"[
                                href = output_page.input_file.frontmatter.source.url.clone()
                            ] { "View Source" }
                            a."edit"[
                                href = format!("{GITHUB_LINK}/edit/main/{}", output_page.input_file.path)
                            ] { "Edit on GitHub" }
                        }

                    }

                    main {
                        @if let Some(html_ast) = &output_page.html_ast {
                            @markup::raw(html_ast.render())
                        } else if output_page.input_file.todo {
                            section {
                                p {
                                    br {}
                                    "A transcript for this interview doesn't yet exist." 

                                    @if let Some(t) = &output_page.input_file.frontmatter.transcription {
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
