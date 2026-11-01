// Potion tints derived from the Mojang/bedrock-samples v1.26.60.5 pre-coloured potion textures. For most
// effects, bottle + potion_overlay * tint reproduces the Bedrock texture exactly.

package portmypack

import "image/color"

// potionTint is the overlay colour Bedrock bakes into each pre-coloured potion texture, and the bottle
// kinds Bedrock ships for it.
type potionTint struct {
	color                    color.NRGBA
	drink, splash, lingering bool
}

// potionTints is keyed by the Bedrock effect suffix used in items/potion_bottle_*<effect>.png.
var potionTints = map[string]potionTint{
	"absorption":     {color.NRGBA{60, 133, 255, 255}, true, true, false},
	"blindness":      {color.NRGBA{71, 76, 71, 255}, true, true, false},
	"confusion":      {color.NRGBA{138, 48, 121, 255}, true, true, false},
	"damageBoost":    {color.NRGBA{252, 197, 0, 255}, true, true, true},
	"digSlowdown":    {color.NRGBA{121, 107, 38, 255}, true, true, false},
	"digSpeed":       {color.NRGBA{255, 255, 109, 255}, true, true, false},
	"fireResistance": {color.NRGBA{252, 151, 0, 255}, true, true, true},
	"harm":           {color.NRGBA{167, 100, 105, 255}, true, true, true},
	"heal":           {color.NRGBA{245, 35, 34, 255}, true, true, true},
	"healthBoost":    {color.NRGBA{255, 203, 57, 255}, true, true, false},
	"hunger":         {color.NRGBA{143, 191, 135, 255}, true, true, false},
	"infested":       {color.NRGBA{139, 153, 139, 255}, true, true, true},
	"invisibility":   {color.NRGBA{243, 243, 243, 255}, true, true, true},
	"jump":           {color.NRGBA{250, 252, 131, 255}, true, true, true},
	"levitation":     {color.NRGBA{250, 252, 131, 255}, true, true, false},
	"luck":           {color.NRGBA{88, 191, 6, 255}, false, false, true},
	"moveSlowdown":   {color.NRGBA{137, 173, 222, 255}, true, true, true},
	"moveSpeed":      {color.NRGBA{50, 232, 252, 255}, true, true, true},
	"nightVision":    {color.NRGBA{192, 252, 101, 255}, true, true, true},
	"oozing":         {color.NRGBA{151, 252, 161, 255}, true, true, true},
	"poison":         {color.NRGBA{134, 161, 98, 255}, true, true, true},
	"regeneration":   {color.NRGBA{203, 91, 169, 255}, true, true, true},
	"resistance":     {color.NRGBA{252, 151, 0, 255}, true, true, false},
	"saturation":     {color.NRGBA{245, 35, 34, 255}, true, true, false},
	"slowFall":       {color.NRGBA{243, 207, 185, 255}, true, true, true},
	"turtleMaster":   {color.NRGBA{139, 129, 227, 255}, true, true, true},
	"waterBreathing": {color.NRGBA{150, 215, 190, 255}, true, true, true},
	"weakness":       {color.NRGBA{71, 76, 71, 255}, true, true, true},
	"weaving":        {color.NRGBA{119, 104, 89, 255}, true, true, true},
	"windCharged":    {color.NRGBA{186, 198, 251, 255}, true, true, true},
	"wither":         {color.NRGBA{115, 97, 86, 255}, true, true, true},
}
