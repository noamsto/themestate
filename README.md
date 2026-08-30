# themestate

`themestate` provides a small, shared contract for desktop theme state.

The state file is `theme-state.json` under `$XDG_STATE_HOME`. When
`XDG_STATE_HOME` is unset or empty, it is under `$HOME/.local/state`. Its JSON
contents have a `theme` field, for example:

```json
{"theme":"dark"}
```

The package exports three functions:

- `Path()` resolves the state-file path.
- `Detect()` reads and parses the theme. It returns the non-empty `theme` value
  unchanged, and returns `"dark"` for any read or parse error or an empty value.
- `ModTime(path)` returns a file's modification time and the error from `os.Stat`.

The package does not watch the file. Callers that need live updates own the
polling policy and can use `ModTime` to detect changes.
