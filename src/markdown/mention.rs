use std::fmt::Debug;
use markdown_it::{
    MarkdownIt, Node, NodeValue, Renderer,
    parser::inline::{InlineRule, InlineState},
};
use scraper::{Html, Selector};
use crate::{InputFile, scraper::Scraper, MENTION_SLUG};
use serde::{Serialize, Deserialize};

#[derive(PartialEq)]
pub enum Phase {
    Author,
    MentioableKind,
    MentionableDefinition,
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

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub struct Author {
    pub cardinal: String,
    pub prefix: Option<String>,
}

impl Author {
    pub fn id(&self) -> String {
        if let Some(prefix) = &self.prefix {
            format!("{}-{}", prefix, self.cardinal)
                .to_lowercase()
        } else {
            self.cardinal.clone()
        }.to_lowercase().replace(" ", "-")
    }

    pub fn display_text(&self) -> String  {
        if let Some(prefix) = &self.prefix {
            format!("{} {}", prefix, self.cardinal)
        } else {
            self.cardinal.clone()
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub enum WorkKind {
    Book,
    Article,
    Paper,
    Url,
}

impl WorkKind {
    fn id(&self) -> String {
        match self {
            WorkKind::Book => "book",
            WorkKind::Article => "article",
            WorkKind::Url => "url",
            WorkKind::Paper => "paper",
        }.into()
    }
}

#[derive(Debug, Clone)]
pub struct Work {
    pub kind: WorkKind,
    pub signature: String,
    pub title: String,
}

impl Work {
    fn id(&self) -> String {
        format!("{}-{}", self.signature, self.kind.id())
            .to_lowercase()
            .replace(" ", "-")
    }
}

#[derive(Debug, Clone)]
pub struct Mention {
    pub input_file: InputFile,
    pub position: u32,
    pub author: Author,
    pub work: Option<Work>,
}

impl Mention {
    pub fn id(&self) -> String {
        let id = {
            if let Some(work) = &self.work {
                format!("{}-{}", self.author.id(), work.id())
            } else {
                self.author.id()
            }
        };

        if self.position > 1 {
            format!("{}-{}", id, self.position)
        } else {
            id
        }
    }

    pub fn slug(&self) -> String {
        format!("{}#{}", self.input_file.slug, self.id())
    }

    pub fn more_details_slug(&self) -> String {
        if let Some(work) = &self.work {
            format!("/{MENTION_SLUG}/{}#{}", self.author.id(), work.id())
        } else {
            format!("/{MENTION_SLUG}/{}", self.author.id())
        }
    }

    pub fn kind(&self) -> String {
        if let Some(work) = &self.work {
            work.kind.id()
        } else {
            "author".into()
        }
    }

    pub fn display_text(&self) -> String {
        if let Some(work) = &self.work {
            work.title.clone()
        } else {
            self.author.display_text()
        }
    }
}

impl NodeValue for Mention {
    fn render(&self, node: &Node, fmt: &mut dyn Renderer) {
        let mut attrs = node.attrs.clone();

        attrs.push(("id", self.id().to_string()));
        attrs.push(("class", format!("mention {}", self.kind())));
        attrs.push(("href", self.more_details_slug()));
        attrs.push(("data-position", self.position.to_string()));
        attrs.push(("data-display-text", self.display_text()));

        fmt.open("a", &attrs);

        if node.children.is_empty() {
            fmt.text(self.display_text().as_str());
        } else {
            fmt.contents(&node.children);
        }

        fmt.close("a");
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct DoiData {
    pub title: String,
}

#[derive(Debug, Hash, PartialEq, Eq, Clone)]
pub struct MentionDeclaration {
    author_cardinal: String,
    author_prefix: String,
    work_kind: Option<WorkKind>,
    work_signature: String,
}

impl MentionDeclaration {
    pub fn as_mention(&self, input_file: InputFile, position: u32, scraper: &mut Scraper) -> Mention {
        let author = Author {
            cardinal: self.author_cardinal.clone(),
            prefix: {
                if self.author_prefix.is_empty() {
                    None
                } else {
                    Some(self.author_prefix.clone())
                }
            }
        };
        Mention {
            input_file,
            position,
            author,
            work: match &self.work_kind {
                None => None,
                Some(work_kind) => Some(Work {
                    kind: work_kind.clone(),
                    signature: self.work_signature.clone(),
                    title: {
                        match work_kind {
                            WorkKind::Url => scraper.get(
                                "title", 
                                |client| client.get(self.work_signature.clone()).build().unwrap(), 
                                |url, text| {
                                    let html = Html::parse_document(text.as_str());
                                    let title = html
                                        .select(&Selector::parse("head title").unwrap())
                                        .next();

                                    match title {
                                        Some(title) => title.inner_html().clone().trim().to_string(),
                                        None => url,
                                    }
                                }
                            ),
                            WorkKind::Paper => scraper.get(
                                "title",
                                |client| client
                                    .get(format!("https://doi.org/{}", self.work_signature.clone()))
                                    .header("Accept", "application/json; charset=utf-8")
                                    .build()
                                    .expect(format!("Failed to build HTTP request for {}", self.work_signature).as_str()),
                                |url, text| serde_json::from_str::<DoiData>(text.as_str())
                                    .expect(format!("Failed to deserialize JSON HTTP response for {}", url).as_str())
                                    .title.trim().to_string(),
                            ),
                            WorkKind::Book|WorkKind::Article => self.work_signature.clone(),
                        }
                    }
                })
            }
        }
    }
}

impl NodeValue for MentionDeclaration {}


struct MentionInlineScanner {}

impl MentionInlineScanner {
    fn consume_mention(state: &mut InlineState) -> Option<Node> {
        if !state.src.get(state.pos..state.pos_max)?.starts_with("[[") {
            return None;
        }

        state.pos += 2;

        let mut phase = Phase::Author;

        let mut author_cardinal = String::new();
        let mut author_comma_found = false;
        let mut author_prefix = String::new();

        let mut work_kind = String::new();
        let mut work_signature = String::new();

        let mut display_text_start: Option<usize> = None;
        let mut display_text_end: Option<usize> = None;

        while state.pos < state.pos_max {
            let chunk = get_chunk(state)?;

            match phase {
                Phase::Author => {
                    // [["Quoted Authors, Ignore Commas"...
                    if chunk.starts_with('"') {
                        let closer = state.src.get(state.pos+1..state.pos_max)?.find('"')? + 1;
                        let quoted_text = state.src.get(state.pos+1..state.pos+closer)?;
                        if author_comma_found {
                            author_prefix.push_str(quoted_text)
                        } else {
                            author_cardinal.push_str(quoted_text)
                        }
                        state.pos += closer + 1;
                    } else if chunk.starts_with("[") {
                        phase = Phase::MentioableKind;
                    } else if chunk.starts_with("|") {
                        phase = Phase::DisplayText;
                    } else if is_closer(state) {
                        phase = Phase::Done;
                    } else if chunk.starts_with(",") {
                        author_comma_found = true;
                        state.pos += chunk.len();

                    // [[Whybrow, Marcus...
                    } else {
                        if author_comma_found {
                            author_prefix.push_str(chunk);
                        } else {
                            author_cardinal.push_str(chunk);
                        }
                        state.pos += chunk.len();
                    }
                },

                Phase::MentioableKind => {
                    // [[Whybrow, Marcus, [KIND]..
                    if chunk.starts_with("[") && work_kind.is_empty() {
                        let closer = state.src.get(state.pos+1..state.pos_max)?.find("]")? + 1;
                        work_kind = state.src.get(state.pos+1..state.pos+closer)?.to_string();
                        state.pos += closer + 1;
                        phase = Phase::MentionableDefinition;
                    } else {
                        return None;
                    }
                },

                Phase::MentionableDefinition => {
                    if chunk.starts_with('"') {
                        let closer = state.src.get(state.pos+1..state.pos_max)?.find('"')? + 1;
                        let quoted_text = state.src.get(state.pos+1..state.pos+closer)?;
                        work_signature.push_str(quoted_text);
                        state.pos += closer + 1;

                    } else if chunk.starts_with("|") {
                        phase = Phase::DisplayText;

                    } else if is_closer(state) {
                        phase = Phase::Done;

                    } else {
                        work_signature.push_str(chunk.get(0..1)?);
                        state.pos += 1;
                    }
                },

                Phase::DisplayText => {
                    if chunk.starts_with("|") {
                        display_text_start = Some(state.pos+1);
                        state.pos += chunk.len();
                    } else if is_closer(state) {
                        display_text_end = Some(state.pos);
                        phase = Phase::Done;
                    } else {
                        state.md.inline.skip_token(state)
                    }
                },

                Phase::Done => {
                    if is_closer(state) {
                        state.pos += 2;
                        break;
                    }
                }
            }
        }

        let mut node = Node::new(MentionDeclaration {
            author_cardinal: author_cardinal.trim().to_string(),
            author_prefix: author_prefix.trim().to_string(),
            work_kind: match work_kind.trim() {
                "url" => Some(WorkKind::Url),
                "book" => Some(WorkKind::Book),
                "article" => Some(WorkKind::Article),
                "paper" => Some(WorkKind::Paper),
                "" => None,
                _ => return None,
            },
            work_signature: work_signature.trim().to_string(),
        });

        if let (Some(start), Some(end)) = (display_text_start, display_text_end) {
            let orig_node = std::mem::replace(&mut state.node, node);
            let pos = std::mem::replace(&mut state.pos, start);
            let pos_max = std::mem::replace(&mut state.pos_max, end);

            state.md.inline.tokenize(state);

            state.pos = pos;
            state.pos_max = pos_max;
            node = std::mem::replace(&mut state.node, orig_node);
        }


        Some(node)
    }
}

impl InlineRule for MentionInlineScanner {
    const MARKER: char = '[';

    fn run(state: &mut InlineState) -> Option<(Node, usize)> {
        let start = state.pos;

        if let Some(node) = MentionInlineScanner::consume_mention(state) {
            let consumed = state.pos - start;
            state.pos = start;
            return Some((node, consumed))
        }

        state.pos = start;
        return None
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<MentionInlineScanner>();
}
