package proxysql

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/juju/errors"
)

//Servers 定义服务
type Servers struct {
	HostGroupID       uint64 `db:"hostgroup_id,omitempty" json:"hostgroup_id"`
	HostGroupName     string `db:"hostgroup_name,omitempty" json:"hostgroup_name"`
	HostName          string `db:"hostname" json:"hostname"`
	Port              uint64 `db:"port" json:"port"`
	Status            string `db:"status" json:"status"`
	Weight            uint64 `db:"weight" json:"weight"`
	Compression       uint64 `db:"compression" json:"compression"`
	MaxConnections    uint64 `db:"max_connections" json:"max_connections"`
	MaxReplicationLag uint64 `db:"max_replication_lag" json:"max_replication_lag"`
	UseSsl            uint64 `db:"use_ssl" json:"use_ssl"`
	MaxLatencyMs      uint64 `db:"max_latency_ms" json:"max_latency_ms"`
	Comment           string `db:"comment" json:"comment"`
}

//定义查询规则
const (
	/*add a new backends.*/
	StmtAddOneServers = `
	INSERT 
	INTO 
		mysql_servers(
			hostgroup_id,
			hostname,
			port,
			max_connections
		) 
	VALUES(%d,%q,%d,%d)`

	/*delete a backend*/
	StmtDeleteOneServers = `
	DELETE 
	FROM 
		mysql_servers 
	WHERE 
		hostgroup_id=%d 
	AND hostname=%q 
	AND port=%d`

	/*update a backends*/
	StmtUpdateOneServer = `
	UPDATE 
		mysql_servers 
	SET 
		status=%q,
		weight=%d,
		compression=%d,
		max_connections=%d,
		max_replication_lag=%d,
		use_ssl=%d,
		max_latency_ms=%d,
		comment=%q 
	WHERE 
		hostgroup_id=%d 
	AND hostname=%q 
	AND port=%d`

	/*list all mysql_servers*/
	StmtFindAllServer = `
	SELECT 
		ifnull(hostgroup_id,0) as hostgroup_id,
		ifnull(hostname,"") as hostname,
		ifnull(port,0) as port,
		ifnull(status,"") as status,
		ifnull(weight,0) as weight,
		ifnull(compression,0) as compression,
		ifnull(max_connections,0) as max_connections,
		ifnull(max_replication_lag,0) as max_replication_lag,
		ifnull(use_ssl,0) as use_ssl,
		ifnull(max_latency_ms,0) as max_latency_ms,
		ifnull(comment,"") as comment 
	FROM 
		mysql_servers 
	LIMIT %d 
	OFFSET %d`

	/*list all runtime_mysql_servers*/
	StmtFindAllRunTimeServer = `
	SELECT 
		ifnull(hostgroup_id,0) as hostgroup_id,
		ifnull(hostname,"") as hostname,
		ifnull(port,0) as port,
		ifnull(status,"") as status,
		ifnull(weight,0) as weight,
		ifnull(compression,0) as compression,
		ifnull(max_connections,0) as max_connections,
		ifnull(max_replication_lag,0) as max_replication_lag,
		ifnull(use_ssl,0) as use_ssl,
		ifnull(max_latency_ms,0) as max_latency_ms,
		ifnull(comment,"") as comment 
	FROM 
		runtime_mysql_servers 
	LIMIT %d 
	OFFSET %d`
)

// FindAllServerInfo list all mysql_servers
func FindAllServerInfo(db *sql.DB, limit uint64, skip uint64) ([]Servers, error) {

	var allserver []Servers

	Query := fmt.Sprintf(StmtFindAllServer, limit, skip)

	rows, err := db.Query(Query)
	if err != nil {
		return []Servers{}, errors.Trace(err)
	}
	defer rows.Close()

	for rows.Next() {

		var tmpserver Servers

		err = rows.Scan(
			&tmpserver.HostGroupID,
			&tmpserver.HostName,
			&tmpserver.Port,
			&tmpserver.Status,
			&tmpserver.Weight,
			&tmpserver.Compression,
			&tmpserver.MaxConnections,
			&tmpserver.MaxReplicationLag,
			&tmpserver.UseSsl,
			&tmpserver.MaxLatencyMs,
			&tmpserver.Comment,
		)

		if err != nil {
			continue
		}

		allserver = append(allserver, tmpserver)
	}

	return allserver, nil
}

// FindAllRunTimeServerInfo list all runtime_mysql_servers
func FindAllRunTimeServerInfo(db *sql.DB, limit uint64, skip uint64) ([]Servers, error) {

	var allserver []Servers
	Query := fmt.Sprintf(StmtFindAllRunTimeServer, limit, skip)

	rows, err := db.Query(Query)
	if err != nil {
		return []Servers{}, errors.Trace(err)
	}
	defer rows.Close()

	for rows.Next() {

		var tmpserver Servers

		err = rows.Scan(
			&tmpserver.HostGroupID,
			&tmpserver.HostName,
			&tmpserver.Port,
			&tmpserver.Status,
			&tmpserver.Weight,
			&tmpserver.Compression,
			&tmpserver.MaxConnections,
			&tmpserver.MaxReplicationLag,
			&tmpserver.UseSsl,
			&tmpserver.MaxLatencyMs,
			&tmpserver.Comment,
		)

		if err != nil {
			continue
		}

		allserver = append(allserver, tmpserver)
	}

	return allserver, nil
}

//NewServer init a new servers.
func NewServer(hostgroupID uint64, hostname string, port uint64) (*Servers, error) {
	newsrv := new(Servers)

	newsrv.HostGroupID = hostgroupID
	newsrv.HostName = hostname
	newsrv.Port = port

	newsrv.Status = "ONLINE"
	newsrv.Weight = 1000
	newsrv.Compression = 0
	newsrv.MaxConnections = 10000
	newsrv.MaxReplicationLag = 0
	newsrv.UseSsl = 0
	newsrv.MaxLatencyMs = 0
	newsrv.Comment = ""

	return newsrv, nil
}

//SetServerStatus set servers status
func (srvs *Servers) SetServerStatus(status string) {
	switch status {
	case "ONLINE":
		srvs.Status = "ONLINE"
	case "SHUNNED":
		srvs.Status = "SHUNNED"
	case "OFFLINE_SOFT":
		srvs.Status = "OFFLINE_SOFT"
	case "OFFLINE_HARD":
		srvs.Status = "OFFLINE_HARD"
	default:
		srvs.Status = "ONLINE"
	}
}

//SetServerWeight set servers weight
func (srvs *Servers) SetServerWeight(weight uint64) {
	srvs.Weight = weight
}

//SetServerCompression set servers compression
func (srvs *Servers) SetServerCompression(compression uint64) {
	srvs.Compression = compression
}

//SetServerMaxConnection set servers max_connections
func (srvs *Servers) SetServerMaxConnection(maxConnections uint64) {
	if maxConnections >= 10000 {
		srvs.MaxConnections = 10000
	} else {
		srvs.MaxConnections = maxConnections
	}
}

//SetServerMaxReplicationLag set servers max_replication_lag
func (srvs *Servers) SetServerMaxReplicationLag(maxReplicationLag uint64) {
	if maxReplicationLag > 126144000 {
		srvs.MaxReplicationLag = 1261440000
	} else {
		srvs.MaxReplicationLag = maxReplicationLag
	}
}

//SetServerUseSSL set servers use_ssl
func (srvs *Servers) SetServerUseSSL(useSSL uint64) {
	if useSSL >= 1 {
		srvs.UseSsl = 1
	} else {
		srvs.UseSsl = 0
	}
}

//SetServerMaxLatencyMs set servers max_latency_ms
func (srvs *Servers) SetServerMaxLatencyMs(maxLatencyMS uint64) {
	srvs.MaxLatencyMs = maxLatencyMS
}

//SetServersComment set servers comment
func (srvs *Servers) SetServersComment(comment string) {
	srvs.Comment = comment
}

//AddOneServers add a new backend
func (srvs *Servers) AddOneServers(db *sql.DB) error {

	Query := fmt.Sprintf(StmtAddOneServers, srvs.HostGroupID, srvs.HostName, srvs.Port, srvs.MaxConnections)

	_, err := db.Exec(Query)
	if err != nil {
		switch {
		case err.(*mysql.MySQLError).Number == 1045:
			return errors.NewAlreadyExists(err, strconv.Itoa(int(srvs.HostGroupID))+"-"+srvs.HostName+"-"+strconv.Itoa(int(srvs.Port)))
		default:
			return errors.Trace(err) //add server failed
		}
	}

	LoadServerToRuntime(db)
	SaveServerToDisk(db)

	return nil
}

//DeleteOneServers delete a backend
func (srvs *Servers) DeleteOneServers(db *sql.DB) error {

	Query := fmt.Sprintf(StmtDeleteOneServers, srvs.HostGroupID, srvs.HostName, srvs.Port)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(strconv.Itoa(int(srvs.HostGroupID)) + "-" + srvs.HostName + "-" + strconv.Itoa(int(srvs.Port)))
	}

	LoadServerToRuntime(db)
	SaveServerToDisk(db)

	return nil
}

//UpdateOneServerInfo 更新后端服务全部信息
func (srvs *Servers) UpdateOneServerInfo(db *sql.DB) error {

	Query := fmt.Sprintf(StmtUpdateOneServer,
		srvs.Status,
		srvs.Weight,
		srvs.Compression,
		srvs.MaxConnections,
		srvs.MaxReplicationLag,
		srvs.UseSsl,
		srvs.MaxLatencyMs,
		srvs.Comment,
		srvs.HostGroupID,
		srvs.HostName,
		srvs.Port)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(strconv.Itoa(int(srvs.HostGroupID)) + "-" + srvs.HostName + "-" + strconv.Itoa(int(srvs.Port)))
	}

	LoadServerToRuntime(db)
	SaveServerToDisk(db)

	return nil
}
