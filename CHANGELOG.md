# Changelog

## [Unreleased]

## [0.3.1] - 2026-09-25

### Fixed

- `hookBeforeExport`'s asset-error path printed a `WARN` status line
  while returning a non-zero exit code (`exitcode.ErrDiagnostics`) —
  gdt's V2 hook protocol treats any non-zero exit as fatal regardless
  of the printed line, so `gdt export` was actually aborting despite
  the "WARN" framing suggesting otherwise. Now prints `FAIL` to match
  the exit code's actual (intentional, since v0.2.0) behavior.
- CI's `govulncheck` job pinned `golang.org/x/vuln/cmd/govulncheck` to
  `v1.7.0` instead of `@latest` — `v1.8.0` requires Go >=1.26, but this
  project targets Go 1.25.
- `refs repair`'s planning pass now warns on stderr when a `.tscn`/
  `.tres` file can't be read instead of silently skipping it —
  previously an unreadable file meant fewer repairs than actually
  needed, with no indication.
- `gdt-assets scan` (and every other command going through the shared
  scanner) now warns on stderr when a path can't be read during the
  filesystem walk (permission denied, broken symlink, etc.) instead of
  silently skipping it and under-reporting the asset inventory.
- Six diagnostics (`image.pot`, `audio.format`, `model.fbx`,
  `naming.no_spaces`, `naming.case`, `structure.folder`) were missing
  the `Category` field that the rest of the diagnostic taxonomy sets
  (`optimization`/`metadata`/`naming`/`structure`) — affects the
  documented JSON output schema.

### Changed

- `plugin.toml`'s `requires_gdt` lowered from the scaffold-default
  `>=1.0` to `>=0.3.0`, matching the actual gdt version this plugin has
  been verified against (gdt does not yet enforce this constraint, but
  the declared value should reflect reality once it does).
- `rename.ToSnakeCase` collapses consecutive underscores via a single
  regex pass instead of a `Contains`/`ReplaceAll` loop — same output,
  consistent with the file's other regex-based transforms.

## [0.3.0] - 2026-07-21

### Changed

- Build with Go 1.25 (latest patch via `go.mod` + `check-latest`) for
  standard-library security fixes
- CI now runs `golangci-lint`, `govulncheck`, and a gofmt check across
  linux/macOS/windows, in addition to build/vet/test
- Release binaries now ship a `checksums.txt` signed with keyless cosign
  (`checksums.txt.bundle`, Sigstore)
- Bumped all GitHub Actions off the deprecated Node 20 runtime; added Dependabot
- Use `strings.ReplaceAll` in refs repair (idiomatic; flagged by staticcheck)

### Fixed

- `gofmt` formatting in `internal/cmd/init.go` and `internal/policy/policy.go`

## [0.2.0] - 2026-03-24

### Added
- Image metadata extraction (width, height, power-of-two) via stdlib `image.DecodeConfig`
- Scanner auto-populates image dimensions during scan
- Power-of-two validation for images with UI path exemption (`/ui/`, `/gui/`, `/hud/`)
- Diagnostic `category` field for grouping (naming, structure, metadata, optimization)
- Structured exit codes: 0 (OK), 1 (usage), 2 (policy), 3 (errors), 4 (blockers)
- Rename collision detection (multi-source and existing-target)
- Rollback manifests for rename operations (`--rollback` flag to undo)
- Conservative `refs repair` command via `--from-manifest` for rewriting `res://` paths
- Policy profile system with named profiles and selective overrides
- `--profile` flag on `lint` command
- `GDT_ASSETS_PROFILE` env var for `before_export` hook

### Changed
- Analyzers now use policy-driven max size thresholds instead of hardcoded values
- Doctor command shows lint summary, broken ref count, and asset type breakdown
- Image policy fields renamed: `max_size_default` -> `max_size_default_kb`, `max_size_ui` -> `max_size_ui_kb`
- Audio and model policies gain `max_size_kb` field
- Asset struct now includes JSON tags for clean serialized output
- `refs check` and `refs repair` use shared `FindEngineFiles` helper

### Breaking Changes
- Policy TOML keys renamed: `max_size_default` -> `max_size_default_kb`, `max_size_ui` -> `max_size_ui_kb`
- JSON output field names changed to snake_case (e.g., `AbsPath` -> `abs_path`)
- `hookBeforeExport` now returns exit code 3 on errors (was 0)

## [0.1.0] - 2026-03-23

### Added
- Asset scanning with type detection and SHA256 hashing
- Policy-driven linting (naming, structure, images, audio, models)
- Batch rename with snake_case, kebab-case, lowercase support
- Broken `res://` reference detection
- Duplicate detection (hash-based and name-based)
- Package validation (forbid source files in builds)
- GDT lifecycle hooks (after_new, before_run, before_export)
- Doctor integration
- Shell completions for bash, zsh, and fish
- Multiple output formats (table, JSON, CSV, Markdown)
- GitHub Actions CI and cross-platform release workflow
