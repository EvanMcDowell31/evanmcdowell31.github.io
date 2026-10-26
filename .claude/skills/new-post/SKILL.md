---
name: new-post
description: Draft a new Story post (the blog tab) for Evan's Hugo site (evanmcdowell31.github.io). Use when asked to write, start, draft or add a story, blog post, article or note, or to turn an idea, LinkedIn post or outline into one.
---

# New story post

The Story tab is the blog. Posts live in `content/story/` as Markdown with TOML front matter (`+++`). The newest one shows as the featured "Latest chapter" card, and the rest are grouped by year.

## Steps

1. Pick a short kebab-case slug from the topic (e.g. `agent-skills-at-scale`). Check `content/story/` so you don't reuse one.
2. Create the file from the archetype so the front matter is consistent:
   ```sh
   hugo new content story/<slug>.md
   ```
   This sets `title`, `date` and `draft = true`.
3. Fill in the front matter:
   - `title`: sentence case, no trailing period.
   - `summary`: one sentence, under about 160 characters. It's used as the meta description and in RSS.
   - `tags`: one to three lowercase tags, reusing existing ones where possible (`ls public/tags` after a build). Tags power the "wander by tag" chips on Search.
   - Leave `draft = true` until Evan says it's ready to publish.
4. Write the body in Markdown. Start with `##` headings, because the page already renders `title` as the `h1`. Fenced code blocks need a language tag for syntax highlighting.
5. Images go in a page bundle: turn the post into `content/story/<slug>/index.md` and put the images beside it, then reference them relatively (`![alt](diagram.png)`). Always write meaningful alt text.
6. Preview with drafts included and check the post renders:
   ```sh
   hugo server -D
   ```
   Open `http://localhost:1313/story/<slug>/`. The first paragraph gets a drop cap, so open with a real sentence, not a heading or an image.

## Voice

Write the way Evan's LinkedIn writing reads: first person, direct, short paragraphs, concrete numbers over adjectives. Avoid hype words and filler intros.

## Publishing

Change `draft = false` only when Evan asks, then use the `publish-site` skill. Don't commit or push from this skill.

## Privacy

This site is public. When adapting work content (Walmart, internal tools, teammates), don't add internal names, metrics or people beyond what Evan has already made public on LinkedIn, and point out anything borderline before publishing.
