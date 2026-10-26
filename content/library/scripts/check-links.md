+++
title = "Find broken external links"
date = 2026-10-26
summary = "Build a Hugo site and list every external link that doesn't return 200."
tags = ["hugo", "bash"]
+++

```sh
hugo --gc --minify -d /tmp/site-check
grep -rhoE 'href="?https?://[^" >]+' /tmp/site-check | sed -E 's/href="?//' | sort -u \
  | while read -r u; do
      code=$(curl -s -o /dev/null -L -w '%{http_code}' --max-time 10 "$u")
      [ "$code" = 200 ] || echo "$code $u"
    done
```

LinkedIn answers bots with `999`, so those links usually aren't really broken.
