package nodeAPI

import (
	"PCLive_project/util"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoopReceiveDataFilesOfSpecificInterest(interestName string) {
	fmt.Println("try receiving", interestName)
	CastInterestAndReceiveDataFiles(interestName)
}
func BatchCatchChunks(fileNames []string) {
	for i := range fileNames {
		if util.CheckFileValidity(util.AbsPath+fileNames[i]) == true {
			continue
		} else {
			go func(index int) {
				start := time.Now()
				count := 0
				for {
					if util.CheckFileValidity(util.AbsPath+fileNames[index]) == false {
						_, err := util.GetBashResult("ndncatchunks /ndn/demo/data/" + fileNames[index] + " > " + util.AbsPath + fileNames[index])
						if err != nil {
							os.Remove(util.AbsPath + fileNames[index])
							time.Sleep(time.Millisecond * 100)
							if count >= 5 {
								break
							}
						} else {
						}
					} else {
						sizeString, _ := util.GetFileSize(util.AbsPath + fileNames[index])
						sizeFloat, _ := strconv.ParseFloat(sizeString, 64)
						durationTime := time.Since(start)
						durationFloat := float64(durationTime)
						fmt.Println(fileNames[index], "speed         ---", sizeFloat/(durationFloat/1000000000)*8, "bit/s")
						fmt.Println()
						dur := time.Now().Sub(util.RecTime)
						if dur.Minutes() >= 10 {
							filePath := "/root/record.txt"
							file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
							if err != nil {
								fmt.Println("文件打开失败", err)
							}
							defer file.Close()
							avgSpeed := util.RecordedSpeedSum / util.RecordedTimeCount
							write := bufio.NewWriter(file)
							sprintf := fmt.Sprintln(avgSpeed, "---bit/s")
							write.WriteString(sprintf)
							write.Flush()
							util.RecordedSpeedSum = 0
							util.RecordedTimeCount = 0
							util.RecTime = time.Now()
						} else {
							util.RecordedSpeedSum += sizeFloat / (durationFloat / 1000000000) * 8
							util.RecordedTimeCount += 1
						}
						return
					}
					count++
				}
			}(i)
		}
	}
}
func CastInterestAndReceiveDataFiles(interest string) {
	for {
		SendInterestContinuously(interest)
		EraseInterestInCS()
		tsFiles := CatM3u8ChunkAndResolve(interest)
		BatchCatchChunks(tsFiles)
	}
}
func RecoverM3u8(interest string, content string) {
	filePath := util.AbsPath + interest + ".m3u8"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(content)
	write.Flush()
}
func CatM3u8ChunkAndResolve(interest string) []string {
	originalContent, err := util.ReadFile(util.AbsPath + interest + ".m3u8")
	for {
		util.RunCommand("bash", "-c", "ndncatchunks /ndn/demo/data/"+interest+".m3u8"+" > "+util.AbsPath+interest+".m3u8")
		newContent, _ := util.ReadFile(util.AbsPath + interest + ".m3u8")
		if newContent == "" {
			RecoverM3u8(interest, originalContent)
			time.Sleep(time.Millisecond * 500)
			continue
		} else {
			break
		}
	}
	_, tsLines, err := util.ResolveM3u8(interest, util.AbsPath)
	if err != nil {
		return nil
	} else {
		return tsLines
	}
}
func BatchPutChunks(fileNames []string) {
	for i := range fileNames {
		go PutChunk(fileNames[i], fileNames[i])
	}
}
func GetNDNPutChunksPIDs(interest string) []string {
	result, err := util.GetCommandRunningResult("bash", "-c", "ps -ef |grep 'ndnputchunks' | grep '/ndn/demo/data/"+interest+"' | grep -v 'grep' | awk '{print $2}'")
	if err != nil {
		return nil
	} else {
		split := strings.Split(result, "\n")
		return split
	}
}
func PutChunk(ndnName string, fileName string) {
	util.RunCommand("bash", "-c", "ndnputchunks -s 1000 /ndn/demo/data/"+ndnName+" < "+util.AbsPath+fileName+" &")
}
func KillNDNPutChunksProcesses(interest string) {
	pids := GetNDNPutChunksPIDs(interest)
	if pids != nil {
		for i := range pids {
			util.RunCommand("bash", "-c", "kill "+pids[i])
		}
	}
}
func PutM3u8ChunkIfExisted(interest string) {
	havingM3u8 := CheckHavingM3u8(interest)
	if havingM3u8 == true {
		PutChunk(interest+".m3u8", interest+".m3u8")
	} else {
	}
}
func CatchChunk(ndnName string, fileName string) {
	err := util.RunCommand("bash", "-c", "ndncatchunks /ndn/demo/"+ndnName+" > "+util.AbsPath+fileName)
	if err != nil {
		fmt.Println("fail to catch chunks")
	}
}
