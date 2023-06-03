package main

import (
	"flag"
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
	glog.Info("Start shm client")

	shmclient := shmclient.NewShmClient()

	glog.Info("Init shm")

	err := shmclient.InitShm()
	if err != nil {
		glog.Errorf("Err:%v",err)
		return
	}

	glog.Info("Write to shm")

}
