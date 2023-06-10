package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "shm-client/shmh/dfxp_shm_common.h"
#include "shm-client/shmh/dfxp_shm_client.h"

*/
import "C"
import (
	"flag"
	"fmt"

	shmclient "dflux.io/shm-go/shm-client"
	"github.com/sirupsen/logrus"
)

var shmCfg = shmclient.ShmConfig{}

func main() {
	var err error
	cmdPtr := flag.Int("cmdFlag", int(0), "shm cmd")
	flag.Parse()

	shmclient := shmclient.NewShmClient()
	cmdname := shmclient.GetShmCmdName(*cmdPtr)

	logrus.Infof("shmcmd:%d:%s", *cmdPtr, cmdname)

	err = shmclient.ShmCmdValidation(*cmdPtr)
	if err != nil {
		logrus.Errorf("shm validation failed.Err:%v", err)
		return
	}
	err = shmclient.InitShm()
	if err != nil {
		logrus.Errorf("shm Init failed.Err:%v", err)
		return
	}

	err = configShm(*cmdPtr, &shmCfg)
	if err != nil {
		logrus.Errorf("shm config failed. Err:%v", err)
		return
	}

	err = shmclient.ShmRunCmd(*cmdPtr, &shmCfg)
	if err != nil {
		logrus.Errorf("Failed to run shm cmd.Err :%v", err)
	}

}

func configShm(cmd int, cfg *shmclient.ShmConfig) error {

	shmcmd := C.dfxp_shm_cmd(cmd)

	logrus.Infof("configShm.cmd:%d", shmcmd)
	switch shmcmd {
	case C.DFXP_SHM_CMD_CONFIG_TRAFFIC:
		traffic := &shmclient.ShmConfigTraffic{}
		cfg.CfgTraffic = traffic
	case C.DFXP_SHM_CMD_CONFIG_PORTS:
		ports := &shmclient.ShmConfigPorts{}
		cfg.CfgPorts = ports
	case C.DFXP_SHM_CMD_ADD_IP_GTP:
		tunnels := &shmclient.ShmConfigTunnels{}
		cfg.CfgTunnels = tunnels
	case C.DFXP_SHM_CMD_DEL_IP_GTP:
		tunnels := &shmclient.ShmConfigTunnels{}
		cfg.CfgTunnels = tunnels

	default:
		fmt.Errorf("Wrong shm config cmd:%d", cmd)
	}
	return nil
}
