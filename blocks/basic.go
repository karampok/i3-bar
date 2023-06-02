package blocks

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"barista.run/bar"
	"barista.run/colors"
	"barista.run/modules/battery"
	"barista.run/modules/bluetooth"
	"barista.run/modules/volume"
	"barista.run/outputs"
	"barista.run/pango"
	"barista.run/pango/icons/material"
	"github.com/glebtv/custom_barista/kbdlayout"
)

var (
	spacer       = pango.Text("   ").XXSmall()
	defaultColor colors.ColorfulColor
	focusTime    bool
)

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
		"white":    "#ffffff",
		"black":    "#000000",
	})
	defaultColor = colors.Scheme("#dim-icon")

	// if err:=	typicons.Load(home(".icon/typicons.font")); err!=nil{
	// panic(err)
	// }
	if err := material.Load(home(".icons/material-design-icons")); err != nil {
		panic(err)
	}
}

// Clock ...
func Clock(now time.Time) bar.Output {
	ic := pango.Icon("material-schedule")
	out := new(pango.Node).Concat(ic).ConcatText(now.Format("Mon 2 Jan "), now.Format("15:04:05"))
	if focusTime {
		return out
	}
	cl := colors.Scheme("dim-icon")
	return out.Color(cl)
}

// Bat ...
func Bat(i battery.Info) bar.Output {
	if i.Status == battery.Disconnected || i.Status == battery.Unknown {
		return outputs.Textf("%v", i.Status).Urgent(true)
	}

	txt := ""
	if i.RemainingPct() < 90 {
		txt = fmt.Sprintf("%d%%", i.RemainingPct())
	}
	if i.Power > 3.0 && i.Status == battery.Discharging {
		txt += fmt.Sprintf("(%2.1f Watt)", i.Power)
	}

	disp := pango.Text(txt)
	iconName := "material-battery-full"
	icon := pango.Icon(iconName).Color(colors.Scheme("dim-icon"))
	cl := colors.Scheme("dim-icon")

	if i.Status == battery.Charging {
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
	ic := pango.Icon("material-language")
	if la := strings.ToLower(i.Layout); la != "us" {
		return outputs.Pango(ic, fmt.Sprintf("%s", la)).Color(colors.Scheme("bad")).OnClick(m.Click)
	}
	return nil
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
