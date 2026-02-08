//go:build linux
// +build linux

package input

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

type linuxListener struct {
	stopChan chan struct{}
}

type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

func newLinuxListener() (Listener, error) {
	return &linuxListener{
		stopChan: make(chan struct{}),
	}, nil
}

func (l *linuxListener) Start(events chan<- Event) error {
	go l.readDevices(events)
	return nil
}

func (l *linuxListener) readDevices(events chan<- Event) {
	entries, _ := filepath.Glob("/dev/input/event*")

	for _, path := range entries {
		go l.readDevice(path, events)
	}
}

func (l *linuxListener) readDevice(path string, events chan<- Event) {
	file, err := os.Open(path)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("LOG: Input failed to open device %s: %v\n", path, err))
		return
	}
	os.Stderr.WriteString(fmt.Sprintf("LOG: Input successfully opened device %s\n", path))
	defer file.Close()

	buffer := make([]byte, 24)

	for {
		select {
		case <-l.stopChan:
			return
		default:
			n, err := file.Read(buffer)
			if err != nil || n < 24 {
				continue
			}

			var ev inputEvent
			ev.Time.Sec = int64(binary.LittleEndian.Uint32(buffer[0:4]))
			ev.Time.Usec = int64(binary.LittleEndian.Uint32(buffer[4:8]))
			ev.Type = binary.LittleEndian.Uint16(buffer[8:10])
			ev.Code = binary.LittleEndian.Uint16(buffer[10:12])
			ev.Value = int32(binary.LittleEndian.Uint32(buffer[12:16]))

			if ev.Type == 1 && ev.Value == 1 {
				event := Event{
					Type:  KeyEvent,
					Value: mapKeyCode(ev.Code),
				}
				if event.Value != "" {
					select {
					case events <- event:
					default:
					}
				}
			}
		}
	}
}

func mapKeyCode(code uint16) string {
	keyMap := map[uint16]string{
		1:   "Esc",
		2:   "1",
		3:   "2",
		4:   "3",
		5:   "4",
		6:   "5",
		7:   "6",
		8:   "7",
		9:   "8",
		10:  "9",
		11:  "0",
		12:  "-",
		13:  "=",
		14:  "Back",
		15:  "Tab",
		16:  "Q",
		17:  "W",
		18:  "E",
		19:  "R",
		20:  "T",
		21:  "Y",
		22:  "U",
		23:  "I",
		24:  "O",
		25:  "P",
		26:  "[",
		27:  "]",
		28:  "Enter",
		29:  "LCtrl",
		30:  "A",
		31:  "S",
		32:  "D",
		33:  "F",
		34:  "G",
		35:  "H",
		36:  "J",
		37:  "K",
		38:  "L",
		39:  ";",
		40:  "'",
		41:  "`",
		42:  "LShift",
		43:  "\\",
		44:  "Z",
		45:  "X",
		46:  "C",
		47:  "V",
		48:  "B",
		49:  "N",
		50:  "M",
		51:  ",",
		52:  ".",
		53:  "/",
		54:  "RShift",
		55:  "*",
		56:  "LAlt",
		57:  "Space",
		58:  "Caps",
		59:  "F1",
		60:  "F2",
		61:  "F3",
		62:  "F4",
		63:  "F5",
		64:  "F6",
		65:  "F7",
		66:  "F8",
		67:  "F9",
		68:  "F10",
		87:  "F11",
		88:  "F12",
		96:  "NumpadEnter",
		97:  "RCtrl",
		100: "RAlt",
		102: "Home",
		103: "Up",
		104: "PgUp",
		105: "Left",
		106: "Right",
		107: "End",
		108: "Down",
		109: "PgDn",
		110: "Ins",
		111: "Del",
		125: "LWin",
		126: "RWin",
	}

	if key, ok := keyMap[code]; ok {
		return key
	}

	return fmt.Sprintf("Key%d", code)
}

func (l *linuxListener) Stop() error {
	close(l.stopChan)
	return nil
}

func init() {
	os.Stderr.WriteString("LOG: input/linux.go init called\n")
	runtime.LockOSThread()
}
