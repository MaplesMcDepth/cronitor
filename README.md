# cronitor

![CI](https://github.com/MaplesMcDepth/cronitor/actions/workflows/ci.yml/badge.svg)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)


Cron expression validator and explainer.

## Install

```bash
go install github.com/MaplesMcDepth/cronitor/cmd/cronitor@latest
```

## Commands

### `validate` — Check if expression is valid
```bash
cronitor validate "*/5 * * * *"
cronitor validate "0 9 * * 1-5"
```

### `explain` — Convert to human-readable text
```bash
cronitor explain "0 9 * * 1-5"
# → At 0 minute, 9 hour, from 1-5
```

### `next` — Show next execution times
```bash
cronitor next "0 */6 * * *"
cronitor next -n 10 "0 0 * * 0"
```

### `list` — Parse and list fields
```bash
cronitor list "0 0 * * 0"
# Field           Value
# ------------------------------
# Minute          0
# Hour            0
# Day of Month    *
# Month           *
# Day of Week     0
```

## Common use cases

### Check a weekday business-hours schedule
```bash
cronitor validate "0 9-17 * * 1-5"
# Valid cron expression

cronitor explain "0 9-17 * * 1-5"
# At 0 minute, from 9-17, from 1-5
```

### Preview the next times a job will run
```bash
cronitor next -n 3 "0 */6 * * *"
# Next 3 executions:
#
# 1. 2026-07-05 06:00:00 AEST
# 2. 2026-07-05 12:00:00 AEST
# 3. 2026-07-05 18:00:00 AEST
```

### Inspect an unfamiliar expression field-by-field
```bash
cronitor list "30 2 1 * *"
# Field           Value
# ------------------------------
# Minute          30
# Hour            2
# Day of Month    1
# Month           *
# Day of Week     *
```

## Tips and tricks

- Quote cron expressions so your shell does not expand `*` before `cronitor` sees it.
- Use `validate` first when debugging a schedule copied from docs or dashboards.
- Pair `explain` with `next` to sanity-check both the wording and actual upcoming run times.
