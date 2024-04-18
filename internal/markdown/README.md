This project is using the [Goldmark](https://github.com/yuin/goldmark) markdown parser. I
t's the only parser that supports custom markdown extensions. However, it is a little verbose.

Following Goldmark convesion, each extension is split over multiple directories.

- `extension` contains the entry point for each extension, just a few lines of code that's called from `./cmd/ray-peat-rodeo/main.go`.
- `parser` defines under what conditions an extension get's to execute logic. Goldmark is an Abstract Syntax Tree based parser. This means it scans the input markdown one character at a time from top to bottom, creating nodes with nodes. This is where we register our own logic to create our own nodes.
- `ast` defines the Abstract Syntax Tree nodes that the parser code will create.
- `render` defines how our custom nodes render themselves into HTML.
- `transformer` has code that operates on the completed AST. Useful for small adjustments to existing nodes otherwise outside our control, e.g. links.
