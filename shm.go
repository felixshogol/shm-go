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
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"bufio"
	"net"

	shmclient "dflux.io/shm-go/shm-client"
	"github.com/sirupsen/logrus"
)

const (
	STR_EMPTY      = "empty"
	STR_QUIT       = "quit"
	STR_GTPENABLE  = "gtpenable"
	STR_GTPDISABLE = "gtpdisable"
)

var shmCfg = shmclient.ShmConfig{}
var gtpEnable bool = false
var protocol uint8 = 17

func main() {
	var err error
	//cmdPtr := flag.Int("cmdFlag", int(0), "shm cmd")
	//flag.Parse()

	logrus := logrus.New()

	logrus.SetReportCaller(true)

	shmclient := shmclient.NewShmClient()
	//cmdname := shmclient.GetShmCmdName(*cmdPtr)
	//logrus.Infof("shmcmd:%d:%s", *cmdPtr, cmdname)

	err = shmclient.InitShm()
	if err != nil {
		logrus.Errorf("shm Init failed.Err:%v", err)
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		text, err := scanConsole(scanner)
		logrus.Info(text)
		if err != nil {
			logrus.Errorf("shm config failed. Err:%v", err)
		}
		if text == STR_EMPTY {
			logrus.Warn("empty line")
			usage()
			continue
		}
		if text == STR_QUIT {
			logrus.Warn("shm-go quit")
			return
		}
		if text == STR_GTPENABLE {
			gtpEnable = true
			continue
		}

		if text == STR_GTPDISABLE {
			gtpEnable = false
			continue
		}
		if text == "tcp" {
			protocol = 6
			continue
		}
		if text == "udp" {
			protocol = 17
			continue
		}

		cmd, err := strconv.Atoi(text)

		if err != nil {
			fmt.Println("Error during conversion")
			continue
		}
		err = shmclient.ShmCmdValidation( /**cmdPtr*/ cmd)
		if err != nil {
			logrus.Errorf("shm validation failed.Err:%v", err)
			continue
		}
		err = configShm(shmclient /*cmdPtr*/, cmd, &shmCfg)
		if err != nil {
			logrus.Errorf("configShm failed.Err:%v", err)
			continue
		}
		err = shmclient.ShmRunCmd( /**cmdPtr*/ cmd, &shmCfg)
		if err != nil {
			logrus.Errorf("Failed to run shm cmd.Err :%v", err)
			continue
		}
	}

}

func configShm(client *shmclient.ShmClient, cmd int, shmcfg *shmclient.ShmConfig) error {

	shmcmd := C.dfxp_shm_cmd(cmd)

	logrus.Infof("configShm.cmd:%d", shmcmd)
	cfg := (*C.dfxp_shm_t)(unsafe.Pointer(&shmcfg.Cfg))

	switch shmcmd {
	case C.DFXP_SHM_CMD_CONFIG_TRAFFIC:
		traffic := (*C.dfxp_traffic_config_t)(unsafe.Pointer(&cfg.cfgTraffic))
		err := configTraffic(traffic,gtpEnable)
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
		tunnels := (*C.dfxp_shm_ip_gtps_t)(unsafe.Pointer(&cfg.cfgIpGtps))
		err := configTunnels(tunnels)
		if err != nil {
			return err
		}

	case C.DFXP_SHM_CMD_DEL_IP_GTP:
		tunnels := (*C.dfxp_shm_ip_gtps_t)(unsafe.Pointer(&cfg.cfgIpGtps))
		deleteTunnels(tunnels)

	case C.DFXP_SHM_CMD_CLEAR_CONFIG:
		logrus.Info("Clear dfxp config")

	case C.DFXP_SHM_CMD_GET_STATS:
		logrus.Info("Get statistics")
	default:
		fmt.Errorf("Wrong shm config cmd:%d", cmd)
	}

	return nil
}

func configTraffic(traffic *C.dfxp_traffic_config_t, gtpenable bool) error {
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
	if gtpenable {
		traffic.gtpu_enable = true
	} else {
		traffic.gtpu_enable = false
	}
	traffic.protocol = C.uint8_t(protocol)
	return nil
}

func configPorts(ports *C.dfxp_ports_t) error {
	ports.port_num = C.int(1)
	str := (*C.char)(unsafe.Pointer(&ports.ports[0].pci))
	C.strcpy(str, (*C.char)(C.CString("0000:00:08.0")))

	str = (*C.char)(unsafe.Pointer(&ports.ports[0].server_ip))
	C.strcpy(str, (*C.char)(C.CString("192.168.1.224")))

	str = (*C.char)(unsafe.Pointer(&ports.ports[0].gateway_ip))
	C.strcpy(str, (*C.char)(C.CString("192.168.1.1")))

	str = (*C.char)(unsafe.Pointer(&ports.ports[0].local_ip))
	C.strcpy(str, (*C.char)(C.CString("192.168.1.240")))

	return nil
}

func configTunnels(tunnels *C.dfxp_shm_ip_gtps_t) error {

	tunnels.num = C.int(2)
	ue1str := "10.0.0.1"
	ue2str := "10.0.0.3"
	upfstr := "106.10.138.240"

	ue1 := (*C.char)(unsafe.Pointer(&tunnels.ip_gtp[0].address))
	C.strcpy(ue1, (*C.char)(C.CString(ue1str)))

	ue2 := (*C.char)(unsafe.Pointer(&tunnels.ip_gtp[1].address))
	C.strcpy(ue2, (*C.char)(C.CString(ue2str)))

	ip := net.ParseIP(upfstr)
	upf, err := shmclient.IPv4ToInt(ip)
	if err != nil {
		return err
	}

	//UE1
	tunnels.ip_gtp[0].tunnel.upf_ipv4 = C.uint32_t(upf)
	ip = net.ParseIP(ue1str)
	ipu, err := shmclient.IPv4ToInt(ip)
	if err != nil {
		return err
	}
	tunnels.ip_gtp[0].tunnel.id = C.uint32_t(0)
	tunnels.ip_gtp[0].tunnel.ue_ipv4 = C.uint32_t(ipu)
	tunnels.ip_gtp[0].tunnel.teid_in = C.uint32_t(10)
	tunnels.ip_gtp[0].tunnel.teid_out = C.uint32_t(1010)

	//UE2
	tunnels.ip_gtp[1].tunnel.id = C.uint32_t(1)
	tunnels.ip_gtp[1].tunnel.upf_ipv4 = C.uint32_t(upf)
	ip = net.ParseIP(ue2str)
	ipu, err = shmclient.IPv4ToInt(ip)
	if err != nil {
		return err
	}
	tunnels.ip_gtp[1].tunnel.ue_ipv4 = C.uint32_t(ipu)
	tunnels.ip_gtp[1].tunnel.teid_in = C.uint32_t(11)
	tunnels.ip_gtp[1].tunnel.teid_out = C.uint32_t(1011)

	return nil
}

func deleteTunnels(tunnels *C.dfxp_shm_ip_gtps_t) error {

	tunnels.num = C.int(1)
	ue1str := "10.0.0.1"
	//ue2str := "10.0.0.3"

	ue1 := (*C.char)(unsafe.Pointer(&tunnels.ip_gtp[0].address))
	C.strcpy(ue1, (*C.char)(C.CString(ue1str)))

	// ue2 := (*C.char)(unsafe.Pointer(&tunnels.ip_gtp[1].address))
	// C.strcpy(ue2, (*C.char)(C.CString(ue2str)))

	return nil
}

func scanConsole(scanner *bufio.Scanner) (string, error) {

	fmt.Print("-> ")
	scanner.Scan()
	// Holds the string that scanned
	text := scanner.Text()
	if len(text) != 0 {
		text = strings.Replace(text, "\n", "", -1)
		fmt.Println(text)
	} else {
		return STR_EMPTY, nil
	}
	// handle error
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
		return "", scanner.Err()
	}
	return text, nil
}

func usage() {
	fmt.Println("uage:")
	fmt.Println("gtpenable - enable gtp")
	fmt.Println("gtpdisable - disable gtp")
	fmt.Println("tcp - protocol tcp")
	fmt.Println("udp - protocol udp")
	fmt.Println("quit- shm-go quit")
	fmt.Println("1 - DFXP_SHM_CMD_CONFIG_TRAFFIC")
	fmt.Println("2 - DFXP_SHM_CMD_CONFIG_PORTS")
	fmt.Println("3 - DFXP_SHM_CMD_START")
	fmt.Println("4 - DFXP_SHM_CMD_STOP")
	fmt.Println("5 - DFXP_SHM_CMD_SHUTDOWN")
	fmt.Println("6 - DFXP_SHM_CMD_ADD_IP_GTP")
	fmt.Println("7 - DFXP_SHM_CMD_DEL_IP_GTP")
	fmt.Println("8 - DFXP_SHM_CMD_GET_STATS")
	fmt.Println("8 - DFXP_SHM_CMD_CLEAR_CONFIG")
}
