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
	"barista.run/modules/gsuite/gmail"
	"barista.run/modules/netinfo"
	"barista.run/modules/shell"
	"barista.run/modules/volume"
	"barista.run/modules/volume/pulseaudio"
	"barista.run/modules/wlan"
	"barista.run/oauth"
	"github.com/glebtv/custom_barista/kbdlayout"
	"github.com/martinohmann/barista-contrib/modules/ip"
	"github.com/martinohmann/barista-contrib/modules/ip/ipify"

	"github.com/karampok/i3-bar/blocks"
	"github.com/karampok/i3-bar/module"
	"github.com/karampok/i3-bar/xbacklight"
)

func main() {
	//System stuff
	var cl, gm bar.Module
	if out, err := setupGSuiteCreds(); err == nil {
		cl = calendar.New(out).Output(blocks.GCal).TimeWindow(8 * time.Hour)
		gm = gmail.New(out, "INBOX").Output(blocks.GMail)
	} else {
		cl = module.NewDummyModule("calendar error")
		gm = module.NewDummyModule("gmail error")
	}

	//System stuff
	lly := kbdlayout.New().Output(blocks.Layout) // TODO: get one that does not crash on setup-auth
	br := xbacklight.New().Output(blocks.Brightness)
	snd := volume.New(pulseaudio.DefaultSink()).Output(blocks.Snd)
	bat := battery.All().Output(blocks.Bat)

	//Bluetooth stuff
	adapter, mac, _ := "hci0", "4C:87:5D:58:8B:C2", "bluez_sink.4C_87_5D_58_8B_C2.headset_head_unit"
	blD := bluetooth.Device(adapter, mac).Output(blocks.PerBlueDevice("QC35"))
	bl := bluetooth.DefaultAdapter().Output(blocks.Bluetooth)

	//Net stuff
	wi := wlan.Named("wifi").Output(blocks.WLAN)
	tvpn := netinfo.Interface("tailscale0").Output(blocks.PerVPN("TS"))
	rvpn := netinfo.Interface("redhat0").Output(blocks.PerVPN("RH"))
	online := ip.New(ipify.Provider).Output(blocks.Online).Every(time.Minute)
	via := shell.New("bash", "-c", "ip route get 8.8.8.8 | grep -Po '(?<=dev )(\\S+)'").
		Output(blocks.ViaInterface).Every(time.Minute)

	ti := clock.Local().Output(time.Second, blocks.Clock)

	panic(barista.Run(
		cl, gm, lly, br, snd, bat, bl, blD, wi, tvpn, rvpn, online, via, ti,
	))
}

func setupGSuiteCreds() ([]byte, error) {
	fetch := func(path string) ([]byte, error) {
		return exec.Command("gopass", "show", "-o", path).CombinedOutput()
	}
	masterKey, err := fetch("services/barista")
	if err != nil {
		return masterKey, err
	}
	oauth.SetEncryptionKey(masterKey)
	return fetch("services/i3bar-rh-gsuite")
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
