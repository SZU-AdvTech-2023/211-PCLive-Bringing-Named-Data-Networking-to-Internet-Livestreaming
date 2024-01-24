package nodeAPI

import (
	"PCLive_project/util"
	"fmt"
	"strings"
	"time"
)

func CheckHavingFile(subName string) bool {
	fileNames := util.GetFileNamesWithCondInDir(subName, util.AbsPath)
	if len(fileNames) == 0 {
		return false
	}
	return true
}
func CheckHavingM3u8(dataName string) bool {
	fileNames := util.GetFileNamesWithCondInDir(dataName+".m3u8", util.AbsPath)
	if len(fileNames) == 0 {
		return false
	}
	return true
}
func LoopHandleReceivedInterest() {
	for {
		interest := ReceiveInterestContinuously()
		havingM3u8 := CheckHavingM3u8(interest)
		if havingM3u8 == true {
			KillNDNPutChunksProcesses(interest)
			PutChunk(interest+".m3u8", interest+".m3u8")
			_, tsFileNames, err := util.ResolveM3u8(interest, util.AbsPath)
			if err != nil {
				fmt.Println(err)
			} else {
				go BatchPutChunks(tsFileNames)
			}
			time.Sleep(time.Second)
		} else {
		}
	}
}
func SendInterestContinuously(interest string) {
	i := 0
	for {
		err := SendInterest(interest)
		if err != nil {
			i++
			if i >= 5 {
				dataNames := GetRecordedExistedDataNames()
				newDataNames := DeleteOneExistedDataName(dataNames, interest)
				SaveExistedDataNames(newDataNames)
				return
			}
			continue
		} else {
			return
		}
	}
}
func SendInterest(interest string) error {
	_, err := util.GetCommandRunningResult("bash", "-c", "echo '"+interest+"' | ndnpoke '/ndn/demo/interest' -w 1000")
	return err
}
func ReceiveInterest() (string, error) {
	EraseInterestInCS()
	result, err := util.GetBashResult("ndnpeek -p '/ndn/demo/interest'")
	result = strings.Replace(result, "\n", "", -1)
	return result, err
}
func EraseInterestInCS() {
	err := util.RunCommand("bash", "-c", "nfdc cs erase /ndn/demo/interest")
	if err != nil {
		fmt.Println(err)
	}
}
func EraseAllFileNamesToBeReceivedInCS() {
	err := util.RunCommand("bash", "-c", "nfdc cs erase /ndn/demo/fileNamesToBeReceived")
	if err != nil {
		fmt.Println(err)
	}
}
func ReceiveInterestContinuously() string {
	for {
		time.Sleep(time.Second)
		interest, err := ReceiveInterest()
		if err != nil {
			continue
		} else {
			fmt.Println("received interest---", interest)
			return interest
		}
	}
}
func HandleReceivedInterest(interestFileName string) {
	if CheckHavingFile(interestFileName) {
		PutChunk(interestFileName, interestFileName)
	} else {
		SendInterest(interestFileName)
	}
}
func CastFileNamesOfInterest(interest string) {
	havingData := CheckHavingM3u8(interest + ".m3u8")
	if havingData {
		fileName, tsLines, err := util.ResolveM3u8(interest, util.AbsPath)
		if err != nil {
			return
		} else {
			fileNames := fileName
			for i := range tsLines {
				fileNames = fileNames + ";" + tsLines[i]
			}
			util.RunCommand("bash", "-c", "echo '"+fileNames+"' | ndnpoke '/ndn/demo/fileNamesToBeReceived/'"+interest+" -w 2000")
		}
	} else {
		fmt.Println("no such data")
	}
}
func ReceiveCastingFileNamesOfInterest(interest string) ([]string, error) {
	EraseAllFileNamesToBeReceivedInCS()
	result, err := util.GetCommandRunningResult("bash", "-c", "ndnpeek -p '/ndn/demo/fileNamesToBeReceived/'"+interest)
	EraseAllFileNamesToBeReceivedInCS()
	if err != nil {
		return nil, err
	} else {
		names := strings.Split(result, ";")
		return names, nil
	}
}
