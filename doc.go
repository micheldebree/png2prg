package png2prg

import (
	"flag"
	"fmt"
)

func PrintUsage() {
	fmt.Println("usage: ./png2prg [-help -h -d -q -v -bpc 0,6,14,3 -o outfile.prg -td testdata] FILE [FILE..]")
}

func PrintHelp() {
	fmt.Println()
	fmt.Printf("# PNG2PRG %v by Burglar\n", Version)
	fmt.Println()
	fmt.Println("Png2prg converts a 320x200 image (png/gif/jpeg) to a c64 hires or")
	fmt.Println("multicolor bitmap, charset or sprites. It will find the best matching palette")
	fmt.Println("and backgroundcolor automatically, no need to modify your source images or")
	fmt.Println("configure a palette.")
	fmt.Println("Vice screenshots with default borders (384x272) are automatically cropped.")
	fmt.Println("Vice's main screen's offset is at x=32, y=35.")
	fmt.Println("Images in sprite dimensions will be converted to sprites.")
	fmt.Println()
	fmt.Println("The resulting .prg includes the 2-byte start address and optional displayer.")
	fmt.Println("The displayers for koala, hires, mcibitmap and animations include fullscreen")
	fmt.Println("fade-in/out and optionally a .sid tune.")
	fmt.Println()
	fmt.Println("This tool can be used in all buildchains on most platforms.")
	fmt.Println()
	fmt.Println("## What it is *not*")
	fmt.Println()
	fmt.Println("Png2prg is not a tool to wire fullcolor images. It needs input images to")
	fmt.Println("already be compliant with c64 color and size restrictions.")
	fmt.Println("In verbose mode (-v) it outputs locations of color clashes, if any.")
	fmt.Println()
	fmt.Println("If you do need to wire fullcolor images, check out Youth's [Retropixels](https://www.micheldebree.nl/retropixels/).")
	fmt.Println()
	fmt.Println("## Supported Graphics Modes")
	fmt.Println()
	fmt.Println("    koala:     multicolor bitmap (max 4 colors per char)")
	fmt.Println("    hires:     singlecolor bitmap (max 2 colors per char)")
	fmt.Println("    mccharset: multicolor charset (max 4 colors)")
	fmt.Println("    sccharset: singlecolor charset (max 2 colors)")
	fmt.Println("    mcsprites: multicolor sprites (max 4 colors)")
	fmt.Println("    scsprites: singlecolor sprites (max 2 colors)")
	fmt.Println("    mcibitmap: 320x200 multicolor interlace bitmap (max 4 colors per char)")
	fmt.Println()
	fmt.Println("Png2prg is mostly able to autodetect the correct graphics mode, but you can")
	fmt.Println("also force a specific graphics mode with the -mode flag:")
	fmt.Println()
	fmt.Println("    ./png2prg -m koala image.png")
	fmt.Println()
	fmt.Println("## Koala or Hires Bitmap")
	fmt.Println()
	fmt.Println("    Bitmap: $2000 - $3f3f")
	fmt.Println("    Screen: $3f40 - $4327")
	fmt.Println("    D020:   $4328         (singlecolor only)")
	fmt.Println("    D800:   $4328 - $470f (multicolor only)")
	fmt.Println("    D021:   $4710         (multicolor only, low-nibble)")
	fmt.Println("    D020:   $4710         (multicolor only, high-nibble)")
	fmt.Println()
	fmt.Println("## Multicolor Interlace Bitmap")
	fmt.Println()
	fmt.Println("Experimental interlace support, you can supply one 320x200 multicolor")
	fmt.Println("images with max 4 colors per 8x8 pixel char:")
	fmt.Println()
	fmt.Println("    ./png2prg testdata/madonna/cjam_pure_madonna.png")
	fmt.Println()
	fmt.Println("Or supply both frames in regular koala specs (-interlace required):")
	fmt.Println()
	fmt.Println("    ./png2prg -i testdata/madonna/frame_0.png testdata/madonna/frame_1.png")
	fmt.Println()
	fmt.Println("### Drazlace (shared screenram and colorram for both frames")
	fmt.Println()
	fmt.Println("    D800:    $5800 - $5be7")
	fmt.Println("    Screen:  $5c00 - $5fe7")
	fmt.Println("    Bitmap1: $6000 - $7f3f")
	fmt.Println("    D021:    $7f40         (low-nibble)")
	fmt.Println("    D020:    $7f40         (high-nibble)")
	fmt.Println("    D016Offset: $7f42")
	fmt.Println("    Bitmap2: $8000 - $9f3f")
	fmt.Println()
	fmt.Println("### Multicolor Interlace (shared colorram, true paint .mci format)")
	fmt.Println()
	fmt.Println("    Screen1: $9c00 - $9fe7")
	fmt.Println("    D021:    $9fe8         (low-nibble)")
	fmt.Println("    D020:    $9fe8         (high-nibble)")
	fmt.Println("    D016Offset: $9fe9")
	fmt.Println("    Bitmap1: $a000 - $bf3f")
	fmt.Println("    Bitmap2: $c000 - $df3f")
	fmt.Println("    Screen2: $e000 - $e3e7")
	fmt.Println("    D800:    $e400 - $e7e7")
	fmt.Println()
	fmt.Println("## Single or Multicolor Charset")
	fmt.Println()
	fmt.Println("Currently only images with max 4 colors can be converted into a charset.")
	fmt.Println("Support for individual d800 colors and mixed single/multicolor chars may be")
	fmt.Println("added in a future release, if the need arises.")
	fmt.Println()
	fmt.Println("By default charsets are packed, they only contain unique characters.")
	fmt.Println("If you do not want charpacking, eg for a 1x1 charset, please use -no-pack")
	fmt.Println()
	fmt.Println("    Charset:   $2000-$27ff")
	fmt.Println("    Screen:    $2800-$2be7")
	fmt.Println("    CharColor: $2be8")
	fmt.Println("    D021:      $2be9")
	fmt.Println("    D022:      $2bea       (multicolor only)")
	fmt.Println("    D023:      $2beb       (multicolor only)")
	fmt.Println("    D020:      $2bec       (multicolor only)")
	fmt.Println()
	fmt.Println("## Single or Multicolor Sprites")
	fmt.Println()
	fmt.Println("If the source image size is a multiple of a 24x21 pixel sprite,")
	fmt.Println("the image is considered to contain sprites.")
	fmt.Println()
	fmt.Println("The image will be converted left to right, top to bottom.")
	fmt.Println()
	fmt.Println("    Sprite 1: $2000-$203f")
	fmt.Println("    Sprite 2: $2040-$207f")
	fmt.Println("    ...")
	fmt.Println()
	fmt.Println("## Bitpair Colors")
	fmt.Println()
	fmt.Println("By default, png2prg guesses bitpair colors by itself. In most cases you")
	fmt.Println("don't need to configure anything. It will provide a mostly normalized image")
	fmt.Println("which should yield good pack results, but your miles may vary.")
	fmt.Println()
	fmt.Println("To give you more control, you can force/prefer a specific bitpair")
	fmt.Println("color-order. Use c64 colors, so 0 for black, 1 for white, 2 for red, etc.")
	fmt.Println()
	fmt.Println("The following example will force background color 0 for bitpair 00 and")
	fmt.Println("prefer colors 6,14,3 for bitpairs 01,10,11:")
	fmt.Println()
	fmt.Println("    ./png2prg -bitpair-colors 0,6,14,3 image.png")
	fmt.Println()
	fmt.Println("It's also possible to explicitly skip certain bitpair preferences with -1:")
	fmt.Println()
	fmt.Println("    ./png2prg -bitpair-colors 0,-1,-1,3 image.png")
	fmt.Println()
	fmt.Println("## Sprite Animation")
	fmt.Println()
	fmt.Println("Each frame will be concatenated in the output .prg.")
	fmt.Println("You can supply an animated .gif or multiple image files.")
	fmt.Println()
	fmt.Println("## Bitmap Animation (only koala and hires)")
	fmt.Println()
	fmt.Println("If multiple files are added, they are treated as animation frames.")
	fmt.Println("You can also supply an animated .gif.")
	fmt.Println("The first image will be exported with all framedata appended.")
	fmt.Println("Koala animation frames start at $4711, hires at $4329.")
	fmt.Println()
	fmt.Println("The frame files are following this format.")
	fmt.Println("Each frame consists of 1 or more chunks. A chunk looks like this:")
	fmt.Println()
	fmt.Println("    .byte $03    // number of chars in this chunk")
	fmt.Println("                 // $00 marks end of frame")
	fmt.Println("                 // $ff marks end of all frames")
	fmt.Println("    .word bitmap // bitmap address of this chunk (the high byte is <$20)")
	fmt.Println("    .word screen // screenram address (the high byte is <$04)")
	fmt.Println()
	fmt.Println("    For each char in this chunk:")
	fmt.Println()
	fmt.Println("      .byte 0,31,15,7,8,34,0,128 // pixels")
	fmt.Println("      .byte $64                  // screenram colors")
	fmt.Println("      .byte $01                  // colorram color (koala only)")
	fmt.Println("      ...                        // next char(s)")
	fmt.Println()
	fmt.Println("    ...          // next chunks")
	fmt.Println("    .byte 0      // end of frame")
	fmt.Println("    ...          // next frame(s)")
	fmt.Println("    .byte $ff    // end of all frames")
	fmt.Println()
	fmt.Println("## Displayer")
	fmt.Println()
	fmt.Println("The -d or -display flag will link displayer code infront of the picture.")
	fmt.Println("By default it will also crunch the resulting file with Antonio Savona's")
	fmt.Println("[TSCrunch](https://github.com/tonysavon/TSCrunch/) with a couple of changes in my own [fork](https://github.com/staD020/TSCrunch/).")
	fmt.Println()
	fmt.Println("For hires, koala and koala-anim the displayer also supports adding a .sid.")
	fmt.Println("Multispeed sids are supported as long as the .sid initializes the CIA timers")
	fmt.Println("correctly.")
	fmt.Println()
	fmt.Println("You can use sids located from $0d00-$1fff or $9000+ in hires/koala displayers.")
	fmt.Println("For animation displayers use $0e00-$1fff, $4900-$88ff and >$e000-$fff9.")
	fmt.Println("Note that animation frames will be loaded to $4800 and up and could overload")
	fmt.Println("the sid.")
	fmt.Println()
	fmt.Println("NB: For hires anims, $4500-$abff is free and anim is loaded to $4400.")
	fmt.Println("If needed, you can relocate most sids using lft's [sidreloc](http://www.linusakesson.net/software/sidreloc/index.php).")
	fmt.Println("In general $0e00-$1fff and $6000-$88ff are pretty safe.")
	fmt.Println()
	fmt.Println("Zeropages $08-$0f are used in the animation displayers, while none are used")
	fmt.Println("in hires/koala displayers, increasing sid compatibility.")
	fmt.Println()
	fmt.Println("## Examples")
	fmt.Println()
	fmt.Println("This release contains examples with all assets included for you to test with.")
	fmt.Println("Also included are the assets of [Évoluer](https://csdb.dk/release/?id=220170) by The Sarge and Flotsam.")
	fmt.Println()
	fmt.Println("## Changes for version 1.3")
	fmt.Println()
	fmt.Println(" - Support for even more far-out palette ranges (thanks Perplex).")
	fmt.Println(" - Now throws an error if the palette can't be detected properly, this should")
	fmt.Println("   never happen. Please let me know if you run into this error.")
	fmt.Println(" - Separated library and cli tool.")
	fmt.Println(" - Library supports the standard [io.Reader](https://pkg.go.dev/io@go1.19.1#Reader) and [io.Writer](https://pkg.go.dev/io@go1.19.1#Writer) interfaces.")
	fmt.Println(" - Slightly improved packratio for koala images, in most cases.")
	fmt.Println(" - Patched [TSCrunch](https://github.com/staD020/TSCrunch/) further to increase crunch speed and use less memory.")
	fmt.Println(" - Added -parallel and -worker flags to treat each input file as standalone")
	fmt.Println("   and convert all files in parallel. Gifs with multiple frames are still")
	fmt.Println("   treated as animations.")
	fmt.Println(" - Stop relying on .gif filename extension, detect it.")
	fmt.Println(" - When conversion fails for images in vice 384x272 resolution (offset 32, 35),")
	fmt.Println("   png2prg also tries common wrong offset 32, 36, used by multiple graphicians.")
	fmt.Println(" - Add -symbols flag to write symbols to a .sym file.")
	fmt.Println(" - Experimental interlace support, mci bitmap (drazlace and truepaint).")
	fmt.Println()
	fmt.Println("## Changes for version 1.2")
	fmt.Println()
	fmt.Println(" - Added displayer for koala animations.")
	fmt.Println(" - Added displayer for hires animations.")
	fmt.Println(" - Added -frame-delay flag for animation displayers.")
	fmt.Println(" - Added -wait-seconds flag for animation displayers.")
	fmt.Println(" - Fixed bug in koala/hires displayers not allowing sids to overlap $c000-$c7ff.")
	fmt.Println(" - Expanding wildcards: using pic??.png or pic*.png now also works on Windows.")
	fmt.Println(" - Set bank via $dd00 in displayers.")
	fmt.Println()
	fmt.Println("## Changes for version 1.0")
	fmt.Println()
	fmt.Println(" - Added fullscreen fade in/out to koala and hires displayers.")
	fmt.Println(" - Added optional .sid support for koala and hires displayers.")
	fmt.Println(" - Added optional crunching for all displayers using TSCrunch.")
	fmt.Println()
	fmt.Println("## Credits")
	fmt.Println()
	fmt.Println("Png2prg was created by Burglar, using the following third-party libraries:")
	fmt.Println()
	fmt.Println("[TSCrunch 1.3](https://github.com/tonysavon/TSCrunch/) by Antonio Savona for optional crunching when exporting")
	fmt.Println("an image with a displayer.")
	fmt.Println()
	fmt.Println("[Colfade Doc](https://csdb.dk/release/?id=132276) by Veto for the colfade")
	fmt.Println("tables used in the koala and hires displayers.")
	fmt.Println()
	fmt.Println("[Kick Assembler](http://www.theweb.dk/KickAssembler/) by Slammer to compile the displayers.")
	fmt.Println()
	fmt.Println("[Go](https://go.dev/) by The Go Authors is the programming language used to create png2prg.")
	fmt.Println()
	fmt.Println("## Options")
	fmt.Println()
	fmt.Println("```")
	flag.PrintDefaults()
	fmt.Println("```")
	fmt.Println()
	return
}
