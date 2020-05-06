package mysql

import (
	"bou.ke/monkey"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"reflect"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var cmdExecutor executor.Executor

func init() {
	//initialize executor
	cmdExecutor, _ = executor.NewLocalExecutor()
	if cmdExecutor == nil {
		fmt.Print("executor initialize error")
	}
}

func TestRegistOrchestrator(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamHost:      "192.168.0.181",
		CmdParamPort:      "3000",
		CmdParamHostName:  "tsql01",
		CmdParamMySQLPort: "3306",
	}

	result := RegistOrchestrator(cmdExecutor, &params)
	t.Log("result=", result)
}

func TestStopMysqlInstance(t *testing.T) {
	ex, err := executor.NewSSHAgentExecutor("192.168.11.221", "root", "root123", 22)
	if err != nil {
		t.Errorf("new executor failed :%s", err.Error())
	}
	params := executor.ExecutorCmdParams{
		CmdParamHost:          "192.168.11.221",
		CmdParamMysqlPath:     "/home/zmysql/db/186userdb/186userdb01/mysql",
		CmdParamMySQLPort:     "13306",
		CmdParamMysqlUser:     "root",
		CmdParamMysqlPassword: "Root@123",
	}

	result := StopMysqlInstance(ex, &params)

	if result.Successful == false {
		t.Error(result.Message)
	} else {
		t.Log("ok")
	}
}

func TestMySQLClusterHasMaster(t *testing.T) {
	Convey("Given a valid cluster info", t, func() {
		params := executor.ExecutorCmdParams{
			CmdParamMemberHost:     "192.168.11.223,192.168.11.221,192.168.11.222",
			CmdParamMemberPort:     "9901,9901,9901",
			CmdParamMemberUser:     "root,root,root",
			CmdParamMemberPassword: "123456,123456,123456",
		}

		ex, err := executor.NewLocalExecutor()
		So(err, ShouldBeNil)
		So(ex, ShouldNotBeNil)
		Convey("When start get mysql cluster has master", func() {
			//mock all need to connect database functions
			monkey.Patch(getMySQLVersion, func(info *InstanceInfo) (version string, err error) {
				t.Log("do getMySQLVersion for mock")
				if info != nil {
					version = fmt.Sprintf("%d", 5)
				}
				return
			})
			monkey.Patch(existMasterInfo, func(info *InstanceInfo, version string) (masterExist bool, err error) {
				t.Log("do existMasterInfo for mock")
				if info != nil {
					if strings.Contains(version, "5") {
						masterExist = true
					} else {
						masterExist = true
					}
				}
				return
			})
			defer monkey.UnpatchAll()

			result := ClusterHasMaster(ex, &params)
			t.Log("result = ", result)
			Convey("Than get the mysql cluster has master successful", func() {
				So(result.Successful, ShouldBeTrue)
			})

		})

		Convey("When start get mysql cluster has no master", func() {
			//mock all need to connect database functions
			monkey.Patch(getMySQLVersion, func(info *InstanceInfo) (version string, err error) {
				t.Log("do getMySQLVersion for mock")
				if info != nil {
					version = fmt.Sprintf("%d", 5)
				}
				return
			})
			monkey.Patch(existMasterInfo, func(info *InstanceInfo, version string) (masterExist bool, err error) {
				t.Log("do existMasterInfo for mock")
				if info != nil {
					if strings.Contains(version, "5") {
						masterExist = false
					} else {
						masterExist = false
					}
				}
				return
			})
			defer monkey.UnpatchAll()
			result := ClusterHasMaster(ex, &params)
			t.Log("result = ", result)
			Convey("Than get the mysql cluster has no master successful", func() {
				So(result.Successful, ShouldBeTrue)
			})
		})
	})
}

func TestSelectMySQLClusterMasterInfo(t *testing.T) {
	Convey("Given a valid cluster info", t, func() {
		params := executor.ExecutorCmdParams{
			CmdParamMemberHost:     "192.168.11.223,192.168.11.221,192.168.11.222",
			CmdParamMemberPort:     "9901,9901,9901",
			CmdParamMemberUser:     "root,root,root",
			CmdParamMemberPassword: "123456,123456,123456",
		}

		ex, err := executor.NewLocalExecutor()
		So(err, ShouldBeNil)
		So(ex, ShouldNotBeNil)
		Convey("When start get mysql cluster master", func() {

			//mock all need to connect database functions
			monkey.Patch(getGTID, func(info *InstanceInfo) (gitID string, err error) {
				t.Log("do getGTID for mock")
				if info != nil {
					gitID = fmt.Sprintf("%d", 5)
				}
				return
			})
			monkey.Patch(getMasterByGTID, func(info *InstanceInfo, firstGTID, secondGTID string) (isMaster bool, err error) {
				t.Log("do getMasterByGTID for mock")
				if info != nil {
					if firstGTID == secondGTID {
						isMaster = true
					}
				}
				return
			})
			defer monkey.UnpatchAll()

			result := SelectMySQLClusterMasterInfo(ex, &params)
			t.Log("result = ", result)
			Convey("Than get the mysql cluster  master successful", func() {
				So(result.Successful, ShouldBeTrue)
			})
		})
	})
}

func TestMySQLInstanceAlive(t *testing.T) {
	params := executor.ExecutorCmdParams{}
	params[CmdParamHost] = "192.168.11.33"
	params[CmdParamPort] = 3306
	params[CmdParamMysqlUser] = "root"
	params[CmdParamMysqlPassword] = "123456"
	type args struct {
		e      executor.Executor
		params *executor.ExecutorCmdParams
	}
	tests := []struct {
		name string
		args args
		want executor.ExecuteResult
	}{
		{name: "success ", args: args{params: &params}, want: executor.SuccessulExecuteResultNoData("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MySQLInstanceAlive(tt.args.e, tt.args.params); !reflect.DeepEqual(got.Successful, tt.want.Successful) {
				t.Errorf("MySQLInstanceAlive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckMySQLCommunityClientExist(t *testing.T) {
	ex, err := executor.NewSSHAgentExecutor("192.168.11.221", "root", "root123", 22)
	if err != nil {
		t.Errorf("new executor failed :%s", err.Error())
	}
	ex1, err := executor.NewSSHAgentExecutor("192.168.11.168", "root", "root123", 22)
	if err != nil {
		t.Errorf("new executor failed :%s", err.Error())
	}
	type args struct {
		e   executor.Executor
		in1 *executor.ExecutorCmdParams
	}
	tests := []struct {
		name string
		args args
		want string
	}{

		{name: "has", args: args{e: ex}, want: "MySQL Community Client Exist"},
		{name: "no has", args: args{e: ex1}, want: "MySQL Community Client No Exist"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckMySQLCommunityClientExist(tt.args.e, tt.args.in1); !reflect.DeepEqual(got.Message, tt.want) {
				t.Errorf("CheckMySQLCommunityClientExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOneSchemaTables(t *testing.T) {
	host := "192.168.11.169"
	port := 3306
	username := "root"
	password := "Root@123"
	schema := "mydata_service_test"

	tables, err := getOneSchemaTables(host, username, password, schema, port)
	t.Log("tables = ", tables)
	t.Log("err = ", err)
}

func Test_renameSchema(t *testing.T) {
	host := "192.168.11.169"
	port := 3306
	username := "root"
	password := "Root@123"
	schema := "mydata_service_test"
	err := renameOneSchema(schema, "mydata_service_test_old", host, username, password, []string{}, port, true)
	t.Log("err = ", err)
}

func TestMySQLInstanceAliveEx(t *testing.T) {
	host := "192.168.11.221"
	username := "root"
	password := "root123"
	port := 22
	e, err := executor.NewSSHAgentExecutor(host, username, password, port)
	if err != nil {
		t.Error("executor.NewSSHAgentExecutor failed %v", err)
	}
	params := executor.ExecutorCmdParams{
		"mysqlPath":  "/home/zmysql/db/mytest/mytest01/mysql",
		"socketFile": "/home/zmysql/db/mytest/mytest01/run/mysql.sock",
	}

	er := MySQLInstanceAlive(e, &params)
	t.Log("er = ", er)
}

func TestTestSelectMySQLClusterMasterInfoExt(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamMemberHost:     "192.168.66.137,192.168.66.135,192.168.66.136",
		CmdParamMemberPort:     "8722,8722,8722",
		CmdParamMemberUser:     "zcloud_platform,zcloud_platform,zcloud_platform",
		CmdParamMemberPassword: "xF7UDA62LZQ19jl_,xF7UDA62LZQ19jl_,xF7UDA62LZQ19jl_",
	}

	ex, err := executor.NewLocalExecutor()
	if err != nil {
		t.Errorf("NewLocalExecutor failed %v", err)
		return
	}

	result := SelectMySQLClusterMasterInfo(ex, &params)
	t.Log("result = ", result)
}

///home/zmysql/db/cat/cat01/conf/my.cnf

func TestInterruptDatabaseRemoteService(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamMysqlParameterFile: "/home/zmysql/db/cat/cat01/conf/my.cnf",
	}

	ex, err := executor.NewSSHAgentExecutor("192.168.11.222", "root", "root123", 22)
	if err != nil {
		t.Errorf("NewLocalExecutor failed %v", err)
		return
	}

	result := InterruptDatabaseRemoteService(ex, &params)
	t.Log("result = ", result)

	result = RecoverDatabaseRemoteService(ex, &params)
	t.Log("result = ", result)
}

func TestShutDownMySQLInstance(t *testing.T) {
	params := executor.ExecutorCmdParams{
		"dbName":       "single_recover_test",
		"instanceName": "single_recover_test01",
	}

	ex, err := executor.NewSSHAgentExecutor("192.168.11.221", "root", "root123", 22)
	if err != nil {
		t.Errorf("NewLocalExecutor failed %v", err)
		return
	}

	result := ShutDownMySQLInstance(ex, &params)
	t.Log("result = ", result)
}

func TestStartupMySQLInstance(t *testing.T) {
	params := executor.ExecutorCmdParams{
		"dbName":       "single_recover_test",
		"instanceName": "single_recover_test01",
	}

	ex, err := executor.NewSSHAgentExecutor("192.168.11.221", "root", "root123", 22)
	if err != nil {
		t.Errorf("NewLocalExecutor failed %v", err)
		return
	}

	result := StartMySQLService(ex, &params)
	t.Log("result = ", result)
}