use std::fmt::Debug;
use markdown_it::{
    MarkdownIt, Node, NodeValue, Renderer,
    parser::inline::{InlineRule, InlineState},
};


#[derive(Debug)]
pub enum Mention {
    Visible(Mentionable),
    Invisible(Mentionable),
}

impl NodeValue for Mention {
    fn render(&self, node: &Node, fmt: &mut dyn Renderer) {
        let mut render_link = |class: &str, default_alt_text: String| {
            fmt.open("a", &{
                let mut attrs = node.attrs.clone();
                attrs.push(("class", vec!["citation", class].join(" ")));
                attrs.push(("target", "_blank".into()));
                attrs.push(("href", "".into()));
                attrs
            });

            if node.children.len() > 0 {
                fmt.contents(&node.children);
            } else {
                fmt.text(default_alt_text.as_str());
            }

            fmt.close("a");
        };

        use Mention::*;
        use Mentionable::*;

        match self {
            Invisible(_) => return,
            Visible(mentionable) => match mentionable {
                Empty => fmt.contents(&node.children),
                Book { title, primary_author } => render_link("book", title.clone()),
                Person { first_names, last_name } => render_link("person", format!("{} {}", first_names, last_name)),
                Paper { doi } => render_link("paper", doi.clone()),
                Link { url } => render_link("link", url.to_string()),
            },
        };
    }
}


#[derive(Debug)]
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
    Empty,
}


impl Mentionable {
    /// Creates Metionable from mention signature:
    ///
    /// [[mention signature|optional alt text\]\]
    /// --^^^^^^^^^^^^^^^^^
    fn from_signature(mention_signature: &str) -> Mentionable {
        if let Some(doi) = mention_signature.strip_prefix("doi:") {
            return Mentionable::Paper { doi: doi.to_string() };
        } else if let Ok(url) = url::Url::parse(mention_signature) {
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

    /// Converts Metionable into Markdown AST node. If optional alt text is found, it is tokenize
    /// as inline markdown and everything is appened to the Mentionable node as children.
    ///
    /// Expects state.pos to pointing to the character **before** the alt text
    fn into_node(self, state: &mut InlineState) -> Option<Node> {
        if state.src.chars().nth(state.pos)? != '|' {
            state.pos += 1;
            return match self {
                // [[]]
                // --^^
                Mentionable::Empty => None,

                // [[mention signature]]
                // -------------------^^
                _ => Some(Node::new(Mention::Visible(self))),
            };
        }

        state.pos += 1;
        let orig_pos = state.pos;

        while state.pos < state.pos_max {
            if state.src.get(state.pos-1..state.pos+1)? == "]]" {
                return Some({
                    if state.pos <= orig_pos + 2 {
                        match self {
                            // [[|]]
                            // --^-^
                            Mentionable::Empty => return None,

                            // [[mention signature|]]
                            // -------------------^-^
                            _ => Node::new(Mention::Invisible(self)),
                        }
                    } else {
                        // [[mention signature|display text]]
                        // -------------------^-------------^
                        inline_tokenize(Node::new(Mention::Visible(self)), state, orig_pos, state.pos - 1)
                    }
                });
            }
            state.md.inline.skip_token(state);
        }

        return None;
    }
}




/// Candidate for reuse in a utils module
fn inline_tokenize(node: Node, state: &mut InlineState, start: usize, end: usize) -> Node {
    if start >= end { return node; }

    let orig = (state.pos, state.pos_max);

    let original_node = std::mem::replace(&mut state.node, node);
    state.pos = start;
    state.pos_max = end;

    state.md.inline.tokenize(state);

    (state.pos, state.pos_max) = orig;

    std::mem::replace(&mut state.node, original_node)
}


struct MentionInlineScanner;

impl InlineRule for MentionInlineScanner{
    const MARKER: char = '[';

    fn run(state: &mut InlineState) -> Option<(Node, usize)> {
        let orig_pos = state.pos;

        state.pos += {
            if state.src.get(state.pos..state.pos+2)? != "[[" { return None; };
            2
        };

        let node = 'out: {
            let orig_pos = state.pos;

            while state.pos < state.pos_max {
                if state.src.chars().nth(state.pos)? == '|' || state.src.get(state.pos..state.pos+2)? == "]]" {
                    break 'out {
                        if state.pos <= orig_pos {
                            Mentionable::Empty
                        } else {
                            Mentionable::from_signature(&state.src[orig_pos..state.pos])
                        }
                    }.into_node(state)?;
                }
                state.pos += 1;
            }

            return None;
        };

        let consumed = state.pos + 1 - orig_pos;
        state.pos = orig_pos;
        return Some((node, consumed));
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<MentionInlineScanner>();
}
