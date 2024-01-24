package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func CheckFileValidity(filePath string) bool {
	size, _ := GetFileSize(filePath)
	if size == "" || size == "0" {
		return false
	} else {
		return true
	}
}
func GetFileSize(filePath string) (string, error) {
	result, err := GetBashResult("stat " + filePath + " -t | awk '{print $2}'")
	result = strings.Replace(result, "\n", "", -1)
	return result, err
}
func ResolveM3u8(fileName string, location string) (string, []string, error) {
	filePath := location + fileName + ".m3u8"
	lines, err := ReadLines(filePath)
	tsLines := make([]string, 0)
	for i := range lines {
		if strings.Contains(lines[i], ".ts") {
			tsLines = append(tsLines, lines[i])
		}
	}
	if err != nil {
		return fileName, nil, err
	}
	return fileName, tsLines, nil
}
func GetFilePathsWithDirAndSubName(dirPath string, subName string) ([]string, error) {
	fileNames := make([]string, 0)
	rd, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return fileNames, err
	}
	for _, fi := range rd {
		if !fi.IsDir() && strings.Contains(fi.Name(), subName) {
			fullName := dirPath + "/" + fi.Name()
			fileNames = append(fileNames, fullName)
		}
	}
	return fileNames, nil
}
func GetFileNamesWithCondInDir(subName string, dirname string) []string {
	files, _ := ioutil.ReadDir(dirname)
	fileNames := make([]string, 0)
	for _, f := range files {
		if strings.Contains(f.Name(), subName) && !strings.Contains(f.Name(), ".ts.tmp") {
			fileNames = append(fileNames, f.Name())
		}
	}
	return fileNames
}
func GetAllFileNamesInDir(dirname string) []string {
	files, _ := ioutil.ReadDir(dirname)
	fileNames := make([]string, 0)
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	return fileNames
}
func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func ReadFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), err
}
func GetMD5sum(filePath string) string {
	tmp, _ := GetCommandRunningResult("md5sum", filePath)
	md5sum := (strings.Split(tmp, "  "))[0]
	return md5sum
}
func EmptyHlsFileDir() {
	RunCommand("bash", "-c", "rm -rf /root/hls_file/**")
	RunCommand("bash", "-c", "rm -rf /home/aaa/hls_file/**")
}
