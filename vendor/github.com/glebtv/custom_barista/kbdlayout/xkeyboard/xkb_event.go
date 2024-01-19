// Based on https://github.com/BurntSushi/xgb/blob/master/xtest/xtest.go
package xkeyboard

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

const XkbNewKeyboardNotify = 0
const XkbMapNotify = 1
const XkbStateNotify = 2
const XkbControlsNotify = 3
const XkbIndicatorStateNotify = 4
const XkbIndicatorMapNotify = 5
const XkbNamesNotify = 6
const XkbCompatMapNotify = 7
const XkbBellNotify = 8
const XkbActionMessage = 9
const XkbAccessXNotify = 10
const XkbExtensionDeviceNotify = 11

func init() {
	xgb.NewEventFuncs[85] = XkbEventNew
}

type XkbEvent struct {
	Sequence   uint16
	Type       byte
	Time       xproto.Timestamp
	Root       xproto.Window
	Event      xproto.Window
	Child      xproto.Window
	RootX      int16
	RootY      int16
	EventX     int16
	EventY     int16
	State      uint16
	SameScreen bool
	// padding: 1 bytes
}

// KeyPressEventNew constructs a KeyPressEvent value that implements xgb.Event from a byte slice.
func XkbEventNew(buf []byte) xgb.Event {
	//spew.Dump(buf)
	v := XkbEvent{}
	b := 1 // don't read event number

	v.Type = buf[b]
	b += 1

	// TODO: THIS IS PROBABLY NOT CORRECTLY PARSED
	v.Sequence = xgb.Get16(buf[b:])
	b += 2

	v.Time = xproto.Timestamp(xgb.Get32(buf[b:]))
	b += 4

	v.Root = xproto.Window(xgb.Get32(buf[b:]))
	b += 4

	v.Event = xproto.Window(xgb.Get32(buf[b:]))
	b += 4

	v.Child = xproto.Window(xgb.Get32(buf[b:]))
	b += 4

	v.RootX = int16(xgb.Get16(buf[b:]))
	b += 2

	v.RootY = int16(xgb.Get16(buf[b:]))
	b += 2

	v.EventX = int16(xgb.Get16(buf[b:]))
	b += 2

	v.EventY = int16(xgb.Get16(buf[b:]))
	b += 2

	v.State = xgb.Get16(buf[b:])
	b += 2

	if buf[b] == 1 {
		v.SameScreen = true
	} else {
		v.SameScreen = false
	}
	b += 1

	//b += 1 // padding

	//log.Println(v.String())

	return v
}

// Not implemented
func (v XkbEvent) Bytes() []byte {
	return []byte{}
}

// SequenceId returns the sequence id attached to the KeyPress event.
// Events without a sequence number (KeymapNotify) return 0.
// This is mostly used internally.
func (v XkbEvent) SequenceId() uint16 {
	return v.Sequence
}

func (v XkbEvent) String() string {
	fieldVals := make([]string, 0, 12)
	fieldVals = append(fieldVals, xgb.Sprintf("Sequence: %d", v.Sequence))
	fieldVals = append(fieldVals, xgb.Sprintf("Type: %d", v.Type))
	fieldVals = append(fieldVals, xgb.Sprintf("Time: %d", v.Time))
	fieldVals = append(fieldVals, xgb.Sprintf("Root: %d", v.Root))
	fieldVals = append(fieldVals, xgb.Sprintf("Event: %d", v.Event))
	fieldVals = append(fieldVals, xgb.Sprintf("Child: %d", v.Child))
	fieldVals = append(fieldVals, xgb.Sprintf("RootX: %d", v.RootX))
	fieldVals = append(fieldVals, xgb.Sprintf("RootY: %d", v.RootY))
	fieldVals = append(fieldVals, xgb.Sprintf("EventX: %d", v.EventX))
	fieldVals = append(fieldVals, xgb.Sprintf("EventY: %d", v.EventY))
	fieldVals = append(fieldVals, xgb.Sprintf("State: %d", v.State))
	fieldVals = append(fieldVals, xgb.Sprintf("SameScreen: %t", v.SameScreen))
	return "XkbEvent {" + xgb.StringsJoin(fieldVals, ", ") + "}"
}
