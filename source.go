package config

import (
	"context"
	"time"

	"github.com/rusriver/config/v2/deepcopy"
)

type Source struct {
	Config        *Config
	ChCmd         chan *MsgCmd
	ChFlushSignal chan *MsgFlushSignal
	Opts          *NewSource_Options
}

type MsgCmd struct {
	Command  Command
	FullPath []string
	V        any
	Err      error
}

type Command int

const (
	Command_Set Command = iota
)

type MsgFlushSignal struct {
	ChDown chan struct{}
}

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
		ChFlushSignal: make(chan *MsgFlushSignal, opts.CommandBufferSize/10),
		Opts:          opts,
	}

	s.Config.Source = s

	go s.theWriteBackUpdaterG()

	return
}

func (s *Source) theWriteBackUpdaterG() {
	tick := time.NewTicker(s.Opts.UpdatePeriod)
	defer tick.Stop()

	executeTheQueue := func() {
		// "What the hell?" - would yell the avid adept of the Go Memory Model document.
		// This is RCU. Or RMW.

		c2 := s.Config.ChildCopy()
		c2.DataSubTree = deepcopy.Copy(c2.DataSubTree)

		qLen := len(s.ChCmd)
		for i := 0; i < qLen; i++ {
			msg := <-s.ChCmd
			switch msg.Command {
			case Command_Set:
				c2.NonThreadSafe_Set(msg.FullPath, msg.V)
			}
		}

		s.Config = c2
	}

	for {
		select {
		case <-s.Opts.Context.Done():
			// fmt.Println("+++DONE")
			return
		case msg := <-s.ChFlushSignal:
			// fmt.Println("+++FLUSH-SIG")
			executeTheQueue()
			// drain the queue, and broadcast notifications
			if msg.ChDown != nil {
				msg.ChDown <- struct{}{}
			}
			qLen := len(s.ChFlushSignal)
			for i := 0; i < qLen; i++ {
				msg = <-s.ChFlushSignal
				if msg.ChDown != nil {
					msg.ChDown <- struct{}{}
				}
			}
		case <-tick.C:
			// fmt.Println("+++TICK")
			executeTheQueue()
		}
	}
}
