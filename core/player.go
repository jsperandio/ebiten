package core

import (
	"bytes"
	"fmt"
	"image"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

type Hunger int

const (
	Full Hunger = iota
	Satisfied
	Hungry
	Starving
)

func (h Hunger) String() string {
	switch h {
	case Starving:
		return "Starving"
	case Hungry:
		return "Hungry"
	case Full:
		return "Full"
	default:
		return "Satisfied"
	}
}

type Player struct {
	CoreAnimatedImage
	FaceDirection Direction
	Hunger        Hunger
	MoveSpeed     float64
	LifePoints    int
	DrawEvents    []*DrawableEvent

	velocityX    float64
	velocityY    float64
	jumpVelocity float64
	gameTick     *int

	spaceObject *resolv.Object
	collision   *resolv.Collision
	observers   map[string]Observer
}

func NewPlayer(playerSheet []byte, x, y float64, gameTick *int) *Player {
	cmi := NewCoreAnimatedImage(playerSheet, x, y, 0, 0, 32, 25, 4)
	p := &Player{
		CoreAnimatedImage: *cmi,
		FaceDirection:     Left,
		MoveSpeed:         1,
		LifePoints:        100,
		Hunger:            Full,
		spaceObject:       resolv.NewObject(x, y, 32, 25, "player"),
		collision:         nil,
		gameTick:          gameTick,
		velocityX:         0.0,
		velocityY:         0.0,
		jumpVelocity:      2.5,
		observers:         make(map[string]Observer),
	}
	return p
}

func (p *Player) Describe(screen *ebiten.Image) {
	xpos := 320
	mousex, mousey := ebiten.CursorPosition()

	ebitenutil.DebugPrintAt(screen, "Direction   : "+p.FaceDirection.String(), xpos, 0)
	ebitenutil.DebugPrintAt(screen, "PosX        : "+fmt.Sprintf("%f", p.PosX)+" Object X :"+fmt.Sprintf("%f", p.spaceObject.X), xpos, 10)
	ebitenutil.DebugPrintAt(screen, "PosY        : "+fmt.Sprintf("%f", p.PosY)+" Object Y :"+fmt.Sprintf("%f", p.spaceObject.Y), xpos, 20)
	ebitenutil.DebugPrintAt(screen, "MoveSpeed   : "+fmt.Sprintf("%f", p.MoveSpeed), xpos, 30)
	ebitenutil.DebugPrintAt(screen, "LifePoints  : "+fmt.Sprintf("%d", p.LifePoints), xpos, 40)
	ebitenutil.DebugPrintAt(screen, "Hunger      : "+p.Hunger.String(), xpos, 50)
	ebitenutil.DebugPrintAt(screen, "In colision : "+fmt.Sprintf("%t %v", p.InCollision(), p.CollisonTags()), xpos, 60)
	ebitenutil.DebugPrintAt(screen, "Collison    : "+fmt.Sprintf("%v", p.collision), xpos, 70)
	ebitenutil.DebugPrintAt(screen, "Tick        : "+fmt.Sprintf("%v", *p.gameTick), xpos, 80)
	ebitenutil.DebugPrintAt(screen, "FPS         : "+fmt.Sprintf("%v", ebiten.ActualFPS()), xpos, 90)
	ebitenutil.DebugPrintAt(screen, "TPS         : "+fmt.Sprintf("%v", ebiten.ActualTPS()), xpos, 100)
	ebitenutil.DebugPrintAt(screen, "MOUSE POS   : "+fmt.Sprintf("x=%v y=%v", mousex, mousey), xpos, 110)
}

func (p *Player) SetImagePosition(x, y float64) {
	p.CoreMoveableImage.SetImagePosition(x, y)
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.AnimatedDraw(screen, *p.gameTick, p.FaceDirection)
	p.runDrawEvents(screen)
	if DEBUG {
		p.Describe(screen)
	}
	p.debugSpaceBounds(screen, p.spaceObject, &ebiten.DrawImageOptions{})
}

func (p *Player) runDrawEvents(screen *ebiten.Image) {
	for _, e := range p.DrawEvents {
		e.Draw(screen, *p.gameTick)
	}
}

func (p *Player) appendDrawEvent(de *DrawableEvent) bool {
	p.DrawEvents = append(p.DrawEvents, de)
	return true
}

func (p *Player) debugSpaceBounds(screen *ebiten.Image, spaceObject *resolv.Object, imgOptions *ebiten.DrawImageOptions) {
	if !DEBUG {
		return
	}

	width := spaceObject.W
	height := spaceObject.H
	rectImg := NewCoreRectImage(ebiten.NewImage(int(width), int(height)), spaceObject.X, spaceObject.Y, RectCoreSpaceBorder, nil)
	screen.DrawImage(rectImg.Img, rectImg.ImgOpts)
}

// ------------------------------------------------------------------------------
//
//	Movement functions
//
// ------------------------------------------------------------------------------
func (p *Player) Update(stageGravity float64) {
	// stageGravity = 0
	if p.IsDead() {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.MoveRight()
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.MoveLeft()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustReleased(ebiten.KeyD) || inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyA) {
		p.StopHorizontal()
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.Jump()
	}

	p.ApplyGravity(stageGravity)
	p.updatePlayerPosition(p.spaceObject.X+p.velocityX, p.spaceObject.Y+p.velocityY)
}

func (p *Player) updatePlayerPosition(x, y float64) {
	p.spaceObject.X = x
	p.spaceObject.Y = y
	p.spaceObject.Update()
	p.updateSpritePosition(p.spaceObject.X+16, p.spaceObject.Y+12)
}

func (p *Player) updateSpritePosition(x, y float64) {
	p.PosX = x
	p.PosY = y
}

func (p *Player) ApplyGravity(gravity float64) {
	floorCollision, c := p.WillColide(0, p.velocityY+gravity)
	if gravity <= 0 || floorCollision {
		p.StopVertical()
	} else {
		p.velocityY += gravity
	}
	p.collision = c
}

func (p *Player) MoveRight() {
	p.FaceDirection = Right
	if _, p.collision = p.WillColide(float64(Right), 0); p.InCollision() {
		p.StopHorizontal()
		return
	}
	p.velocityX = p.MoveSpeed
}

func (p *Player) MoveLeft() {
	p.FaceDirection = Left

	if _, p.collision = p.WillColide(float64(Left), 0); p.InCollision() {
		p.StopHorizontal()
		return
	}

	p.velocityX = -p.MoveSpeed
}

func (p *Player) MoveUp() {
	if _, p.collision = p.WillColide(0, float64(Up)); p.InCollision() {
		return
	}
}

// func (p *Player) MoveDown(qty float64) {
// 	if _, p.collision = p.WillColide(0, float64(Down)); p.InCollision() {
// 		return
// 	}

// 	if qty == 0 {
// 		qty = p.MoveSpeed
// 	}

// 	p.CoreMoveableImage.MoveDown(qty)
// 	p.spaceObject.Y = p.PosY
// }

func (p *Player) Jump() {
	if p.velocityY == 0 {
		p.velocityY = -p.jumpVelocity
	}
}

// func (p *Player) ApplyGravity(gvt float64) {
// 	if p != nil {
// 		if (!p.InCollisionWithTag("floor")) && (*p.gameTick%2 == 0) {
// 			p.MoveDown(gvt)
// 		}
// 	}
// }

func (p *Player) StopHorizontal() {
	p.velocityX = 0
}

func (p *Player) StopVertical() {
	p.velocityY = 0
}

//------------------------------------------------------------------------------
// 							     Actions
//------------------------------------------------------------------------------

// WillColide returns true and the objects if the player will colide with the given direction
// with a Wall and other objects.
func (p *Player) WillColide(moveX, moveY float64) (bool, *resolv.Collision) {
	// Wall hit (-4 and -12 just for aesthetic reasons)
	// if ((p.spaceObject.X-12)+moveX <= 0) || ((p.spaceObject.X+p.spaceObject.W-4)+moveX >= ScreenWidth) {
	if ((p.spaceObject.X)+moveX <= 0) || ((p.spaceObject.X+p.spaceObject.W)+moveX >= ScreenWidth) {
		log.Debug("Collision: ", "Wall hit")
		return true, resolv.NewCollision()
	}

	c := p.spaceObject.Check(moveX, moveY)
	if c == nil {
		return false, nil
	}

	log.Debug("<WillColide> Collision: ", c.Objects[0].Tags())
	return true, c
}

func (p *Player) InCollisionWithTag(tag string) bool {
	if !p.InCollision() {
		return false
	}

	if p.collision.HasTags(tag) {
		return true
	}

	return false
}

func (p *Player) CollisonTags() string {
	if !p.InCollision() {
		return ""
	}

	var tags string
	for _, o := range p.collision.Objects {
		tags = tags + (strings.Join(o.Tags(), ","))
	}

	return tags
}

// SufferDamage will reduce the player's life points by the given amount and
// check if the player is dead
func (p *Player) SufferDamage(damage int) bool {
	p.LifePoints -= damage
	p.Notify("GUI", p.LifePoints)

	de := &DrawableEvent{
		DrawItems: map[string]Drawable{
			"AscendText": NewAscendingText(fmt.Sprint(damage), p.PosX, p.PosY-100, FontSizeMedium, ClrRed, 1),
		},
		durationTicks: 100,
		eventType:     "Damage Taken",
	}
	de.SetDrawAction(func(screen *ebiten.Image, tick int) {
		de.DrawItems["AscendText"].Draw(screen, tick)
	})
	p.appendDrawEvent(de)

	if p.IsDead() {
		p.LifePoints = 0
		p.Die()
		return false
	}

	return true
}

// IncreaseHunger increases the hunger of the player by desired amount
func (p *Player) IncreaseHunger(value int) {
	if p.Hunger != Starving {
		p.Hunger = Hunger(int(p.Hunger) + value)
	}
}

func (p *Player) InCollision() bool {
	return p.collision != nil
}

// //Check collision with top of the player
// func (p *Player) CheckCollisionTop() bool {
// 	return p.SpaceObject.SharesCells()
// }

// IsHungry returns true if the player is Hungry or Starving.
func (p *Player) IsHungry() bool {
	return p.Hunger > Satisfied
}

// IsDead returns true if the player is Dead.
func (p *Player) IsDead() bool {
	return p.LifePoints <= 0
}

// Die make all the player's actions when he dies.
func (p *Player) Die() {
	log.Info("Player died")

	p.setDeathAnimation()
	time.Sleep(400 * time.Millisecond)
	p.CoreAnimatedImage.Img.Clear()

	// Notify GameFlow that player died
	p.Notify("GameFlow", "PlayerDied")
}

// Get player Positions
func (p *Player) GetPosition() image.Point {
	return image.Point{int(p.PosX), int(p.PosY)}
}

// setDeathAnimation sets the death animation of the player with current sprite sheet.
func (p *Player) setDeathAnimation() {
	p.CoreAnimatedImage.FrameOX = 32
	p.CoreAnimatedImage.FrameOY = 50
}

// RestoreImage restores the player image from the raw sprite sheet bytes
func (p *Player) RestoreImage() {
	if len(p.RawImage) == 0 {
		return
	}

	// Recreate the image from raw bytes
	sceneDecoded, _, err := image.Decode(bytes.NewReader(p.RawImage))
	if err != nil {
		log.Error("Failed to restore player image: ", err)
		return
	}

	p.Img = ebiten.NewImageFromImage(sceneDecoded)
}

func (p *Player) Register(id string, observer Observer) {
	p.observers[id] = observer
}

func (p *Player) Deregister(id string) {
	delete(p.observers, id)
}

func (p *Player) Notify(id string, data interface{}) {
	p.observers[id].OnNotify(data)
}
