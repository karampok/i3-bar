package blocks

import (
	"barista.run/bar"
	"barista.run/colors"
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

	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-music-note")

	out := new(outputs.SegmentGroup)
	if m.Playing() {
		out.Append(ic).OnClick(ifLeft(m.Pause))
	} else {
		out.Append(ic).OnClick(ifLeft(m.Play)).Color(cl)
	}
	//	out.Append(outputs.Text(">|").OnClick(ifLeft(m.Next)))
	return out
}
