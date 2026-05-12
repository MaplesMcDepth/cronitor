# cronitor

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
```
