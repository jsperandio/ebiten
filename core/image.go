package core

import (
	"bytes"
	"image"
	"image/color"

	log "github.com/sirupsen/logrus"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// #############################################################################
// ===============================CoreImage=====================================
// #############################################################################
const (
	RectCoreFill = iota
	RectCoreBorder
	RectCoreSpaceBorder
	RectCore
)

type Drawable interface {
	Draw(screen *ebiten.Image, tick int)
}

type CoreImage struct {
	RawImage []byte
	Img      *ebiten.Image
	ImgOpts  *ebiten.DrawImageOptions
}

func NewCoreImage(ri []byte, x, y float64) *CoreImage {
	sceneDecoded, _, err := image.Decode(bytes.NewReader(ri))
	if err != nil {
		log.Fatal(err)
	}

	ci := &CoreImage{
		RawImage: ri,
		Img:      ebiten.NewImageFromImage(sceneDecoded),
		ImgOpts:  &ebiten.DrawImageOptions{},
	}

	ci.SetImagePosition(x, y)

	return ci
}

func NewCoreImageFromImage(img *ebiten.Image, x, y float64) *CoreImage {
	ci := &CoreImage{
		RawImage: nil,
		Img:      img,
		ImgOpts:  &ebiten.DrawImageOptions{},
	}

	ci.SetImagePosition(x, y)

	return ci
}

func NewCoreTextRectImage(image *ebiten.Image, x, y float64, clr color.Color, text string, textclr color.Color, fontSize FontSize) *CoreImage {
	ci := &CoreImage{
		RawImage: nil,
		Img:      image,
		ImgOpts:  &ebiten.DrawImageOptions{},
	}
	ci.Img.Fill(clr)

	textRect, _ := font.BoundString(FontMap[fontSize], text)
	centerX := (ci.Width() / 2) - (textRect.Max.X.Ceil() / 2)
	centerY := (ci.Height() / 2) - (textRect.Min.Y.Floor() / 2)
	DrawText(ci.Img, text, centerX, centerY, textclr, fontSize)

	ci.SetImagePosition(x, y)
	return ci
}

func NewCoreRectImage(gf *ebiten.Image, x, y float64, fillType int, clr color.Color) *CoreImage {
	ci := &CoreImage{
		RawImage: nil,
		Img:      gf,
		ImgOpts:  &ebiten.DrawImageOptions{},
	}
	rect := ci.Img.Bounds()

	switch fillType {
	case RectCoreFill:
		ci.Img.Fill(clr)

	case RectCoreBorder:
		boundsColor := ClrDebugColor
		var lineWidth float32 = 2.0

		// top-left to top-right
		vector.StrokeLine(ci.Img, 0, 0, float32(rect.Dx()), 0, lineWidth, boundsColor, false)
		// top-right to bottom-right
		vector.StrokeLine(ci.Img, float32(rect.Dx()), 0, float32(rect.Dx()), float32(rect.Dy()), lineWidth, boundsColor, false)
		// bottom-right to bottom-left
		vector.StrokeLine(ci.Img, float32(rect.Dx()), float32(rect.Dy()), 0, float32(rect.Dy()), lineWidth, boundsColor, false)
		// bottom-left to top-left
		vector.StrokeLine(ci.Img, 0, float32(rect.Dy()), 0, 0, lineWidth, boundsColor, false)
	case RectCoreSpaceBorder:
		boundsColor := ClrYellow
		var lineWidth float32 = 2.0

		// top-left to top-right
		vector.StrokeLine(ci.Img, 0, 0, float32(rect.Dx()), 0, lineWidth, boundsColor, false)
		// top-right to bottom-right
		vector.StrokeLine(ci.Img, float32(rect.Dx()), 0, float32(rect.Dx()), float32(rect.Dy()), lineWidth, boundsColor, false)
		// bottom-right to bottom-left
		vector.StrokeLine(ci.Img, float32(rect.Dx()), float32(rect.Dy()), 0, float32(rect.Dy()), lineWidth, boundsColor, false)
		// bottom-left to top-left
		vector.StrokeLine(ci.Img, 0, float32(rect.Dy()), 0, 0, lineWidth, boundsColor, false)

	}

	ci.SetImagePosition(x, y)
	return ci
}

func (c *CoreImage) SetImagePosition(x, y float64) {
	c.ImgOpts.GeoM.Translate(x, y)
}

func (c *CoreImage) Width() int {
	return c.Img.Bounds().Dx()
}

func (c *CoreImage) Height() int {
	return c.Img.Bounds().Dy()
}

func (c *CoreImage) X0Y0() image.Point {
	return image.Point{
		X: c.Img.Bounds().Min.X,
		Y: c.Img.Bounds().Min.Y,
	}
}

func (c *CoreImage) X1Y0() image.Point {
	return image.Point{
		X: c.Img.Bounds().Max.X,
		Y: c.Img.Bounds().Min.Y,
	}
}

func (c *CoreImage) X1Y1() image.Point {
	return image.Point{
		X: c.Img.Bounds().Max.X,
		Y: c.Img.Bounds().Max.Y,
	}
}

func (c *CoreImage) X0Y1() image.Point {
	return image.Point{
		X: c.Img.Bounds().Min.X,
		Y: c.Img.Bounds().Max.Y,
	}
}

func (c *CoreImage) SimpleDraw(screen *ebiten.Image) {
	screen.DrawImage(c.Img, c.ImgOpts)
	c.debugBounds(screen, c.Img, c.ImgOpts)
}

func (c *CoreImage) debugBounds(screen *ebiten.Image, refImg *ebiten.Image, imgOptions *ebiten.DrawImageOptions) {
	if !DEBUG {
		return
	}

	width := float64(refImg.Bounds().Dx())
	height := float64(refImg.Bounds().Dy())
	rectImg := NewCoreRectImage(ebiten.NewImage(refImg.Bounds().Dx(), refImg.Bounds().Dy()), width, height, RectCoreBorder, nil)
	screen.DrawImage(rectImg.Img, imgOptions)
}

// #############################################################################
// ===========================CoreMoveableImage=================================
// #############################################################################

type CoreMoveableImage struct {
	CoreImage
	PosX float64
	PosY float64
}

func NewCoreMoveableImage(ri []byte, x, y float64) *CoreMoveableImage {
	cmi := &CoreMoveableImage{
		CoreImage: *NewCoreImage(ri, x, y),
		PosX:      x,
		PosY:      y,
	}

	return cmi
}

func NewCoreMoveableImageFromImage(img *ebiten.Image, x, y float64) *CoreMoveableImage {
	cmi := &CoreMoveableImage{
		CoreImage: *NewCoreImageFromImage(img, x, y),
		PosX:      x,
		PosY:      y,
	}

	return cmi
}

func (cmi *CoreMoveableImage) MoveRight(ms float64) {
	// cmi.PosX = (float64(Right) * ms) + cmi.PosX
	// cmi.PosY = 0 + cmi.PosY
	cmi.PosX += (float64(Right) * ms)
	cmi.SetImagePosition(cmi.PosX, cmi.PosY)
}

func (cmi *CoreMoveableImage) MoveLeft(ms float64) {
	cmi.PosX = (float64(Left) * ms) + cmi.PosX
	cmi.PosY = 0 + cmi.PosY
	cmi.SetImagePosition(cmi.PosX, cmi.PosY)
}

func (cmi *CoreMoveableImage) MoveUp(ms float64) {
	cmi.PosX = 0 + cmi.PosX
	cmi.PosY = (float64(Up) * ms) + cmi.PosY
	cmi.SetImagePosition(cmi.PosX, cmi.PosY)
}

func (cmi *CoreMoveableImage) MoveDown(ms float64) {
	cmi.PosX = 0 + cmi.PosX
	cmi.PosY = (float64(Down) * ms) + cmi.PosY
	cmi.SetImagePosition(cmi.PosX, cmi.PosY)
}

// #############################################################################
// ===========================CoreAnimatedImage=================================
// #############################################################################

type CoreAnimatedImage struct {
	CoreMoveableImage
	FrameOX     int
	FrameOY     int
	FrameWidth  int
	FrameHeight int
	FrameNum    int
}

func NewCoreAnimatedImage(ri []byte, x, y float64, frameOX, frameOY, frameWidth, frameHeight, frameNum int) *CoreAnimatedImage {
	cai := &CoreAnimatedImage{
		CoreMoveableImage: *NewCoreMoveableImage(ri, x, y),
		FrameOX:           frameOX,
		FrameOY:           frameOY,
		FrameWidth:        frameWidth,
		FrameHeight:       frameHeight,
		FrameNum:          frameNum,
	}

	return cai
}

func (cai *CoreAnimatedImage) AnimatedDraw(screen *ebiten.Image, frame int, imgDirection Direction) {
	op := &ebiten.DrawImageOptions{}
	// remove img size from point
	op.GeoM.Translate(-float64(cai.FrameWidth)/2, -float64(cai.FrameHeight)/2)

	// Flip image if needed
	if imgDirection == Left {
		op.GeoM.Scale(1, 1)
	}
	if imgDirection == Right {
		op.GeoM.Scale(-1, 1)
	}

	// set image position
	op.GeoM.Translate(cai.PosX, cai.PosY)

	// set image frame update rate
	i := (frame / 10) % cai.FrameNum
	// set image frame on sheet
	sx, sy := cai.FrameOX+i*cai.FrameWidth, cai.FrameOY
	// final sheet image frame
	simg := cai.Img.SubImage(image.Rect(sx, sy, sx+cai.FrameWidth, sy+cai.FrameHeight)).(*ebiten.Image)
	// Draw sub image from img Sheet
	screen.DrawImage(simg, op)
	cai.debugBounds(screen, simg, op)
}
