# netsuite-docs

> Scraper that brings Oracle Netsuite documentation into an LLM-friendly format.

## Usage

```sh
go run .
```

The scraper outputs the following files:

- `out.json` - a list of json objects containing content and metadata for each page of documentation.
- `volumes/volume-*.md` - chunks of documentation pages combined into various volumes in the markdown format, useful for use with Google's Notebooklm.

