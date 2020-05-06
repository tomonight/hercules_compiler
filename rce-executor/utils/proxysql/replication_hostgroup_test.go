package proxysql

import (
	"flag"
	"testing"
)

func TestQueryAllRHG(t *testing.T) {

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

	allrhg, err := QueryAllRHG(db, 1, 0)
	if err != nil {
		t.Error(allrhg, err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestAddOneRHG(t *testing.T) {

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

	newrhg, err := NewRHG(0, 1)
	if err != nil {
		t.Error(err)
	}

	err = newrhg.AddOneRHG(db)
	if err != nil {
		t.Error(err)
	}

	newrhg.SetWriterHostGroup(0)
	newrhg.SetReaderHostGroup(2)
	newrhg.SetComment("rhg2")

	err = newrhg.UpdateOneRHG(db)
	if err != nil {
		t.Error(err)
	}

	err = newrhg.DeleteOneRHG(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
