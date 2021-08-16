package main

import (
	"os/exec"
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
	"github.com/martinohmann/barista-contrib/modules/ip"
	"github.com/martinohmann/barista-contrib/modules/ip/ipify"
)

func main() {
	var cl bar.Module
	if out, err := setupCalendarCreds(); err == nil {
		cl = calendar.New(out).Output(blocks.GCal).TimeWindow(8 * time.Hour)
	} else {
		cl = module.NewDummyModule("calendar error")
	}

	// TODO: get one that does not crash on setup-auth
	lly := kbdlayout.New().Output(blocks.Layout)
	// ly := module.NewXkblayoutState("us", "gr").Output(func(l keyboard.Layout) bar.Output {
	// 	return outputs.Textf("%s", strings.ToUpper(l.Name)).OnClick(func(e bar.Event) {
	// 		switch e.Button {
	// 		case bar.ButtonLeft:
	// 			l.Next()
	// 		case bar.ButtonRight:
	// 			l.Previous()
	// 		}
	// 	})
	// })
	br := xbacklight.New().Output(blocks.Brightness)
	snd := volume.New(pulseaudio.DefaultSink()).Output(blocks.Snd)
	bat := battery.All().Output(blocks.Bat)
	wi := wlan.Named("wifi").Output(blocks.WLAN)
	// TODO: make this to print default route interface
	gateway := ip.New(ipify.Provider).Output(func(info ip.Info) bar.Output {
		if info.Connected() {
			return outputs.Textf("online: %s", info.IP)
		}

		return outputs.Text("offline")
	})
	tvpn := netinfo.Interface("tailscale0").Output(blocks.PerVPN("TS"))
	rvpn := netinfo.Interface("redhat0").Output(blocks.PerVPN("RH"))
	ti := clock.Local().Output(time.Second, blocks.Clock)

	//adapter, mac, _ := "hci0", "09:A5:C1:A6:5C:77", "bluez_sink.09_A5_C1_A6_5C_77.headset_head_unit"
	adapter, mac, _ := "hci0", "4C:87:5D:58:8B:C2", "bluez_sink.4C_87_5D_58_8B_C2.headset_head_unit"
	//snd2 := volume.Sink(_).Output(blocks.Snd2)
	blD := bluetooth.Device(adapter, mac).Output(blocks.Blue)
	bl := bluetooth.DefaultAdapter().Output(blocks.Bluetooth)

	panic(barista.Run(
		cl, blD, lly, br, snd, bat, wi, gateway, tvpn, rvpn, bl, ti,
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
	// return of type []byte(`{"installed": {
	// 	"client_id":"%%GOOGLE_CLIENT_ID%%",
	// 	"project_id":"i3-barista",
	// 	"auth_uri":"https://accounts.google.com/o/oauth2/auth",
	// 	"token_uri":"https://www.googleapis.com/oauth2/v3/token",
	// 	"auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
	// 	"client_secret":"%%GOOGLE_CLIENT_SECRET%%",
	// 	"redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
	// }}`)
}
