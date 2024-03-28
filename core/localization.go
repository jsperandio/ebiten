package core

type Direction float64

const (
	Right        Direction = 2.0
	Left         Direction = -2.0
	Up           Direction = -2.0
	Down         Direction = 2.0
	Static       Direction = 0.0
	ScreenWidth            = 630
	ScreenHeight           = 500
)

func (d Direction) Invert() Direction {
	return -d
}

func (d Direction) String() string {
	if d == Static {
		return "Static"
	} else if d == Right {
		return "Right"
	} else if d == Left {
		return "Left"
	} else {
		return "Unknown"
	}
}
