package script_engine

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

var (
	ifExpress      = []string{"if", "endif"}
	forExpress     = []string{"for", "endfor"}
	finallyExpress = []string{"finally", "endfinally"}
)

func isCommentLine(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return true
	}
	if line[0] == '#' {
		return true
	}
	return false
}
func tokensToParams(tc *scriptContext, threadContext *scriptThreadContext, tokens []string, startIdx int) (map[string]string, error) {
	//log.Printf("tokens=%q\n", tokens)

	result := make(map[string]string)
	for i := startIdx; i < len(tokens); i++ {
		pairs := strings.SplitN(tokens[i], "=", 2)
		if len(pairs) != 2 {
			//log.Printf("tokens=%q\n", tokens)
			return nil, fmt.Errorf("parse token '%s' to statement params error", tokens[i])
		}

		result[pairs[0]] = replaceToken(tc, threadContext, pairs[1])
	}
	return result, nil
}

//将token中的占位符进行替换
func replaceToken(tc *scriptContext, threadContext *scriptThreadContext, token string) string {

	re, err := regexp.Compile(`\$\{[\w]+\}`)
	//正常不应该编译出错，所以日志也不会打印
	if err != nil {
		log.Printf("engine replace token for variables error: %s\n", err.Error())
		return token
	}
	vars := re.FindAllString(token, -1)
	//数据替换区分大小写

	for _, varTag := range vars {
		varName := varTag[2 : len(varTag)-1]
		if len(varName) == 0 {
			continue
		}

		//先替换变量
		value, exists := threadContext.vars[varName]
		if exists {
			token = strings.Replace(token, varTag, value, -1)
			continue
		}
		value = tc.params[varName]
		token = strings.Replace(token, varTag, value, -1)
	}

	if threadContext.resultData == nil {
		return token
	}
	re, err = regexp.Compile(`\$\{\{[\w]+\}\}`)
	//正常不应该编译出错，所以日志也不会打印
	if err != nil {
		log.Printf("engine replace token for result data error: %s\n", err.Error())
		return token
	}
	vars = re.FindAllString(token, -1)

	//数据替换区分大小写
	for _, varTag := range vars {
		varName := varTag[3 : len(varTag)-2]
		if len(varName) == 0 {
			continue
		}
		//先替换变量
		value, _ := threadContext.resultData[varName]
		// if exists {
		token = strings.Replace(token, varTag, value, -1)
		// }
	}
	return token
}

//将token中的占位符进行替换
func replaceAllToken(tc *scriptContext, threadContext *scriptThreadContext, tokens []string) (returnTokens []string) {
	for _, v := range tokens {
		token := replaceToken(tc, threadContext, v)
		token = lengthFunc(tc, threadContext, token)
		returnTokens = append(returnTokens, token)
	}
	return returnTokens
}

//针对list类型的数据，做length操作
func lengthFunc(tc *scriptContext, threadContext *scriptThreadContext, token string) string {
	re, err := regexp.Compile(`length\([\w]+\)`)
	//正常不应该编译出错，所以日志也不会打印
	if err != nil {
		log.Printf("engine replace token for variables error: %s\n", err.Error())
		return token
	}
	vars := re.FindString(token)
	if vars == "" {
		return token
	}
	varName := vars[7 : len(vars)-1]
	if len(varName) == 0 {
		return token
	}
	value, exists := tc.params[varName]
	if exists {
		listValue := []map[string]string{}
		err = json.Unmarshal([]byte(value), &listValue)
		if err != nil {
			log.Printf("engine replace token for variables error: list variables is need")
			return ""
		}
		return fmt.Sprintf("%d", len(listValue))
	}

	return token
}

//replaceList 针对list类型的数据,替换值
func replaceList(tc *scriptContext, threadContext *scriptThreadContext, token string) (list []map[string]string, err error) {
	params := replaceToken(tc, threadContext, token)
	if params == "" {
		return nil, nil
	}
	listValue := []map[string]string{}
	err = json.Unmarshal([]byte(params), &listValue)
	if err != nil {
		log.Printf("engine replace token for variables error: list variables is need")
		return nil, err
	}
	return listValue, nil

}

func statementParamExists(params map[string]string, paramName string) bool {
	_, exists := params[paramName]
	return exists
}

func parseScript(lines []string, tc *scriptContext) ([]scriptLineTokens, error, int) {
	allTokens := []scriptLineTokens{}
	ifExpressLineNo := []int{}
	forExpressLineNo := []int{}
	finallyExpressLineNo := []int{}
	for lineNo := 0; lineNo < len(lines); {
		if isCommentLine(lines[lineNo]) {
			lineNo++
			continue
		}
		tokens, err, newLineNo := parseLine(lines, lineNo)
		if err != nil {
			return nil, err, newLineNo
		}
		if len(tokens) > 0 {
			allTokens = append(allTokens, scriptLineTokens{tokens: tokens, lineNo: lineNo})
		}
		isExpress, isEnd := isExpressFunc(lines[lineNo])
		if len(ifExpressLineNo) != 0 {
			for i := 0; i < len(ifExpressLineNo); i++ {
				scriptBlockToken := tc.expressBlock[ifExpressLineNo[i]]
				thisToken := &scriptLineTokens{}
				thisToken.lineNo = lineNo
				thisToken.tokens = tokens
				scriptBlockToken.lines = append(scriptBlockToken.lines, thisToken)
			}
		}
		if len(forExpressLineNo) != 0 {
			for i := 0; i < len(forExpressLineNo); i++ {
				scriptBlockToken := tc.expressBlock[forExpressLineNo[i]]
				thisToken := &scriptLineTokens{}
				thisToken.lineNo = lineNo
				thisToken.tokens = tokens
				scriptBlockToken.lines = append(scriptBlockToken.lines, thisToken)
			}
		}
		if len(finallyExpressLineNo) != 0 {
			for i := 0; i < len(finallyExpressLineNo); i++ {
				scriptBlockToken := tc.expressBlock[finallyExpressLineNo[i]]
				thisToken := &scriptLineTokens{}
				thisToken.lineNo = lineNo
				thisToken.tokens = tokens
				scriptBlockToken.lines = append(scriptBlockToken.lines, thisToken)
			}
		}
		if isExpress != 0 {
			if isEnd {
				if isExpress == EXP_IF {
					tc.expressBlock[ifExpressLineNo[len(ifExpressLineNo)-1]].endLineNo = lineNo
					ifExpressLineNo = removeListSuffix(ifExpressLineNo)
				}
				if isExpress == EXP_FOR {
					tc.expressBlock[forExpressLineNo[len(forExpressLineNo)-1]].endLineNo = lineNo
					forExpressLineNo = removeListSuffix(forExpressLineNo)
				}
				if isExpress == EXP_finally {
					tc.expressBlock[finallyExpressLineNo[len(finallyExpressLineNo)-1]].endLineNo = lineNo
					finallyExpressLineNo = removeListSuffix(finallyExpressLineNo)
				}
			} else {
				if isExpress == EXP_IF {
					ifExpressLineNo = append(ifExpressLineNo, lineNo)
				}
				if isExpress == EXP_FOR {
					forExpressLineNo = append(forExpressLineNo, lineNo)
				}
				if isExpress == EXP_finally {
					finallyExpressLineNo = append(finallyExpressLineNo, lineNo)
				}
				scriptBlockToken := &scriptBlockTokens{}
				scriptBlockToken.beginLineNo = lineNo
				scriptBlockToken.exp = isExpress
				scriptBlockToken.lines = []*scriptLineTokens{}
				tc.expressBlock[lineNo] = scriptBlockToken
			}
		}

		lineNo = newLineNo
	}
	//log.Printf("%q\n", allTokens)
	return allTokens, nil, len(lines)
}

//解析代码，从某一行开始，返回token以及是否有错误，以及返回下一行的位置
func parseLine(lines []string, lineNo int) ([]string, error, int) {
	line := lines[lineNo]
	newLineNo := lineNo + 1
	isQuote := false
	token := ""
	tokens := []string{}
	lineLen := len(line)
	//log.Printf("Starting parse line %d: '%s'", lineNo, line)
	for i := 0; i < lineLen; i++ {
		c := line[i]
		switch c {
		case ' ', '\t':
			if isQuote {
				token += string(c)
				continue
			}
			if token != "" {
				//log.Printf("token=%s\n", token)
				tokens = append(tokens, token)
				token = ""
			}
		case '"':
			if isQuote {
				// if token == "" && i == 2 {
				// 	return nil, &ScriptError{ErrorMessage: "can not be empty quote", ErrorLineNo: lineNo + 1, ColumnNo: i + 1, Staging: TCS_PARSING}, lineNo
				// }
				tokens = append(tokens, token)
				token = ""
				isQuote = false
				continue
			}
			//if isQuote==false
			isQuote = true

		case '`':
			//如果在双引号后面
			if isQuote {
				token += string(c)
				continue
			}
			//处理行剩余部分：
			//行剩余部分应该为空，而不应该有其他字符
			rightLine := strings.TrimSpace(string(line[i+1:]))
			if len(rightLine) != 0 {
				return nil, &ScriptError{ErrorMessage: "multiple line start tag char must be at line end", ErrorLineNo: lineNo + 1, ColumnNo: i + 1, Staging: TCS_PARSING}, lineNo
			}
			for newLineNo = lineNo + 1; newLineNo < len(lines); newLineNo++ {
				newLine := lines[newLineNo]
				if strings.TrimSpace(newLine) == "`" {
					if token == "" {
						return nil, &ScriptError{ErrorMessage: "can not be empty multiple line token", ErrorLineNo: newLineNo + 1, ColumnNo: len(newLine), Staging: TCS_PARSING}, lineNo
					}
					tokens = append(tokens, token)
					return tokens, nil, newLineNo + 1
				}
				token = token + "\n" + newLine
			}
			//正常来说代码不应该到这里，因为多行毕竟有一个结束符
			return nil, &ScriptError{ErrorMessage: "multiple line end tag char not found", ErrorLineNo: len(lines), ColumnNo: len(lines[len(lines)-1]), Staging: TCS_PARSING}, lineNo

		default:
			token += string(c)
			//log.Printf("token=%s\n", token)
		}
	} //end for i
	//如果行末不是双引号，报错
	if isQuote {
		return nil, &ScriptError{ErrorMessage: "quote is not completed", ErrorLineNo: lineNo + 1, ColumnNo: len(line), Staging: TCS_PARSING}, lineNo
	}
	if token != "" {
		tokens = append(tokens, token)
	}
	//log.Printf("end parse line %d: '%s',tokens:\n", lineNo, line)
	//log.Printf("%q\n", tokens)
	return tokens, nil, newLineNo
}

//isExpress #if this line is express,like if, for, else if
func isExpressFunc(line string) (expresstype int, isEnd bool) {
	expresstype = 0
	lowerLine := strings.ToLower(strings.TrimSpace(line))
	for _, v := range ifExpress {
		if strings.HasPrefix(lowerLine, v) {
			expresstype = EXP_IF
		}
	}
	for _, v := range forExpress {
		if strings.HasPrefix(lowerLine, v) {
			expresstype = EXP_FOR
		}
	}
	for _, v := range finallyExpress {
		if strings.HasPrefix(lowerLine, v) {
			expresstype = EXP_finally
		}
	}

	if expresstype == 0 {
		return expresstype, false
	}
	if strings.HasPrefix(lowerLine, "end") {
		return expresstype, true
	}
	return expresstype, false
}

//getBiggerOne get the bigger one from the two number
func getBiggerOne(exp1, exp2 int) (exp int) {

	if exp1 > exp2 {
		return exp1
	} else if exp1 < exp2 {
		return exp2
	} else {
		return 0
	}
}
