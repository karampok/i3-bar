// XKEYBOARD-Request(135,5):
// LatchLockState opcode=0x87 opcode2=0x05 unparsed-data=0x00,0x01,0x00,0x00,0x01,0x01,0x00,0x00,0x1f,0x00,0x00,0x00;
// LatchLockState opcode=0x87 opcode2=0x05 unparsed-data=0x00,0x01,0x00,0x00,0x01,0x00,0x00,0x00,0x1f,0x00,0x00,0x00;

package xkeyboard

import (
	"github.com/BurntSushi/xgb"
)

type LatchLockStateCookie struct {
	*xgb.Cookie
}

// SelectEvents selects events from X11
// X11 sends no reply for this command
func LatchLockState(c *xgb.Conn, data []byte) LatchLockStateCookie {
	c.ExtLock.RLock()
	defer c.ExtLock.RUnlock()
	if _, ok := c.Extensions["XKEYBOARD"]; !ok {
		panic("Cannot issue request 'LatchLockState' using the uninitialized extension 'XKEYBOARD'. xkeyboard.Init(connObj) must be called first.")
	}
	cookie := c.NewCookie(false, false)
	c.NewRequest(latchLockStateRequest(c, data), cookie)
	return LatchLockStateCookie{cookie}
}

// Write request to wire for SelectEvents
// selectEventsRequest writes a SelectEvents request to a byte slice.
func latchLockStateRequest(c *xgb.Conn, data []byte) []byte {
	size := 4
	b := 0
	buf := make([]byte, size)

	c.ExtLock.RLock()
	buf[b] = c.Extensions["XKEYBOARD"]
	c.ExtLock.RUnlock()
	b += 1

	buf[b] = 5
	b += 1

	xgb.Put16(buf[b:], uint16((size+len(data))/4)) // write request size in 4-byte units
	b += 2

	buf = append(buf, data...)
	//spew.Dump(buf, len(buf))
	return buf
}
