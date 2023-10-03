pub mod timecode;
pub mod speaker;
pub mod sidenote;
pub mod mention;
pub mod github;

use std::collections::BTreeMap;
use serde::{Serialize, Deserialize};

#[derive(Debug, Clone, Hash, PartialEq, Eq, Serialize, Deserialize)]
pub struct Transcription {
    pub url: Option<String>,
    pub author: Option<String>,
    pub date: Option<String>,
}

#[derive(Debug, Clone, Hash, PartialEq, Eq, Serialize, Deserialize)]
pub struct Source {
    pub title: String,
    pub series: Option<String>,
    pub url: String,
    pub duration: Option<String>,
}

#[derive(Debug, Clone, Hash, PartialEq, Eq, Serialize, Deserialize)]
pub struct Frontmatter {
    pub source: Source,
    pub speakers: BTreeMap<String, String>, 
    pub transcription: Option<Transcription>,
}
