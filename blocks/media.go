package blocks

import (
	"barista.run/bar"
	"barista.run/modules/media"
	"barista.run/outputs"
	"barista.run/pango"
)

func ifLeft(dofn func()) func(bar.Event) {
	return func(e bar.Event) {
		if e.Button == bar.ButtonLeft {
			dofn()
		}
	}
}

func Media(m media.Info) bar.Output {
	if !m.Connected() {
		return nil
	}

	out := new(outputs.SegmentGroup)
	if m.Playing() {
		ic := pango.Icon("material-play-arrow")
		out.Append(ic).OnClick(ifLeft(m.Pause))
	} else {
		ic := pango.Icon("material-pause")
		out.Append(ic).OnClick(ifLeft(m.Play))
	}
	//	out.Append(outputs.Text(">|").OnClick(ifLeft(m.Next)))
	return out
}
