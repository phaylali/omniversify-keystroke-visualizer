//go:build linux
// +build linux

package input

type EventType string

const (
	KeyEvent   EventType = "key"
	ClickEvent EventType = "click"
)

type Event struct {
	Type  EventType
	Value string
}

type Listener interface {
	Start(events chan<- Event) error
	Stop() error
}

func NewListener() (Listener, error) {
	return newLinuxListener()
}
