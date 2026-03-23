# Changelog

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
