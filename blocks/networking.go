package blocks

import (
	"fmt"

	"barista.run/bar"
	"barista.run/base/watchers/netlink"
	"barista.run/colors"
	"barista.run/modules/netinfo"
	"barista.run/modules/wlan"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/martinohmann/barista-contrib/modules/ip"
)

// Online ...
func Online(info ip.Info) bar.Output {
	cl := colors.Scheme("dim-icon")
	disp := pango.Textf("online")
	if !info.Connected() {
		cl = colors.Scheme("bad")
		disp = pango.Textf("offline")
	}
	return outputs.Pango(disp).Color(cl)
}

// ViaInterface ...
func ViaInterface(intf string) bar.Output {

	// TODO just send via wire or wifi
	// enp4s10f1                        pci 0000:04:0a.1
	// | | |  |                                |  |  | |
	// | | |  |                   domain <- 0000  |  | |
	// | | |  |                                   |  | |
	// en| |  |  --> ethernet                     |  | |
	// | |  |                                   |  | |
	// p4|  |  --> prefix/bus number (4)   <-- 04  | |
	// |  |                                      | |
	// s10|  --> slot/device number (10) <--    10 |
	// |                                        |
	// f1 --> function number (1)     <--       1

	cl := colors.Scheme("dim-icon")
	txt := fmt.Sprintf("via %s", intf)
	return outputs.Pango(txt).Color(cl)
}

// PerVPN returns custom per vpn interface output function.
func PerVPN(name string) func(s netinfo.State) bar.Output {
	disp := pango.Text(name)
	ret := func(s netinfo.State) bar.Output {
		cl := colors.Scheme("dim-icon")
		ic := pango.Icon("material-vpn-lock")

		if len(s.IPs) >= 1 {
			cl = colors.Scheme("good")
		}
		return outputs.
			Pango(ic, spacer, disp).Color(cl)
	}
	return ret
}

// WLAN ...
func WLAN(i wlan.Info) bar.Output {
	disp := pango.Textf(fmt.Sprintf("%s", i.SSID))
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-signal-wifi-4-bar")

	switch {
	case i.State == netlink.Down:
		disp = pango.Textf("")
		cl = colors.Scheme("degraded")
	case !i.Enabled():
		return nil
	case i.Connecting():
		return outputs.Text("W: ...")
	case !i.Connected():
		return outputs.Text("W: down").Color(colors.Scheme("degraded"))
	}

	return outputs.
		Pango(ic, disp).Color(cl)
}
