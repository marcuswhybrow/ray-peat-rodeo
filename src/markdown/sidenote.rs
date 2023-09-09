use markdown_it::{
    MarkdownIt, Node, NodeValue, Renderer,
    parser::inline::{InlineRule, InlineState},
    parser::extset::RootExt,
};


#[derive(Debug)]
pub struct InlineSidenote {
    position: u32,
}

impl NodeValue for InlineSidenote {
    fn render(&self, node: &Node, fmt: &mut dyn Renderer) {
        let id = format!("sidenote-{}", self.position);

        let mut attrs = node.attrs.clone();
        attrs.push(("for", id.clone()));
        attrs.push(("class", "sidenote-toggle sidenote-number".into()));

        fmt.open("label", &attrs);
        fmt.close("label");
        fmt.self_close("input", &[
            ("type", "checkbox".into()),
            ("id", id),
            ("class", "sidenote-toggle".into()),
        ]);
        fmt.open("span", &[("class", "sidenote".into())]);
        fmt.contents(&node.children);
        fmt.close("span");
    }
}

#[derive(Debug)]
struct Position(u32);

impl RootExt for Position {}

struct SidenodeInlineScanner;

impl InlineRule for SidenodeInlineScanner {
    const MARKER: char = '{';

    fn run(state: &mut InlineState) -> Option<(Node, usize)> {
        if state.level > 0 || state.src.get(state.pos..state.pos+1)? != "{" {
            return None;
        }

        let starting_pos = state.pos;
        state.pos += 1;

        while state.pos < state.pos_max {
            // Increase state.pos to overcome the next token.
            // This prevents finding a closing "}" that matches a nested rule
            state.md.inline.skip_token(state);

            // Keep skipping tokens until we find the closing "}"
            if state.src.get(state.pos..state.pos+1)? != "}" {
                continue;
            }

            let consumed = state.pos + 1 - starting_pos;

            // Create a new InlineSidenote node, and use the existing parser 
            // state to tokenize the body.
            let node = {
                let node = std::mem::replace(&mut state.node, Node::new(InlineSidenote {
                    position: {
                        let position = state.root_ext.get_or_insert(Position(0));
                        position.0 += 1;
                        position.0
                    }
                }));
                let pos = std::mem::replace(&mut state.pos, starting_pos + 1);
                let pos_max = std::mem::replace(&mut state.pos_max, pos);

                state.md.inline.tokenize(state);

                state.pos = starting_pos;
                state.pos_max = pos_max;
                std::mem::replace(&mut state.node, node)
            };

            return Some((node, consumed));
        }

        None
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<SidenodeInlineScanner>();
}
