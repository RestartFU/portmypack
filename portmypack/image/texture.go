package image

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
)

// Texture represents a Minecraft texture.
type Texture struct {
	image.Image
	Name string
}

// NewTexture creates a new Texture from a name and reader.
func NewTexture(name string, reader io.Reader, alphaFix bool) (Texture, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return Texture{}, err
	}

	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return Texture{}, err
	}

	if alphaFix {
		img = fixAlpha(img)
	}

	t := Texture{
		Image: img,
		Name:  path.Base(name),
	}
	return t, nil
}

// NewTextureFS reads the PNG at filePath in f.
func NewTextureFS(f fs.FS, filePath string, alphaFix bool) (Texture, error) {
	reader, err := f.Open(filePath)
	if err != nil {
		return Texture{}, err
	}
	defer reader.Close()

	return NewTexture(filePath, reader, alphaFix)
}

// fixAlpha makes every partially transparent pixel fully transparent, because Bedrock renders items and
// HUD icons with alpha testing.
func fixAlpha(img image.Image) image.Image {
	newImage := image.NewNRGBA(img.Bounds())

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			pixel := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			if pixel.A != 0 && pixel.A != 255 {
				pixel = color.NRGBA{}
			}
			newImage.SetNRGBA(x, y, pixel)
		}
	}
	return newImage
}

// Empty reports whether the texture has no image or no visible pixel.
func (t Texture) Empty() bool {
	if t.Image == nil {
		return true
	}
	b := t.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if _, _, _, a := t.At(x, y).RGBA(); a != 0 {
				return false
			}
		}
	}
	return true
}

// Crop returns a copy of the r region of the texture, relative to its top-left corner, under a new name.
// Parts of r outside the texture stay transparent.
func (t Texture) Crop(r image.Rectangle, name string) Texture {
	out := image.NewNRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	if t.Image != nil {
		min := t.Bounds().Min
		draw.Draw(out, out.Bounds(), t.Image, r.Min.Add(min), draw.Src)
	}
	return Texture{Image: out, Name: name}
}

// FirstFrame returns the top square frame of a vertical animation strip, or the texture itself when it
// is not a strip.
func (t Texture) FirstFrame() Texture {
	if t.Image == nil {
		return t
	}
	w, h := t.Bounds().Dx(), t.Bounds().Dy()
	if w == 0 || h < 2*w || h%w != 0 {
		return t
	}
	return t.Crop(image.Rect(0, 0, w, w), t.Name)
}

// Rename returns the texture under a new name.
func (t Texture) Rename(name string) Texture {
	t.Name = name
	return t
}

// Tint multiplies the colour of every pixel by c, the way Minecraft tints item layers.
func (t Texture) Tint(c color.NRGBA) Texture {
	b := t.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			p := color.NRGBAModel.Convert(t.At(x, y)).(color.NRGBA)
			p.R = uint8(uint16(p.R) * uint16(c.R) / 255)
			p.G = uint8(uint16(p.G) * uint16(c.G) / 255)
			p.B = uint8(uint16(p.B) * uint16(c.B) / 255)
			out.SetNRGBA(x-b.Min.X, y-b.Min.Y, p)
		}
	}
	return Texture{Image: out, Name: t.Name}
}

// Over draws top over the texture, scaling nothing: both are aligned at their top-left corners.
func (t Texture) Over(top Texture, name string) Texture {
	b := t.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), t.Image, b.Min, draw.Src)
	draw.Draw(out, out.Bounds(), top.Image, top.Bounds().Min, draw.Over)
	return Texture{Image: out, Name: name}
}

// Canvas is a transparent texture that other textures are pasted into.
type Canvas struct {
	*image.NRGBA
}

// NewCanvas creates a transparent canvas of the given size.
func NewCanvas(w, h int) Canvas {
	return Canvas{NRGBA: image.NewNRGBA(image.Rect(0, 0, w, h))}
}

// Paste copies src onto the canvas with its top-left corner at p.
func (c Canvas) Paste(src Texture, p image.Point) {
	if src.Image == nil {
		return
	}
	r := image.Rectangle{Min: p, Max: p.Add(src.Bounds().Size())}
	draw.Draw(c.NRGBA, r, src.Image, src.Bounds().Min, draw.Src)
}

// Texture returns the canvas contents as a named texture.
func (c Canvas) Texture(name string) Texture {
	return Texture{Image: c.NRGBA, Name: name}
}

// WriteFile writes the texture to a PNG file.
func (t Texture) WriteFile(to string) error {
	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, t)
}

// Write writes the texture as a PNG at path in a zip archive.
func (t Texture) Write(w *zip.Writer, path string) error {
	writer, err := w.Create(path)
	if err != nil {
		return err
	}

	return png.Encode(writer, t)
}

// DirectoryTexturesFS reads every PNG below dirPath in f. Unreadable PNGs are skipped and reported in the
// returned error alongside the textures that could be read.
func DirectoryTexturesFS(f fs.FS, dirPath string, alphaFix bool) (textures []Texture, err error) {
	var decodeErrs []error
	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".png") {
			return nil
		}

		texture, err := NewTextureFS(f, path, alphaFix)
		if err != nil {
			decodeErrs = append(decodeErrs, fmt.Errorf("%s: %w", path, err))
			return nil
		}

		textures = append(textures, texture)
		return nil
	}

	err = fs.WalkDir(f, dirPath, walkDirFunc)
	return textures, errors.Join(append([]error{err}, decodeErrs...)...)
}
