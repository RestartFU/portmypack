package portmypack

import (
	img "image"
	"strconv"

	"github.com/restartfu/portmypack/portmypack/image"
)

// spriteUI maps Java 1.20.2+ HUD sprites (paths below gui/sprites/hud) to the individual textures the
// Bedrock HUD draws, which replaced gui/icons.png for hearts, hunger, armour and air.
var spriteUI = map[string]string{
	"heart/container":          "textures/ui/heart_background.png",
	"heart/container_blinking": "textures/ui/heart_blink.png",
	"heart/full":               "textures/ui/heart.png",
	"heart/half":               "textures/ui/heart_half.png",
	"heart/full_blinking":      "textures/ui/heart_flash.png",
	"heart/half_blinking":      "textures/ui/heart_flash_half.png",

	"heart/poisoned_full":          "textures/ui/poison_heart.png",
	"heart/poisoned_half":          "textures/ui/poison_heart_half.png",
	"heart/poisoned_full_blinking": "textures/ui/poison_heart_flash.png",
	"heart/poisoned_half_blinking": "textures/ui/poison_heart_flash_half.png",

	"heart/withered_full":          "textures/ui/wither_heart.png",
	"heart/withered_half":          "textures/ui/wither_heart_half.png",
	"heart/withered_full_blinking": "textures/ui/wither_heart_flash.png",
	"heart/withered_half_blinking": "textures/ui/wither_heart_flash_half.png",

	"heart/absorbing_full": "textures/ui/absorption_heart.png",
	"heart/absorbing_half": "textures/ui/absorption_heart_half.png",

	"heart/frozen_full":          "textures/ui/freeze_heart.png",
	"heart/frozen_half":          "textures/ui/freeze_heart_half.png",
	"heart/frozen_full_blinking": "textures/ui/freeze_heart_flash.png",
	"heart/frozen_half_blinking": "textures/ui/freeze_heart_flash_half.png",

	"heart/hardcore_full":                   "textures/ui/hardcore/heart.png",
	"heart/hardcore_half":                   "textures/ui/hardcore/heart_half.png",
	"heart/hardcore_full_blinking":          "textures/ui/hardcore/heart_flash.png",
	"heart/hardcore_half_blinking":          "textures/ui/hardcore/heart_flash_half.png",
	"heart/poisoned_hardcore_full":          "textures/ui/hardcore/poison_heart.png",
	"heart/poisoned_hardcore_half":          "textures/ui/hardcore/poison_heart_half.png",
	"heart/poisoned_hardcore_full_blinking": "textures/ui/hardcore/poison_heart_flash.png",
	"heart/poisoned_hardcore_half_blinking": "textures/ui/hardcore/poison_heart_flash_half.png",
	"heart/withered_hardcore_full":          "textures/ui/hardcore/wither_heart.png",
	"heart/withered_hardcore_half":          "textures/ui/hardcore/wither_heart_half.png",
	"heart/withered_hardcore_full_blinking": "textures/ui/hardcore/wither_heart_flash.png",
	"heart/withered_hardcore_half_blinking": "textures/ui/hardcore/wither_heart_flash_half.png",
	"heart/absorbing_hardcore_full":         "textures/ui/hardcore/absorption_heart.png",
	"heart/absorbing_hardcore_half":         "textures/ui/hardcore/absorption_heart_half.png",
	"heart/frozen_hardcore_full":            "textures/ui/hardcore/freeze_heart.png",
	"heart/frozen_hardcore_half":            "textures/ui/hardcore/freeze_heart_half.png",
	"heart/frozen_hardcore_full_blinking":   "textures/ui/hardcore/freeze_heart_flash.png",
	"heart/frozen_hardcore_half_blinking":   "textures/ui/hardcore/freeze_heart_flash_half.png",

	"food_empty":        "textures/ui/hunger_background.png",
	"food_full":         "textures/ui/hunger_full.png",
	"food_half":         "textures/ui/hunger_half.png",
	"food_empty_hunger": "textures/ui/hunger_effect_background.png",
	"food_full_hunger":  "textures/ui/hunger_effect_full.png",
	"food_half_hunger":  "textures/ui/hunger_effect_half.png",

	"armor_empty": "textures/ui/armor_empty.png",
	"armor_full":  "textures/ui/armor_full.png",
	"armor_half":  "textures/ui/armor_half.png",

	"air":          "textures/ui/bubble.png",
	"air_bursting": "textures/ui/bubble_pop.png",
	"air_empty":    "textures/ui/bubble_empty.png",
}

// rect is a region of a 256x256 legacy GUI atlas.
type rect struct{ x, y, w, h int }

// legacyIconRects locates each 1.20.2+ HUD sprite inside the gui/icons.png used by packs made for
// 1.8-1.20.1. The positions are exact matches between the 1.20.1 atlas and the 1.20.2 sprites.
var legacyIconRects = map[string]rect{
	"crosshair":                             {0, 0, 15, 15},
	"heart/container":                       {16, 0, 9, 9},
	"heart/container_blinking":              {25, 0, 9, 9},
	"heart/full":                            {52, 0, 9, 9},
	"heart/half":                            {61, 0, 9, 9},
	"heart/full_blinking":                   {70, 0, 9, 9},
	"heart/half_blinking":                   {79, 0, 9, 9},
	"heart/poisoned_full":                   {88, 0, 9, 9},
	"heart/poisoned_half":                   {97, 0, 9, 9},
	"heart/poisoned_full_blinking":          {106, 0, 9, 9},
	"heart/poisoned_half_blinking":          {115, 0, 9, 9},
	"heart/withered_full":                   {124, 0, 9, 9},
	"heart/withered_half":                   {133, 0, 9, 9},
	"heart/withered_full_blinking":          {142, 0, 9, 9},
	"heart/withered_half_blinking":          {151, 0, 9, 9},
	"heart/absorbing_full":                  {160, 0, 9, 9},
	"heart/absorbing_half":                  {169, 0, 9, 9},
	"heart/frozen_full":                     {178, 0, 9, 9},
	"heart/frozen_half":                     {187, 0, 9, 9},
	"armor_empty":                           {16, 9, 9, 9},
	"armor_half":                            {25, 9, 9, 9},
	"armor_full":                            {34, 9, 9, 9},
	"heart/vehicle_container":               {52, 9, 9, 9},
	"heart/vehicle_full":                    {88, 9, 9, 9},
	"heart/vehicle_half":                    {97, 9, 9, 9},
	"air":                                   {16, 18, 9, 9},
	"air_bursting":                          {25, 18, 9, 9},
	"food_empty":                            {16, 27, 9, 9},
	"food_full":                             {52, 27, 9, 9},
	"food_half":                             {61, 27, 9, 9},
	"food_full_hunger":                      {88, 27, 9, 9},
	"food_half_hunger":                      {97, 27, 9, 9},
	"food_empty_hunger":                     {133, 27, 9, 9},
	"heart/container_hardcore":              {16, 45, 9, 9},
	"heart/container_hardcore_blinking":     {25, 45, 9, 9},
	"heart/hardcore_full":                   {52, 45, 9, 9},
	"heart/hardcore_half":                   {61, 45, 9, 9},
	"heart/hardcore_full_blinking":          {70, 45, 9, 9},
	"heart/hardcore_half_blinking":          {79, 45, 9, 9},
	"heart/poisoned_hardcore_full":          {88, 45, 9, 9},
	"heart/poisoned_hardcore_half":          {97, 45, 9, 9},
	"heart/poisoned_hardcore_full_blinking": {106, 45, 9, 9},
	"heart/poisoned_hardcore_half_blinking": {115, 45, 9, 9},
	"heart/withered_hardcore_full":          {124, 45, 9, 9},
	"heart/withered_hardcore_half":          {133, 45, 9, 9},
	"heart/withered_hardcore_full_blinking": {142, 45, 9, 9},
	"heart/withered_hardcore_half_blinking": {151, 45, 9, 9},
	"heart/absorbing_hardcore_full":         {160, 45, 9, 9},
	"heart/absorbing_hardcore_half":         {169, 45, 9, 9},
	"heart/frozen_hardcore_full":            {178, 45, 9, 9},
	"heart/frozen_hardcore_half":            {187, 45, 9, 9},
	"experience_bar_background":             {0, 64, 182, 5},
	"experience_bar_progress":               {0, 69, 182, 5},
	"jump_bar_cooldown":                     {0, 74, 182, 5},
	"jump_bar_background":                   {0, 84, 182, 5},
	"jump_bar_progress":                     {0, 89, 182, 5},
	"hotbar_attack_indicator_background":    {0, 94, 18, 18},
	"hotbar_attack_indicator_progress":      {18, 94, 18, 18},
	"crosshair_attack_indicator_background": {36, 94, 16, 4},
	"crosshair_attack_indicator_progress":   {52, 94, 16, 4},
	"crosshair_attack_indicator_full":       {68, 94, 16, 16},
}

// Legacy gui/widgets.png regions.
var (
	legacyHotbarRect    = rect{0, 0, 182, 22}
	legacySelectionRect = rect{0, 22, 24, 24}
)

// atlasScale returns how many pixels of an HD atlas cover one pixel of its 256-pixel vanilla layout.
func atlasScale(t image.Texture) int {
	if s := t.Bounds().Dx() / 256; s > 1 {
		return s
	}
	return 1
}

// crop cuts r, scaled by s, out of an atlas texture.
func (r rect) crop(t image.Texture, s int, name string) image.Texture {
	return t.Crop(img.Rect(r.x*s, r.y*s, (r.x+r.w)*s, (r.y+r.h)*s), name)
}

// legacySprites cuts the HUD sprites out of a legacy gui/icons.png, skipping regions the pack left
// empty (for example frozen hearts, which packs made before 1.17 do not have).
func legacySprites(icons image.Texture) map[string]image.Texture {
	sprites := map[string]image.Texture{}
	if icons.Image == nil {
		return sprites
	}
	s := atlasScale(icons)
	for name, r := range legacyIconRects {
		if t := r.crop(icons, s, name+".png"); !t.Empty() {
			sprites[name] = t
		}
	}
	return sprites
}

// hudTextures converts HUD sprites, a hotbar strip and a hotbar selection into Bedrock UI textures.
func hudTextures(sprites map[string]image.Texture, hotbar, selection image.Texture) map[string]image.Texture {
	out := map[string]image.Texture{}
	for name, p := range spriteUI {
		if t, ok := sprites[name]; ok {
			out[p] = t
		}
	}
	if hotbar.Image != nil {
		// The 182-pixel Java hotbar is a 1-pixel cap, nine 20-pixel slots and another 1-pixel cap;
		// Bedrock stores each piece separately.
		s := hotbar.Bounds().Dx() / 182
		if s < 1 {
			s = 1
		}
		h := hotbar.Bounds().Dy()
		out["textures/ui/hotbar_start_cap.png"] = hotbar.Crop(img.Rect(0, 0, s, h), "hotbar_start_cap.png")
		for i := 0; i < 9; i++ {
			x := (1 + 20*i) * s
			name := "hotbar_" + strconv.Itoa(i) + ".png"
			out["textures/ui/"+name] = hotbar.Crop(img.Rect(x, 0, x+20*s, h), name)
		}
		out["textures/ui/hotbar_end_cap.png"] = hotbar.Crop(img.Rect(181*s, 0, 182*s, h), "hotbar_end_cap.png")
	}
	if selection.Image != nil {
		out["textures/ui/selected_hotbar_slot.png"] = squareSelection(selection)
	}
	return out
}

// squareSelection pads the 24x23 hotbar selection of 1.20.2+ to Bedrock's 24x24 by repeating its top
// border row at the bottom, which is how the 24x24 legacy and Bedrock textures are drawn.
func squareSelection(t image.Texture) image.Texture {
	b := t.Bounds()
	w, h := b.Dx(), b.Dy()
	if w == h || w == 0 {
		return t.Rename("selected_hotbar_slot.png")
	}
	s := w / 24
	if s < 1 {
		s = 1
	}
	c := image.NewCanvas(w, w)
	c.Paste(t, img.Point{})
	for y := h; y < w; y++ {
		c.Paste(t.Crop(img.Rect(0, (y-h)%s, w, (y-h)%s+1), ""), img.Pt(0, y))
	}
	return c.Texture("selected_hotbar_slot.png")
}

// iconsFromSprites rebuilds a gui/icons.png in the legacy layout from 1.20.2+ sprites so that the
// crosshair, which Bedrock still reads from that atlas, is ported. Regions without a sprite stay empty.
func iconsFromSprites(sprites map[string]image.Texture) (image.Texture, bool) {
	cross, ok := sprites["crosshair"]
	if !ok {
		return image.Texture{}, false
	}
	s := cross.Bounds().Dx() / 15
	if s < 1 {
		s = 1
	}
	c := image.NewCanvas(256*s, 256*s)
	for name, r := range legacyIconRects {
		t, ok := sprites[name]
		if !ok || t.Bounds().Dx() != r.w*s || t.Bounds().Dy() != r.h*s {
			continue
		}
		c.Paste(t, img.Pt(r.x*s, r.y*s))
	}
	return c.Texture("icons.png"), true
}
