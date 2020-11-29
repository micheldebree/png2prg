package main

import (
	"fmt"
	"image"
	"image/gif"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func newSourceImage(filename string) (*sourceImage, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not os.Open file %q: %v", filename, err)
	}
	defer r.Close()

	img := &sourceImage{sourceFilename: filename}

	if err = img.setPreferredBitpairColors(bitPairColors); err != nil {
		return nil, fmt.Errorf("setPreferredBitpairColors failed: %v", err)
	}

	if img.image, _, err = image.Decode(r); err != nil {
		return nil, fmt.Errorf("image.Decode %q failed: %v", filename, err)
	}

	if err = img.checkBounds(); err != nil {
		return nil, fmt.Errorf("img.checkBounds error %q: %v", filename, err)
	}
	if verbose && (img.xOffset != 0 || img.yOffset != 0) {
		log.Printf("img.xOffset, yOffset = %d, %d\n", img.xOffset, img.yOffset)
	}

	return img, nil
}

func newSourceImages(filename string) (imgs []sourceImage, err error) {
	if strings.ToLower(filepath.Ext(filename)) != ".gif" {
		return nil, fmt.Errorf("%q is not a .gif", filename)
	}

	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("os.Open could not open file %q: %v", filename, err)
	}
	defer r.Close()

	g, err := gif.DecodeAll(r)
	if err != nil {
		return nil, fmt.Errorf("gif.DecodeAll %q failed: %v", filename, err)
	}

	if len(g.Image) < 2 {
		return nil, fmt.Errorf("file %q is not an animation, frames found: %d", filename, len(g.Image))
	}
	if verbose {
		log.Printf("found animated gif %q with %d frames", filename, len(g.Image))
	}

	for _, img := range g.Image {
		width, height := img.Bounds().Max.X-img.Bounds().Min.X, img.Bounds().Max.Y-img.Bounds().Min.Y
		s := sourceImage{
			sourceFilename: filename,
			image:          img,
		}
		if err = s.setPreferredBitpairColors(bitPairColors); err != nil {
			return nil, fmt.Errorf("setPreferredBitpairColors failed: %v", err)
		}
		imgs = append(imgs, s)
		fmt.Printf("img %T, width %d x height %d\n", img, width, height)
	}
	return imgs, nil
}

func (img *sourceImage) setPreferredBitpairColors(v string) (err error) {
	if v != "" {
		if img.preferredBitpairColors, err = parseBitPairColors(v); err != nil {
			return fmt.Errorf("parseBitPairColors %q failed: %v", bitPairColors, err)
		}
		if verbose {
			log.Printf("will prefer bitpair colors: %v", img.preferredBitpairColors)
		}
	}
	return nil
}

func (img *sourceImage) checkBounds() error {
	width, height := img.image.Bounds().Max.X-img.image.Bounds().Min.X, img.image.Bounds().Max.Y-img.image.Bounds().Min.Y
	img.xOffset, img.yOffset = img.image.Bounds().Min.X, img.image.Bounds().Min.Y

	switch {
	case (width == 320) && (height == 200):
		return nil
	case (width == 384) && (height == 272):
		// default screenshot size in vice with default borders
		img.xOffset += (384 - 320) / 2
		img.yOffset += ((272 - 200) / 2) - 1
		return nil
	}
	return fmt.Errorf("image %q is not 320x200 or 384x272 pixels, but %d x %d pixels", img.sourceFilename, width, height)
}

func (img *sourceImage) analyze() error {
	img.analyzePalette()
	err := img.makeCharColors()
	if err != nil {
		return err
	}

	max, _ := img.maxColorsPerChar()
	if verbose {
		log.Printf("max colors per char: %d\n", max)
	}
	numColors, colorIndexes := img.countColors()
	if verbose {
		log.Printf("total colors: %d (%v)\n", numColors, colorIndexes)
	}
	switch {
	case max < 2:
		return fmt.Errorf("max colors per char %q < 2, is this a blank image?", max)
	case numColors == 2:
		img.graphicsType = singleColorCharset
		if verbose {
			log.Println("singleColorCharset found")
		}
	case max == 2:
		img.graphicsType = singleColorBitmap
		if verbose {
			log.Println("singleColorBitmap found")
		}
	case numColors == 3 || numColors == 4:
		img.graphicsType = multiColorCharset
		if verbose {
			log.Println("multiColorCharset found")
		}
	case max > 2:
		img.graphicsType = multiColorBitmap
		img.findBackgroundColor()
		if verbose {
			log.Println("multiColorBitmap found")
		}
	}

	return nil
}

func (img *sourceImage) countColors() (int, []byte) {
	m := make(map[RGB]byte, 16)
	for i := range img.charColors {
		for rgb, colorIndex := range img.charColors[i] {
			m[rgb] = colorIndex
		}
	}
	ci := []byte{}
	for _, v := range m {
		ci = append(ci, v)
	}
	sort.Slice(ci, func(i, j int) bool {
		return ci[i] < ci[j]
	})
	return len(m), ci
}

func (img *sourceImage) maxColorsPerChar() (max int, m map[RGB]byte) {
	for i := range img.charColors {
		if len(img.charColors[i]) > max {
			max = len(img.charColors[i])
			m = img.charColors[i]
		}
	}
	return max, m
}

func (img *sourceImage) findBackgroundColorCandidates() {
	backgroundCharColors := []map[RGB]byte{}
	for _, v := range img.charColors {
		if len(v) == 4 {
			backgroundCharColors = append(backgroundCharColors, v)
		}
	}

	// need to copy the map, as we delete stuff to eliminate false candidates
	candidates := make(map[RGB]byte, 16)
	switch {
	case len(backgroundCharColors) > 0:
		for k, v := range backgroundCharColors[0] {
			candidates[k] = v
		}
	default:
		for k, v := range img.palette {
			candidates[k] = v
		}
	}

	if verbose {
		log.Printf("all bgcol candidates: %v", candidates)
	}

	for _, charcolormap := range backgroundCharColors {
		for rgb, _ := range candidates {
			if _, ok := charcolormap[rgb]; !ok {
				if verbose {
					log.Printf("not a bgcol candidate, delete: %v", rgb)
				}
				delete(candidates, rgb)
			}
		}
	}
	img.backgroundCandidates = candidates
	if verbose {
		log.Printf("final bgcol candidates = %v", img.backgroundCandidates)
	}
	return
}

func (img *sourceImage) findBackgroundColor() {
	if img.backgroundCandidates == nil {
		img.findBackgroundColorCandidates()
	}
	var rgb RGB
	var colorIndex byte
	forceBgCol := -1
	if len(img.preferredBitpairColors) > 0 {
		forceBgCol = int(img.preferredBitpairColors[0])
	}

	for rgb, colorIndex = range img.backgroundCandidates {
		switch {
		case forceBgCol < 0:
			img.backgroundColor = colorInfo{rgb: rgb, colorIndex: colorIndex}
			return
		default:
			if colorIndex == byte(forceBgCol) {
				if verbose {
					log.Printf("findBackgroundColor: successfully found forced background color %d\n", forceBgCol)
				}
				img.backgroundColor = colorInfo{rgb: rgb, colorIndex: colorIndex}
				return
			}
		}
	}
	log.Printf("findBackgroundColor: error, dont think we found anything?")
	img.backgroundColor = colorInfo{rgb: rgb, colorIndex: colorIndex}
	return
}

func (img *sourceImage) makeCharColors() error {
	forceBgCol := -1
	if len(img.preferredBitpairColors) > 0 {
		forceBgCol = int(img.preferredBitpairColors[0])
	}
	fatalError := false
	for char := 0; char < 1000; char++ {
		charColors := img.colorMapFromChar(char)
		if forceBgCol >= 0 && len(charColors) == 4 {
			found := false
			for _, val := range charColors {
				if val == byte(forceBgCol) {
					found = true
					break
				}
			}
			if !found {
				x, y := xyFromChar(char)
				log.Printf("forced bgcol %d not possible in char %v (x=%d, y=%d)", forceBgCol, char, x, y)
				fatalError = true
			}
		}
		if len(charColors) > 4 {
			count := make(map[byte]byte, 16)
			for _, indexcolor := range charColors {
				count[indexcolor] = 1
			}
			if len(count) > 4 {
				x, y := xyFromChar(char)
				log.Printf("amount of colors in char %v (x=%d, y=%d) %d > 4 : %v", char, x, y, len(count), count)
				fatalError = true
			}
		}

		img.charColors[char] = charColors
	}
	if fatalError {
		return fmt.Errorf("fatal error: unable to convert %q", img.sourceFilename)
	}
	return nil
}

func (img *sourceImage) colorMapFromChar(char int) map[RGB]byte {
	charColors := make(map[RGB]byte, 16)
	x, y := xyFromChar(char)
	x += img.xOffset
	y += img.yOffset
	for j := 0; j < 32; j++ {
		pixelx, pixely := xyFromPixel(j)
		r, g, b, _ := img.image.At(x+pixelx, y+pixely).RGBA()
		rgb := RGB{byte(r), byte(g), byte(b)}
		if _, ok := charColors[rgb]; !ok {
			charColors[rgb] = img.palette[rgb]
		}
	}
	return charColors
}

func xyFromChar(i int) (int, int) {
	return 8*i - (320 * int(math.Floor(float64(i/40)))),
		8 * int(math.Floor(float64(i/40)))
}

func xyFromPixel(i int) (x, y int) {
	return i << 1 & 7, i >> 2
}

// analyzePalette finds the closest paletteMap and sets img.palette
func (img *sourceImage) analyzePalette() {
	minDistance := 9e9
	paletteName := ""
	paletteMap := make(map[RGB]byte, 16)
	for name, palette := range c64palettes {
		distance, curMap := img.distanceAndMap(palette)

		if verbose {
			log.Printf("color distance: %v => %v\n", name, distance)
		}
		if distance < minDistance {
			paletteMap, paletteName, minDistance = curMap, name, distance
		}
		if distance == 0 {
			break
		}
	}
	if verbose {
		log.Printf("%v palette found: %v distance: %v\n", img.sourceFilename, paletteName, minDistance)
	}
	img.palette = paletteMap
	return
}

func (img *sourceImage) distanceAndMap(palette [16]C64RGB) (float64, map[RGB]byte) {
	curMap := make(map[RGB]byte, 16)
	totalDistance := 0.0
	for x := 0; x < 320; x += 2 {
		for y := 0; y < 200; y++ {
			r, g, b, _ := img.image.At(img.xOffset+x, img.yOffset+y).RGBA()
			rgb := RGB{byte(r), byte(g), byte(b)}
			if _, ok := curMap[rgb]; !ok {
				d := 0.0
				curMap[rgb], d = rgb.colorIndexAndDistance(palette)
				totalDistance += d
			}
		}
	}
	return totalDistance, curMap
}

func (r RGB) colorIndexAndDistance(palette [16]C64RGB) (byte, float64) {
	distance := r.distanceTo(palette[0].RGB)
	closestColorIndex := 0
	for i := 1; i < len(palette); i++ {
		d := r.distanceTo(palette[i].RGB)
		if d < distance {
			distance = d
			closestColorIndex = i
		}
	}
	return byte(closestColorIndex), distance
}

func (r RGB) distanceTo(r2 RGB) float64 {
	dr := math.Abs(float64(r.R) - float64(r2.R))
	dg := math.Abs(float64(r.G) - float64(r2.G))
	db := math.Abs(float64(r.B) - float64(r2.B))
	return dr + dg + db
}
