package blocks

import (
	"fmt"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/modules/battery"
	"barista.run/modules/netinfo"
	"barista.run/modules/shell"
	"barista.run/modules/volume"
	"barista.run/modules/wlan"
	"barista.run/outputs"
	"barista.run/pango"
	"barista.run/pango/icons/material"
	"github.com/glebtv/custom_barista/kbdlayout"
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
	})

	material.Load(home(".icons/material-design-icons"))
	//	material.Load(home(".icons/material-design-icons"))
	// mdi.Load(home("~/.icons/MaterialDesign-Webfont"))
	// typicons.Load(home("~/.icons/typicons.font"))
	//	ionicons.LoadMd(home("~/.icons/ionicons"))
	// fontawesome.Load(home("~/.icons/Font-Awesome"))

}

// Clock ...
func Clock(now time.Time) bar.Output {
	return outputs.Pango(
		pango.Icon("material-today").Color(colors.Scheme("dim-icon")),
		pango.Icon("material-access-time").Color(colors.Scheme("dim-icon")),
		now.Format("Mon 2 Jan "),
		now.Format("15:04:05"),
	).OnClick(click.RunLeft("gsimplecal"))
}

// Bat ...
func Bat(i battery.Info) bar.Output {
	disp := pango.Textf("Bat: %d%%", i.RemainingPct())
	if i.Status == battery.Disconnected || i.Status == battery.Unknown {
		return nil
	}
	iconName := "material-battery-std"
	icon := pango.Icon(iconName).Color(colors.Scheme("dim-icon"))
	if i.Status == battery.Charging {
		iconName = "material-battery-charging-full"
		icon = pango.Icon(iconName)
	}
	out := outputs.Pango(icon, disp)
	switch {
	case i.RemainingPct() <= 10:
		exec.Command("notify-send", "-t", "2000", "battery", "very low", "-u", "critical").Run()
		out.Urgent(true)
	case i.RemainingPct() <= 20:
		exec.Command("notify-send", "-t", "2000", "battery", "low", "-u", "normal").Run()
		out.Color(colors.Scheme("bad"))
	case i.RemainingPct() <= 30:
		out.Color(colors.Scheme("degraded"))
	}
	return out
}

// Snd ...
func Snd(v volume.Volume) bar.Output {
	iconName := "mute"
	pct := v.Pct()
	if v.Mute {
		return outputs.
			Pango(pango.Icon("material-volume-off"), spacer, pango.Textf("%2d%%", pct)).
			Color(colors.Scheme("degraded"))
	}
	if pct > 66 {
		iconName = "up"
	} else if pct > 33 {
		iconName = "down"
	}
	return outputs.Pango(
		pango.Icon("material-volume-"+iconName),
		spacer,
		pango.Textf("%2d%%", pct),
	)
}

// Xblight ...
var Xblight = shell.New("xbacklight").
	Every(time.Second).
	Output(func(s string) bar.Output {
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return outputs.Textf("%s", s)
		}
		return outputs.
			Pango(pango.Icon("material-brightness-medium"), spacer, fmt.Sprintf("%2.0f%%", i))
	})

// Layout ...
func Layout(m *kbdlayout.Module, i kbdlayout.Info) bar.Output {
	la := strings.ToLower(i.Layout)
	c := colors.Scheme("dim-icon")
	if la != "us" {
		c = colors.Scheme("degraded")
	}
	return outputs.Pango(
		pango.Icon("material-language"), spacer,
		fmt.Sprintf("%s", la),
	).OnClick(m.Click).Color(c)
}

// Net ...
func Net(s netinfo.State) bar.Output {
	if len(s.IPs) < 1 {
		return outputs.Text("No network").Color(colors.Scheme("bad"))
	}
	return outputs.
		Pango(pango.Icon("material-settings-ethernet"), spacer, fmt.Sprintf("%s:%v", s.Name, s.IPs[0]))
}

// Yubi ...
func Yubi(x bool, y bool) bar.Output {
	if x {
		return outputs.Textf("GPG").Background(colors.Scheme("bad"))
	}
	if y {
		return outputs.Textf("       ").Background(colors.Scheme("bad"))
	}
	return nil
}

// WLAN ...
func WLAN(i wlan.Info) bar.Output {
	switch {
	case !i.Enabled():
		return nil
	case i.Connecting():
		return outputs.Text("W: ...")
	case !i.Connected():
		return outputs.Text("W: down")
	default:
		return outputs.
			Pango(pango.Icon("material-signal-wifi-4-bar"), spacer, fmt.Sprintf("%s", i.SSID))
	}
}
