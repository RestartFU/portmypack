package bedrock

import (
	"archive/zip"
	"encoding/json"

	"github.com/google/uuid"
)

// ManifestFormatVersion is the manifest format portmypack writes. Bedrock 1.13 introduced format 2 and
// current 1.26 clients still accept it; they warn that format 0 is deprecated, and format 3 is only
// understood by recent clients.
const ManifestFormatVersion = 2

// DefaultMinEngineVersion is the min_engine_version written by default: the first Bedrock release that
// reads format 2 manifests. Clients older than this value refuse the pack, newer ones accept it.
var DefaultMinEngineVersion = [3]int{1, 13, 0}

// Manifest represents a Minecraft resource pack manifest.
type Manifest struct {
	// FormatVersion is the format version of the manifest.
	FormatVersion int `json:"format_version"`
	// Header is the header of the manifest.
	Header ManifestHeader `json:"header"`
	// Modules is the modules of the resource pack.
	Modules []ManifestModule `json:"modules"`
}

// ManifestHeader represents a Minecraft resource pack manifest header.
type ManifestHeader struct {
	// Description is the description of the resource pack.
	Description string `json:"description"`
	// Name is the name of the resource pack.
	Name string `json:"name"`
	// UUID is the UUID of the resource pack.
	UUID string `json:"uuid"`
	// Version is the version of the resource pack.
	Version [3]int `json:"version"`
	// MinEngineVersion is the oldest Bedrock version allowed to load the pack.
	MinEngineVersion [3]int `json:"min_engine_version"`
}

// ManifestModule represents a Minecraft resource pack manifest module.
type ManifestModule struct {
	// Description is the description of the module.
	Description string `json:"description"`
	// Type is the type of the module.
	Type string `json:"type"`
	// UUID is the UUID of the module.
	UUID string `json:"uuid"`
	// Version is the version of the module.
	Version [3]int `json:"version"`
}

// NewManifest creates a new resource pack manifest that loads on Bedrock DefaultMinEngineVersion and later.
func NewManifest(name, description string) Manifest {
	return NewManifestWithEngine(name, description, DefaultMinEngineVersion)
}

// NewManifestWithEngine creates a new resource pack manifest with an explicit min_engine_version.
func NewManifestWithEngine(name, description string, minEngine [3]int) Manifest {
	header := ManifestHeader{
		Description:      description,
		Name:             name,
		UUID:             uuid.NewString(),
		Version:          [3]int{1, 0, 0},
		MinEngineVersion: minEngine,
	}
	modules := []ManifestModule{
		{
			Description: description,
			Type:        "resources",
			UUID:        uuid.NewString(),
			Version:     [3]int{1, 0, 0},
		},
	}
	return Manifest{
		FormatVersion: ManifestFormatVersion,
		Header:        header,
		Modules:       modules,
	}
}

// generateManifest generates the manifest for the resource pack.
func (r ResourcePack) generateManifest(w *zip.Writer) error {
	description := r.Description
	if description == "" {
		description = "github.com/restartfu"
	}
	minEngine := r.MinEngineVersion
	if minEngine == ([3]int{}) {
		minEngine = DefaultMinEngineVersion
	}
	manifest := NewManifestWithEngine(r.Name, description, minEngine)
	f, err := w.Create("manifest.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(manifest)
}
