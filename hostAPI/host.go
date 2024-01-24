package hostAPI

import (
	"PCLive_project/nodeAPI"
	"PCLive_project/util"
	"fmt"
	"time"
)

func DistributeM3u8Continuously(m3u8Location string, containerName string) {
	go func() {
		for {
			M3u8Distributing(m3u8Location, containerName)
			time.Sleep(time.Second)
		}
	}()
}
func M3u8Distributing(m3u8Location string, containerName string) {
	names := nodeAPI.GetLocalExistedDataNames(util.HostAbsPath)
	if len(names) > 0 {
		for i := range names {
			_, tsFileNames, err := util.ResolveM3u8(names[i], util.HostAbsPath)
			if err != nil {
				fmt.Println(err)
			}
			filePaths := make([]string, 0)
			for i := range tsFileNames {
				filePaths = append(filePaths, m3u8Location+tsFileNames[i])
			}
			filePaths = append(filePaths, m3u8Location+names[i]+".m3u8")
			for i := range filePaths {
				util.RunCommand("bash", "-c", "lxc file push -r "+filePaths[i]+" "+containerName+util.AbsPath)
			}
		}
	}
}
