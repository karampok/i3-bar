package main

import (
	"time"

	"barista.run"
	"barista.run/modules/battery"
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

	panic(barista.Run(
		yu, ly, br, snd, bat, wi, nt, ti,
	))

}
