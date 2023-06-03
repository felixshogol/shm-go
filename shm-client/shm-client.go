package shmclient

/*
#cgo CFLAGS: -g -I./shmh
#cgo LDFLAGS: -L ./shmlib -lshmclient -lpthread -lrt -lm
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "shmh/dfxp_shm_common.h"
#include "shmh/dfxp_shm_client.h"

*/
import "C"
import (
	"fmt"

	"github.com/golang/glog"
)

type ShmClient struct {

}

func NewShmClient () *ShmClient {
	glog.Info("Create ShmClient")
	return &ShmClient{}
}

func (shmclient ShmClient) InitShm ( ) error {
	glog.Info("InitShm")

	shmn := C.CString("/dfxp-shm")
	ret := C.ShmInit(shmn,2,0)
    
	if ret != 0 {
		return fmt.Errorf("ShmInit failed")
	}

	return nil 
}
