package main

import (
	"time"

	"barista.run"
	"barista.run/modules/battery"
	"barista.run/modules/bluetooth"
	"barista.run/modules/clock"
	"barista.run/modules/netinfo"
	"barista.run/modules/volume"
	"barista.run/modules/volume/pulseaudio"
	"barista.run/modules/wlan"

	"github.com/glebtv/custom_barista/kbdlayout"
	"github.com/karampok/i3-bar/blocks"
	"github.com/karampok/i3-bar/xbacklight"
)

func main() {
	//	yu := yubikey.New().Output(blocks.Yubi)
	ly := kbdlayout.New().Output(blocks.Layout)
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
		//blD, ly, br, snd, bat, wi, nt1, nt2, bl, ti,
		blD, ly, br, snd, bat, wi, nt1, nt2, nt3, bl, ti,
	))
}
