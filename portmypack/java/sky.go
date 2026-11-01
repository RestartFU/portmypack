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

// skyProperty is the part of an OptiFine sky .properties file portmypack reads. Both keys are optional.
type skyProperty struct {
	Rotate bool   `properties:"rotate,default=true"`
	Source string `properties:"source,default="`
}

// skySourcePath resolves an OptiFine sky "source" relative to the sky folder, accepting "./x.png",
// "x.png", "optifine/sky/world0/x.png" and namespaced "minecraft:optifine/sky/world0/x.png".
func skySourcePath(dir, source string) string {
	source = strings.TrimPrefix(source, "minecraft:")
	if rest, ok := strings.CutPrefix(source, "./"); ok {
		return dir + "/" + rest
	}
	if strings.Contains(source, "/") {
		if i := strings.Index(dir, "assets/minecraft/"); i >= 0 {
			return dir[:i] + "assets/minecraft/" + source
		}
		return source
	}
	return dir + "/" + source
}

// findSkies reads the cube-map images referenced by the OptiFine sky properties in path.
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
			if err != nil {
				return err
			}
			var sky skyProperty
			if err = p.Decode(&sky); err != nil {
				return err
			}

			skyProperties = append(skyProperties, sky)
			return nil
		}
		return nil
	}

	err := fs.WalkDir(f, path, walkDirFunc)

	var skies []image.Texture
	for _, sky := range skyProperties {
		if sky.Source == "" {
			continue
		}
		cubemaps, err := f.Open(skySourcePath(path, sky.Source))
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
