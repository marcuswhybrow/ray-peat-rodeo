use markdown_it::{MarkdownIt, Node, NodeValue, Renderer};
use markdown_it::parser::inline::{InlineRule, InlineState};
use crate::MarkdownFile;
use maud::html;

const LINK_SVG: maud::PreEscaped<&str> = maud::PreEscaped(r#"
    <svg width="16" height="16" version="1.1" viewBox="0 0 383.028 383.027">
        <path d="M361.213,244.172l-71.073-71.073c-16.042-16.042-37.632-23.216-58.648-21.562c1.653-21.019-5.521-42.609-21.563-58.651
            l-71.073-71.073c-29.084-29.084-76.408-29.083-105.492,0L21.814,33.361c-29.084,29.084-29.084,76.408,0,105.493l71.073,71.073
            c16.042,16.042,37.632,23.217,58.651,21.563c-1.654,21.02,5.52,42.607,21.563,58.65l71.073,71.073
            c29.084,29.084,76.408,29.083,105.492,0l11.548-11.548C390.297,320.58,390.297,273.256,361.213,244.172z M136.174,161.292
            l29.458,29.458c-14.997,8.932-34.734,6.955-47.629-5.94l-71.073-71.073c-15.233-15.234-15.233-40.022,0-55.258l11.549-11.548
            c15.235-15.235,40.023-15.235,55.259,0l71.072,71.073c12.896,12.895,14.873,32.632,5.94,47.63l-29.458-29.458
            c-6.937-6.937-18.181-6.937-25.117,0S129.238,154.354,136.174,161.292z M336.095,324.547l-11.548,11.548
            c-15.234,15.235-40.022,15.234-55.258,0l-71.073-71.073c-12.895-12.895-14.873-32.632-5.938-47.629l29.458,29.458
            c6.936,6.938,18.181,6.938,25.116,0c6.937-6.937,6.938-18.181,0-25.115l-29.458-29.459c14.998-8.934,34.735-6.956,47.631,5.939
            l71.072,71.073C351.331,284.523,351.331,309.312,336.095,324.547z"/>
    </svg>
"#);

#[derive(Debug)]
pub struct InlineTimecode {
    pub url: url::Url,
    pub hours: u8,
    pub minutes: u8,
    pub seconds: u8,
}

impl NodeValue for InlineTimecode {
    fn render(&self, node: &Node, fmt: &mut dyn Renderer) {
        let timecode = {
            let timecode = format!("{:0>2}:{:0>2}", self.minutes, self.seconds);

            if self.hours > 0 {
                format!("{:0>2}:{timecode}", self.hours)
            } else {
                timecode
            }
        };

        let href = format!("{}#t={}", self.url, 'out: {
            let Some(host) = self.url.host() else { break 'out timecode.clone() };
            let host = host.to_string();

            if host.ends_with("youtube.com") || host.ends_with("youtu.be") {
                format!("{:0>2}h{:0>2}m{:0>2}s", self.hours, self.minutes, self.seconds)
            } else {
                timecode.clone()
            }
        });

        fmt.open("span", &{
            let mut attrs = node.attrs.clone();
            attrs.push(("class", "timecode".into()));
            attrs
        });

        fmt.text_raw(html! {
            a.internal href={ "#t=" (timecode) } { (LINK_SVG) }
            a.external id={ "t=" (timecode) } target="_blank" href=(href) { (timecode) }
        }.into_string().as_str());

        fmt.close("span");
    }
}

struct TimecodeInlineScanner;

impl InlineRule for TimecodeInlineScanner {
    const MARKER: char = '[';

    // [00:00] or [00:00:00] with variable length digits in each section
    fn run(state: &mut InlineState) -> Option<(Node, usize)> {
        let mut sections: Vec<Vec<char>> = Vec::from([Vec::new()]);

        let mut chars = state.src[state.pos..state.pos_max].chars();

        if chars.next().unwrap() != '[' { return None; }

        for (i, char) in chars.enumerate() {
            if char.is_digit(10) {
                // [00:00:00]...
                // -^^-^^-^^-
                let last = sections.len() - 1;
                sections[last].push(char.clone());
            } else if char == ':' {
                // [00:00:00]...
                // ---^--^---
                if sections.len() >= 3 { return None }
                sections.push(Vec::new());
            } else if char == ']' {
                // [00:00:00]...
                // ---------^
                if sections.len() < 2 || sections.len() > 3 { return None };
                return Some((
                    Node::new(InlineTimecode {
                        url: url::Url::parse(
                            state.md.ext.get::<MarkdownFile>().unwrap().frontmatter.source.as_str()
                        ).unwrap(),
                        seconds: sections.pop().unwrap_or(Vec::new()).iter().collect::<String>().parse().unwrap_or(0),
                        minutes: sections.pop().unwrap_or(Vec::new()).iter().collect::<String>().parse().unwrap_or(0),
                        hours: sections.pop().unwrap_or(Vec::new()).iter().collect::<String>().parse().unwrap_or(0),
                    }),
                    i + 2,
                ));
            } else {
                return None;
            }
        }

        return None;
    }
}

pub fn add(md: &mut MarkdownIt) {
    md.inline.add_rule::<TimecodeInlineScanner>();
}
