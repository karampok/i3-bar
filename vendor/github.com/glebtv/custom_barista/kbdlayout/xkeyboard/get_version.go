package xkeyboard

import (
	"github.com/BurntSushi/xgb"
)

// GetVersionCookie is a cookie used only for GetVersion requests.
type GetVersionCookie struct {
	*xgb.Cookie
}

// GetVersion sends a checked request.
// If an error occurs, it will be returned with the reply by calling GetVersionCookie.Reply()
func GetVersion(c *xgb.Conn, MajorVersion byte, MinorVersion uint16) GetVersionCookie {
	c.ExtLock.RLock()
	defer c.ExtLock.RUnlock()
	if _, ok := c.Extensions["XKEYBOARD"]; !ok {
		panic("Cannot issue request 'GetVersion' using the uninitialized extension 'XKEYBOARD'. xkeyboard.Init(connObj) must be called first.")
	}
	cookie := c.NewCookie(true, true)
	c.NewRequest(getVersionRequest(c, MajorVersion, MinorVersion), cookie)
	return GetVersionCookie{cookie}
}

// GetVersionUnchecked sends an unchecked request.
// If an error occurs, it can only be retrieved using xgb.WaitForEvent or xgb.PollForEvent.
func GetVersionUnchecked(c *xgb.Conn, MajorVersion byte, MinorVersion uint16) GetVersionCookie {
	c.ExtLock.RLock()
	defer c.ExtLock.RUnlock()
	if _, ok := c.Extensions["XKEYBOARD"]; !ok {
		panic("Cannot issue request 'GetVersion' using the uninitialized extension 'XKEYBOARD'. xkeyboard.Init(connObj) must be called first.")
	}
	cookie := c.NewCookie(false, true)
	c.NewRequest(getVersionRequest(c, MajorVersion, MinorVersion), cookie)
	return GetVersionCookie{cookie}
}

// GetVersionReply represents the data returned from a GetVersion request.
type GetVersionReply struct {
	Sequence     uint16 // sequence number of the request for this reply
	Length       uint32 // number of bytes in this reply
	MajorVersion byte
	MinorVersion uint16
}

// Reply blocks and returns the reply data for a GetVersion request.
func (cook GetVersionCookie) Reply() (*GetVersionReply, error) {
	buf, err := cook.Cookie.Reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return getVersionReply(buf), nil
}

// getVersionReply reads a byte slice into a GetVersionReply value.
func getVersionReply(buf []byte) *GetVersionReply {
	v := new(GetVersionReply)
	b := 1 // skip reply determinant

	v.MajorVersion = buf[b]
	b += 1

	v.Sequence = xgb.Get16(buf[b:])
	b += 2

	v.Length = xgb.Get32(buf[b:]) // 4-byte units
	b += 4

	v.MinorVersion = xgb.Get16(buf[b:])
	b += 2

	return v
}

// Write request to wire for GetVersion
// getVersionRequest writes a GetVersion request to a byte slice.
func getVersionRequest(c *xgb.Conn, MajorVersion byte, MinorVersion uint16) []byte {
	size := 8
	b := 0
	buf := make([]byte, size)

	c.ExtLock.RLock()
	buf[b] = c.Extensions["XKEYBOARD"]
	c.ExtLock.RUnlock()
	b += 1

	buf[b] = 0 // request opcode
	b += 1

	xgb.Put16(buf[b:], uint16(size/4)) // write request size in 4-byte units
	b += 2

	buf[b] = MajorVersion
	b += 1

	b += 1 // padding

	xgb.Put16(buf[b:], MinorVersion)
	b += 2

	return buf
}
