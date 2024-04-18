This project is using the [Goldmark](https://github.com/yuin/goldmark) markdown parser. It's the only parser that supports custom markdown extensions. However, it is a little verbose.

Following Goldmark convesion, each extension is split over multiple directories.

- `extension` contains the entry point for each extension, just a few lines of code that's called from `./cmd/ray-peat-rodeo/main.go`.
- `parser` defines under what conditions an extension get's to execute logic. Goldmark is an Abstract Syntax Tree based parser. This means it scans the input markdown one character at a time from top to bottom, creating nodes, e.g. paragraphs, links, blockquotes. This is where we register our own logic to create our own custom nodes. This part can be quite difficult to visualise, and takes some trial and error.
- `ast` defines the Abstract Syntax Tree nodes that the parser code will create. Each node contains all the data it represents. For example timestamps contain the time.
- `render` defines how our custom nodes render themselves into HTML, e.g. mentions render HTML to create a popup that summarises all other mentions.
- `transformer` has code that operates on the completed AST. Useful for small adjustments to existing nodes otherwise outside our control, e.g. links.
