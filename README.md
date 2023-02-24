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

- *People* - Link to people by surrounding their full name with double square brackets `[[William Blake]]`.
- *Books* - Link to books with the title and the primary author's full name ``[[Jerusalem -by- William Blake]]``. The `-by-` separator, and exactly one author is required. Display text (see below) defaults to Book title.
- *URLs* - Link to external URLs ``[[https://www.youtube.com/watch?v=lDr71LHO0Jo]]``.
- *DOIs* - Link to scientific papers by their DOI ``[[doi:10.5860/choice.37-5129]]`` The `doi:` prefix is required.

All `[[Links]]` may optionally override the display text with the pipe sufffix `[[William Blake|a poet]]`. Hidden links (that produce no markup) are created with an empty display text string `[[William Blake|]]`. Missing links can be declared by omitting everything before the pipe `[[|text that will eventually link to something]]`.

## Style Guide

- **Em Dashes** - Long dashses, or em dashes, (Windows ALT code 0151) when used for parenthesis contain no spaces `While I was shopping—wandering aimlessly up and down the aisles, actually—I ran into our old neighbor.