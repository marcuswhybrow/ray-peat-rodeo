use std::fmt::Debug;
use std::hash::{Hash, Hasher};
use std::collections::{BTreeMap, hash_map::DefaultHasher};
use base64::{Engine, engine::general_purpose};
use markdown_it::{
    MarkdownIt, Node, NodeValue, Renderer,
    parser::inline::{InlineRule, InlineState, Text},
    parser::extset::RootExt

};

#[derive(Debug, Clone, Hash, serde::Serialize)]
pub enum Mention {
    Placeholder {
        fragment: Fragment
    },
    Hidden {
        mentionable: Mentionable,
        occurance: u32,
    },
    Normal {
        mentionable: Mentionable,
        occurance: u32,
        fragment: Option<Fragment>,
    },
}

impl Mention {
    fn new(state: &mut InlineState, signature: MentionSignature, alt_text: AltText) -> Option<Mention> {
        use MentionSignature::*;
        use Mention::*;

        match (signature, alt_text) {
            (SignifiedButNotProvided, AltText::SignifiedAndProvided(fragment)) => Some(Placeholder { fragment }),
            (SignifiedAndProvided(signature), AltText::NotSignified|AltText::SignifiedButNotProvided) => {
                let (mentionable, occurance) = Mentionable::new(state, signature);
                Some(Hidden { mentionable, occurance })
            },
            (SignifiedAndProvided(signature), alt_text) => {
                let (mentionable, occurance) = Mentionable::new(state, signature);
                let fragment = Some(alt_text.to_fragment()?);
                Some(Normal { mentionable, occurance, fragment })
            },
            (_, _) => None,
        }
    }

    pub fn base64_hash(&self) -> String {
        let mut hasher = DefaultHasher::new();
        self.hash(&mut hasher);
        general_purpose::URL_SAFE_NO_PAD.encode(hasher.finish().to_ne_bytes())
    }

    fn into_node(self, state: &mut InlineState) -> Option<Node> {
        use Mention::*;

        let mut node = Node::new(self.clone());

        match self {
            Hidden { mentionable: _, occurance: _ } => (),
            Normal { mentionable, occurance: _, fragment: None } => {
                node.children.push(
                    Node::new(Text {
                        content: mentionable.default_display_text()
                    })
                );
            },
            Placeholder { fragment } |
            Normal { mentionable: _, occurance: _, fragment: Some(fragment) } => {
                let node = std::mem::replace(&mut state.node, node);
                let pos = std::mem::replace(&mut state.pos, fragment.start);
                let pos_max = std::mem::replace(&mut state.pos_max, fragment.end);

                state.md.inline.tokenize(state);

                state.pos = pos;
                state.pos_max = pos_max;
                return Some(std::mem::replace(&mut state.node, node));
            }
        };

        Some(node)
    }
}

impl NodeValue for Mention {
    fn render(&self, node: &Node, fmt: &mut dyn Renderer) {
        use Mention::*;

        match self {
            Hidden { mentionable: _, occurance: _ } => return,
            Placeholder { fragment: _ } => fmt.contents(&node.children),
            Normal { mentionable, occurance: _, fragment: _ } => {
                let mut attrs = node.attrs.clone();

                use Mentionable::*;

                let class = match mentionable {
                    Book(_) => "book",
                    Person(_) => "person",
                    Paper(_) => "paper",
                    Link(_) => "link",
                };

                // Purposfully unfriendly id to indicate they're unreliable
                attrs.push(("id", self.base64_hash()));

                attrs.push(("class", vec!["citation", class].join(" ")));
                attrs.push(("target", "_blank".into()));
                attrs.push(("href", match mentionable {
                    Book(b) => format!("https://google.com/search?q={} {}", b.title, b.primary_author), 
                    Person(p) => format!("https://google.com/search?q={} {}", p.first_names, p.last_name),
                    Paper(p) => format!("https://doi.org/{}", p.doi),
                    Link(l) => l.url.to_string(),
                }));

                fmt.open("a", &attrs);
                fmt.contents(&node.children);
                fmt.close("a");
            },
        };
    }
}

#[derive(Debug, Default)]
struct MentionableOccurances(BTreeMap<Mentionable, u32>);

impl RootExt for MentionableOccurances {}

#[derive(Debug, Clone, Eq, PartialEq, PartialOrd, Ord, Hash, serde::Serialize)]
pub struct Person {
    first_names: String,
    last_name: String,
}

#[derive(Debug, Clone, Eq, PartialEq, PartialOrd, Ord, Hash, serde::Serialize)]
pub struct Book {
    title: String,
    primary_author: String,
}

#[derive(Debug, Clone, Eq, PartialEq, PartialOrd, Ord, Hash, serde::Serialize)]
pub struct Paper {
    doi: String,
}

#[derive(Debug, Clone, Eq, PartialEq, PartialOrd, Ord, Hash, serde::Serialize)]
pub struct Link {
    url: String,
}

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord, Hash, serde::Serialize)]
pub enum Mentionable {
    Person(Person),
    Book(Book),
    Paper(Paper),
    Link(Link),
}

impl Mentionable {
    fn new(state: &mut InlineState, mention_signature: String) -> (Mentionable, u32) {
        let mentionable = { 
            if let Some(doi) = mention_signature.strip_prefix("doi:") {
                Mentionable::Paper(Paper { doi: doi.to_string() })
            } else if let Ok(url) = url::Url::parse(mention_signature.as_str()) {
                Mentionable::Link(Link { url: mention_signature })
            } else if mention_signature.contains("-by-") {
                let mut segments: Vec<&str> = mention_signature.split("-by-").map(|x| x.trim()).collect();
                Mentionable::Book(Book {
                    primary_author: segments.pop().unwrap().to_string(),
                    title: segments.join(" "),
                })
            } else {
                let mut names = mention_signature.split(' ').map(|x| x.trim()).collect::<Vec<&str>>();
                Mentionable::Person(Person {
                    last_name: names.pop().unwrap_or("").to_string(),
                    first_names: names.join(" "),
                })
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
            Person(person) => vec![person.first_names.clone(), person.last_name.clone()].join(" "),
            Book(book) => book.title.clone(),
            Paper(paper) => paper.doi.clone(),
            Link(link)=> link.url.to_string(),
        }
    }
}

#[derive(Debug)]
enum MentionSignature {
    SignifiedAndProvided(String),
    SignifiedButNotProvided,
}

impl MentionSignature {
    fn new(state: &mut InlineState) -> Option<MentionSignature> {
        let start = state.pos;

        while state.pos < state.pos_max {
            if state.src.get(state.pos..state.pos+1)? == "|" || state.src.get(state.pos..state.pos+2)? == "]]" {
                if state.pos <= start {
                    return Some(MentionSignature::SignifiedButNotProvided);
                } else {
                    return Some(MentionSignature::SignifiedAndProvided(state.src.get(start..state.pos)?.to_string()));
                }
            }

            state.pos += 1;
        }

        return None;
    }
}

#[derive(Debug)]
enum AltText {
    SignifiedAndProvided(Fragment),
    SignifiedButNotProvided,
    NotSignified,
}

impl AltText {
    fn new(state: &mut InlineState) -> Option<AltText> {
        if state.src.get(state.pos..state.pos+1)? != "|" {
            if state.src.get(state.pos..state.pos+2)? != "]]" {
                panic!("AltText must only be called whilst state.pos is pointing to '|' or the beginning of ']]'");
            } else {
                state.pos += 2;
            }

            return Some(AltText::NotSignified);
        } else {
            state.pos += 1;
        }

        let start = state.pos;

        while state.pos < state.pos_max {
            if state.src.get(state.pos-1..state.pos+1)? == "]]" {
                let final_pos = state.pos;
                state.pos += 1;

                if final_pos <= start + 2 {
                    return Some(AltText::SignifiedButNotProvided);
                } else {
                    return Some(AltText::SignifiedAndProvided(Fragment { start, end: state.pos - 2 }));
                }
            }

            state.md.inline.skip_token(state);
        }

        return None;
    }

    fn to_fragment(self) -> Option<Fragment> {
        match self {
            AltText::SignifiedAndProvided(fragment) => Some(fragment.clone()),
            _ => None,
        }
    }
}

#[derive(Debug, Copy, Clone, Hash, serde::Serialize)]
pub struct Fragment {
    start: usize,
    end: usize,
}

struct MentionInlineScanner;

impl InlineRule for MentionInlineScanner{
    const MARKER: char = '[';

    fn run(state: &mut InlineState) -> Option<(Node, usize)> {
        let start = state.pos;

        if state.src.get(state.pos..state.pos+2)? != "[[" { return None; }
        state.pos += 2;

        let mention_signature = MentionSignature::new(state)?;
        let alt_text = AltText::new(state)?;

        let mention = Mention::new(state, mention_signature, alt_text)?;

        Some((
            mention.into_node(state)?,
            std::mem::replace(&mut state.pos, start) - start,
        ))
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<MentionInlineScanner>();
}
