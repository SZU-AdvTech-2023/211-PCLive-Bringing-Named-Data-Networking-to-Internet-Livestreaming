package util

import (
	"strings"
)

func CombineExistedDataNameStrings(textA string, textB string) string {
	if textA == "" {
		return textB
	}
	if textB == "" {
		return textA
	}
	var names1 []string = nil
	var names2 []string = nil
	textA = strings.Replace(textA, "\n", "", -1)
	textB = strings.Replace(textB, "\n", "", -1)
	if textA != "" {
		names1 = strings.Split(textA, ";")
	}
	if textB != "" {
		names2 = strings.Split(textB, ";")
	}
	namesMap := make(map[string]int)
	for i := range names1 {
		namesMap[names1[i]] = 1
	}
	for i := range names2 {
		namesMap[names2[i]] = 1
	}
	var namesArr []string
	for k, _ := range namesMap {
		namesArr = append(namesArr, k)
	}
	res := strings.Join(namesArr, ";")
	return res
}
func LinkTextsWithSemicolon(texts []string) string {
	switch len(texts) {
	case 0:
		return ""
	case 1:
		return texts[0]
	default:
		{
			str := texts[0]
			for i := 1; i < len(texts); i++ {
				str += texts[i]
			}
			return str
		}
	}
}
