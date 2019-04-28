package main

import (
	"time"

	"barista.run"
	"barista.run/modules/battery"
	"barista.run/modules/bluetooth"
	"barista.run/modules/clock"
	"barista.run/modules/netinfo"
	"barista.run/modules/volume"
	"barista.run/modules/wlan"

	"github.com/glebtv/custom_barista/kbdlayout"
	"github.com/karampok/i3-bar/blocks"
	"github.com/karampok/i3-bar/xbacklight"
	"github.com/karampok/i3-bar/yubikey"
)

func main() {
	yu := yubikey.New().Output(blocks.Yubi)
	ly := kbdlayout.New().Output(blocks.Layout)
	br := xbacklight.New().Output(blocks.Brightness)
	snd := volume.DefaultMixer().Output(blocks.Snd)
	bat := battery.All().Output(blocks.Bat)
	wi := wlan.Named("wifi").Output(blocks.WLAN)
	nt := netinfo.New().Output(blocks.Net)
	ti := clock.Local().Output(time.Second, blocks.Clock)

	adapter, mac, _ := "hci0", "09:A5:C1:A6:5C:77", "bluez_sink.09_A5_C1_A6_5C_77.headset_head_unit"
	bl := bluetooth.Device(adapter, mac).Output(blocks.Blue)
	//snd2 := volume.Sink(sink).Output(blocks.Snd2)

	panic(barista.Run(
		yu, ly, br, snd, bat, wi, nt, bl, ti,
	))

}
