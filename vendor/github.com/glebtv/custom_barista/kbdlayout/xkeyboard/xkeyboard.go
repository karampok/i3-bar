// Based on https://github.com/BurntSushi/xgb/blob/master/xtest/xtest.go
package xkeyboard

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func Init(c *xgb.Conn) error {
	reply, err := xproto.QueryExtension(c, 9, "XKEYBOARD").Reply()
	//spew.Dump(reply, err)

	switch {
	case err != nil:
		return err
	case !reply.Present:
		return xgb.Errorf("No extension named XKEYBOARD could be found on on the server.")
	}

	c.ExtLock.Lock()
	c.Extensions["XKEYBOARD"] = reply.MajorOpcode
	c.ExtLock.Unlock()
	return nil
}
