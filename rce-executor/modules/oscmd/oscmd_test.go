package oscmd

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	downloadURL      = "http://nginx.org/download/nginx-1.14.0.tar.gz"
	downloadFilename = "/tmp/nginx-1.14.0.tar.gz"
	downloadFileMD5  = "2d856aca3dfe1d32e3c9f8c4cac0cc95"
	downloadErrMsg   = "download file error"

	executableCmd = "Python -m SimpleHTTPServer"
)

var cmdExecutor executor.Executor

func init() {
	//initialize executor
	//cmdExecutor, _ = modules.NewExecutor()
	//if cmdExecutor == nil {
	//	fmt.Print("executor initialize error")
	//}

}

func download() executor.ExecuteResult {
	params := executor.ExecutorCmdParams{
		CmdParamURL:            downloadURL,
		CmdParamOutputFilename: downloadFilename,
	}
	return DownloadFile(cmdExecutor, &params)
}

func TestMakeDir(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamPath: "/tmp/new/",
	}
	Convey("test make dir", t, func() {
		er := MakeDir(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestChangeOwnAndGroup(t *testing.T) {
	// download file first
	downloadRes := download()
	if !downloadRes.Successful {
		t.Error(downloadErrMsg)
	}

	params := executor.ExecutorCmdParams{
		CmdParamOwn:             "root",
		CmdParamGroup:           "root",
		CmdParamFilenamePattern: downloadFilename,
		CmdParamRecursiveChange: true,
	}
	Convey("test change own and group", t, func() {
		er := ChangeOwnAndGroup(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestUnzipFile(t *testing.T) {
	// download file first
	downloadRes := download()
	if !downloadRes.Successful {
		t.Error(downloadErrMsg)
	}

	params := executor.ExecutorCmdParams{
		CmdParamFilename:  downloadFilename,
		CmdParamDirectory: "/tmp/",
	}
	Convey("test unzip file", t, func() {
		er := UnzipFile(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestDownloadFile(t *testing.T) {
	Convey("test download file", t, func() {
		er := download()
		So(er.Successful, ShouldBeTrue)
	})
}

func TestMD5Sum(t *testing.T) {
	// download file first
	downloadRes := download()
	if !downloadRes.Successful {
		t.Error("download file error")
	}

	params := executor.ExecutorCmdParams{
		CmdParamFilename: downloadFilename,
	}
	Convey("test MD5 sum", t, func() {
		er := MD5Sum(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
		So(er.ResultData[ResultDataKeyMD5], ShouldEqual, downloadFileMD5)
	})
}

func TestChangeMode(t *testing.T) {
	// download file first
	downloadRes := download()
	if !downloadRes.Successful {
		t.Error("download file error")
	}

	params := executor.ExecutorCmdParams{
		CmdParamModeExp:         "u=rwx,g=rwx,o=rwx",
		CmdParamFilenamePattern: downloadFilename,
		CmdParamRecursiveChange: false,
	}
	Convey("test change mode", t, func() {
		er := ChangeMode(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestMove(t *testing.T) {
	// download file first
	downloadRes := download()
	if !downloadRes.Successful {
		t.Error("download file error")
	}

	params := executor.ExecutorCmdParams{
		CmdParamSource: downloadFilename,
		CmdParamTarget: fmt.Sprintf("%s.changed", downloadFilename),
	}
	Convey("test move", t, func() {
		er := Move(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestCopy(t *testing.T) {
	// download file first
	downloadRes := download()
	if !downloadRes.Successful {
		t.Error("download file error")
	}

	params := executor.ExecutorCmdParams{
		CmdParamSource: downloadFilename,
		CmdParamTarget: fmt.Sprintf("%s.copied", downloadFilename),
	}
	Convey("test copy", t, func() {
		er := Copy(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestRemove(t *testing.T) {
	// download file first
	downloadRes := download()
	if !downloadRes.Successful {
		t.Error("download file error")
	}

	params := executor.ExecutorCmdParams{
		CmdParamFilenamePattern: downloadFilename,
		CmdParamRecursiveRemove: false,
	}
	Convey("test remove", t, func() {
		er := Remove(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func nohup() executor.ExecuteResult {
	params := executor.ExecutorCmdParams{
		CmdParamExecutable: "python -m SimpleHTTPServer 9900",
		CmdParamLogFile:    "/tmp/python-simple-server.log",
		CmdParamPIDFile:    "/tmp/python-simple-server.pid",
	}
	return Nohup(cmdExecutor, &params)
}

func processStatus() executor.ExecuteResult {
	params := executor.ExecutorCmdParams{
		CmdParamProcessName: "SimpleHTTPServer",
	}
	return ProcessStatus(cmdExecutor, &params)
}

func killProcessByPID(pid string) executor.ExecuteResult {
	params := executor.ExecutorCmdParams{
		CmdParamPID:       pid,
		CmdParamForceKill: true,
	}
	return KillProcessByPID(cmdExecutor, &params)
}

func TestNohup(t *testing.T) {
	er := nohup()
	Convey("test nohup", t, func() {
		So(er.Successful, ShouldBeTrue)
	})
	killProcessByPID(er.ResultData[ResultDataKeyPID])
}

func TestProcessStatus(t *testing.T) {
	// nohup first
	nohupRes := nohup()
	if !nohupRes.Successful {
		t.Error("nohup error")
	}

	params := executor.ExecutorCmdParams{
		CmdParamProcessName: "SimpleHTTPServer",
	}
	Convey("test process status", t, func() {
		er := ProcessStatus(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
		So(er.ResultData[ResultDataKeyCommand], ShouldContainSubstring, "SimpleHTTPServer")
		So(er.ResultData[ResultDataKeyPID], ShouldEqual, nohupRes.ResultData[ResultDataKeyPID])
	})
	killProcessByPID(nohupRes.ResultData[ResultDataKeyPID])
}

func TestKillProcessByPID(t *testing.T) {
	// nohup first
	nohupRes := nohup()
	if !nohupRes.Successful {
		t.Error("nohup error")
	}

	pid := nohupRes.ResultData[ResultDataKeyPID]
	params := executor.ExecutorCmdParams{
		CmdParamPID:       pid,
		CmdParamForceKill: true,
	}
	Convey("test kill process by pid", t, func() {
		er := KillProcessByPID(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestDisableLinuxSE(t *testing.T) {
	cmdExecutor, err := executor.NewSSHAgentExecutor("192.168.0.121", "root", "root123", 22)
	if err != nil {
		log.Debug("NewSSHAgentExecutor err %v", err)
		return
	}
	if cmdExecutor == nil {
		t.Log("executor is nil ")
		return
	}
	result := DisableSELinux(cmdExecutor, nil)
	t.Log("result = ", result)
}

func TestAddSSHAuthorizedKeys(t *testing.T) {
	params := executor.ExecutorCmdParams{}
	Convey("test add SSH authorized keys", t, func() {
		er := AddSSHAuthorizedKeys(cmdExecutor, &params)
		t.Log(er.Message)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestTextToFile(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamFilename:  "/root/file",
		CmdParamOutText:   "word1",
		CmdParamOverwrite: false,
	}
	Convey("test write text to file", t, func() {
		er := TextToFile(cmdExecutor, &params)
		t.Log(er.Message)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestGetCPUInformation(t *testing.T) {
	params := executor.ExecutorCmdParams{}
	Convey("test get cpu information", t, func() {
		er := GetCPUInformation(cmdExecutor, &params)
		t.Log(er.ResultData)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestDisplayFileSystem(t *testing.T) {
	params := executor.ExecutorCmdParams{}
	Convey("test display file system", t, func() {
		er := DisplayFileSystem(cmdExecutor, &params)
		t.Log(er.ResultData)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestYumPackageInstall(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamSoftNames: "socat libaio libaio-devel perl-Time-HiRes perl-DBD-MySQL perl-Digest-MD5",
	}
	Convey("test install software", t, func() {
		er := YumPackageInstall(cmdExecutor, &params)
		t.Log(er.ResultData)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestBackupYumRepos(t *testing.T) {
	params := executor.ExecutorCmdParams{}
	Convey("test backup yum repos", t, func() {
		er := BackupYumRepos(cmdExecutor, &params)
		t.Log(er.Message)
		So(er.Successful, ShouldBeTrue)
	})
}

//func TestGetSShPath(t *testing.T) {
//	tests := []struct {
//		name     string
//		wantPath string
//		wantErr  bool
//	}{
//		{name:"TestGetSShPath",wantPath:"/MyData/zmysql"},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotPath, err := GetSShPath()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetSShPath() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if gotPath != tt.wantPath {
//				t.Errorf("GetSShPath() = %v, want %v", gotPath, tt.wantPath)
//			}
//		})
//	}
//}

func TestUnInstallSoftware(t *testing.T) {

	ex, err := executor.NewSSHAgentExecutor("192.168.11.167", "root", "123456", 22)
	if err != nil {
		t.Errorf("new executor failed :%s", err.Error())
	}
	params := executor.ExecutorCmdParams{
		CmdParamProcessName: "node_exporter",
		CmdParamPath:        "/home/zmysql/product/node_exporter1",
		CmdParamOSVersion:   "7",
	}
	type args struct {
		e      executor.Executor
		params *executor.ExecutorCmdParams
	}
	tests := []struct {
		name string
		args args
		want executor.ExecuteResult
	}{
		{args: args{e: ex, params: &params}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnInstallSoftware(tt.args.e, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnInstallSoftware() = %v, want %v", got, tt.want)
			}
		})
	}
}
