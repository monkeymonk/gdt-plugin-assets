package policy

func Default() Policy {
	return Policy{
		Version: 1,
		Naming: NamingPolicy{
			Case:         "snake",
			AllowSpaces:  false,
			AllowedChars: "a-z0-9_-/",
		},
		Folders: FolderPolicy{
			Images:  "assets/images",
			Audio:   "assets/audio",
			Models:  "assets/models",
			Vectors: "assets/vectors",
			Fonts:   "assets/fonts",
			Source:  "source_assets",
		},
		Images: ImagePolicy{
			MaxSizeDefault:    4096,
			MaxSizeUI:         2048,
			RequirePowerOfTwo: true,
			AllowNonPotForUI:  true,
			AllowedFormats:    []string{"png", "webp", "jpg"},
		},
		Audio: AudioPolicy{
			PreferredFormats:   []string{"ogg", "wav"},
			AllowedSampleRates: []int{44100, 48000},
		},
		Models: ModelPolicy{
			PreferredFormats: []string{"glb", "gltf"},
			WarnOnFBX:        true,
		},
		Animations: AnimationPolicy{
			ClipCase:    "snake",
			BaselineFPS: 30,
		},
		Package: PackagePolicy{
			Release: ReleasePolicy{
				ForbidSourceFiles: true,
			},
		},
	}
}
