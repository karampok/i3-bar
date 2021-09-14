package blocks

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/base/watchers/netlink"
	"barista.run/colors"
	"barista.run/modules/battery"
	"barista.run/modules/bluetooth"
	"barista.run/modules/netinfo"
	"barista.run/modules/volume"
	"barista.run/modules/wlan"
	"barista.run/outputs"
	"barista.run/pango"
	"barista.run/pango/icons/material"
	"github.com/glebtv/custom_barista/kbdlayout"
	"github.com/martinohmann/barista-contrib/modules/ip"
)

var spacer = pango.Text("   ").XXSmall()

func home(path string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, path)
}

func init() {
	colors.LoadFromMap(map[string]string{
		"good":     "#6d6",
		"degraded": "#dd6",
		"bad":      "#d66",
		"dim-icon": "#777",
		"dimmed":   "#777",
		"black":    "#ffffff",
		"white":    "#000000",
	})

	// if err:=	typicons.Load(home(".icon/typicons.font")); err!=nil{
	// panic(err)
	// }
	if err := material.Load(home(".icons/material-design-icons")); err != nil {
		panic(err)
	}
}

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

// Snd2 ...
func Snd2(intf string) bar.Output {
	space := pango.Text(" /  ").XXSmall()
	out := new(pango.Node)

	i, o := parsePulsemixer(intf)
	out.Concat(i.PangoNode()).Concat(space)
	out.Concat(o.PangoNode())
	// cl := colors.Scheme("dim-icon")
	// if i.mute{
	// cl = colors.Scheme("bad")
	// }
	// disp := pango.Textf("input: %s| output: %s", i.alias, o.alias)
	// return outputs.Pango(disp).Color(cl)
	return out
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
	disp := pango.Textf("via %s", intf)
	return outputs.Pango(disp).Color(cl)
}

// Clock ...
func Clock(now time.Time) bar.Output {
	return outputs.Pango(
		pango.Icon("material-access-time"),
		now.Format("Mon 2 Jan "),
		now.Format("15:04:05"),
	).OnClick(click.RunLeft("gsimplecal")).Color(colors.Scheme("dim-icon"))
}

// Bat ...
func Bat(i battery.Info) bar.Output {
	if i.Status == battery.Disconnected || i.Status == battery.Unknown {
		return outputs.Textf("%v", i.Status).Urgent(true)
	}

	disp := pango.Textf("%d%% (%2.1f Watt)", i.RemainingPct(), i.Power)
	iconName := "material-battery-std"
	icon := pango.Icon(iconName).Color(colors.Scheme("dim-icon"))
	cl := colors.Scheme("dim-icon")

	if i.Status == battery.Charging {
		disp = pango.Textf("Bat: %d%%", i.RemainingPct())
		iconName = "material-battery-charging-full"
		icon = pango.Icon(iconName)
	}

	var urgent bool
	switch {
	case i.RemainingPct() <= 10:
		urgent = true
	case i.RemainingPct() <= 20:
		cl = colors.Scheme("bad")
	case i.RemainingPct() <= 30:
		cl = colors.Scheme("degraded")
	}
	return outputs.Pango(icon, disp).Color(cl).Urgent(urgent)
}

// Snd ...
func Snd(v volume.Volume) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-volume-down")
	pct := v.Pct()
	if pct > 66 {
		ic = pango.Icon("material-volume-up")
		cl = colors.Scheme("degraded")
	}
	if v.Mute {
		ic = pango.Icon("material-volume-off")
		cl = colors.Scheme("dim-icon")
	}

	return outputs.
		Pango(ic, pango.Textf("%2d%%", pct)).Color(cl)
}

// Brightness ...
func Brightness(i int) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-brightness-medium")
	if i > 50 {
		cl = colors.Scheme("degraded")
		ic = pango.Icon("material-brightness-high")
	}
	return outputs.
		Pango(ic, fmt.Sprintf("%2.0d%%", i)).Color(cl)
}

// Layout ...
func Layout(m *kbdlayout.Module, i kbdlayout.Info) bar.Output {
	la := strings.ToLower(i.Layout)
	c := colors.Scheme("dim-icon")
	if la != "us" {
		c = colors.Scheme("degraded")
	}
	return outputs.Pango(
		pango.Icon("material-language"),
		fmt.Sprintf("%s", la),
	).OnClick(m.Click).Color(c)
}

// Net ...
func Net(s netinfo.State) bar.Output {
	disp := pango.Text("no network")
	cl := colors.Scheme("bad")
	ic := pango.Icon("material-settings-ethernet")

	if len(s.IPs) >= 1 {
		disp = pango.Textf(fmt.Sprintf("%s:%v", s.Name, s.IPs[0]))
		cl = colors.Scheme("dim-icon")
		return outputs.
			Pango(ic, spacer, disp).Color(cl)
	}
	return nil
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

// Yubi ...
func Yubi(x bool, y bool) bar.Output {
	if x {
		return outputs.Textf("g").Background(colors.Scheme("dim-icon")).MinWidth(200)
	}
	if y {
		return outputs.Textf("s").Background(colors.Scheme("dim-icon")).MinWidth(200)
	}
	return nil
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

// Bluetooth ...
func Bluetooth(s bluetooth.AdapterInfo) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-bluetooth")

	if !s.Powered {
		cl = colors.Scheme("degraded")
		return outputs.
			Pango(ic).Color(cl)
	}
	return nil
}

// Blue ...
func PerBlueDevice(alias string) func(bluetooth.DeviceInfo) bar.Output {
	return func(i bluetooth.DeviceInfo) bar.Output {
		dp := pango.Textf(fmt.Sprintf("%v", alias))
		cl := colors.Scheme("good")
		ic := pango.Icon("material-bluetooth-connected")

		if !i.Connected {
			return nil
		}
		return outputs.
			Pango(ic, dp).Color(cl)
	}
}

// Snd2 ...
// func Snd2(v volume.Volume) bar.Output {
// 	cl := colors.Scheme("good")
// 	ic := pango.Icon("material-volume-up")
// 	pct := v.Pct()
// 	if v.Mute {
// 		ic = pango.Icon("material-volume-off")
// 		cl = colors.Scheme("bad")
// 	}

// 	return outputs.
// 		Pango(ic, spacer, pango.Textf("%2d%%", pct)).Color(cl)
// }
