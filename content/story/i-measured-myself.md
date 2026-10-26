+++
title = "I had no gauge, so I measured myself"
date = 2026-10-20T18:05:00-04:00
summary = "A year of merged PRs, product work only: from the 57th–85th percentile to the 99.95th, and what actually changed."
tags = ["ai", "agents", "developer-productivity"]
[params]
  linkedin = "https://www.linkedin.com/feed/update/urn:li:activity:7518432134043299840/"
+++

As a contractor, there's rarely a ranking or calibration to tell you how you compare. I've also spent the past few months building AI tooling for engineering enablement, and I wanted hard numbers on whether it's actually worth it.

So I used the one test subject I had full access to: me.

I tracked my merged PRs per month for a year, product work only. Tooling repos didn't count.

## What I found

- Most of the year, I was already above average, between the 57th and 85th percentile.
- In August, I hit the top 10% for the first time.
- In September, the 99.95th percentile.
- Across those two months, about 9× my own baseline.

I didn't suddenly get better at engineering. I got better at using AI.

## What changed

- Parallel agents instead of one chat at a time
- Version-controlled playbooks, so no session starts from zero
- Deterministic tools for routine work like tickets, chat and PR checks
- Subagents for heavy research, and reusing the agent that already has the context. Most of my token spend was re-reading context, so keeping it small mattered more than picking a cheaper model.

For context, public studies (GitHub, Google, Microsoft/Accenture, METR) put AI's effect on developer productivity anywhere from about 20% slower to about 55% faster. Having access to AI isn't what makes the difference. How you use it does.

## The caveats

PR count measures volume, not impact. Two months is a short run, and it's a sample of one. I built the tools, so I'm also their most practiced user.

That's why the next step matters more than the chart. I'm packaging this workflow so any engineer can install it and get the same setup without building it from scratch. Then I'll measure again.

One engineer going faster is nice. A whole team getting even part of that lift is the real impact.
