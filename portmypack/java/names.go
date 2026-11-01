package java

import "strings"

// BedrockBlockPath returns where a block texture read from a pack with the given layout belongs in a
// Bedrock pack: a path relative to textures/ without extension, such as "blocks/planks_oak".
func BedrockBlockPath(layout Layout, fileName string) string {
	return bedrockPath(layout, "block", "blocks", fileName)
}

// BedrockItemPath returns where an item texture read from a pack with the given layout belongs in a
// Bedrock pack: a path relative to textures/ without extension, such as "items/apple_golden".
func BedrockItemPath(layout Layout, fileName string) string {
	return bedrockPath(layout, "item", "items", fileName)
}

// bedrockPath looks fileName up in the table for layout, defaulting to the same name in bedrockDir.
func bedrockPath(layout Layout, flatDir, bedrockDir, fileName string) string {
	name := strings.TrimSuffix(fileName, ".png")
	if layout == LayoutFlattened {
		if p, ok := flattenedNames[flatDir+"/"+name]; ok {
			return p
		}
		return bedrockDir + "/" + name
	}
	if p, ok := legacyNames[bedrockDir+"/"+name]; ok {
		return p
	}
	return bedrockDir + "/" + name
}
