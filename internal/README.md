`internal`, is the [official convention](https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages) for code that's reusable within this project, but not outside of it. In our case, this code does the heavy lifting for the user-facing commands in [`./cmd`](https://github.com/marcuswhybrow/ray-peat-rodeo/tree/main/cmd).

`./assets/` is a dumb directory of files that get's directly copied to `./build/assets/`. Once `./build` is deployed as `https://raypeat.rodeo` our assets are accessible at `https://raypeat.rodeo/assets/`. So we can drop any images, or CSS styles we may need in here.

`./blog` is a collection of markdown files, naming scheme `YYYY-MM-DD-slug.md`. `./cmd/ray-peat-rodeo/` converts each one to HTML and outputs it to `./build/blog/`, along with a list of all blog posts at `./build/blog/index.html`

`./cache` is a wrapper around Go's `net/http` module that caches data fetched over HTTP to a YAML file that can be commited to source control. Take the example of a markdown file in `./assets` that contains the text `[[https://raypeat.com]]`. The double square brackets are a custom markdown extension marking this URL as a "mention". 

The point of a mention is to build a global index of all mentions, so the URL is a great unique identifier for a web page. Ray Peat Rodeo automatically scrapes the title of mentioned URLs to display to readers as well. 

Scraping the title of every mentioned URL on every build would slow the dev server to a crawl, and make GitHub deployments non-deterministic. i.e. if a website returns a valid title from my home internet during development, it may be block requests from GitHub's servers when they build the same code for deployment to GitHub pages.

The solution is to cache all scraped data to a file during development, commit it into source control, then use that (deterministic/static) data during the build.

`./global` contains a few centralised constants for use by all other modules.

`./markdown/` contains custom extensions such as the aforementioned `[[Mentions]]`, as well as `[12:23]` timecodes, `{sidenotes}`, `{#12}` to reference GitHub issues as a sidenote, and `RP: Hi \n\n MW: Hello` for speaker demarcation. See [`./markdown/README.md`](https://github.com/marcuswhybrow/ray-peat-rodeo/tree/main/internal/markdown) for more detail.

`./http_cache.yml` is the aforementioned source controlled YAML cache of scraped HTTP data.
