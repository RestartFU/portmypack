package bedrock

import (
	"strconv"

	img "image"

	"github.com/restartfu/portmypack/portmypack/image"
)

func CubemapsFromTexture(cubemaps image.Texture) ([]image.Texture, error) {
	textures := make([]image.Texture, 6)

	cubeMapHeight := cubemaps.Bounds().Dy() / 2
	cubeMapWidth := cubemaps.Bounds().Dx() / 3

	for i := 0; i < 6; i++ {
		rgba := img.NewRGBA(img.Rect(0, 0, cubeMapWidth, cubeMapHeight))

		for y := 0; y < cubeMapHeight; y++ {
			for x := 0; x < cubeMapWidth; x++ {
				rgba.Set(x, y, cubemaps.At(x+(i%3)*cubeMapWidth, y+(i/3)*cubeMapHeight))
			}
		}

		var rotation int
		switch i {
		case 0:
			rotation = 5
		case 1:
			rotation = 4
		case 2:
			rotation = 2
		case 3:
			rotation = 3
		case 4:
			rotation = 0
		case 5:
			rotation = 1
		}

		texture := image.Texture{
			Image: rgba,
			Name:  "cubemap_" + strconv.Itoa(rotation) + ".png",
		}
		textures[rotation] = texture
	}

	return textures, nil
}
