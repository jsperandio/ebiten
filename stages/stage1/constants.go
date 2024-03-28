package stage1

import "time"

const (
	LifeDecreaseDamageMultiplier = 2
	LifeDecreaseInterval         = 1 * time.Second

	HungerIncreaseValue    = 2
	HungerIncreaseInterval = 5 * time.Second

	CeilObjectLabel  = "ceil"
	FloorObjectLabel = "floor"
)
