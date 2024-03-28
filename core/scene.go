package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type Scene struct {
	CoreImage
}

func NewScene(ri []byte, x, y float64) *Scene {
	ci := NewCoreImage(ri, x, y)

	return &Scene{
		CoreImage: *ci,
	}
}

func NewSceneFromImage(i *ebiten.Image, x, y float64) *Scene {
	ci := NewCoreImageFromImage(i, x, y)

	return &Scene{
		CoreImage: *ci,
	}
}

func (s *Scene) SetImagePosition(x, y float64) {
	s.CoreImage.SetImagePosition(x, y)
}

func (s *Scene) Draw(screen *ebiten.Image) {
	s.SimpleDraw(screen)
}

type SceneObject struct {
	CoreImage
	Space   *resolv.Object
	Visible bool
}

func (so *SceneObject) SpaceObject() *resolv.Object {
	return so.Space
}

func (so *SceneObject) Draw(screen *ebiten.Image) {
	if so.Visible {
		screen.DrawImage(so.Img, so.ImgOpts)
	}
}
