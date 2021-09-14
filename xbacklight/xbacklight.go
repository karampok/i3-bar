package xbacklight

import (
	"bytes"
	"os/exec"
	"strconv"
	"time"

	"barista.run/bar"
	"barista.run/base/value"
	"barista.run/outputs"
	"barista.run/timing"
)

// Module ...
type Module struct {
	level     value.Value
	outf      value.Value // of func(int) bar.Output
	scheduler *timing.Scheduler
}

func (*Module) get() int {
	b, err := exec.Command("light", "-G").CombinedOutput()
	if err != nil {
		return -10
	}

	bb := bytes.TrimSpace(b)
	s, err := strconv.ParseFloat(string(bb), 64)
	if err != nil {
		return -10
	}
	return int(s)
}

// New constructs a new xbacklight module.
func New() *Module {
	m := &Module{}
	m.scheduler = timing.NewScheduler()
	m.scheduler.Every(time.Second)
	return m
}

// Output sets the output format for the module.
func (m *Module) Output(f func(int) bar.Output) *Module {
	m.outf.Set(f)
	return m
}

// Stream starts the module.
func (m *Module) Stream(s bar.Sink) {
	outf := m.outf.Get().(func(int) bar.Output)
	value := m.get()

	changes, done := m.level.Subscribe()
	defer done()

	for {
		s.Output(outputs.Group(outf(value)).OnClick(m.click))
		select {
		case <-changes:
			value = m.level.Get().(int)
			if value == 0 {
				exec.Command("light", "-S", "1").Run()
			}
		case <-m.scheduler.C:
			m.level.Set(m.get())
		}
	}
}

func (m *Module) click(e bar.Event) {
	switch e.Button {
	case bar.ButtonLeft, bar.ScrollDown, bar.ScrollLeft, bar.ButtonBack:
		exec.Command("light", "-U", "1").Run()
	case bar.ButtonRight, bar.ScrollUp, bar.ScrollRight, bar.ButtonForward:
		exec.Command("light", "-A", "1").Run()
	}
	m.level.Set(m.get())
}
