package bedrock

import (
	"strconv"

	img "image"

	"github.com/restartfu/portmypack/portmypack/image"
)

var cubemapRotations = []int{5, 4, 2, 3, 0, 1}

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
		rotation := cubemapRotations[i]

		texture := image.Texture{
			Image: rgba,
			Name:  "cubemap_" + strconv.Itoa(rotation) + ".png",
		}
		textures[rotation] = texture
	}

	return textures, nil
}
