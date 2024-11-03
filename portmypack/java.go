package portmypack

import (
	"fmt"
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

func PortJavaEditionPackAndExtract(pack java.ResourcePack, outputDirector string) {
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
	newPack.CubeMaps, _ = bedrock.CubemapsFromTexture(pack.Skies[0])

	tmp := "portmypack-" + strconv.Itoa(int(rand.IntN(99999)))
	tmpPath := "tmp/" + tmp + ".mcpack"

	os.Mkdir("tmp", os.ModePerm)
	newPack.WriteZip(tmpPath)
	fmt.Println("mcpack file written to:", tmpPath)

	output := outputDirector + "/" + tmp
	fsutil.Unzip(tmpPath, output)
	fmt.Println("extracted mcpack to:", output)
}

func PortJavaEditionPack(pack java.ResourcePack, output string) {
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
	newPack.CubeMaps, _ = bedrock.CubemapsFromTexture(pack.Skies[0])

	os.Mkdir("tmp", os.ModePerm)
	newPack.WriteZip(output)
	fmt.Println("mcpack file written to:", output)
}
