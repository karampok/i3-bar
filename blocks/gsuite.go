package blocks

import (
	"fmt"
	"time"

	"barista.run/bar"
	"barista.run/colors"
	"barista.run/modules/gsuite/calendar"
	"barista.run/modules/gsuite/gmail"
	"barista.run/outputs"
	"barista.run/pango"
)

func GCal(evts calendar.EventList) bar.Output {
	cec := colors.Scheme("bad")
	aec := colors.Scheme("degraded")
	uec := colors.Scheme("dim-icon")
	ic := pango.Icon("material-today")
	// urgent := false

	if outsideWorkingHours(time.Now()) {
		return nil
	}

	output := func(e calendar.Event, str string) string {
		resp := "UNKNOWN"
		switch e.Response {
		case calendar.StatusConfirmed:
			resp = "C"
		case calendar.StatusTentative:
			resp = "A"
		case calendar.StatusDeclined:
			resp = "D"
		case calendar.StatusUnresponded:
			resp = "U"
		}
		return fmt.Sprintf("%s (%v) %s", e.Summary, resp, str)
	}

	out := new(pango.Node).Concat(ic)
	space := pango.Text("  /  ")

	for _, evt := range evts.InProgress {
		txt := output(evt, fmt.Sprintf("ends at %v", evt.End))
		out.ConcatText(txt).Concat(space).Color(cec).Bold()
	}
	for _, evt := range evts.Alerting {
		txt := output(evt, fmt.Sprintf("starts at %v", evt.Start))
		out.ConcatText(txt).Append(space).Color(aec)
	}
	for _, evt := range evts.Upcoming {
		txt := output(evt, fmt.Sprintf("starts at %v", evt.Start))
		out.ConcatText(txt).Concat(space).Color(uec)
	}

	return outputs.Repeat(func(time.Time) bar.Output {
		return out
	}).Every(time.Minute)
}

func GMail(n gmail.Info) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-email")
	urgent := false

	if outsideWorkingHours(time.Now()) {
		return nil
	}

	v := n.Unread["INBOX"]
	if v > 0 {
		urgent = true
	}

	return outputs.Pango(ic, spacer, v).Color(cl).Urgent(urgent)
}

func outsideWorkingHours(t time.Time) bool {
	if d := t.Weekday(); d == time.Saturday || d == time.Sunday {
		return true
	}

	if h := t.Hour(); h > 17 || h < 9 {
		return true
	}

	return false
}
