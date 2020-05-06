package proxysql

import (
	"flag"
	"testing"
)

func TestFindAllUsers(t *testing.T) {

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

	allusers, err := FindAllUserInfo(db, 1, 0)
	if err != nil {
		t.Error(allusers, err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestAddOneUser(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetUserActive(1)

	err = newuser.AddOneUser(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneUserFastForwardEnable(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetFastForward(1)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneUserFastForwardDisable(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetFastForward(0)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneUserMaxConnections(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetMaxConnections(999)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneUserActive(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetUserActive(1)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneUserDisactive(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetUserActive(0)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneUserUseSslEnable(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetUseSSL(1)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
func TestUpdateOneUserUseSslDisable(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetUseSSL(0)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneUserShcemaLockedEnable(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetSchemaLocked(1)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
func TestUpdateOneUserShcemaLockedDisable(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetSchemaLocked(0)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneUserTransactionPersistentEnable(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetTransactionPersistent(1)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
func TestUpdateOneUserTransactionPersistentDisable(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}
	newuser.SetTransactionPersistent(0)

	err = newuser.UpdateOneUserInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestDeleteOneUser(t *testing.T) {

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

	newuser, err := NewUser("devtest", "devtest", 0, "dev")
	if err != nil {
		t.Error(err)
	}

	err = newuser.DeleteOneUser(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
