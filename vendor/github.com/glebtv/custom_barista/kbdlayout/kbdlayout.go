package kbdlayout

import (
	"log"
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/glebtv/custom_barista/kbdlayout/xkeyboard"
)

var Layouts []string
var Layout string
var Mods uint8
var Group byte
var Callback func(string, uint8)

func parseLayoutNames(names string) []string {
	//log.Println("parsing layout names:", names)
	parts := strings.Split(names, "+")
	ret := make([]string, 0)
	for i, part := range parts {
		//log.Println(i, part)
		if i == 0 {
			continue
		}
		if i == 1 {
			ret = append(ret, part)
		}
		if i == 2 {
			pt := strings.Split(part, ":")
			if len(pt) > 0 {
				ret = append(ret, pt[0])
			}
		}
	}
	Layouts = ret
	//spew.Dump(ret)
	return ret
}

var X *xgbutil.XUtil

func init() {
	var err error

	X, err = xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	conn := X.Conn()

	err = xkeyboard.Init(conn)
	if err != nil {
		log.Fatal(err)
	}

	// this is really UseExtension message
	vresp := xkeyboard.GetVersion(conn, 1, 0)
	//spew.Dump(vresp.Cookie.Reply())
	_, err = vresp.Reply()
	if err != nil {
		log.Fatal(err)
	}
}

func GetLayout() (string, uint8, error) {
	conn := X.Conn()

	nresp := xkeyboard.GetNames(X.Conn(), xkeyboard.XkbSymbolsNameMask)
	//log.Println("reply:")
	//spew.Dump(nresp)
	repl, err := nresp.Reply()
	if err != nil {
		return "", 0, err
	}

	//spew.Dump(vresp.Reply())
	//atom, err := xprop.Atom(X, "THE_ATOM_NAME", false)
	//if err == nil {
	//println("The atom number: ", atom.Atom)
	//}

	// GetAtomName for atom=0x1f2 (kbdDescPtr->names->symbols atom)
	anresp := xproto.GetAtomName(conn, xproto.Atom(repl.ValueList))
	anreply, err := anresp.Reply()
	if err != nil {
		return "", 0, err
	}
	//log.Println("layout names:", anreply.Name)
	names := parseLayoutNames(anreply.Name)
	//log.Println("parsed layout names:")
	//spew.Dump(names)

	sresp := xkeyboard.GetState(conn, 3)
	sreply, err := sresp.Reply()
	if err != nil {
		return "", 0, err
	}
	//spew.Dump(sreply)
	//log.Println("getstate reply, group:", sreply.Group)
	if len(names)-1 < int(sreply.Group) {
		log.Println("no group number", sreply.Group, "found in layout names", names)
		Layout = "?"
	} else {
		Layout = names[sreply.Group]
	}
	//spew.Dump(sreply)
	Mods = sreply.LatchedMods
	Group = sreply.Group
	return Layout, Mods, nil
}

func Subscribe(callback func(string, uint8)) {
	Callback = callback
	updateMaps := func() {
		layout, mods, err := GetLayout()
		//log.Println(layout, mods)
		if err != nil {
			panic(err)
		}
		callback(layout, mods)
	}

	// Doesn't work with XKEYBOARD and current X11 - see below
	//xevent.MappingNotifyFun(updateMaps).Connect(X, xevent.NoWindow)
	//xevent.KeymapNotifyFun(updateMaps).Connect(X, xevent.NoWindow)

	// @See https://stackoverflow.com/a/48407438/679778
	// this is XkbSelectEvents(d, XkbUseCoreKbd, XkbAllEventsMask, XkbAllEventsMask);
	// Dumped via xtrace
	// This messages have no reply from X11
	//xkeyboard.SelectEvents(X.Conn(), []byte{0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x05, 0x00})
	//xkeyboard.SelectEvents(X.Conn(), []byte{0x03, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x00, 0x07, 0x00})
	xkeyboard.SelectEvents(X.Conn(), []byte{0x00, 0x01, 0xff, 0x0f, 0x00, 0x00, 0xff, 0x0f, 0xff, 0x00, 0xff, 0x00})
	//log.Println("subscribed to mappings notification")

	// xevent doesn't support events we need, so loop manually
	//xevent.Main(X)
	go MainLoop(X, updateMaps)
}

func Switch(groupNum byte) {
	xkeyboard.LatchLockState(X.Conn(), []byte{0x00, 0x01, 0x00, 0x00, 0x01, groupNum, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00})
	if Callback != nil {
		layout, mods, err := GetLayout()
		if err != nil {
			panic(err)
		}
		Callback(layout, mods)
	}
}

func SwitchToNext() {
	next := Group + 1
	if int(next) > len(Layouts) {
		next = 0
	}
	Switch(next)
}

func MainLoop(xu *xgbutil.XUtil, updateMaps func()) {
	for {
		if xevent.Quitting(xu) {
			break
		}

		// Gobble up as many events as possible (into the queue).
		// If there are no events, we block.
		xevent.Read(xu, true)

		for !xevent.Empty(xu) {
			if xevent.Quitting(xu) {
				return
			}
			ev, err := xevent.Dequeue(xu)
			if err != nil {
				log.Fatal(err)
			}
			switch v := ev.(type) {
			case xkeyboard.XkbEvent:
				if v.Type == xkeyboard.XkbExtensionDeviceNotify {
					updateMaps()
				}
				//spew.Dump(ev)
			}
		}
	}
}
