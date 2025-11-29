package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GUI struct {
	Player    *Player
	HealthBar *HealthBar
}

func NewGUI(p *Player) *GUI {
	gui := &GUI{
		Player:    p,
		HealthBar: NewHealthBar(p),
	}

	p.Register("GUI", gui)

	return gui
}

func (g *GUI) Draw(screen *ebiten.Image, tick int) {
	g.HealthBar.Draw(screen, tick)
}

func (g *GUI) OnNotify(value interface{}) {
	g.HealthBar.UpdateHitpoints(value.(int))
}

type HealthBar struct {
	CoreImage
	Hitpoints int
}

func NewHealthBar(p *Player) *HealthBar {
	posx := 10.0
	posy := 10.0

	w := 200
	h := 10

	return &HealthBar{
		Hitpoints: p.LifePoints,
		CoreImage: *NewCoreRectImage(ebiten.NewImage(int(w), int(h)), posx, posy, RectCoreFill, ClrGreen),
	}
}

func (hb *HealthBar) Draw(screen *ebiten.Image, tick int) {
	hb.CoreImage.SimpleDraw(screen)
}

func (hb *HealthBar) UpdateHitpoints(value int) {
	hb.Hitpoints = value
	hb.redraw()
}

func (hb *HealthBar) redraw() {
	// Refill with green background
	hb.Img.Fill(ClrGreen)

	// Draw red line for damage if hitpoints are less than 100
	if hb.Hitpoints < 100 {
		midY := float32((hb.X0Y1().Y - hb.X0Y0().Y) / 2.0)
		endX := float32(hb.X1Y1().X)
		rprDmg := endX - float32(2*(100-hb.Hitpoints))
		vector.StrokeLine(hb.Img, endX, midY, rprDmg, midY, 10, ClrRed, false)
	}
}

func (hb *HealthBar) Reset() {
	hb.Hitpoints = 100
	hb.redraw()
}
