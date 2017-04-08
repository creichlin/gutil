package run

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestFileTriggerRunner(t *testing.T) {
	countChannel := make(chan int, 10)
	fn := NewFileTriggerRunner("testdata/t1", func() error {
		countChannel <- 1
		return nil
	})

	go fn.Start()
	time.Sleep(time.Millisecond * 20)
	if len(countChannel) != 1 {
		t.Errorf("Func should be executed initially once but was %v", len(countChannel))
	}

	ioutil.WriteFile("testdata/t1/a.foo", []byte{0}, 0644)
	time.Sleep(time.Millisecond * 110)
	if len(countChannel) != 2 {
		t.Errorf("Func should be executed after file create but was %v", len(countChannel))
	}

	ioutil.WriteFile("testdata/t1/b.foo", []byte{100}, 0644)
	time.Sleep(time.Millisecond * 50)
	ioutil.WriteFile("testdata/t1/a.foo", []byte{100}, 0644)
	time.Sleep(time.Millisecond * 110)
	if len(countChannel) != 3 {
		t.Errorf("Func should be executed after file create and modify once but was %v", len(countChannel))
	}
}
