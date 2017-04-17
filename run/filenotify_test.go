package run

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// nolint: errcheck
func TestFileTriggerRunnerStop(t *testing.T) {
	os.MkdirAll("testdata/t1/a", 0755)
	defer os.RemoveAll("testdata/t1")

	fn := NewFileTriggerRunner("testdata/t1/a", false, func() error {
		return nil
	})

	go func() {
		time.Sleep(time.Millisecond * 10)
		fn.Stop()
	}()

	err := fn.Start()
	if err != nil {
		t.Errorf("Error should be nil but is %v", err)
	}
}

// nolint: errcheck
func TestFileTriggerRunnerStopByError(t *testing.T) {
	os.MkdirAll("testdata/t1/a", 0755)
	defer os.RemoveAll("testdata/t1")

	fn := NewFileTriggerRunner("testdata/t1/a", false, func() error {
		return errors.New("expected error")
	})

	err := fn.Start()
	if err.Error() != "expected error" {
		t.Errorf("Expected err is not actual err")
	}
}

// nolint: errcheck
func TestFileTriggerRunner(t *testing.T) {
	os.MkdirAll("testdata/t1/a", 0755)
	defer os.RemoveAll("testdata/t1")

	countChannel := make(chan int, 10)
	fn := NewFileTriggerRunner("testdata/t1", true, func() error {
		fmt.Printf("Yeah\n")
		countChannel <- 1
		return nil
	})

	go func() {
		t.Error(fn.Start())
	}()
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

	ioutil.WriteFile("testdata/t1/.a.foo", []byte{0}, 0644)
	time.Sleep(time.Millisecond * 110)
	if len(countChannel) != 3 {
		t.Errorf(". files should be ignored, count is %v", len(countChannel))
	}

	os.MkdirAll("testdata/t1/.b", 0755)
	time.Sleep(time.Millisecond * 110)
	if len(countChannel) != 3 {
		t.Errorf(". dirs should be ignored, count is %v", len(countChannel))
	}

	ioutil.WriteFile("testdata/t1/a/.foo", []byte{0}, 0644)
	time.Sleep(time.Millisecond * 110)
	if len(countChannel) != 3 {
		t.Errorf(". files in dirs should be ignored, count is %v", len(countChannel))
	}
}
