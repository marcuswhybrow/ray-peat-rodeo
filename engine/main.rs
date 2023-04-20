use std::fs;
use std::path::Path;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let out = Path::new("build");

    if out.exists() {
        for entry in fs::read_dir(out)? {
            let entry = entry?;
            let path = entry.path();

            if entry.file_type()?.is_dir() {
                fs::remove_dir_all(path)?;
            } else {
                fs::remove_file(path)?;
            }
        }
    } else {
        fs::create_dir("build")?;
    }


    fs::write(out.join("index.html"), r#"
        <html>
            <head></head>
            <body>
                <h1>Ray Peat Rodeo 1</h1>
            </body>
        </html>
    "#)?;

    println!("Ray Peat Rodeo");

    Ok(())
}
