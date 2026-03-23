package policy

import "fmt"

func ResolveProfile(base *Policy, name string) (*Policy, error) {
	if name == "" {
		return base, nil
	}

	overrides, ok := base.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("unknown profile: %q", name)
	}

	resolved := *base

	if o := overrides.Images; o != nil {
		if o.MaxSizeDefaultKB != nil {
			resolved.Images.MaxSizeDefaultKB = *o.MaxSizeDefaultKB
		}
		if o.MaxSizeUIKB != nil {
			resolved.Images.MaxSizeUIKB = *o.MaxSizeUIKB
		}
		if o.RequirePowerOfTwo != nil {
			resolved.Images.RequirePowerOfTwo = *o.RequirePowerOfTwo
		}
		if o.AllowedFormats != nil {
			resolved.Images.AllowedFormats = *o.AllowedFormats
		}
	}

	if o := overrides.Audio; o != nil {
		if o.MaxSizeKB != nil {
			resolved.Audio.MaxSizeKB = *o.MaxSizeKB
		}
	}

	if o := overrides.Models; o != nil {
		if o.MaxSizeKB != nil {
			resolved.Models.MaxSizeKB = *o.MaxSizeKB
		}
	}

	return &resolved, nil
}
