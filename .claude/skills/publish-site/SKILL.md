---
name: publish-site
description: Build, commit, push and verify a deploy of Evan's Hugo site to GitHub Pages (https://evanmcdowell31.github.io). Use when asked to publish, deploy, ship, push the site, or get changes live.
---

# Publish the site

Deploys are automatic: a push to `main` runs `.github/workflows/hugo.yml`, which builds with Hugo and deploys to GitHub Pages. This skill makes sure what goes out is correct, then confirms it landed.

## Steps

1. **Build locally.** Stop if there are errors or warnings:
   ```sh
   hugo --gc --minify
   ```
2. **Review what's going out:**
   - `git status` and `git diff`: summarize the change for Evan.
   - Any post that should go live has `draft = false`. Mention any drafts that are staying unpublished.
   - Scan the diff for anything that shouldn't be public: emails, internal URLs, tokens, `*-state.json` files, private metrics.
3. **Commit.** Use a concise message that describes the content change (e.g. "Add post on agent skills"). Follow the session's commit attribution rules.
4. **Push** with `git push`. The remote is `git@github.com:EvanMcDowell31/evanmcdowell31.github.io.git` over SSH, so no `gh` login is needed.
5. **Watch the deploy.** The repo is public, so the API needs no auth:
   ```sh
   curl -s "https://api.github.com/repos/EvanMcDowell31/evanmcdowell31.github.io/actions/runs?per_page=1" \
     | python3 -c "import json,sys;r=json.load(sys.stdin)['workflow_runs'][0];print(r['status'],r['conclusion'],r['html_url'])"
   ```
   Wait in the background with an `until` loop until the status is `completed`, and report a `failure` with the run URL. To read the logs, check the job steps through `/actions/runs/<id>/jobs`.
6. **Verify live.** `curl` the changed URLs on `https://evanmcdowell31.github.io` and confirm they return 200 with the new content. The CDN can lag the deploy by a minute or two.

## Gotchas

- The Pages source must stay set to **GitHub Actions** (repo Settings → Pages). If it's switched back to "Deploy from a branch", pushes will publish raw files instead of the Hugo build.
- `react-resume` is a different repo: Evan's old 2018 site, which was renamed from this repo's name. Never push the site there.
- To redeploy without content changes, push an empty commit (`git commit --allow-empty`).
