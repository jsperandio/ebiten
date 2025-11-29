package core

import (
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
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
	face := textv2.NewGoXFace(FontMap[fontSize])
	
	// text/v2 uses different coordinate system - need to adjust Y position
	// The old text API used baseline as Y coordinate, text/v2 uses top-left of rendering region
	// Convert from baseline to top by subtracting HAscent
	metrics := face.Metrics()
	adjustedY := float64(y) - metrics.HAscent
	
	opts := &textv2.DrawOptions{}
	opts.GeoM.Translate(float64(x), adjustedY)
	opts.ColorScale.ScaleWithColor(clr)
	textv2.Draw(target, txt, face, opts)
}

func DrawTextWithOptions(target *ebiten.Image, txt string, x, y int, opts DrawTextOptions) {
	f := FontMap[opts.FontSize]
	c := opts.Color
	r, _ := font.BoundString(f, txt)
	x2 := 0.0
	y2 := 0.0
	if opts.HAligh == HAlignLeft {
		x2 = float64(int(-r.Min.X) + x)
		y2 = float64(int(-r.Min.Y) + y)
	} else if opts.HAligh == HAlignCenter {
		x2 = float64(x - int((r.Max.X-r.Min.X)/2))
		y2 = float64(y + int((r.Max.Y-r.Min.Y)/2))
	}
	face := textv2.NewGoXFace(f)
	// Convert from baseline (old API) to top-left (text/v2) by subtracting HAscent
	metrics := face.Metrics()
	adjustedY2 := y2 - metrics.HAscent
	drawOpts := &textv2.DrawOptions{}
	drawOpts.GeoM.Translate(x2, adjustedY2)
	drawOpts.ColorScale.ScaleWithColor(c)
	textv2.Draw(target, txt, face, drawOpts)
}
