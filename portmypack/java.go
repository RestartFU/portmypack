package portmypack

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path"
	"path/filepath"
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

	"copper_layer_1", "copper_1",
	"copper_layer_2", "copper_2",

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

	"turtle_layer_1", "turtle_1",
)

// Option configures a port.
type Option func(*bedrock.ResourcePack)

// WithMinEngineVersion sets the oldest Bedrock version allowed to load the ported pack. The default,
// bedrock.DefaultMinEngineVersion, keeps the pack loadable from Bedrock 1.13 up to current releases.
func WithMinEngineVersion(v [3]int) Option {
	return func(p *bedrock.ResourcePack) { p.MinEngineVersion = v }
}

// ConvertJavaEditionPack converts a Java Edition pack, in either the legacy (1.8-1.12) or the flattened
// (1.13 to current) layout, into a Bedrock Edition pack.
func ConvertJavaEditionPack(pack java.ResourcePack, opts ...Option) bedrock.ResourcePack {
	newPack := bedrock.ResourcePack{
		Name:        pack.Name,
		Description: pack.Description,
		Textures:    map[string]image.Texture{},
	}
	for _, opt := range opts {
		opt(&newPack)
	}

	pack.PackIcon.Name = "pack_icon.png"
	newPack.PackIcon = pack.PackIcon

	if len(pack.Skies) > 0 {
		newPack.CubeMaps, _ = bedrock.CubemapsFromTexture(pack.Skies[0])
	}

	set := func(p string, t image.Texture) {
		if t.Image != nil {
			newPack.Textures[p] = t.Rename(path.Base(p))
		}
	}
	setAll := func(textures map[string]image.Texture) {
		for p, t := range textures {
			set(p, t)
		}
	}

	set("textures/gui/icons.png", pack.Icons)
	set("textures/particle/particles.png", pack.Particles)

	for _, t := range pack.Blocks {
		setTexture(set, java.BedrockBlockPath(pack.Layout, t.Name), t)
	}
	items := map[string]image.Texture{}
	for _, t := range pack.Items {
		items[strings.TrimSuffix(t.Name, ".png")] = t
		setTexture(set, java.BedrockItemPath(pack.Layout, t.Name), t)
	}
	setAll(potionTextures(items, pack.Layout))
	setAll(atlasTextures(items))

	for _, a := range pack.Armors {
		set("textures/models/armor/"+bedrockReplacer.Replace(a.Name), a)
	}

	// HUD: crop legacy atlases first so that 1.20.2+ sprites, when present, take precedence.
	sprites := legacySprites(pack.Icons)
	for name, t := range pack.Sprites {
		sprites[name] = t
	}
	var hotbar, selection image.Texture
	if pack.Widgets.Image != nil {
		s := atlasScale(pack.Widgets)
		hotbar = legacyHotbarRect.crop(pack.Widgets, s, "hotbar.png")
		selection = legacySelectionRect.crop(pack.Widgets, s, "hotbar_selection.png")
	}
	if t, ok := pack.Sprites["hotbar"]; ok {
		hotbar = t
	}
	if t, ok := pack.Sprites["hotbar_selection"]; ok {
		selection = t
	}
	setAll(hudTextures(sprites, hotbar, selection))
	if pack.Icons.Image == nil {
		if icons, ok := iconsFromSprites(pack.Sprites); ok {
			set("textures/gui/icons.png", icons)
		}
	}
	return newPack
}

// setTexture stores a block or item texture at its Bedrock path (relative to textures/, without
// extension). Animation strips are cut to their first frame unless Bedrock animates that texture.
func setTexture(set func(string, image.Texture), bedrockPath string, t image.Texture) {
	if !bedrock.IsFlipbook(bedrockPath) {
		t = t.FirstFrame()
	}
	set("textures/"+bedrockPath+".png", t)
}

// PortJavaEditionPackAndExtract ports a Java Edition pack, writes it to tmp/portmypack-<n>.mcpack and
// extracts it to outputDirector/portmypack-<n>.
func PortJavaEditionPackAndExtract(pack java.ResourcePack, outputDirector string, opts ...Option) error {
	tmp := "portmypack-" + strconv.Itoa(int(rand.IntN(99999)))
	tmpPath := filepath.Join("tmp", tmp+".mcpack")

	if err := os.MkdirAll("tmp", os.ModePerm); err != nil {
		return err
	}
	if err := PortJavaEditionPack(pack, tmpPath, opts...); err != nil {
		return err
	}

	output := filepath.Join(outputDirector, tmp)
	if err := fsutil.Unzip(tmpPath, output); err != nil {
		return err
	}
	fmt.Println("extracted mcpack to:", output)
	return nil
}

// PortJavaEditionPack ports a Java Edition pack and writes it as a .mcpack archive at output.
func PortJavaEditionPack(pack java.ResourcePack, output string, opts ...Option) error {
	newPack := ConvertJavaEditionPack(pack, opts...)
	if err := newPack.WriteZip(output); err != nil {
		return err
	}
	fmt.Println("mcpack file written to:", output)
	return nil
}
