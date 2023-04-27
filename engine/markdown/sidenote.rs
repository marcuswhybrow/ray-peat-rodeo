use markdown_it::{
    MarkdownIt, Node, NodeValue, Renderer,
    parser::inline::{InlineRule, InlineState},
    parser::core::CoreRule,
};


#[derive(Debug)]
pub struct InlineSidenoteWithoutPosition;

impl NodeValue for InlineSidenoteWithoutPosition {
    fn render(&self, _: &Node, _: &mut dyn Renderer) {
        panic!("TempInlineSidenote must be replaced with InlineSidenote before rendering");
    }
}

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

struct SidenodeInlineScanner;

impl InlineRule for SidenodeInlineScanner {
    const MARKER: char = '{';

    fn run(state: &mut InlineState) -> Option<(Node, usize)> {
        if state.level > 0 { return None } // (nested) inside img or link

        let starting_pos = state.pos;

        let Some(char) = state.src.chars().nth(state.pos) else { return None };
        if char != '{' { return None };

        state.pos += 1;

        while state.pos < state.pos_max {
            let char = {
                state.md.inline.skip_token(state); // increases state.pos to overcome the next token 

                if state.pos >= state.pos_max { return None }

                match state.src.chars().nth(state.pos) {
                    Some(char) => char,
                    None => return None,
                }
            };

            if char != '}' { continue }

            let Some(next_char) = state.src.chars().nth(state.pos + 1) else { return None };
            if !next_char.is_whitespace() { return None }

            let consumed = state.pos + 1 - starting_pos;

            // Inspiration: https://github.com/rlidwka/markdown-it.rs/blob/eb5459039685d19cefd0361859422118d08d35d4/src/generics/inline/full_link.rs#L124-L136
            let node = {
                let original_node = std::mem::replace(&mut state.node, Node::new(InlineSidenoteWithoutPosition));

                let original_pos_max = state.pos_max;
                state.pos_max = state.pos;
                state.pos = starting_pos + 1;
                state.md.inline.tokenize(state);

                state.pos = starting_pos;
                state.pos_max = original_pos_max;

                std::mem::replace(&mut state.node, original_node)
            };

            return Some((node, consumed));
        }

        None
    }
}

struct SidenoteCalcPositionRule;

impl CoreRule for SidenoteCalcPositionRule {
    fn run(root: &mut Node, _: &MarkdownIt) {
        let mut counter = 1u32;

        root.walk_mut(|node, _depth| {
            if node.is::<InlineSidenoteWithoutPosition>() {
                node.replace(InlineSidenote {
                    position: counter,
                });
                counter = counter + 1;
            }
        });
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<SidenodeInlineScanner>();
    md.add_rule::<SidenoteCalcPositionRule>();
}
