use std::fmt::Debug;
use markdown_it::{
    MarkdownIt, Node, NodeValue, Renderer,
    parser::inline::{InlineRule, InlineState, Text},

};

#[derive(Debug,Clone)]
pub enum Mention {
    Placeholder(Fragment),
    Hidden(Mentionable),
    Normal(Mentionable, Option<Fragment>),
}

impl Mention {
    fn new(signature: MentionSignature, alt_text: AltText) -> Option<Mention> {
        use MentionSignature::*;
        use Mention::*;

        match (signature, alt_text) {
            (SignifiedButNotProvided, AltText::SignifiedAndProvided(fragment)) => Some(Placeholder(fragment)),
            (SignifiedAndProvided(signature), AltText::NotSignified|AltText::SignifiedButNotProvided) => Some(Hidden(Mentionable::new(signature))),
            (SignifiedAndProvided(signature), alt_text) => Some(Normal(Mentionable::new(signature), Some(alt_text.to_fragment()?))),
            (_, _) => None,
        }
    }

    fn into_node(self, state: &mut InlineState) -> Option<Node> {
        use Mention::*;

        let mut node = Node::new(self.clone());

        match self {
            Hidden(_) => (),
            Normal(mentionable, None) => node.children.push(Node::new(Text { content: mentionable.default_display_text() })),
            Placeholder(fragment) |
            Normal(_, Some(fragment)) => {
                let node = std::mem::replace(&mut state.node, node);
                let pos = std::mem::replace(&mut state.pos, fragment.start);
                let pos_max = std::mem::replace(&mut state.pos_max, fragment.end);

                state.md.inline.tokenize(state);

                state.pos = pos;
                state.pos_max = pos_max;
                return Some(std::mem::replace(&mut state.node, node));
            }
        }

        Some(node)
    }
}

impl NodeValue for Mention {
    fn render(&self, node: &Node, fmt: &mut dyn Renderer) {
        match self {
            Mention::Hidden(_) => return,
            Mention::Placeholder(_) => fmt.contents(&node.children),
            Mention::Normal(mentionable, _) => {
                let mut attrs = node.attrs.clone();

                use Mentionable::*;

                let class = match mentionable {
                    Book { title: _, primary_author: _ } => "book",
                    Person { first_names: _, last_name: _ } => "person",
                    Paper { doi: _ } => "paper",
                    Link { url: _ } => "link",
                };

                attrs.push(("class", vec!["citation", class].join(" ")));
                attrs.push(("target", "_blank".into()));
                attrs.push(("href", match mentionable {
                    Book { title, primary_author } => format!("https://google.com/search?q={} {}", title, primary_author), 
                    Person { first_names, last_name } => format!("https://google.com/search?q={} {}", first_names, last_name),
                    Paper { doi } => format!("https://doi.org/{}", doi),
                    Link { url } => url.to_string(),
                }));

                fmt.open("a", &attrs);
                fmt.contents(&node.children);
                fmt.close("a");
            },
        };
    }
}


#[derive(Debug,Clone)]
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
    fn new(mention_signature: String) -> Mentionable {
        if let Some(doi) = mention_signature.strip_prefix("doi:") {
            return Mentionable::Paper { doi: doi.to_string() };
        } else if let Ok(url) = url::Url::parse(mention_signature.as_str()) {
            return Mentionable::Link { url };
        } else if mention_signature.contains("-by-") {
            let mut segments: Vec<&str> = mention_signature.split("-by-").map(|x| x.trim()).collect();
            return Mentionable::Book {
                primary_author: segments.pop().unwrap().to_string(),
                title: segments.join(" "),
            };
        } else {
            let mut names = mention_signature.split(' ').map(|x| x.trim()).collect::<Vec<&str>>();
            return Mentionable::Person {
                last_name: names.pop().unwrap_or("").to_string(),
                first_names: names.join(" "),
            };
        }
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

#[derive(Debug,Copy,Clone)]
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

        Some((
            Mention::new(
                MentionSignature::new(state)?,
                AltText::new(state)?,
            )?.into_node(state)?,
            std::mem::replace(&mut state.pos, start) - start,
        ))
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<MentionInlineScanner>();
}
