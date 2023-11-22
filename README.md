[raypeat.rodeo](https://raypeat.rodeo) is the open-source effort to transcribe
the public works of Ray Peat.

[![Deploy to GitHub Pages](https://github.com/marcuswhybrow/ray-peat-rodeo/actions/workflows/gh-pages.yml/badge.svg)](https://github.com/marcuswhybrow/ray-peat-rodeo/actions/workflows/gh-pages.yml)

![banner](https://raw.githubusercontent.com/marcuswhybrow/ray-peat-rodeo/back-to-go/internal/assets/docs/ray-peat-rodeo-banner.png)

# Usage

Requirements: [Nix Package Manager](https://nixos.org/download.html#download-nix)

```bash
git clone git@github.com:marcuswhybrow/ray-peat-rodeo.git
cd ray-peat-rodeo
nix develop -c watch-and-serve

# In another terminal add a new transcript.
# See existing files in ./content for examples
touch ./content/YYYY-MM-DD-title-of-new-transcription.md
```
