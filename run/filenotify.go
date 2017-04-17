package run

import (
	"github.com/rjeczalik/notify"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// Provides an easy way to run code triggered by changes on the
// filesystem
type FileTriggerRunner struct {
	folder     string
	absFolder  string
	action     Action
	recursive  bool
	events     chan notify.EventInfo
	stopSignal chan int
}

type Action func() error

// NewFileTriggerRunner constructor needs a path as first argument with
// no trailing slash.
// files in the given directory-root, starting with a . are ignored
func NewFileTriggerRunner(folder string, recursive bool, action Action) *FileTriggerRunner {
	return &FileTriggerRunner{
		folder:     folder,
		action:     action,
		recursive:  recursive,
		events:     make(chan notify.EventInfo, 1),
		stopSignal: make(chan int),
	}
}

func (ftr *FileTriggerRunner) Start() error {
	absFolder, err := filepath.Abs(ftr.folder)
	if err != nil {
		return err
	}
	ftr.absFolder = absFolder

	log.Printf("Folder to watch is %v", ftr.absFolder)
	if ftr.recursive {
		err = notify.Watch(ftr.absFolder+"/...", ftr.events, notify.All)
	} else {
		err = notify.Watch(ftr.absFolder, ftr.events, notify.All)
	}
	if err != nil {
		return err
	}

	defer notify.Stop(ftr.events)

	// run it once initially
	err = ftr.action()
	if err != nil {
		return err
	}

	for {
		select {
		case ev := <-ftr.events:
			err := ftr.eventIn(ev)
			if err != nil {
				return err
			}
		case <-ftr.stopSignal:
			return nil
		}
	}
}

func (ftr *FileTriggerRunner) eventIn(ev notify.EventInfo) error {
	relPath, err := filepath.Rel(ftr.absFolder, ev.Path())
	if err != nil {
		return err
	}

	// if one of the path parts is hidden (starts with a .) we will ignore this event
	for _, relItem := range strings.Split(relPath, string(filepath.Separator)) {
		if strings.HasPrefix(relItem, ".") {
			return nil
		}
	}

	log.Printf("Filewatcher triggered by %v", relPath)
	// file-changes often come in batches, as to not run it for every event once
	// we first wait
	time.Sleep(time.Millisecond * 100)
	// then drain channel
	for len(ftr.events) > 0 {
		<-ftr.events
	}
	// and after the draining we run it, while new events
	// might already be collecting
	return ftr.action()
}

func (ftr *FileTriggerRunner) Stop() {
	ftr.stopSignal <- 0
}
