# gdt-assets

A [GDT](https://github.com/monkeymonk/gdt) plugin that provides policy-driven asset operations, hygiene, and validation for game development teams.

`gdt-assets` helps teams inspect, standardize, reorganize, validate, and prepare assets for use in Godot projects. It catches common asset mistakes before they reach runtime or shipping builds, and provides deterministic, reproducible rules across teams and contributors.

## Features

- **Policy-driven** -- all checks configured via a single TOML file with named profiles
- **Asset scanning** -- inventory all project assets with type detection, hashing, and image metadata
- **Linting** -- naming conventions, folder structure, format restrictions, size thresholds, power-of-two checks
- **Duplicate detection** -- find exact duplicates (by hash) or likely duplicates (by filename)
- **Reference checking** -- detect broken `res://` paths in `.tscn` and `.tres` files
- **Reference repair** -- conservatively rewrite `res://` paths from rename manifests
- **Batch renaming** -- convert filenames to snake_case, kebab-case, or lowercase with collision detection
- **Rollback support** -- rename operations produce JSON manifests for safe undo
- **Package validation** -- block source files (.blend, .psd, etc.) from shipping builds
- **Lifecycle hooks** -- automatically validate assets before export and during development
- **Multiple output formats** -- table, JSON, CSV, Markdown for CI integration
- **Policy profiles** -- named overrides (mobile, release, etc.) for different build targets

## Supported Asset Types

| Category | Extensions |
|----------|-----------|
| Images | `.png`, `.jpg`, `.jpeg`, `.tga`, `.webp`, `.exr`, `.bmp`, `.hdr` |
| Vectors | `.svg` |
| Audio | `.wav`, `.ogg`, `.mp3`, `.flac`, `.opus` |
| Video | `.mp4`, `.webm`, `.mov`, `.ogv` |
| 3D Models | `.glb`, `.gltf`, `.obj`, `.fbx`, `.dae` |
| Animations | `.anim`, `.animlib` |
| Fonts | `.ttf`, `.otf`, `.woff`, `.woff2`, `.fnt` |
| Shaders | `.gdshader`, `.gdshaderinc`, `.glsl`, `.hlsl` |
| Documents | `.blend`, `.psd`, `.kra`, `.aseprite`, `.ase`, `.xcf`, `.ai`, `.afdesign` |
| Engine | `.tscn`, `.tres`, `.import`, `.godot` |

## Installation

### Via GDT (recommended)

```bash
gdt plugin install monkeymonk/gdt-plugin-assets
```

This automatically downloads a pre-built binary for your platform from GitHub Releases. If no release binary is available, gdt builds from source (requires Go 1.23+).

Verify the installation:

```bash
gdt plugin list          # should show "assets"
gdt doctor               # should show asset policy checks
```

### Manual install

```bash
git clone https://github.com/monkeymonk/gdt-plugin-assets.git
cd gdt-plugin-assets
make install             # builds and copies binary + manifest to ~/.gdt/plugins/gdt-plugin-assets/
```

### From source (development)

```bash
git clone https://github.com/monkeymonk/gdt-plugin-assets.git
cd gdt-plugin-assets
make build               # produces the "assets" binary locally
```

## Requirements

- [gdt](https://github.com/monkeymonk/gdt) >= 1.0
- Go 1.23+ (build from source only; pre-built binaries available via releases)

## Quick Start

```bash
# Set up asset policy for your project
gdt assets init

# See what assets you have
gdt assets scan

# Check for problems
gdt assets lint all

# Full health report
gdt assets report

# Health check
gdt assets doctor check
```

## Commands

### `gdt assets init`

Creates an `assets.policy.toml` file in your project root with sensible defaults.

```bash
gdt assets init
gdt assets init --with-sample-folders
```

The `--with-sample-folders` flag also creates the canonical asset directories defined in the policy.

### `gdt assets scan`

Discovers and lists all recognized assets in your project. Image assets include extracted metadata (dimensions, power-of-two status).

```bash
gdt assets scan                          # all assets, table format
gdt assets scan --type image             # images only
gdt assets scan --type image,audio       # images and audio
gdt assets scan --format json            # JSON output
gdt assets scan --format csv             # CSV output
gdt assets scan --format md              # Markdown table
gdt assets scan --hash                   # include SHA256 hashes
```

Example output:

```
  audio    assets/audio/blast_sfx.wav                         12 KB
  image    assets/images/hero_sprite.png                      245 KB
  model    assets/models/crate.glb                            1.2 MB

Total: 3 assets (1.4 MB)
  audio: 1, image: 1, model: 1
```

### `gdt assets lint`

Runs policy checks on assets without modifying anything. Returns structured exit codes suitable for CI gates.

```bash
gdt assets lint all                      # run all analyzers
gdt assets lint names                    # naming conventions only
gdt assets lint structure                # folder structure only
gdt assets lint images                   # image-specific checks
gdt assets lint audio                    # audio-specific checks
gdt assets lint models                   # 3D model checks
gdt assets lint all --format json        # machine-readable output
gdt assets lint all --profile mobile     # lint with mobile profile overrides
```

**Analyzers included:**

| Analyzer | What it checks | Category |
|----------|---------------|----------|
| `names` | Filename casing (snake_case by default), spaces in filenames | naming |
| `structure` | Assets placed in correct folders per policy | structure |
| `images` | Allowed formats, policy-driven size thresholds, power-of-two dimensions | metadata, optimization |
| `audio` | Preferred formats, policy-driven size thresholds | metadata, optimization |
| `models` | FBX format warnings, policy-driven size thresholds | metadata, optimization |

**Severity levels:**

| Level | Meaning | Affects exit code |
|-------|---------|-------------------|
| `INFO` | Informational | No |
| `WARNING` | Should be addressed | No |
| `ERROR` | Must be fixed | Yes (exit 3) |
| `BLOCKER` | Prevents export | Yes (exit 4) |

Example output:

```
  WARNING  assets/images/HeroBanner.png: filename "HeroBanner.png" does not match snake convention
  ERROR    assets/images/huge texture.png: filename contains spaces: "huge texture.png"
  WARNING  assets/images/legacy.bmp: format .bmp not in allowed list [png webp jpg]
  WARNING  assets/images/bg.png: dimensions 300x300 are not power-of-two
  WARNING  assets/audio/music.mp3: format .mp3 not in preferred list [ogg wav]
  WARNING  assets/models/Hero.fbx: FBX format detected; preferred formats: [glb gltf]

0 info, 5 warnings, 1 errors, 0 blockers
```

### `gdt assets report`

Generates a combined inventory and lint report.

```bash
gdt assets report                        # table format
gdt assets report --format json          # JSON for CI artifacts
gdt assets report --format md            # Markdown for documentation
gdt assets report --hash                 # include file hashes
```

### `gdt assets rename`

Batch renames assets to match the naming convention defined in your policy. Defaults to dry-run mode for safety. Detects collisions before applying and writes a rollback manifest for safe undo.

```bash
gdt assets rename                        # preview (dry-run is default)
gdt assets rename --dry-run              # explicit preview
gdt assets rename --apply                # execute renames (writes rollback manifest)
gdt assets rename --rollback <manifest>  # undo a previous rename
```

Example dry-run output:

```
  assets/images/HeroBanner.png -> assets/images/hero_banner.png
  assets/images/huge texture.png -> assets/images/huge_texture.png
  assets/models/Hero.fbx -> assets/models/hero.fbx

3 file(s) to rename

Dry run. Use --apply to execute.
```

When `--apply` is used, a rollback manifest (`.assets-rollback-YYYYMMDD-HHMMSS.json`) is written to the project root before any files are moved. Use `--rollback` with this file to undo.

Supported conventions: `snake` (default), `kebab`, `lower`.

### `gdt assets refs`

Checks and repairs asset references in Godot scene and resource files.

```bash
gdt assets refs check                                      # find broken res:// references
gdt assets refs repair --from-manifest <rollback.json>     # preview ref repairs
gdt assets refs repair --from-manifest <rollback.json> --apply  # apply ref repairs
```

Scans `.tscn`, `.tres`, and `.godot` files for `res://` paths that don't resolve to existing files.

The `repair` subcommand conservatively rewrites `res://` references based on a rename manifest. It only performs exact `path="res://..."` string replacements in known text-based engine files.

Example check output:

```
  BROKEN  level.tscn:3 -> res://assets/images/missing_icon.png
  BROKEN  level.tscn:4 -> res://scenes/enemy.tscn

2 broken reference(s)
```

### `gdt assets dedupe`

Finds duplicate assets in the project.

```bash
gdt assets dedupe                        # exact duplicates (by SHA256 hash)
gdt assets dedupe --name                 # likely duplicates (by filename)
```

Example output:

```
  Duplicate group [55bd32b35c482516]:
    assets/images/hero.png
    assets/images/hero_copy.png
    assets/old/hero_backup.png

1 duplicate group(s)
```

### `gdt assets package`

Validates the project for release readiness. Checks that no forbidden source files are present in the project tree.

```bash
gdt assets package
```

Example output when a `.blend` file is found:

```
  BLOCKER  source_assets/hero.blend: source/authoring file present in project tree

0 info, 0 warnings, 0 errors, 1 blockers

Package validation FAILED
```

### `gdt assets policy`

Inspect and validate the policy file.

```bash
gdt assets policy show                   # display current policy
gdt assets policy validate               # check policy file syntax
```

### `gdt assets doctor check`

Shows comprehensive plugin health: policy validity, asset inventory with type breakdown, lint summary, and broken reference count.

```bash
gdt assets doctor check
```

Example output:

```
OK   assets.policy.toml valid
OK   found 42 assets (image: 20, audio: 8, model: 5, font: 3, engine: 6)
OK   3 warning(s), no errors
OK   no broken references
```

## Policy Configuration

All behavior is driven by `assets.policy.toml` in your project root. Run `gdt assets init` to generate it with defaults.

### Full reference

```toml
version = 1

[naming]
case = "snake"                           # snake, kebab, camel, pascal
allow_spaces = false                     # reject filenames with spaces
allowed_chars = "a-z0-9_-/"              # allowed character set

[folders]
images = "assets/images"                 # expected location for images
audio = "assets/audio"                   # expected location for audio
models = "assets/models"                 # expected location for 3D models
vectors = "assets/vectors"               # expected location for vectors
fonts = "assets/fonts"                   # expected location for fonts
source = "source_assets"                 # authoring source files

[images]
max_size_default_kb = 4096               # max image file size in KB (default 4 MB)
max_size_ui_kb = 2048                    # max UI image file size in KB (default 2 MB)
require_power_of_two = true              # require POT dimensions for textures
allow_non_pot_for_ui = true              # exempt UI textures from POT requirement
allowed_formats = ["png", "webp", "jpg"] # accepted image formats

[audio]
preferred_runtime_formats = ["ogg", "wav"]  # preferred audio formats
allowed_sample_rates = [44100, 48000]       # valid sample rates
max_size_kb = 51200                         # max audio file size in KB (default 50 MB)

[models]
preferred_formats = ["glb", "gltf"]      # preferred 3D formats
warn_on_fbx = true                       # flag FBX usage
max_size_kb = 102400                     # max model file size in KB (default 100 MB)

[animations]
clip_case = "snake"                      # animation clip naming convention
baseline_fps = 30                        # expected animation FPS

[package.release]
forbid_source_files = true               # block .blend/.psd/etc in builds

[profiles.mobile]
[profiles.mobile.images]
max_size_default_kb = 1024               # tighter limit for mobile
require_power_of_two = true

[profiles.mobile.audio]
max_size_kb = 10240                      # 10 MB limit for mobile audio

[profiles.release]
[profiles.release.models]
max_size_kb = 51200                      # 50 MB limit for release builds
```

### Overriding defaults

Only set the fields you want to change. Unset fields keep their defaults. For example, to allow BMP images and MP3 audio:

```toml
version = 1

[images]
allowed_formats = ["png", "webp", "jpg", "bmp"]

[audio]
preferred_runtime_formats = ["ogg", "wav", "mp3"]
```

### Profiles

Named profiles let you apply different thresholds for different build targets. Use `--profile` with lint or set `GDT_ASSETS_PROFILE` for hooks:

```bash
gdt assets lint all --profile mobile     # lint with mobile limits
GDT_ASSETS_PROFILE=release gdt export    # before_export uses release profile
```

Profile overrides are merged on top of the base policy. Only specified fields are overridden; everything else keeps its base value.

## GDT Lifecycle Hooks

The plugin integrates with GDT's hook system to provide automatic validation at key moments.

### `after_new`

Triggered after creating a new project with `gdt new`. Suggests running `gdt assets init` to set up asset policy.

### `before_run`

Triggered before `gdt run`. Checks for broken asset references and warns if any are found, without blocking execution.

### `before_export`

Triggered before `gdt export`. Runs the full analyzer suite against project policy. Supports `GDT_ASSETS_PROFILE` env var for profile selection.

- **FAIL** (blocks export) if any blocker-level issues are found
- **WARN** (allows export) if errors are found but no blockers
- **OK** if all assets pass validation

### `doctor`

Contributes to `gdt doctor` output. Shows:

- Policy file validity
- Asset count with type breakdown
- Lint summary (blockers, errors, warnings)
- Broken reference count

## CI Integration

`gdt-assets` is designed for CI pipelines. All commands support non-interactive operation, deterministic exit codes, and machine-readable output.

### Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success, no issues |
| `1` | Operational failure or invalid usage |
| `2` | Policy file invalid |
| `3` | Diagnostics at error threshold |
| `4` | Diagnostics at blocker threshold |

### Example GitHub Actions step

```yaml
- name: Lint assets
  run: gdt assets lint all --format json > asset-lint.json

- name: Lint assets (mobile profile)
  run: gdt assets lint all --profile mobile

- name: Check references
  run: gdt assets refs check

- name: Validate package
  run: gdt assets package
```

### Example GitLab CI

```yaml
asset-validation:
  script:
    - gdt assets lint all
    - gdt assets refs check
    - gdt assets package
  artifacts:
    when: always
    paths:
      - asset-lint.json
```

## Architecture

```
gdt-assets (binary)
  main.go                  CLI entry point
  internal/
    cmd/                   Command handlers (init, scan, lint, ...)
    asset/                 Asset types, extension detection, image metadata extraction
    policy/                TOML policy loading with defaults and profile resolution
    scanner/               Filesystem walker with filtering, hashing, and metadata
    diagnostic/            Severity levels, categories, and diagnostic collections
    exitcode/              Structured exit code constants
    analyzer/              Pluggable analyzers (name, structure, image, audio, model)
    rename/                Batch rename with collision detection, plans, and rollback
    refs/                  Godot res:// reference scanner and conservative repair
    dedupe/                Hash and name-based duplicate detection
    report/                Multi-format output (table, JSON, CSV, Markdown)
```

## Development

```bash
# Build
make build

# Run tests
make test

# Install to GOPATH
make install
```

## License

See [LICENSE](LICENSE) for details.
