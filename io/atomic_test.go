package io

import (
	"io/ioutil"
	"testing"
)

// This test writes simultaneously many times to the same file
// The read checks if corrupt data arrives
// If i use ioutil.WriteFile() the method itself failes
// because the permissions are set wrong and also the read gets corrupt data
func TestWriteFileAtomic(t *testing.T) {
	for i := 0; i < 100; i++ {
		i := i
		go func() {
			for o := i; o < 100000; o++ {
				err := WriteFileAtomic("lala", []byte{byte(i), byte(i >> 8), byte(i), byte(i >> 8)}, 0x666)
				if err != nil {
					t.Error(err)
				}
			}
		}()
	}

	for i := 0; i < 5000; i++ {
		data, err := ioutil.ReadFile("lala")
		if err != nil {
			t.Error(err)
		}
		if len(data) != 4 || data[0] != data[2] || data[1] != data[3] {
			t.Errorf("Read corrupt data")
		}
	}

}
