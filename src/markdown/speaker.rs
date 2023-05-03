use markdown_it::{
    MarkdownIt, Node, NodeValue, Renderer,
    parser::inline::InlineRoot,
    parser::block::{BlockRule, BlockState},
    parser::extset::RootExt,
    plugins::cmark::block::paragraph::Paragraph,
};
use crate::MarkdownFile;

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

/// Retrieves the path for this markdown document
fn path<'a>(state: &'a BlockState<'a, 'a>) -> &'a str {
    state.md.ext.get::<MarkdownFile>().unwrap().path.to_str().unwrap()
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
        } else if char.is_alphanumeric() {
            continue;
        } else {
            return None;
        }
    }
    return None; 
}


#[derive(Debug)]
struct InsideSpeakerSection(String);

impl RootExt for InsideSpeakerSection {}


/// Like a normal ParagraphScanner but strips speaker shortname definition
/// from the start of the paragraph.
struct SpeakerParagraphScanner;

impl BlockRule for SpeakerParagraphScanner {
    fn run(state: &mut BlockState) -> Option<(Node, usize)> {

        let start_line = state.line;
        let mut next_line = start_line;

        loop {
            next_line += 1;

            if next_line >= state.line_max || state.is_empty(next_line) { break; }

            // The logic in this codeblock are lifted directly from the MarkdownIt ParagraphScanner
            // https://github.com/rlidwka/markdown-it.rs/blob/eb5459039685d19cefd0361859422118d08d35d4/src/plugins/cmark/block/paragraph.rs#L43-L59
            {
                if state.line_indent(next_line) >= 4 { continue; }
                if state.line_offsets[next_line].indent_nonspace < 0 { continue; }

                let any_rule_applies = {
                    let old_state_line = state.line;
                    state.line = next_line;
                    let any_rule_applies = state.test_rules_at_line();
                    state.line = old_state_line;
                    any_rule_applies
                };

                if any_rule_applies { break; }
            }
        }

        let section_shortname = match state.root_ext.get::<InsideSpeakerSection>() {
            Some(iss) => &iss.0,
            None => return None,
        };

        let (content, mapping) = {
            let (content, mapping) = state.get_lines(start_line, next_line, state.blk_indent, false);

            let shortname = {
                let Some(shortname) = get_speaker_shortname(content.as_str()) else { return None };

                if shortname != section_shortname {
                    panic!("Speaker shortname {shortname} found within speaker section {section_shortname} in {}", path(state));
                }

                shortname
            };

            let new_content = content[shortname.len()+1..].trim_start().to_string();
            let reduction = content.len() - new_content.len();

            let mut new_mapping: Vec<(usize, usize)> = vec![];
            for map in mapping {
                new_mapping.push((
                    if map.0 <= reduction { 0 } else { map.0 - reduction },
                    map.1,
                ));
            }
            (new_content, new_mapping)
        };

        let mut node = Node::new(Paragraph);
        node.children.push(Node::new(InlineRoot::new(content, mapping)));
        Some((node, next_line - start_line))
    }
}




/// Rule that handles paragraphs beginning with a speaker shortname
/// e.g. "RP: Ray peat says this or that"
/// 
/// It consumes every line until a different shortnames is defined.
/// All lines are tokenized as blocks and made it's children.
struct SpeakerSectionBlockScanner;

impl BlockRule for SpeakerSectionBlockScanner {
    fn run(state: &mut BlockState) -> Option<(Node, usize)> {
        match state.root_ext.get::<InsideSpeakerSection>() {
            Some(_) => return None,
            None => (),
        };

        let start_line = state.line;
        let mut next_line = start_line;

        let shortname = {
            let Some(shortname) = get_speaker_shortname(state.get_line(start_line)) else { return None };
            shortname.to_string()
        };

        let speakers = state.md.ext.get::<MarkdownFile>().unwrap().frontmatter.speakers.clone();
        
        let longname = speakers.get(&shortname.to_string())
            .expect(format!("Speaker shortname \"{shortname}\" not found in \"speakers\" in YAML frontmatter in {}", path(state)).as_str())
            .to_string();

        loop {
            next_line += 1;
            if next_line >= state.line_max { break }
            let Some(subsequent_shortname) = get_speaker_shortname(state.get_line(next_line)) else { continue; };
            if subsequent_shortname != shortname { break }
        }

        let speaker_section = {
            let original_node = std::mem::replace(&mut state.node, Node::new(SpeakerSection {
                shortname: shortname.to_string(), longname,
            }));

            let original_line_max = state.line_max;
            state.line_max = next_line;
            state.line = start_line;

            state.root_ext.insert(InsideSpeakerSection(shortname));
            state.md.block.tokenize(state);
            state.root_ext.remove::<InsideSpeakerSection>();

            state.line = start_line;
            state.line_max = original_line_max;

            std::mem::replace(&mut state.node, original_node)
        };

        let consumed = next_line - start_line;

        Some((speaker_section, consumed))
    }
}

/// Looks for paragraphs beginning with a speaker definition: Any alphanumeric characters followed
/// by a colon and then some whitespace as the first text in paragraph. For example...
///
/// RP: Hello, my name is Ray Peat.
///
/// A SpeakerSection node is created that consumes this line and all subsequent lines until a
/// different speaker definition is made.
///
/// For a particular SpeakerSection all lines are block tokenized, which is to say parsed as
/// normal. And speaker definitions are removed from the beginning of paragraphs within a speaker
/// section.
///
/// Each SpeakerSection translates the speaker definition from it's shortname, to the full speaker
/// name defined in each document's frontmatter, which must be passed to this parser like so...
///
/// markdown_parser.ext.insert(crate::markdown::Speakers(frontmatter.speakers));
pub fn add(md: &mut MarkdownIt) {
    md.block.add_rule::<SpeakerSectionBlockScanner>();
    md.block.add_rule::<SpeakerParagraphScanner>();
}
