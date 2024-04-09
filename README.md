[raypeat.rodeo](https://raypeat.rodeo) is the open-source effort to transcribe
the public works of Ray Peat.

[![Deploy to GitHub Pages](https://github.com/marcuswhybrow/ray-peat-rodeo/actions/workflows/gh-pages.yml/badge.svg)](https://github.com/marcuswhybrow/ray-peat-rodeo/actions/workflows/gh-pages.yml)

![banner](https://raw.githubusercontent.com/marcuswhybrow/ray-peat-rodeo/back-to-go/internal/assets/docs/ray-peat-rodeo-banner.png)

# Developing

Requirements: [Nix Package Manager](https://nixos.org/download.html#download-nix)

```bash
git clone git@github.com:marcuswhybrow/ray-peat-rodeo.git
cd ray-peat-rodeo

nix develop -c modd # Starts auto-reloading dev server
```

In a second terminal...

```bash
nix develop -c tailwind-watch # Allows tailwind to rebuild CSS classes
```

In a third terminal edit a transcription file... 

```bash
vim ./assets/YYYY-MM-DD-title-of-new-transcription.md
```

# AI Transcription

This project aims to create transcripts of Ray Peat's interviews in order to 
then augment them with metadata and annotations, this is slow.

A similar project exists called [The Ray Peat Archive](https://github.com/0x2447196/raypeatarchive) 
which is much faster. It uses AI to transcribe the audio, and it's fast.
RPA has the far superior quantity of transcription.

This project will catch up by copying their approach. `./assets/todo` contains
every transcript-less interview for which a URL is known. What follows is an 
example of my AI workflow to grab the audio from that URL, use AI to transcribe 
it, and wrangle it into custom markdown.

```bash
# Inside project directory, launch shell with necessary tools in environment
nix develop 

# Pick a markdown file from ./assets/todo and copy the source URL.
# This command downloads that URL as an audio file called "source-audio"
yt-dlp -x "https://website.com/some-video-or-audio-file-url" -o source-audio

# Next we get OpenAI's Whisper to transcribe the audio file into a JSON file.
# This takes a while, and creates "source-audio.json"
whisper --language English --output_format json source-audio

# Finally, transform the resulting JSON into markdown.
# This command outputs markdown to stdout, so here we redirect stdout to append 
# to whichever markdown file we took the original audio URL from
whisper-json2md source-audio.json >> ./assets/todo/2022-08-19-example.md
```
