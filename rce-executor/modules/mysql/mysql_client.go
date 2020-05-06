package mysql

import (
	"fmt"
	"hercules_compiler/rce-executor/utils"
	"regexp"

	"github.com/astaxie/beego/orm"

	"hercules_compiler/rce-executor/log"
	"strings"
)

//ExecuteProxySQLCommand 执行ProxySQL命令
func ExecuteProxySQLCommand(ip, username, password, sql string, port int) error {
	if ip == "" || username == "" || password == "" || port <= 0 {
		return fmt.Errorf("输入ProxySQL服务信息不合法")
	}
	if sql == "" {
		return fmt.Errorf("SQL语句不能为空")
	}
	log.Debug("execute sql =%s", sql)
	sqls := strings.Split(sql, ";")
	log.Debug("execute sqls =%v and len=%d taget host=%s and port=%d", sqls, len(sqls), ip, port)
	mysqlClient := utils.NewMySQLClient(ip, username, password, utils.MySQLClientForProxSQL, uint(port))
	err := mysqlClient.Connect()
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

	for _, curSQL := range sqls {
		if strings.TrimSpace(curSQL) == "" {
			continue
		}
		log.Debug("execute curSQL %s", curSQL)
		err = mysqlClient.DoExecuteSQL(curSQL)
		if err != nil {
			return err
		}
	}
	return err
}

//ExecuteMySQLCommand 执行MySQL命令
func ExecuteMySQLCommand(ip, username, password, sql string, port int) error {
	if ip == "" || username == "" || password == "" || port <= 0 {
		return fmt.Errorf("输入MySQL服务信息不合法")
	}
	if sql == "" {
		return fmt.Errorf("SQL语句不能为空")
	}
	log.Debug("execute sql =%s", sql)
	sqls := strings.Split(sql, ";")
	log.Debug("execute sqls =%v and len=%d", sqls, len(sqls))
	mysqlClient := utils.NewMySQLClient(ip, username, password, utils.MySQLClientForSource, uint(port))
	err := mysqlClient.Connect()
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

	for _, curSQL := range sqls {
		if strings.TrimSpace(curSQL) == "" {
			continue
		}
		log.Debug("execute curSQL = %s", curSQL)
		err = mysqlClient.DoExecuteSQL(curSQL)
		if err != nil {
			err = fmt.Errorf("execute sql:%s failed %v ", curSQL, err)
			return err
		}
	}
	return err
}

//ExecuteMySQLCommandWithResult 执行MySQL命令，并返回值
func ExecuteMySQLCommandWithResult(ip, username, password, sql string, port int) (string, error) {
	if ip == "" || username == "" || password == "" || port <= 0 {
		return "", fmt.Errorf("输入MySQL服务信息不合法")
	}
	if sql == "" {
		return "", fmt.Errorf("SQL语句不能为空")
	}
	log.Debug("execute sql =%s", sql)
	sqls := strings.Split(sql, ";")
	log.Debug("execute sqls =%v and len=%d", sqls, len(sqls))
	mysqlClient := utils.NewMySQLClient(ip, username, password, utils.MySQLClientForSource, uint(port))
	err := mysqlClient.Connect()
	if err != nil {
		return "", err
	}
	defer mysqlClient.Close()

	if strings.TrimSpace(sql) == "" {
		return "", nil
	}
	log.Debug("execute curSQL = %s", sql)
	return mysqlClient.ExcuteSQLWithSingleResult(sql)
}

//ExecuteMySQLPing 查看mysql是否能连接
func ExecuteMySQLPing(ip, username, password string, port int) error {
	if ip == "" || username == "" || password == "" || port <= 0 {
		return fmt.Errorf("输入MySQL服务信息不合法")
	}
	mysqlClient := utils.NewMySQLClient(ip, username, password, utils.MySQLClientForSource, uint(port))
	err := mysqlClient.Connect()
	if err != nil {
		return err
	}
	defer mysqlClient.Close()
	return mysqlClient.Ping()
}

//ExecutePorxySQLCommandQuery 执行ProxySQL查询命令
func ExecutePorxySQLCommandQuery(ip, username, password, sql string, port int, res interface{}) error {
	if ip == "" || username == "" || password == "" || port <= 0 {
		return fmt.Errorf("输入ProxySQL服务信息不合法")
	}
	if sql == "" {
		return fmt.Errorf("SQL语句不能为空")
	}
	log.Debug("Query SQL = %s", sql)
	mysqlClient := utils.NewMySQLClient(ip, username, password, utils.MySQLClientForProxSQL, uint(port))
	err := mysqlClient.Connect()
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

	sqlDB, err := mysqlClient.GetSQLDB()
	if err != nil {
		return err
	}
	aliasName := fmt.Sprintf("%s:%d", ip, port)
	o, err := orm.NewOrmWithDB("mysql", aliasName, sqlDB)
	if err != nil {
		return err
	}
	err = o.Raw(sql).QueryRow(res)
	return err
}

//ExecuteMySQLCommandQuery 执行MySQL查询命令
func ExecuteMySQLCommandQuery(ip, username, password, sql string, port int, mutil bool, res interface{}) error {
	if ip == "" || username == "" || password == "" || port <= 0 {
		return fmt.Errorf("输入MySQL服务信息不合法")
	}
	if sql == "" {
		return fmt.Errorf("SQL语句不能为空")
	}
	log.Debug("Query SQL = %s", sql)
	mysqlClient := utils.NewMySQLClient(ip, username, password, utils.MySQLClientForSource, uint(port))
	err := mysqlClient.Connect()
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

	sqlDB, err := mysqlClient.GetSQLDB()
	if err != nil {
		return err
	}
	aliasName := fmt.Sprintf("%s:%d", ip, port)
	o, err := orm.NewOrmWithDB("mysql", aliasName, sqlDB)
	if err != nil {
		return err
	}
	if mutil {
		_, err = o.Raw(sql).QueryRows(res)
	} else {
		err = o.Raw(sql).QueryRow(res)
	}

	return err
}

//KillThreads kill实例进程
func KillThreads(ip, username, password string, port int) error {
	var ids []int
	sql := `SELECT Id FROM information_schema.PROCESSLIST WHERE Command != 'binlog dump' AND User != 'system user' AND Id != CONNECTION_ID() and user !='orchestrator'
and id in (select trx_mysql_thread_id from INFORMATION_SCHEMA.INNODB_TRX where trx_mysql_thread_id is not null)`
	err := ExecuteMySQLCommandQuery(ip, username, password, sql, port, true, &ids)
	log.Debug("ExecuteMySQLCommandQuery err = %v ids=%v", err, ids)
	if err != nil {
		return err
	}
	killSQLList := []string{}
	for _, id := range ids {
		killSQL := fmt.Sprintf("KILL %d", id)
		killSQLList = append(killSQLList, killSQL)
	}
	if len(killSQLList) > 0 {
		sql = strings.Join(killSQLList, ";")
		return ExecuteMySQLCommand(ip, username, password, sql, port)
	}
	return nil

}

func getMySQLVersion(instance *InstanceInfo) (version string, err error) {
	//+-----------+
	//| @@version |
	//	+-----------+
	//| 8.0.11    |
	//	+-----------+
	sql := "select @@version;"
	err = ExecuteMySQLCommandQuery(instance.Host, instance.Username, instance.Password, sql, instance.Port, false, &version)
	if err != nil && err != orm.ErrNoRows {
		return
	}
	if version == "" || err == orm.ErrNoRows {
		err = fmt.Errorf("实例(%s:%d)未查询到MySQL版本信息", instance.Host, instance.Port)
	}
	return
}

func getGTID(instance *InstanceInfo) (gitID string, err error) {
	//	sql := "select @@gtid_executed;"
	sql := "select @@GLOBAL.gtid_executed;"
	err = ExecuteMySQLCommandQuery(instance.Host, instance.Username, instance.Password, sql, instance.Port, false, &gitID)
	if err != nil && err != orm.ErrNoRows {
		return
	}
	if gitID == "" || err == orm.ErrNoRows {
		err = fmt.Errorf("实例(%s:%d)不存在gtid", instance.Host, instance.Port)
	}
	return
}

func getMasterByGTID(instance *InstanceInfo, firstGTID, secondGTID string) (isMaster bool, err error) {
	sql := "select GTID_SUBSET('%s','%s') as is_master;"
	doSQL := fmt.Sprintf(sql, firstGTID, secondGTID)
	err = ExecuteMySQLCommandQuery(instance.Host, instance.Username, instance.Password, doSQL, instance.Port, false, &isMaster)
	return
}

func getMaster(instanceList []*InstanceInfo) (master *InstanceInfo, err error) {
	count := len(instanceList)
	if count > 0 {
		lastMasterInfo := instanceList[0]
		if count > 1 {
			for index := 1; index < count; index++ {
				compInstanceInfo := instanceList[index]
				lastMasterInfo, err = whoIsMaster(compInstanceInfo, lastMasterInfo)
				if err != nil {
					return
				}
			}
		}
		master = lastMasterInfo
		return
	}
	err = fmt.Errorf("输入的实例列表为空")
	return
}

//whoIsMaster Who is the primary node
func whoIsMaster(first, second *InstanceInfo) (master *InstanceInfo, err error) {
	if first == nil || second == nil {
		err = fmt.Errorf("输入实例信息为空")
		return
	}

	firstGTID, err := getGTID(first)
	if err != nil {
		return
	}

	secondGTID, err := getGTID(second)
	if err != nil {
		return
	}

	isMaster, err := getMasterByGTID(first, secondGTID, firstGTID)
	if err != nil {
		return
	}

	if isMaster {
		master = first
	} else {
		isMaster, err = getMasterByGTID(first, firstGTID, secondGTID)
		if err != nil {
			return
		}
		if isMaster {
			master = second
		}
	}
	return
}

//existMasterInfo check instance exist master info or not
func existMasterInfo(instance *InstanceInfo, version string) (masterExist bool, err error) {
	//version regexp
	version5 := regexp.MustCompile("^5")
	version8 := regexp.MustCompile("^8")
	masterInfo := struct {
		VariableName string `orm:"column(Variable_name)"`
		Value        string `orm:"column(Value)"`
		MemberHost   string `orm:"column(MEMBER_HOST)"`
		MemberPort   string `orm:"column(MEMBER_PORT)"`
	}{}

	var sql string
	if version5.MatchString(version) {
		sql = "show global status like 'group_replication_primary_member';"
	} else if version8.MatchString(version) {
		sql = "select MEMBER_HOST, MEMBER_PORT from performance_schema.replication_group_members where MEMBER_ROLE='PRIMARY';"
	} else {
		err = fmt.Errorf("不支持的MySQL版本%s", version)
		return
	}

	err = ExecuteMySQLCommandQuery(instance.Host, instance.Username, instance.Password, sql, instance.Port, false, &masterInfo)
	if err != nil && err != orm.ErrNoRows {
		return
	}
	err = nil
	if masterInfo.MemberHost != "" || (masterInfo.Value != "" && masterInfo.Value != "UNDEFINED") {
		masterExist = true
	}
	return
}
