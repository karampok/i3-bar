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
	c := colors.Scheme("dim-icon")
	ic := pango.Icon("material-today")

	var e calendar.Event
	switch {
	case len(evts.InProgress) > 0:
		e = evts.InProgress[0]
	case len(evts.Alerting) > 0:
		e = evts.Alerting[0]
	case len(evts.Upcoming) > 0:
		e = evts.Upcoming[0]
	default:
		return outputs.Pango(ic, "empty").Color(c)
	}
	untilStart := e.UntilStart()
	minus := ""
	if untilStart < 0 {
		untilStart = -untilStart
		minus = "-"
	}
	txt := strings.ToLower(e.Summary)
	return outputs.Repeat(func(time.Time) bar.Output {
		return outputs.Pango(ic, fmt.Sprintf("%s  %s%dh%dm", txt,
			minus, int(untilStart.Hours()), int(untilStart.Minutes())%60)).Color(c)
	}).Every(time.Minute)
}
