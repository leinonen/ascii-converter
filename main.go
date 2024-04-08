package main

import (
	"ascii-converter/cp437"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ascii to image converter")
		fmt.Println("Usage", os.Args[0], "filename")
		return
	}

	artFile := os.Args[1]
	imgFile := replaceExtension(artFile, "png")
	lines, err := readAscii(artFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Lines:", len(lines))

	// Set up the image dimensions
	fnt := loadFont()
	advance, _ := fnt.GlyphAdvance(' ')
	characterWidth := advance.Ceil()
	width := 79 * characterWidth
	height := len(lines) * fnt.Metrics().Height.Ceil()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	drawLines(lines, fnt, img)

	saveImage(img, imgFile)
}

func replaceExtension(filename, newExt string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return filename + "." + newExt
	}
	return filename[:len(filename)-len(ext)] + "." + newExt
}

func drawLines(lines []string, fnt font.Face, img *image.RGBA) {
	y := 0
	for _, line := range lines {
		y += fnt.Metrics().Height.Ceil()
		drawString(img, line, 0, y, fnt, color.White)
	}
}

func readAscii(fileName string) ([]string, error) {
	contentBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("error reading ascii file" + fileName)
	}

	lines := strings.Split(string(contentBytes), "\r\n")

	var result []string
	for _, line := range lines {
		result = append(result, cp437.String([]byte(line)))
	}
	return result, nil
}

func saveImage(img *image.RGBA, outFile string) {
	imgFile, err := os.Create(outFile)
	if err != nil {
		fmt.Println("Error creating image file:", err)
		return
	}
	defer imgFile.Close()

	err = png.Encode(imgFile, img)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	fmt.Println("Image saved as", outFile)
}

func drawString(img draw.Image, s string, x, y int, fnt font.Face, clr color.Color) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(clr),
		Face: fnt,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	d.DrawString(s)
}

func loadFont() font.Face {
	fontBytes := gomono.TTF
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	// Create a Face for the loaded font
	return truetype.NewFace(font, &truetype.Options{
		// Size: float64(config.FontSize),
	})
}
