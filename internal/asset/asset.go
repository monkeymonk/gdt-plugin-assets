package asset

import "fmt"

// AssetType classifies an asset by its general family.
type AssetType int

const (
	TypeUnknown AssetType = iota
	TypeImage
	TypeVector
	TypeAudio
	TypeVideo
	TypeModel
	TypeAnimation
	TypeFont
	TypeShader
	TypeDocument
	TypeEngineResource
)

func (t AssetType) String() string {
	switch t {
	case TypeImage:
		return "image"
	case TypeVector:
		return "vector"
	case TypeAudio:
		return "audio"
	case TypeVideo:
		return "video"
	case TypeModel:
		return "model"
	case TypeAnimation:
		return "animation"
	case TypeFont:
		return "font"
	case TypeShader:
		return "shader"
	case TypeDocument:
		return "document"
	case TypeEngineResource:
		return "engine"
	default:
		return "unknown"
	}
}

// HumanSize formats a byte count for display.
func HumanSize(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.0f KB", float64(b)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

// Asset represents a discovered project asset.
type Asset struct {
	Path    string    // relative path from project root
	AbsPath string    // absolute filesystem path
	Type    AssetType
	Size    int64
	Hash    string   // sha256, populated on demand
	Tags    []string
}
