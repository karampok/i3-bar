package blocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsemixer(t *testing.T) {

	input := `Sink:            ID: sink-125, Name: Comet Lake PCH-LP cAVS HDMI / DisplayPort 3 Output, Mute: 0, Channels: 2, Volumes: ['69%', '69%']
Sink:            ID: sink-103, Name: Comet Lake PCH-LP cAVS HDMI / DisplayPort 2 Output, Mute: 0, Channels: 2, Volumes: ['74%', '74%']
Sink:            ID: sink-82, Name: Comet Lake PCH-LP cAVS HDMI / DisplayPort 1 Output, Mute: 0, Channels: 2, Volumes: ['74%', '74%']
Sink:            ID: sink-107, Name: Comet es, Mute: 0, Channels: 2, Volumes: ['69%', '69%'], Default
Source:          ID: source-114, Name: HD Pro Webcam C920 Pro, Mute: 1, Channels: 2, Volumes: ['100%', '100%'], Default
Source:          ID: source-65661, Name: Monitor of Comet Lake PCH-LP cAVS HDMI / DisplayPort 3 Output, Mute: 0, Channels: 2, Volumes: ['100%', '100%']
Source:          ID: source-65639, Name: Monitor of Comet Lake PCH-LP cAVS HDMI / DisplayPort 2 Output, Mute: 0, Channels: 2, Volumes: ['100%', '100%']
Source:          ID: source-65618, Name: Monitor of Comet Lake PCH-LP cAVS HDMI / DisplayPort 1 Output, Mute: 0, Channels: 2, Volumes: ['100%', '100%']
Source:          ID: source-65643, Name: Monitor of Comet Lake PCH-LP cAVS Speaker + Headphones, Mute: 0, Channels: 2, Volumes: ['100%', '100%']
Source:          ID: source-86, Name: Comet Lake PCH-LP cAVS Headphones Stereo Microphone, Mute: 0, Channels: 2, Volumes: ['99%', '99%']
Source:          ID: source-93, Name: Comet Lake PCH-LP cAVS Digital Microphone, Mute: 0, Channels: 2, Volumes: ['100%', '100%']
Source output:   ID: source-output-100, Name: PulseAudio Volume Control, Mute: 0, Channels: 1, Volumes: ['100%']
Source output:   ID: source-output-162, Name: PulseAudio Volume Control, Mute: 0, Channels: 1, Volumes: ['100%']
Source output:   ID: source-output-133, Name: PulseAudio Volume Control, Mute: 0, Channels: 1, Volumes: ['100%']
Source output:   ID: source-output-101, Name: PulseAudio Volume Control, Mute: 0, Channels: 1, Volumes: ['100%']
Source output:   ID: source-output-73, Name: PulseAudio Volume Control, Mute: 0, Channels: 1, Volumes: ['100%']`

	tests := map[string]struct {
		input                 string
		wantInput, wantOutput *device
	}{
		"t14s": {input,
			&device{name: "Comet es", dfault: true, volume: "69%"},
			&device{name: "HD Pro Webcam C920 Pro", alias: "C920Pa", mute: true, dfault: true, volume: "100%"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			i, o := parsePulsemixer(tc.input)
			assert.Equal(t, tc.wantInput, i)
			assert.Equal(t, tc.wantOutput, o)
		})
	}
}
