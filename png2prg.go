package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tscrunch "github.com/staD020/TSCrunch"
)

const version = "0.9-dev"

type RGB struct {
	R, G, B byte
}

type colorInfo struct {
	ColorIndex byte
	RGB        RGB
}

type graphicsType byte

const (
	unknownGraphicsType graphicsType = iota
	singleColorBitmap
	multiColorBitmap
	singleColorCharset
	multiColorCharset
	singleColorSprites
	multiColorSprites
)

func stringToGraphicsType(s string) graphicsType {
	switch s {
	case "koala":
		return multiColorBitmap
	case "hires":
		return singleColorBitmap
	case "sccharset":
		return singleColorCharset
	case "mccharset":
		return multiColorCharset
	case "scsprites":
		return singleColorSprites
	case "mcsprites":
		return multiColorSprites
	}
	return unknownGraphicsType
}

func (t graphicsType) String() string {
	switch t {
	case singleColorBitmap:
		return "hires"
	case multiColorBitmap:
		return "koala"
	case singleColorCharset:
		return "singlecolor charset"
	case multiColorCharset:
		return "multicolor charset"
	case singleColorSprites:
		return "singlecolor sprites"
	case multiColorSprites:
		return "multicolor sprites"
	default:
		return "unknown"
	}
}

type bitpairColors []byte

func (b bitpairColors) String() (s string) {
	for i, v := range b {
		s = s + strconv.Itoa(int(v))
		if i < len(b)-1 {
			s = s + ","
		}
	}
	return s
}

type sourceImage struct {
	sourceFilename         string
	image                  image.Image
	xOffset                int
	yOffset                int
	width                  int
	height                 int
	palette                map[RGB]byte
	colors                 []RGB
	charColors             [1000]map[RGB]byte
	backgroundCandidates   map[RGB]byte
	backgroundColor        colorInfo
	borderColor            colorInfo
	preferredBitpairColors bitpairColors
	graphicsType           graphicsType
}

type MultiColorChar struct {
	CharIndex       int
	Bitmap          [8]byte
	BackgroundColor byte
	ScreenColor     byte
	D800Color       byte
}

type Koala struct {
	SourceFilename  string
	Bitmap          [8000]byte
	ScreenColor     [1000]byte
	D800Color       [1000]byte
	BackgroundColor byte
	BorderColor     byte
}

type Hires struct {
	SourceFilename string
	Bitmap         [8000]byte
	ScreenColor    [1000]byte
	BorderColor    byte
}

type MultiColorCharset struct {
	SourceFilename  string
	Bitmap          [0x800]byte
	Screen          [1000]byte
	CharColor       byte
	BackgroundColor byte
	D022Color       byte
	D023Color       byte
	BorderColor     byte
}

type SingleColorCharset struct {
	SourceFilename  string
	Bitmap          [0x800]byte
	Screen          [1000]byte
	CharColor       byte
	BackgroundColor byte
	BorderColor     byte
}

type SingleColorSprites struct {
	SourceFilename  string
	Bitmap          []byte
	SpriteColor     byte
	BackgroundColor byte
	Columns         byte
	Rows            byte
}

type MultiColorSprites struct {
	SourceFilename  string
	Bitmap          []byte
	SpriteColor     byte
	BackgroundColor byte
	D025Color       byte
	D026Color       byte
	Columns         byte
	Rows            byte
}

var displayers = make(map[graphicsType][]byte, 0)

//go:embed "display_koala.prg"
var koalaDisplay []byte

//go:embed "display_hires.prg"
var hiresDisplay []byte

//go:embed "display_mc_charset.prg"
var mcCharsetDisplay []byte

//go:embed "display_sc_charset.prg"
var scCharsetDisplay []byte

//go:embed "display_mc_sprites.prg"
var mcSpritesDisplay []byte

//go:embed "display_sc_sprites.prg"
var scSpritesDisplay []byte

func initDisplayers() {
	displayers[multiColorBitmap] = koalaDisplay
	displayers[singleColorBitmap] = hiresDisplay
	displayers[multiColorCharset] = mcCharsetDisplay
	displayers[singleColorCharset] = scCharsetDisplay
	displayers[multiColorSprites] = mcSpritesDisplay
	displayers[singleColorSprites] = scSpritesDisplay
}

func processFiles(filenames ...string) (err error) {
	if len(filenames) < 1 {
		log.Println("no files supplied, nothing to do.")
		return err
	}

	imgs, err := newSourceImages(filenames...)
	switch {
	case err != nil:
		return fmt.Errorf("newSourceImages failed: %v", err)
	case len(imgs) == 0:
		return fmt.Errorf("no images found")
	case len(imgs) > 1:
		if err = handleAnimation(imgs); err != nil {
			return fmt.Errorf("handleAnimation failed: %v", err)
		}
		return nil
	}

	img := imgs[0]
	if verbose {
		log.Printf("processing file %q", img.sourceFilename)
	}
	if err = img.analyze(); err != nil {
		return fmt.Errorf("analyze %q failed: %v", img.sourceFilename, err)
	}

	var c io.WriterTo
	switch img.graphicsType {
	case multiColorBitmap:
		if c, err = img.convertToKoala(); err != nil {
			return fmt.Errorf("convertToKoala %q failed: %v", img.sourceFilename, err)
		}
	case singleColorBitmap:
		if c, err = img.convertToHires(); err != nil {
			return fmt.Errorf("convertToHires %q failed: %v", img.sourceFilename, err)
		}
	case singleColorCharset:
		if c, err = img.convertToSingleColorCharset(); err != nil {
			if graphicsMode != "" {
				return fmt.Errorf("convertToSingleColorCharset %q failed: %v", img.sourceFilename, err)
			}
			if !quiet {
				fmt.Printf("falling back to %s because convertToSingleColorCharset %q failed: %v\n", singleColorBitmap, img.sourceFilename, err)
			}
			img.graphicsType = singleColorBitmap
			if c, err = img.convertToHires(); err != nil {
				return fmt.Errorf("convertToHires %q failed: %v", img.sourceFilename, err)
			}
		}
	case multiColorCharset:
		if c, err = img.convertToMultiColorCharset(); err != nil {
			if graphicsMode != "" {
				return fmt.Errorf("convertToMultiColorCharset %q failed: %v", img.sourceFilename, err)
			}
			if !quiet {
				fmt.Printf("falling back to %s because convertToMultiColorCharset %q failed: %v\n", multiColorBitmap, img.sourceFilename, err)
			}
			img.graphicsType = multiColorBitmap
			err = img.findBackgroundColor()
			if err != nil {
				return fmt.Errorf("findBackgroundColor %q failed: %v", img.sourceFilename, err)
			}
			if c, err = img.convertToKoala(); err != nil {
				return fmt.Errorf("convertToKoala %q failed: %v", img.sourceFilename, err)
			}
		}
	case singleColorSprites:
		if c, err = img.convertToSingleColorSprites(); err != nil {
			return fmt.Errorf("convertToSingleColorSprites %q failed: %v", img.sourceFilename, err)
		}
	case multiColorSprites:
		if c, err = img.convertToMultiColorSprites(); err != nil {
			return fmt.Errorf("convertToMultiColorSprites %q failed: %v", img.sourceFilename, err)
		}
	default:
		return fmt.Errorf("unsupported graphicsType for %q", img.sourceFilename)
	}

	if display {
		c, err = injectCrunch(c)
		if err != nil {
			return fmt.Errorf("injectCrunch failed: %v", err)
		}
	}

	destFilename := destinationFilename(img.sourceFilename)
	f, err := os.Create(destFilename)
	if err != nil {
		return fmt.Errorf("os.Create %q failed: %v", destFilename, err)
	}
	defer f.Close()

	if _, err = c.WriteTo(f); err != nil {
		return fmt.Errorf("WriteTo %q failed: %v", destFilename, err)
	}
	if !quiet {
		fmt.Printf("converted %q to %q in %q format\n", img.sourceFilename, destFilename, img.graphicsType)
	}

	return nil
}

func injectCrunch(in io.WriterTo) (io.WriterTo, error) {
	buf := &bytes.Buffer{}
	if _, err := in.WriteTo(buf); err != nil {
		return nil, fmt.Errorf("WriteTo buffer failed: %v", err)
	}
	conf := tscrunch.Config{
		PRG:     true,
		QUIET:   true,
		INPLACE: false,
		JumpTo:  "$0819",
	}
	out, err := tscrunch.New(conf, buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("tscrunch.New failed: %v", err)
	}
	return out, nil
}

// defaultHeader returns the startaddress of a file located at 0x2000
func defaultHeader() []byte {
	return []byte{0x00, 0x20}
}

func (k Koala) WriteTo(w io.Writer) (n int64, err error) {
	header := defaultHeader()
	if display {
		header = displayers[multiColorBitmap]
		fill := 0x2000 - 0x7ff - len(header)
		for i := 0; i < fill; i++ {
			header = append(header, 0)
		}
	}
	bgBorder := k.BackgroundColor | k.BorderColor<<4
	return writeData(w, [][]byte{header, k.Bitmap[:], k.ScreenColor[:], k.D800Color[:], {bgBorder}})
}

func (h Hires) WriteTo(w io.Writer) (n int64, err error) {
	header := defaultHeader()
	if display {
		header = displayers[singleColorBitmap]
		fill := 0x2000 - 0x7ff - len(header)
		for i := 0; i < fill; i++ {
			header = append(header, 0)
		}
	}
	return writeData(w, [][]byte{header, h.Bitmap[:], h.ScreenColor[:], {h.BorderColor}})
}

func (c MultiColorCharset) WriteTo(w io.Writer) (n int64, err error) {
	header := defaultHeader()
	if display {
		header = displayers[multiColorCharset]
	}
	return writeData(w, [][]byte{header, c.Bitmap[:], c.Screen[:], {c.CharColor, c.BackgroundColor, c.D022Color, c.D023Color, c.BorderColor}})
}

func (c SingleColorCharset) WriteTo(w io.Writer) (n int64, err error) {
	header := defaultHeader()
	if display {
		header = displayers[singleColorCharset]
	}
	return writeData(w, [][]byte{header, c.Bitmap[:], c.Screen[:], {c.CharColor, c.BackgroundColor}})
}

func (s SingleColorSprites) WriteTo(w io.Writer) (n int64, err error) {
	header := defaultHeader()
	if display {
		header = displayers[singleColorSprites]
		header = append(header, s.Columns, s.Rows, s.BackgroundColor, s.SpriteColor)
	}
	return writeData(w, [][]byte{header, s.Bitmap[:]})
}

func (s MultiColorSprites) WriteTo(w io.Writer) (n int64, err error) {
	header := defaultHeader()
	if display {
		header = displayers[multiColorSprites]
		header = append(header, s.Columns, s.Rows, s.BackgroundColor, s.D025Color, s.SpriteColor, s.D026Color)
	}
	return writeData(w, [][]byte{header, s.Bitmap[:]})
}

func writeData(w io.Writer, data [][]byte) (n int64, err error) {
	for _, d := range data {
		var m int
		m, err = w.Write(d)
		n += int64(m)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func destinationFilename(filename string) (destfilename string) {
	if len(targetdir) > 0 {
		destfilename = filepath.Dir(targetdir+string(os.PathSeparator)) + string(os.PathSeparator)
	}
	switch {
	case len(outfile) > 0:
		return destfilename + outfile
	case len(outfile) == 0:
		return destfilename + filepath.Base(strings.TrimSuffix(filename, filepath.Ext(filename))+".prg")
	}
	return destfilename
}
