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

func configShm(client *shmclient.ShmClient, cmd int, shmcfg *shmclient.ShmConfig) error {

	shmcmd := C.dfxp_shm_cmd(cmd)

	logrus.Infof("configShm.cmd:%d", shmcmd)
	cfg := (*C.dfxp_shm_t)(unsafe.Pointer(&shmcfg.Cfg))

	// client.DumpCfg(shmcfg)

	switch shmcmd {
	case C.DFXP_SHM_CMD_CONFIG_TRAFFIC:
		traffic := (*C.dfxp_traffic_config_t)(unsafe.Pointer(&cfg.cfgTraffic))
		err := configTraffic(traffic)
		if err != nil {
			return err
		}
		client.DumpTraffic(shmcfg)

	case C.DFXP_SHM_CMD_CONFIG_PORTS:
		ports := (*C.dfxp_ports_t)(unsafe.Pointer(&cfg.cfgPorts))
		err := configPorts(ports)
		if err != nil {
			return err
		}
		client.DumpPorts(shmcfg)
	case C.DFXP_SHM_CMD_ADD_IP_GTP:
		tunnels := *(*C.dfxp_shm_ip_gtps_t)(unsafe.Pointer(&cfg.cfgIpGtps))
		tunnels.num = C.int(1)

	case C.DFXP_SHM_CMD_DEL_IP_GTP:
		tunnels := *(*C.dfxp_shm_ip_gtps_t)(unsafe.Pointer(&cfg.cfgIpGtps))
		tunnels.num = C.int(1)

	default:
		fmt.Errorf("Wrong shm config cmd:%d", cmd)
	}
	return nil
}

func configTraffic(traffic *C.dfxp_traffic_config_t) error {
		//required
		traffic.duration = C.int(120) // seconds
		traffic.server = C.bool(false)
		traffic.listen = C.int(5678)
		traffic.listen_num = C.int(1)
		traffic.cps = C.int(50000) // 50k
		traffic.cpu[0] = C.int(1)
		traffic.cpu_num = C.int(1)
		traffic.lport_min = C.int(2020)
		traffic.lport_max = C.int(2030)

		return nil
}

func configPorts (ports *C.dfxp_ports_t) error {
	ports.port_num = C.int(1)
	str :=  (*C.char)	(unsafe.Pointer(&ports.ports[0].pci))
	C.strcpy(str, (*C.char)(C.CString("0000:00:08.0")))

	str =  (*C.char)	(unsafe.Pointer(&ports.ports[0].server_ip))
	C.strcpy(str, (*C.char)(C.CString("192.168.1.224")))
    
	str =  (*C.char)	(unsafe.Pointer(&ports.ports[0].gateway_ip))
	C.strcpy(str, (*C.char)(C.CString("192.168.1.1")))

	str =  (*C.char)	(unsafe.Pointer(&ports.ports[0].local_ip))
	C.strcpy(str, (*C.char)(C.CString("192.168.1.240")))

	return nil
}