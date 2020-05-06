package zdata

import (
	"errors"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/oscmd"
	"hercules_compiler/rce-executor/modules/osconfig"
	"testing"
	"time"
)

var e executor.Executor
var err error

func init() {
	e, err = executor.NewSSHAgentExecutor(
		"192.168.0.204",
		"root",
		"root123",
		22)
	if err != nil {
		log.Error("connect target failed %v", err)
		e = nil
	}
}

func setRunLevel() error {
	params := executor.ExecutorCmdParams{
		oscmd.CmdParamRunLevel: 3}

	result := oscmd.SetRunLevel(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}

	result = oscmd.IsolateRunLevel(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}

	return nil
}

func setHosts() error {
	hostName := "rac1"
	params := executor.ExecutorCmdParams{
		osconfig.CmdParamHostname: hostName}

	result := osconfig.SetHostname(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	params[osconfig.CmdParamIPAddr] = "192.168.0.201"
	result = osconfig.SetHostsFile(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

func setYum() error {
	params := executor.ExecutorCmdParams{
		CmdParamYumSourceAddr: "http://wwww.baidu.com"}
	result := SetYumSource(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

func installRPM() error {
	publicRpmPackages := "binutils ksh gcc gcc-c++ glibc "
	publicRpmPackages += "glibc-common glibc-devel glibc-headers libaio libaio-devel libgcc "
	publicRpmPackages += "libstdc++ libstdc++-devel make sysstat compat-libcap1 cpp mpfr libXp "
	publicRpmPackages += "xorg-x11-utils xorg-x11-xauth smartmontools ethtool bind-utils pam "
	publicRpmPackages += "numactl libaio mutt lm_sensors lvm2 "
	publicRpmPackages += "tuna ksh tigervnc pciutils unzip parted mlocate ntp gtk2 atk "
	publicRpmPackages += "cairo elfutils-libelf-devel lrzsz patch kernel-devel perl-devel lsscsi "
	publicRpmPackages += "tk gcc-gfortran openssh-clients dstat zlib-devel "
	publicRpmPackages += "keyutils-libs e2fsprogs-devel libsepol-devel libselinux-devel krb5-devel "
	publicRpmPackages += "openssl-devel tcl-devel rpm-build redhat-rpm-config device-mapper-multipath "
	publicRpmPackages += "libtool bison flex glib2-devel tree sg3_utils iotop ncurses-devel ftp "
	publicRpmPackages += "sysfsutils python-devel gcc-c++ tcl gcc-gfortran libxml2-python bc tcsh tk "
	publicRpmPackages += "mlocate nfs-utils lsof bash libgudev1-devel "
	publicRpmPackages += "telnet net-tools setuptool perf unixODBC-devel unixODBC pyparted python-six "
	publicRpmPackages += "wget perl-ExtUtils-MakeMaker sos psmisc "
	publicRpmPackages += "iptraf-ng lvm2-python-libs bash-completion"

	params := executor.ExecutorCmdParams{
		oscmd.CmdParamSoftNames: publicRpmPackages}

	result := oscmd.YumPackageInstall(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	for _, value := range result.ResultData {
		log.Debug("install log:", value)
	}
	return nil
}

// change_zone
func changeZone() error {
	zone := "Asia/Shanghai"
	params := executor.ExecutorCmdParams{
		osconfig.CmdParamTimezone: zone}
	result := osconfig.SetTimezone(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

// settime
func setTime() error {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	log.Debug("set time %s", timeStr)
	params := executor.ExecutorCmdParams{
		osconfig.CmdParamDateTime: timeStr}
	result := osconfig.SetDateTime(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

//disableFireWall
func disableFirewall() error {
	params := executor.ExecutorCmdParams{
		oscmd.CmdParamVersion: "7.0"}

	result := oscmd.DisableFireWall(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

//disableSeLinux
func diableSelinux() error {
	params := executor.ExecutorCmdParams{
		oscmd.CmdParamVersion: "7.0"}

	result := oscmd.DisableSELinux(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

//disableService
func disableService() error {
	params := executor.ExecutorCmdParams{
		CmdParamIsMonitor: true}

	result := DisableMonitorServices(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

//
func changeLvmConf() error {
	//params := executor.ExecutorCmdParams{
	//	oscmd.CmdParamVersion: "7.0"}
	//ChangeLimitsConf()
	return nil
}

func installStorage() error {
	params := executor.ExecutorCmdParams{}
	result := InitStorageNode(e, &params)
	log.Debug("result = %v", result)
	log.Debug("executeHasError: %v", result.ExecuteHasError)
	log.Debug("exit code = %d", result.ExitCode)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

//RemoveLinuxModule
func removeLinuxModule() error {
	params := executor.ExecutorCmdParams{
		CmdParamsLinuxModName: "daiwei",
	}
	result := RmMpd(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

//generate client conf
func gernateClientConf() error {
	params := executor.ExecutorCmdParams{
		CmdParamEtcdIpList:   "192.168.0.201, 192.168.0.202,192.168.0.203",
		CmdParamEtcdPortList: "23233,32445,54566",
	}
	result := GenerateClientConf(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

//generate client conf
func gernateAgentConf() error {
	params := executor.ExecutorCmdParams{
		CmdParamMonitorIp:   "192.168.0.201",
		CmdParamMonitorPort: "23233",
		CmdParamNodeType:    1,
	}
	result := GenerateAgentConf(e, &params)
	if result.ExecuteHasError {
		return errors.New(result.Message)
	}
	return nil
}

//TestInstallZdata 安装zdata
func TestInstallZdata(t *testing.T) {
	//check zdata

	//	//设置runlevel
	//	err = setRunLevel()
	//	if err != nil {
	//		t.Error("setRunLevel failed ", err)
	//		return
	//	}
	//	log.Debug("runlelve done")
	//	//set hosts
	//	err = setHosts()
	//	if err != nil {
	//		t.Error("setHosts failed ", err)
	//		return
	//	}
	//	log.Debug("sethosts done")
	//	//setlocal yum
	//	err = setYum()
	//	if err != nil {
	//		t.Error("setYum failed ", err)
	//		return
	//	}
	//	log.Debug("set local yum done")

	//	//installRPM
	//	err = installRPM()
	//	if err != nil {
	//		t.Error("installRPM failed ", err)
	//		return
	//	}
	//	log.Debug("install rpm done")

	//	//changeZone
	//	err = changeZone()
	//	if err != nil {
	//		t.Error("changeZone failed ", err)
	//		return
	//	}
	//	log.Debug("changeZone done")

	//	//setTime
	//	err = setTime()
	//	if err != nil {
	//		t.Error("setTime failed ", err)
	//		return
	//	}
	//	log.Debug("setTime done")

	//	//disableFirewall
	//	err = disableFirewall()
	//	if err != nil {
	//		t.Error("disableFirewall failed ", err)
	//		return
	//	}
	//	log.Debug("disable firewall done")
	//	//diableSelinux
	//	err = diableSelinux()
	//	if err != nil {
	//		t.Error("diableSelinux failed ", err)
	//		return
	//	}
	//	log.Debug("disable selinux done")

	//	//disableService
	//	err = disableService()
	//	if err != nil {
	//		t.Error("disable monitor service failed ", err)
	//		return
	//	}
	//	log.Debug("disable monitor service done")

	//	err := installStorage()
	//	if err != nil {
	//		t.Error("installStorage failed ", err)
	//		return
	//	}
	//	log.Debug("installStorage done")

	//	err := removeLinuxModule()

	//	if err != nil {
	//		t.Error("removeLinuxModule failed ", err)
	//		return
	//	}
	//	log.Debug("removeLinuxModule done")

	err := gernateClientConf()
	if err != nil {
		t.Error("gernateClientConf failed ", err)
		return
	}
	log.Debug("gernateClientConf done")

	err = gernateAgentConf()
	if err != nil {
		t.Error("gernateAgentConf")
		return
	}
	log.Debug("gernateAgentConf done")
}
