+++
title = "Watch a GitHub Pages deploy"
date = 2026-10-26
summary = "Poll the latest Actions run for a public repo until it finishes. No auth needed."
tags = ["github", "bash"]
+++

```sh
repo="EvanMcDowell31/evanmcdowell31.github.io"
until [ "$(curl -s "https://api.github.com/repos/$repo/actions/runs?per_page=1" \
  | python3 -c 'import json,sys;print(json.load(sys.stdin)["workflow_runs"][0]["status"])')" = completed ]; do
  sleep 15
done
curl -s "https://api.github.com/repos/$repo/actions/runs?per_page=1" \
  | python3 -c 'import json,sys;r=json.load(sys.stdin)["workflow_runs"][0];print(r["conclusion"], r["html_url"])'
```

Prints `success` or `failure` with a link to the run.
