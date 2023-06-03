package main

import (
	"flag"
	"time"
	//"fmt"
	//log "github.com/sirupsen/logrus"
	"os"

	shmclient "dflux.io/shm-go/shm-client"
	"github.com/golang/glog"
)

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "WARNING")
	flag.Set("v", "2")
	// This is wa
	flag.Parse()
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var err error 
	glog.Info("Start shm client")

	shmclient := shmclient.NewShmClient()

	glog.Info("Init shm")

	err = shmclient.InitShm()
	if err != nil {
		glog.Errorf("Err:%v",err)
		return
	}

	glog.Info("Write start to shm")
	err = shmclient.ShmWriteSart()
	if err != nil {
		glog.Errorf("start Err:%v",err)
		return
	}

	time.Sleep(1 *time.Second)

	glog.Info("Write stop to shm")
	err = shmclient.ShmWriteStop()
	if err != nil {
		glog.Errorf("stop Err:%v",err)
		return
	}
	time.Sleep(1 *time.Second)
	

	glog.Info("Write shutdown to shm")
	err = shmclient.ShmWriteShutdown()
	if err != nil {
		glog.Errorf("shutdown Err:%v",err)
		return
	}
	time.Sleep(1 *time.Second)

}
