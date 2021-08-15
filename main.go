package main

import (
	"os/exec"
	"strings"
	"time"

	"barista.run"
	"barista.run/bar"
	"barista.run/modules/battery"
	"barista.run/modules/bluetooth"
	"barista.run/modules/clock"
	"barista.run/modules/gsuite/calendar"
	"barista.run/modules/netinfo"
	"barista.run/modules/volume"
	"barista.run/modules/volume/pulseaudio"
	"barista.run/modules/wlan"
	"barista.run/oauth"
	"barista.run/outputs"
	"github.com/glebtv/custom_barista/kbdlayout"
	"github.com/karampok/i3-bar/blocks"
	"github.com/karampok/i3-bar/module"
	"github.com/karampok/i3-bar/xbacklight"
	"github.com/martinohmann/barista-contrib/modules/keyboard"
	"github.com/martinohmann/barista-contrib/modules/keyboard/xkbmap"
)

func main() {
	var cl bar.Module
	if out, err := setupCalendarCreds(); err == nil {
		cl = calendar.New(out).Output(blocks.GCal)
	} else {
		cl = module.NewDummyModule("calendar error")
	}

	// TODO: get one that does not crash on setup-auth
	lly := kbdlayout.New().Output(blocks.Layout)
	ly := xkbmap.New("us", "gr").Output(func(l keyboard.Layout) bar.Output {
		return outputs.Textf("%s", strings.ToUpper(l.Name)).OnClick(func(e bar.Event) {
			switch e.Button {
			case bar.ButtonLeft:
				l.Next()
			case bar.ButtonRight:
				l.Previous()
			}
		})
	})
	br := xbacklight.New().Output(blocks.Brightness)
	snd := volume.New(pulseaudio.DefaultSink()).Output(blocks.Snd)
	bat := battery.All().Output(blocks.Bat)
	wi := wlan.Named("wifi").Output(blocks.WLAN)
	nt1 := netinfo.Interface("wifi").Output(blocks.Net)
	nt2 := netinfo.Interface("net").Output(blocks.Net)
	nt3 := netinfo.Interface("tailscale0").Output(blocks.Net)
	ti := clock.Local().Output(time.Second, blocks.Clock)

	adapter, mac, _ := "hci0", "09:A5:C1:A6:5C:77", "bluez_sink.09_A5_C1_A6_5C_77.headset_head_unit"
	//snd2 := volume.Sink(_).Output(blocks.Snd2)
	blD := bluetooth.Device(adapter, mac).Output(blocks.Blue)
	bl := bluetooth.DefaultAdapter().Output(blocks.Bluetooth)

	panic(barista.Run(
		blD, cl, lly, ly, br, snd, bat, wi, nt1, nt2, nt3, bl, ti,
	))
}

func setupCalendarCreds() ([]byte, error) {
	fetch := func(path string) ([]byte, error) {
		return exec.Command("gopass", "show", "-o", path).CombinedOutput()
	}
	masterKey, err := fetch("services/barista")
	if err != nil {
		return masterKey, err
	}
	oauth.SetEncryptionKey(masterKey)
	return fetch("services/rh-calendar")
}
