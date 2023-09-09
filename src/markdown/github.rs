use markdown_it::{
    MarkdownIt, Node,
    parser::inline::{InlineRule, InlineState},
    plugins::cmark::inline::link::Link,
    parser::inline::Text,
    plugins::html::html_inline::HtmlInline,
};

use crate::markdown::sidenote::{
    InlineSidenote, Position,
};

struct GitHubIssueInlineScanner;

impl InlineRule for GitHubIssueInlineScanner {
    const MARKER: char = '{';

    fn run(state: &mut InlineState) -> Option<(Node, usize)> {
        if state.level > 0 || state.src.get(state.pos..state.pos+2)? != "{#" {
            return None;
        }

        let start = state.pos;
        state.pos += 2;

        while state.pos < state.pos_max {
            state.pos += 1;

            if state.src.get(state.pos..state.pos+1)? != "}" {
                continue;
            }

            let consumed = state.pos + 1 - start;

            let id = state.src.get(start+2..state.pos)?;

            let mut link = Node::new(Link {
                url: String::from(
                    std::format!("https://github.com/marcuswhybrow/ray-peat-rodeo/issues/{}", id)
                ),
                title: Some(String::from("GitHub Issue")),
            });

            link.children.push(Node::new(Text {
                content: String::from(format!("GitHub Issue #{}", id)),
            }));

            let mut node = Node::new(
                InlineSidenote::new(
                    state.root_ext.get_or_insert(Position(0))
                )
            );
            node.children.push(Node::new(Text {
                content: String::from("This mention needs indentifying. If you know more, let us know on "),
            }));
            node.children.push(link);
            node.children.push(Node::new(Text {
                content: String::from("."),
            }));

            state.pos = start;

            return Some((node, consumed));
        }

        None
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<GitHubIssueInlineScanner>();
}
