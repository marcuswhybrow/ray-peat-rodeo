use std::collections::BTreeMap;
use markdown_it::parser::extset::MarkdownItExt;
use url::Url;

pub mod timecode;
pub mod speaker;
pub mod sidenote;

#[derive(Debug)]
pub struct Speakers(pub BTreeMap<String, String>);

impl MarkdownItExt for Speakers {}


#[derive(Debug)]
pub struct Source(pub Url);

impl MarkdownItExt for Source {}


#[derive(Debug)]
pub struct Path(pub std::path::PathBuf);

impl MarkdownItExt for Path {}
