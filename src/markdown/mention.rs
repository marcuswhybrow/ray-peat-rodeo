use std::fmt::Debug;
use std::hash::{Hash, Hasher};
use std::collections::{BTreeMap, hash_map::DefaultHasher};
use base64::{Engine, engine::general_purpose};
use markdown_it::{
    MarkdownIt, Node, NodeValue, Renderer,
    parser::inline::{InlineRule, InlineState, Text},
    parser::extset::RootExt

};

// A Mention is the data structure describing a bespoke markdown inline element
// consumed by InlineMentionScanner, which looks like this:
//
// [[Mentionable Signature|Alternative Display Text]]
//
// The mention signature is the name of a person, the title of a book, a URL,
// or the DOI of a scientific paper. This string is the key by which this
// particular Mention will be grouped with other Mentions in other documents
// that refer to the same Mentionable Signature. By this means a global lookup
// table can be constructed with links to, say, all mentions of "William Blake"
// should that be the "mention signature".
//
// The "alternative display text" is anthing to the right of an (optional) pipe
// character, and allows the document author to customise the link text of this
// particular mention. If not provided, then a sensible default is used.
//
// Both the mentionable signature, and the alternative display text are 
// fields of the Mention struct, and both are derived from the text explicitly
// defined by the author of the markdown document. It's also useful to know if 
// this is the first, second, or nth, mention of a particular mentionable
// signature within this document.
//
// The mention "occurance" field is implicity determined by the number of 
// preceding mentions with the same mention signature. Thus, the first 
// occurance may be directly linked to, or indeed, any occurance.
//
#[derive(Debug, Clone, Hash)]
pub struct Mention {
    // The person, book, URL, or DOI derrived from the mention signature
    mentionable: Mentionable,

    // The position of this mention relative to other metions in the same 
    // document that have the same mention signature.
    occurance: u32,

    // In order to support markdown foratting of the alternative display text,
    // the alt_text must first be captured and consumed by InlineMentionScanner
    // and later explicitly resubmitted to the parser to evince the less greedy
    // rules such as bold, italic, underline, etc.
    alt_text_fragment: Option<Fragment>,
}

impl Mention {
    fn new(state: &mut InlineState, signature: String, alt_text: AltText) -> Option<Mention> {
        let (mentionable, occurance) = Mentionable::new(state, signature);
        Some(Mention { mentionable, occurance, alt_text_fragment: alt_text.to_fragment() })
    }

    fn base64_hash(&self) -> String {
        let mut hasher = DefaultHasher::new();
        self.hash(&mut hasher);
        general_purpose::URL_SAFE_NO_PAD.encode(hasher.finish().to_ne_bytes())
    }
}

impl NodeValue for Mention {
    fn render(&self, node: &Node, fmt: &mut dyn Renderer) {
        let mut attrs = node.attrs.clone();

        use Mentionable::*;

        let class = match self.mentionable {
            Book { title: _, primary_author: _ } => "book",
            Person { first_names: _, last_name: _ } => "person",
            Paper { doi: _ } => "paper",
            Link { url: _ } => "link",
        };

        // Purposfully unfriendly id to indicate they're unreliable
        attrs.push(("id", self.base64_hash()));

        attrs.push(("class", vec!["citation", class].join(" ")));
        attrs.push(("target", "_blank".into()));

        let href = match self.mentionable.clone() {
            Book { title, primary_author } => format!("https://google.com/search?q={} {}", title, primary_author), 
            Person { first_names, last_name } => format!("https://google.com/search?q={} {}", first_names, last_name),
            Paper { doi } => format!("https://doi.org/{}", doi),
            Link { url } => url.to_string(),
        };

        attrs.push(("href", href));

        fmt.open("a", &attrs);
        fmt.contents(&node.children);
        fmt.close("a");
    }
}

#[derive(Debug, Default)]
struct MentionableOccurances(BTreeMap<Mentionable, u32>);

impl RootExt for MentionableOccurances {}


#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord, Hash)]
pub enum Mentionable {
    Person {
        first_names: String,
        last_name: String,
    },
    Book {
        title: String,
        primary_author: String,
    },
    Paper {
        doi: String,
    },
    Link {
        url: url::Url,
    },
}

impl Mentionable {
    fn new(state: &mut InlineState, mention_signature: String) -> (Mentionable, u32) {
        let mentionable = { 
            if let Some(doi) = mention_signature.strip_prefix("doi:") {
                Mentionable::Paper { doi: doi.to_string() }
            } else if let Ok(url) = url::Url::parse(mention_signature.as_str()) {
                Mentionable::Link { url }
            } else if mention_signature.contains("-by-") {
                let mut segments: Vec<&str> = mention_signature.split("-by-").map(|x| x.trim()).collect();
                Mentionable::Book {
                    primary_author: segments.pop().unwrap().to_string(),
                    title: segments.join(" "),
                }
            } else {
                let mut names = mention_signature.split(' ').map(|x| x.trim()).collect::<Vec<&str>>();
                Mentionable::Person {
                    last_name: names.pop().unwrap_or("").to_string(),
                    first_names: names.join(" "),
                }
            }
        };

        let occurance = *state.root_ext
            .get_or_insert_default::<MentionableOccurances>().0
            .entry(mentionable.clone())
            .and_modify(|occurances| *occurances += 1)
            .or_insert(1u32);

        (mentionable, occurance)
    }

    fn default_display_text(&self) -> String {
        use Mentionable::*;

        match self {
            Person { first_names, last_name } => vec![first_names.clone(), last_name.clone()].join(" "),
            Book { title, primary_author: _ } => title.clone(),
            Paper { doi } => doi.clone(),
            Link { url } => url.to_string(),
        }
    }
}


#[derive(Debug)]
enum AltText {
    Provided(Fragment),
    NotProvided,
}

impl AltText {
    fn new(state: &mut InlineState) -> Option<AltText> {
        if state.src.get(state.pos..state.pos+1)? != "|" {
            if state.src.get(state.pos..state.pos+2)? != "]]" {
                panic!("AltText must only be called whilst state.pos is pointing to '|' or the beginning of ']]'");
            } else {
                state.pos += 2;
            }

            return Some(AltText::NotProvided);
        } else {
            state.pos += 1;
        }

        let start = state.pos;

        while state.pos < state.pos_max {
            if state.src.get(state.pos-1..state.pos+1)? == "]]" {
                state.pos += 1;
                return Some(AltText::Provided(Fragment { start, end: state.pos - 2 }));
            }

            state.md.inline.skip_token(state);
        }

        return None;
    }

    fn to_fragment(self) -> Option<Fragment> {
        match self {
            AltText::Provided(fragment) => Some(fragment.clone()),
            _ => None,
        }
    }
}

#[derive(Debug, Copy, Clone, Hash)]
pub struct Fragment {
    start: usize,
    end: usize,
}

impl Fragment {
    fn parse(self, state: &mut InlineState, node: Node) -> Node {
        let node = std::mem::replace(&mut state.node, node);
        let pos = std::mem::replace(&mut state.pos, self.start);
        let pos_max = std::mem::replace(&mut state.pos_max, self.end);

        state.md.inline.tokenize(state);

        state.pos = pos;
        state.pos_max = pos_max;
        std::mem::replace(&mut state.node, node)
    }
}


fn consume_mention_signature(state: &mut InlineState) -> Option<String> {
    let start = state.pos;

    while state.pos < state.pos_max {
        if state.src.get(state.pos..state.pos+1)? == "|" || state.src.get(state.pos..state.pos+2)? == "]]" {
            if state.pos <= start {
                return None;
            } else {
                return Some(state.src.get(start..state.pos)?.to_string());
            }
        }

        state.pos += 1;
    }

    return None;
}


struct MentionInlineScanner;

impl InlineRule for MentionInlineScanner{
    const MARKER: char = '[';

    fn run(state: &mut InlineState) -> Option<(Node, usize)> {
        let start = state.pos;

        if state.src.get(state.pos..state.pos+2)? != "[[" { return None; }
        state.pos += 2;

        let mention_signature = consume_mention_signature(state)?;
        let alt_text = AltText::new(state)?;
        let mention = Mention::new(state, mention_signature, alt_text)?;

        let mut node = Node::new(mention.clone());

        if let Some(fragment) = mention.alt_text_fragment {
            node = fragment.parse(state, node);
        } else {
            node.children.push(Node::new(Text { content: mention.mentionable.default_display_text() }));
        }

        Some((node, std::mem::replace(&mut state.pos, start) - start))
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<MentionInlineScanner>();
}
