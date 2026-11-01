package bedrock

import (
	"errors"
	"io/fs"
	"path"
	"sort"
	"strconv"

	"github.com/restartfu/portmypack/portmypack/image"
	"github.com/restartfu/portmypack/portmypack/internal/fsutil"
	"github.com/restartfu/portmypack/portmypack/internal/logger"
)

// ResourcePack is a Bedrock Edition resource pack.
type ResourcePack struct {
	Name string
	// Description is written to the manifest; a default is used when it is empty.
	Description string
	// MinEngineVersion is written to the manifest; DefaultMinEngineVersion is used when it is zero.
	MinEngineVersion [3]int

	r fs.FS

	CubeMaps []image.Texture

	PackIcon  image.Texture
	Icons     image.Texture
	Particles image.Texture

	Items  []image.Texture
	Blocks []image.Texture
	Armors []image.Texture

	// Textures holds any other texture keyed by its path inside the pack, such as
	// "textures/ui/heart.png". Entries here replace the same path from the fields above.
	Textures map[string]image.Texture
}

// NewResourcePack reads a Bedrock Edition resource pack from a .mcpack/.zip archive or a folder.
func NewResourcePack(filePath string) (ResourcePack, error) {
	pck := ResourcePack{}
	pck.Name = fsutil.PackName(filePath)

	archive, err := fsutil.OpenPack(filePath)
	if err != nil {
		return ResourcePack{}, errors.New("error opening bedrock pack: " + err.Error())
	}
	r := archive
	if root, found := fsutil.FindPackRoot(archive, "manifest.json"); found && root != "." {
		if r, err = fs.Sub(archive, root); err != nil {
			return ResourcePack{}, err
		}
	}

	path, found := fsutil.FindDirectory(r, "overworld_cubemap")
	if found {
		pck.CubeMaps = make([]image.Texture, 6)
		for i := 0; i < 6; i++ {
			pck.CubeMaps[i], _ = image.NewTextureFS(r, path+"/cubemap_"+strconv.Itoa(i)+".png", false)
		}
	}

	pck.PackIcon, err = image.NewTextureFS(r, "pack_icon.png", false)
	if err != nil {
		logger.Debugf("error trying to find icon: %s", err)
	}

	texturesPath := "textures"
	if !fsutil.IsDir(r, texturesPath) {
		return ResourcePack{}, errors.New("could not find textures path, stopping progress")
	}

	pck.Icons, err = image.NewTextureFS(r, texturesPath+"/gui/icons.png", true)
	if err != nil {
		logger.Debugf("error trying to find icon: %s", err)
	}
	pck.Particles, err = image.NewTextureFS(r, texturesPath+"/particle/particles.png", false)
	if err != nil {
		logger.Debugf("error trying to find particles: %s", err)
	}

	pck.Items, err = image.DirectoryTexturesFS(r, texturesPath+"/items", true)
	if err != nil {
		logger.Debugf("error trying to find items: %s", err)
	}
	pck.Blocks, err = image.DirectoryTexturesFS(r, texturesPath+"/blocks", false)
	if err != nil {
		logger.Debugf("error trying to find blocks: %s", err)
	}
	pck.Armors, err = image.DirectoryTexturesFS(r, texturesPath+"/models/armor", false)
	if err != nil {
		logger.Debugf("error trying to find armors: %s", err)
	}

	pck.r = r
	return pck, nil
}

// IsFlipbook reports whether vanilla Bedrock animates the texture at p (relative to textures/, without
// extension, such as "blocks/fire_0") as a vertical strip of frames.
func IsFlipbook(p string) bool {
	_, ok := flipbooks[p]
	return ok
}

// Files returns every texture of the pack keyed by its path inside the pack.
func (pck ResourcePack) Files() map[string]image.Texture {
	files := map[string]image.Texture{}
	add := func(p string, t image.Texture) {
		if t.Image != nil {
			files[p] = t
		}
	}
	addAll := func(dir string, textures []image.Texture) {
		for _, t := range textures {
			add(path.Join(dir, t.Name), t)
		}
	}
	add("pack_icon.png", pck.PackIcon)
	addAll("textures/environment/overworld_cubemap", pck.CubeMaps)
	add("textures/gui/icons.png", pck.Icons)
	add("textures/particle/particles.png", pck.Particles)
	addAll("textures/items", pck.Items)
	addAll("textures/blocks", pck.Blocks)
	addAll("textures/models/armor", pck.Armors)
	for p, t := range pck.Textures {
		add(p, t)
	}
	return files
}

// WriteZip writes the pack, with a generated manifest, as a zip archive such as a .mcpack file.
func (pck ResourcePack) WriteZip(output string) error {
	f, w, err := fsutil.CreateZip(output)
	if err != nil {
		return errors.New("error creating new zip file: " + err.Error())
	}
	defer f.Close()

	if err = pck.generateManifest(w); err != nil {
		return errors.New("error generating manifest: " + err.Error())
	}

	files := pck.Files()
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		if err := files[p].Write(w, p); err != nil {
			return errors.New("error writing " + p + ": " + err.Error())
		}
	}

	if err := w.Close(); err != nil {
		return err
	}
	return f.Close()
}
