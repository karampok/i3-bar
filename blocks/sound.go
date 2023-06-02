package blocks

import (
	"fmt"
	"os/exec"
	"strings"

	"barista.run/bar"
	"barista.run/colors"
	"barista.run/outputs"
	"barista.run/pango"
)

// PulseAudio ...
func PulseAudio(intf string) bar.Output {
	out := new(outputs.SegmentGroup)
	i, o := parsePulsemixer(intf)
	out.Append(i.Output())
	out.Append(o.Output())
	return out
}

var aliases = map[string]string{
	"HD Pro Webcam C920 Pro":                                        "C920Pa",
	"Comet Lake PCH-LP cAVS Speaker + Headphones":                   "jack",
	"ThinkPad Thunderbolt 3 Dock USB Audio Digital Stereo (IEC958)": "dock",
	"ThinkPad Thunderbolt 3 Dock USB Audio Mono":                    "dock",
	"Comet Lake PCH-LP cAVS Digital Microphone":                     "-",
	"Comet Lake PCH-LP cAVS Headphones Stereo Microphone":           "jack",
	"Jabra Elite 75t":                  "j75t",
	"HD Pro Webcam C920 Analog Stereo": "C920",
}

type device struct {
	kind            string // source/sink
	name, id, alias string
	mute, dfault    bool
	volume          string
}

func (d *device) Output() *bar.Segment {
	ic := pango.Icon("material-volume-up")
	cl := colors.Scheme("dim-icon")

	name := d.name
	if d.alias != "" {
		name = d.alias
	}

	toggleMuteMic := func(bar.Event) {
		exec.Command("pulsemixer", "--toggle-mute", "--id", d.id).CombinedOutput()
		v, _ := exec.Command("pulsemixer", "--get-mute", "--id", d.id).CombinedOutput()
		if string(v) == "1\n" {
			exec.Command("light", "-s", "sysfs/leds/platform::micmute", "-S", "100").CombinedOutput()
			return
		}
		exec.Command("light", "-s", "sysfs/leds/platform::micmute", "-S", "0").CombinedOutput()
	}

	if d.kind == "input" && d.mute {
		ic = pango.Icon("material-mic-off")
		return outputs.Pango(ic, name).Color(cl).OnClick(toggleMuteMic)
	}

	if d.kind == "input" && !d.mute {
		ic = pango.Icon("material-mic")
		return outputs.Pango(ic, name).OnClick(toggleMuteMic)
	}

	txt := fmt.Sprintf("%s %s", name, d.volume)
	toggleMute := func(bar.Event) {
		exec.Command("pulsemixer", "--toggle-mute", "--id", d.id).CombinedOutput()
		v, _ := exec.Command("pulsemixer", "--get-mute", "--id", d.id).CombinedOutput()
		if string(v) == "1\n" {
			exec.Command("light", "-s", "sysfs/leds/platform::mute", "-S", "100").CombinedOutput()
			return
		}
		exec.Command("light", "-s", "sysfs/leds/platform::mute", "-S", "0").CombinedOutput()
	}
	if d.kind == "output" && d.mute {
		ic = pango.Icon("material-volume-off")
		if d.mute {
			ic = pango.Icon("material-volume-off")
		}
		return outputs.Pango(ic, txt).OnClick(toggleMute)
	}
	out := outputs.Pango(ic, txt).Color(cl).OnClick(toggleMute)
	return out
}

func parsePulsemixer(in string) (input *device, output *device) {
	for _, l := range strings.Split(in, "\n") {
		x := &device{
			name: "toBeDefined",
		}
		//Sink: ID: sink-107, Name: Comet es, Mute: 0, Channels: 2, Volumes: ['69%', '69%'], Default
		if !strings.Contains(l, "Default") {
			continue
		}
		strings.TrimSuffix(l, ", Default")
		x.dfault = true
		//Sink: ID: sink-107, Name: Comet es, Mute: 0, Channels: 2, Volumes: ['69%', '69%']

		ll := strings.SplitN(l, ":", 2)
		kind := ll[0] // Sink
		rest := ll[1] //ID: sink-107, Name: Come es, Mute: 0, Channels: 2, Volumes: ['69%', '69%']

		m := make(map[string]string)
		for _, entry := range strings.SplitN(rest, ",", 5) {
			split := strings.Split(entry, ":")
			k, v := split[0], split[1]
			m[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
		//[ID]: sink-107,
		//[Name]: Come es,
		//[Mute]: 0,
		//[Channels]: 2,
		//[Volumes]: ['69%', '69%']

		if v, ok := m["ID"]; ok {
			x.id = v
		}
		if v, ok := m["Name"]; ok {
			x.name = v
		}
		if v, ok := aliases[x.name]; ok {
			x.alias = v
		}

		x.mute = m["Mute"] == "1"

		vols := strings.Split(m["Volumes"], "'")
		x.volume = vols[1] // TODO: do regex!
		switch kind {
		case "Sink":
			x.kind = "output"
			output = x
		case "Source":
			x.kind = "input"
			input = x
		}
	}
	return
}
