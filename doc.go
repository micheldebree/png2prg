package png2prg

import (
	"flag"
	"fmt"
)

func PrintUsage() {
	fmt.Println("usage: ./png2prg [-help -h -d -q -v -bpc 0,6,14,3 -o outfile.prg -td testdata] FILE [FILE..]")
}

func PrintHelp() {
	fmt.Printf("# PNG2PRG %v by burg\n", Version)
	fmt.Println()
	fmt.Println("Png2prg converts a 320x200 image (png/gif/jpeg) to a c64 hires or")
	fmt.Println("multicolor bitmap, charset, petscii, ecm or sprites prg. It will find the best")
	fmt.Println("matching palette and background/bitpair-colors automatically, no need to modify")
	fmt.Println("your source images or configure a palette.")
	fmt.Println()
	fmt.Println("Vice screenshots with default borders (384x272) are automatically cropped.")
	fmt.Println("Vice's main screen offset is at x=32, y=35.")
	fmt.Println("Images in sprite dimensions will be converted to sprites.")
	fmt.Println()
	fmt.Println("The resulting .prg includes the 2-byte start address and optional displayer.")
	fmt.Println("The displayers can optionally play a .sid tune.")
	fmt.Println()
	fmt.Println("This tool can be used in all buildchains on all common platforms.")
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
	fmt.Println("    koala:        multicolor bitmap (max 4 colors per char)")
	fmt.Println("    hires:        singlecolor bitmap (max 2 colors per char)")
	fmt.Println("    mixedcharset: multicolor charset (max 4 colors per char (fixed bgcol, d022, d023))")
	fmt.Println("    mccharset:    multicolor charset (max 4 colors)")
	fmt.Println("    sccharset:    singlecolor charset (max 2 colors per char (fixed bgcol))")
	fmt.Println("    petscii:      singlecolor rom charset (max 2 colors per char (fixed bgcol))")
	fmt.Println("    ecm:          singlecolor charset (max 2 colors per char (4 fixed bgcolors), max 64 chars)")
	fmt.Println("    mcsprites:    multicolor sprites (max 4 colors)")
	fmt.Println("    scsprites:    singlecolor sprites (max 2 colors)")
	fmt.Println("    mcibitmap:    320x200 multicolor interlace bitmap (max 4 colors per char/frame)")
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
	fmt.Println("You can supply one 320x200 multicolor image with max 4 colors per 8x8 pixel")
	fmt.Println("char per frame of which at least 2 are shared (the D021 and D800 colors).")
	fmt.Println()
	fmt.Println("Or supply both frames in regular koala specs (-interlace flag required).")
	fmt.Println("When making screenshots in vice, please disable the d016 pixel shift manually.")
	fmt.Println()
	fmt.Println("    ./png2prg -i testdata/madonna/frame_0.png testdata/madonna/frame_1.png")
	fmt.Println()
	fmt.Println("### Drazlace (shared screenram and colorram for both frames)")
	fmt.Println()
	fmt.Println("    ./png2prg testdata/madonna/cjam_pure_madonna.png")
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
	fmt.Println("    ./png2prg -i -d016 1 testdata/mcinterlace/parriot?.png")
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
	fmt.Println("## Singlecolor, PETSCII or ECM Charset (individual d800 colors)")
	fmt.Println()
	fmt.Println("By default charsets are packed, they only contain unique characters.")
	fmt.Println("If you do not want charpacking, eg for a 1x1 charset, please use -no-pack.")
	fmt.Println()
	fmt.Println("With ECM -bitpair-colors can be used to force d021-d024 colors.")
	fmt.Println()
	fmt.Println("NB: individual d800 colors are not supported with -no-pack.")
	fmt.Println()
	fmt.Println("    ./png2prg -m sccharset testdata/hirescharset/ohno_logo.png")
	fmt.Println("    ./png2prg -m petscii testdata/petscii/hein_hibiscus.png")
	fmt.Println("    ./png2prg -m ecm testdata/ecm/xpardey.png")
	fmt.Println("    ./png2prg -m ecm testdata/ecm/shampoo.png")
	fmt.Println("    ./png2prg -m ecm -bpc 2,7,14,0 testdata/ecm/orion.png")
	fmt.Println()
	fmt.Println("    Charset:   $2000-$27ff (omitted for petscii)")
	fmt.Println("    Screen:    $2800-$2be7")
	fmt.Println("    D800:      $2c00-$2fe7")
	fmt.Println("    D020:      $2fe8")
	fmt.Println("    D021:      $2fe9")
	fmt.Println("    D022:      $2fea (ecm only)")
	fmt.Println("    D023:      $2feb (ecm only)")
	fmt.Println("    D024:      $2fec (ecm only)")
	fmt.Println()
	fmt.Println("## Mixed Multi/Singlecolor Charset (individual d800 colors)")
	fmt.Println()
	fmt.Println("Png2prg tries to figure out the right -bitpair-colors and auto-corrects")
	fmt.Println("where it can, but there still are edge-cases like the ones below.")
	fmt.Println("If an impossible color is found, an error will be displayed.")
	fmt.Println("Swap some -bpc colors around and retry.")
	fmt.Println("There can also be cases where manual -bpc colors can influence char-count or")
	fmt.Println("packed size.")
	fmt.Println()
	fmt.Println("    ./png2prg -m mixedcharset testdata/mixedcharset/hein_neo.png")
	fmt.Println("    ./png2prg -m mixedcharset testdata/mixedcharset/huntress.gif")
	fmt.Println("    ./png2prg -m mixedcharset -bpc 3 testdata/mixedcharset/shine.png")
	fmt.Println("    ./png2prg -m mixedcharset -bpc 0 testdata/mixedcharset/charsetcompo.png")
	fmt.Println()
	fmt.Println("    Charset:   $2000-$27ff")
	fmt.Println("    Screen:    $2800-$2be7")
	fmt.Println("    D800:      $2c00-$2fe7")
	fmt.Println("    D020:      $2fe8")
	fmt.Println("    D021:      $2fe9")
	fmt.Println("    D022:      $2fea")
	fmt.Println("    D023:      $2feb")
	fmt.Println()
	fmt.Println("## Single or Multicolor Sprites")
	fmt.Println()
	fmt.Println("If the source image size is a multiple of a 24x21 pixel sprite,")
	fmt.Println("the image is considered to contain sprites.")
	fmt.Println()
	fmt.Println("The image will be converted from left to right, top to bottom.")
	fmt.Println()
	fmt.Println("    ./png2prg image.png")
	fmt.Println("    ./png2prg -m scsprites image.png")
	fmt.Println("    ./png2prg -m mcsprites image.png")
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
	fmt.Println("All displayers except for sprites support adding a .sid.")
	fmt.Println("Multispeed sids are supported as long as the .sid initializes the CIA timers")
	fmt.Println("correctly.")
	fmt.Println()
	fmt.Println("You can use sids located from $0e00-$1fff or $e000+ in the displayers.")
	fmt.Println("More areas may be free depending on graphics type.")
	fmt.Println("A memory usage map is shown on error and in -vv (very verbose) mode.")
	fmt.Println()
	fmt.Println("If needed, you can relocate most sids using lft's [sidreloc](http://www.linusakesson.net/software/sidreloc/index.php).")
	fmt.Println()
	fmt.Println("Zeropages $08-$0f are used in the animation displayers, while none are used")
	fmt.Println("in hires/koala displayers, increasing sid compatibility.")
	fmt.Println()
	fmt.Println("## Brute Force Mode and Pack Optimization")
	fmt.Println()
	fmt.Println("By default png2prg 1.8 does a pretty good job at optimizing the resulting prg")
	fmt.Println("for crunchers and packers. It is not enough to beat [SPOT 1.3](https://csdb.dk/release/?id=242492).")
	fmt.Println()
	fmt.Println("The optimization techniques used by png2prg are also responsible for cleaning")
	fmt.Println("up the bitmap, making it ideal for animations and color effects.")
	fmt.Println()
	fmt.Println("### -brute-force (-bf)")
	fmt.Println()
	fmt.Println("Iterates are over many -bitpair-colors permutations automatically, packs")
	fmt.Println("with the built in TSCrunch and selects the shortest.")
	fmt.Println()
	fmt.Println("    ./png2prg -bf image.png")
	fmt.Println()
	fmt.Println("The -brute-force mode can be used in combination with additional flags.")
	fmt.Println()
	fmt.Println("### -no-bitpair-counters (-nbc)")
	fmt.Println()
	fmt.Println("Disable counting of bitpairs per color to guess a bitpair for a color.")
	fmt.Println()
	fmt.Println("    ./png2prg -bf -nbc image.png")
	fmt.Println()
	fmt.Println("### -no-prev-char-colors (-npcc)")
	fmt.Println()
	fmt.Println("Disable lookback to previous char's charcolors to guess a bitpair for a color.")
	fmt.Println()
	fmt.Println("    ./png2prg -bf -npcc image.png")
	fmt.Println()
	fmt.Println("Since TSCrunch is optimized for speed, packing with Dali can give varying")
	fmt.Println("results. This is also the reason for not including these options in the")
	fmt.Println("brute force permutations automatically.")
	fmt.Println()
	fmt.Println("## Benchmark")
	fmt.Println()
	fmt.Println("The [koala otpimizing thread](https://csdb.dk/forums/?roomid=13&topicid=38311&showallposts=1) on csdb has gained some interest in the scene.")
	fmt.Println("Since Sparta released [SPOT](https://csdb.dk/release/?id=242492) it has been the best optimizer available.")
	fmt.Println()
	fmt.Println("Png2prg 1.8 has improved optimization techniques but requires -brute-force")
	fmt.Println("mode to beat SPOT 1.3. Manual flags can optimize even better in some cases.")
	fmt.Println()
	fmt.Println("All koalas are packed with [Dali 0.3.2](https://csdb.dk/release/?id=223584).")
	fmt.Println()
	fmt.Println("    +---------+--------+----------+------------+--------+")
	fmt.Println("    | spot1.3 | p2p1.8 | p2p1.8bf | p2p1.8best | p2p1.6 |")
	fmt.Println("    +---------+--------+----------+------------+--------+")
	fmt.Println("    |    7332 |   7372 |     7332 |       7324 |   7546 | Untitled/Floris")
	fmt.Println("    |    5136 |   5190 |     5149 |         bf |   5464 | Song of the Sunset/Mermaid")
	fmt.Println("    |    5968 |   5998 |     5963 |         bf |   6155 | Short Circuit/Karen Davies")
	fmt.Println("    |    3618 |   3647 |     3616 |       3589 |   3830 | Portait L+D/Sander")
	fmt.Println("    |    5094 |   5080 |     5083 |       5078 |   5320 | Weee/Mermaid")
	fmt.Println("    |    7497 |   7471 |     7458 |         bf |   7612 | Deadlock/Robin Levy")
	fmt.Println("    |    8068 |   8097 |     8046 |       8038 |   8227 | Room with a view/Veto")
	fmt.Println("    |    7445 |   7490 |     7432 |         bf |   7582 | Vangelis/Talent")
	fmt.Println("    |    6759 |   6739 |     6737 |         bf |   6963 | Temple of Tears/Hend")
	fmt.Println("    |    7859 |   7848 |     7839 |       7821 |   7998 | Thanos/JonEgg")
	fmt.Println("    |    4859 |   4849 |     4782 |         bf |   4983 | Solar-Sonar/Leon")
	fmt.Println("    |    5640 |   5671 |     5613 |         bf |   5869 | Cisco Heat/Alan Grier")
	fmt.Println("    |    6243 |   6286 |     6228 |         bf |   6430 | Daylight/Sulevi")
	fmt.Println("    |    2850 |   2884 |     2848 |         bf |   3092 | Yie Ar Kung Fu/Steve Wahid")
	fmt.Println("    |    6727 |   6721 |     6730 |       6711 |   6901 | Lee/The Sarge")
	fmt.Println("    |    7837 |   7828 |     7798 |         bf |   7960 | Parrot/Mirage")
	fmt.Println("    +---------+--------+----------+------------+--------+")
	fmt.Println("    |   98932 |  99171 |    98654 |      98569 | 101932 | Total")
	fmt.Println("    +---------+--------+----------+------------+--------+")
	fmt.Println()
	fmt.Println(" - p2p1.8: default png2prg result w/o options")
	fmt.Println(" - p2p1.8bf: -brute-force mode")
	fmt.Println(" - p2p1.8best: hand-picked -bitpair-colors, or bruteforced with -npcc and/or -nbc flags")
	fmt.Println(" - p2p1.6: default png2prg 1.6 result w/o options")
	fmt.Println()
	fmt.Println("## Examples")
	fmt.Println()
	fmt.Println("This release contains examples with all assets included for you to test with.")
	fmt.Println("Also included are the assets of [Évoluer](https://csdb.dk/release/?id=220170) by The Sarge and Flotsam.")
	fmt.Println("A larger set of testdata can be found in the [github repo](https://github.com/staD020/png2prg/tree/master/testdata).")
	fmt.Println()
	fmt.Println("## Install from source")
	fmt.Println()
	fmt.Println("Png2prg was built on Linux, building on Mac should work out of the box.")
	fmt.Println("For Windows, try out Windows Subsystem Linux (WSL), works pretty well.")
	fmt.Println("However, natively building on Windows should be easy enough, look at")
	fmt.Println("Compiling without Make below.")
	fmt.Println()
	fmt.Println("The compiled displayer prgs are included in the repo to ease building")
	fmt.Println("and importing png2prg as a library. Java is only required to build")
	fmt.Println("the displayers with KickAssembler (included in the repo).")
	fmt.Println()
	fmt.Println("But first [install Go 1.20 or higher](https://go.dev/dl/).")
	fmt.Println()
	fmt.Println("### Simple install")
	fmt.Println()
	fmt.Println("    go install -v github.com/staD020/png2prg.git@master")
	fmt.Println()
	fmt.Println("### Compiling with Make (recommended)")
	fmt.Println()
	fmt.Println("    git clone https://github.com/staD020/png2prg.git")
	fmt.Println("    cd png2prg")
	fmt.Println("    make -j")
	fmt.Println()
	fmt.Println("Build for all common targets:")
	fmt.Println()
	fmt.Println("    make all -j")
	fmt.Println()
	fmt.Println("### Compiling without Make")
	fmt.Println()
	fmt.Println("    go build ./cmd/png2prg")
	fmt.Println()
	fmt.Println("## Install and use as library")
	fmt.Println()
	fmt.Println("In your Go project's path, go get the library:")
	fmt.Println()
	fmt.Println("    go get github.com/staD020/png2prg")
	fmt.Println()
	fmt.Println("In essence png2prg implements the [io.WriterTo](https://pkg.go.dev/io#WriterTo) interface.")
	fmt.Println("Typical usage could look like below. A more complex example can be found")
	fmt.Println("in the [source](https://github.com/staD020/png2prg/blob/master/cmd/png2prg/main.go) of the cli tool.")
	fmt.Println()
	fmt.Println("```go")
	fmt.Println("import (")
	fmt.Println("	\"fmt\"")
	fmt.Println("	\"io\"")
	fmt.Println("	\"github.com/staD020/png2prg\"")
	fmt.Println(")")
	fmt.Println()
	fmt.Println("func convertPNG(w io.Writer, png io.Reader) (int64, error) {")
	fmt.Println("	p, err := png2prg.New(png2prg.Options{}, png)")
	fmt.Println("	if err != nil {")
	fmt.Println("		return 0, fmt.Errorf(\"png2prg.New failed: %w\", err)")
	fmt.Println("	}")
	fmt.Println("	return p.WriteTo(w)")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println()
	fmt.Printf("## Changes for version %s\n", Version)
	fmt.Println()
	fmt.Println(" - Improve crunchiness by re-using the previous char's bitpair-colors.")
	fmt.Println(" - Add -no-prev-char-colors flag to disable re-use of the previous char's")
	fmt.Println("   bitpair-colors, in some cases this optimization causes worse pack results.")
	fmt.Println(" - Add -brute-force mode to find bitpair color combinations with better")
	fmt.Println("   crunchiness. Burns some CPU for a couple seconds.")
	fmt.Println(" - Add -no-bitpair-counters flag to disable using bitpair counters per color")
	fmt.Println("   for color guessing.")
	fmt.Println(" - Added multi-frame support for mccharset, where all frames use the same")
	fmt.Println("   charset.")
	fmt.Println(" - Add support for any centered fullscreen image resolution bigger than")
	fmt.Println("   320x200 and other than 384x272.")
	fmt.Println(" - Add support for Marq's PETSCII tool .png resolution 352x232 (thanks jab).")
	fmt.Println(" - Bugfix: docs fixes related to installation from source (thanks jab).")
	fmt.Println(" - Bugfix: hide findECMColors log behind -verbose mode (thanks jab).")
	fmt.Println(" - Docs fix: add a bit more info for sprites (thanks fungus).")
	fmt.Println()
	fmt.Println("## Changes for version 1.6")
	fmt.Println()
	fmt.Println(" - Added -mode mixedcharset for mixed multicolor/singlecolor and")
	fmt.Println("   individual d800 colors per char.")
	fmt.Println(" - Modified -mode sccharset to use individual d800 colors per char.")
	fmt.Println(" - Added -mode petscii.")
	fmt.Println(" - Added -mode ecm.")
	fmt.Println(" - Added -no-pack-empty to skip packing empty chars to filled chars to re-use")
	fmt.Println("   for different colors. Only for mixed and ecm charsets.")
	fmt.Println(" - Added -force-pack-empty for singlecolor and multicolor charset, may save")
	fmt.Println("   a char, but usually pack-ratio is worse due to increased d800 color usage.")
	fmt.Println(" - Improved auto-detection of graphics modes, including various charset modes.")
	fmt.Println(" - Added sid support to charset displayers.")
	fmt.Println(" - Added fullscreen fade in/out to charset displayers.")
	fmt.Println(" - Bug Fix: -force-border-color for singlecolor charset (thanks Raistlin).")
	fmt.Println(" - Bug Fix: do not write empty .prg file on error.")
	fmt.Println(" - Standardized d02x colors in output.prg for charset modes.")
	fmt.Println()
	fmt.Println("## Changes for version 1.4")
	fmt.Println()
	fmt.Println(" - Support for even more far-out palette ranges (thanks Perplex).")
	fmt.Println(" - Now throws an error if the palette can't be detected properly, this should")
	fmt.Println("   never happen. Please let me know if you run into this error.")
	fmt.Println(" - Separated library and cli tool.")
	fmt.Println(" - Library supports the standard [io.Reader](https://pkg.go.dev/io#Reader) and [io.Writer](https://pkg.go.dev/io#Writer) interfaces.")
	fmt.Println(" - Patched [TSCrunch](https://github.com/staD020/TSCrunch/) further to increase crunch speed and use less memory.")
	fmt.Println(" - Added -parallel and -worker flags to treat each input file as standalone")
	fmt.Println("   and convert all files in parallel. Gifs with multiple frames are still")
	fmt.Println("   treated as animations.")
	fmt.Println(" - Stop relying on .gif filename extension, detect it.")
	fmt.Println(" - Add -alt-offset flag to force screenshot offset 32, 36), used by a few")
	fmt.Println("   graphicians. Though, please switch to the correct 32, 35.")
	fmt.Println(" - Add -symbols flag to write symbols to a .sym file.")
	fmt.Println(" - Interlace support for mcibitmap (drazlace and truepaint).")
	fmt.Println(" - Bugfix: allow blank images input (thanks Spider-J).")
	fmt.Println(" - Allow colors not present in the image as -bitpair-colors (thanks Map).")
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
}
