package core

import (
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	localfonts "github.com/jsperandio/ebiten/assets/fonts"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/font"
)

type FontSize int

const (
	FontSizeXLarge FontSize = iota
	FontSizeLarge
	FontSizeMedium
	FontSizeSmall

	FontSizeXLargeAbsolute = 45
	FontSizeLargeAbsolute  = 38
	FontSizeMediumAbsolute = 24
	FontSizeSmallAbsolute  = 12
)

type HAlign int

const (
	HAlignLeft HAlign = iota
	HAlignCenter
)

var FontMap = map[FontSize]font.Face{}

func LoadFonts() {
	tt, err := truetype.Parse(localfonts.ELECSTROM_REGULAR)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72

	FontMap[FontSizeXLarge] = truetype.NewFace(tt, &truetype.Options{
		Size:    45,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	FontMap[FontSizeLarge] = truetype.NewFace(tt, &truetype.Options{
		Size:    38,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	FontMap[FontSizeMedium] = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	FontMap[FontSizeSmall] = truetype.NewFace(tt, &truetype.Options{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

type DrawTextOptions struct {
	Color    color.Color
	Width    int
	HAligh   HAlign
	FontSize FontSize
}

func DrawText(target *ebiten.Image, txt string, x, y int, clr color.Color, fontSize FontSize) {
	text.Draw(target, txt, FontMap[fontSize], x, y, clr)
}

func DrawTextWithOptions(target *ebiten.Image, txt string, x, y int, opts DrawTextOptions) {
	f := FontMap[opts.FontSize]
	c := opts.Color
	r, _ := font.BoundString(f, txt)
	x2 := 0
	y2 := 0
	if opts.HAligh == HAlignLeft {
		x2 = int(-r.Min.X) + x
		y2 = int(-r.Min.Y) + y
	} else if opts.HAligh == HAlignCenter {
		x2 = x - int((r.Max.X-r.Min.X)/2)
		y2 = y + int((r.Max.Y-r.Min.Y)/2)
	}
	text.Draw(target, txt, f, x2, y2, c)
}
