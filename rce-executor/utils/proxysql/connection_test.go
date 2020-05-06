package proxysql

import (
	"flag"
	"testing"
)

var proxysqlAddr = flag.String("addr", "127.0.0.1", "proxysql listen address.default 127.0.0.1")
var proxysqlPort = flag.Uint64("port", 6032, "proxysql listen port,default 6032")
var proxysqlUser = flag.String("user", "admin", "proxysql administrator name.default admin")
var proxysqlPass = flag.String("pass", "admin", "proxysql administrator password.default admin")

func TestNewConn(t *testing.T) {

	flag.Parse()
	conn, err := NewConn(*proxysqlAddr, *proxysqlPort, *proxysqlUser, *proxysqlPass)
	if err != nil {
		t.Error(conn, err)
	}

	conn.SetCharset("utf8")
	conn.SetCollation("utf8_general_ci")
	conn.MakeDBI()

	db, err := conn.OpenConn()
	if err != nil {
		t.Error(db, err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
