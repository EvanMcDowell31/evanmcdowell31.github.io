+++
title = "Most engineers aren't short on AI tools"
date = 2026-10-19T14:20:00-04:00
summary = "They're short on AI tools that are easy to run and worth what they cost. Here's the layer I built to fix that."
tags = ["ai", "agents", "developer-productivity"]
[params]
  linkedin = "https://www.linkedin.com/feed/update/urn:li:activity:7518013105889492992/"
+++

Most engineers aren't short on AI tools. They're short on AI tools that are easy to run and worth what they cost.

Since July I've been building a layer on top of Walmart's Wibey CLI. It packages advanced agent workflows so engineers can use them without having to become AI experts first.

## What it does

- **Concurrent agents on a standard laptop.** Several agents work in parallel on a feature, tests, a bug and PR checklists, all from a 36GB MacBook. No special hardware, no cluster.
- **Skills, not prompts.** 170+ version-controlled playbooks cover how we release, debug, test and ship. Engineers say what they want, and the right playbook loads. Expertise that took months to build up can be installed in seconds.
- **Deterministic tools for routine work.** Jira, Slack, PR preflight and release monitoring each take a single predictable call. Every session gets the same result, with no tokens spent rediscovering how things work.
- **Cost-aware model routing.** A frontier model plans and reviews. Search and routine code go to models that cost as little as a fifth as much. Every handoff is logged and graded, so ROI is measured, not assumed.

## Beyond my own stack

I'm a web frontend engineer. With the same setup, I've also shipped Sparky UI work in the Glass iOS app and am now contributing to Glass on Android, in codebases, languages and build systems I hadn't worked in before. The skills hold each platform's conventions, build steps and gotchas, so getting productive on a new platform took days instead of months.

## The return so far

- ~3x merged PRs per month (Aug–Sep vs. the prior 10 months)
- 2,000+ automated test cases in one 30-day stretch
- Shipping on web, iOS and Android, as one engineer
- All of it on hardware and tooling we already own

## What I'd tell leaders investing in GenAI for engineering

1. **Adoption depends on ease of use.** If an advanced capability takes an expert to run, most engineers won't use it. Package it.
2. **Route work to the right model.** Paying frontier prices for grep is the easiest money to save.
3. **Instrument everything.** If you can't show cost per outcome, you can't defend the investment, or scale it.

The upside isn't one engineer going faster. It's making that capability easy for everyone to use.
