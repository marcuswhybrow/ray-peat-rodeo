use markdown_it::{MarkdownIt, Node, NodeValue, Renderer};
use markdown_it::parser::block::{BlockRule, BlockState};
use markdown_it::common::sourcemap::SourcePos;

#[derive(Debug)]
pub struct TempSpeakerSection {
    pub shortname: String,
}

impl NodeValue for TempSpeakerSection {
    fn render(&self, _: &Node, _: &mut dyn Renderer) {
        panic!("TempSpeakerSection must be replace with SpeakerSection before rendering");
    }
}

#[derive(Debug)]
pub struct SpeakerSection {
    pub shortname: String,
    pub longname: String,
}

impl NodeValue for SpeakerSection {
    fn render(&self, node: &Node, fmt: &mut dyn Renderer) {
        let mut attrs = node.attrs.clone();
        attrs.push(("class", "speaker".into()));
        attrs.push(("data-shortname", self.shortname.clone()));
        attrs.push(("data-longname", self.longname.clone()));
        fmt.cr();
        fmt.open("div", &attrs);
        fmt.text_raw(format!("<span class=\"speaker-name\">{}</span>", self.longname).as_str());
        fmt.contents(&node.children);
        fmt.close("div");
        fmt.cr();
    }
}


fn get_speaker_shortname(line: &str) -> Option<&str> {
    let mut chars = line.chars().enumerate();
    for (idx, char) in chars.by_ref() {
        if char == ':' {
            if idx <= 0 { return None }

            let (_, next) = match chars.next() {
                Some(c) => c,
                None => return None,
            };

            if !next.is_whitespace() { return None }

            return Some(&line[..idx]);
        } else if !char.is_whitespace() {
            continue;
        } else {
            return None;
        }
    }
    return None; 
}


struct SpeakerSectionBlockScanner;

impl BlockRule for SpeakerSectionBlockScanner {
    fn run(state: &mut BlockState) -> Option<(Node, usize)> {
        let start_line = state.line;
        let mut next_line = start_line;
        let mut content: Vec<&str> = vec![];

        fn trim_shortname<'a>(sn: &str, line: &'a str) -> &'a str {
            line.clone()[sn.len()+1..].trim_start()
        }

        let shortname = {
            let line = state.get_line(start_line);
            match get_speaker_shortname(line) {
                Some(sn) => {
                    content.push(trim_shortname(sn, line));
                    sn
                },
                None => return None,
            }
        };

        loop {
            next_line += 1;

            if next_line >= state.line_max { break }

            match get_speaker_shortname(state.get_line(next_line)) {
                Some(sn) => {
                    if sn != shortname { break }
                    content.push(trim_shortname(shortname, state.get_line(next_line)));
                },
                None => content.push(state.get_line(next_line)),
            };
        }

        let mut node = state.md.parse(content.join("\n").as_str());
        node.replace(TempSpeakerSection {
            shortname: shortname.to_string(),
        });
        node.srcmap = Some(SourcePos::new(start_line, next_line));

        Some((node, next_line - start_line))
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.block.add_rule::<SpeakerSectionBlockScanner>();
}
