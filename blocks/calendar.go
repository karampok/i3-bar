package blocks

import (
	"fmt"
	"strings"
	"time"

	"barista.run/bar"
	"barista.run/colors"
	"barista.run/modules/gsuite/calendar"
	"barista.run/outputs"
	"barista.run/pango"
)

func GCal(evts calendar.EventList) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-today")

	var e calendar.Event
	switch {
	case len(evts.InProgress) > 0:
		cl = colors.Scheme("bad")
		e = evts.InProgress[0]
	case len(evts.Alerting) > 0:
		e = evts.Alerting[0]
	case len(evts.Upcoming) > 0:
		e = evts.Upcoming[0]
	default:
		return outputs.Pango(ic, "empty").Color(cl)
	}
	untilStart := e.UntilStart()
	if untilStart < time.Hour*1 {
		cl = colors.Scheme("degraded")
	}
	minus := ""
	if untilStart < 0 {
		untilStart = -untilStart
		minus = "-"
	}
	txt := strings.ToLower(e.Summary)
	return outputs.Repeat(func(time.Time) bar.Output {
		return outputs.Pango(ic, spacer, fmt.Sprintf("%s (%v)(%v)  %s%dh%dm", txt, e.Response, e.EventStatus,
			minus, int(untilStart.Hours()), int(untilStart.Minutes())%60), "           ").Color(cl)
	}).Every(time.Minute)
}
