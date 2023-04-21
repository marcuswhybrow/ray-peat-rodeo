#[macro_use]
extern crate lazy_static;

use std::fs;
use std::path::Path;
use clap::Parser;

#[derive(Parser, Debug)]
#[command(name = "Ray Peat Rodeo Engine")]
#[command(author = "Marcus Whybrow <marcus@whybrow.uk>")]
#[command(about = "Builds Ray Peat Rodeo into HTML from source")]
#[command(long_about = None)]
struct Args {
    /// The directory in which Ray Peat Rodeo should build HTML
    out_path: Option<String>,

    /// Whether files and directories inside of OUT_PATH should be deleted
    /// before building
    #[arg(short, long, default_value_t = false)]
    clean: bool,
}

fn main() {
    let args = Args::parse();
    let out_path_arg: String = match args.out_path {
        Some(path) => path,
        None => String::from("build"),
    };
    let out_path = Path::new(out_path_arg.as_str());
    let canonical_out_path = out_path.canonicalize().unwrap();

    println!("Building Ray Peat Rodeo in {}", canonical_out_path.display());

    if !out_path.exists() {
        println!("Creating directory");
        fs::create_dir(out_path).unwrap();
    } else {
        if args.clean {
            println!("Clean option enabled. \
                Deleting files and directories inside {:?}",
                canonical_out_path);

            for entry in fs::read_dir(out_path).unwrap() {
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

    lazy_static! {
        pub static ref TEMPLATES: tera::Tera = {
            match tera::Tera::new("templates/**/*") {
                Ok(t) => t,
                Err(e) => {
                    println!("Parsing error(s): {}", e);
                    std::process::exit(1);
                }
            }
        };
    }

    fn write_from_template(
        template_name: &str,
        cx: &tera::Context,
        out_path: &std::path::PathBuf
    ) {
        TEMPLATES.render_to(
            template_name,
            &cx,
            fs::File::create(out_path).unwrap()
        ).unwrap();
        println!("Wrote {:?}", out_path.canonicalize().unwrap());
    }

    let mut cx = tera::Context::new();

    cx.insert("global_project_link", "https://github.com/marcuswhybrow/ray-peat-rodeo");
    cx.insert("global_contact_link", "https://raypeat.rodeo/contact");

    write_from_template("index.html", &cx, &out_path.join("index.html"));
    write_from_template("style.css", &cx, &out_path.join("style.css"));

    println!("Done");
}
