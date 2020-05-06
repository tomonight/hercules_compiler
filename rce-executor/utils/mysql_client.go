package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/utils/proxysql"

	//import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

//定义数据库类型
const (
	MySQLClientForSource  = "mysql"    //mysql数据库
	MySQLClientForProxSQL = "proxysql" //proxysql数据库
	DefaultUserCount      = 1000
)

//ErrMySQLClientNotInit mysql客户端未初始化
var ErrMySQLClientNotInit = errors.New("mysql client not init")

//ErrMySQLClientNotConnected  mysql客户端未连接
var ErrMySQLClientNotConnected = errors.New("mysql client not connected")

//ErrMySQLUnSupportDriver 不支持的客户端驱动
var ErrMySQLUnSupportDriver = errors.New("unsupport driver name")

//ErrMySQLHandlerNil mysql处理器为空
var ErrMySQLHandlerNil = errors.New("mysql handler is nil")

//MemberInfor 定义成员信息
type MemberInfor struct {
	HostGroupID uint64
	HostName    string
	Port        uint64
	Status      string
}

//MySQLClient 定义mysql客户端
type MySQLClient struct {
	username  string
	password  string
	host      string
	port      uint
	connected bool
	init      bool
	lastErr   string
	dbHandler *sql.DB
	driver    string
	clientID  string
}

//NewMySQLClient 新建MySQL客户端
func NewMySQLClient(host, username, password, driver string, port uint) *MySQLClient {
	client := new(MySQLClient)
	client.username = username
	client.password = password
	client.host = host
	client.port = port
	client.connected = false
	client.init = true
	if driver == "" {
		driver = MySQLClientForSource
	}
	client.driver = driver
	return client
}

//Invalid  判断参数是否合法
func (client *MySQLClient) Invalid() bool {
	if client.username == "" || client.password == "" ||
		client.host == "" || client.port == 0 {
		return true
	}
	return false
}

//SetClientID 设置客户端id
func (client *MySQLClient) SetClientID(clientID string) {
	client.clientID = clientID
}

//GetClientID  获取客户端id
func (client *MySQLClient) GetClientID() string {
	return client.clientID
}

func (client *MySQLClient) connectMySQL() error {
	link := fmt.Sprintf("%s:%s@tcp(%s:%d)/", client.username, client.password, client.host, client.port)
	db, err := sql.Open("mysql", link)
	if err != nil {
		return err
	}
	client.connected = true
	client.dbHandler = db
	return nil
}

func (client *MySQLClient) userDatabase(dbName string) error {
	if dbName == "" {
		dbName = "mysql"
	}
	_, err := client.dbHandler.Exec("USE " + dbName)
	return err
}

func (client *MySQLClient) connectProxySQL() error {
	conn, err := proxysql.NewConn(client.host, uint64(client.port), client.username, client.password)
	if err != nil {
		return err
	}

	conn.SetCharset("utf8")
	conn.SetCollation("utf8_general_ci")
	conn.MakeDBI()

	db, err := conn.OpenConn()
	if err != nil {
		fmt.Println("err = ", err)
		return err
	}

	client.connected = true
	client.dbHandler = db
	return nil
}

//Connect 连接mysql
func (client *MySQLClient) Connect() error {
	if client.init {
		if !client.connected {
			switch client.driver {
			case MySQLClientForProxSQL:
				return client.connectProxySQL()
			case MySQLClientForSource:
				return client.connectMySQL()
			default:
				return ErrMySQLUnSupportDriver
			}
		} else {
			return ErrMySQLClientNotConnected
		}
	} else {
		return ErrMySQLClientNotInit
	}
}

//GetSQLDB 获取mysql原始连接信息
func (client *MySQLClient) GetSQLDB() (*sql.DB, error) {
	if client.connected {
		if client.dbHandler != nil {
			return client.dbHandler, nil
		}
		return nil, ErrMySQLClientNotConnected

	}
	return nil, ErrMySQLClientNotInit
}

//Ping 判断连接是否可用
func (client *MySQLClient) Ping() error {
	if !client.connected {
		return ErrMySQLClientNotConnected
	}
	return client.dbHandler.Ping()
}

//Close 关闭连接
func (client *MySQLClient) Close() error {
	if client.connected {
		if client.dbHandler != nil {
			return client.dbHandler.Close()
		}
	}
	return nil
}

//DoExecuteSQL 执行SQL语句
func (client *MySQLClient) DoExecuteSQL(sql string) error {
	if client.connected {
		if client.dbHandler != nil {
			log.Debug("start execute sql %s", sql)
			_, err := client.dbHandler.Exec(sql)
			return err
		}
		return ErrMySQLClientNotConnected
	}
	return ErrMySQLClientNotInit
}

func (client *MySQLClient) getMySQLUserList() ([]string, error) {
	userList := []string{}
	if client.connected {
		if client.dbHandler != nil {
			if err := client.userDatabase(""); err != nil {
				return userList, err
			}

			queryStr := "select user from user"
			stmt, err := client.dbHandler.Prepare(queryStr)
			if err != nil {
				return userList, err
			}
			defer stmt.Close()

			rows, err := stmt.Query()
			if err != nil {
				return userList, err
			}

			for rows.Next() {
				var username string
				err = rows.Scan(&username)
				if err != nil {
					return userList, err
				}
				if username != "" {
					userList = append(userList, username)
				}
			}
		} else {
			return userList, ErrMySQLClientNotConnected
		}
	} else {
		return userList, ErrMySQLClientNotInit
	}

	return userList, nil
}

//ExcuteSQLWithSingleResult ExcuteSQLWithSingleResult
func (client *MySQLClient) ExcuteSQLWithSingleResult(sql string) (string, error) {
	result := ""
	if client.connected {
		if client.dbHandler != nil {
			if err := client.userDatabase(""); err != nil {
				return result, err
			}

			stmt, err := client.dbHandler.Prepare(sql)
			if err != nil {
				return result, err
			}
			defer stmt.Close()

			rows, err := stmt.Query()
			if err != nil {
				return result, err
			}

			for rows.Next() {
				err = rows.Scan(&result)
				if err != nil {
					return result, err
				}
			}
		} else {
			return result, ErrMySQLClientNotConnected
		}
	} else {
		return result, ErrMySQLClientNotInit
	}

	return result, nil
}

func (client *MySQLClient) getProxySQLUserList() ([]string, error) {
	fmt.Println("start do getProxySQLUserList")
	userList := []string{}
	if client.connected {
		if client.dbHandler != nil {
			allusers, err := proxysql.FindAllUserInfo(client.dbHandler, DefaultUserCount, 0)
			if err != nil {
				return userList, err
			}
			for _, user := range allusers {
				fmt.Println("userInfor = ", user)
				userList = append(userList, user.Username)
			}

		} else {
			return userList, ErrMySQLClientNotConnected
		}
	} else {
		return userList, ErrMySQLClientNotInit
	}
	return userList, nil
}

//addProxySQLUser 添加proxysql 用户信息
func (client *MySQLClient) addProxySQLUser(username, password string) error {
	newuser, err := proxysql.NewUser(username, password, 2, "")
	if err != nil {
		return err
	}

	err = newuser.AddOneUser(client.dbHandler)
	if err != nil {
		return err
	}
	return nil
}

//addMySQLUser 添加mysql 用户信息
func (client *MySQLClient) addMySQLUser(username, password string) error {
	return nil
}

//removeProxySQLUser 移除proxysql 用户信息
func (client *MySQLClient) removeProxySQLUser(username string) error {
	fmt.Println("start remove proxysql")
	allusers, err := proxysql.FindAllUserInfo(client.dbHandler, DefaultUserCount, 0)
	if err != nil {
		return err
	}

	for _, user := range allusers {
		if user.Username == username {
			return user.DeleteOneUser(client.dbHandler)
		}
	}
	return nil
}

func (client *MySQLClient) removeMySQLUser(username string) error {
	return nil
}

//AddUser 添加用户
func (client *MySQLClient) AddUser(username, password string) error {
	if client.connected {
		if client.dbHandler != nil {
			switch client.driver {
			case MySQLClientForProxSQL:
				return client.addProxySQLUser(username, password)
			case MySQLClientForSource:
				return client.addMySQLUser(username, password)
			default:
				return ErrMySQLUnSupportDriver
			}
		} else {
			return ErrMySQLHandlerNil
		}
	} else {
		return ErrMySQLClientNotConnected
	}
}

//RemoveUser 移除用户
func (client *MySQLClient) RemoveUser(username string) error {
	if client.connected {
		if client.dbHandler != nil {
			switch client.driver {
			case MySQLClientForProxSQL:
				return client.removeProxySQLUser(username)
			case MySQLClientForSource:
				return client.removeMySQLUser(username)
			default:
				return ErrMySQLUnSupportDriver
			}
		} else {
			return ErrMySQLHandlerNil
		}
	} else {
		return ErrMySQLClientNotConnected
	}
}

//GetUserList 获取用户列表
func (client *MySQLClient) GetUserList() ([]string, error) {
	userList := []string{}
	if client.connected {
		if client.dbHandler != nil {
			switch client.driver {
			case MySQLClientForProxSQL:
				return client.getProxySQLUserList()
			case MySQLClientForSource:
				return client.getMySQLUserList()
			default:
				return userList, ErrMySQLUnSupportDriver
			}
		} else {
			return userList, ErrMySQLHandlerNil
		}
	} else {
		return userList, ErrMySQLClientNotConnected
	}
}

//GetProxySQLMemberInfor 获取Proxysql成员信息
func (client *MySQLClient) GetProxySQLMemberInfor() ([]MemberInfor, error) {
	memberInfors := []MemberInfor{}
	if client.connected {
		if client.dbHandler != nil {
			switch client.driver {
			case MySQLClientForProxSQL:
				memebers, err := proxysql.FindAllRunTimeServerInfo(client.dbHandler, DefaultUserCount, 0)
				if err != nil {
					return memberInfors, err
				}
				for _, memeber := range memebers {
					mInfor := MemberInfor{}
					mInfor.HostGroupID = memeber.HostGroupID
					mInfor.Status = memeber.Status
					mInfor.HostName = memeber.HostName
					mInfor.Port = memeber.Port
					memberInfors = append(memberInfors, mInfor)
				}
				return memberInfors, nil

			default:
				return memberInfors, ErrMySQLUnSupportDriver
			}
		} else {
			return memberInfors, ErrMySQLHandlerNil
		}
	} else {
		return memberInfors, ErrMySQLClientNotConnected
	}
}
