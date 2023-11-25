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
