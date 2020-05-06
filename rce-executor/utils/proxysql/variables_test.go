package proxysql

import (
	"flag"
	"testing"
)

func TestGetConfigs(t *testing.T) {

	flag.Parse()
	conn, err := NewConn(* proxysqlAddr, *proxysqlPort, *proxysqlUser, *proxysqlPass)
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

	allusers, err := GetConfig(db)
	if err != nil {
		t.Error(allusers, err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneConfigs(t *testing.T) {

	flag.Parse()
	conn, err := NewConn(* proxysqlAddr, *proxysqlPort, *proxysqlUser, *proxysqlPass)
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

	err = UpdateOneConfig(db, "mysql-max_connections", "99999")
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
