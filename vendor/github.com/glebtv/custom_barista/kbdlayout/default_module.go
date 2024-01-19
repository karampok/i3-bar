package kbdlayout

import (
	"strings"

	"barista.run/bar"
	"barista.run/colors"
)

func Get() bar.Module {
	return New().Output(func(m *Module, i Info) bar.Output {
		out := KbdOut{}
		la := strings.ToUpper(i.Layout)
		lseg := bar.PangoSegment(la).OnClick(m.Click)
		if la != "US" {
			lseg.Color(colors.Scheme("bad"))
		}
		out.Seg = append(out.Seg, lseg)
		for _, mod := range i.GetMods() {
			s := bar.PangoSegment(mod)
			if mod == "CAPS" {
				s.Color(colors.Scheme("bad"))
			}
			out.Seg = append(out.Seg, s)
		}
		return out
	})
}
