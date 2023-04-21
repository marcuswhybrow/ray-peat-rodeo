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

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();
    let out_path_arg: String = match args.out_path {
        Some(path) => path,
        None => String::from("build"),
    };
    let out_path = Path::new(out_path_arg.as_str());
    let canonical_out_path = out_path.canonicalize()?;

    println!("Building Ray Peat Rodeo in {}", canonical_out_path.display());

    if !out_path.exists() {
        println!("Creating directory");
        fs::create_dir(out_path)?;
    } else {
        if args.clean {
            println!("Cleaning directory");

            for entry in fs::read_dir(out_path)? {
                let entry = entry?;
                let path = entry.path();

                if entry.file_type()?.is_dir() {
                    fs::remove_dir_all(path)?;
                } else {
                    fs::remove_file(path)?;
                }
            }
        }
    }

    println!("Writing index.html");
    fs::write(out_path.join("index.html"), r#"
        <html>
            <head></head>
            <body>
                <h1>Ray Peat Rodeo</h1>
            </body>
        </html>
    "#)?;

    println!("Done");

    Ok(())
}
