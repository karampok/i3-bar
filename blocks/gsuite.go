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

func GCal(evts calendar.EventList) bar.Output {
	ic := pango.Icon("material-today")
	out := new(pango.Node).Concat(ic).Concat(spacer)

	//cec := colors.Scheme("bad")
	aec := colors.Scheme("degraded")
	uec := colors.Scheme("dim-icon")
	space := pango.Text("  /  ").Color(aec)

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
	// TODO: spacer to have black background.
	for _, evt := range evts.InProgress {
		txt := output(evt, fmt.Sprintf("until %v", evt.End.Format("15:04")))
		out.ConcatText(txt).Color(colors.Scheme("white")).Background(colors.Scheme("black")).Concat(space)
	}
	for _, evt := range evts.Alerting {
		txt := output(evt, fmt.Sprintf("@ %v", evt.Start.Format("15:04")))
		out.ConcatText(txt).Append(space).Color(aec)
	}
	for _, evt := range evts.Upcoming {
		txt := output(evt, fmt.Sprintf("@ %v", evt.Start.Format("15:04")))
		out.ConcatText(txt).Concat(space).Color(uec)
	}

	return outputs.Repeat(func(time.Time) bar.Output {
		return out
	}).Every(time.Minute)
}

func GMail(n gmail.Info) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-email")

	if outsideWorkingHours(time.Now()) {
		return nil
	}

	v := n.Unread["INBOX"]
	ret := outputs.Pango(spacer, ic, v, spacer).Color(cl)
	if v > 0 {
		return ret.Color(colors.Scheme("white")).Background(colors.Scheme("black"))
	}
	return ret

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
