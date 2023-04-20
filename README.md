[raypeat.rodeo](https://raypeat.rodeo) is the open-source effort to transcribe
the public works of Ray Peat.

You are looking at the `rust` git branch: an in progress reimplementation of
the Ray Peat Rodeo golang codebase in Rust.

# Reimplement

- [x] modd/devd development environment
- [ ] HTML Templating
- [ ] Parse transcripts from Markdown to HTML
  - [ ] Standard Markdown parsing
  - [ ] Block interview syntax 
  - [ ] Inline citations
    - [ ] People
    - [ ] Books
    - [ ] ISBN
    - [ ] External Links
  - [ ] Citation metadata scraping
  - [ ] Timecodes
- [ ] Pagefind static search

# Usage

```bash
nix develop -c serve   # build and serve site over HTTP (auto-reloads)
```
