use std::collections::BTreeMap;
use markdown_it::parser::extset::MarkdownItExt;
use url::Url;
use crate::markdown::mention::Mentionable;

pub mod timecode;
pub mod speaker;
pub mod sidenote;
pub mod mention;

#[derive(Debug)]
pub struct Speakers(pub BTreeMap<String, String>);

impl MarkdownItExt for Speakers {}


#[derive(Debug)]
pub struct Source(pub Url);

impl MarkdownItExt for Source {}


#[derive(Debug)]
pub struct Path(pub std::path::PathBuf);

impl MarkdownItExt for Path {}


#[derive(Debug)]
pub struct Mentions {
    pub people: BTreeMap<Mentionable, u32>,
    pub books: BTreeMap<Mentionable, u32>,
    pub papers: BTreeMap<Mentionable, u32>,
    pub links: BTreeMap<Mentionable, u32>,
}

impl MarkdownItExt for Mentions {}
