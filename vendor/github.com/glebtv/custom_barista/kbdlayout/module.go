package kbdlayout

import (
	"strings"
	"sync"
	"time"

	"barista.run/bar"
)

const NUM_LOCK = 16
const CAPS_LOCK = 2

var lock sync.Mutex

type Info struct {
	// Layout code
	Layout string
	// Mods - modifier map
	Mods uint8
}

func (i Info) GetMods() []string {
	ret := make([]string, 0)
	if i.Mods&NUM_LOCK == NUM_LOCK {
		ret = append(ret, "NUM")
	}
	if i.Mods&CAPS_LOCK == CAPS_LOCK {
		ret = append(ret, "CAPS")
	}
	//log.Println("getmods", mods, ret)
	return ret
}

type Module struct {
	bar.Module
	bar.Sink
	outputFunc func(*Module, Info) bar.Output
}

func (m *Module) Stream(s bar.Sink) {
	forever := make(chan struct{})
	m.Sink = s
	<-forever
}

func (m *Module) Output(outputFunc func(*Module, Info) bar.Output) *Module {
	m.outputFunc = outputFunc
	return m
}

type KbdOut struct {
	Seg []*bar.Segment
}

func (k KbdOut) Segments() []*bar.Segment {
	return k.Seg
}

// New constructs an instance of the clock module with a default configuration.
func New() *Module {
	m := &Module{}
	// Default output template

	Subscribe(func(layout string, mods uint8) {
		i := Info{Layout: layout, Mods: mods}
		m.Sink.Output(m.outputFunc(m, i))
	})
	m.Output(func(m *Module, i Info) bar.Output {
		out := []*bar.Segment{}
		lseg := bar.TextSegment(strings.ToUpper(i.Layout)).OnClick(m.Click)
		out = append(out, lseg)
		for _, mod := range i.GetMods() {
			out = append(out, bar.TextSegment(mod))
		}
		return KbdOut{Seg: out}
	})
	go func() {
		time.Sleep(1 * time.Second)
		m.update()
	}()
	return m
}

func (m *Module) update() {
	layout, mods, err := GetLayout()
	if err != nil {
		layout = err.Error()
		mods = 0
	}
	i := Info{Layout: layout, Mods: mods}
	m.Sink.Output(m.outputFunc(m, i))
}

func (m *Module) Click(e bar.Event) {
	if e.Button == bar.ButtonLeft {
		SwitchToNext()
		m.update()
	}
}
