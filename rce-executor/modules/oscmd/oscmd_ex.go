package oscmd

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	ResTupleJoinFlag = "*"
)

//文件单位
const (
	//SizeTypeB byte
	SizeTypeB = 1
	//SizeTypeKB kb
	SizeTypeKB = 2
	//SizeTypeMB mb
	SizeTypeMB = 3
)

//default info
const (
	DefaultRecvPort = 12808
)

// 函数名常量定义
const (
	CmdNameLsbRelease             = "LsbRelease"
	CmdNameCat                    = "Cat"
	CmdNameList                   = "List"
	CmdNameSed                    = "Sed"
	CmdNameLinuxDist              = "LinuxDist"
	CmdNameYumPackageInstall      = "YumPackageInstall"
	CmdNameDisableFirewall        = "DisableFireWall"
	CmdNameDisableSelinux         = "DisableSELinux"
	CmdNameGetMemorySize          = "GetMemorySize"
	CmdNameReplaceFileText        = "ReplaceFileText"
	CmdNameTouch                  = "Touch"
	CmdNameSplitTextGetIndexValue = "SplitTextGetIndexValue"
	CmdNameGetCPUInformation      = "GetCPUInformation"
	CmdNameReadFile               = "ReadFile"
	CmdNamePsEfInformation        = "GetPsInfoByPid"
	CmdNameYumRemove              = "YumRemove"
	CmdNameBackupYumRepos         = "BackupYumRepos"
	CmdNameDirectoryCanDo         = "DirectoryCanDo"
	CmdNamePortValid              = "PortValid"
	CmdNameTransferFileToTarget   = "TransferFileToTarget"
	CmdNameOSVersionAvailable     = "OSVersionAvailable"
)

// 命令参数常量定义
const (
	CmdParamBasic          = "basic"
	CmdParamVersion        = "version"
	CmdParamID             = "id"
	CmdParamDesc           = "description"
	CmdParamRelease        = "release"
	CmdParamCode           = "code"
	CmdParamAll            = "all"
	CmdParamHelp           = "help"
	CmdParamText           = "text"
	CmdParamFlag           = "flag"
	CmdParamIndex          = "index"
	CmdParamSoftNames      = "softNames"
	CmdParamRepoName       = "repoName"
	CmdParamSubString      = "subString"
	CmdParamReplaceString  = "replaceString"
	CmdParamPackageNames   = "packageNames"
	CmdParamPackageURLs    = "packageURLs"
	CmdParamIntallPackaget = "intallPath"
	CmdParamPort           = "port"
	CmdParamOSVersion      = "osversion"
)

// 结果集键定义
const (
	ResultDataKeyBasic             = "basic"
	ResultDataKeyPaths             = "paths"
	ResultDataKeyVersion           = "version"
	ResultDataKeyText              = "text"
	ResultDataKeyIndexValue        = "indexValue"
	ResultDataKeyDist              = "dist"
	ResultDataKeyPsuedoname        = "psuedoname"
	ResultDataKeyArch              = "architecture"
	ResultDataKeyKernel            = "kernel"
	ResultDataMemorySize           = "memorySize"
	ResultDataCPUProcessorCount    = "cpuProcessorCount"
	ResultDataCPUPhysicalCoreCount = "cpuPhysicalCoreCount"
	ResultDataCPUModelName         = "cpuModelName"
)

// GetLinuxDistVer 获取系统版本系统
func GetLinuxDistVer(params *executor.ExecutorCmdParams) (int, error) {
	iVer := 0
	ver, err := executor.ExtractCmdFuncStringParam(params, CmdParamVersion)
	if err != nil {
		return iVer, err
	}

	verArray := strings.Split(ver, ".")
	if len(verArray) > 0 {
		iVer, err = strconv.Atoi(verArray[0])
		if err != nil {
			return iVer, err
		}
	}
	return iVer, nil
}

// LsbRelease lsb_release 命令
func LsbRelease(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	basicParams := map[string]string{
		CmdParamVersion: "-v",
		CmdParamID:      "-i",
		CmdParamDesc:    "-d",
		CmdParamRelease: "-r",
		CmdParamCode:    "-c",
		CmdParamAll:     "-a",
		CmdParamHelp:    "-h"}

	if params == nil {
		executor.ErrorExecuteResult(errors.New("params is nil"))
	}

	//@remark lsb_release 只接受一个参数
	if len(*params) > 1 {
		executor.ErrorExecuteResult(errors.New("too many params"))
	}

	lsbParam, err := executor.ExtractCmdFuncStringParam(params, CmdParamBasic)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	var (
		paramFlag string
		ok        bool
	)
	paramFlag, ok = basicParams[lsbParam]
	if !ok {
		return executor.ErrorExecuteResult(errors.New("unrecognized params"))
	}

	cmdstr := fmt.Sprintf("%s %s", "lsb_release", paramFlag)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "Linux LSB Release get successul")
		var value string
		if len(es.Stdout) > 0 {
			value = strings.Join(es.Stdout, ResTupleJoinFlag)
		} else {
			if len(es.Stdout) == 1 {
				value = es.Stdout[0]
			}
		}
		resultData := map[string]string{ResultDataKeyVersion: value}
		er.ResultData = resultData
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)

}

// Cat cat 命令
func Cat(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if filePath == "" {
		return executor.ErrorExecuteResult(errors.New("file path is nil"))
	}

	cmdstr := fmt.Sprintf("%s %s", "cat", filePath)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "file "+filePath+" content get successul")
		var value string
		er.ResultData = make(map[string]string)
		for i, v := range es.Stdout {
			er.ResultData[fmt.Sprintf("%d", i)] = v
			value = value + v
		}
		er.ResultData[ResultDataKeyText] = value
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// Touch touch 命令
func Touch(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if filePath == "" {
		return executor.ErrorExecuteResult(errors.New("file path is nil"))
	}

	cmdStr := fmt.Sprintf("ls %s", filePath)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("%s executor failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	if err == nil {
		if len(es.Stdout) > 0 {
			if filePath == strings.TrimSpace(es.Stdout[0]) {
				er := executor.SuccessulExecuteResult(es, false, " file "+filePath+"alreay exist")
				return er
			}
		}
	}

	cmdstr := fmt.Sprintf("%s %s", "touch", filePath)
	log.Debug("touch cmd string %s", cmdstr)
	es, err = e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "create file "+filePath+" successul")
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// List ls 命令
func List(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if filePath == "" {
		return executor.ErrorExecuteResult(errors.New("file path is nil"))
	}

	//@remark 可选参数不做错误判断
	baisc, _ := executor.ExtractCmdFuncStringParam(params, CmdParamBasic)

	var cmdstr string
	if baisc == "" {
		cmdstr = fmt.Sprintf("%s %s", "ls", filePath)
	} else {
		cmdstr = fmt.Sprintf("%s %s %s", "ls", baisc, filePath)
	}

	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "file list for "+filePath+" successful")
		var value string
		if len(es.Stdout) > 0 {
			value = strings.Join(es.Stdout, ResTupleJoinFlag)
		} else {
			if len(es.Stdout) == 1 {
				value = es.Stdout[0]
			}
		}
		resultData := map[string]string{ResultDataKeyPaths: value}
		er.ResultData = resultData
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// Sed sed 命令
func Sed(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	//@remark 基础命令 可选项不做判断
	basic, _ := executor.ExtractCmdFuncStringParam(params, CmdParamBasic)

	filePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if filePath == "" {
		return executor.ErrorExecuteResult(errors.New("file path is nil"))
	}

	text, err := executor.ExtractCmdFuncStringParam(params, CmdParamText)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	var cmdstr string
	if text == "" {
		return executor.ErrorExecuteResult(errors.New("text is nil"))
	}

	if basic == "" {
		cmdstr = fmt.Sprintf("%s %s %s", "sed", text, filePath)
	} else {
		cmdstr = fmt.Sprintf("%s %s %s %s", "sed", basic, text, filePath)
	}

	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "sed on file "+filePath+" successful")
		var value string
		if len(es.Stdout) > 0 {
			value = strings.Join(es.Stdout, ResTupleJoinFlag)
		} else {
			if len(es.Stdout) == 1 {
				value = es.Stdout[0]
			}
		}
		if value != "" {
			resultData := map[string]string{ResultDataKeyBasic: value}
			er.ResultData = resultData
		}

		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// LinuxDist 获取Linux版本发行信息, 只适用linux版本
func LinuxDist(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	var (
		redhat   = "/etc/redhat-release"
		suse     = "/etc/SUSE-release"
		mandrake = "/etc/mandrake-release"
		debian   = "/etc/debian_version"
	)

	var (
		kernel       string
		architecture string
		dist         string
		version      string
		psuedoname   string
	)

	//@remark do redhat check
	//PSUEDONAME=`cat /etc/redhat-release | sed s/.*\(// | sed s/\)//`
	sedParamHead := "sed s/.*\\(//" //替换 ( 的前半部分为空字符串
	sedParamBack := "sed s/\\)//"   //替换后）的半部分为空
	redhatCmd := fmt.Sprintf("cat %s | %s | %s", redhat, sedParamHead, sedParamBack)
	//fmt.Println("cmd =", redhatCmd)
	res, err := e.ExecShell(redhatCmd)
	//fmt.Println("res=", res, " err = ", err)
	if err == nil {
		dist = "RedHat"
		if stdoutLen := len(res.Stdout); stdoutLen > 0 {
			//fmt.Println("stdout=", res.Stdout)
			psuedoname = res.Stdout[0]
		}
		//REV=`cat /etc/redhat-release | sed s/.*release\ // | sed s/\ .*//`
		sedParamHead = "sed s/.*release\\ //" //替换relase 前半部分为空字符串
		sedParamBack = "sed s/\\ .*//"
		redhatCmd = fmt.Sprintf("cat %s | %s | %s", redhat, sedParamHead, sedParamBack)
		res, err = e.ExecShell(redhatCmd)
		if err == nil {
			if stdoutLen := len(res.Stdout); stdoutLen > 0 {
				version = res.Stdout[0]
			}
		} else {
			return executor.ErrorExecuteResult(err)
		}

	} else {
		//@remark do suse check
		//DIST=`cat /etc/SUSE-release | tr "\n" ' '| sed s/VERSION.*//`
		trHead := "tr \"\n\" ' '"         //换行符替换为空格
		sedParamBack = "sed s/VESION.*//" //VERSION后面的全部替换为空字符
		suseCmd := fmt.Sprintf("cat %s | %s | %s", suse, trHead, sedParamBack)
		res, err = e.ExecShell(suseCmd)
		if err == nil {
			if stdoutLen := len(res.Stdout); stdoutLen > 0 {
				dist = res.Stdout[0]
			}

			//REV=`cat /etc/SUSE-release | tr "\n" ' ' | sed s/.*=\ //`
			sedParamBack = "sed s/.*=\\ //"
			suseCmd = fmt.Sprintf("cat %s | %s | %s", suse, trHead, sedParamBack)
			res, err = e.ExecShell(suseCmd)
			if err == nil {

				if stdoutLen := len(res.Stdout); stdoutLen > 0 {
					version = res.Stdout[0]
				}
			} else {
				return executor.ErrorExecuteResult(err)
			}

		} else {
			//@remark do  mandrake check
			//PSUEDONAME=`cat /etc/mandrake-release | sed s/.*\(// | sed s/\)//`
			mandrakeCmd := fmt.Sprintf("cat %s | %s | %s", mandrake, sedParamHead, sedParamBack)
			res, err = e.ExecShell(mandrakeCmd)
			if err == nil {
				dist = "Mandrake"
				if stdoutLen := len(res.Stdout); stdoutLen > 0 {
					psuedoname = res.Stdout[0]
				}

				//REV=`cat /etc/mandrake-release | sed s/.*release\ // | sed s/\ .*//`
				sedParamHead = "s/.*release\\ //" //替换relase 前半部分为空字符串
				sedParamBack = "s/\\ .*//"
				mandrakeCmd = fmt.Sprintf("cat %s | %s | %s", mandrake, sedParamHead, sedParamBack)
				res, err = e.ExecShell(mandrakeCmd)
				if err == nil {
					if stdoutLen := len(res.Stdout); stdoutLen > 0 {
						version = res.Stdout[0]
					}
				} else {
					return executor.ErrorExecuteResult(err)
				}

			} else {
				//@remark do debian check
				//DIST = "Debian `cat /etc/debian_version`"
				debianCmd := fmt.Sprintf("cat %s", debian)
				res, err = e.ExecShell(debianCmd)
				if err == nil {
					if stdoutLen := len(res.Stdout); stdoutLen > 0 {
						dist = "Debian " + res.Stdout[0]
					}
					version = ""
				} else {
					return executor.ErrorExecuteResult(err)
				}

			}
		}
	}

	machCmd := "uname -m"
	res, err = e.ExecShell(machCmd)
	if err == nil {
		if stdoutLen := len(res.Stdout); stdoutLen > 0 {
			architecture = res.Stdout[0]
		}
	} else {
		return executor.ErrorExecuteResult(err)
	}

	kernelCmd := "uname -r"
	res, err = e.ExecShell(kernelCmd)
	if err == nil {
		if stdoutLen := len(res.Stdout); stdoutLen > 0 {
			kernel = res.Stdout[0]
		}
	} else {
		return executor.ErrorExecuteResult(err)
	}

	result := new(executor.ExecuteResult)
	executor.ExecutedStatus2ExecuteResult(result, res)
	result.Changed = false
	result.Successful = true
	result.ExecuteHasError = false
	resultData := make(map[string]string)
	resultData[ResultDataKeyDist] = strings.ToLower(dist)
	resultData[ResultDataKeyVersion] = strings.ToLower(version)
	resultData[ResultDataKeyPsuedoname] = strings.ToLower(psuedoname)
	resultData[ResultDataKeyArch] = strings.ToLower(architecture)
	resultData[ResultDataKeyKernel] = strings.ToLower(kernel)
	result.Message = "linux dist get successful"
	result.ResultData = resultData
	return *result
}

// YumPackageInstall 使用`yum install`命令安装软件
func YumPackageInstall(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	softNames, err := executor.ExtractCmdFuncStringParam(params, CmdParamSoftNames)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	repoName, err := executor.ExtractCmdFuncStringParam(params, CmdParamRepoName)
	if repoName == "" {
		//做之前mydata兼容适配
		repoName = "MyData"
	}

	var errMsg string
	var es *executor.ExecutedStatus
	tries := 1
	for tries <= 2 {
		cmdstr := fmt.Sprintf("yum install --disablerepo=* --enablerepo=%s -y %s", repoName, softNames)
		log.Debug(cmdstr)
		es, err = e.ExecShell(cmdstr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 {
			return executor.SuccessulExecuteResult(es, true, fmt.Sprintf("rpm package '%s' installed", softNames))
		}

		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
			break
		} else {
			errMsg = es.Stderr[0]
			// 遇到 Repodata is over 2 weeks old. Install yum-cron? Or run: yum makecache fast 错误
			// 重新执行一次安装
			if strings.Index(errMsg, "Repodata is over") >= 0 && strings.Index(errMsg, "old") >= 0 && tries < 2 {
				log.Debug("Encountered", errMsg, ", retring...")
				cmdstr = "yum makecache fast"
				es, err = e.ExecShell(cmdstr)
				if err != nil {
					return executor.ErrorExecuteResult(err)
				}
			}
		}
		tries++
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//定义yum 命令可忽略的错误
func yumCommandIgnoreText() []string {
	ignoreTexts := []string{
		"不删除任何软件包",
		"无须任何处理",
		"Nothing to do",
		"altered outside of yum",
	}
	return ignoreTexts
}

//判断yum 执行结果中的错误是否可忽略
func yumErrorCanIgnore(result string) bool {
	ignoreTexts := yumCommandIgnoreText()
	for _, value := range ignoreTexts {
		if strings.Contains(result, value) {
			return true
		}
	}
	return false
}

// YumRemove 使用`yum remove`命令安装软件
func YumRemove(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	softNames, err := executor.ExtractCmdFuncStringParam(params, CmdParamSoftNames)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//参数并未赋值不执行卸载命令
	if strings.Contains(softNames, "$") {
		return executor.SuccessulExecuteResultNoData("no need to uninstall rpm package ")
	}

	cmdstr := fmt.Sprintf("yum remove -y %s", softNames)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, fmt.Sprintf("rpm package '%s' uninstalled", softNames))
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	if yumErrorCanIgnore(errMsg) {
		return executor.SuccessulExecuteResult(es, true, fmt.Sprintf("rpm package '%s' uninstalled", softNames))
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

func getPackageURL(packageName string, urlAddrList []string) string {
	for _, urlAddr := range urlAddrList {
		if strings.Contains(urlAddr, packageName) {
			return urlAddr
		}
	}
	return ""
}

// YumLocalPackagesInstall 使用yum 命令安装本地目录下的软件包
func YumLocalPackagesInstall(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	packageNames, err := executor.ExtractCmdFuncStringListParam(params, CmdParamPackageNames, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("packageNames = %v and len =%d", packageNames, len(packageNames))

	packageURLs, err := executor.ExtractCmdFuncStringListParam(params, CmdParamPackageURLs, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("packageURLs = %v and len =%d", packageURLs, len(packageURLs))
	//check packageName exist or not
	//if pakcageName not exist ,then start download package
	for _, packageName := range packageNames {
		ext := filepath.Ext(packageName)
		log.Debug("pacakgeName %s ext = %s", packageName, ext)
		if packageName == "" || ext != ".rpm" {
			return executor.ErrorExecuteResult(fmt.Errorf("package %s ext invalid", packageName))
		}
		if exist, _ := FileExist(e, packageName); !exist {
			//start download
			urlAddr := getPackageURL(packageName, packageURLs)
			if urlAddr == "" {
				return executor.ErrorExecuteResult(fmt.Errorf("package %s download url not found", packageName))
			}

			downloadParams := executor.ExecutorCmdParams{
				CmdParamURL:            urlAddr,
				CmdParamOutputFilename: packageName,
			}
			es := DownloadFile(e, &downloadParams)
			if !es.Successful {
				return es
			}
		}
	}
	//yum install --disablerepo=*  -q -y socat-1.7.2.2-5.el7.x86_64.rpm
	cmdstr := fmt.Sprintf("yum install --disablerepo=*  -q -y %s", strings.Join(packageNames, " "))
	log.Debug("command = %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, fmt.Sprintf("rpm package '%s' install", strings.Join(packageNames, " ")))
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	if yumErrorCanIgnore(errMsg) {
		return executor.SuccessulExecuteResult(es, true, fmt.Sprintf("rpm package '%s' install", strings.Join(packageNames, " ")))
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// DisableFireWall 停止防火墙
func DisableFireWall(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {
		var (
			es  *executor.ExecutedStatus
			err error
			ver int
		)

		er := LinuxDist(e, params)
		if !er.ExecuteHasError {
			version, ok := er.ResultData[ResultDataKeyVersion]
			if ok {
				(*params)[CmdParamVersion] = version
			}
		} else {
			return executor.ErrorExecuteResult(errors.New(er.Message))
		}

		ver, err = GetLinuxDistVer(params)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		var cmdList []string
		if ver >= 7 {
			cmdList = []string{
				"systemctl disable firewalld",
				"systemctl stop firewalld"}
		} else {
			cmdList = []string{
				"chkconfig iptables off",
				"/etc/init.d/iptables stop",
				"chkconfig ip6tables off",
				"/etc/init.d/ip6tables stop"}
		}

		for _, cmd := range cmdList {
			es, err = e.ExecShell(cmd)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}

			err = executor.GetExecResult(es)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
		}
		return executor.SuccessulExecuteResult(es, true, "Firewall disabled")

	} else {
		return executor.ErrorExecuteResult(errors.New("executor is nil"))
	}
}

// DisableSELinux 停止SELinux
func DisableSELinux(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {
		var (
			es     *executor.ExecutedStatus
			err    error
			status string
		)
		//@remark 通过getenforce 获取linuxse状态
		log.Debug("%s", "start getenforce")
		cmdStr := "getenforce"
		es, err = e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		err = executor.GetExecResult(es)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		if len(es.Stdout) > 0 {
			status = es.Stdout[0]
			status = strings.ToLower(status)
			log.Debug("Current selinux status %s", status)
			ignoreStatus := []string{"permissive", "disabled"}

			if strings.Contains(strings.Join(ignoreStatus, ""), status) {
				//if status == "permissive" {
				return executor.SuccessulExecuteResult(es, false, "SELinux already disabled")
			}
		} else {
			return executor.ErrorExecuteResult(errors.New(executor.ErrMsgUnknow))
		}

		log.Debug("%s", "start set selinux to permissive")

		cmdStr = "setenforce 0"
		es, err = e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		err = executor.GetExecResult(es)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		log.Debug("%s", "start set selinux forever disable")
		selinuxConfigFile := "/etc/selinux/config"
		cmdStr = "sed -i s/SELINUX=enforcing/SELINUX=disabled/g " + selinuxConfigFile
		log.Debug("cmd=%s", cmdStr)
		es, err = e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		err = executor.GetExecResult(es)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		if len(es.Stdout) > 0 {
			log.Debug("last text:%s", es.Stdout[0])
		} else {
			log.Debug("last stdou:%i", es.Stdout)
		}

		return executor.SuccessulExecuteResult(es, true, "SELinux disabled")

	} else {
		return executor.ErrorExecuteResult(errors.New("executor is nil"))
	}
}

// GetMemorySize 获取目标机器内存大小
func GetMemorySize(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdstr := fmt.Sprintf("%s %s", "cat", "/proc/meminfo |grep MemTotal")
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "get memory size success ")
		var value string
		if len(es.Stdout) == 1 {
			value = es.Stdout[0]
		}
		//MemTotal:         999696 kB
		valueArray := strings.Split(value, ":")
		if len(valueArray) == 2 {
			value = valueArray[1]
			value = strings.Replace(value, "kB", "", -1)
			value = strings.TrimSpace(value)
		}

		resultData := map[string]string{ResultDataMemorySize: value}
		er.ResultData = resultData
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//replaceFileText 替换文件中的文本
func replaceFileText(e executor.Executor, fileName, oldText, newText string) (*executor.ExecutedStatus, error) {
	//sed -i 's/book/books/g' file
	cmdStr := fmt.Sprintf("sed -i s/\"%s\"/\"%s\"/g %s", oldText, newText, fileName)
	log.Debug("cmdStr= %v", cmdStr)
	es, err := e.ExecShell(cmdStr)
	if err == nil {
		err = executor.GetExecResult(es)

	}
	return es, err
}

//ReplaceFileText 替换文件文本 对外接口
func ReplaceFileText(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {
		filePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilename)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		oldText, err := executor.ExtractCmdFuncStringParam(params, CmdParamSubString)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		newText, err := executor.ExtractCmdFuncStringParam(params, CmdParamReplaceString)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		es, err := replaceFileText(e, filePath, oldText, newText)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		err = executor.GetExecResult(es)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		return executor.SuccessulExecuteResult(es, true, "replace susccess")

	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}

//SplitTextGetIndexValue SplitTextGetIndexValue
func SplitTextGetIndexValue(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	text, err := executor.ExtractCmdFuncStringParam(params, CmdParamText)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	flag, err := executor.ExtractCmdFuncStringParam(params, CmdParamFlag)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	index, err := executor.ExtractCmdFuncIntParam(params, CmdParamIndex)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	log.Debug("text=%s flag=%s index=%d", text, flag, index)
	lastText := ""
	if text != "" && index > 0 {
		textArray := []string{}
		if strings.TrimSpace(flag) == "" {
			textArray = strings.Fields(text)
		} else {
			textArray = strings.Split(text, flag)
		}

		log.Debug("textArray=%v", textArray)
		count := len(textArray)
		if index <= count {
			lastText = textArray[index-1]
		}
	}
	log.Debug("lastText=%s", lastText)
	es := executor.ExecutedStatus{}
	es.ExitCode = 0
	er := executor.SuccessulExecuteResult(&es, false, "text split get successul")
	resultData := map[string]string{ResultDataKeyIndexValue: lastText}
	er.ResultData = resultData
	return er
}

// GetCPUInformation 获取CPU信息
func GetCPUInformation(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	processorCountCmd := "cat /proc/cpuinfo| grep 'processor'| wc -l"
	physicalCoreCountCmd := "cat /proc/cpuinfo| grep 'physical id'| sort| uniq| wc -l"
	modelNameCmd := "cat /proc/cpuinfo | grep 'model name' | cut -f2 -d: | uniq -c"
	cmdstr := fmt.Sprintf("%s;%s;%s", processorCountCmd, physicalCoreCountCmd, modelNameCmd)
	log.Debug("cmdstr: ", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "get cpu info success")
		processorCount := es.Stdout[0]
		physicalCoreCount := es.Stdout[1]
		modelName := strings.TrimSpace(es.Stdout[2])
		resultData := map[string]string{
			ResultDataCPUProcessorCount:    processorCount,
			ResultDataCPUPhysicalCoreCount: physicalCoreCount,
			ResultDataCPUModelName:         modelName,
		}
		er.ResultData = resultData
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)

}

// 读取指定目录的文件
// tail 命令封装
func ReadFile(e executor.Executor, params *executor.ExecutorCmdParams) (result executor.ExecuteResult) {
	filePath, err := executor.ExtractCmdFuncStringParam(params, "filePath")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	length, err := executor.ExtractCmdFuncIntParam(params, "length")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	cmdstr := fmt.Sprintf("tail -n %d %s", length, filePath)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "get logs success ")
		resultData := make(map[string]string)
		for k, v := range es.Stdout {
			resultData[strconv.Itoa(k)] = v
		}
		er.ResultData = resultData
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// GetPsInfoByPid 根据pid获取进程参数信息
func GetPsInfoByPid(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	pid, err := executor.ExtractCmdFuncIntParam(params, "pid")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	cmdstr := fmt.Sprintf("ps -p %d -h 2>/dev/null|grep %d", pid, pid)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "get ps info success ")
		var value string
		if len(es.Stdout) == 1 {
			value = es.Stdout[0]
		}

		resultData := map[string]string{"info": value}
		er.ResultData = resultData
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// BackupYumRepos 备份yum源
func BackupYumRepos(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdstr := "mkdir -p /etc/yum.repos.d/backup && mv /etc/yum.repos.d/*.repo /etc/yum.repos.d/backup"
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "backup yum repos success")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
		if strings.Index(errMsg, "No such file or directory") >= 0 {
			return executor.SuccessulExecuteResult(es, false, "no yum repos found, no need backup")
		}
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//DirectoryCanDo 目标主机目录可操作
func DirectoryCanDo(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	directoryExistText := "directory already exist"
	directoryNotExistText := "directory not exist"
	directory, err := executor.ExtractCmdFuncStringParam(params, CmdParamDirectory)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if directory == "" {
		return executor.ErrorExecuteResult(fmt.Errorf("input directory is empty"))
	}
	//多个目录之间以;分隔
	//(test -d a|| test -d b || test -d /home) && echo file exist
	directoryList := strings.Split(directory, ";")
	cmdStr := "("
	dirCount := len(directoryList)
	for index, dir := range directoryList {
		if strings.TrimSpace(dir) == "" { //忽略为空的目录
			continue
		}
		if index+1 == dirCount {
			// cmdStr += "test -d " + dir
			cmdStr += fmt.Sprintf("test -d %s", dir)
		} else {
			// cmdStr += "test -d " + dir + "||"
			cmdStr += fmt.Sprintf("test -d %s ||", dir)
		}
	}
	// cmdStr += ") && echo " + directoryExistText
	cmdStr += fmt.Sprintf(") && echo %s || echo %s", directoryExistText, directoryNotExistText)
	log.Debug("commond = %s", cmdStr)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		if strings.Join(es.Stdout, "") == directoryExistText {
			errText := "目录 " + directory + " 已经存在，不能继续安装"
			return executor.NotSuccessulExecuteResult(es, errText)
		} else {
			return executor.SuccessulExecuteResult(es, true, "directory can do ")
		}
	}
	log.Debug("exit code = %d es.stderr = %v", es.ExitCode, es.Stderr)
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// du -s 命令
func DirUseSpace(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if filePath == "" {
		return executor.ErrorExecuteResult(errors.New("file path is nil"))
	}
	cmdstr := fmt.Sprintf("du -s %s", filePath)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "file "+filePath+" usespace get successul")
		er.ResultData = make(map[string]string)
		value := strings.Split(es.Stdout[0], "\t")
		er.ResultData["space"] = value[0]
		return er
	}
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//PortValid 适用 ss -nap|grep '^tcp'|awk '{print $5}'|grep ':54092'|wc -l 命令 检查端口是否被可用
func PortValid(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	port, err := executor.ExtractCmdFuncStringParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if port == "" {
		return executor.ErrorExecuteResult(errors.New("input port is nil"))
	}
	//ss -nap -Atcp|awk '{print $4}'|grep ':12345$'|wc -l
	// command := fmt.Sprintf("ss -nap|grep '^tcp'|awk '{print $5}'|grep ':%s$'|wc -l", port)
	// ss -nap -Atcp 2>/dev/null |awk '{print $4}'|grep ':33066$'|wc -l
	//command := fmt.Sprintf("ss -nap -Atcp|awk '{print $4}'|grep ':%s$'|wc -l", port)
	//ss -nap -Atcp 2>/dev/null |awk '{print $4}'|grep ':33066$'|wc -l
	// 端口检查方案优化
	//command := fmt.Sprintf("ss -nap -Atcp 2>/dev/null |awk '{print $4}'|grep ':%s$'|wc -l", port)
	command := fmt.Sprintf("ss -na -Atcp 2>/dev/null |awk '{print $4}'|grep ':%s$'|wc -l", port)
	log.Debug("command = %s", command)
	es, err := e.ExecShell(command)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		portCountText := strings.Join(es.Stdout, "")
		log.Debug("port count = %s", portCountText)
		portCount, err := strconv.Atoi(portCountText)
		if err != nil {
			return executor.NotSuccessulExecuteResult(es, err.Error())
		}
		if portCount > 0 {
			return executor.NotSuccessulExecuteResult(es, "port already used")
		}
		er := executor.SuccessulExecuteResult(es, false, "port can do")
		return er
	}
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//GrepLine grep 命令
func GrepLine(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamDirectory)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if filePath == "" {
		return executor.ErrorExecuteResult(errors.New("file path is nil"))
	}
	str, err := executor.ExtractCmdFuncStringParam(params, CmdParamValue)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if str == "" {
		return executor.ErrorExecuteResult(errors.New("search string is nil"))
	}
	cmdstr := fmt.Sprintf("grep \"%s\" %s", str, filePath)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	fmt.Println(es.Stdout)
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "file "+filePath+" usespace get successul")
		er.ResultData = make(map[string]string)
		if len(es.Stdout) > 0 {
			er.ResultData["value"] = strings.Join(es.Stdout, ";")
		} else {
			er.ResultData["value"] = ""
		}
		return er
	}
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//TransferFileToTarget  传输文件到目标主机
func TransferFileToTarget(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	needTransfer, err := executor.ExtractCmdFuncBoolParam(params, CmdParamNeedTransfer)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if !needTransfer {
		return executor.SuccessulExecuteResultNoData("不需要文件传输到远程服务器")
	}

	//get agent type exclude ZCAgent
	agentType, _ := executor.ExtractCmdFuncStringParam(params, CmdParamAgentType)

	pathList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamPath, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	remoteHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	transferHost, _ := executor.ExtractCmdFuncStringParam(params, CmdParamTransferHost)
	if transferHost == "" {
		transferHost = remoteHost
	}

	remotePort, err := executor.ExtractCmdFuncIntParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	remoteUserName, err := executor.ExtractCmdFuncStringParam(params, CmdParamUserName)
	if err != nil && agentType == "" {
		return executor.ErrorExecuteResult(err)
	}

	remotePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamRemotePath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	recvPort, _ := executor.ExtractCmdFuncIntParam(params, CmdParamRemoteRecvPort)
	if recvPort == 0 {
		recvPort = DefaultRecvPort
	}

	for _, path := range pathList {
		err = socatTransferFiles(e, agentType, remoteHost, transferHost, remoteUserName, path, remotePath, uint(remotePort), uint(recvPort))
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
	}
	return executor.SuccessulExecuteResultNoData("文件传输到远程服务器成功")
}

//GrepPidOf pidof 命令
func GrepPidOf(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	name, err := executor.ExtractCmdFuncStringParam(params, CmdParamProcessName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if name == "" {
		return executor.ErrorExecuteResult(errors.New("search name is nil"))
	}
	cmdstr := fmt.Sprintf("pidof %s", name)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "pid get successul")
		er.ResultData = make(map[string]string)
		if len(es.Stdout) > 0 {
			er.ResultData["value"] = es.Stdout[0]
		} else {
			er.ResultData["value"] = ""
		}
		return er
	}
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//GetPortByPID 根据pid获取端口 命令
func GetPortByPID(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	pid, err := executor.ExtractCmdFuncIntParam(params, CmdParamPID)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if pid == 0 {
		return executor.ErrorExecuteResult(errors.New("search pid is nil"))
	}
	// cmdstr := fmt.Sprintf("ss -nltp|grep pid=%d|awk '{print $4}'", pid)
	//ss -nltp|egrep '=14744,|,14744,'|awk '{print $4}'
	cmdstr := fmt.Sprintf("ss -nltp|egrep '=%d,|,%d,'|awk '{print $4}'", pid, pid)
	log.Debug("command = %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("result  = %v", es.Stdout)
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "port get successul")
		er.ResultData = make(map[string]string)
		if len(es.Stdout) > 0 {
			portValue := []string{}
			for _, v := range es.Stdout {
				r, _ := regexp.Compile(":[0-9]+")
				if r.FindString(v) != "" {
					portValue = append(portValue, r.FindString(v)[1:])
				}
			}
			er.ResultData["value"] = strings.Join(portValue, ",")
		} else {
			er.ResultData["value"] = ""
		}
		return er
	}
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//GetMysqldPath 根据pid获取mysqld路径 命令
func GetMysqldPath(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	pid, err := executor.ExtractCmdFuncIntParam(params, CmdParamPID)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if pid == 0 {
		return executor.ErrorExecuteResult(errors.New("search pid is nil"))
	}
	cmdstr := fmt.Sprintf("ls -l /proc/%d/exe|awk '{print $11}'", pid)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "mysqld path get successul")
		er.ResultData = make(map[string]string)
		if len(es.Stdout) > 0 {
			er.ResultData["value"] = es.Stdout[0]
		} else {
			er.ResultData["value"] = ""
		}
		return er
	}
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//GetMysqlServerID 根据conf文件获取mysql的serverid
func GetMysqlServerID(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	path, err := executor.ExtractCmdFuncStringParam(params, "path")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if path == "" {
		return executor.ErrorExecuteResult(errors.New("search conf path is nil"))
	}
	cmdstr := fmt.Sprintf("cat %s | grep -E -i '^server_id|^server-id'", path)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "server id get successul")
		er.ResultData = make(map[string]string)
		str := []string{}
		for _, v := range es.Stdout {
			serverID := strings.Replace(strings.Split(v, "=")[1], " ", "", -1)
			str = append(str, serverID)
		}
		er.ResultData["value"] = strings.Join(str, ",")
		return er
	}
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//GetSlaveIP 获取从库IP
func GetSlaveIP(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	port, err := executor.ExtractCmdFuncStringParam(params, "port")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if port == "" {
		return executor.ErrorExecuteResult(errors.New("search conf path is nil"))
	}
	cmdstr := fmt.Sprintf("ss -antp|awk '{print $5}'|grep '%s$'", port)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "mysqld path get successul")
		er.ResultData = make(map[string]string)
		if len(es.Stdout) > 0 {
			str := es.Stdout[0]
			re, _ := regexp.Compile(`\d+\.\d+\.\d+\.\d+`)
			er.ResultData["value"] = re.FindString(str)
		} else {
			er.ResultData["value"] = ""
		}
		return er
	}
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//OSVersionAvailable 目标主机操作系统是否支持
func OSVersionAvailable(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	command := fmt.Sprintf("uname -a")
	log.Debug("command = %s", command)
	es, err := e.ExecShell(command)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	available := false
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		log.Debug("stdout = %v ", es.Stdout)
		for _, v := range es.Stdout {
			if strings.Contains(v, "el7.x86_64") {
				available = true
				break
			}
			if strings.Contains(v, "el6.x86_64") {
				available = true
				break
			}
		}
	}
	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	log.Debug("available = %t", available)
	if available {
		return executor.SuccessulExecuteResultNoData("os version can do mydata")
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}
