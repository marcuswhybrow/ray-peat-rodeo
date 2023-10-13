use std::{collections::HashMap, path::Path, fmt::Debug};
use markdown_it::{parser::extset::MarkdownItExt, MarkdownIt};
use tokio::task::JoinSet;

use crate::{content::ContentFile, markdown::ProjectParser};

type ResponseHandler = Box<dyn FnOnce(String, String) -> String + Send + Sync + 'static>;

/// An HTTP request to be made in the future. Whereupon the result will be 
/// handled and a String returned, for a given key.
pub struct StashRequest {
    request: reqwest::Request, 
    handler: ResponseHandler,
    key: String, 
}

/// request name to request result. For example "title" to "Page Title"
type StashResponse = HashMap<String, String>;

/// Each stash request, and whether it was fullfilled (true) or mocked (false)
type StashLogEntry = (StashRequest, bool);

type Url = String;


#[derive(Debug)]
pub enum StashMode {
    PanicOnMisses,
    MockOnMisses,
}

pub struct Stash {
    mode: StashMode,
    stash: HashMap<Url, StashResponse>,
    log: Vec<StashLogEntry>,
    client: reqwest::Client,
}

impl MarkdownItExt for Stash {}

impl Debug for Stash {
    fn fmt(&self, _f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        // Overriding Debug for scraper, because I don't know how to 
        // implement Debug for ResponseHandler.
        Ok(())
    }
}

impl From<&Path> for Stash {
    fn from(path: &Path) -> Self {
        Self {
            mode: StashMode::PanicOnMisses,
            stash: serde_yaml::from_str(
                std::fs::read_to_string(path.clone())
                .unwrap_or_else(|_| panic!("Failed to read stash {:?}", path))
                .as_str()
            ).unwrap_or_else(|_| panic!("Failed to deserialize stash {:?}. Consider regenerating the cache.", path)),
            log: vec![],
            client: reqwest::Client::new(),
        }
    }
}

impl Stash {
    /// Parse all the input files using a stash that will mock missing values.
    /// A call to `Stash::write`
    pub fn analyse(mut self, content_files: Vec<ContentFile>) -> Stash {
        self.mode = StashMode::MockOnMisses;

        let parser = &mut MarkdownIt::new_project_parser();

        for content_file in content_files {
            content_file.as_content(parser, &mut self);
        }

        self
    }

    fn get_from_stash(&self, stash_request: &StashRequest) -> Option<String> {
        Some(
            self.stash.get(&stash_request.request.url().to_string())?
            .get(&stash_request.key)?
            .clone()
        )
    }

    /// For example
    ///
    /// ```
    /// stash.get(
    ///     "body",
    ///     |client| client.get("https://raypeat.rodeo"),
    ///     |_url, body| body,
    /// )
    /// ```
    pub fn get(
        &mut self, 
        key: &str, 
        request_from_client: impl FnOnce(&reqwest::Client) -> reqwest::Request, 
        string_from_response: impl FnOnce(String, String) -> String + Send + Sync + 'static
    ) -> String {
        let req = StashRequest {
            request: request_from_client(&self.client),
            handler: Box::new(string_from_response),
            key: key.to_string(),
        };

        if let Some(string) = self.get_from_stash(&req) {
            self.log.push((req, true));
            return string;
        }

        match self.mode {
            StashMode::PanicOnMisses => {
                panic!("Haven't stashed {:?} for {:?}", req.key, req.request.url());
            },
            StashMode::MockOnMisses => {
                self.log.push((req, false));
                return "".into();
            },
        }
    }

    pub async fn write(mut self, path: &Path) {
        let mut join_set = JoinSet::new();

        let unstashed_requests = self.log.into_iter()
            .filter(|entry| !entry.1)
            .map(|entry| entry.0);

        for stash_request in unstashed_requests {
            join_set.spawn(async move {
                let url = stash_request.request.url().to_string();

                match reqwest::Client::new().execute(stash_request.request).await {
                    Ok(resp) => {
                        let status = resp.status();

                        if !(status.is_success() || status.is_redirection()) {
                            eprintln!("    {status} for {url} seeking \"{}\"", stash_request.key);
                        }

                        let body = resp.text().await
                            .unwrap_or_else(|_| panic!("Failed to extract body from HTTP response from {}", url));

                        (url.clone(), stash_request.key, (stash_request.handler)(url.clone(), body))
                    },

                    Err(e) => {
                        // TODO returning URL on error does always make sense
                        // Consider adding an error_handler closure argument.
                        eprintln!("    Bad request for {url} seeking \"{}\". Returning \"{url}\" instead. The error was:", stash_request.key);
                        eprintln!("      {e}");
                        (url.clone(), stash_request.key, url.clone())
                    }
                }
            });
        }

        while let Some(result) = join_set.join_next().await {
            let (url, attribute, parsed_response) = result.expect("Error getting HTTP result");

            self.stash
                .entry(url).or_default()
                .entry(attribute).or_insert(parsed_response);
        }

        let mut text = String::from(r#"
# This stash was autogenerated by the --update-stash or --wipe-stash flag.
#
# Ray Peat Rodeo creates a nice development experience by automatically pulling
# data from the internet for the rendering of GitHub issues and mentions.
#
# However, it would be dangerous to assume that this remote data will be the 
# same at deploy time as it is during development and testing. Indeed Nix, the
# package manager we're using, intentionally prevents HTTP request without 
# explicit guarantees in place.
# 
# The solution we're using, is to stash all the data pulled in during 
# development in this file and check it into the source code. During 
# development, the file can be updated using the --update-stash flag. At deploy 
# time we're simply reading from a file.
"#);

        text.push_str(
            serde_yaml::to_string(&self.stash)
            .expect("Failed to serialize scraper to YAML")
            .as_str()
        );

        std::fs::write(path, text).expect("Error scraper to disk");
    }
}
