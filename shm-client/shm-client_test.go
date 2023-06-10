package shmclient

import "testing"
import "github.com/stretchr/testify/assert"

func TestShmCmds(t *testing.T) {
	tests := []struct {
		testName string
		cmd      int
		param    interface{}
	}{
		{
			testName: "test Config traffic",
			cmd: 1, 
		},
	}

	shmclient := NewShmClient()
	for _, test := range tests {
		err := shmclient.InitShm()
		assert.NoError(t, err)
		t.Logf("t:%v",test)


	}
}
