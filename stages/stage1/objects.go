package stage1

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jsperandio/ebiten/core"
	"github.com/solarlune/resolv"
)

// ----------------------------------------------------------------------------------------
// Ceil
// This object is the ceil of the stage.
// ----------------------------------------------------------------------------------------
type ceilSceneObject struct {
	name string
	*core.SceneObject
}

func NewCeilSceneObject() *ceilSceneObject {
	w := float64(core.ScreenWidth)
	h := 75.0
	posx := 0.0
	posy := 0.0

	return &ceilSceneObject{
		name: CeilObjectLabel,
		SceneObject: &core.SceneObject{
			CoreImage: *core.NewCoreRectImage(ebiten.NewImage(int(w), int(h)), posx, posy, core.RectCoreBorder, nil),
			Visible:   core.DEBUG,
			Space:     resolv.NewObject(posx, posy, w, h, CeilObjectLabel),
		},
	}
}

func (o *ceilSceneObject) Name() string {
	return o.name
}

// ----------------------------------------------------------------------------------------
// Floor
// This object is the floor of the stage.
// ----------------------------------------------------------------------------------------
type floorSceneObject struct {
	name string
	*core.SceneObject
}

func NewFloorSceneObject() *floorSceneObject {
	w := float64(core.ScreenWidth)
	h := 130.0
	posx := 0.0
	posy := 340.0

	return &floorSceneObject{
		name: FloorObjectLabel,
		SceneObject: &core.SceneObject{
			CoreImage: *core.NewCoreRectImage(ebiten.NewImage(int(w), int(h)), posx, posy, core.RectCoreBorder, nil),
			Visible:   core.DEBUG,
			Space:     resolv.NewObject(posx, posy, w, h, FloorObjectLabel),
		},
	}
}

func (o *floorSceneObject) Name() string {
	return o.name
}
