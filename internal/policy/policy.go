package policy

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
)

const FileName = "assets.policy.toml"

type Policy struct {
	Version    int                        `toml:"version"`
	Naming     NamingPolicy               `toml:"naming"`
	Folders    FolderPolicy               `toml:"folders"`
	Images     ImagePolicy                `toml:"images"`
	Audio      AudioPolicy                `toml:"audio"`
	Models     ModelPolicy                `toml:"models"`
	Animations AnimationPolicy            `toml:"animations"`
	Package    PackagePolicy              `toml:"package"`
	Profiles   map[string]ProfileOverrides `toml:"profiles,omitempty"`
}

type NamingPolicy struct {
	Case         string `toml:"case"`
	AllowSpaces  bool   `toml:"allow_spaces"`
	AllowedChars string `toml:"allowed_chars"`
}

type FolderPolicy struct {
	Images  string `toml:"images"`
	Audio   string `toml:"audio"`
	Models  string `toml:"models"`
	Vectors string `toml:"vectors"`
	Fonts   string `toml:"fonts"`
	Source  string `toml:"source"`
}

type ImagePolicy struct {
	MaxSizeDefaultKB  int      `toml:"max_size_default_kb"`
	MaxSizeUIKB       int      `toml:"max_size_ui_kb"`
	RequirePowerOfTwo bool     `toml:"require_power_of_two"`
	AllowNonPotForUI  bool     `toml:"allow_non_pot_for_ui"`
	AllowedFormats    []string `toml:"allowed_formats"`
}

type AudioPolicy struct {
	PreferredFormats   []string `toml:"preferred_runtime_formats"`
	AllowedSampleRates []int    `toml:"allowed_sample_rates"`
	MaxSizeKB          int      `toml:"max_size_kb"`
}

type ModelPolicy struct {
	PreferredFormats []string `toml:"preferred_formats"`
	WarnOnFBX        bool     `toml:"warn_on_fbx"`
	MaxSizeKB        int      `toml:"max_size_kb"`
}

type AnimationPolicy struct {
	ClipCase    string `toml:"clip_case"`
	BaselineFPS int    `toml:"baseline_fps"`
}

type PackagePolicy struct {
	Release ReleasePolicy `toml:"release"`
}

type ReleasePolicy struct {
	ForbidSourceFiles bool `toml:"forbid_source_files"`
}

type ProfileOverrides struct {
	Images *ImageOverrides `toml:"images,omitempty"`
	Audio  *AudioOverrides `toml:"audio,omitempty"`
	Models *ModelOverrides `toml:"models,omitempty"`
}

type ImageOverrides struct {
	MaxSizeDefaultKB  *int      `toml:"max_size_default_kb,omitempty"`
	MaxSizeUIKB       *int      `toml:"max_size_ui_kb,omitempty"`
	RequirePowerOfTwo *bool     `toml:"require_power_of_two,omitempty"`
	AllowedFormats    *[]string `toml:"allowed_formats,omitempty"`
}

type AudioOverrides struct {
	MaxSizeKB *int `toml:"max_size_kb,omitempty"`
}

type ModelOverrides struct {
	MaxSizeKB *int `toml:"max_size_kb,omitempty"`
}

func Load(path string) (*Policy, error) {
	p := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := toml.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func MarshalDefault() ([]byte, error) {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(Default()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func LoadOrDefault(path string) *Policy {
	p, err := Load(path)
	if err != nil {
		d := Default()
		return &d
	}
	return p
}
