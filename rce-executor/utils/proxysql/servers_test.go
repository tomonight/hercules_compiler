package proxysql

import (
	"flag"
	"testing"
)

func TestFindAllServers(t *testing.T) {

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

	allservers, err := FindAllServerInfo(db, 1, 0)
	if err != nil {
		t.Error(allservers, err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestAddOneServer(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)

	err = newsrv.AddOneServers(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneServerStatusToOnline(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerStatus("ONLINE")

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneServerStatusToOfflineSoft(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerStatus("OFFLINE_SOFT")

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneServerStatusToOfflineHard(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerStatus("OFFLINE_HARD")

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneServerWeight(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerWeight(1000)

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneServerCompressionEnable(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerCompression(1)

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
func TestUpdateOneServerCompressionDisable(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerCompression(0)

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
func TestUpdateOneServerMaxConnection(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerMaxConnection(9999)

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneServerMaxReplication(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerMaxReplicationLag(1000)

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneServerUseSslEnable(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerUseSSL(1)

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
func TestUpdateOneServerUseSslDisable(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerUseSSL(0)

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
func TestUpdateOneServerMaxLatencyMs(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServerMaxLatencyMs(3000)

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateOneServerComment(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)
	newsrv.SetServersComment("test hostgroup")

	err = newsrv.UpdateOneServerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestDeleteOneServer(t *testing.T) {

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

	newsrv, err := NewServer(1, "192.168.100.111", 6032)

	err = newsrv.DeleteOneServers(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
