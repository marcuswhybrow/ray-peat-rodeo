mod markdown;

use std::{fs, include_str, path::Path, collections::BTreeMap, collections::HashMap};
use clap::Parser;
use serde::{Serialize, Deserialize};
use extract_frontmatter::{Extractor, config::Splitter::EnclosingLines};
use markdown_it::parser::extset::MarkdownItExt;
use crate::markdown::mention::{Mention, Mentionable, Person, Link, Book, Paper};


#[derive(Parser, Debug)]
#[command(name = "Ray Peat Rodeo Engine")]
#[command(author = "Marcus Whybrow <marcus@whybrow.ustatusListk>")]
#[command(about = "Builds Ray Peat Rodeo into HTML from source")]
#[command(long_about = None)]
struct Args {
    /// The input path containing markdown content
    #[arg(short, long, default_value_t = String::from("./content"))]
    input: String,

    /// The output path in which Ray Peat Rodeo should build HTML
    #[arg(short, long, default_value_t = String::from("./build"))]
    output: String,

    /// Whether files and directories inside of OUT_PATH should be deleted
    /// before building
    #[arg(short, long, default_value_t = false)]
    clean: bool,
}

struct Renderer {
    tera: tera::Tera,
    global_context: tera::Context,
    out_path: std::path::PathBuf,
}

impl Renderer {
    fn new(out_path: std::path::PathBuf) -> Self {
        Renderer {
            tera: {
                let mut tera = tera::Tera::default();
                tera.add_raw_templates(vec![
                   ("base.html", include_str!("./templates/base.html")),
                   ("page.html", include_str!("./templates/page.html")),
                   ("index.html", include_str!("./templates/index.html")),
                   ("style.css", include_str!("./templates/style.css")),
                ]).expect("Unable to load template");
                tera
            },
            global_context: {
                let mut gcx = tera::Context::new();
                gcx.insert("global_project_link", "https://github.com/marcuswhybrow/ray-peat-rodeo");
                gcx.insert("global_contact_link", "https://raypeat.rodeo/contact");
                gcx
            },
            out_path,
        }
    }

    fn render_with_context(&self, template: &str, path: &str, cx: tera::Context) {
        let final_path = self.out_path.join(path);
        std::fs::create_dir_all(&final_path.parent().unwrap()).unwrap();
        let mut new_cx = self.global_context.clone();
        new_cx.extend(cx);
        self.tera.render_to(
            template,
            &new_cx,
            std::fs::File::create(&final_path).unwrap(),
        ).unwrap();
        println!("Wrote {:?}", final_path.to_str().unwrap());

    }

    fn render_from(&self, renderable: Renderable) {
        match renderable {
            Renderable::Document(document) => {
                let mut cx = tera::Context::new();
                cx.insert("title", &document.parsed_markdown_file.markdown_file.frontmatter.title);
                cx.insert("contents", &document.parsed_markdown_file.ast.render());
                cx.insert("mentions", &document.mentions);
                cx.extend(self.global_context.clone());

                let path = format!("{}/index.html", document.parsed_markdown_file.markdown_file.slug);

                self.render_with_context("page.html", path.as_str(), cx);
            },
            Renderable::TemplateAsIs(template) => {
                self.render_with_context(template.as_str(), template.as_str(), tera::Context::new());
            },
            Renderable::TemplateAsIsWithContext(template, cx) => {
                self.render_with_context(template.as_str(), template.as_str(), cx);
            },
        }
    }
}

#[derive(Debug)]
struct Document {
    parsed_markdown_file: ParsedMarkdownFile,
    mentions: Vec<Mention>,
}

impl PartialEq for Document {
    fn eq(&self, other: &Self) -> bool {
        self.parsed_markdown_file.markdown_file.path == other.parsed_markdown_file.markdown_file.path
    }
}

impl Eq for Document {}

impl PartialOrd for Document {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        self.parsed_markdown_file.markdown_file.path.partial_cmp(&other.parsed_markdown_file.markdown_file.path)
    }
}

impl Ord for Document {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        self.parsed_markdown_file.markdown_file.path.cmp(&other.parsed_markdown_file.markdown_file.path)
    }
}

enum Renderable<'a> {
    Document(&'a Document),
    TemplateAsIs(String),
    TemplateAsIsWithContext(String, tera::Context),
}

#[derive(Debug, Clone, Hash, Eq, PartialEq, PartialOrd, Ord, serde::Serialize)]
struct MarkdownFile {
    path: std::path::PathBuf,
    slug: String,
    frontmatter: Frontmatter,
    markdown: String,
}

impl MarkdownItExt for MarkdownFile {}

impl MarkdownFile {
    fn new(path: std::path::PathBuf) -> Self {
        let text = fs::read_to_string(&path)
            .expect(format!("Could not read {:?}", &path).as_str());

        let (frontmatter_text, markdown) = Extractor::new(EnclosingLines("---")).extract(text.as_str());

        MarkdownFile {
            path: path.clone(),
            slug: path.file_stem().unwrap().to_str().unwrap().split_at(11).1.to_string(),
            markdown: markdown.to_string(),
            frontmatter: serde_yaml::from_str(&frontmatter_text)
                .expect(format!("Invalid YAML frontmatter in {:?}", path).as_str()),
        }
    }
}

#[derive(Debug)]
struct ParsedMarkdownFile {
    markdown_file: MarkdownFile,
    ast: markdown_it::Node,
}

impl ParsedMarkdownFile {
    fn new(markdown_file: MarkdownFile) -> Self {
        let markdown_parser = &mut markdown_it::MarkdownIt::new();
        markdown_it::plugins::cmark::add(markdown_parser);
        markdown::timecode::add(markdown_parser);
        markdown::speaker::add(markdown_parser);
        markdown::sidenote::add(markdown_parser);
        markdown::mention::add(markdown_parser);

        markdown_parser.ext.insert(markdown_file.clone());

        let markdown = markdown_file.markdown.clone();

        ParsedMarkdownFile {
            markdown_file,
            ast: markdown_parser.parse(markdown.as_str()),
        }
    }

    fn get_mentions(&self) -> Vec<Mention> {
        let mut mentions: Vec<Mention> = vec![];
        self.ast.walk(|node, _depth| {
            if let Some(mention) = node.cast::<Mention>() {
                mentions.push(mention.clone());
            }
        });
        mentions
    }
}

#[derive(Serialize, Deserialize, Debug, Clone, Hash, PartialEq, Eq, PartialOrd, Ord)]
struct Transcription {
    source: Option<String>,
    author: Option<String>,
    date: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone, Hash, PartialEq, PartialOrd, Ord, Eq)]
struct Frontmatter {
    title: String,
    series: Option<String>,
    speakers: BTreeMap<String, String>, 
    source: String,
    source_duration: Option<String>,
    transcription: Option<Transcription>,
}

fn main() {
    let args = Args::parse();
    let input = &Path::new(&args.input).canonicalize().expect("Input path not found");
    let output = &Path::new(&args.output).canonicalize().expect("Output path not found");

    println!("Building Ray Peat Rodeo");
    println!("Input: {:?}", input);
    println!("Output: {:?}", output);

    if !output.exists() {
        println!("Creating directory");
        fs::create_dir(output).unwrap();
    } else {
        if args.clean {
            println!("Clean option enabled. \
                Deleting files and directories inside {:?}",
                output);

            for entry in fs::read_dir(output).unwrap() {
                let entry = entry.unwrap();
                let path = entry.path();

                if entry.file_type().unwrap().is_dir() {
                    fs::remove_dir_all(path).unwrap();
                } else {
                    fs::remove_file(path).unwrap();
                }
            }
        }
    }

    let renderer = Renderer::new(output.clone());

    let documents = fs::read_dir(input)
        .unwrap()
        .map(|entry| {
            let parsed_markdown_file = ParsedMarkdownFile::new(MarkdownFile::new(entry.unwrap().path()));
            let mentions = parsed_markdown_file.get_mentions();
            let document = Document { parsed_markdown_file, mentions };
            renderer.render_from(Renderable::Document(&document));
            document
        });
    
    let mut people: HashMap<Person, HashMap<MarkdownFile, Vec<Mention>>> = HashMap::new();
    let mut books: HashMap<Book, HashMap<MarkdownFile, Vec<Mention>>> = HashMap::new();
    let mut links: HashMap<Link, HashMap<MarkdownFile, Vec<Mention>>> = HashMap::new();
    let mut papers: HashMap<Paper, HashMap<MarkdownFile, Vec<Mention>>> = HashMap::new();
    for document in documents {
        for mention in document.mentions {
            match mention.clone() {
                Mention::Normal { mentionable, occurance: _, fragment: _ } => {
                    let markdown_file = document.parsed_markdown_file.markdown_file.clone();
                    match mentionable {
                        Mentionable::Person(person) => people
                            .entry(person).or_insert(HashMap::new())
                            .entry(markdown_file).or_insert(vec![])
                            .push(mention),
                        Mentionable::Book(book) => books
                            .entry(book).or_insert(HashMap::new())
                            .entry(markdown_file).or_insert(vec![])
                            .push(mention),
                        Mentionable::Link(link) => links
                            .entry(link).or_insert(HashMap::new())
                            .entry(markdown_file).or_insert(vec![])
                            .push(mention),
                        Mentionable::Paper(paper) => papers
                            .entry(paper).or_insert(HashMap::new())
                            .entry(markdown_file).or_insert(vec![])
                            .push(mention),
                    }
                },
                _ => (),
            }
        }
    }

    let mut cx = tera::Context::new();
    cx.insert("people", &people);
    cx.insert("books", &books);
    cx.insert("links", &links);
    cx.insert("papers", &papers);

    renderer.render_from(Renderable::TemplateAsIsWithContext("index.html".into(), cx));
    renderer.render_from(Renderable::TemplateAsIs("style.css".into()));

    println!("Done");
}
