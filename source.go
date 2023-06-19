package config

import (
	"context"
	"time"
)

type Source struct {
	Config        *Config
	ChCmd         chan *MsgCmd
	ChFlushSignal chan struct{}
}

type MsgCmd struct {
	Command Command
	ChDown  chan *MsgCmd
}

type Command int

const (
	Command_Set Command = iota
)

type NewSource_Options struct {
	Config            *Config
	Context           context.Context
	CommandBufferSize int
	UpdatePeriod      time.Duration
}

func NewSource(f ...func(opts *NewSource_Options)) (s *Source) {
	opts := &NewSource_Options{
		Context:           context.Background(),
		CommandBufferSize: 500,
		UpdatePeriod:      time.Second,
	}
	for _, f := range f {
		f(opts)
	}

	s = &Source{
		Config:        opts.Config,
		ChCmd:         make(chan *MsgCmd, opts.CommandBufferSize),
		ChFlushSignal: make(chan struct{}, 1),
	}

	s.Config.Source = s

	go s.theWriteBackUpdaterG()

	return
}

func (s *Source) theWriteBackUpdaterG() {

}
