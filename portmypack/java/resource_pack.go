package java

import (
	"archive/zip"
	"errors"
	"path/filepath"

	"github.com/restartfu/portmypack/portmypack/image"
	"github.com/restartfu/portmypack/portmypack/internal/fsutil"
	"github.com/restartfu/portmypack/portmypack/internal/logger"
)

type ResourcePack struct {
	Name string
	r    *zip.Reader

	Skies []image.Texture

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
		return ResourcePack{}, errors.New("error opening java pack zip: " + err.Error())
	}

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
