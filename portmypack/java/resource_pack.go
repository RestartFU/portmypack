package java

import (
	"encoding/json"
	"errors"
	"io/fs"
	"path"
	"strings"

	"github.com/restartfu/portmypack/portmypack/image"
	"github.com/restartfu/portmypack/portmypack/internal/fsutil"
	"github.com/restartfu/portmypack/portmypack/internal/logger"
)

// Layout identifies how a Java Edition pack names its texture folders.
type Layout int

const (
	// LayoutLegacy is the 1.8-1.12 layout: textures/blocks, textures/items and pre-flattening names.
	LayoutLegacy Layout = iota
	// LayoutFlattened is the 1.13+ layout: textures/block, textures/item and flattened names.
	LayoutFlattened
)

// String returns a short description of the layout.
func (l Layout) String() string {
	if l == LayoutFlattened {
		return "flattened (1.13+)"
	}
	return "legacy (1.8-1.12)"
}

// ResourcePack is the subset of a Java Edition resource pack that can be ported to Bedrock Edition.
type ResourcePack struct {
	Name string
	// Description is the plain-text pack description from pack.mcmeta.
	Description string
	// Format is the pack format from pack.mcmeta: pack_format, or the major min_format of 1.21.9+
	// packs. It is 0 when pack.mcmeta is missing or unreadable.
	Format int
	// Layout is the texture folder layout the pack uses.
	Layout Layout

	r fs.FS

	Skies []image.Texture

	PackIcon  image.Texture
	Icons     image.Texture
	Particles image.Texture
	// Widgets is gui/widgets.png, which holds the hotbar in packs made before 1.20.2.
	Widgets image.Texture

	Items  []image.Texture
	Blocks []image.Texture
	// Armors holds worn armour layers named <material>_layer_1.png and <material>_layer_2.png,
	// whichever Java layout they were read from.
	Armors []image.Texture
	// Sprites holds the 1.20.2+ HUD sprites keyed by their path below gui/sprites/hud without the
	// extension, such as "heart/full" or "hotbar".
	Sprites map[string]image.Texture
}

// FS returns the pack contents, rooted at the folder that holds pack.mcmeta.
func (pck ResourcePack) FS() fs.FS {
	return pck.r
}

// NewResourcePack reads a Java Edition resource pack from a zip archive or a folder.
func NewResourcePack(filePath string) (ResourcePack, error) {
	pck := ResourcePack{}
	pck.Name = fsutil.PackName(filePath)

	archive, err := fsutil.OpenPack(filePath)
	if err != nil {
		return ResourcePack{}, errors.New("error opening java pack: " + err.Error())
	}
	r, err := packRoot(archive)
	if err != nil {
		return ResourcePack{}, err
	}
	pck.r = r

	pck.readMeta()

	path, found := fsutil.FindDirectory(r, "world0")
	if found {
		pck.Skies, err = findSkies(r, path)
		if err != nil {
			logger.Debugf("error trying to find skies: %s", err)
		}
	}

	pck.PackIcon, err = image.NewTextureFS(r, "pack.png", false)
	if err != nil {
		logger.Debugf("error trying to find icon: %s", err)
	}

	texturesPath := "assets/minecraft/textures"
	if !fsutil.IsDir(r, texturesPath) {
		var found bool
		texturesPath, found = fsutil.FindDirectory(r, "textures")
		if !found {
			return ResourcePack{}, errors.New("could not find textures path, stopping progress")
		}
	}
	pck.Layout = detectLayout(r, texturesPath, pck.Format)
	logger.Debugf("pack format %d, %s layout", pck.Format, pck.Layout)

	pck.Icons, _ = optionalTexture(r, texturesPath+"/gui/icons.png", true)
	pck.Widgets, _ = optionalTexture(r, texturesPath+"/gui/widgets.png", false)
	pck.Particles, _ = optionalTexture(r, texturesPath+"/particle/particles.png", false)

	itemsDir, blocksDir := "/items", "/blocks"
	if pck.Layout == LayoutFlattened {
		itemsDir, blocksDir = "/item", "/block"
	}
	pck.Items = optionalDirectory(r, texturesPath+itemsDir, true)
	pck.Blocks = optionalDirectory(r, texturesPath+blocksDir, false)
	pck.Armors = readArmors(r, texturesPath)
	pck.Sprites = readSprites(r, texturesPath+"/gui/sprites/hud")

	return pck, nil
}

// packRoot returns the sub-filesystem rooted where pack.mcmeta (or, failing that, assets/) lives.
func packRoot(archive fs.FS) (fs.FS, error) {
	root, found := fsutil.FindPackRoot(archive, "pack.mcmeta")
	if !found {
		var assets string
		if assets, found = fsutil.FindDirectory(archive, "assets"); found {
			root = path.Dir(assets)
		}
	}
	if !found || root == "." {
		return archive, nil
	}
	return fs.Sub(archive, root)
}

// packMeta is the part of pack.mcmeta portmypack reads.
type packMeta struct {
	Pack struct {
		PackFormat  int             `json:"pack_format"`
		MinFormat   json.RawMessage `json:"min_format"`
		Description json.RawMessage `json:"description"`
	} `json:"pack"`
}

// readMeta fills Format and Description from pack.mcmeta when it can be read.
func (pck *ResourcePack) readMeta() {
	buf, err := fs.ReadFile(pck.r, "pack.mcmeta")
	if err != nil {
		logger.Debugf("error reading pack.mcmeta: %s", err)
		return
	}
	var meta packMeta
	if err := json.Unmarshal(buf, &meta); err != nil {
		logger.Debugf("error decoding pack.mcmeta: %s", err)
		return
	}
	pck.Format = meta.Pack.PackFormat
	if pck.Format == 0 {
		pck.Format = majorFormat(meta.Pack.MinFormat)
	}
	pck.Description = textComponent(meta.Pack.Description)
}

// majorFormat decodes a 1.21.9+ pack format, written as a number or as [major, minor].
func majorFormat(raw json.RawMessage) int {
	var n int
	if json.Unmarshal(raw, &n) == nil {
		return n
	}
	var v []int
	if json.Unmarshal(raw, &v) == nil && len(v) > 0 {
		return v[0]
	}
	return 0
}

// textComponent flattens a Java text component (string, object or array) into plain text without
// formatting codes.
func textComponent(raw json.RawMessage) string {
	var b strings.Builder
	var walk func(v any)
	walk = func(v any) {
		switch t := v.(type) {
		case string:
			b.WriteString(t)
		case []any:
			for _, e := range t {
				walk(e)
			}
		case map[string]any:
			walk(t["text"])
			walk(t["extra"])
		}
	}
	var v any
	if json.Unmarshal(raw, &v) == nil {
		walk(v)
	}
	s := []rune(b.String())
	var out []rune
	for i := 0; i < len(s); i++ {
		if s[i] == '§' && i+1 < len(s) {
			i++
			continue
		}
		out = append(out, s[i])
	}
	return strings.TrimSpace(string(out))
}

// firstFlattenedFormat is the pack format of 1.13, which renamed blocks/ and items/.
const firstFlattenedFormat = 4

// detectLayout picks the texture layout from the folders present, falling back to the pack format.
func detectLayout(r fs.FS, texturesPath string, format int) Layout {
	switch {
	case fsutil.IsDir(r, texturesPath+"/block"), fsutil.IsDir(r, texturesPath+"/item"):
		return LayoutFlattened
	case fsutil.IsDir(r, texturesPath+"/blocks"), fsutil.IsDir(r, texturesPath+"/items"):
		return LayoutLegacy
	case format >= firstFlattenedFormat:
		return LayoutFlattened
	}
	return LayoutLegacy
}

// optionalTexture reads a texture that packs may omit, logging only unexpected errors.
func optionalTexture(r fs.FS, p string, alphaFix bool) (image.Texture, bool) {
	t, err := image.NewTextureFS(r, p, alphaFix)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			logger.Debugf("error reading %s: %s", p, err)
		}
		return image.Texture{}, false
	}
	return t, true
}

// optionalDirectory reads every PNG below a folder that packs may omit.
func optionalDirectory(r fs.FS, p string, alphaFix bool) []image.Texture {
	if !fsutil.IsDir(r, p) {
		return nil
	}
	textures, err := image.DirectoryTexturesFS(r, p, alphaFix)
	if err != nil {
		logger.Debugf("error reading %s: %s", p, err)
	}
	return textures
}

// readArmors reads worn armour from the pre-1.21.2 models/armor folder and from the 1.21.2+
// entity/equipment folders, naming every layer <material>_layer_<n>.png.
func readArmors(r fs.FS, texturesPath string) []image.Texture {
	armors := optionalDirectory(r, texturesPath+"/models/armor", false)
	for dir, layer := range map[string]string{"humanoid": "1", "humanoid_leggings": "2"} {
		for _, t := range optionalDirectory(r, texturesPath+"/entity/equipment/"+dir, false) {
			material := strings.TrimSuffix(t.Name, ".png")
			suffix := ""
			if base, ok := strings.CutSuffix(material, "_overlay"); ok {
				material, suffix = base, "_overlay"
			}
			if material == "turtle_scute" {
				material = "turtle"
			}
			armors = append(armors, t.Rename(material+"_layer_"+layer+suffix+".png"))
		}
	}
	return armors
}

// readSprites reads the 1.20.2+ HUD sprites keyed by their path below dir without extension.
func readSprites(r fs.FS, dir string) map[string]image.Texture {
	if !fsutil.IsDir(r, dir) {
		return nil
	}
	sprites := map[string]image.Texture{}
	_ = fs.WalkDir(r, dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".png") {
			return nil
		}
		t, err := image.NewTextureFS(r, p, true)
		if err != nil {
			logger.Debugf("error reading %s: %s", p, err)
			return nil
		}
		key := strings.TrimSuffix(strings.TrimPrefix(p, dir+"/"), ".png")
		sprites[key] = t
		return nil
	})
	return sprites
}
