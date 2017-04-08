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
	folder string
	action Action
	events chan notify.EventInfo
}

type Action func() error

// NewFileTriggerRunner constructor needs a path as first argument with
// no trailing slash.
// files in the given directory-root, starting with a . are ignored
func NewFileTriggerRunner(folder string, action Action) *FileTriggerRunner {
	return &FileTriggerRunner{
		folder: folder,
		action: action,
		events: make(chan notify.EventInfo, 1),
	}
}

func (ftr *FileTriggerRunner) Start() error {
	absFolder, err := filepath.Abs(ftr.folder)
	if err != nil {
		return err
	}

	log.Printf("Folder to watch is %v", absFolder)

	err = notify.Watch(absFolder+"/...", ftr.events, notify.All)
	if err != nil {
		return err
	}

	defer notify.Stop(ftr.events)

	// run it once initially
	err = ftr.action()
	if err != nil {
		return err
	}

eventLoop:
	for ev := range ftr.events {
		relPath, err := filepath.Rel(absFolder, ev.Path())
		if err != nil {
			return err
		}

		for _, relItem := range strings.Split(relPath, string(filepath.Separator)) {
			if strings.HasPrefix(relItem, ".") {
				continue eventLoop
			}
		}

		log.Printf("Filewatcher triggered by %v", relPath)
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
