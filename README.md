# gdt-assets

A [GDT](https://github.com/monkeymonk/gdt) plugin that provides policy-driven asset operations, hygiene, and validation for game development teams.

`gdt-assets` helps teams inspect, standardize, reorganize, validate, and prepare assets for use in Godot projects. It catches common asset mistakes before they reach runtime or shipping builds, and provides deterministic, reproducible rules across teams and contributors.

## Features

- **Policy-driven** -- all checks configured via a single TOML file
- **Asset scanning** -- inventory all project assets with type detection and hashing
- **Linting** -- naming conventions, folder structure, format restrictions, size thresholds
- **Duplicate detection** -- find exact duplicates (by hash) or likely duplicates (by filename)
- **Reference checking** -- detect broken `res://` paths in `.tscn` and `.tres` files
- **Batch renaming** -- convert filenames to snake_case, kebab-case, or lowercase with dry-run preview
- **Package validation** -- block source files (.blend, .psd, etc.) from shipping builds
- **Lifecycle hooks** -- automatically validate assets before export and during development
- **Multiple output formats** -- table, JSON, CSV, Markdown for CI integration

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

Discovers and lists all recognized assets in your project.

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

Runs policy checks on assets without modifying anything. Returns exit code 1 if errors or blockers are found, making it suitable for CI gates.

```bash
gdt assets lint all                      # run all analyzers
gdt assets lint names                    # naming conventions only
gdt assets lint structure                # folder structure only
gdt assets lint images                   # image-specific checks
gdt assets lint audio                    # audio-specific checks
gdt assets lint models                   # 3D model checks
gdt assets lint all --format json        # machine-readable output
```

**Analyzers included:**

| Analyzer | What it checks |
|----------|---------------|
| `names` | Filename casing (snake_case by default), spaces in filenames |
| `structure` | Assets placed in correct folders per policy |
| `images` | Allowed formats (png/webp/jpg default), file size thresholds (>10MB) |
| `audio` | Preferred formats (ogg/wav default), file size thresholds (>50MB) |
| `models` | FBX format warnings, file size thresholds (>100MB) |

**Severity levels:**

| Level | Meaning | Affects exit code |
|-------|---------|-------------------|
| `INFO` | Informational | No |
| `WARNING` | Should be addressed | No |
| `ERROR` | Must be fixed | Yes |
| `BLOCKER` | Prevents export | Yes |

Example output:

```
  WARNING  assets/images/HeroBanner.png: filename "HeroBanner.png" does not match snake convention
  ERROR    assets/images/huge texture.png: filename contains spaces: "huge texture.png"
  WARNING  assets/images/legacy.bmp: format .bmp not in allowed list [png webp jpg]
  WARNING  assets/audio/music.mp3: format .mp3 not in preferred list [ogg wav]
  WARNING  assets/models/Hero.fbx: FBX format detected; preferred formats: [glb gltf]

0 info, 4 warnings, 1 errors, 0 blockers
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

Batch renames assets to match the naming convention defined in your policy. Defaults to dry-run mode for safety.

```bash
gdt assets rename                        # preview (dry-run is default)
gdt assets rename --dry-run              # explicit preview
gdt assets rename --apply                # execute renames
```

Example dry-run output:

```
  assets/images/HeroBanner.png -> assets/images/hero_banner.png
  assets/images/huge texture.png -> assets/images/huge_texture.png
  assets/models/Hero.fbx -> assets/models/hero.fbx

3 file(s) to rename

Dry run. Use --apply to execute.
```

Supported conventions: `snake` (default), `kebab`, `lower`.

### `gdt assets refs`

Checks and repairs asset references in Godot scene and resource files.

```bash
gdt assets refs check                    # find broken res:// references
gdt assets refs repair                   # interactive repair (planned)
```

Scans `.tscn`, `.tres`, and `.godot` files for `res://` paths that don't resolve to existing files.

Example output:

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
max_size_default = 4096                  # max texture dimension (px)
max_size_ui = 2048                       # max UI texture dimension (px)
require_power_of_two = true              # require POT dimensions
allow_non_pot_for_ui = true              # exempt UI textures from POT
allowed_formats = ["png", "webp", "jpg"] # accepted image formats

[audio]
preferred_runtime_formats = ["ogg", "wav"]  # preferred audio formats
allowed_sample_rates = [44100, 48000]       # valid sample rates

[models]
preferred_formats = ["glb", "gltf"]      # preferred 3D formats
warn_on_fbx = true                       # flag FBX usage

[animations]
clip_case = "snake"                      # animation clip naming convention
baseline_fps = 30                        # expected animation FPS

[package.release]
forbid_source_files = true               # block .blend/.psd/etc in builds
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

## GDT Lifecycle Hooks

The plugin integrates with GDT's hook system to provide automatic validation at key moments.

### `after_new`

Triggered after creating a new project with `gdt new`. Suggests running `gdt assets init` to set up asset policy.

### `before_run`

Triggered before `gdt run`. Checks for broken asset references and warns if any are found, without blocking execution.

### `before_export`

Triggered before `gdt export`. Runs the full analyzer suite against project policy:

- **FAIL** (blocks export) if any blocker-level issues are found
- **WARN** (allows export) if errors are found but no blockers
- **OK** if all assets pass validation

### `doctor`

Contributes to `gdt doctor` output. Checks:

- Whether `assets.policy.toml` exists and is valid
- Whether the project asset scan completes successfully

## CI Integration

`gdt-assets` is designed for CI pipelines. All commands support non-interactive operation, deterministic exit codes, and machine-readable output.

### Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success, no errors |
| `1` | Errors or blockers found |

### Example GitHub Actions step

```yaml
- name: Lint assets
  run: gdt assets lint all --format json > asset-lint.json

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
    asset/                 Asset types and extension-based detection
    policy/                TOML policy loading with defaults
    scanner/               Filesystem walker with filtering and hashing
    diagnostic/            Severity levels and diagnostic collections
    analyzer/              Pluggable analyzers (name, structure, image, audio, model)
    rename/                Batch rename with case conversion
    refs/                  Godot res:// reference scanner
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
