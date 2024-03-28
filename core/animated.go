package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type CascadeText struct {
	CoreMoveableImage
	Text  string
	Speed int
	Color color.Color

	yInital       int
	fontSize      FontSize
	durationTicks int
	x             float64
	y             float64
}

func NewCascadeText(text string, x, y float64, fontSize FontSize, clr color.Color, speed int) *CascadeText {
	return &CascadeText{
		Text:          text,
		yInital:       initialPosYByFontSize(fontSize),
		fontSize:      fontSize,
		x:             x,
		y:             y,
		Speed:         speed,
		durationTicks: 50,
		Color:         clr,
	}
}

func (ct *CascadeText) Draw(screen *ebiten.Image, gameTick int) {
	if ct.durationTicks <= 0 {
		return
	}

	bkg := makeBackgroundImageFromText(ct.fontSize, ct.Text)
	ct.yInital += ct.Speed

	DrawText(bkg, ct.Text, 0, ct.yInital, ct.Color, ct.fontSize)

	ct.CoreMoveableImage = *NewCoreMoveableImageFromImage(bkg, ct.x, ct.y)
	ct.SimpleDraw(screen)
	ct.durationTicks--
}

type AscendingText struct {
	CoreMoveableImage
	Text  string
	Speed int
	Color color.Color

	yInital       int
	fontSize      FontSize
	durationTicks int
	x             float64
	y             float64
}

func NewAscendingText(txt string, x, y float64, fs FontSize, clr color.Color, spd int) *AscendingText {
	return &AscendingText{
		Text:          txt,
		yInital:       100 - initialPosYByFontSize(fs),
		fontSize:      fs,
		x:             x,
		y:             y,
		Speed:         spd,
		durationTicks: 50,
		Color:         clr,
	}
}

func (at *AscendingText) Draw(screen *ebiten.Image, tck int) {
	if at.durationTicks < 0 {
		return
	}

	bkg := makeBackgroundImageFromText(at.fontSize, at.Text)
	at.yInital -= at.Speed

	DrawText(bkg, at.Text, 0, at.yInital, at.Color, at.fontSize)

	at.CoreMoveableImage = *NewCoreMoveableImageFromImage(bkg, at.x, at.y)
	at.SimpleDraw(screen)
	at.durationTicks--
}

func makeBackgroundImageFromText(fs FontSize, txt string) *ebiten.Image {
	txtBound, _ := font.BoundString(FontMap[fs], txt)
	return ebiten.NewImage(txtBound.Max.X.Ceil(), 100)
}

func initialPosYByFontSize(fontSize FontSize) int {
	switch fontSize {
	case FontSizeXLarge:
		return 33
	case FontSizeLarge:
		return 28
	case FontSizeMedium:
		return 15
	case FontSizeSmall:
		return 8
	default:
		return 0
	}
}
