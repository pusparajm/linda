package event

const (
	TypeMessage = iota
	TypeStatusChange
)

// Abstract event (message or presence change) containing all necessary info for reply
type Event struct {
	Type      int
	Channel   string
	Status    string
	Text      string
	Timestamp string
	UserId    string
	Username  string
}

// Create new event instance
func New(eventType int) *Event {
	event := new(Event)
	event.Type = eventType
	return event
}
