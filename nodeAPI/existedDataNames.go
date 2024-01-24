package nodeAPI

import (
	"PCLive_project/util"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func LoopReceiveCastingExistedDataNamesAndRecordThem() {
	for {
		time.Sleep(time.Millisecond * 300)
		ReceiveCastingExistedDataNamesAndRecordThem()
	}
}
func LoopCastLocalExistedDataNames() {
	for {
		CastLocalExistedDataNames()
	}
}
func CastLocalExistedDataNames() {
	dataNames := GetLocalExistedDataNames(util.AbsPath)
	LocalExistedDataNameString := util.LinkTextsWithSemicolon(dataNames)
	time.Sleep(time.Second)
	util.RunCommand("bash", "-c", "echo '"+LocalExistedDataNameString+"' | ndnpoke '/ndn/demo/existedData' -w 10000")
}
func DeleteOneExistedDataName(dataNames []string, dataNameToBeDeleted string) string {
	for i := range dataNames {
		if dataNameToBeDeleted == dataNames[i] {
			dataNames = append(dataNames[:i], dataNames[i+1:]...)
			break
		}
	}
	dataNameString := strings.Join(dataNames, ";")
	return dataNameString
}
func GetRecordedExistedDataNames() []string {
	ExistedDataNameString, _ := util.ReadFile(util.AbsPath + "existedDataNames")
	ExistedDataNameString = strings.Replace(ExistedDataNameString, "\n", "", -1)
	var ExistedDataNames []string
	if ExistedDataNameString == "" {
		ExistedDataNames = nil
	} else {
		ExistedDataNames = strings.Split(ExistedDataNameString, ";")
	}
	return ExistedDataNames
}
func GetLocalExistedDataNames(location string) []string {
	fileNames, err := util.GetFilePathsWithDirAndSubName(location, ".m3u8")
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		for i := range fileNames {
			fileNames[i] = fileNames[i][strings.LastIndex(fileNames[i], "/")+1 : strings.Index(fileNames[i], ".")]
		}
	}
	return fileNames
}
func SaveExistedDataNames(dataNameString string) {
	filePath := util.AbsPath + "existedDataNames"
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(dataNameString)
	write.Flush()
}
func EraseCachedExistedDataNamesInCS() {
	err := util.RunCommand("bash", "-c", "nfdc cs erase /ndn/demo/existedData")
	if err != nil {
		fmt.Println(err)
	}
}
func ReceiveCastingExistedDataNamesAndRecordThem() {
	LocalSavedExistedDataString, err := util.GetBashResult("cat " + util.AbsPath + "existedDataNames")
	if err != nil {
		util.RunCommand("mkdir", "-p", util.AbsPath)
		util.RunCommand("touch", util.AbsPath+"existedDataNames")
		LocalSavedExistedDataString = ""
	}
	EraseCachedExistedDataNamesInCS()
	CastingExistedDataString, _ := util.GetBashResult("ndnpeek -p '/ndn/demo/existedData'")
	NewExistedDataNameString := util.CombineExistedDataNameStrings(LocalSavedExistedDataString, CastingExistedDataString)
	NewExistedDataNameString = strings.Replace(NewExistedDataNameString, "\n", "", -1)
	SaveExistedDataNames(NewExistedDataNameString)
}
func AddOneLocalExistedDataName(newDataName string) {
	dataNames := GetRecordedExistedDataNames()
	dataNames = append(dataNames, newDataName)
	dataNameString := strings.Join(dataNames, ";")
	SaveExistedDataNames(dataNameString)
}
