use std::{fmt::Debug, path::Path, collections::HashMap, mem::replace};
use markdown_it::{
    MarkdownIt, Node, NodeValue, Renderer,
    parser::inline::{InlineRule, InlineState},
};
use scraper::{Html, Selector};
use url::{Url, Host};
use crate::{stash::Stash, MENTION_SLUG, write, BasePage, content::ContentFile};
use serde::{Serialize, Deserialize};

use itertools::Itertools;

#[derive(PartialEq, Clone)]
pub enum Phase {
    SubMention,
    DisplayText,
    Done,
}

fn get_chunk<'a>(state: &'a InlineState) -> Option<&'a str> {
    for x in state.pos+1..state.pos_max {
        if let Some(chunk) = state.src.get(state.pos..x) {
            return Some(chunk);
        }
    }
    None
}

fn is_closer(state: &InlineState) -> bool {
    if let Some(s) = state.src.get(state.pos..state.pos+2) {
        if s.starts_with("]]") {
            return true;
        }
    }
    false
}

markup::define! {
    MentionPage<'a>(mentionable: &'a Mentionable, direct_mentions: DirectMentions<'a>, sub_mentions: SubMentions<'a>) {
        @BasePage {
            title: Some(mentionable.human_readable().as_str()),
            content: markup::new! {
                article.mentions {
                    section."popup-select" {
                        h1 { 
                            a [href = mentionable.permalink()] {
                                @mentionable.human_readable() 
                            }
                        }

                        @if direct_mentions.len() > 0 {
                            ul {
                                @for (content_file, mentions) in direct_mentions
                                    .iter().sorted_by_key(|x| x.0.frontmatter.source.title.clone()) 
                                {
                                    li.content {
                                        a[href = mentions.first().unwrap().permalink()] {
                                            @content_file.frontmatter.source.title
                                        }
                                    }
                                }
                            }
                        }

                        @for (sub_mentionable, content_file_to_mentions) in sub_mentions
                            .iter().sorted_by_key(|x| x.0.default_display_text.clone()) 
                        {
                            h2.mentionable {
                                @match sub_mentionable.as_url() {
                                    Some(url) => {
                                        a [ href = url.to_string() ] {
                                            @sub_mentionable.default_display_text
                                        }
                                    }
                                    None => { @sub_mentionable.default_display_text }
                                }
                            }

                            ul {
                                @for (content_file, mentions) in content_file_to_mentions
                                    .iter().sorted_by_key(|x| x.0.frontmatter.source.title.clone()) 
                                {
                                    li.content {
                                        a[href = format!("/{}", mentions.first().unwrap().permalink())] {
                                            @content_file.frontmatter.source.title
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            },
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct DoiData {
    pub title: String,
}

#[derive(Debug, Clone)]
pub struct Mention {
    pub content_file: ContentFile,
    pub position: u32,
    pub mentionable: Mentionable,
    pub sub_mentionable: Option<Mentionable>,
}

impl Mention {
    pub fn id(&self) -> String {
        let id = match &self.sub_mentionable {
            Some(sub_mention) => format!("{}-{}", self.mentionable.id(), sub_mention.id()),
            None => self.mentionable.id(),
        };

        if self.position == 1 {
            return id;
        }

        format!("{}-{}", id, self.position)
    }

    /// The absolute URL to the place in a page where this mention is made
    pub fn permalink(&self) -> String {
        format!("{}#{}", self.content_file.permalink(), self.id())
    }

    /// The absolute URL to more details about the thing being mentioned
    pub fn mentionable_permalink(&self) -> String {
        match &self.sub_mentionable {
            None  => self.mentionable.permalink(),
            Some(sub_mentionable) => self.mentionable.sub_mentionable_permalink(sub_mentionable),
        }
    }

    pub fn ultimate(&self) -> &Mentionable {
        match &self.sub_mentionable {
            Some(sub_mention) => sub_mention,
            None => &self.mentionable,
        }
    }

    pub fn as_url(&self) -> Option<Url> {
        self.ultimate().as_url()
    }

}

impl NodeValue for Mention {
    fn render(&self, node: &Node, fmt: &mut dyn Renderer) {
        fmt.text_raw(markup::new! {
            span [
                id = self.id(),
                class = "mention",
                "hx-trigger" = "mouseenter",
                "hx-target" = "find .popup-card",
                "hx-get" = self.mentionable_permalink(),
                "hx-swap" = "innerHTML",
                "hx-select" = ".popup-select",
            ] {
                a [ href = self.mentionable_permalink() ] {
                    @match node.children.is_empty() {
                        false => { @node.collect_text() }
                        true => { @{self.ultimate().default_display_text.clone()} }
                    }
                }

                span."popup-card" {}
            }
        }.to_string().as_str());
    }
}

type DirectMentions<'a> = HashMap<&'a ContentFile, Vec<&'a Mention>>;
type SubMentions<'a> = HashMap<&'a &'a Mentionable, HashMap<&'a ContentFile, Vec<&'a &'a Mention>>>;

#[derive(Debug, Hash, PartialEq, Eq, PartialOrd, Ord, Clone)]
pub struct MentionableDeclaration {
    pub cardinal: String,
    pub prefix: String,
}

impl MentionableDeclaration {
    pub fn has_prefix(&self) -> bool {
        !self.prefix.is_empty()
    }

    #[allow(dead_code)]
    pub fn cardinal_first(&self) -> String {
        let mut result = self.cardinal.clone();
        if self.has_prefix() {
            result.push_str(", ");
            result.push_str(self.prefix.as_str());
        }
        result
    }

    pub fn human_readable(&self) -> String {
        format!("{} {}", self.prefix, self.cardinal)
            .trim()
            .to_string()
    }


    pub fn as_url(&self) -> Option<Url> {
        if self.has_prefix() {
            None
        } else {
            Url::parse(self.cardinal.as_str()).ok()
        }
    }

    pub fn as_mentionable(&self, stash: &mut Stash) -> Mentionable {
        Mentionable {
            cardinal: self.cardinal.clone(),
            prefix: self.prefix.clone(),
            default_display_text: match self.as_url() {
                Some(url) => match url.host() {
                    Some(Host::Domain("doi.org")) => {
                        stash.get(
                            "title",
                            |client| client
                                .get(url.clone())
                                .header("Accept", "application/json; charset=utf-8")
                                .build()
                                .unwrap_or_else(|_| panic!("Failed to build HTTP request for {url}")),
                            |url, text| serde_json::from_str::<DoiData>(text.as_str())
                                .unwrap_or_else(|_| panic!("Failed to deserialize JSON HTTP response for {url}"))
                                .title.trim().to_string(),
                        )
                    },
                    Some(_) => {
                        stash.get(
                            "title", 
                            |client| client.get(url.clone()).build()
                                .unwrap_or_else(|_| panic!("Failed to build HTTP request for {url}")), 
                            |url, text| {
                                let html = Html::parse_document(text.as_str());
                                let title = html
                                    .select(
                                        &Selector::parse("head title")
                                            .unwrap_or_else(|_| panic!("Failed to parse title selector"))
                                    )
                                    .next();

                                match title {
                                    Some(title) => title.inner_html().clone().trim().to_string(),
                                    None => url,
                                }
                            }
                        )
                    },
                    None => self.human_readable(),
                }
                None => self.human_readable(),
            }
        }
    }
}

#[derive(Debug, PartialEq, Eq, Hash, PartialOrd, Ord, Clone)]
pub struct Mentionable {
    pub cardinal: String,
    pub prefix: String,
    pub default_display_text: String,
}

impl Mentionable {
    pub fn has_prefix(&self) -> bool {
        !self.prefix.is_empty()
    }

    pub fn cardinal_first(&self) -> String {
        let mut result = self.cardinal.clone();
        if self.has_prefix() {
            result.push_str(", ");
            result.push_str(self.prefix.as_str());
        }
        result
    }

    pub fn human_readable(&self) -> String {
        format!("{} {}", self.prefix, self.cardinal)
            .trim()
            .to_string()
    }

    pub fn id(&self) -> String {
        self.human_readable().replace(" ", "-").to_lowercase()
    }

    pub fn as_url(&self) -> Option<Url> {
        if self.has_prefix() {
            None
        } else {
            Url::parse(self.cardinal.as_str()).ok()
        }
    }

    pub fn permalink(&self) -> String {
        format!("/{MENTION_SLUG}/{}", self.id())
    }

    pub fn sub_mentionable_permalink(&self, sub_mentionable: &Mentionable) -> String {
        format!("{}#{}", self.permalink(), sub_mentionable.id())
    }

    pub fn write(&self, output: &Path, mentions: Vec<&Mention>) {
        let mentions = mentions.into_iter()
            .sorted_by_key(|m| m.position);

        write(&output.join(format!("{MENTION_SLUG}/{}/index.html", self.id())), markup::new! {
            @MentionPage {
                mentionable: &self,
                direct_mentions: mentions.clone()
                    .filter(|m| m.sub_mentionable.is_none())
                    .into_group_map_by(|m| &m.content_file),
                sub_mentions: mentions.clone()
                    .filter(|m| m.sub_mentionable.is_some())
                    .into_group_map_by(|m| m.sub_mentionable.as_ref()
                        .unwrap_or_else(|| panic!("Expect Sub Mentionable to exist")))
                    .iter().map(|(mentionable, mentions)| {
                        let x = 
                        (mentionable, mentions.iter().into_group_map_by(|m| &m.content_file))
                        ; x
                    }).collect::<SubMentions>(),

            }
        });
    }
}

#[derive(Debug, Eq, PartialEq, Hash, Clone)]
pub struct MentionDeclaration {
    mentionable_declaration: MentionableDeclaration,
    sub_mentionable_declaration: Option<MentionableDeclaration>,
}

impl MentionDeclaration {
    pub fn as_mention(&self, content_file: ContentFile, count: u32, stash: &mut Stash) -> Mention {
        Mention {
            position: count,
            content_file,
            mentionable: self.mentionable_declaration.as_mentionable(stash),
            sub_mentionable: self.sub_mentionable_declaration.as_ref()
                .and_then(|m| Some(m.as_mentionable(stash))),
        }
    }
}

impl NodeValue for MentionDeclaration {}


type SrcFragment = (usize, usize);


struct MentionInlineScanner {}

impl MentionInlineScanner {
    fn consume_quoted_text(state: &mut InlineState) -> Option<String> {
        let (pos, max) = (state.pos, state.pos_max);

        let closer = state.src.get(pos+1..max)?.find('"')? + 1;
        state.pos += closer + 1;

        let quoted_text = state.src.get(pos+1..pos+closer)?;
        Some(quoted_text.to_string())
    }

    fn consume_mentionable(state: &mut InlineState, until: &[(&str, Phase)]) -> Option<(MentionableDeclaration, Phase)> {
        let mut cardinal = String::new();
        let mut prefix = String::new();
        let mut target = &mut cardinal;

        while state.pos < state.pos_max {
            let chunk = get_chunk(state)?;

            for (terminator, phase) in until {
                 if state.src.get(state.pos..state.pos_max)?.starts_with(terminator) {
                     state.pos += terminator.len();
                     return Some((
                         MentionableDeclaration { 
                             cardinal: cardinal.trim().into(),
                             prefix: prefix.trim().into() 
                         }, 
                         (*phase).clone()
                     ));
                 }
            }

            if chunk.starts_with('"') {
                target.push_str(Self::consume_quoted_text(state)?.as_str());
            } else if chunk.starts_with(",") {
                target = &mut prefix;
                state.pos += chunk.len();
            } else {
                target.push_str(chunk);
                state.pos += chunk.len();
            }
        }

        None
    }

    fn consume_display_text(state: &mut InlineState) -> Option<(usize, usize)> {
        let start: usize = state.pos;

        while state.pos < state.pos_max {
            if is_closer(state) {
                state.pos += 2;
                return Some((start, state.pos-2));
            }

            state.md.inline.skip_token(state);
        }

        None
    }

    fn consume_mention(state: &mut InlineState) -> Option<Node> {
        if !state.src.get(state.pos..state.pos_max)?.starts_with("[[") {
            return None;
        }

        state.pos += 2;

        let (mentionable, mut phase) = Self::consume_mentionable(state, &[
            (">", Phase::SubMention),
            ("|", Phase::DisplayText),
            ("]]", Phase::Done),
        ])?;
        let mut sub_mentionable: Option<MentionableDeclaration> = None;
        let mut display_text: Option<SrcFragment> = None;

        while state.pos < state.pos_max {
            match phase {
                Phase::SubMention => {
                    let (mentionable, next_phase) = Self::consume_mentionable(state, &[
                        ("|", Phase::DisplayText),
                        ("]]", Phase::Done),
                    ])?;
                    sub_mentionable = Some(mentionable);
                    phase = next_phase;
                },

                Phase::DisplayText => {
                    display_text = Some(Self::consume_display_text(state)?);
                    phase = Phase::Done;
                },

                Phase::Done => {
                    break;
                }
            }
        }

        let mut node = Node::new(MentionDeclaration {
            mentionable_declaration: mentionable,
            sub_mentionable_declaration: sub_mentionable,
        });

        if let Some((start, end)) = display_text {
            let orig_node = replace(&mut state.node, node);
            let pos = replace(&mut state.pos, start);
            let pos_max = replace(&mut state.pos_max, end);

            state.md.inline.tokenize(state);

            state.pos = pos;
            state.pos_max = pos_max;
            node = replace(&mut state.node, orig_node);
        }


        Some(node)
    }
}

impl InlineRule for MentionInlineScanner {
    const MARKER: char = '[';

    fn run(state: &mut InlineState) -> Option<(Node, usize)> {
        let start = state.pos;

        match Self::consume_mention(state) {
            Some(node) => {
                let consumed = state.pos - start;
                state.pos = start;
                Some((node, consumed))
            },
            None => {
                state.pos = start;
                None
            }
        }
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<MentionInlineScanner>();
}
