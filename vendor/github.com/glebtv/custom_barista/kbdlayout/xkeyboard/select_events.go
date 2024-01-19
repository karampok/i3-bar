package xkeyboard

import (
	"github.com/BurntSushi/xgb"
)

type SelectEventsCookie struct {
	*xgb.Cookie
}

const XkbAllEventsMask = 4095

// SelectEvents selects events from X11
// X11 sends no reply for this command
func SelectEvents(c *xgb.Conn, data []byte) SelectEventsCookie {
	c.ExtLock.RLock()
	defer c.ExtLock.RUnlock()
	if _, ok := c.Extensions["XKEYBOARD"]; !ok {
		panic("Cannot issue request 'SelectEvents' using the uninitialized extension 'XKEYBOARD'. xkeyboard.Init(connObj) must be called first.")
	}
	cookie := c.NewCookie(false, false)
	c.NewRequest(selectEventsRequest(c, data), cookie)
	return SelectEventsCookie{cookie}
}

// Write request to wire for SelectEvents
// selectEventsRequest writes a SelectEvents request to a byte slice.
func selectEventsRequest(c *xgb.Conn, data []byte) []byte {
	size := 4
	b := 0
	buf := make([]byte, size)

	c.ExtLock.RLock()
	buf[b] = c.Extensions["XKEYBOARD"]
	c.ExtLock.RUnlock()
	b += 1

	buf[b] = 1 // request opcode
	b += 1

	xgb.Put16(buf[b:], uint16((size+len(data))/4)) // write request size in 4-byte units
	b += 2

	buf = append(buf, data...)
	//spew.Dump(buf, len(buf))
	return buf
}
