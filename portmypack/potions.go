package portmypack

import (
	img "image"
	"strconv"

	"github.com/restartfu/portmypack/portmypack/image"
	"github.com/restartfu/portmypack/portmypack/java"
)

// potionBottles names the Java bottle textures for each Bedrock bottle kind, per layout.
var potionBottles = map[java.Layout]map[string]string{
	java.LayoutLegacy:    {"": "potion_bottle_drinkable", "splash_": "potion_bottle_splash", "lingering_": "potion_bottle_lingering"},
	java.LayoutFlattened: {"": "potion", "splash_": "splash_potion", "lingering_": "lingering_potion"},
}

// potionTextures builds Bedrock's pre-coloured potion textures. Java draws a potion as its bottle over a
// tinted potion_overlay, while Bedrock ships one texture per effect, so each is composed the same way.
// Nothing is produced for a bottle kind unless the pack overrides both that bottle and the overlay.
func potionTextures(items map[string]image.Texture, layout java.Layout) map[string]image.Texture {
	out := map[string]image.Texture{}
	overlay, ok := items["potion_overlay"]
	if !ok {
		return out
	}
	for kind, bottleName := range potionBottles[layout] {
		bottle, ok := items[bottleName]
		if !ok || bottle.Bounds().Size() != overlay.Bounds().Size() {
			continue
		}
		for effect, tint := range potionTints {
			if (kind == "" && !tint.drink) || (kind == "splash_" && !tint.splash) || (kind == "lingering_" && !tint.lingering) {
				continue
			}
			name := "potion_bottle_" + kind + effect + ".png"
			out["textures/items/"+name] = overlay.Tint(tint.color).Over(bottle, name)
		}
	}
	return out
}

// itemAtlases lists the Java animation frames Bedrock stores as one vertical atlas.
var itemAtlases = []struct {
	prefix, atlas string
	frames        int
}{
	{"clock_", "watch_atlas", 64},
	{"compass_", "compass_atlas", 32},
	{"recovery_compass_", "recovery_compass_atlas", 32},
}

// atlasTextures stitches clock and compass frames into the Bedrock atlases when the pack overrides every
// frame at one size; a partial set would leave holes in the atlas.
func atlasTextures(items map[string]image.Texture) map[string]image.Texture {
	out := map[string]image.Texture{}
	for _, a := range itemAtlases {
		var frames []image.Texture
		for i := 0; i < a.frames; i++ {
			n := strconv.Itoa(i)
			if i < 10 {
				n = "0" + n
			}
			t, ok := items[a.prefix+n]
			if !ok || (len(frames) > 0 && t.Bounds().Size() != frames[0].Bounds().Size()) {
				frames = nil
				break
			}
			frames = append(frames, t)
		}
		if frames == nil {
			continue
		}
		w, h := frames[0].Bounds().Dx(), frames[0].Bounds().Dy()
		c := image.NewCanvas(w, h*len(frames))
		for i, f := range frames {
			c.Paste(f, img.Pt(0, i*h))
		}
		out["textures/items/"+a.atlas+".png"] = c.Texture(a.atlas + ".png")
	}
	return out
}
