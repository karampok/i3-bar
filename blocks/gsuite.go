package blocks

import (
	"fmt"
	"strings"
	"time"

	"barista.run/bar"
	"barista.run/colors"
	"barista.run/modules/gsuite/calendar"
	"barista.run/modules/gsuite/gmail"
	"barista.run/outputs"
	"barista.run/pango"
)

func outsideWorkingHours() bool {
	if h := time.Now().Hour(); h > 17 || h < 9 {
		return true
	}
	return false
}

func GCal(evts calendar.EventList) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-today")
	urgent := false

	if outsideWorkingHours() {
		return nil
	}

	var e calendar.Event
	switch {
	case len(evts.InProgress) > 0:
		e = evts.InProgress[0]
	case len(evts.Alerting) > 0:
		e = evts.Alerting[0]
	case len(evts.Upcoming) > 0:
		e = evts.Upcoming[0]
	default:
		return outputs.Pango(ic, "empty").Color(cl)
	}
	untilStart := e.UntilStart()
	if untilStart < 15*time.Minute {
		cl = colors.Scheme("degraded")
	}
	minus := ""
	if untilStart < 0 {
		cl = colors.Scheme("bad")
		urgent = true
		untilStart = -untilStart
		minus = "-"
	}
	txt := strings.ToLower(e.Summary)
	return outputs.Repeat(func(time.Time) bar.Output {
		return outputs.Pango(ic, spacer, fmt.Sprintf("%s (%v)  %s%dh%dm", txt, e.Response,
			minus, int(untilStart.Hours()), int(untilStart.Minutes())%60)).Color(cl).Urgent(urgent)
	}).Every(time.Minute)
}

func GMail(n gmail.Info) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-email")
	urgent := false

	if outsideWorkingHours() {
		return nil
	}

	v := n.Unread["INBOX"]
	if v > 0 {
		urgent = true
	}

	return outputs.Pango(ic, spacer, v).Color(cl).Urgent(urgent)
}
