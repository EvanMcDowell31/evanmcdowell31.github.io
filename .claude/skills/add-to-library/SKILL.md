---
name: add-to-library
description: Add a document, bookmark or script to the Library tab of Evan's Hugo site. Use when asked to save, bookmark, add a link, add a snippet or script, file a doc, or "put this in the library".
---

# Add to the library

The Library has three shelves, each a Hugo section under `content/library/`. Every item is one Markdown file with TOML front matter. Search picks it up automatically from `search-index.json`.

| Shelf | Folder | Shows as | Opens to |
|---|---|---|---|
| Documents | `documents/` | a title and summary row linking to its own page | the full page |
| Bookmarks | `bookmarks/` | a title and domain row | the summary, tags and a Visit button. Search results go straight to the link. |
| Scripts | `scripts/` | a title and summary row | the code inline, with a copy button |

## Front matter

```toml
+++
title = "Short, specific title"
date = 2026-10-26          # today
summary = "One line, under about 120 characters. It's always visible, so make it useful."
tags = ["lowercase", "reuse-existing"]
# bookmarks only:
[params]
  link = "https://example.com/page"
+++
```

## Per shelf

- **Bookmark:** the body can stay empty. Write the summary in Evan's voice, saying why it's worth keeping rather than what it is. Check that the URL returns 200 with `curl -sIL`.
- **Script:** the body is one fenced code block with a language tag, plus a one-line note on what it prints or any gotcha. Keep it copy-pasteable: no prompts, and no secrets or internal hostnames.
- **Document:** a normal Markdown page that starts with `##` headings. Documents don't get the drop cap.

## Files and checks

- Name files with a kebab-case slug of the title, e.g. `bookmarks/react-docs.md`.
- Run `hugo --gc --minify` and confirm the item appears in `public/search-index.json`. Preview at `http://localhost:1313/library/<shelf>/`.

## Privacy

Public site: no Walmart-internal links (`*.walmart.com` internal tools, Confluence, Jira, GitHub Enterprise), no tokens, and nothing from private notes unless Evan says it's fine to share.

Don't commit or push from this skill. Use `publish-site`.
