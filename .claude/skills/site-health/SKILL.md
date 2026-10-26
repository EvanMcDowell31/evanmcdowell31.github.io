---
name: site-health
description: Check and maintain Evan's Hugo site — broken links, build warnings, Hugo version upgrades, workflow action versions, and stale content. Use when asked to check, audit, maintain, upgrade or clean up the site, or when a deploy fails.
---

# Site health check

Run these checks and report the results as a short list (OK, fixed, or needs Evan). Fix mechanical issues; ask before changing content.

## 1. Build

```sh
hugo --gc --minify 2>&1 | grep -iE 'error|warn'
```
Any output needs fixing. Deprecation warnings name their replacement; this site already uses `locale`, `site.Language.Locale` and `hugo.Data`.

## 2. Hugo version

- Local version: `hugo version`. Latest: `brew info hugo`.
- The CI version is `HUGO_VERSION` in `.github/workflows/hugo.yml`. Keep local and CI on the same version. When upgrading, bump both, rebuild, and check the deprecation warnings.
- The templates use the Hugo 0.146+ layout system (`layouts/home.html`, `layouts/_partials/`), so don't go older.

## 3. Workflow actions

Check that `actions/checkout`, `configure-pages`, `upload-pages-artifact` and `deploy-pages` in `hugo.yml` are on current major versions. GitHub deprecates old ones, which fails deploys.

## 4. Links

Build, then check internal links and external links:
```sh
hugo --gc --minify -d /tmp/site-check
grep -rhoE 'href="?https?://[^" >]+' /tmp/site-check | sed -E 's/href="?//' | sort -u \
  | while read -r u; do printf '%s %s\n' "$(curl -s -o /dev/null -L -w '%{http_code}' --max-time 10 "$u")" "$u"; done | grep -v '^200'
```
LinkedIn returns 999 to bots, so it isn't really broken. Check internal links by serving `/tmp/site-check` or running `hugo server` and curling each page.

## 5. Live site

`curl` `/`, `/about/`, `/posts/` and `/index.xml` on `https://evanmcdowell31.github.io`. Each should return 200. Check the latest Actions run succeeded (see the `publish-site` skill).

## 6. Content freshness

Flag but don't edit:
- placeholder text such as "Hello, world" or "Replace with"
- `content/about.md` roles that may be out of date. Offer the `update-profile` skill.
- drafts (`draft = true`) older than a month
- the newest post being more than 3 months old, if Evan wants a regular cadence

If fixes were made, offer to publish them with `publish-site`.
