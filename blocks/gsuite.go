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

// GCal ...
func GCal(evts calendar.EventList) bar.Output {
	ic := pango.Icon("material-today")
	out := new(outputs.SegmentGroup)

	if outsideWorkingHours(time.Now()) {
		return nil
	}

	output := func(e calendar.Event, str string) string {
		resp := ""
		switch e.Response {
		case calendar.StatusConfirmed:
			resp = "\\C"
		case calendar.StatusTentative:
			resp = "\\A"
		case calendar.StatusDeclined:
			return ""
		case calendar.StatusUnresponded:
			resp = "\\U"
		}
		txt := strings.ToLower(e.Summary)
		if len(txt) > 20 {
			txt = txt[0:20]
		}
		return fmt.Sprintf("%s %s %s", resp, txt, str)
	}

	// TODO: allow only 3 events.
	total := 0
	if len(evts.InProgress) > 0 || len(evts.Alerting) > 0 {
		focusTime = true
	} else {
		focusTime = false
	}

	for _, evt := range evts.InProgress {
		total++
		txt := output(evt, fmt.Sprintf(" until %v", evt.End.Format("15:04")))
		f := func(bar.Event) {
			// how to hide it?
		}
		s := new(outputs.SegmentGroup).Append(pango.Text(txt).Medium()).
			Color(colors.Scheme("black")).Background(colors.Scheme("white")).OnClick(f)
		out.Append(s)
	}
	for _, evt := range evts.Alerting {
		total++
		if total > 2 {
			continue
		}
		txt := output(evt, fmt.Sprintf("at %v", evt.Start.Format("15:04")))
		s := outputs.Pango(ic, txt)
		out.Append(s)
	}
	for _, evt := range evts.Upcoming {
		total++
		if total > 2 {
			continue
		}
		if ignorePersonal(evt.Summary) {
			continue
		}
		txt := output(evt, fmt.Sprintf("@ %v", evt.Start.Format("15:04")))
		s := outputs.Pango(ic, txt).Color(colors.Scheme("dim-icon"))
		out.Append(s)
	}

	return outputs.Repeat(func(time.Time) bar.Output {
		return out
	}).Every(time.Minute)
}

func ignorePersonal(n string) bool {
	if strings.Contains(n, "AFK") || strings.Contains(n, "EOD") {
		return true
	}
	return false
}

// GMail ...
func GMail(n gmail.Info) bar.Output {
	if outsideWorkingHours(time.Now()) {
		return nil
	}

	v := n.Unread["INBOX"]
	if v > 0 {
		ic := pango.Icon("material-mail")
		return outputs.Pango("rh", ic)
	}
	return nil
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
