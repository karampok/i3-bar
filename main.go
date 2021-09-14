package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	_ "net/http/pprof"

	"barista.run"

	"barista.run/bar"
	"barista.run/modules/battery"
	"barista.run/modules/bluetooth"
	"barista.run/modules/clock"
	"barista.run/modules/netinfo"
	"barista.run/modules/shell"
	"barista.run/modules/wlan"
	"barista.run/oauth"

	"barista.run/modules/gsuite/calendar"
	"barista.run/modules/gsuite/gmail"

	"github.com/glebtv/custom_barista/kbdlayout"
	"github.com/karampok/i3-bar/blocks"
	"github.com/karampok/i3-bar/module"
	"github.com/karampok/i3-bar/xbacklight"
	"github.com/martinohmann/barista-contrib/modules/ip"
	"github.com/martinohmann/barista-contrib/modules/ip/ipify"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
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

	bat := battery.All().Output(blocks.Bat)

	//Bluetooth stuff
	qc35, qc25mac, _ := "hci0", "4C:87:5D:58:8B:C2", "bluez_sink.4C_87_5D_58_8B_C2.headset_head_unit"
	jabra75t5, jabra75mac, _ := "hci0", "70:BF:92:B2:95:E9", "bluez_sink.70_BF_92_B2_95_E9.headset_head_unit"
	blDq := bluetooth.Device(qc35, qc25mac).Output(blocks.PerBlueDevice("QC35"))
	blDj := bluetooth.Device(jabra75t5, jabra75mac).Output(blocks.PerBlueDevice("j75t"))
	bl := bluetooth.DefaultAdapter().Output(blocks.Bluetooth)

	//Net stuff
	wi := wlan.Named("wlp0s20f3").Output(blocks.WLAN)
	tvpn := netinfo.Interface("tailscale0").Output(blocks.PerVPN("TS"))
	rvpn := netinfo.Interface("tun0").Output(blocks.PerVPN("RH"))
	online := ip.New(ipify.Provider).Output(blocks.Online).Every(time.Minute)
	via := shell.New("bash", "-c", "ip route get 8.8.8.8 | grep -Po '(?<=dev )(\\S+)'").Output(blocks.ViaInterface).Every(time.Minute)
	ti := clock.Local().Output(time.Second, blocks.Clock)

	snd := shell.New("bash", "-c", "pulsemixer --list").Output(blocks.Snd2).Every(time.Second)

	panic(barista.Run(
		//cl, gm, lly, br, snd, bat, bl, blD, wi, tvpn, rvpn, online, via, ti,
		cl, gm, lly, snd, br, bat, bl, blDq, blDj, wi, online, tvpn, rvpn, via, ti,
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
