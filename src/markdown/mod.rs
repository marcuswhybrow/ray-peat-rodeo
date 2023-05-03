use std::collections::BTreeMap;
use markdown_it::parser::extset::MarkdownItExt;
use crate::markdown::mention::Mentionable;

pub mod timecode;
pub mod speaker;
pub mod sidenote;
pub mod mention;

#[derive(Debug)]
pub struct Mentions {
    pub people: BTreeMap<Mentionable, u32>,
    pub books: BTreeMap<Mentionable, u32>,
    pub papers: BTreeMap<Mentionable, u32>,
    pub links: BTreeMap<Mentionable, u32>,
}

impl MarkdownItExt for Mentions {}
