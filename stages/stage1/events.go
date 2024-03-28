package stage1

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jsperandio/ebiten/core"
)

// Life Decrease Event
// This event will decrease the player's life by his hunger level.
// As a Slime the player always need to eat to live.
func LifeDecreaseEvent() core.Event {
	lifeDecreaseEvent := core.NewEvent("Life Decrease", LifeDecreaseInterval, true)

	lifeDecreaseEvent.SetAction(func(objs ...interface{}) bool {
		plyr := objs[0].(*core.Player)
		// screen := objs[1].(*ebiten.Image)

		log.Debug("Life Decrease Event Fired")
		if plyr.IsHungry() {
			log.Info("Life Decrease")
			// if screen != nil {
			// 	log.Info("Drawing Life Decrease")
			// 	core.DrawText(screen, fmt.Sprintf("%d", int(plyr.Hunger)*2), int(plyr.PosX+10), int(plyr.PosY+10), color.White, core.FontSizeLarge)
			// }
			return plyr.SufferDamage(int(plyr.Hunger) * LifeDecreaseDamageMultiplier)
		}
		return true
	})

	return lifeDecreaseEvent
}

// Hunger Increase Event
// This event will increase the Player's hunger level by 1 every event time.
// As a Slime the player always need to eat to live.
func HungerIncreaseEvent() core.Event {
	hungerIncreaseEvent := core.NewEvent("Hunger Increase", 5*time.Second, true)

	hungerIncreaseEvent.SetAction(func(objs ...interface{}) bool {
		plyr := objs[0].(*core.Player)
		if plyr != nil && !plyr.IsDead() {
			log.Debug("Hunger Increase Event Fired")
			plyr.IncreaseHunger(HungerIncreaseValue)
		}

		return true
	})

	return hungerIncreaseEvent
}
