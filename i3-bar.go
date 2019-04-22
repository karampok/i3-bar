// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// simple demonstrates a simpler i3bar built using barista.
// Serves as a good starting point for building custom bars.
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
	"github.com/karampok/i3-bar/yubikey"
)

func main() {
	yu := yubikey.New().Output(blocks.Yubi)
	ly := kbdlayout.New().Output(blocks.Layout)
	br := blocks.Xblight
	snd := volume.DefaultMixer().Output(blocks.Snd)
	bat := battery.All().Output(blocks.Bat)
	wi := wlan.Named("wifi").Output(blocks.WLAN)
	nt := netinfo.New().Output(blocks.Net)
	ti := clock.Local().Output(time.Second, blocks.Clock)

	panic(barista.Run(
		yu, ly, br, snd, bat, wi, nt, ti,
	))

}
