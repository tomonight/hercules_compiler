package proxysql

// proxysql admin commands

import (
	"database/sql"

	"github.com/juju/errors"
)

// define commands
const (
	CmdProxyReadOnly               = `PROXYSQL READONLY`
	CmdProxyReadWrite              = `PROXYSQL READWRITE`
	CmdProxyStart                  = `PROXYSQL START`
	CmdProxyRestart                = `PROXYSQL RESTART`
	CmdProxyStop                   = `PROXYSQL STOP`
	CmdProxyPause                  = `PROXYSQL PAUSE`
	CmdProxyResume                 = `PROXYSQL RESUME`
	CmdProxyShutdown               = `PROXYSQL SHUTDOWN`
	CmdProxyFlushLogs              = `PROXYSQL FLUSH LOGS`
	CmdProxyKill                   = `PROXYSQL KILL`
	CmdLoadUserToRuntime           = `LOAD MYSQL USERS TO RUNTIME`
	CmdSaveUserToDisk              = `SAVE MYSQL USERS TO DISK`
	CmdLoadServerToRuntime         = `LOAD MYSQL SERVERS TO RUNTIME`
	CmdSaveServerToDisk            = `SAVE MYSQL SERVERS TO DISK`
	CmdLoadQueryRulesToRuntime     = `LOAD MYSQL QUERY RULES TO RUNTIME`
	CmdSaveQueryRulesToDisk        = `SAVE MYSQL QUERY RULES TO DISK`
	CmdLoadSchedulerToRuntime      = `LOAD SCHEDULER TO RUNTIME`
	CmdSaveSchedulerToDisk         = `SAVE SCHEDULER TO DISK`
	CmdLoadMySQLVariablesToRuntime = `LOAD MYSQL VARIABLES TO RUNTIME`
	CmdSaveMySQLVariablesToDisk    = `SAVE MYSQL VARIABLES TO DISK`
	CmdLoadAdminVariablesToRuntime = `LOAD ADMIN VARIABLES TO RUNTIME`
	CmdSaveAdminVariablesToDisk    = `SAVE ADMIN VARIABLES TO DISK`
)

//ProxyReadOnly set proxysql to readonly mode.
func ProxyReadOnly(db *sql.DB) error {
	_, err := db.Exec(CmdProxyReadOnly)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//ProxyReadWrite set proxysql to readwrite mode.
func ProxyReadWrite(db *sql.DB) error {
	_, err := db.Exec(CmdProxyReadWrite)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//ProxyStart start proxysql child process.
func ProxyStart(db *sql.DB) error {
	_, err := db.Exec(CmdProxyStart)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//ProxyRestart restart proxysql process.
func ProxyRestart(db *sql.DB) error {
	_, err := db.Exec(CmdProxyRestart)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//ProxyStop stop proxysql child process.
func ProxyStop(db *sql.DB) error {
	_, err := db.Exec(CmdProxyStop)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//ProxyPause pause proxysql
func ProxyPause(db *sql.DB) error {
	_, err := db.Exec(CmdProxyPause)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//ProxyResume resume proxysql
func ProxyResume(db *sql.DB) error {
	_, err := db.Exec(CmdProxyResume)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//ProxyShutdown shutdown proxysql
func ProxyShutdown(db *sql.DB) error {
	_, err := db.Exec(CmdProxyShutdown)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//ProxyFlushLogs flush proxysql logs to file
func ProxyFlushLogs(db *sql.DB) error {
	_, err := db.Exec(CmdProxyFlushLogs)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//ProxyKill kill child process.
func ProxyKill(db *sql.DB) error {
	_, err := db.Exec(CmdProxyKill)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//LoadUserToRuntime execute load mysql users to runtime.
func LoadUserToRuntime(db *sql.DB) error {
	_, err := db.Exec(CmdLoadUserToRuntime)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//SaveUserToDisk execute save mysql users to disk.
func SaveUserToDisk(db *sql.DB) error {
	_, err := db.Exec(CmdSaveUserToDisk)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//LoadServerToRuntime execute load mysql servers to runtime.
func LoadServerToRuntime(db *sql.DB) error {
	_, err := db.Exec(CmdLoadServerToRuntime)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//SaveServerToDisk execute save mysql servers to disk.
func SaveServerToDisk(db *sql.DB) error {
	_, err := db.Exec(CmdSaveServerToDisk)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//LoadQueryRulesToRuntime execute load mysql query rules to runtime.
func LoadQueryRulesToRuntime(db *sql.DB) error {
	_, err := db.Exec(CmdLoadQueryRulesToRuntime)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//SaveQueryRulesToDisk execute save mysql query rules to disk.
func SaveQueryRulesToDisk(db *sql.DB) error {
	_, err := db.Exec(CmdSaveQueryRulesToDisk)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//LoadSchedulerToRuntime execute load schedulers to runtime.
func LoadSchedulerToRuntime(db *sql.DB) error {
	_, err := db.Exec(CmdLoadSchedulerToRuntime)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//SaveSchedulerToDisk execute save schedulers to disk.
func SaveSchedulerToDisk(db *sql.DB) error {
	_, err := db.Exec(CmdSaveSchedulerToDisk)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//LoadMySQLVariablesToRuntime execute  load mysql variables to runtime.
func LoadMySQLVariablesToRuntime(db *sql.DB) error {
	_, err := db.Exec(CmdLoadMySQLVariablesToRuntime)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//LoadAdminVariablesToRuntime execute load admin variables to runtime.
func LoadAdminVariablesToRuntime(db *sql.DB) error {
	_, err := db.Exec(CmdLoadAdminVariablesToRuntime)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//SaveMySQLVariablesToDisk execute save mysql variables to runtime.
func SaveMySQLVariablesToDisk(db *sql.DB) error {
	_, err := db.Exec(CmdSaveMySQLVariablesToDisk)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

//SaveAdminVariablesToDisk execute save admin variables to disk.
func SaveAdminVariablesToDisk(db *sql.DB) error {
	_, err := db.Exec(CmdSaveAdminVariablesToDisk)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}
