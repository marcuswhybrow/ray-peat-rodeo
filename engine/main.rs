use std::fs;
use std::path::Path;
use clap::Parser;

#[derive(Parser, Debug)]
#[command(name = "Ray Peat Rodeo Engine")]
#[command(author = "Marcus Whybrow <marcus@whybrow.uk>")]
#[command(about = "Builds Ray Peat Rodeo into HTML from source")]
#[command(long_about = None)]
struct Args {
    /// The input path containing markdown content
    #[arg(short, long, default_value_t = std::string::String::from("content"))]
    input: String,

    /// The output path in which Ray Peat Rodeo should build HTML
    #[arg(short, long, default_value_t = std::string::String::from("build"))]
    output: String,

    /// Whether files and directories inside of OUT_PATH should be deleted
    /// before building
    #[arg(short, long, default_value_t = false)]
    clean: bool,
}

fn main() {

    // CLI Arguments and Options

    let args = Args::parse();
    let output = &Path::new(&args.output).canonicalize().unwrap();

    println!("Building Ray Peat Rodeo in {:?}", output);

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

    // Templating

    let tera = tera::Tera::new("templates/**/*").unwrap();

    let mut gcx = tera::Context::new();
    gcx.insert("global_project_link", "https://github.com/marcuswhybrow/ray-peat-rodeo");
    gcx.insert("global_contact_link", "https://raypeat.rodeo/contact");

    let render = |template, context: &tera::Context, path: &str| {
        let final_path = output.join(path);
        std::fs::create_dir_all(&final_path.parent().unwrap()).unwrap();
        tera.render_to(
            template,
            &context,
            std::fs::File::create(&final_path).unwrap(),
        ).unwrap();
        println!("Wrote {:?}", final_path);
    };

    // Render Specific Pages

    render("index.html", &gcx, "index.html");
    render("style.css", &gcx, "style.css");

    // Render Content

    let input = &Path::new(&args.input).canonicalize().unwrap();

    use extract_frontmatter::{Extractor, config::Splitter::EnclosingLines};
    let frontmatter_extractor = Extractor::new(EnclosingLines("---"));

    let markdown_parser = &mut markdown_it::MarkdownIt::new();
    markdown_it::plugins::cmark::add(markdown_parser);

    for entry in fs::read_dir(input).unwrap() {
        let entry = entry.unwrap();
        let path = entry.path();
        let text = std::fs::read_to_string(&path).unwrap();

        let (_, data) = frontmatter_extractor.extract(text.as_str());

        let html = markdown_parser.parse(data).render();
        let (_, slug) = path.file_stem().unwrap().to_str().unwrap().split_at(11);
        let out_name = &format!("{}/index.html", &slug);

        let mut cx = tera::Context::new();
        cx.insert("contents", html.as_str());
        cx.extend(gcx.clone());

        render("page.html", &cx, out_name);
    }

    println!("Done");
}
