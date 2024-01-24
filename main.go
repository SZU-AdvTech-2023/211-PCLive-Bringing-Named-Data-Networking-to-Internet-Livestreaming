package main

import (
	"PCLive_project/hostAPI"
	"PCLive_project/nodeAPI"
	"PCLive_project/util"
	"flag"
	"fmt"
	"time"
)

func main() {
	{
		util.RecTime = time.Now()
		util.RecordedSpeedSum = 0
		util.RecordedTimeCount = 0
		var NeedHelp = flag.Bool("h", false, "to get documentation")
		var IsHost = flag.Bool("d", false, "running on the host for distributing created HLS files")
		var InterestName = flag.String("i", "", "interest name of requested data")
		var ToInit = flag.Bool("init", false, "init the environment of containers")
		var ToKillProcesses = flag.Bool("k", false, "kill all processes made by ndn")
		flag.Parse()
		if *NeedHelp == true {
			assistance := "default  ---   serve as a normal node, listening coming interest\n"
			assistance += "-d       ---   running on the host, distribute created HLS files\n"
			assistance += "-i       ---   interest name of requested data\n"
			assistance += "-init    ---   init the environment of containers\n"
			assistance += "-k       ---   kill all processes made by ndn\n"
			fmt.Println(assistance)
		} else if *ToInit == true {
			util.InitContainer(105)
		} else if *ToKillProcesses == true {
			util.KillAllProcessesMadeByNDNInContainers()
		} else {
			if *IsHost == true {
				fmt.Println("the host is distributing files to the container")
				hostAPI.DistributeM3u8Continuously(util.HostAbsPath, "u1")
			} else if *InterestName != "" {
				fmt.Println("waiting for data response")
				nodeAPI.LoopReceiveDataFilesOfSpecificInterest(*InterestName)
			} else {
				fmt.Println("listenning requests")
				go nodeAPI.LoopHandleReceivedInterest()
				go nodeAPI.LoopCastLocalExistedDataNames()
				go nodeAPI.LoopReceiveCastingExistedDataNamesAndRecordThem()
			}
			ch := make(chan int)
			<-ch
		}
	}
}
