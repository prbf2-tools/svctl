package fsm

type Event string

const (
	EventStart   Event = "start"
	EventStop    Event = "stop"
	EventRestart Event = "restart"
	EventReset   Event = "reset"
)
