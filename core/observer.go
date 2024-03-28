package core

type Observable interface {
	Register(id string, observer Observer)
	Deregister(id string)
	Notify(id string, data interface{})
}

type Observer interface {
	OnNotify(value interface{})
}
