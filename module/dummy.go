package module

import (
	"barista.run/bar"
	"barista.run/base/value"
	"barista.run/outputs"
)

type DummyModule struct {
	outputFunc value.Value
}

func NewDummyModule(info string) *DummyModule {
	m := &DummyModule{}
	m.Output(func() bar.Output {
		return outputs.Textf(info)
	})
	return m

}

// Output sets the output format for the module.
func (m *DummyModule) Output(outputFunc func() bar.Output) *DummyModule {
	m.outputFunc.Set(outputFunc)
	return m
}

// Stream starts the module.
func (m DummyModule) Stream(sink bar.Sink) {
	outf := m.outputFunc.Get().(func() bar.Output)
	nextOutputFunc, done := m.outputFunc.Subscribe()
	defer done()
	for {
		sink.Output(outf())
		select {
		case <-nextOutputFunc:
			outf = m.outputFunc.Get().(func() bar.Output)
		}
	}
}
