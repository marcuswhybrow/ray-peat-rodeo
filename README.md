[raypeat.rodeo](https://raypeat.rodeo) is the open-source effort to transcribe
the public works of Ray Peat.

[![Deploy to GitHub Pages](https://github.com/marcuswhybrow/ray-peat-rodeo/actions/workflows/gh-pages.yml/badge.svg)](https://github.com/marcuswhybrow/ray-peat-rodeo/actions/workflows/gh-pages.yml)
[![Built with Nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

![banner](https://raw.githubusercontent.com/marcuswhybrow/ray-peat-rodeo/back-to-go/internal/assets/docs/ray-peat-rodeo-banner.png)

# Getting Started

Get [Nix Package Manager](https://nixos.org/download.html#download-nix), then 
clone this repository and start the auto-reloading dev server:

```bash
git clone git@github.com:marcuswhybrow/ray-peat-rodeo.git
cd ray-peat-rodeo
nix develop -c modd
```

- `./assets` contains a markdown file for each article or transcription.
- `./assets/todo` contains a markdown file for assets without transcriptions.
- `./github/workflows/gh-pages` auto deploys this repo to https://raypeat.rodeo
- `./cmd/ray-peat-rodeo` is the code that builds the website from the markdown 
assets.
- `./cmd/whisper-json2md` is a custom tool to massage AI transcripts to 
markdown (see [here](#ai-transcription)).
- `./internal` contains this projects unique features, especially 
`./internal/markdown` and `./internal/cache/` which provide custom markdown 
plugins and automatic caching of remote data (e.g. GitHub issue titles).

Those are the main things, and editing any of them trigger the dev server to 
auto reload, and your browser should hot-load the changes and auto refresh, 
making a near instant dev cycle.

- `./flake.nix` & `./flake.lock` tell the `nix` command how to do everything 
for us, for example the dev server (lauched by `nix develop -c modd`), tells
nix to examine `./flake.nix` enter the custom shell environment defined 
there, and run the command `modd` which is our dev server of choice.
- `./gomod2nix.toml` in conjunction with `./flake.nix` helps the `nix` command 
build this project. It's autogenerated by running `nix develop -c gomod2nix`.
- `./modd.conf` tells `modd` how to behave, such as running Tailwind CCS 
process automatically.
- `./tailwind.config.js` tells tailwind how to do it's thing.

And finally, you may wish to use [direnv](https://direnv.net) and 
[nix-direnv](https://github.com/nix-community/nix-direnv) to automatically load
all project dependencies and tools into your shell environment whilst you are 
inside the project directory (auto unloads when you leave it). In the project 
directory:

```bash
direnv allow
```

# Project Goals

1. Round up every Ray Peat interview, article, newsletter and book.
2. Use AI to quickly transcribe interviews.
3. Store each interview (etc.) as human readable markdown.
4. Generate a website from those markdown source files.
5. Site-wide search of all assets.
6. Tooltips for all mentioned topics and people linking to all other mentions.
7. Timestamps linking to specific times in original audio or video.
8. Sidenote annotations for clarifications and issues to be resolved.

What this amounts to is using AI to quickly transcribe all interviews, then 
storing the results in markdown. Next, one improves and augments each markdown 
file with corrections, formatting and tagging all mentions and timecodes. 

Formatting is part of the markdown standard, but what I'm calling "mentions", 
"timecodes", and "issues" are extensions to the markdown syntax written 
specifically for this project. With custom markdown syntax any functionality
can be realised whilst keeping the markdown documents human readible for 
archival purposes, and portability to other projects.

# Adding "Todo" Ray Peat Interviews

Go to `./assets/todo`. Every file in this directory is a 
[Markdown](https://www.markdownguide.org) file. Each one repesenting a unique
Ray Peat interview awaiting transcription. Each filename is formatted in 
[Kebab Case](https://developer.mozilla.org/en-US/docs/Glossary/Kebab_case) and 
begins with the date as `YYYY-MM-DD` (ISO 8601 format) followed by the title of 
the interview.

Ray Peat Rodeo will respect whatever date is declared in the file name, and use 
it across the website. The title portion, verbatim, becomes the URL at which 
this interview will exist.

For example...

```bash
touch ./assets/todo/2008-07-02-an-example.md
```

... will become a web page accessible at `raypeat.rodeo/an-example/` and will 
appear in the 2008 section, as having taken place on July 2nd. The contents of 
the file must begin with the following [YAML](https://yaml.org/) frontmatter.

```markdown 
---
source:
    series: The name of the show Ray is appearing on 
    title: Human readable title (similar to filename title but more flexible)
    url: https://example.com/the-original-audio-or-video
    kind: audio
---
```

- The `series` is used to group interviews by the show/host, so make sure you 
match the series _exactly_ to existing series in other interviews.
- The `title` can contain any characters and appears at the top of the 
interview page, and on the homepage listing. 
- The `url` is used to link to the original source URL, and to constuct 
"timestamp" links to allow readers to click through from a given point in the 
interview diectly to that time in the source audio or video. 
- `kind` can be either `audio` or `video` and is offered as a filter when 
searching Ray Peat Rodeo.

**Done**. Next one may use the `transcribe` tool to automatically add an AI 
transcription to this file (see [AI Transcription](#ai-transcription)).

# AI Transcription

`flake.nix` packages a `bash` script named `transcribe`. It downloads the 
source audio of any file in `./assets/todo`, transcribes it, then updates the 
asset with the transcription, and updates the frontmatter data to reflect this 
change.

1. Argument **#1** is the markdown file to transcribe and update.
2. Argument **#2** is your name, to log in the assets metadata.

```bash
nix run github:marcuswhybrow/ray-peat-rodeo#transcribe -- ./assets/todo/2024-10-12-example.md "Marcus Whybrow"
```

**Done**. Once you've added the AI transcript it's contents will be available to 
the site-wide search engine, helping readers to further explore Ray's ideas. 
Finally, and optionally, one may augment the transcript with special formatting 
to take it to the next level (see 
[Augmenting and Completing A Transcript](#augmenting-and-completing-a-transcript)).

# Augmenting and Completing A Transcript

## Who's Speaking?

Prefixing sentences with the speakers initials when the speaker changes, such 
as `RP:` for Ray Peat, allows Ray Peat Rodeo to separate the transcript into 
different speach bubbles.

Make sure to define the full name for each initials used in the YAML 
frontmatter at the top of the markdown file like so:

```markdown
---
speakers:
    RP: Ray Peat
    MW: Marcus Whybrow
---

MW: Hi Ray, how are you?

RP: Very good, thank you.
```

## What Time Is It?

Interspersing timestamps within the transcript, such as `[12:34]`, allows 
readers to jump staight to that point in the original source audio or video. 
I like to use timestamps sparingly to indicate a change in topic or a new 
question being asked. For example...

```markdown 
MW: That's great. [12:34] And what do you think about that, Ray?

RP: I think...
```

*Tip: Timestamps can express hours too: `[2:01:12]`*

## Mentions

When a person, topic, chemical, hormone, book, website, or any *thing* is 
mentioned, marking it as a "mention" gives readers a little popup bubble that 
provides a mini summary of where else it's been discussed. Surround the 
mentioned thing in double square bracets like this...

```markdown 
RP: The history of [[Estrogen]] reesearch...
```

For mentioned people, put their surname first, then a comma, then their given 
names (without titles such as Sir or Doctor). For example...

```markdown
MW: [[Wodehouse, Pelham Grenville]] was the creater of Jeeves and Wooster...
```

This backwards convension helps Ray Peat Rodeo know how to order every mention 
alphabetically. RPR is smart enough to output the name the right way around to 
the reader...

> **Pelham Grenville Wodehouse** was the creator of Jeeves and Wooster...

*Note: The first comma always has this effect. Commas must be otherwise 
avoided in mention names.*

To tailor the displayed text to your liking use the `|` character...

```markdown 
MW: I've been reading [[Blake, William|an author]] that...
```

Which becomes...

> I've been reading **an author** that...

And finally, you can associate books with their authors using the `>` 
character. For example...

```markdown
MW: and I discovered he wrote [[Blake, William > Jerusalem]] around then...
```

`Jerusalem` is known as a "sub mention", and it'll be included included in the 
popup summary for William Blake, and *vise versa*. Sub mentions are a powereful 
way to help new readers explore Ray's influences by hopping around these 
associations bound together via unique conversations. 

*Tip: A mention or submention may be a URL or email address. In these speacial 
cases, the popup summary will also contain a direct link to the URL, or a 
"mailto" link to open the reader's email client directly.*

*Tip: Mentioning a scientific paper by it's DOI URL (https://doi.org/...) 
automatically grabs the papers full title from the DOI database to display to 
the reader. See this [real example](https://raypeat.rodeo/john-william-gofman/#https%3A%2F%2Fdoi.org%2F10.5860%2Fchoice.37-5129).*

## Is That Clear?

When someone new to Ray Peat may not understand a reference or term, one can 
add a sidenote, using curly brackets, that appears distinct from the main text
in a little bubble. For example...

```markdown 
RP: PUFA {Polyunsaturated Fats} were originally...
``` 

I like to clarify a term this way the first time it's used in a transcript, 
then trust the reader to recall it's definition, or refer back to it. 
This serves to keep interruptions to a minimum and let Ray take center stage.

## Huh?

Sometimes Ray's mentions are ambiguous, or the full name of a paper or 
person is unclear. In this case one may  
[create an issue](https://github.com/marcuswhybrow/ray-peat-rodeo/issues/new)
in the GitHub project and title it as a question to which others may know the 
answer. For example "Which 1986 biology paper is Ray refering to?" Add to the 
issue's description any pertenant context and submit the issue.

Once created, take note of the issue's unique numerical ID displayed near the 
issue's title. Refer to this ID using a `#` inside of a sidenote:

```markdown
RP: In 1986 they showed {#51} that even though...
```

When a sidenote contains a `#` and a number, a golden, call to action, issue 
bubble containing the issue title will be shown to readers. In this case the 
bubble will read "#51 Which 1986 biology paper is Ray refering to?" Clicking 
the bubble takes readers to the GitHub issue itself.

GitHub issues are a great way to keep track of opportunities for improving
the clarity of readability for new readers, and serve to invite and organise 
the expertese of those who might fill in the gaps.

I like to use issues liberally. If I'm unsure of a mention, or don't know how 
to word a sidenote, I create an issue and move on. This keeps transcription 
fluent, leaving future me, or someone better educated, to fix the issue later.

# Similar Projects 

- [The Ray Peat Archive](https://github.com/0x2447196/raypeatarchive) has a very large collection of AI transcriptions stored in a plain text subtitle format callled WebVTT.
- [Bioenergetic Life](https://bioenergetic.life/) is an interactive search engine for The Ray Peat Archive's data set, with side-by-side text and audio snippets.

[Open an issue](https://github.com/marcuswhybrow/ray-peat-rodeo/issues/new) if I've missed a similar project.
