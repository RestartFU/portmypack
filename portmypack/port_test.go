package portmypack

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/restartfu/portmypack/portmypack/bedrock"
	"github.com/restartfu/portmypack/portmypack/java"
)

// solid returns a w x h PNG filled with c.
func solid(w, h int, c color.NRGBA) []byte {
	m := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			m.SetNRGBA(x, y, c)
		}
	}
	return encode(m)
}

// encode encodes m as PNG.
func encode(m image.Image) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// markedAtlas returns a size x size image, transparent except for the given regions, each filled with
// its colour. Regions are given in 256-pixel atlas coordinates and scaled with the atlas.
func markedAtlas(size int, regions map[rect]color.NRGBA) []byte {
	s := size / 256
	m := image.NewNRGBA(image.Rect(0, 0, size, size))
	for r, c := range regions {
		for y := r.y * s; y < (r.y+r.h)*s; y++ {
			for x := r.x * s; x < (r.x+r.w)*s; x++ {
				m.SetNRGBA(x, y, c)
			}
		}
	}
	return encode(m)
}

// writeZip writes files into a new zip archive at p.
func writeZip(t *testing.T, p string, files map[string][]byte) {
	t.Helper()
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	for name, data := range files {
		fw, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fw.Write(data); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

// writeDir writes files below dir.
func writeDir(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()
	for name, data := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// readMcpack returns the files of a written .mcpack.
func readMcpack(t *testing.T, p string) map[string][]byte {
	t.Helper()
	r, err := zip.OpenReader(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	out := map[string][]byte{}
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		out[f.Name], _ = io.ReadAll(rc)
		rc.Close()
	}
	return out
}

// pixel decodes a PNG from files and returns its size and the colour at (x, y).
func pixel(t *testing.T, files map[string][]byte, name string, x, y int) (image.Point, color.NRGBA) {
	t.Helper()
	data, ok := files[name]
	if !ok {
		t.Fatalf("%s missing from pack", name)
	}
	m, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return m.Bounds().Size(), color.NRGBAModel.Convert(m.At(x, y)).(color.NRGBA)
}

// port reads a Java pack at in and returns the files of the ported .mcpack.
func port(t *testing.T, in string, opts ...Option) (java.ResourcePack, map[string][]byte) {
	t.Helper()
	pack, err := java.NewResourcePack(in)
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "out.mcpack")
	if err := PortJavaEditionPack(pack, out, opts...); err != nil {
		t.Fatal(err)
	}
	return pack, readMcpack(t, out)
}

var (
	red    = color.NRGBA{255, 0, 0, 255}
	green  = color.NRGBA{0, 255, 0, 255}
	blue   = color.NRGBA{0, 0, 255, 255}
	white  = color.NRGBA{255, 255, 255, 255}
	yellow = color.NRGBA{255, 255, 0, 255}
)

func TestPortLegacyPack(t *testing.T) {
	// A 1.8-style pack with an HD (2x) icons.png.
	icons := markedAtlas(512, map[rect]color.NRGBA{
		legacyIconRects["crosshair"]:  white,
		legacyIconRects["heart/full"]: red,
		legacyIconRects["food_full"]:  yellow,
	})
	widgets := markedAtlas(256, map[rect]color.NRGBA{
		{0, 0, 1, 22}:    red,
		{1, 0, 20, 22}:   green,
		{161, 0, 20, 22}: blue,
		{0, 22, 24, 24}:  white,
	})
	overlay := solid(16, 16, white)
	bottle := image.NewNRGBA(image.Rect(0, 0, 16, 16)) // transparent bottle: overlay shows through
	files := map[string][]byte{
		"pack.mcmeta": []byte(`{"pack":{"pack_format":1,"description":"§6Legacy §rPvP"}}`),
		"pack.png":    solid(16, 16, blue),
		"assets/minecraft/textures/blocks/planks_oak.png":              solid(16, 16, green),
		"assets/minecraft/textures/blocks/destroy_stage_3.png":         solid(16, 16, green),
		"assets/minecraft/textures/items/fish_cod_raw.png":             solid(16, 16, blue),
		"assets/minecraft/textures/items/diamond_sword.png":            solid(16, 16, blue),
		"assets/minecraft/textures/items/potion_overlay.png":           overlay,
		"assets/minecraft/textures/items/potion_bottle_splash.png":     encode(bottle),
		"assets/minecraft/textures/models/armor/diamond_layer_1.png":   solid(64, 32, blue),
		"assets/minecraft/textures/models/armor/chainmail_layer_2.png": solid(64, 32, red),
		"assets/minecraft/textures/gui/icons.png":                      icons,
		"assets/minecraft/textures/gui/widgets.png":                    widgets,
		"assets/minecraft/textures/particle/particles.png":             solid(128, 128, red),
		"assets/minecraft/mcpatcher/sky/world0/sky1.properties":        []byte("source=./sky1.png\n"),
		"assets/minecraft/mcpatcher/sky/world0/sky1.png":               solid(48, 32, blue),
	}
	in := filepath.Join(t.TempDir(), "Legacy PvP.zip")
	writeZip(t, in, files)

	pack, out := port(t, in)
	if pack.Layout != java.LayoutLegacy || pack.Format != 1 {
		t.Fatalf("layout %v format %d, want legacy 1", pack.Layout, pack.Format)
	}

	var manifest bedrock.Manifest
	if err := json.Unmarshal(out["manifest.json"], &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.FormatVersion != 2 || manifest.Header.MinEngineVersion != bedrock.DefaultMinEngineVersion {
		t.Fatalf("manifest format %d engine %v", manifest.FormatVersion, manifest.Header.MinEngineVersion)
	}
	if manifest.Header.Name != "Legacy PvP" || manifest.Header.Description != "Legacy PvP" {
		t.Fatalf("manifest header %+v", manifest.Header)
	}

	for _, name := range []string{
		"pack_icon.png",
		"textures/blocks/planks_oak.png",
		"textures/environment/destroy_stage_3.png",
		"textures/items/fish_raw.png",
		"textures/items/diamond_sword.png",
		"textures/models/armor/diamond_1.png",
		"textures/models/armor/chain_2.png",
		"textures/gui/icons.png",
		"textures/particle/particles.png",
		"textures/environment/overworld_cubemap/cubemap_0.png",
	} {
		if _, ok := out[name]; !ok {
			t.Errorf("%s missing", name)
		}
	}

	// Splash heal potion: transparent bottle over a white overlay tinted with the heal colour.
	_, c := pixel(t, out, "textures/items/potion_bottle_splash_heal.png", 8, 8)
	if want := potionTints["heal"].color; c != want {
		t.Errorf("splash heal pixel %v, want %v", c, want)
	}
	if _, ok := out["textures/items/potion_bottle_heal.png"]; ok {
		t.Error("drinkable potions generated without a drinkable bottle")
	}

	// HUD icons cut from the 2x atlas.
	size, c := pixel(t, out, "textures/ui/heart.png", 0, 0)
	if size != image.Pt(18, 18) || c != red {
		t.Errorf("heart %v %v", size, c)
	}
	if _, c = pixel(t, out, "textures/ui/hunger_full.png", 5, 5); c != yellow {
		t.Errorf("hunger_full %v", c)
	}
	if _, ok := out["textures/ui/freeze_heart.png"]; ok {
		t.Error("empty atlas region exported")
	}

	// Hotbar sliced from widgets.png.
	if size, c = pixel(t, out, "textures/ui/hotbar_start_cap.png", 0, 0); size != image.Pt(1, 22) || c != red {
		t.Errorf("start cap %v %v", size, c)
	}
	if size, c = pixel(t, out, "textures/ui/hotbar_0.png", 0, 0); size != image.Pt(20, 22) || c != green {
		t.Errorf("hotbar_0 %v %v", size, c)
	}
	if _, c = pixel(t, out, "textures/ui/hotbar_8.png", 19, 21); c != blue {
		t.Errorf("hotbar_8 %v", c)
	}
	if size, _ = pixel(t, out, "textures/ui/selected_hotbar_slot.png", 0, 0); size != image.Pt(24, 24) {
		t.Errorf("selection %v", size)
	}
}

func TestPortCurrentPack(t *testing.T) {
	sprites := "assets/minecraft/textures/gui/sprites/hud/"
	selection := image.NewNRGBA(image.Rect(0, 0, 24, 23))
	for x := 0; x < 24; x++ {
		selection.SetNRGBA(x, 0, white)
		selection.SetNRGBA(x, 22, red)
	}
	strip := image.NewNRGBA(image.Rect(0, 0, 16, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 16; x++ {
			c := red
			if y >= 16 {
				c = blue
			}
			strip.SetNRGBA(x, y, c)
		}
	}
	files := map[string][]byte{
		"pack.mcmeta": []byte(`{"pack":{"description":{"text":"Current ","extra":[{"text":"Pack"}]},"min_format":[97,1],"max_format":97}}`),
		"pack.png":    solid(16, 16, blue),
		"assets/minecraft/textures/block/oak_planks.png":                             solid(16, 16, green),
		"assets/minecraft/textures/block/oak_leaves.png":                             solid(16, 16, green),
		"assets/minecraft/textures/block/fire_0.png":                                 encode(strip),
		"assets/minecraft/textures/block/sea_lantern.png":                            encode(strip),
		"assets/minecraft/textures/block/stone.png":                                  encode(strip),
		"assets/minecraft/textures/item/golden_apple.png":                            solid(16, 16, yellow),
		"assets/minecraft/textures/item/diamond_sword.png":                           encode(strip),
		"assets/minecraft/textures/item/potion.png":                                  encode(image.NewNRGBA(image.Rect(0, 0, 16, 16))),
		"assets/minecraft/textures/item/potion_overlay.png":                          solid(16, 16, white),
		"assets/minecraft/textures/entity/equipment/humanoid/diamond.png":            solid(64, 32, blue),
		"assets/minecraft/textures/entity/equipment/humanoid/turtle_scute.png":       solid(64, 32, green),
		"assets/minecraft/textures/entity/equipment/humanoid_leggings/chainmail.png": solid(64, 32, red),
		sprites + "heart/full.png":                                                   solid(9, 9, red),
		sprites + "heart/hardcore_full.png":                                          solid(9, 9, yellow),
		sprites + "food_full_hunger.png":                                             solid(9, 9, green),
		sprites + "air.png":                                                          solid(9, 9, blue),
		sprites + "crosshair.png":                                                    solid(15, 15, white),
		sprites + "hotbar.png":                                                       solid(182, 22, green),
		sprites + "hotbar_selection.png":                                             encode(selection),
	}
	for i := 0; i < 64; i++ {
		name := "assets/minecraft/textures/item/clock_" + string(rune('0'+i/10)) + string(rune('0'+i%10)) + ".png"
		files[name] = solid(16, 16, color.NRGBA{uint8(i * 4), 0, 0, 255})
	}
	// Zipped together with its enclosing folder, as many downloaded packs are.
	nested := map[string][]byte{}
	for name, data := range files {
		nested["Current Pack/"+name] = data
	}
	in := filepath.Join(t.TempDir(), "current.zip")
	writeZip(t, in, nested)

	pack, out := port(t, in, WithMinEngineVersion([3]int{1, 21, 0}))
	if pack.Layout != java.LayoutFlattened || pack.Format != 97 || pack.Description != "Current Pack" {
		t.Fatalf("layout %v format %d description %q", pack.Layout, pack.Format, pack.Description)
	}
	var manifest bedrock.Manifest
	if err := json.Unmarshal(out["manifest.json"], &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Header.MinEngineVersion != [3]int{1, 21, 0} {
		t.Errorf("min engine %v", manifest.Header.MinEngineVersion)
	}

	for _, name := range []string{
		"textures/blocks/planks_oak.png",
		"textures/blocks/leaves_oak.png",
		"textures/items/apple_golden.png",
		"textures/models/armor/diamond_1.png",
		"textures/models/armor/turtle_1.png",
		"textures/models/armor/chain_2.png",
		"textures/ui/hardcore/heart.png",
		"textures/ui/hunger_effect_full.png",
		"textures/ui/bubble.png",
		"textures/ui/hotbar_4.png",
		"textures/items/potion_bottle_heal.png",
		"textures/items/potion_bottle_drinkable.png",
	} {
		if _, ok := out[name]; !ok {
			t.Errorf("%s missing", name)
		}
	}
	for name := range out {
		if strings.Contains(name, "oak_planks") || strings.Contains(name, "golden_apple") {
			t.Errorf("unmapped Java name %s written", name)
		}
	}

	// Bedrock animates fire_0 and sea_lantern but not stone or items: other strips keep frame one.
	for name, h := range map[string]int{
		"textures/blocks/fire_0.png":       64,
		"textures/blocks/sea_lantern.png":  64,
		"textures/blocks/stone.png":        16,
		"textures/items/diamond_sword.png": 16,
	} {
		if size, _ := pixel(t, out, name, 0, 0); size.Y != h {
			t.Errorf("%s height %d, want %d", name, size.Y, h)
		}
	}

	if _, c := pixel(t, out, "textures/ui/heart.png", 4, 4); c != red {
		t.Errorf("heart %v", c)
	}
	size, c := pixel(t, out, "textures/ui/selected_hotbar_slot.png", 0, 23)
	if size != image.Pt(24, 24) || c != white {
		t.Errorf("selection %v bottom row %v, want 24x24 and the top border colour", size, c)
	}
	// Crosshair placed into a rebuilt icons.png; other regions stay empty.
	if size, c = pixel(t, out, "textures/gui/icons.png", 7, 7); size != image.Pt(256, 256) || c != white {
		t.Errorf("icons crosshair %v %v", size, c)
	}
	if _, c = pixel(t, out, "textures/gui/icons.png", 55, 3); c != red {
		t.Errorf("icons heart %v", c)
	}
	if size, c = pixel(t, out, "textures/items/watch_atlas.png", 0, 16*10); size != image.Pt(16, 1024) || c.R != 40 {
		t.Errorf("watch atlas %v %v", size, c)
	}
}

func TestBedrockReadBack(t *testing.T) {
	dir := t.TempDir()
	writeDir(t, dir, map[string][]byte{
		"pack.mcmeta": []byte(`{"pack":{"pack_format":34,"description":"dir"}}`),
		"assets/minecraft/textures/block/oak_planks.png":       solid(16, 16, green),
		"assets/minecraft/optifine/sky/world0/sky1.properties": []byte("source=minecraft:optifine/sky/world0/sky.png\n"),
		"assets/minecraft/optifine/sky/world0/sky.png":         solid(48, 32, blue),
	})
	pack, err := java.NewResourcePack(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(pack.Skies) != 1 {
		t.Fatalf("skies %d, want 1", len(pack.Skies))
	}
	out := filepath.Join(t.TempDir(), "dir.mcpack")
	if err := PortJavaEditionPack(pack, out); err != nil {
		t.Fatal(err)
	}
	b, err := bedrock.NewResourcePack(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(b.Blocks) != 1 || b.Blocks[0].Name != "planks_oak.png" || len(b.CubeMaps) != 6 {
		t.Fatalf("read back %d blocks, %d cubemaps", len(b.Blocks), len(b.CubeMaps))
	}
}
