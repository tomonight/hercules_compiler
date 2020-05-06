package proxysql

import (
	"flag"
	"testing"
)

func TestFindAllSchedulers(t *testing.T) {

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

	allservers, err := FindAllSchedulerInfo(db, 1, 0)
	if err != nil {
		t.Error(allservers, err)
	}
	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}

}

func TestAddOneSchedulers(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	err = sched.AddOneScheduler(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOneSchedulersActiveEnable(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	sched.Active = 1

	err = sched.UpdateOneSchedulerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOneSchedulersActiveDisable(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	sched.Active = 0

	err = sched.UpdateOneSchedulerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneSchedulersIntervalMs(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	sched.SetSchedulerIntervalMs(3000)

	err = sched.UpdateOneSchedulerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOneSchedulersArg1(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	sched.SetSchedulerArg1("-a")

	err = sched.UpdateOneSchedulerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneSchedulersArg2(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	sched.SetSchedulerArg2("-b")

	err = sched.UpdateOneSchedulerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateOneSchedulersArg3(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	sched.SetSchedulerArg3("-c")

	err = sched.UpdateOneSchedulerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOneSchedulersArg4(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	sched.SetSchedulerArg4("-d")

	err = sched.UpdateOneSchedulerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOneSchedulersArg5(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	sched.SetSchedulerArg5("-e")

	err = sched.UpdateOneSchedulerInfo(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteOneSchedulers(t *testing.T) {

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

	sched, err := NewSch("/bin/bash", 1000)
	if err != nil {
		t.Error(err)
	}

	err = sched.DeleteOneScheduler(db)
	if err != nil {
		t.Error(err)
	}

	err = conn.CloseConn(db)
	if err != nil {
		t.Error(err)
	}
}
