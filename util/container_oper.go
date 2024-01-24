package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func RestartAllContainers() {
	containerNames := GetAllContainerNames()
	for i := range containerNames {
		fmt.Println(containerNames[i] + "---stopping")
		RunCommand("bash", "-c", "lxc stop "+containerNames[i])
		fmt.Println(containerNames[i] + "---starting")
		RunCommand("bash", "-c", " lxc start "+containerNames[i])
	}
	time.Sleep(time.Second)
	result, _ := GetCommandRunningResult("bash", "-c", "lxc list")
	fmt.Println(result)
}
func GetAllContainerNames() []string {
	result, _ := GetCommandRunningResult("bash", "-c", "lxc list")
	containerNames := make([]string, 0)
	i := 1
	for {
		containerName := "u" + strconv.Itoa(i)
		if strings.Contains(result, containerName) {
			containerNames = append(containerNames, containerName)
			i++
		} else {
			break
		}
	}
	return containerNames
}
func KillAllProcessesMadeByNDNInContainers() {
	RunCommand("bash", "-c", "ps -ef |grep ndnput | grep -v grep | awk '{print $2}'|xargs kill")
}
