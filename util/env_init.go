package util

import (
	"fmt"
	"strconv"
	"strings"
)

func InitContainer(beginIP int) {
	err := RunCommand("bash", "-c", "rm -rf "+HostAbsPath+"*")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("host hls_file emptied")
	}
	RestartAllContainers()
	names := GetAllContainerNames()
	for i := 0; i < len(names); i++ {
		result, _ := GetBashResult("lxc exec " + names[i] + " nfdc cs config serve off")
		fmt.Println(names[i] + "  ---  " + result)
	}
	CreateHlsDirForAllContainer()
	fmt.Println("successfully set the environment of containers!")
}
func CreateHlsDirForAllContainer() {
	i := 1
	for {
		exist, err := GetCommandRunningResult("bash", "-c", "lxc exec u"+strconv.Itoa(i)+" ls | grep hls_file | wc -l")
		if err != nil {
			fmt.Println(err)
			break
		} else if strings.Contains(exist, "Error: Instance not found") {
			break
		} else if strings.Index(exist, "0") == 0 {
			_, err := GetCommandRunningResult("bash", "-c", "lxc exec u"+strconv.Itoa(i)+" mkdir /root/hls_file")
			if err != nil {
				if strings.Contains(err.Error(), "exit status 1") {
					break
				}
			} else {
				fmt.Print("u"+strconv.Itoa(i), "---successfully create hls_file directory\n\n")
			}
		} else if strings.Index(exist, "1") == 0 {
			fmt.Print("u"+strconv.Itoa(i), "---successfully create hls_file directory\n\n")
		} else {
			fmt.Println(exist)
		}
		i++
	}
}
