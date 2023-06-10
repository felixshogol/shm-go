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

type ShmConfig struct {
	cmd int
}

type ShmClient struct {
}

func NewShmClient() *ShmClient {
	return &ShmClient{}
}

func (shmclient *ShmClient) GetShmCmdName(cmd int) string {
	cmdstr := C.ShmGetCmdName(C.dfxp_shm_cmd(cmd))
	return C.GoString(cmdstr)

}

func (shmClient *ShmClient) ShmCmdValidation(cmd int) error {

	switch C.dfxp_shm_cmd(cmd) {
	case C.DFXP_SHM_CMD_NONE:
	case C.DFXP_SHM_CMD_CONFIG_TRAFFIC:
	case C.DFXP_SHM_CMD_CONFIG_PORTS:
	case C.DFXP_SHM_CMD_START:
	case C.DFXP_SHM_CMD_STOP:
	case C.DFXP_SHM_CMD_SHUTDOWN:
	case C.DFXP_SHM_CMD_ADD_IP_GTP:
	case C.DFXP_SHM_CMD_DEL_IP_GTP:
	case C.DFXP_SHM_CMD_GET_STATS:

	default:
		fmt.Errorf("Wrong shm cmd:%d", cmd)
	}
	return nil
}

func (shmClient *ShmClient) ShmRunCmd(cmd int, cfg *ShmConfig) error {
	switch C.dfxp_shm_cmd(cmd) {
	case C.DFXP_SHM_CMD_NONE:
	case C.DFXP_SHM_CMD_CONFIG_TRAFFIC:
		cfg.cmd = int(C.DFXP_SHM_CMD_CONFIG_TRAFFIC)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_CONFIG_PORTS:
		cfg.cmd = int(C.DFXP_SHM_CMD_CONFIG_PORTS)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_START:
		cfg.cmd = int(C.DFXP_SHM_CMD_START)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_STOP:
		cfg.cmd = int(C.DFXP_SHM_CMD_STOP)
		return shmClient.ShmWrite(cfg)

	case C.DFXP_SHM_CMD_SHUTDOWN:
		cfg.cmd = int(C.DFXP_SHM_CMD_SHUTDOWN)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_ADD_IP_GTP:
	case C.DFXP_SHM_CMD_DEL_IP_GTP:
	case C.DFXP_SHM_CMD_GET_STATS:

	default:
		fmt.Errorf("Wrong shm cmd:%d", cmd)
	}
	return nil
}

func (shmclient *ShmClient) InitShm() error {
	glog.Info("InitShm")

	shmn := C.CString("/dfxp-shm")
	ret := C.ShmInit(shmn, 2, 0)

	if ret != 0 {
		return fmt.Errorf("ShmInit failed")
	}

	return nil
}

func (shmclient *ShmClient) ShmWriteSart(cfg *ShmConfig) error {
	return shmclient.ShmWrite(cfg)
}

func (shmclient *ShmClient) ShmWriteStop(cfg *ShmConfig) error {
	return shmclient.ShmWrite(cfg)
}

func (shmclient *ShmClient) ShmWriteShutdown(cfg *ShmConfig) error {
	return shmclient.ShmWrite(cfg)
}

func (shmclient *ShmClient) ShmWriteConfig(cfg *ShmConfig) error {
	return shmclient.ShmWrite(cfg)
}

func (shmclient *ShmClient) ShmWrite(cfg *ShmConfig) error {

	shmcfg := &C.dfxp_shm_t{}
	//cdata := C.GoBytes(unsafe.Pointer(shmcfg), C.sizeof_dfxp_shm_t)

	shmcfg.cmd = C.dfxp_shm_cmd(cfg.cmd)

	glog.Infof("Write cmd:%d", shmcfg.cmd)
	ret := C.ShmWrite(shmcfg)
	if ret != 0 {
		return fmt.Errorf("ShmWrite failed")
	}

	return nil
}
