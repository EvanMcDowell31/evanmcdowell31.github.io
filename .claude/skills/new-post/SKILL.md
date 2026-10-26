---
name: new-post
description: Draft a new blog post for Evan's Hugo site (evanmcdowell31.github.io). Use when asked to write, start, draft or add a blog post, article or note, or to turn an idea, LinkedIn post or outline into a post.
---

# New blog post

Posts live in `content/posts/` as Markdown with TOML front matter (`+++`).

## Steps

1. Pick a short kebab-case slug from the topic (e.g. `agent-skills-at-scale`). Check `content/posts/` so you don't reuse one.
2. Create the file from the archetype so the front matter is consistent:
   ```sh
   hugo new content posts/<slug>.md
   ```
   This sets `title`, `date` and `draft = true`.
3. Fill in the front matter:
   - `title`: sentence case, no trailing period.
   - `summary`: one sentence, under about 160 characters. It's used as the meta description and in RSS.
   - Leave `draft = true` until Evan says it's ready to publish.
4. Write the body in Markdown. Start with `##` headings, because the page already renders `title` as the `h1`. Fenced code blocks need a language tag for syntax highlighting.
5. Images go in a page bundle: turn the post into `content/posts/<slug>/index.md` and put the images beside it, then reference them relatively (`![alt](diagram.png)`). Always write meaningful alt text.
6. Preview with drafts included and check the post renders:
   ```sh
   hugo server -D
   ```
   Open `http://localhost:1313/posts/<slug>/`.

## Voice

Write the way Evan's LinkedIn writing reads: first person, direct, short paragraphs, concrete numbers over adjectives. Avoid hype words and filler intros.

## Publishing

Change `draft = false` only when Evan asks, then use the `publish-site` skill. Don't commit or push from this skill.

## Privacy

This site is public. When adapting work content (Walmart, internal tools, teammates), don't add internal names, metrics or people beyond what Evan has already made public on LinkedIn, and point out anything borderline before publishing.
