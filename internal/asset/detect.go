package asset

import (
	"path/filepath"
	"strings"
)

var extTypeMap = map[string]AssetType{
	// Images
	".png": TypeImage, ".jpg": TypeImage, ".jpeg": TypeImage,
	".tga": TypeImage, ".webp": TypeImage, ".exr": TypeImage,
	".bmp": TypeImage, ".hdr": TypeImage,
	// Vectors
	".svg": TypeVector,
	// Audio
	".wav": TypeAudio, ".ogg": TypeAudio, ".mp3": TypeAudio,
	".flac": TypeAudio, ".opus": TypeAudio,
	// Video
	".mp4": TypeVideo, ".webm": TypeVideo, ".mov": TypeVideo,
	".ogv": TypeVideo,
	// 3D Models
	".glb": TypeModel, ".gltf": TypeModel, ".obj": TypeModel,
	".fbx": TypeModel, ".dae": TypeModel,
	// Animations
	".anim": TypeAnimation, ".animlib": TypeAnimation,
	// Fonts
	".ttf": TypeFont, ".otf": TypeFont, ".woff": TypeFont,
	".woff2": TypeFont, ".fnt": TypeFont,
	// Shaders
	".gdshader": TypeShader, ".gdshaderinc": TypeShader,
	".glsl": TypeShader, ".hlsl": TypeShader,
	// Documents / source sidecars
	".blend": TypeDocument, ".psd": TypeDocument,
	".kra": TypeDocument, ".aseprite": TypeDocument,
	".ase": TypeDocument, ".xcf": TypeDocument,
	".ai": TypeDocument, ".afdesign": TypeDocument,
	// Engine resources
	".tscn": TypeEngineResource, ".tres": TypeEngineResource,
	".import": TypeEngineResource, ".godot": TypeEngineResource,
}

func DetectType(path string) AssetType {
	ext := strings.ToLower(filepath.Ext(path))
	if t, ok := extTypeMap[ext]; ok {
		return t
	}
	return TypeUnknown
}
