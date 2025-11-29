package gameover

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jsperandio/ebiten/core"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font"
)

type menuButton struct {
	*core.SceneObject
	name        string
	label       string
	observers   map[string]core.Observer
	normalColor color.Color
	hoverColor  color.Color
	textColor   color.Color
	fontSize    core.FontSize
	hovered     bool
}

func NewReturnToMenuObject() *menuButton {
	w := float64(core.ScreenWidth / 4)
	h := 50.0
	posx := float64(core.ScreenWidth/2) - (w / 2)
	posy := float64(core.ScreenHeight*2/3) - (h / 2)

	return newMenuButton("ReturnToMenuButton", "RETURN TO MENU", "ReturnToMenu", posx, posy, w, h)
}

func newMenuButton(name, label, spaceTag string, x, y, w, h float64) *menuButton {
	btn := &menuButton{
		name:        name,
		label:       label,
		observers:   make(map[string]core.Observer),
		normalColor: core.ClrGray,
		hoverColor:  core.ClrTeal,
		textColor:   core.ClrWhite,
		fontSize:    core.FontSizeMedium,
	}

	btn.SceneObject = &core.SceneObject{
		CoreImage: *core.NewCoreTextRectImage(
			ebiten.NewImage(int(w), int(h)),
			x,
			y,
			btn.normalColor,
			btn.label,
			btn.textColor,
			btn.fontSize,
		),
		Visible: true,
		Space:   resolv.NewObject(x, y, w, h, spaceTag),
	}

	return btn
}

func (b *menuButton) Name() string {
	return b.name
}

func (b *menuButton) OnClick() {
	b.Notify("Stage", b.name)
}

func (b *menuButton) OnRightClick() {}

func (b *menuButton) OnMiddleClick() {}

func (b *menuButton) Register(id string, observer core.Observer) {
	b.observers[id] = observer
}

func (b *menuButton) Deregister(id string) {
	delete(b.observers, id)
}

func (b *menuButton) Notify(id string, data interface{}) {
	b.observers[id].OnNotify(data)
}

func (b *menuButton) OnMouseEnter() {
	b.setHovered(true)
}

func (b *menuButton) OnMouseLeave() {
	b.setHovered(false)
}

func (b *menuButton) setHovered(state bool) {
	if b.hovered == state {
		return
	}

	b.hovered = state
	b.redrawFace()
}

func (b *menuButton) redrawFace() {
	fillColor := b.normalColor
	if b.hovered {
		fillColor = b.hoverColor
	}

	b.SceneObject.Img.Fill(fillColor)

	textRect, _ := font.BoundString(core.FontMap[b.fontSize], b.label)
	centerX := (b.SceneObject.Img.Bounds().Dx() / 2) - (textRect.Max.X.Ceil() / 2)
	centerY := (b.SceneObject.Img.Bounds().Dy() / 2) - (textRect.Min.Y.Floor() / 2)

	core.DrawText(b.SceneObject.Img, b.label, centerX, centerY, b.textColor, b.fontSize)
}
