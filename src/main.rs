mod markdown;
mod scraper;

use std::{
    fs,
    path::{Path, PathBuf}, collections::HashMap,
};
use clap::Parser;
use fs_extra::dir::CopyOptions;
use markdown::{
    Frontmatter,
    mention::{MentionDeclaration, Mention, Author}, github::{GitHubIssueDeclaration, GitHubIssue}
};
use markdown_it::{MarkdownIt, parser::extset::MarkdownItExt, Node};
use markup::DynRender;
use extract_frontmatter::{Extractor, config::Splitter::EnclosingLines};
use scraper::Scraper;

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
        let mut scraper = Scraper::new_scraper(cache_path.to_path_buf());
        parse_input_files(input.to_path_buf(), &mut scraper);
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

    let input_file_results = parse_input_files(input.to_path_buf(), &mut Scraper::new_fulfiller(cache_path.to_path_buf()));

    for (input_file, ast, _mentions, _issues) in input_file_results.iter() {
        render(&output.join(format!("{}/index.html", input_file.slug)), markup::new! {
            @Page { input_file: &input_file, html: ast.render() }
        });
    }

    {
        let mut authors: HashMap<Author, Vec<Mention>> = HashMap::new();
        for (_, _, mentions, _) in input_file_results.iter() {
            for mention in mentions {
                authors.entry(mention.author.clone()).or_insert(vec![]).push(mention.clone());
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
        for (input_file, _, _, _) in input_file_results.iter() {
            if let Some(series) = &input_file.frontmatter.source.series {
                series_map.entry(series).or_insert(vec![]).push(input_file);
            }
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
        @Homepage { 
            input_files: input_file_results
                .iter()
                .map(|(input_file, _ast, mentions, _issues)| {
                    (
                        input_file.clone(), 
                        mentions.iter()
                        .filter(|mention| mention.position == 1)
                        .collect()
                    )
                })
                .collect() 
        }
    });

    copy_dir(&PathBuf::from("./src/assets"), &output.join("assets"));

    println!("\nDone");

    Ok(())
}

/// An input markdownfile from ./content
#[derive(Debug, Hash, PartialEq, Eq, Clone)]
pub struct InputFile {
    date: String,
    path: String,
    slug: String,
    frontmatter: Frontmatter,
    markdown: String,
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



fn parse_input_files(input: PathBuf, scraper: &mut Scraper) -> Vec<(InputFile, Node, Vec<Mention>, Vec<GitHubIssue>)> {
    let parser = &mut MarkdownIt::new();

    // Standard markdown parsing rules
    markdown_it::plugins::cmark::add(parser);

    // Custom markdown parsing rules
    markdown::timecode::add(parser);
    markdown::speaker::add(parser);
    markdown::github::add(parser); // must apply before sidenote rules
    markdown::sidenote::add(parser);
    markdown::mention::add(parser);

    let mut input_file_results = vec![];
    for entry in fs::read_dir(input.clone()).unwrap() {
        let entry = entry.unwrap();
        let path = entry.path();

        let extension = path.extension();
        if path.is_dir() || extension.is_none() || extension.is_some_and(|ext| ext != "md") {
            continue;
        }

        let text = std::fs::read_to_string(&path).unwrap();

        let (raw_frontmatter, markdown) = Extractor::new(EnclosingLines("---"))
            .extract(text.as_str());

        let (date, slug) = path.file_stem().unwrap().to_str().unwrap().split_at(11);

        let input_file = InputFile {
            date: date.split_at(10).0.into(),
            path: path.strip_prefix(input.as_path()).unwrap().to_str().unwrap().into(),
            markdown: markdown.into(),
            slug: slug.into(),
            frontmatter: match serde_yaml::from_str(&raw_frontmatter) {
                Ok(f) => f,
                Err(e) => panic!("Invalid YAML frontmatter in {:?}\n{e}", path),
            },
        };

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

        input_file_results.push((input_file, ast, mentions, issues));
    }

    input_file_results
}


#[derive(Debug)]
struct GlobalScraper(Scraper);
impl MarkdownItExt for GlobalScraper {}


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
        link[href = "/_pagefind/pagefind-ui.css", rel = "stylesheet"] {}
        script[src = "/_pagefind/pagefind-ui.js", type = "text/javascript"] {}
        div #search {}
        script { r#"
            window.addEventListener('DOMContentLoaded', (event) => {
                new PagefindUI({
                    element: '#search',
                    translations: {
                        placeholder: 'Search Ray Peat Rodeo'
                    }
                });
            });
        "# }
    }

    Sidenote<'a>(content: DynRender<'a>) {
        span .sidenote."sidenote-standalone" {
            @content
        }
    }

    Homepage<'a>(input_files: Vec<(InputFile, Vec<&'a Mention>)>) {
        @Base {
            title: Some("Ray Peat Rodeo"),
            content: markup::new! {
                article {
                    section { @PageFind {} }

                    section {
                        @Sidenote {
                            content: markup::new! {
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
                                via "#
                                a [href = GITHUB_LINK] { "GitHub" } "."
                            }
                        }

                        div #documents {
                            @for (input_file, mentions) in input_files {
                                div .document {
                                    div ."document-header" {
                                        div ."document-hud" {
                                            span ."document-date" { @input_file.date }
                                            @if let Some(series) = &input_file.frontmatter.source.series {
                                                a ."document-series" [
                                                    href = format!("/series/{}", series.to_lowercase().replace(" ", "-"))
                                                ] {
                                                    @series
                                                }
                                            }
                                        }

                                        a ."document-title" [href = input_file.slug.clone()] {
                                            @input_file.frontmatter.source.title
                                        }

                                    }

                                    div ."document-body" {
                                        @for mention in mentions.iter() {
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
                                    a[href = format!("/{}", input_file.slug)] {
                                        @input_file.frontmatter.source.title.clone()
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
    Page<'a>(input_file: &'a InputFile, html: String) {
        @Base {
            title: Some(input_file.frontmatter.source.title.as_str()),
            content: markup::new! {
                article .interview {
                    header {
                        h1.title { @input_file.frontmatter.source.title }

                        div.hud {
                            @if let Some(series) = &input_file.frontmatter.source.series {
                                a.series [
                                    href = format!("/series/{}", series.to_lowercase().replace(" ", "-"))
                                ] {
                                    @series
                                }
                            }
                        }

                        span .sidenote."sidenote-meta" {
                            { input_file.frontmatter.source.series.clone().unwrap_or("Someone".into()) }
                            a [href = input_file.frontmatter.source.url.clone(), target = "_bank"] {
                                " originally published "
                            }

                            " this interview on " @input_file.date "."

                            @if let Some(transcription) = input_file.frontmatter.transcription.clone()  {
                                @if let Some(author) = transcription.author.clone() {
                                    @if let Some(url) = transcription.url.clone() {
                                        " Thank you to " @author " who "
                                        a[href = url] {
                                            "published"
                                        }
                                        " this transcript "
                                    } else {
                                        " Transcribed by " @author ", "
                                    }

                                    @ if let Some(date) = transcription.date.clone() {
                                        @date "."
                                    }

                                }
                            } 

                            " "

                            a[href = format!("{GITHUB_LINK}/edit/main/content/{}", input_file.path), target="_blank"] {
                                "Edit"
                            }

                            " this page on GitHub."
                        }
                    }

                    main ["data-pagefind-body"] {
                        @markup::raw(html)
                    }
                }
            },
        }
    }
}
