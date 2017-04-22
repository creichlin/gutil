package run

import (
	"github.com/rjeczalik/notify"
	"os"
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
	accept     map[string]bool
	events     chan notify.EventInfo
	stopSignal chan int
}

type Action func(event, path string) error

// NewFileTriggerRunner constructor needs a path as first argument with
// no trailing slash.
// files in the given directory-root, starting with a . are ignored
// if path is a file, the parent base will be resitered and the file
// is used to filter out events. This allow for deleting and recreating
// the file or renaming it.
func NewFileTriggerRunner(folder string, recursive bool, action Action) *FileTriggerRunner {
	return &FileTriggerRunner{
		folder:     folder,
		action:     action,
		recursive:  recursive,
		events:     make(chan notify.EventInfo, 1),
		stopSignal: make(chan int),
		accept:     map[string]bool{},
	}
}

func (ftr *FileTriggerRunner) Start() error {
	absFolder, err := filepath.Abs(ftr.folder)
	if err != nil {
		return err
	}

	node, err := os.Stat(absFolder)
	if err != nil {
		return err
	}

	if node.IsDir() {
		ftr.absFolder = absFolder
		if ftr.recursive {
			err = notify.Watch(ftr.absFolder+"/...", ftr.events, notify.All)
		} else {
			err = notify.Watch(ftr.absFolder, ftr.events, notify.All)
		}
		if err != nil {
			return err
		}
	} else { // node is file
		// we watch the parent folder and add the filename to the accept list
		// this way we can rename the file and create a new one and it still triggers
		// similar behaviour is done quite often by editors to allow atomic file writes
		dir, file := filepath.Split(absFolder)
		ftr.absFolder = dir
		ftr.accept[file] = true

		err = notify.Watch(ftr.absFolder, ftr.events, notify.All)
		if err != nil {
			return err
		}
	}

	defer notify.Stop(ftr.events)

	// run it once initially
	err = ftr.action("initial", ftr.absFolder)
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

	if len(ftr.accept) > 0 {
		if _, found := ftr.accept[relPath]; !found {
			return nil
		}
	} else {
		if dirIsHidden(relPath) {
			return nil
		}
	}

	// file-changes often come in batches, as to not run it for every event once
	// we first wait
	time.Sleep(time.Millisecond * 100)
	// then drain channel
	for len(ftr.events) > 0 {
		<-ftr.events
	}
	// and after the draining we run it, while new events
	// might already be collecting
	return ftr.action(ev.Event().String(), ev.Path())
}

func (ftr *FileTriggerRunner) Stop() {
	ftr.stopSignal <- 0
}

func dirIsHidden(path string) bool {
	// if one of the path parts is hidden (starts with a .) we will ignore this event
	for _, relItem := range strings.Split(path, string(filepath.Separator)) {
		if strings.HasPrefix(relItem, ".") {
			return true
		}
	}
	return false
}
