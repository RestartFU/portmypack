package java

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/magiconair/properties"
	"github.com/restartfu/portmypack/portmypack/image"
)

type skyProperty struct {
	Rotate bool   `properties:"rotate"`
	Source string `properties:"source"`
}

func findSkies(f fs.FS, path string) ([]image.Texture, error) {
	var skyProperties []skyProperty

	walkDirFunc := func(filePath string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(filePath, ".properties") {
			f, err := f.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()

			buf, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			p, err := properties.Load(buf, properties.UTF8)
			var sky skyProperty
			err = p.Decode(&sky)
			if err != nil {
				return err
			}

			skyProperties = append(skyProperties, sky)
			return nil
		}
		return nil
	}

	err := fs.WalkDir(f, path, walkDirFunc)
	replacer := strings.NewReplacer("./", "")

	var skies []image.Texture
	for _, sky := range skyProperties {
		if sky.Source == "" {
			continue
		}
		cubemaps, err := f.Open(path + "/" + replacer.Replace(sky.Source))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}
		defer cubemaps.Close()
		img, err := image.NewTexture(sky.Source, cubemaps, false)
		if err != nil {
			return nil, err
		}
		skies = append(skies, img)
	}
	return skies, err
}
