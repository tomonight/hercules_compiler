package proxysql

import (
	"flag"
	"testing"
)

func TestFindAllQr(t *testing.T) {

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

	allqr, err := FindAllQr(db, 1, 0)
	if err != nil {
		t.Error(allqr, err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}
func TestAddOneQr(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)

	err = newqr.AddOneQr(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOneQrActive(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrActive(1)

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrDisactive(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrActive(0)

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrApply(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrApply(1)

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrNotapply(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrApply(0)

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrSchemaname(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrSchemaName("devschema")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrFlagIN(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrFlagIN(1)

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrClientAddr(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrClientAddr("192.168.200.101")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrProxyAddr(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrProxyAddr("192.168.200.1")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrProxyPort(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetProxyPort("9999")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrDigest(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrDigest("0x0D34125A1267E1FE")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrMatchDigest(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrMatchDigest("^SELECT")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrMatchPattern(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrMatchPattern("^UPDATE")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrNegateMatchPatternEnable(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrNegateMatchPattern(1)

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOneQrNegateMatchPatternDisable(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrNegateMatchPattern(0)

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrFlagOut(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrFlagOut("2")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrReplacePattern(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrReplacePattern("^UPDATE3")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrDestHostgroup(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrDestHostGroup("100")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrCacheTTL(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrCacheTTL("5000")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrReconnectEnable(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrReconnect("1")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrReconnectDisable(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrReconnect("0")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrTimeOut(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrTimeOut("600")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrRetries(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrRetries("3")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrDelay(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrDelay("10")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrMirrorFlagOut(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrMirrorFlagOUT("1")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrMirrorHostgroup(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrMirrorHostgroup("100")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrErrorMsg(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrErrorMsg("fuck errors.")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrLogEnable(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrLog("1")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneQrLogDisable(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}
	newqr.SetQrRuleid(1)
	newqr.SetQrLog("0")

	err = newqr.UpdateOneQrInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteOneQr(t *testing.T) {

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

	newqr, err := NewQr("dev")
	if err != nil {
		t.Error(err)
	}

	newqr.SetQrRuleid(1)
	err = newqr.DeleteOneQr(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
