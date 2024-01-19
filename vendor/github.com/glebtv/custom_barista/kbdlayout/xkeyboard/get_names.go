package xkeyboard

import (
	"github.com/BurntSushi/xgb"
)

// GetVersionCookie is a cookie used only for GetVersion requests.
type GetNamesCookie struct {
	*xgb.Cookie
}

const XkbSymbolsNameMask = (1 << 2)
const XkbGroupNamesMask = (1 << 12)

// GetVersion sends a checked request.
// If an error occurs, it will be returned with the reply by calling GetVersionCookie.Reply()
func GetNames(c *xgb.Conn, which uint32) GetNamesCookie {
	c.ExtLock.RLock()
	defer c.ExtLock.RUnlock()
	if _, ok := c.Extensions["XKEYBOARD"]; !ok {
		panic("Cannot issue request 'GetVersion' using the uninitialized extension 'XTEST'. xtest.Init(connObj) must be called first.")
	}
	cookie := c.NewCookie(true, true)
	c.NewRequest(getNamesRequest(c, which), cookie)
	return GetNamesCookie{cookie}
}

//BYTE	type;
//BYTE	deviceID;
//CARD16	sequenceNumber B16;
//CARD32	length B32;
//CARD32	which B32;
//KeyCode	minKeyCode;
//KeyCode	maxKeyCode;
//CARD8	nTypes;
//CARD8	groupNames;
//CARD16	virtualMods B16;
//KeyCode	firstKey;
//CARD8	nKeys;
//CARD32	indicators B32;
//CARD8	nRadioGroups;
//CARD8	nKeyAliases;
//CARD16	nKTLevels B16;
//CARD32	pad3 B32;

// GetNamesReply represents the data returned from a GetNames request.
type GetNamesReply struct {
	DeviceId     byte
	Sequence     uint16 // sequence number of the request for this reply
	Length       uint32 // number of bytes in this reply
	Which        uint32
	MinKeyCode   byte
	MaxKeyCode   byte
	NTypes       byte
	GroupNames   byte
	VirtualMods  uint16
	FirstKey     byte
	NKeys        byte
	Indicators   uint32
	NRadioGroups byte
	NKeyAliases  byte
	NKTLevels    uint16
	ValueList    uint32
}

// Reply blocks and returns the reply data for a GetNames request.
func (cook GetNamesCookie) Reply() (*GetNamesReply, error) {
	buf, err := cook.Cookie.Reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return getNamesReply(buf), nil
}

// getNamesReply reads a byte slice into a GetNamesReply value.
func getNamesReply(buf []byte) *GetNamesReply {
	v := new(GetNamesReply)
	b := 1

	v.DeviceId = buf[b]
	b += 1

	v.Sequence = xgb.Get16(buf[b:])
	b += 2

	v.Length = xgb.Get32(buf[b:]) // 4-byte units
	b += 4

	v.Which = xgb.Get32(buf[b:]) // 4-byte units
	b += 4

	v.MinKeyCode = buf[b]
	b += 1

	v.MaxKeyCode = buf[b]
	b += 1

	v.NTypes = buf[b]
	b += 1

	v.GroupNames = buf[b]
	b += 1

	v.VirtualMods = xgb.Get16(buf[b:])
	b += 2

	v.FirstKey = buf[b]
	b += 1

	v.NKeys = buf[b]
	b += 1

	v.Indicators = xgb.Get32(buf[b:])
	b += 4

	v.NRadioGroups = buf[b]
	b += 1

	v.NKeyAliases = buf[b]
	b += 1

	v.NKTLevels = xgb.Get16(buf[b:])
	b += 2

	//pad3 b32
	b += 4

	v.ValueList = xgb.Get32(buf[b:]) // 4-byte units

	//fmt.Printf("value list: %x", v.ValueList)

	//spew.Dump(buf[b:])

	return v
}

// Write request to wire for GetNames
// getNamesRequest writes a GetNames request to a byte slice.
func getNamesRequest(c *xgb.Conn, which uint32) []byte {
	size := 12
	b := 0
	buf := make([]byte, size)

	c.ExtLock.RLock()
	buf[b] = c.Extensions["XKEYBOARD"]
	c.ExtLock.RUnlock()
	b += 1

	buf[b] = 17 // request opcode X_kbGetNames
	b += 1

	xgb.Put16(buf[b:], uint16(size/4)) // write request size in 4-byte units
	b += 2

	xgb.Put16(buf[b:], uint16(3)) // deviceSpec ?
	b += 2

	xgb.Put16(buf[b:], uint16(0)) // pad ?
	b += 2

	xgb.Put32(buf[b:], which)
	b += 4
	return buf
}
