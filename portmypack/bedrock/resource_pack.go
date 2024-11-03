package bedrock

import (
	"archive/zip"
	"errors"
	"path/filepath"
	"strconv"

	"github.com/restartfu/portmypack/portmypack/image"
	"github.com/restartfu/portmypack/portmypack/internal/fsutil"
	"github.com/restartfu/portmypack/portmypack/internal/logger"
)

type ResourcePack struct {
	Name string
	r    *zip.Reader

	CubeMaps []image.Texture

	PackIcon  image.Texture
	Icons     image.Texture
	Particles image.Texture

	Items  []image.Texture
	Blocks []image.Texture
	Armors []image.Texture
}

func NewResourcePack(filePath string) (ResourcePack, error) {
	pck := ResourcePack{}
	pck.Name = filepath.Base(filePath)

	r, err := fsutil.OpenZip(filePath)
	if err != nil {
		return ResourcePack{}, errors.New("error opening bedrock pack zip: " + err.Error())
	}

	path, found := fsutil.FindDirectory(r, "overworld_cubemap")
	if found {
		for i := 0; i < 6; i++ {
			pck.CubeMaps[i], _ = image.NewTextureFS(r, path+"/cubemap_"+strconv.Itoa(i)+".png", false)
		}
	}

	pck.PackIcon, err = image.NewTextureFS(r, "pack_icon.png", false)
	if err != nil {
		logger.Debugf("error trying to find icon: %s", err)
	}

	texturesPath, found := fsutil.FindDirectory(r, "textures")
	if !found {
		return ResourcePack{}, errors.New("could not find textures path, stopping progress.")
	}

	pck.Icons, err = image.NewTextureFS(r, texturesPath+"/gui/icons.png", true)
	if err != nil {
		logger.Debugf("error trying to find icon: %s", err)
	}
	pck.Particles, err = image.NewTextureFS(r, texturesPath+"/particle/particles.png", false)
	if err != nil {
		logger.Debugf("error trying to find icon: %s", err)
	}

	pck.Items, err = image.DirectoryTexturesFS(r, texturesPath+"/items", true)
	if err != nil {
		logger.Debugf("error trying to find items: %s", err)
	}
	pck.Blocks, err = image.DirectoryTexturesFS(r, texturesPath+"/blocks", false)
	if err != nil {
		logger.Debugf("error trying to find items: %s", err)
	}
	pck.Armors, err = image.DirectoryTexturesFS(r, texturesPath+"/models/armor", false)
	if err != nil {
		logger.Debugf("error trying to find items: %s", err)
	}

	pck.r = r
	return pck, nil
}

func (pck ResourcePack) WriteZip(output string) error {
	f, w, err := fsutil.CreateZip(output)
	if err != nil {
		return errors.New("error creating new zip file: " + err.Error())
	}
	err = pck.generateManifest(w)
	if err != nil {
		return errors.New("error generating manifest: " + err.Error())
	}

	pck.PackIcon.Write(w, "pack_icon.png")
	writeTextures(w, pck.CubeMaps, "textures/environment/overworld_cubemap")
	pck.Icons.Write(w, "textures/gui/icons.png")
	pck.Particles.Write(w, "textures/particle/particles.png")
	writeTextures(w, pck.Items, "textures/items")
	writeTextures(w, pck.Blocks, "textures/blocks")
	writeTextures(w, pck.Armors, "textures/models/armor")

	w.Close()
	f.Close()
	return nil
}

func writeTextures(w *zip.Writer, textures []image.Texture, path string) {
	for _, t := range textures {
		if t == (image.Texture{}) {
			continue
		}
		output := filepath.Join(path, t.Name)
		t.Write(w, output)
	}
}
