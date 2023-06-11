package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdbool.h>
#include "shm-client/shmh/dfxp_shm_common.h"
#include "shm-client/shmh/dfxp_shm_client.h"

*/
import "C"
import (
	"flag"
	"fmt"
	"unsafe"

	shmclient "dflux.io/shm-go/shm-client"
	"github.com/sirupsen/logrus"
)

var shmCfg = shmclient.ShmConfig{}


func main() {
	var err error
	cmdPtr := flag.Int("cmdFlag", int(0), "shm cmd")
	flag.Parse()

	logrus := logrus.New()

	logrus.SetReportCaller(true)

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

	err = configShm(shmclient, *cmdPtr, &shmCfg)
	if err != nil {
		logrus.Errorf("shm config failed. Err:%v", err)
		return
	}
	err = shmclient.ShmRunCmd(*cmdPtr, &shmCfg)
	if err != nil {
		logrus.Errorf("Failed to run shm cmd.Err :%v", err)
	}

}

func configShm(client *shmclient.ShmClient, cmd int /*cfg *C.dfxp_shm_t*/, shmcfg *shmclient.ShmConfig) error {

	shmcmd := C.dfxp_shm_cmd(cmd)

	logrus.Infof("configShm.cmd:%d", shmcmd)
	cfg := (*C.dfxp_shm_t)(unsafe.Pointer(&shmcfg.Cfg))

	client.DumpCfg(shmcfg)

	switch shmcmd {
	case C.DFXP_SHM_CMD_CONFIG_TRAFFIC:
		//required
		traffic := (*C.dfxp_traffic_config_t)(unsafe.Pointer(&cfg.value[0]))
		traffic.duration = C.int(120) // seconds
		traffic.server = C.bool(false)
		traffic.listen = C.int(5678)
		traffic.listen_num = C.int(1)
		traffic.cps = C.int(50000) // 50k
		traffic.cpu[0] = C.int(1)
		traffic.cpu[1] = C.int(20)
		traffic.cpu_num = C.int(2)
		client.DumpTraffic(shmcfg)

	case C.DFXP_SHM_CMD_CONFIG_PORTS:
		ports := *(*C.dfxp_ports_t)(unsafe.Pointer(&cfg.value[1]))
		ports.port_num = C.int(1)
		client.DumpPorts(shmcfg)

	case C.DFXP_SHM_CMD_ADD_IP_GTP:
		tunnels := *(*C.dfxp_shm_ip_gtps_t)(unsafe.Pointer(&cfg.value[2]))
		tunnels.num = C.int(1)

	case C.DFXP_SHM_CMD_DEL_IP_GTP:
		tunnels := *(*C.dfxp_shm_ip_gtps_t)(unsafe.Pointer(&cfg.value[2]))
		tunnels.num = C.int(1)

	default:
		fmt.Errorf("Wrong shm config cmd:%d", cmd)
	}
	return nil
}
