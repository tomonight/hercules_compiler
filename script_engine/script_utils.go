package script_engine

import (
	"encoding/json"
	"errors"
)

const (
	STR_PREFIX = "c4heli2019"
)

//getNextLineNo get next execute line number
func getNextLine(lineNo int, allTokens []scriptLineTokens) (nextNo int, nextLine scriptLineTokens) {
	if lineNo < 0 {
		return 0, allTokens[0]
	}
	index := 0
	for i := 0; i < len(allTokens); i++ {
		if allTokens[i].lineNo == lineNo {
			index = i
			break
		}
	}
	if index == len(allTokens)-1 {
		return -1, nextLine
	}
	return allTokens[index+1].lineNo, allTokens[index+1]
}

//getIfExpressTrueOrFalse check the express true or false
func getIfExpressTrueOrFalse(expressLeft, conditon, expressRight string) (result bool, err error) {
	switch conditon {
	case "==":
		return expressLeft == expressRight, nil
	case "!=":
		return expressLeft != expressRight, nil
	case ">":
		return expressLeft > expressRight, nil
	case "<":
		return expressLeft < expressRight, nil
	case ">=":
		return expressLeft >= expressRight, nil
	case "<=":
		return expressLeft <= expressRight, nil
	default:
		return false, errors.New("invalid express condition")
	}

}

//remove the suffix for list
func removeListSuffix(list []int) (returnList []int) {
	for i := 0; i < len(list)-1; i++ {
		returnList = append(returnList, list[i])
	}
	return returnList
}

//mapToString convert map to string type
func mapToString(data map[string]string) (str string, err error) {
	mjson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	mString := string(mjson)
	reString := STR_PREFIX + mString
	return reString, nil
}
