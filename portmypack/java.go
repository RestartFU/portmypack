package portmypack

import (
	"fmt"
	img "image"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"

	"github.com/restartfu/portmypack/portmypack/bedrock"
	"github.com/restartfu/portmypack/portmypack/image"
	"github.com/restartfu/portmypack/portmypack/internal/fsutil"
	"github.com/restartfu/portmypack/portmypack/java"
)

// bedrockReplacer replaces the texture names with the correct names.
var bedrockReplacer = strings.NewReplacer(
	"chainmail_layer_1", "chain_1",
	"chainmail_layer_2", "chain_2",

	"diamond_layer_1", "diamond_1",
	"diamond_layer_2", "diamond_2",

	"gold_layer_1", "gold_1",
	"gold_layer_2", "gold_2",

	"iron_layer_1", "iron_1",
	"iron_layer_2", "iron_2",

	"leather_layer_1", "leather_1",
	"leather_layer_2", "leather_2",

	"netherite_layer_1", "netherite_1",
	"netherite_layer_2", "netherite_2",
)

func PortJavaEditionPack(pack java.ResourcePack, outputDirector string) {
	newPack := bedrock.ResourcePack{}
	newPack.Name = pack.Name

	pack.PackIcon.Name = "pack_icon.png"
	newPack.PackIcon = pack.PackIcon

	newPack.Icons = pack.Icons
	newPack.Particles = pack.Particles

	newPack.Items = pack.Items
	newPack.Blocks = pack.Blocks

	var newArmors []image.Texture
	for _, a := range pack.Armors {
		a.Name = bedrockReplacer.Replace(a.Name)
		newArmors = append(newArmors, a)
	}
	newPack.Armors = newArmors
	newPack.CubeMaps, _ = splitCubeMaps(pack.Skies[0])

	tmp := "portmypack-" + strconv.Itoa(int(rand.IntN(99999)))
	tmpPath := "tmp/" + tmp + ".mcpack"

	os.Mkdir("tmp", os.ModePerm)
	newPack.WriteZip(tmpPath)
	fmt.Println("mcpack file written to:", tmpPath)

	output := outputDirector + "/" + tmp
	fsutil.Unzip(tmpPath, output)
	fmt.Println("extracted mcpack to:", output)
}

func splitCubeMaps(cubemaps image.Texture) ([]image.Texture, error) {
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
		textures[i] = texture
	}

	return textures, nil
}
