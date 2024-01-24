package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func NDNPutChunkAndGetProcessID(ndnName string, fileName string) int {
	proc := exec.Command("ndnputchunks", "-s", "1000", "/ndn/demo/data/"+ndnName, "<", AbsPath+fileName, "&")
	proc.Start()
	Pid := proc.Process.Pid
	return Pid
}
func RunCommand(instruction string, arg ...string) error {
	cmd := exec.Command(instruction, arg...)
	e := cmd.Run()
	if e != nil {
		return e
	}
	return nil
}
func GetBashResult(arg string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", arg)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	} else {
		return string(out), nil
	}
}
func GetCommandRunningResult(instruction string, arg ...string) (string, error) {
	res := ""
	cmd := exec.Command(instruction, arg...)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		return "", err
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		res += fmt.Sprint(string(tmp))
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		return "", err
	}
	return res, nil
}
func RecordDurationDemo1(t time.Time) {
	fmt.Println("execution time --- ", time.Since(t).Nanoseconds()/(1000*1000))
}
func RecordDurationDemo2() {
	RunCommand("bash", "-c", "ndnputchunks -s 1000 /localhost/demo/gpl3 < /usr/share/common-licenses/GPL-3 &")
	start := time.Now()
	time.Sleep(time.Second * 2)
	dur := time.Since(start)
	durFloat := float64(dur)
	fmt.Println(durFloat / 1000000000)
}
func WriteCommandRunningResultToFile(filePath string, instruction string, arg ...string) error {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buffer := bufio.NewWriter(f)
	cmd := exec.Command(instruction, arg...)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		if strings.Contains(string(tmp), "reason: NoRoute") {
			fmt.Println("no such chunks provided")
			os.Remove(filePath)
			break
		}
		_, wErr := buffer.WriteString(string(tmp))
		if wErr != nil {
			log.Fatal(err)
		}
		if err := buffer.Flush(); err != nil {
			log.Fatal(err)
		}
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}
