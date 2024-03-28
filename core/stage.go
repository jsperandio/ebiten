package core

import (
	log "github.com/sirupsen/logrus"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font"
)

type MouseHandler interface {
	OnClick()
	OnRightClick()
	OnMiddleClick()
}

type StageObject interface {
	Name() string
	SpaceObject() *resolv.Object
	Draw(screen *ebiten.Image)
}

type Stage struct {
	Name        string
	Gravity     float64
	Backgroud   *Scene
	DefaultFont font.Face

	StageSpace   *resolv.Space
	Screen       *ebiten.Image
	Events       []Event
	StageObjects []StageObject

	runningEventCount int
	player            *Player
	observers         map[string]Observer
}

func NewStage(name string, backgroud *Scene, space *resolv.Space, gravity float64) *Stage {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	defaultFont := truetype.NewFace(tt, &truetype.Options{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	return &Stage{
		Name:         name,
		Gravity:      gravity,
		Backgroud:    backgroud,
		Events:       make([]Event, 0),
		DefaultFont:  defaultFont,
		StageSpace:   space,
		StageObjects: []StageObject{},
		observers:    make(map[string]Observer),
	}
}

func (s *Stage) Update(gameTick int) {
	if s.runningEventCount == 0 {
		s.runEvents()
	}

	s.handleMouse()
}

func (s *Stage) Draw(screen *ebiten.Image) {
	s.Backgroud.Draw(screen)

	if len(s.StageObjects) > 0 {
		for _, o := range s.StageObjects {
			o.Draw(screen)
		}
	}

	if s.player != nil {
		s.player.Draw(screen)
	}
}

func (s *Stage) AddEvent(event Event) {
	s.Events = append(s.Events, event)
	log.Debug("Added event: ", event.GetType())
	log.Debug("Event: ", s.Events)
}

func (s *Stage) AddEvents(events []Event) {
	for _, e := range events {
		s.AddEvent(e)
	}
}

func (s *Stage) ExistsEvent(eventType string) bool {
	for _, e := range s.Events {
		if e.GetType() == eventType {
			return true
		}
	}

	return false
}

func (s *Stage) AppendEvent(event Event) bool {
	if !s.ExistsEvent(event.GetType()) {
		s.Events = append(s.Events, event)
		return true
	}
	return false
}

func (s *Stage) AddObject(object StageObject) {
	s.StageObjects = append(s.StageObjects, object)
	s.StageSpace.Add(object.SpaceObject())
}

func (s *Stage) AddObjects(objects []StageObject) {
	for _, o := range objects {
		s.AddObject(o)
		if o, ok := o.(Observable); ok {
			o.Register("Stage", s)
		}
	}
}

func (s *Stage) SetPlayer(player *Player) {
	s.player = player
	s.StageSpace.Add(player.spaceObject)
}

func (s *Stage) OnNotify(value interface{}) {
	s.Notify("GameFlow", value)
}

func (s *Stage) Register(id string, observer Observer) {
	s.observers[id] = observer
}

func (s *Stage) Deregister(id string) {
	delete(s.observers, id)
}

func (s *Stage) Notify(id string, data interface{}) {
	s.observers[id].OnNotify(data)
}

func (s *Stage) runEvents() {
	for _, e := range s.Events {
		e.Run(s.player, s.Screen)
		s.runningEventCount++
	}
}

func (s *Stage) handleMouse() {
	var o *resolv.Object = nil

	cx, cy := ebiten.CursorPosition()
	if o = s.StageSpace.CheckCells(cx, cy, 1, 1); o == nil {
		return
	}

	omh := s.getObjectWithMouseHandler(o.Tags())
	if omh == nil {
		return
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		omh.OnClick()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		omh.OnRightClick()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		omh.OnMiddleClick()
	}
}

// func (s *Stage) ApplyGravity() {
// 	// s.CurrPlayer.ApplyGravity(s.Gravity)
// }

func (s *Stage) getObjectWithMouseHandler(tags []string) MouseHandler {
	for _, o := range s.StageObjects {
		if o.SpaceObject().HasTags(tags...) {
			if mh, ok := o.(MouseHandler); ok {
				return mh
			}
			return nil
		}
	}
	return nil
}
