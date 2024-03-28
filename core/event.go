package core

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
)

type coreEvent struct {
	eventType string
	interval  time.Duration
	keepAlive bool
	ended     bool
	Action    func(obj ...interface{}) bool
	eventChan chan bool
}

type Event interface {
	SetAction(action func(obj ...interface{}) bool)
	Run(player *Player, s *ebiten.Image)
	Stop()
	GetType() string
	GetInterval() time.Duration
	IsEnded() bool
}

func NewEvent(eventType string, interval time.Duration, keepAlive bool) Event {
	return &coreEvent{
		eventType: eventType,
		interval:  interval,
		keepAlive: keepAlive,
	}
}

func (event *coreEvent) Run(player *Player, s *ebiten.Image) {
	log.Debug("Running event:", event.eventType)

	if event.keepAlive {
		event.eventChan = make(chan bool)
	}

	go func() {
		select {
		case <-event.eventChan:
			return
		default:
			ticker := time.NewTicker(event.interval)
			for range ticker.C {
				// log.Println("Range event:", event.Type)
				keeprunning := event.Action(player, s)
				if !keeprunning {
					event.ended = true
					return
				}
			}
		}
	}()
}

func (event *coreEvent) Stop() {
	if event.keepAlive {
		event.eventChan <- true
		event.ended = true
	}
}

func (event *coreEvent) GetType() string {
	return event.eventType
}

func (event *coreEvent) GetInterval() time.Duration {
	return event.interval
}

func (event *coreEvent) SetAction(action func(obj ...interface{}) bool) {
	event.Action = action
}

func (event *coreEvent) IsEnded() bool {
	return event.ended
}

type DrawableEvent struct {
	Fields        map[string]interface{}
	DrawItems     map[string]Drawable
	DrawAction    func(screen *ebiten.Image, tick int)
	durationTicks int

	eventType string
}

func NewDrawableEvent(name string, drawItems map[string]Drawable, fields map[string]interface{}, duration int) *DrawableEvent {
	return &DrawableEvent{
		Fields:        fields,
		DrawItems:     drawItems,
		durationTicks: duration,
		eventType:     name,
	}
}

func (de *DrawableEvent) Draw(screen *ebiten.Image, tick int) {
	if de.durationTicks < 0 {
		return
	}

	de.DrawAction(screen, tick)
	de.durationTicks--
}

func (de *DrawableEvent) SetDrawAction(action func(screen *ebiten.Image, tick int)) {
	de.DrawAction = action
}
