package events

import "time"

type Event interface {
	EventType() string
	OccurredAt() time.Time
}

type BaseEvent struct {
	Type string
	At   time.Time
}

func (e BaseEvent) EventType() string {
	return e.Type
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.At
}
