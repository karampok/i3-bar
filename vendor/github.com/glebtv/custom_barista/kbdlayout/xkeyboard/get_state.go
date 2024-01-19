package xkeyboard

import (
	"github.com/BurntSushi/xgb"
)

// GetStateCookie is a cookie used only for GetState requests.
type GetStateCookie struct {
	*xgb.Cookie
}

// GetState sends a checked request.
// If an error occurs, it will be returned with the reply by calling GetStateCookie.Reply()
func GetState(c *xgb.Conn, deviceSpec uint16) GetStateCookie {
	c.ExtLock.RLock()
	defer c.ExtLock.RUnlock()
	if _, ok := c.Extensions["XKEYBOARD"]; !ok {
		panic("Cannot issue request 'GetState' using the uninitialized extension 'XKEYBOARD'. xkeyboard.Init(connObj) must be called first.")
	}
	cookie := c.NewCookie(true, true)
	c.NewRequest(getStateRequest(c, deviceSpec), cookie)
	return GetStateCookie{cookie}
}

// GetStateUnchecked sends an unchecked request.
// If an error occurs, it can only be retrieved using xgb.WaitForEvent or xgb.PollForEvent.
func GetStateUnchecked(c *xgb.Conn, deviceSpec uint16) GetStateCookie {
	c.ExtLock.RLock()
	defer c.ExtLock.RUnlock()
	if _, ok := c.Extensions["XKEYBOARD"]; !ok {
		panic("Cannot issue request 'GetState' using the uninitialized extension 'XKEYBOARD'. xkeyboard.Init(connObj) must be called first.")
	}
	cookie := c.NewCookie(false, true)
	c.NewRequest(getStateRequest(c, deviceSpec), cookie)
	return GetStateCookie{cookie}
}

// GetStateReply represents the data returned from a GetState request.
type GetStateReply struct {
	Sequence         uint16 // sequence number of the request for this reply
	Length           uint32 // number of bytes in this reply
	Type             byte
	DeviceId         byte
	Mods             byte
	BaseMods         byte
	LatchedMods      byte
	LockedMods       byte
	Group            byte
	LockedGroup      byte
	BaseGroup        uint16
	LatchedGroup     uint16
	CompatState      byte
	GrabMods         byte
	CompatGrabMods   byte
	LookupMods       byte
	CompatLookupMods byte
	PtrBtnState      uint16
}

// Reply blocks and returns the reply data for a GetState request.
func (cook GetStateCookie) Reply() (*GetStateReply, error) {
	buf, err := cook.Cookie.Reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return getStateReply(buf), nil
}

// getStateReply reads a byte slice into a GetStateReply value.
func getStateReply(buf []byte) *GetStateReply {
	v := new(GetStateReply)
	b := 1 // skip reply determinant

	v.Type = buf[b]
	b += 1

	v.Sequence = xgb.Get16(buf[b:])
	b += 2

	v.Length = xgb.Get32(buf[b:])
	b += 4

	v.DeviceId = buf[b]
	b += 1

	v.Mods = buf[b]
	b += 1

	v.BaseMods = buf[b]
	b += 1

	v.LatchedMods = buf[b]
	b += 1

	v.LockedMods = buf[b]
	b += 1

	v.Group = buf[b]
	b += 1

	v.LockedGroup = buf[b]
	b += 1

	v.BaseGroup = xgb.Get16(buf[b:])
	b += 2

	v.LatchedGroup = xgb.Get16(buf[b:])
	b += 2

	v.CompatState = buf[b]
	b += 1

	v.GrabMods = buf[b]
	b += 1

	v.CompatGrabMods = buf[b]
	b += 1

	v.LookupMods = buf[b]
	b += 1

	v.CompatLookupMods = buf[b]
	b += 1

	b += 1 //pad1

	v.PtrBtnState = xgb.Get16(buf[b:])
	b += 2

	b += 2 //pad2
	b += 4 //pad3

	return v
}

// Write request to wire for GetState
// getStateRequest writes a GetState request to a byte slice.
func getStateRequest(c *xgb.Conn, deviceSpec uint16) []byte {
	size := 8
	b := 0
	buf := make([]byte, size)

	c.ExtLock.RLock()
	buf[b] = c.Extensions["XKEYBOARD"]
	c.ExtLock.RUnlock()
	b += 1

	buf[b] = 4 // request opcode
	b += 1

	xgb.Put16(buf[b:], uint16(size/4)) // write request size in 4-byte units
	b += 2

	xgb.Put16(buf[b:], deviceSpec)
	b += 2

	xgb.Put16(buf[b:], uint16(0)) // pad
	b += 2

	return buf
}
