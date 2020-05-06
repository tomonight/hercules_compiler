package script_engine

import (
	"errors"
	"fmt"
	executorPkg "hercules_compiler/rce-executor/executor"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	TCS_PARSING = 1
	TCS_RUNNING = 2

	EFF_CONTINUE     = 2
	EFF_STOP         = 1
	EFS_CONTINUE     = 3
	EFS_STOP         = 4
	EFS_SKIP         = 5
	CFF_SKIPALL_FLAG = 6
	CFF_SKIP_ACTION  = 7

	//callback result type
	CRT_PARSE_COMPELTED     = 1
	CRT_PARSE_FAILED        = 2
	CRT_COMMAND_COMPLETED   = 3
	CRT_COMMAND_FAILED      = 4
	CRT_COMMAND_SKIPED      = 5
	CRT_STATEMENT_COMPLETED = 6
	CRT_STATEMENT_FAILED    = 7
	CRT_SCRIPT_COMPLETED    = 8
	CRT_SCRIPT_FAILED       = 9

	//express code
	EXP_IF      = 1
	EXP_FOR     = 2
	EXP_finally = 3

	FINALLY_ANY     = 0
	FINALLY_SUCCESS = 1
	FINALLY_FAIL    = 2

	EXEC_FAIL    = -1
	EXEC_SUCCESS = 2
	EXEC_BREAK   = 1
)

type engineExecutor struct {
	Host     string
	executor executorPkg.Executor
}

type ScriptError struct {
	ErrorMessage string
	ErrorLineNo  int //注意从1计数，以便正确显示
	ColumnNo     int
	Staging      int
}

func (e *ScriptError) Error() string {
	return e.ErrorMessage
}

type scriptLineTokens struct {
	tokens []string
	lineNo int
}

//scriptBlockTokens 定义脚本代码块
type scriptBlockTokens struct {
	lines       []*scriptLineTokens
	exp         int
	beginLineNo int
	endLineNo   int
	successType int
}
type ScriptEngineCallback func(scriptExecID string, callbackResultType int, callbackResultInfo string)

type scriptThreadContext struct {
	envs                  map[string]string
	vars                  map[string]string
	resultData            map[string]string
	currentTargetName     string
	currentTargetHost     string
	workingPath           string
	workingUser           string
	enableSudo            bool
	executeFailedFlag     int
	executeSuccessfulFlag int
	connectFailedFlag     int
	currentStatementLine  int
	threadTokens          []scriptLineTokens
	executeFailed         bool
	executeSkipNext       bool
}

type scriptContext struct {
	targets              map[string]*engineExecutor
	params               map[string]string
	mainThreadContext    *scriptThreadContext
	threadContexts       map[string]*scriptThreadContext
	staging              int
	script               []string
	scriptName           string
	scriptExecID         string
	currentStatementLine int
	callback             ScriptEngineCallback
	inParallel           bool
	parallelPreparing    bool //并行是否在准备阶段
	prepareTargetName    string
	expressBlock         map[int]*scriptBlockTokens
	returnParams         map[string]string
}

//回调前，屏蔽密码
func decorateMaskPasswordCallback(callback ScriptEngineCallback) ScriptEngineCallback {
	return func(scriptExecID string, callbackResultType int, callbackResultInfo string) {
		re, err := regexp.Compile(`(?i)password[\s|\S]*`)
		if err == nil {
			callbackResultInfo = re.ReplaceAllString(callbackResultInfo, "password ********")
		}
		re, err = regexp.Compile(`(?i)passwd[\s|\S]*`)
		if err == nil {
			callbackResultInfo = re.ReplaceAllString(callbackResultInfo, "passwd ********")
		}

		callback(scriptExecID, callbackResultType, callbackResultInfo)
	}
}
func ExecuteScript(script []string, params map[string]string, scriptExecID string, callback ScriptEngineCallback) {
	tc := scriptContext{}
	tc.params = params
	tc.staging = TCS_PARSING
	tc.scriptExecID = scriptExecID
	tc.mainThreadContext = &scriptThreadContext{}
	currentThreadContext := tc.mainThreadContext
	currentThreadContext.envs = make(map[string]string)
	currentThreadContext.vars = make(map[string]string)
	currentThreadContext.executeFailedFlag = EFF_STOP
	currentThreadContext.executeSuccessfulFlag = EFS_CONTINUE
	currentThreadContext.enableSudo = false

	tc.script = script
	if callback == nil {
		tc.callback = decorateMaskPasswordCallback(defaultCallback)
	} else {
		tc.callback = decorateMaskPasswordCallback(callback)
	}
	tc.targets = make(map[string]*engineExecutor)
	tc.inParallel = false
	tc.threadContexts = nil
	tc.expressBlock = make(map[int]*scriptBlockTokens)
	tc.returnParams = make(map[string]string)
	allTokens, err, _ := parseScript(script, &tc)
	if err != nil {
		switch err.(type) {
		case *ScriptError:
			scriptErr := err.(*ScriptError)
			msg := fmt.Sprintf("parse script script error at line %d column %d:%s, statement='%s'", scriptErr.ErrorLineNo+1, scriptErr.ColumnNo, scriptErr.Error(), script[scriptErr.ErrorLineNo-1])
			tc.callback(scriptExecID, CRT_PARSE_FAILED, msg)
			tc.callback(scriptExecID, CRT_SCRIPT_FAILED, "script run with failed")
			return
		default:
			tc.callback(scriptExecID, CRT_PARSE_FAILED, err.Error())
			tc.callback(scriptExecID, CRT_SCRIPT_FAILED, "script run with failed")
			return
		}
	}
	if len(allTokens) == 0 {
		tc.callback(scriptExecID, CRT_SCRIPT_FAILED, "has no any valid statement")
		return
	}

	for _, lineTokens := range allTokens {
		processFunc, exists := scriptStatementFuncs[strings.ToLower(lineTokens.tokens[0])]
		//如果是 endif 或者endfor 不做任何操作
		if strings.ToLower(lineTokens.tokens[0]) == "endif" || strings.ToLower(lineTokens.tokens[0]) == "endfor" || strings.ToLower(lineTokens.tokens[0]) == "endfinally" {
			continue
		}
		if !exists {
			msg := fmt.Sprintf("parse script script error at line %d: %s, statement='%s'", lineTokens.lineNo+1, "unknown commmand", script[lineTokens.lineNo])
			tc.callback(scriptExecID, CRT_PARSE_FAILED, msg)
			tc.callback(scriptExecID, CRT_SCRIPT_FAILED, "script run with failed")
			return
		}
		_, _, err := processFunc(&tc, currentThreadContext, lineTokens)
		if err != nil {
			msg := fmt.Sprintf("parse script script error at line %d: %s, statement='%s'", lineTokens.lineNo+1, err.Error(), script[lineTokens.lineNo])
			tc.callback(scriptExecID, CRT_PARSE_FAILED, msg)
			tc.callback(scriptExecID, CRT_SCRIPT_FAILED, "script run with failed")
			return
		}
	}
	//在最后不应该还处于并行状态
	if tc.inParallel {
		tc.callback(scriptExecID, CRT_PARSE_FAILED, "script has no end parallel statement")
		return
	}
	tc.callback(scriptExecID, CRT_PARSE_COMPELTED, "script parse completed")
	//正式执行
	tc.staging = TCS_RUNNING
	//重新初始化一些对象
	tc.targets = make(map[string]*engineExecutor)
	tc.inParallel = false
	tc.threadContexts = nil
	tc.parallelPreparing = false
	tc.prepareTargetName = ""

	//关闭所有的连接
	defer func() {
		for _, target := range tc.targets {
			if target != nil {
				target.executor.Close()
			}
		}
	}()
	_, err, str := execTokensFunc(&tc, currentThreadContext, allTokens, len(script)-1)
	defer func() {
		for _, finallyBlock := range tc.expressBlock {
			if finallyBlock.exp == EXP_finally {
				if finallyBlock.successType == FINALLY_ANY || ((finallyBlock.successType == FINALLY_FAIL && err != nil) || (finallyBlock.successType == FINALLY_SUCCESS && err == nil)) {
					ftokens := []scriptLineTokens{}
					for _, v := range finallyBlock.lines {
						ftokens = append(ftokens, *v)
					}
					execTokensFunc(&tc, currentThreadContext, ftokens, finallyBlock.endLineNo)
				}
			}
		}
	}()

	if err != nil {
		tc.callback(scriptExecID, CRT_SCRIPT_FAILED, "script run with failed")
		return
	}
	if str == "" {
		str = "script run completed"
	}
	tc.callback(scriptExecID, CRT_SCRIPT_COMPLETED, str)
}

func connectStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if tc.inParallel {
		return tokenLine.lineNo, "", fmt.Errorf("connect target statement can not in parallel block")
	}
	if len(tokens) < 4 {
		return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement")
	}
	if strings.ToLower(tokens[1]) != "target" {
		return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement")
	}
	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}
	params, err := tokensToParams(tc, threadContext, tokens, 3)
	//log.Printf("connect params=%q\n", params)
	if err != nil {
		return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement,%s", err.Error())
	}
	protocol, exists := params["protocol"]
	if !exists {
		return tokenLine.lineNo, "", fmt.Errorf("connect target:connect protocol not specified")
	}

	delete(params, "protocol")
	targetName := replaceToken(tc, threadContext, tokens[2])

	switch strings.ToLower(protocol) {
	case "ssh":
		if (!statementParamExists(params, "host")) || (!statementParamExists(params, "username")) {
			return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement, require host and username parameters")
		}

		password := ""
		if v, exists := params["password"]; exists {
			password = v
		}
		// path, err := executorPkg.GetSShPath()
		// if err != nil {
		// 	return tokenLine.lineNo, "", fmt.Errorf("connect target:connect protocol not specified")
		// }
		// path := "/MyData/zmysql"
		// path, err = executorPkg.GetSShPath()
		// if err != nil {
		// 	return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement, require host and username parameters")
		// }
		// keyfile := os.ExpandEnv(path + "/data/.ssh/id_rsa")
		// if v, exists := params["keyfile"]; exists {
		// 	keyfile = v
		// }

		port := 22
		if v, exists := params["port"]; exists {
			port, err = strconv.Atoi(v)
			if err != nil {
				return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement, inccorrect port:%s parameter", params["port"])
			}
		}
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			tc.targets[targetName] = nil
			return tokenLine.lineNo, "", nil
		}
		executor, err := executorPkg.NewSSHAgentExecutor(params["host"], params["username"], password, port)
		if err != nil {
			//log.Printf("connect ssh error:%s\n", err.Error())
			if threadContext.connectFailedFlag == CFF_SKIPALL_FLAG {
				threadContext.connectFailedFlag = CFF_SKIP_ACTION
			}

			return tokenLine.lineNo, "", fmt.Errorf("connect ssh target '%s:%d' failed: %s", params["host"], port, err.Error())
		}

		target, exists := tc.targets[targetName]
		if exists {
			if target != nil {
				target.executor.Close()
			}
		}
		newTarget := &engineExecutor{executor: executor, Host: params["host"]}
		tc.targets[targetName] = newTarget
		return tokenLine.lineNo, fmt.Sprintf("target '%s' connected", newTarget.Host), nil
	case "rce":
		if statementParamExists(params, "host") == false {
			return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement, require host parameters")
		}
		port := 5051
		if v, exists := params["port"]; exists {
			port, err = strconv.Atoi(v)
			if err != nil {
				return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement, inccorrect port parameter")
			}
		}
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			tc.targets[targetName] = nil
			return tokenLine.lineNo, "", nil
		}
		executor, err := executorPkg.NewRCEAgentExecutor(params["host"], fmt.Sprintf("%d", port))
		if err != nil {
			return tokenLine.lineNo, "", fmt.Errorf("connect rce target '%s:%d' failed: %s", params["host"], port, err.Error())
		}

		target, exists := tc.targets[targetName]
		if exists {
			if target != nil {
				target.executor.Close()
			}
		}
		newTarget := &engineExecutor{executor: executor, Host: params["host"]}
		tc.targets[targetName] = newTarget
		return tokenLine.lineNo, fmt.Sprintf("target '%s' connected", newTarget.Host), nil

	//support zcloud agent
	case "zcagent":
		if statementParamExists(params, "host") == false {
			return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement, require host parameters")
		}
		//zcloud default agent port is 8100
		port := 8100
		if v, exists := params["port"]; exists {
			port, err = strconv.Atoi(v)
			if err != nil {
				return tokenLine.lineNo, "", fmt.Errorf("connect target:incorrect connect statement, inccorrect port parameter")
			}
		}
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			tc.targets[targetName] = nil
			return tokenLine.lineNo, "", nil
		}
		executor, err := executorPkg.NewZCAgentExecutor(params["host"], port)
		if err != nil {
			return tokenLine.lineNo, "", fmt.Errorf("connect rce target '%s:%d' failed: %s", params["host"], port, err.Error())
		}

		target, exists := tc.targets[targetName]
		if exists {
			if target != nil {
				target.executor.Close()
			}
		}
		newTarget := &engineExecutor{executor: executor, Host: params["host"]}
		tc.targets[targetName] = newTarget
		return tokenLine.lineNo, fmt.Sprintf("target '%s' connected", newTarget.Host), nil

		//连接到本地
	case "local":
		executor, err := executorPkg.NewLocalExecutor()
		if err != nil {
			return tokenLine.lineNo, "", fmt.Errorf("connect local target failed:%s", err.Error())
		}

		target, exists := tc.targets[targetName]
		if exists {
			if target != nil {
				target.executor.Close()
			}
		}
		newTarget := &engineExecutor{executor: executor, Host: "localhost"}
		tc.targets[targetName] = newTarget
		return tokenLine.lineNo, fmt.Sprintf("target '%s' connected", newTarget.Host), nil
	default:
		return tokenLine.lineNo, "", fmt.Errorf("connect target:not supported connect protocol")
	}
	return tokenLine.lineNo, "", fmt.Errorf("connect target failed:unknown error")
}

func executeStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	threadContext.resultData = nil
	if len(tokens) < 2 {
		return tokenLine.lineNo, "", fmt.Errorf("execute command: incorrect execute statement")
	}
	commandName := tokens[1]

	pair := strings.SplitN(commandName, ".", 2)
	if len(pair) != 2 {
		return tokenLine.lineNo, "", fmt.Errorf("execute command: incorrect execute comand, command=%s", commandName)
	}
	cmdFunc, ok := executorPkg.GetCmdByModuleAndName(pair[0], pair[1])

	if !ok {
		return tokenLine.lineNo, "", fmt.Errorf("execute command: can not found command '%s'", commandName)
	}
	params, err := tokensToParams(tc, threadContext, tokens, 2)
	if err != nil {
		return tokenLine.lineNo, "", fmt.Errorf("execute command: incorrect execute statement, parse command parameter %s", err.Error())
	}
	//如果是解析阶段就返回
	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}

	if threadContext.currentTargetName == "" {
		return tokenLine.lineNo, "", fmt.Errorf("execute command: target host has not specified\n")
	}
	target, exists := tc.targets[threadContext.currentTargetName]
	if !exists {
		return tokenLine.lineNo, "", fmt.Errorf("execute command: current target not found,target='%s'\n", threadContext.currentTargetName)
	}
	execParams := executorPkg.ExecutorCmdParams{}
	for k, v := range params {
		execParams[k] = v
	}
	for k, v := range threadContext.envs {
		target.executor.SetEnv(k, v)
	}
	target.executor.SetWorkingPath(threadContext.workingPath)
	target.executor.SetExecuteUser(threadContext.workingUser)
	target.executor.SetSudoEnabled(threadContext.enableSudo)

	//add timeout param
	if timeout, exsit := execParams["timeout"]; exsit {
		timeoutInt, _ := strconv.Atoi(timeout.(string))
		target.executor.SetTimeOut(timeoutInt)
	}
	er := cmdFunc(target.executor, &execParams)

	if !er.Successful {
		if er.ExitCode == 124 {
			er.Message = "execute timeout"
		}
		return tokenLine.lineNo, "", fmt.Errorf("execute command failed: '%s'", er.Message)
	} else {
		threadContext.resultData = er.ResultData
		return tokenLine.lineNo, fmt.Sprintf("%s on target host '%s': %s", commandName, threadContext.currentTargetHost, er.Message), nil
	}
}

func setStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) < 3 {
		return tokenLine.lineNo, "", fmt.Errorf("set statement: incorrect set statement")
	}
	if threadContext.connectFailedFlag == CFF_SKIP_ACTION {
		return tokenLine.lineNo, "skip set statement", nil
	}

	switch strings.ToLower(tokens[1]) {
	case "target":
		if tc.inParallel {
			return tokenLine.lineNo, "", fmt.Errorf("set target statement can not in parallel block")
		}
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		threadContext.currentTargetName = ""
		threadContext.currentTargetHost = ""
		targetName := replaceToken(tc, threadContext, tokens[2])
		target, exists := tc.targets[targetName]
		if !exists {
			return tokenLine.lineNo, "", fmt.Errorf("set target statement: can not found target '%s'", targetName)
		}
		threadContext.currentTargetName = replaceToken(tc, threadContext, tokens[2])
		threadContext.currentTargetHost = target.Host
		return tokenLine.lineNo, fmt.Sprintf("set target statement: switch to target '%s:%s'", targetName, target.Host), nil
	case "var":
		pair, err := tokensToParams(tc, threadContext, tokens, 2)
		if err != nil {
			//log.Printf("tokens=%q\n", tokens)
			return tokenLine.lineNo, "", fmt.Errorf("set var statement: incorrect statement,%s", err.Error())
		}
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		for k, v := range pair {
			threadContext.vars[k] = v
		}
		return tokenLine.lineNo, fmt.Sprintf("set var statement completed"), nil
	case "env":
		pair, err := tokensToParams(tc, threadContext, tokens, 2)
		if err != nil {
			return tokenLine.lineNo, "", fmt.Errorf("set env statement: incorrect statement")
		}
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		for k, v := range pair {
			threadContext.envs[k] = v
		}
		return tokenLine.lineNo, fmt.Sprintf("set env statement completed"), nil
	case "pwd":
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		threadContext.workingPath = replaceToken(tc, threadContext, tokens[2])
		return tokenLine.lineNo, fmt.Sprintf("set pwd statement completed"), nil
	case "user":
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		threadContext.workingUser = replaceToken(tc, threadContext, tokens[2])
		return tokenLine.lineNo, fmt.Sprintf("set user statement completed"), nil
	case "exec":
		if len(tokens) < 4 {
			return tokenLine.lineNo, "", fmt.Errorf("set statement: incorrect set exec statement")
		}
		if tokens[2] != "failed" && tokens[2] != "successful" {
			return tokenLine.lineNo, "", fmt.Errorf("set statement: incorrect set exec statement")
		}
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		if tokens[2] == "failed" {
			switch tokens[3] {
			case "stop":
				threadContext.executeFailedFlag = EFF_STOP
			case "continue":
				threadContext.executeFailedFlag = EFF_CONTINUE
			default:
				return tokenLine.lineNo, "", fmt.Errorf("set statement: incorrect set exec statement")
			}
		}
		if tokens[2] == "successful" {
			switch tokens[3] {
			case "stop":
				threadContext.executeSuccessfulFlag = EFS_STOP
			case "continue":
				threadContext.executeSuccessfulFlag = EFS_CONTINUE
			case "skip":
				threadContext.executeSuccessfulFlag = EFS_SKIP
			default:
				return tokenLine.lineNo, "", fmt.Errorf("set statement: incorrect set exec statement")
			}
		}
		return tokenLine.lineNo, fmt.Sprintf("set exec statement completed"), nil
	case "sudo":
		if len(tokens) != 3 {
			return tokenLine.lineNo, "", fmt.Errorf("set statement: incorrect set sudo statement")
		}
		if tokens[2] != "enable" && tokens[2] != "disable" {
			return tokenLine.lineNo, "", fmt.Errorf("set statement: incorrect set sudo statement")
		}
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		if tokens[2] == "enable" {
			threadContext.enableSudo = true
		} else {
			threadContext.enableSudo = false
		}

	case "connect":
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		if len(tokens) < 4 {
			return tokenLine.lineNo, "", fmt.Errorf("set statement: incorrect set connect statement")
		}
		//目前只处理failed流程，successful目前不做处理，默认继续执行
		if tokens[2] != "failed" {
			return tokenLine.lineNo, "", fmt.Errorf("set statement: incorrect set connect statement")
		}
		if tokens[3] == "skipall" {
			threadContext.connectFailedFlag = CFF_SKIPALL_FLAG
		}
		return tokenLine.lineNo, fmt.Sprintf("set connect statement completed"), nil
	}
	return tokenLine.lineNo, "", fmt.Errorf("not supported set statement")
}

func sleepStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) != 2 {
		return tokenLine.lineNo, "", fmt.Errorf("sleep statement: incorrect sleep statement")
	}
	str := tokens[1]
	value, err := strconv.Atoi(str)
	if err != nil {
		return tokenLine.lineNo, "", fmt.Errorf("sleep statement: incorrect sleep time duration")
	}
	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}
	time.Sleep(time.Duration(value) * time.Millisecond)
	return tokenLine.lineNo, fmt.Sprintf("sleep for %d ms", value), nil
}
func resetStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) < 2 {
		return tokenLine.lineNo, "", fmt.Errorf("reset statement: incorrect reset statement")
	}
	switch strings.ToLower(tokens[1]) {
	case "pwd":
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		threadContext.workingPath = ""
		return tokenLine.lineNo, fmt.Sprintf("reset pwd statement completed"), nil
	case "user":
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		threadContext.workingUser = ""
		return tokenLine.lineNo, fmt.Sprintf("reset user statement completed"), nil
	case "env":
		if len(tokens) != 3 {
			return tokenLine.lineNo, "", fmt.Errorf("reset statement: incorrect reset statement")
		}
		//如果是解析阶段就返回
		if tc.staging == TCS_PARSING {
			return tokenLine.lineNo, "", nil
		}
		delete(threadContext.envs, tokens[2])
		return tokenLine.lineNo, fmt.Sprintf("reset env statement completed"), nil
	}
	return tokenLine.lineNo, "", fmt.Errorf("not supported reset statement")
}

func beginStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	threadContext.resultData = nil
	if (len(tokens) != 2) && (tokens[1] != "parallel") {
		return tokenLine.lineNo, "", fmt.Errorf("begin parallel: incorrect begin parallel statement")
	}
	if tc.inParallel {
		return tokenLine.lineNo, "", fmt.Errorf("begin parallel: parallel can not be nested")
	}
	tc.inParallel = true
	tc.parallelPreparing = true
	tc.threadContexts = make(map[string]*scriptThreadContext)

	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}
	if threadContext.executeSuccessfulFlag == EFS_SKIP {
		threadContext.executeSuccessfulFlag = EFS_CONTINUE
	}
	return tokenLine.lineNo, "begin parallel executing", nil
}

func printStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) != 2 {
		return tokenLine.lineNo, "", fmt.Errorf("print: incorrect printf statement")
	}
	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}
	str := ""
	for i := 1; i < len(tokens); i++ {
		str = str + replaceToken(tc, threadContext, tokens[i]) + " "
	}
	return tokenLine.lineNo, str, nil
}

func parallelStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if (len(tokens) != 3) && (tokens[1] != "target") {
		return tokenLine.lineNo, "", fmt.Errorf("parallel target: incorrect parallel target statement")
	}
	if !tc.inParallel {
		return tokenLine.lineNo, "", fmt.Errorf("parallel target: parallel can not be out of parallel block")
	}

	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}

	targetName := replaceToken(tc, threadContext, tokens[2])

	if _, exists := tc.threadContexts[targetName]; exists {
		return tokenLine.lineNo, "", fmt.Errorf("parallel target: target '%s' duplicated in same parallel block", targetName)
	}

	target, exists := tc.targets[targetName]
	if !exists {
		return tokenLine.lineNo, "", fmt.Errorf("parallel target: target '%s' not defined", targetName)
	}
	tc.prepareTargetName = targetName

	if tc.staging == TCS_PARSING {
		context := scriptThreadContext{currentTargetName: targetName}
		tc.threadContexts[targetName] = &context
		return tokenLine.lineNo, "", nil
	}
	context := scriptThreadContext{currentTargetName: targetName}
	context.executeFailedFlag = threadContext.executeFailedFlag

	if threadContext.executeSuccessfulFlag == EFS_SKIP {
		context.executeSuccessfulFlag = EFS_CONTINUE
	} else {
		context.executeSuccessfulFlag = threadContext.executeSuccessfulFlag
	}
	context.workingPath = threadContext.workingPath
	context.executeSkipNext = false
	context.currentTargetHost = target.Host
	context.workingUser = threadContext.workingUser
	context.enableSudo = threadContext.enableSudo
	context.envs = make(map[string]string)
	context.vars = make(map[string]string)
	context.threadTokens = []scriptLineTokens{}
	for k, v := range threadContext.envs {
		context.envs[k] = v
	}
	for k, v := range threadContext.vars {
		context.vars[k] = v
	}
	threadContext.resultData = nil

	tc.threadContexts[targetName] = &context
	return tokenLine.lineNo, "start parallel on target host " + context.currentTargetHost, nil
}

func endStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	threadContext.resultData = nil
	if (len(tokens) != 2) || (tokens[1] != "parallel") {
		return tokenLine.lineNo, "", fmt.Errorf("end parallel: incorrect end parallel statement")
	}
	if !tc.inParallel {
		return tokenLine.lineNo, "", fmt.Errorf("end parallel: parallel has not started")
	}

	tc.inParallel = false
	tc.prepareTargetName = ""
	tc.parallelPreparing = false

	if tc.staging == TCS_PARSING {
		tc.threadContexts = nil
		return tokenLine.lineNo, "", nil
	}
	threadCount := len(tc.threadContexts)
	var wg sync.WaitGroup
	wg.Add(threadCount)
	// skipFlag := false
	for _, inThreadContext := range tc.threadContexts {
		go func(context *scriptThreadContext) {
			defer wg.Done()
			// for _, lineTokens := range context.threadTokens {

			// 	funcName := lineTokens.tokens[0]
			// 	processFunc, exists := scriptStatementFuncs[strings.ToLower(funcName)]
			// 	if skipFlag {
			// 		msg := fmt.Sprintf("skip execute script statement on target host '%s',  at line %dstatement='%s'", threadContext.currentTargetHost, lineTokens.lineNo+1, tc.script[lineTokens.lineNo])
			// 		tc.callback(tc.scriptExecID, CRT_STATEMENT_FAILED, msg)
			// 		continue
			// 	}
			// 	//编译时确认过，所以不应该有找不到的情况出现
			// 	if !exists {
			// 		msg := fmt.Sprintf("execute script statement on target host '%s', error at line %d,can not find command %s, statement='%s'", threadContext.currentTargetHost, lineTokens.lineNo+1, funcName, tc.script[lineTokens.lineNo])
			// 		tc.callback(tc.scriptExecID, CRT_STATEMENT_FAILED, msg)
			// 		context.executeFailed = true
			// 		return
			// 	}
			// 	if funcName == "execute" && context.executeSkipNext {
			// 		context.executeSkipNext = false
			// 		msg := fmt.Sprintf("execute script statement skip on target host '%s' at line %d: '%s'", threadContext.currentTargetHost, lineTokens.lineNo, tc.script[lineTokens.lineNo])
			// 		tc.callback(tc.scriptExecID, CRT_COMMAND_SKIPED, msg)
			// 		continue
			// 	}
			// 	_, result, err := processFunc(tc, context, lineTokens)
			// 	if err != nil {
			// 		msg := fmt.Sprintf("execute script statement on target host '%s',error at line %d: %s, statement='%s'", threadContext.currentTargetHost, lineTokens.lineNo+1, err.Error(), tc.script[lineTokens.lineNo])

			// 		if context.connectFailedFlag == CFF_SKIP_ACTION {
			// 			tc.callback(tc.scriptExecID, CRT_COMMAND_COMPLETED, msg)
			// 			skipFlag = true
			// 			tc.callback(tc.scriptExecID, CRT_COMMAND_COMPLETED, "connect fail ,skip all statement")
			// 			continue
			// 		}
			// 		if funcName == "execute" {
			// 			if context.executeFailedFlag == EFF_CONTINUE {
			// 				tc.callback(tc.scriptExecID, CRT_COMMAND_FAILED, msg)
			// 			} else {
			// 				tc.callback(tc.scriptExecID, CRT_STATEMENT_FAILED, msg)
			// 				context.executeFailed = true
			// 				return
			// 			}
			// 		}
			// 	} else {
			// 		if funcName == "execute" {
			// 			tc.callback(tc.scriptExecID, CRT_COMMAND_COMPLETED, result)
			// 			if context.executeSuccessfulFlag == EFS_STOP {
			// 				return
			// 			}
			// 			if context.executeSuccessfulFlag == EFS_SKIP {
			// 				context.executeSkipNext = true
			// 				context.executeSuccessfulFlag = EFS_CONTINUE
			// 			}
			// 		} else {
			// 			tc.callback(tc.scriptExecID, CRT_STATEMENT_COMPLETED, result)
			// 		}
			// 	}
			// }
			endLine := context.threadTokens[len(context.threadTokens)-1].lineNo
			_, err, str := execTokensFunc(tc, context, context.threadTokens, endLine)
			if err != nil {
				tc.callback(tc.scriptExecID, CRT_SCRIPT_FAILED, err.Error())
				context.executeFailed = true
			}
			if str == "" {
				str = "script run completed"
			}

		}(inThreadContext)
	} // end for
	wg.Wait()
	for _, inThreadContext := range tc.threadContexts {
		if inThreadContext.executeFailed {
			tc.threadContexts = nil
			return tokenLine.lineNo, "", fmt.Errorf("parallel completed with some failed")
		}
	}
	tc.threadContexts = nil
	return tokenLine.lineNo, "parallel completed", nil
}

func ifStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) < 4 {
		return tokenLine.lineNo, "", fmt.Errorf("incorrect if express")
	}
	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}
	tokens = replaceAllToken(tc, threadContext, tokens)
	express1left := tokens[1]
	express1Condition := tokens[2]
	express1right := tokens[3]
	condition := ""
	express1Result, err := getIfExpressTrueOrFalse(express1left, express1Condition, express1right)
	if err != nil {
		return tokenLine.lineNo, "", err
	}
	express2Result := true
	if len(tokens) > 4 {
		condition = tokens[4]
		express2left := tokens[5]
		express2Condition := tokens[6]
		express2right := tokens[7]
		express2Result, err = getIfExpressTrueOrFalse(express2left, express2Condition, express2right)
		if err != nil {
			return tokenLine.lineNo, "", err
		}
	}
	expressFinalResult := express1Result
	if condition != "" {
		if condition == "&&" || condition == "and" {
			expressFinalResult = express1Result && express2Result
		} else if condition == "||" || condition == "or" {
			expressFinalResult = express1Result || express2Result
		} else {
			return tokenLine.lineNo, "", fmt.Errorf("express condtion is invalid")
		}

	}

	if expressFinalResult {
		return tokenLine.lineNo, "", nil
	} else {
		scriptBlockTokens := tc.expressBlock[tokenLine.lineNo]
		return scriptBlockTokens.endLineNo, "", nil
	}
}

func forStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) != 1 && len(tokens) != 2 {
		return tokenLine.lineNo, "", fmt.Errorf("print: incorrect for statement")
	}
	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}
	if threadContext.resultData == nil {
		threadContext.resultData = make(map[string]string)
	}
	cycletimes := 0
	if len(tokens) == 2 {
		times, err := strconv.Atoi(tokens[1])
		if err == nil {
			cycletimes = times
		}
	}
	allTokens := []scriptLineTokens{}
	block := tc.expressBlock[tokenLine.lineNo]
	for _, value := range block.lines {
		allTokens = append(allTokens, *value)
	}
	returnStr := ""
	if len(tokens) == 1 {
		for true {
			execResult, err, str := execTokensFunc(tc, threadContext, allTokens, block.endLineNo)
			if err != nil {
				return tokenLine.lineNo, "", err
			}
			if execResult == EXEC_BREAK {
				break
			}
			if execResult == EXEC_SUCCESS {
				returnStr = str
				break
			}
		}
	} else if cycletimes > 0 {
		for i := 0; i < cycletimes; i++ {
			execResult, err, str := execTokensFunc(tc, threadContext, allTokens, block.endLineNo)
			if err != nil {
				return tokenLine.lineNo, "", err
			}
			if execResult == EXEC_BREAK {
				break
			}
			if execResult == EXEC_SUCCESS {
				returnStr = str
				break
			}
		}
	} else {
		nparams, err := replaceList(tc, threadContext, tokens[1])
		if err != nil {
			return tokenLine.lineNo, "", err
		}
		if nparams == nil {
			return tokenLine.lineNo, "", nil
		}
		for _, nparam := range nparams {
			for k, v := range nparam {
				threadContext.vars[k] = v
			}
			execResult, err, str := execTokensFunc(tc, threadContext, allTokens, block.endLineNo)
			if err != nil {
				return tokenLine.lineNo, "", err
			}
			if execResult == EXEC_BREAK {
				break
			}
			if execResult == EXEC_SUCCESS {
				returnStr = str
				if returnStr == "" {
					returnStr = "script run complete with return success"
				}
				break
			}
		}
	}
	scriptBlockTokens := tc.expressBlock[tokenLine.lineNo]
	return scriptBlockTokens.endLineNo, returnStr, nil
}

func breakStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) != 1 {
		return tokenLine.lineNo, "", fmt.Errorf("print: incorrect break statement")
	}
	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}
	block := getBreakExpressIn(tc, tokenLine.lineNo)
	return block.endLineNo, "", nil
}

func returnStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) != 1 && len(tokens) != 2 {
		return tokenLine.lineNo, "", fmt.Errorf("print: incorrect return statement")
	}
	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}

	if len(tokens) == 1 {
		reStr := ""
		if len(tc.returnParams) != 0 {
			reStr, _ = mapToString(tc.returnParams)
		}
		return tokenLine.lineNo, reStr, nil
	}
	msg := replaceToken(tc, threadContext, tokens[1])
	return tokenLine.lineNo, "", fmt.Errorf("return with error:%s", msg)
}

func finallyStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) > 2 {
		return tokenLine.lineNo, "", fmt.Errorf("print: incorrect finally statement")
	}
	block := tc.expressBlock[tokenLine.lineNo]
	if len(tokens) == 1 {
		block.successType = FINALLY_ANY
	} else {
		if tokens[1] == "success" {
			block.successType = FINALLY_SUCCESS
		} else {
			block.successType = FINALLY_FAIL
		}
	}
	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}
	return block.endLineNo, "", nil
}

func putStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	tokens := tokenLine.tokens
	if len(tokens) != 3 {
		return tokenLine.lineNo, "", fmt.Errorf("print: incorrect append statement")
	}

	key := replaceToken(tc, threadContext, tokens[1])
	if key == "" {
		return tokenLine.lineNo, "", fmt.Errorf("append command: incorrect apend statement, key can not be empty")
	}

	value := replaceToken(tc, threadContext, tokens[2])
	if value == "" {
		return tokenLine.lineNo, "", fmt.Errorf("append command: incorrect apend statement, value can not be empty")
	}

	if tc.staging == TCS_PARSING {
		return tokenLine.lineNo, "", nil
	}
	tc.returnParams[key] = value
	return tokenLine.lineNo, "", nil
}

func executeForStateMent(tc *scriptContext, threadContext *scriptThreadContext, allTokens []scriptLineTokens, endLine int) (err error) {
	execLine := -1
	for lineNo := 0; lineNo < len(allTokens); lineNo++ {
		if execLine == endLine {
			return
		}
		execLineb, tokenLine := getNextLine(execLine, allTokens)
		execLine = execLineb
		funcName := tokenLine.tokens[0]
		//如果是 endif 或者endfor 不做任何操作
		if strings.ToLower(funcName) == "endif" || strings.ToLower(funcName) == "endfor" || strings.ToLower(funcName) == "endfinally" {
			continue
		}
		processFunc, exists := scriptStatementFuncs[strings.ToLower(funcName)]
		//编译时确认过，所以不应该有找不到的情况出现
		if !exists {
			msg := fmt.Sprintf("execute script statement on target host '%s', error at line %d,can not find command %s, statement='%s'", threadContext.currentTargetHost, tokenLine.lineNo+1, funcName, tc.script[tokenLine.lineNo])
			tc.callback(tc.scriptExecID, CRT_STATEMENT_FAILED, msg)
			return
		}
		if funcName == "execute" && threadContext.executeSkipNext {
			msg := fmt.Sprintf("execute script statement skip on target host '%s' at line %d: '%s'", threadContext.currentTargetHost, tokenLine.lineNo, tc.script[tokenLine.lineNo])
			tc.callback(tc.scriptExecID, CRT_COMMAND_SKIPED, msg)
			return
		}
		result := ""
		var err error
		execLine, result, err = processFunc(tc, threadContext, tokenLine)
		if err != nil {
			msg := fmt.Sprintf("execute script statement on target host '%s',error at line %d: %s, statement='%s'", threadContext.currentTargetHost, tokenLine.lineNo+1, err.Error(), tc.script[tokenLine.lineNo])

			if funcName == "execute" {
				if threadContext.executeFailedFlag == EFF_CONTINUE {
					tc.callback(tc.scriptExecID, CRT_COMMAND_FAILED, msg)
				} else {
					tc.callback(tc.scriptExecID, CRT_STATEMENT_FAILED, msg)
					threadContext.executeFailed = true
					return errors.New("exec fail")
				}
			}
		} else {
			if funcName == "execute" {
				tc.callback(tc.scriptExecID, CRT_COMMAND_COMPLETED, result)
				if threadContext.executeSuccessfulFlag == EFS_STOP {
					return errors.New("exec fail")
				}
				if threadContext.executeSuccessfulFlag == EFS_SKIP {
					threadContext.executeSkipNext = true
					threadContext.executeSuccessfulFlag = EFS_CONTINUE
				}
			} else {
				tc.callback(tc.scriptExecID, CRT_STATEMENT_COMPLETED, result)
			}
		}
	}
	return nil
}

func defaultCallback(scriptExecID string, callbackResultType int, callbackResultInfo string) {
	//log.Printf("scriptExecID-%s:%d:%s", scriptExecID, callbackResultType, callbackResultInfo)
	log.Printf("scriptExecID-%s:%s", scriptExecID, callbackResultInfo)
}

type ScriptStatementProcessFunc func(tc *scriptContext, threadContext *scriptThreadContext, tokens scriptLineTokens) (int, string, error)

var scriptStatementFuncs map[string]ScriptStatementProcessFunc

func init() {
	scriptStatementFuncs = make(map[string]ScriptStatementProcessFunc)
	scriptStatementFuncs["connect"] = connectStatementFunc
	scriptStatementFuncs["execute"] = executeStatementFunc
	scriptStatementFuncs["set"] = setStatementFunc
	scriptStatementFuncs["reset"] = resetStatementFunc
	scriptStatementFuncs["parallel"] = parallelStatementFunc
	scriptStatementFuncs["begin"] = beginStatementFunc
	scriptStatementFuncs["end"] = endStatementFunc
	scriptStatementFuncs["sleep"] = sleepStatementFunc
	scriptStatementFuncs["print"] = printStatementFunc
	scriptStatementFuncs["if"] = ifStatementFunc
	scriptStatementFuncs["for"] = forStatementFunc
	scriptStatementFuncs["break"] = breakStatementFunc
	scriptStatementFuncs["put"] = putStatementFunc
	scriptStatementFuncs["return"] = returnStatementFunc
	scriptStatementFuncs["finally"] = finallyStatementFunc
	scriptStatementFuncs["test"] = testStatementFunc
	scriptStatementFuncs["wtest"] = wtestStatementFunc
}

func testStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	fmt.Println("test::::::::::::", tokenLine.lineNo+1, "--------", replaceToken(tc, threadContext, tokenLine.tokens[1]))

	return tokenLine.lineNo, "", nil
}
func wtestStatementFunc(tc *scriptContext, threadContext *scriptThreadContext, tokenLine scriptLineTokens) (int, string, error) {
	fmt.Println("test::::::::::::", tokenLine.lineNo+1)
	return tokenLine.lineNo, "", nil
}

func execTokensFunc(tc *scriptContext, threadContext *scriptThreadContext, allTokens []scriptLineTokens, endline int) (code int, err error, str string) {
	execLine := -1
	currentThreadContext := threadContext
	for lineNo := 0; lineNo < len(allTokens); lineNo++ {
		execLineb, lineTokens := getNextLine(execLine, allTokens)
		if execLineb == -1 {
			break
		}
		execLine = execLineb
		funcName := lineTokens.tokens[0]
		if tc.inParallel && tc.parallelPreparing && funcName != "parallel" && funcName != "end" {
			if tc.prepareTargetName != "" {
				threadContext, _ := tc.threadContexts[tc.prepareTargetName]
				threadContext.threadTokens = append(threadContext.threadTokens, lineTokens)
				continue
			}
		}
		//如果是 endif 或者endfor 不做任何操作
		if strings.ToLower(funcName) == "endif" || strings.ToLower(funcName) == "endfor" || strings.ToLower(funcName) == "endfinally" {
			continue
		}
		processFunc, exists := scriptStatementFuncs[strings.ToLower(funcName)]
		//编译时确认过，所以不应该有找不到的情况出现
		if !exists {
			msg := fmt.Sprintf("execute script statement error on target host '%s',at line %d,can not find command %s, statement='%s'", currentThreadContext.currentTargetHost, lineTokens.lineNo+1, funcName, tc.script[lineTokens.lineNo])
			tc.callback(tc.scriptExecID, CRT_STATEMENT_FAILED, msg)
			currentThreadContext.executeFailed = true
			tc.callback(tc.scriptExecID, CRT_SCRIPT_FAILED, "script run with failed")
		}
		if funcName == "execute" && currentThreadContext.executeSkipNext {
			currentThreadContext.executeSkipNext = false
			msg := fmt.Sprintf("execute script statement skip on target host '%s', at line %d: '%s'", currentThreadContext.currentTargetHost, lineTokens.lineNo+1, tc.script[lineTokens.lineNo])
			tc.callback(tc.scriptExecID, CRT_COMMAND_SKIPED, msg)
			continue
		}
		result := ""
		replacedStatement := replaceToken(tc, currentThreadContext, tc.script[lineTokens.lineNo])

		execLine, result, err = processFunc(tc, currentThreadContext, lineTokens)
		if funcName == "for" && result != "" {
			return EXEC_SUCCESS, nil, result
		}
		if funcName == "return" {
			if err != nil {
				tc.callback(tc.scriptExecID, CRT_STATEMENT_FAILED, err.Error())
				currentThreadContext.executeFailed = true
				return EXEC_FAIL, err, result
			}
			return EXEC_SUCCESS, nil, result
		}
		if err != nil {
			msg := fmt.Sprintf("execute script statement on target host '%s' error at line %d: %s, statement='%s'",
				currentThreadContext.currentTargetHost, lineTokens.lineNo+1, err.Error(), replacedStatement)
			if currentThreadContext.connectFailedFlag == CFF_SKIP_ACTION {
				tc.callback(tc.scriptExecID, CRT_COMMAND_COMPLETED, msg)
				tc.callback(tc.scriptExecID, CRT_COMMAND_COMPLETED, "connect fail ,skip all statement")
				continue
			}
			if funcName == "execute" {
				if currentThreadContext.executeFailedFlag == EFF_CONTINUE {
					//tc.callback(scriptExecID, CRT_COMMAND_FAILED, msg)
					//忽略错误
					tc.callback(tc.scriptExecID, CRT_COMMAND_COMPLETED, msg)
				} else {
					tc.callback(tc.scriptExecID, CRT_COMMAND_FAILED, msg)
					currentThreadContext.executeFailed = true
					// tc.callback(tc.scriptExecID, CRT_SCRIPT_FAILED, "script run with failed")
					return EXEC_FAIL, fmt.Errorf("execute script statement on target host '%s' error at line %d: %s, statement='%s'",
						currentThreadContext.currentTargetHost, lineTokens.lineNo+1, err.Error(), replacedStatement), ""
				}
			} else { //end if funcName="execute"
				tc.callback(tc.scriptExecID, CRT_STATEMENT_FAILED, msg)
				currentThreadContext.executeFailed = true
				// tc.callback(tc.scriptExecID, CRT_SCRIPT_FAILED, "script run with failed")
				return EXEC_FAIL, fmt.Errorf("execute script statement on target host '%s' error at line %d: %s, statement='%s'",
					currentThreadContext.currentTargetHost, lineTokens.lineNo+1, err.Error(), replacedStatement), ""
			}
		} else {
			if funcName == "execute" {
				tc.callback(tc.scriptExecID, CRT_COMMAND_COMPLETED, result)
				if currentThreadContext.executeSuccessfulFlag == EFS_STOP {
					// tc.callback(tc.scriptExecID, CRT_STATEMENT_FAILED, "")
					return EXEC_SUCCESS, nil, "script run complate with exec successful flag=stop statements"
				}
				if currentThreadContext.executeSuccessfulFlag == EFS_SKIP {
					currentThreadContext.executeSkipNext = true
					currentThreadContext.executeSuccessfulFlag = EFS_CONTINUE
				}
			} else if funcName == "if" || funcName == "endif" || funcName == "for" || funcName == "endfor" || funcName == "finally" || funcName == "endfinally" {
				continue
			} else {
				tc.callback(tc.scriptExecID, CRT_STATEMENT_COMPLETED, result)
			}
		}
		if execLine == endline {
			return EXEC_BREAK, nil, ""
		}
	}
	return 0, nil, ""
}

//getBreakExpressIn
func getBreakExpressIn(tc *scriptContext, lineNo int) (block *scriptBlockTokens) {
	blocks := tc.expressBlock
	forBeginNum := -1
	for _, b := range blocks {
		if b.exp == EXP_FOR {
			if b.beginLineNo < lineNo && b.endLineNo > lineNo && b.beginLineNo > forBeginNum {
				forBeginNum = b.beginLineNo
			}
		}
	}
	if forBeginNum != -1 {
		return tc.expressBlock[forBeginNum]
	}
	return nil
}
