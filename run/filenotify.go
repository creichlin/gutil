package run

import (
	"github.com/rjeczalik/notify"
	"time"
)

// Provides an easy way to run code triggered by changes on the
// filesystem

type FileTriggerRunner struct {
	folder string
	action Action
	events chan notify.EventInfo
}

type Action func() error

func NewFileTriggerRunner(folder string, action Action) *FileTriggerRunner {
	return &FileTriggerRunner{
		folder: folder,
		action: action,
		events: make(chan notify.EventInfo, 1),
	}
}

func (ftr *FileTriggerRunner) Start() error {
	if err := notify.Watch(ftr.folder, ftr.events, notify.All); err != nil {
		return err
	}
	defer notify.Stop(ftr.events)

	// run it once initially
	if err := ftr.action(); err != nil {
		return err
	}
	for range ftr.events {
		// file-changes come in batches, as to not run it for every event once
		// we first wait
		time.Sleep(time.Millisecond * 100)
		// then drain channel
		for len(ftr.events) > 0 {
			<-ftr.events
		}
		// and after the draining we run it, while new events
		// might already be collecting
		if err := ftr.action(); err != nil {
			return err
		}
	}
	return nil
}
