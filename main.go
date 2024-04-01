package main

import (
	"os/exec"
	"time"

	_ "net/http/pprof"

	"barista.run"
	"barista.run/bar"
	"barista.run/modules/battery"
	"barista.run/modules/bluetooth"
	"barista.run/modules/clock"
	"barista.run/modules/gsuite/calendar"
	"barista.run/modules/media"
	"barista.run/modules/netinfo"
	"barista.run/modules/shell"
	"barista.run/modules/wlan"
	"barista.run/oauth"

	"github.com/karampok/i3-bar/blocks"
	"github.com/karampok/i3-bar/module"
	"github.com/karampok/i3-bar/xbacklight"
	"github.com/karampok/i3-bar/yubikey"

	"github.com/glebtv/custom_barista/kbdlayout"
)

func main() {
	//  go func() {
	//    log.Println(http.ListenAndServe("localhost:6060", nil))
	//  }()

	var yb bar.Module
	yb = yubikey.New()
	barista.Add(yb)
	// GSuite stuff
	var cl bar.Module
	if out, err := setupGSuiteCreds(); err == nil {
		cl = calendar.New(out).Output(blocks.GCal).TimeWindow(4 * time.Hour)
		//    gm = gmail.New(out, "INBOX").Output(blocks.GMail)
	} else {
		cl = module.NewDummyModule("calendar error")
		//    gm = module.NewDummyModule("gmail error")
	}
	barista.Add(cl)

	// System stuff
	lly := kbdlayout.New().Output(blocks.Layout) // TODO: get one that does not crash on setup-auth
	barista.Add(lly)
	br := xbacklight.New().Output(blocks.Brightness)
	barista.Add(br)
	bat := battery.All().Output(blocks.Bat)
	barista.Add(bat)
	// audio := shell.New("bash", "-c", "pulsemixer --list").Output(blocks.PulseAudio).Every(time.Second)
	//  barista.Add(audio)
	// Bluetooth stuff
	qc35, qc25mac, _ := "hci0", "4C:87:5D:58:8B:C2", "bluez_sink.4C_87_5D_58_8B_C2.headset_head_unit"
	jabra75t5, jabra75mac, _ := "hci0", "70:BF:92:B2:95:E9", "bluez_sink.70_BF_92_B2_95_E9.headset_head_unit"
	glbud, glbudmac, _ := "hci1", "24:29:34:A0:CD:99", "bluez_sink.24_29_34_A0_CD_99.headset_head_unit"
	blDq := bluetooth.Device(qc35, qc25mac).Output(blocks.PerBlueDevice("QC35"))
	barista.Add(blDq)
	blDj := bluetooth.Device(jabra75t5, jabra75mac).Output(blocks.PerBlueDevice("j75t"))
	barista.Add(blDj)
	blDg := bluetooth.Device(glbud, glbudmac).Output(blocks.PerBlueDevice("gB"))
	barista.Add(blDg)
	// bl := bluetooth.DefaultAdapter().Output(blocks.Bluetooth)
	// barista.Add(bl)

	// Net stuff
	// online := ip.New(ipify.Provider).Output(blocks.Online).Every(time.Minute)
	// barista.Add(online)
	via := shell.New("bash", "-c", "ip -json route get 8.8.8.8  |jq -r .[0].dev").
		Output(blocks.ViaInterface).Every(time.Minute)
	barista.Add(via)
	wifi := wlan.Named("wlp0s20f3").Output(blocks.WLAN)
	barista.Add(wifi)
	tvpn := netinfo.Interface("tailscale0").Output(blocks.PerVPN("TS"))
	barista.Add(tvpn)
	rvpn := netinfo.Interface("tun0").Output(blocks.PerVPN("RH"))
	barista.Add(rvpn)
	bvpn := netinfo.Interface("vpn0").Output(blocks.PerVPN("B"))
	barista.Add(bvpn)

	ti := clock.Local().Output(time.Second, blocks.Clock)
	barista.Add(ti)

	spotify := media.New("spotify").Output(blocks.Media)
	barista.Add(spotify)

	panic(barista.Run())
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
	//  "client_id":"%%GOOGLE_CLIENT_ID%%",
	//  "project_id":"i3-barista",
	//  "auth_uri":"https://accounts.google.com/o/oauth2/auth",
	//  "token_uri":"https://www.googleapis.com/oauth2/v3/token",
	//  "auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
	//  "client_secret":"%%GOOGLE_CLIENT_SECRET%%",
	//  "redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
	// }}`)
}
