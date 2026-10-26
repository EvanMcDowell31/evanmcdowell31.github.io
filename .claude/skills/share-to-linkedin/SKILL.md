---
name: share-to-linkedin
description: Share a published Story post from Evan's Hugo site to LinkedIn via the official API (Go tool + manual GitHub Action), or set up/renew the LinkedIn connection. Use for "post this to LinkedIn", "share the story", "LinkedIn token expired", "reconnect LinkedIn".
---

# Share a story to LinkedIn

To share a story, run the **Share a story on LinkedIn** workflow (`.github/workflows/linkedin.yml`). It builds the site, uploads the post's share card (`public/og/story-<slug>.png`), creates a public article post as Evan, writes the post's URL into the story's front matter as `[params] linkedin`, and redeploys. That URL makes the "Originally posted on LinkedIn" box appear on the page.

## Share a post

1. The story must already be live (`draft = false` and published with `publish-site`). The workflow checks this.
2. Preview first. A dry run prints the title, link, intro and card without posting:
   ```sh
   gh workflow run linkedin.yml -R EvanMcDowell31/evanmcdowell31.github.io -f slug=<slug> -f dry_run=true
   ```
   Or locally: `hugo --gc --minify && (cd tools/linkedin && go run . post -slug <slug> -site ../.. -dry-run)`.
3. Show Evan the preview and **get an explicit yes**. This posts publicly under his name, and a LinkedIn post can't be unposted cleanly.
4. Post with `-f dry_run=false`. Add `-f text="..."` to give the post its own intro; the default intro is the story's `summary`. Write it in Evan's voice: a hook line, short paragraphs, no hashtag spam.
5. Report the post URL from the run summary, then `git pull`, because the workflow pushed a front-matter commit.

The tool refuses to post a story that already has a `linkedin` URL, which is the guard against double-posting. Stories that started on LinkedIn (the Oct 19–22 posts) already have one.

## Setup and renewal (every 60 days)

LinkedIn's self-serve tokens last 60 days and can't be refreshed. A 401 error means it's time to reconnect.

1. **One-time app setup** at https://www.linkedin.com/developers/apps:
   - Create an app. It must be linked to a LinkedIn Page, which Evan creates or picks.
   - Add the products **Sign In with LinkedIn using OpenID Connect** and **Share on LinkedIn**.
   - Under Auth, add the redirect URL `http://localhost:8723/callback`.
2. **Connect.** Evan runs this himself, because it opens his browser to sign in:
   ```sh
   cd ~/Workspaces/Home/website/tools/linkedin
   LINKEDIN_CLIENT_ID=... LINKEDIN_CLIENT_SECRET=... go run . auth
   ```
3. Save the two printed values as repo secrets `LINKEDIN_ACCESS_TOKEN` and `LINKEDIN_PERSON_URN`. The command prints the `gh secret set` lines. Never paste the token into chat or commit it.

## Gotchas

- An API post skips the page's preview tags. The tool sends the card, title and description itself, so changing the page's tags won't change a post that's already out.
- Intro text escapes LinkedIn's reserved characters (`()[]{}<>#@*_~|\`), so they post literally. Hashtags therefore show as plain text, not clickable tags.
- `LINKEDIN_VERSION` (default `202604`) is LinkedIn's monthly API version. If calls start failing with a version error, bump it to a recent `YYYYMM`.
- Don't automate linkedin.com in a browser to post. It breaks LinkedIn's terms and can get the account restricted. Use the API only.
