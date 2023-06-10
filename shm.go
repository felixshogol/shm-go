package main

import (
	"flag"
	"fmt"
	"os"

	shmclient "dflux.io/shm-go/shm-client"
	"github.com/golang/glog"
)

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "WARNING")
}

func usage() {
	fmt.Print("@@@@@@@@@")
	flag.PrintDefaults()
	os.Exit(2)
}

var shmCfg = shmclient.ShmConfig  {}

func main() {
	var err error
	cmdPtr := flag.Int("cmdFlag", int(0), "shm cmd")

	flag.Parse()

	shmclient := shmclient.NewShmClient()
	cmdname := shmclient.GetShmCmdName(*cmdPtr)
	
	glog.Infof("shmcmd:%d:%s", *cmdPtr, cmdname)

	err = shmclient.ShmCmdValidation(*cmdPtr)
	if err != nil {
	}
	err = shmclient.InitShm()
	if err != nil {
		glog.Errorf("Err:%v", err)
		return
	}

	err = shmclient.ShmRunCmd(*cmdPtr, &shmCfg)
	if err != nil {
		glog.Errorf("Failed to run shm cmd.Err :%v",err)
	}


	// glog.Info("Write start to shm")
	// err = shmclient.ShmWriteSart()
	// if err != nil {
	// 	glog.Errorf("start Err:%v", err)
	// 	return
	// }

	// time.Sleep(1 * time.Second)

	// glog.Info("Write stop to shm")
	// err = shmclient.ShmWriteStop()
	// if err != nil {
	// 	glog.Errorf("stop Err:%v", err)
	// 	return
	// }
	// time.Sleep(1 * time.Second)

	// glog.Info("Write shutdown to shm")
	// err = shmclient.ShmWriteShutdown()
	// if err != nil {
	// 	glog.Errorf("shutdown Err:%v", err)
	// 	return
	// }
	// time.Sleep(1 * time.Second)

}
