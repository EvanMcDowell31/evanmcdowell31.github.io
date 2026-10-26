---
name: update-profile
description: Update Evan's bio, experience, education, certifications or project cards on his Hugo site, optionally syncing from his LinkedIn profile with Playwright. Use when asked to update the about page, resume, job history, projects, intro, or to "sync with LinkedIn".
---

# Update profile content

| What | File |
|---|---|
| Home intro | `content/_index.md` |
| About, experience, education, certifications | `content/about.md` |
| Project cards on the home page | `data/projects.toml` (`[[project]]` with `name`, `description`, optional `url`, `tags`) |
| Site description, GitHub/LinkedIn links | `hugo.toml` `[params]` |

For a manual edit, change the right file, run `hugo --gc --minify` to make sure it builds, and check the page with `hugo server`.

## Syncing from LinkedIn

LinkedIn requires sign-in, so this uses Playwright with a saved session.

1. **Workspace:** use a scratchpad directory, never the repo, because the session file is a live login. Install Playwright there if it isn't already:
   ```sh
   cd <scratch>/pw && npm init -y && npm install playwright@1.56.1 && npx playwright install chromium
   ```
   Then copy `scripts/*.mjs` from this skill into `<scratch>/pw` so the scripts resolve the `playwright` package.
2. **Sign in:** run in the background, and tell Evan a browser window is waiting for him:
   ```sh
   node li-login.mjs <scratch>/pw/li-state.json
   ```
   If the session file still works, it exits without prompting.
3. **Read the profile:**
   ```sh
   node li-sections.mjs <scratch>/pw/li-state.json <scratch>/pw/li
   ```
   Exit code 3 means the session expired, so go back to step 2. The `.txt` files include LinkedIn page furniture (footer, language list, "people you may know"). Use only the profile content.
4. **Compare with the site:** compare the scraped data with `content/about.md` and `data/projects.toml`. Show Evan the differences (new roles, changed dates or titles, new certifications) before editing.
5. **Apply** the approved changes, keeping the existing format. Keep the bullets tight and outcome-focused, and rewrite LinkedIn phrasing instead of pasting it.

## Privacy rules

The site is public. Never copy:
- profile analytics, profile viewers, connection counts, or "open to work" settings
- names of teammates or people from posts, unless Evan asks
- internal details beyond what's already in his public LinkedIn experience text. Flag anything that looks confidential.

## Cleanup

When finished, offer to delete `li-state.json`. Never commit it or copy it into the repo. `.gitignore` blocks `*-state.json` as a backstop.
