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
	"unsafe"

	"github.com/sirupsen/logrus"
)

type ShmConfigTraffic struct {
}

type ShmConfigPorts struct {
}

type ShmConfigTunnels struct {
}

type ShmConfig struct {
	Cmd int
	Cfg C.dfxp_shm_t
}

type ShmClient struct {
}

func NewShmClient() *ShmClient {
	return &ShmClient{}
}

func (shmClient *ShmClient) DumpCfg(shmcfg *ShmConfig) {
	logrus.Infof("###  DumpCfg:")
	logrus.Infof("ShmCfg clen:%d", C.ShmSizeofCfg())
	// cfg := (*C.dfxp_shm_t)(unsafe.Pointer(&shmcfg.Cfg))
	// traffic := (*C.dfxp_traffic_config_t) (unsafe.Pointer(&cfg.value[0]))
	shmClient.DumpTraffic(shmcfg)

}

func (shmClient *ShmClient) DumpTraffic(shmcfg *ShmConfig) {
	cfg := (*C.dfxp_shm_t)(unsafe.Pointer(&shmcfg.Cfg))
	traffic := (*C.dfxp_traffic_config_t)(unsafe.Pointer(&cfg.cfgTraffic))

	logrus.Infof("### Traffic required:")
	logrus.Infof("Traffic len:%d clen:%d", unsafe.Sizeof(*traffic), C.ShmSizeofTraffic())
	logrus.Infof("server:%t", bool(traffic.server))
	logrus.Infof("duration:%dsec", int(traffic.duration))
	logrus.Infof("cps:%d", int(traffic.cps))
	logrus.Infof("listen:%d", int(traffic.listen))
	logrus.Infof("listen num:%d", int(traffic.listen_num))
	logrus.Infof("cpu_num:%d", int(traffic.cpu_num))
	for i := 0; i < int(traffic.cpu_num); i++ {
		logrus.Info("cpu:", C.int(traffic.cpu[i]))
	}
	logrus.Infof("cc:%d", int(traffic.cc))

}

func (shmClient *ShmClient) DumpPorts(shmcfg *ShmConfig) {

	cfg := (*C.dfxp_shm_t)(unsafe.Pointer(&shmcfg.Cfg))
	ports := (*C.dfxp_ports_t)(unsafe.Pointer(&cfg.cfgPorts))

	logrus.Infof("### Ports :")
	logrus.Infof("Ports len:%d clen:%d", unsafe.Sizeof(*ports), C.ShmSizeofPorts())
	logrus.Infof("port_num:%d", int(ports.port_num))
	str := (*C.char)(unsafe.Pointer(&ports.ports[0].pci))
	logrus.Infof("pci:%s", C.GoString(str))
	str = (*C.char)(unsafe.Pointer(&ports.ports[0].server_ip))
	logrus.Infof("server_ip:%s", C.GoString(str))
	str = (*C.char)(unsafe.Pointer(&ports.ports[0].gateway_ip))
	logrus.Infof("gateway_ip:%s", C.GoString(str))
	str = (*C.char)(unsafe.Pointer(&ports.ports[0].local_ip))
	logrus.Infof("local_ip:%s", C.GoString(str))

}

func (shmClient *ShmClient) DumpPorts1(cfg *C.dfxp_shm_t) {

	ports := (*C.dfxp_ports_t)(unsafe.Pointer(&cfg.cfgPorts))

	logrus.Infof("### DumpPorts1 Ports :")
	logrus.Infof("Ports len:%d clen:%d", unsafe.Sizeof(*ports), C.ShmSizeofPorts())
	logrus.Infof("port_num:%d", int(ports.port_num))
	str := (*C.char)(unsafe.Pointer(&ports.ports[0].pci))
	logrus.Infof("pci:%s", C.GoString(str))
	str = (*C.char)(unsafe.Pointer(&ports.ports[0].server_ip))
	logrus.Infof("server_ip:%s", C.GoString(str))
	str = (*C.char)(unsafe.Pointer(&ports.ports[0].gateway_ip))
	logrus.Infof("gateway_ip:%s", C.GoString(str))
	str = (*C.char)(unsafe.Pointer(&ports.ports[0].local_ip))
	logrus.Infof("local_ip:%s", C.GoString(str))

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
		cfg.Cmd = int(C.DFXP_SHM_CMD_CONFIG_TRAFFIC)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_CONFIG_PORTS:
		cfg.Cmd = int(C.DFXP_SHM_CMD_CONFIG_PORTS)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_START:
		cfg.Cmd = int(C.DFXP_SHM_CMD_START)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_STOP:
		cfg.Cmd = int(C.DFXP_SHM_CMD_STOP)
		return shmClient.ShmWrite(cfg)

	case C.DFXP_SHM_CMD_SHUTDOWN:
		cfg.Cmd = int(C.DFXP_SHM_CMD_SHUTDOWN)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_ADD_IP_GTP:
		cfg.Cmd = int(C.DFXP_SHM_CMD_ADD_IP_GTP)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_DEL_IP_GTP:
		cfg.Cmd = int(C.DFXP_SHM_CMD_DEL_IP_GTP)
		return shmClient.ShmWrite(cfg)
	case C.DFXP_SHM_CMD_GET_STATS:

	default:
		fmt.Errorf("Wrong shm cmd:%d", cmd)
	}
	return nil
}

func (shmclient *ShmClient) InitShm() error {

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

	// shmcfg := &C.dfxp_shm_t{}
	//cdata := C.GoBytes(unsafe.Pointer(shmcfg), C.sizeof_dfxp_shm_t)
	shmcfg :=(*C.dfxp_shm_t)(unsafe.Pointer(&cfg.Cfg))

	// shmcfg := &cfg.Cfg
	shmcfg.cmd = C.dfxp_shm_cmd(cfg.Cmd)
	shmcfg.status = C.DFXP_SHM_STATUS_WRITTEN_BY_CLIENT

	logrus.Infof("Write cmd:%d", shmcfg.cmd)
	if  shmcfg.cmd  == C.DFXP_SHM_CMD_CONFIG_PORTS{
		shmclient.DumpPorts1(shmcfg)
		
	}
	ret := C.ShmWrite(shmcfg)
	if ret != 0 {
		return fmt.Errorf("ShmWrite failed")
	}

	return nil
}
