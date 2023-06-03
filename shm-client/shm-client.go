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

type SHM_CMD int

const (
	DFPX_SHM_CMD_NONE SHM_CMD = iota
	DFPX_SHM_CMD_CONFIG
	DFPX_SHM_CMD_START
	DFPX_SHM_CMD_STOP
	DFPX_SHM_CMD_SHUTDOWN
	DFPX_SHM_CMD_ADD_IP_GTP
	DFPX_SHM_CMD_DEL_IP_GTP
	DFPX_SHM_CMD_GET_STATS
)

type ShmConfig struct {
	cmd SHM_CMD
}

type ShmClient struct {
}

func NewShmClient() *ShmClient {
	glog.Info("Create ShmClient")
	return &ShmClient{}
}

func (shmclient ShmClient) InitShm() error {
	glog.Info("InitShm")

	shmn := C.CString("/dfxp-shm")
	ret := C.ShmInit(shmn, 2, 0)

	if ret != 0 {
		return fmt.Errorf("ShmInit failed")
	}

	return nil
}

func (shmclient ShmClient) ShmWriteSart() error {

	shmcfg := &ShmConfig{}
	shmcfg.cmd = DFPX_SHM_CMD_START
    return shmclient.ShmWrite(shmcfg)

}

func (shmclient ShmClient) ShmWriteStop() error {

	shmcfg := &ShmConfig{}
	shmcfg.cmd = DFPX_SHM_CMD_STOP
    return shmclient.ShmWrite(shmcfg)

}

func (shmclient ShmClient) ShmWriteShutdown() error {

	shmcfg := &ShmConfig{}
	shmcfg.cmd = DFPX_SHM_CMD_SHUTDOWN
    return shmclient.ShmWrite(shmcfg)

}

func (shmclient ShmClient) ShmWrite(cfg *ShmConfig) error {
	shmcfg := &C.dfxp_shm_t{}

    //C.cmark_list_type(lt)
	shmcfg.cmd = C.dfpx_shm_cmd(cfg.cmd)
    
	glog.Infof("Write cmd:%d",shmcfg.cmd)
	ret := C.ShmWrite(shmcfg)
	if ret != 0 {
		return fmt.Errorf("ShmWrite failed")
	}

	return nil
}
