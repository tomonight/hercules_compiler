package zdata

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	//	"hercules_compiler/rce-executor/modules/oscmd"
	"strings"
)

const (
	CmdParamYumSourceAddr = "yumSourceAddr" //yum源文件地址
)

//SetYumSource 设置yum 源为本地
func SetYumSource(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	log.Debug("%s", "start set zdata local yum source")
	if e != nil {

		sourceFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamYumSourceAddr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		var sourceExist bool
		//		//@remark 获取系统发行版本号
		//		iVer, err := oscmd.GetLinuxDistVer(params)
		//		if err != nil {
		//			return executor.ErrorExecuteResult(err)
		//		}
		log.Debug("%s", "start backup local.repo")
		localYumConf := "/etc/yum.repos.d/local.repo"
		localYumConfBak := "/etc/yum.repos.d/local.repo.bak"
		cmdStr := fmt.Sprintf("ls %s", localYumConf)
		es, err := e.ExecShell(cmdStr)
		if err != nil {
			log.Warn("%s executor failed %v", cmdStr, err)
			return executor.ErrorExecuteResult(err)
		}
		err = executor.GetExecResult(es)
		if err != nil {
			log.Debug("%s executor failed %v", cmdStr, err)

			//start create this file
			cmdStr = "touch " + localYumConf
			es, err = e.ExecShell(cmdStr)
			if err != nil {
				log.Warn("%s executor failed %v", cmdStr, err)
				return executor.ErrorExecuteResult(err)
			}
			err = executor.GetExecResult(es)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}

		} else {
			if len(es.Stdout) > 0 {
				if localYumConf != strings.TrimSpace(es.Stdout[0]) {
					return executor.ErrorExecuteResult(errors.New(localYumConf + "not exist"))
				}
				sourceExist = true
			}
		}

		backExist := false
		if sourceExist {
			cmdStr = fmt.Sprintf("ls %s", localYumConfBak)
			es, err = e.ExecShell(cmdStr)
			if err == nil {
				err = executor.GetExecResult(es)
				if err == nil {
					if len(es.Stdout) > 0 {
						if localYumConfBak == strings.TrimSpace(es.Stdout[0]) {
							backExist = true
						}
					}
				}
			}
		}

		if !backExist && sourceExist {
			cmdStr = fmt.Sprintf("cp %s %s", localYumConf, localYumConfBak)
			es, err = e.ExecShell(cmdStr)
			if err != nil {
				log.Warn("%s copy cmdstr executor failed: %v", cmdStr, err)
				return executor.ErrorExecuteResult(err)
			}
			err = executor.GetExecResult(es)
			if err != nil {
				log.Warn("%s copy cmdstr executor failed: %v", cmdStr, err)
				return executor.ErrorExecuteResult(err)
			}
		}
		log.Debug("%s", "local.repo backup success !")
		log.Debug("%s", "start mkdir /mnt/cdrom")
		mntPath := "/mnt/cdrom"
		cmdStr = fmt.Sprintf("ls %s", mntPath)
		es, err = e.ExecShell(cmdStr)
		mntExist := false
		if err == nil {
			err = executor.GetExecResult(es)
			if err == nil {
				mntExist = true
			}
		}
		if !mntExist {
			//fmt.Println("cmdStr ", cmdStr, " failed ", err)
			//该目录不存在创建该目录
			cmdStr = fmt.Sprintf("mkdir -p %s", mntPath)
			es, err = e.ExecShell(cmdStr)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
			err = executor.GetExecResult(es)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
		}

		log.Debug("%s", "create /mnt/cdrom success start  set localyum")
		localYumContent := "[localyum]\n"
		localYumContent += "name=localyum\n"
		localYumContent += "baseurl=file:///mnt/cdrom\n"
		localYumContent += "gpgcheck=0"

		//echo 'add content' > /etc/yum.repos.d/local.repo 写入内容
		cmdStr = fmt.Sprintf("echo '%s' > %s", localYumContent, localYumConf)
		es, err = e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		err = executor.GetExecResult(es)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		log.Debug("%s", "check local yum is mount")
		checkCmd := "df -h | grep /mnt/cdrom | wc -l"
		es, err = e.ExecShell(checkCmd)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		err = executor.GetExecResult(es)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if len(es.Stdout) > 0 {
			output := es.Stdout[0]
			if strings.Index(output, "1") != -1 { //yum源已经挂载成功
				return executor.SuccessulExecuteResult(es, false, "Mount successful")
			}
		}
		log.Debug("cdrom un mount ,start mount ")
		var cmd string
		//		if iVer >= 7 {
		//cmd = "mount -o loop /soft/rhel-server-7.2-x86_64-dvd.iso /mnt/cdrom"
		cmd = fmt.Sprintf("mount -o loop %s /mnt/cdrom", sourceFile)
		//		} else {
		//			cmd = "mount -o loop /soft/rhel-server-6.8-x86_64-dvd.iso /mnt/cdrom"
		//		}

		es, err = e.ExecShell(cmd)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		err = executor.GetExecResult(es)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		log.Debug("%s", "mount cdrom success")
		log.Debug("%s", "start move unse repo to /etc/")

		mvUnuseRepo := "ls /etc/yum.repos.d|grep -v \"local.repo\"|xargs -i mv /etc/yum.repos.d/{} /etc/"

		es, err = e.ExecShell(mvUnuseRepo)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		err = executor.GetExecResult(es)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		er := executor.SuccessulExecuteResult(es, true, "Yum source set successful")
		return er

	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}
