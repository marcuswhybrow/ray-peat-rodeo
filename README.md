# https://raypeat.rodeo

My effort to catalogue, compile, and transcribe the public works, speeches and interviews of Ray Peat.  
[Open an issue](https://github.com/marcuswhybrow/ray-peat-rodeo/issues) if there's a Ray Peat interview I'm missing.

This repository represents a collection of markdown transcripts built by the static-site generator [11ty](https://www.11ty.dev/).

## Installation

```
npm install
```

## Usage

```
npm start
```

## Interview Syntax

Markdown files in `./src/content/` have additional bespoke template tag shorthands for defining who's speaking, and identifying the people, books, and URLs mentioned by the speakers.

- **Interviewer** - Lines prefixed with `! ` (such as `! Good morning Ray, how are you?`) declare paragraphs spoken by the interviewer. (Support for multiple interviewers forthcoming).
- **Ray Peat** - Lines without the `! ` prefix declare paragraphs spoken by Ray Peat (i.e. `Very well, thank you.`).

People, books, and URLs should be wrapped in double square brackes (`[[Text]]`) as below. Doing so feeds these links into Ray Peat Rodeo's site-wide index.

- *People* - Link to people by surrounding their full name with double square brackets `[[William Blake]]`. Override the display text with the pipe suffix `[[William Blake|an important poet]]`. Or link invisibly by defining an empty display text `[[William Blake|]]`.
- *Books* - Link to books by their 13 (or 10) digit ISBN number ``[[9780385152136|The Complete Poetry & Prose of William Blake]]``. The display text is required. Invisible linking supported.
- *URLs* - Link to external URLs likewise ``[[https://www.youtube.com/watch?v=lDr71LHO0Jo]]``. Display text is supported. Invisible linking is supported.