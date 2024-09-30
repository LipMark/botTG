package events

const (
	Unknown EventType = iota
	Message
)

type EventType int

type Event struct {
	Type EventType
	Text string
	Meta interface{}
}

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Proccesor interface {
	Process(e Event) error
}
