mod markdown;
mod stash;
mod content;

use std::{
    fs,
    path::{Path, PathBuf},
};
use clap::Parser;
use fs_extra::dir::{CopyOptions, remove};
use itertools::Itertools;
use markdown::ProjectParser;
use markdown_it::MarkdownIt;
use markup::DynRender;
use stash::Stash;
use content::{Content, ContentFile};

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

    #[arg(long, default_value_t = String::from("./src/stash.yml"))]
    stash_path: String,

    #[arg(long, default_value_t = false)]
    wipe_stash: bool,

    #[arg(long, default_value_t = false)]
    update_stash: bool,

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
    let stash_path = Path::new(&args.stash_path);

    println!("Ray Peat Rodeo");
    println!("  Input:  {:?}", input);
    println!("  Output: {:?}", output);
    println!("  Cache:  {:?}", stash_path);

    if args.clean {
        remove(output).unwrap();
        println!("  Cleaned. Removed {:?}", output);
    }

    let mut parser = MarkdownIt::new_project_parser();

    let mut content_files: Vec<ContentFile> = fs::read_dir(&input).unwrap().into_iter()
        .chain(fs::read_dir(&input.join("todo")).unwrap().into_iter())
        .filter_map(|e| ContentFile::try_from(e.unwrap()).ok())
        .collect();
    content_files.sort_by_key(|content_file| content_file.filename_date.as_date);
    content_files.reverse();

    if args.wipe_stash {
        fs::remove_file(stash_path).unwrap();
        println!("  Wiping stash. Removed {:?}", stash_path);

        println!("  Regenerating stash...");
        Stash::from(stash_path)
            .analyse(content_files.clone())
            .write(stash_path).await;
        println!("  Built stash at {:?}", stash_path);

    } else if args.update_stash {
        Stash::from(stash_path)
            .analyse(content_files.clone())
            .write(stash_path).await;
        println!("  Stash updating enabled.");
    }

    println!("\nWriting Files");

    let mut done: Vec<ContentFile> = vec![];
    let mut todo: Vec<ContentFile> = vec![];

    for content_file in content_files {
        if content_file.is_todo() {
            todo.push(content_file);
        } else {
            done.push(content_file);
        }
    }

    let mut stash = Stash::from(stash_path);

    let all_content: Vec<Content> = done.into_iter()
        .chain(todo.into_iter())
        .map(|content_file| {
            let content = content_file.as_content(&mut parser, &mut stash);
            content.write(output);
            content
        })
        .collect();

    for (mentionable, mentions) in all_content.iter()
        .map(|c| &c.mentions)
        .flatten()
        .into_group_map_by(|m| &m.mentionable) 
    {
        mentionable.write(output, mentions);
    }

    for (series, content) in all_content.iter()
        .into_group_map_by(|c| &c.file.frontmatter.source.series)
    {
        write(&output.join(format!("series/{}/index.html", series
            .to_lowercase()
            .replace(" ", "-")
        )), markup::new! {
            @SeriesPage { series, content: &content }
        });
    }

    write(&output.join("index.html"), markup::new! {
        @Homepage { all_content: &all_content }
    });

    copy_dir(&PathBuf::from("./src/assets"), &output.join("assets"));

    println!("\nDone");

    Ok(())
}


fn write(path: &PathBuf, content: DynRender) {
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


markup::define! {
    BasePage<'a>(title: Option<&'a str>, content: DynRender<'a>) {
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
                script[src="https://unpkg.com/htmx.org@1.9.6"] {}
            }
            body {
                div #"top-bar" {
                    a [href="/"] { "Ray Peat Rodeo" }

                    @PageFind {}
                }

                @content

                footer {
                    iframe [
                        src="https://github.com/sponsors/marcuswhybrow/button",
                        title="Sponsor marcuswhybrow",
                        height="32",
                        width="114",
                        style="border: 0; border-radius: 6px;",
                    ] {}

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

    Homepage<'a>(all_content: &'a Vec<Content>) {
        @BasePage {
            title: Some("Ray Peat Rodeo"),
            content: markup::new! {
                article.homepage {
                    section {
                        @Sidenote {
                            content: markup::new! {
                                div."homepage-hud" {
                                    a."github-project" [href = GITHUB_LINK] { 
                                        img [
                                            src = "/assets/images/github-mark.svg",
                                            title = "Visit project on GitHub",
                                        ] {}
                                    } 
                                    span."github-sponsor" {
                                        iframe [
                                            src="https://github.com/sponsors/marcuswhybrow/button",
                                            title="Sponsor marcuswhybrow",
                                            height="32",
                                            width="114",
                                            style="border: 0; border-radius: 6px;",
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

                        div #content {
                            @for content in all_content.iter() {
                                div.content {
                                    div.header {
                                        div.hud {
                                            span.date { @content.file.filename_date.to_string() }
                                            a.series [
                                                href = format!("/series/{}", content.file.frontmatter.source.series.to_lowercase().replace(" ", "-"))
                                            ] {
                                                @content.file.frontmatter.source.series
                                            }
                                        }

                                        @if content.file.is_todo() {
                                            a.todo.title [href = content.file.title_slug.clone()] {
                                                @content.file.frontmatter.source.title
                                            }
                                        } else {
                                            a.title [href = content.file.title_slug.clone()] {
                                                @content.file.frontmatter.source.title
                                            }
                                        }

                                    }

                                    @if !content.mentions.is_empty() {
                                        div.body {
                                            @for mention in content.mentions.iter().filter(|m| m.position == 1) {
                                                span.mention [
                                                    "hx-trigger" = "mouseenter",
                                                    "hx-target" = "find .popup-card",
                                                    "hx-get" = mention.mentionable_permalink(),
                                                    "hx-select" = ".popup-select",
                                                    "hx-swap" = "innerHTML",
                                                ]{
                                                    a [
                                                        href = mention.permalink(),
                                                    ] {
                                                        @mention.ultimate().default_display_text
                                                    }

                                                    span."popup-card" {}
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

    SeriesPage<'a>(series: &'a String, content: &'a Vec<&'a Content>) {
        @BasePage {
            title: Some(series),
            content: markup::new! {
                article {
                    section {
                        h1 { @series }

                        p {
                            @content.len() " entries for " @series "."
                            br {}
                        }

                        ul {
                            @for content in content.iter() {
                                li {
                                    @if content.file.is_todo() {
                                        a.todo[href = format!("/{}", content.file.title_slug)] {
                                            @content.file.frontmatter.source.title.clone()
                                        }
                                    } else {
                                        a[href = format!("/{}", content.file.title_slug)] {
                                            @content.file.frontmatter.source.title.clone()
                                        }

                                    }

                                    " ("
                                    a[href = content.file.frontmatter.source.url.clone()] {
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

}
