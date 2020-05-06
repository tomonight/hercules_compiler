package proxysql

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/juju/errors"
)

type (
	//ReplicationHostgroup 定义replication host group
	ReplicationHostgroup struct {
		WriterHostgroup uint64 `json:"writer_hostgroup" db:"writer_hostgroup"`
		ReaderHostgroup uint64 `json:"reader_hostgroup" db:"reader_hostgroup"`
		Comment         string `json:"comment" db:"comment"`
	}
)

//定义查询条件
const (
	// add a new replication hostgroup
	StmtAddOneRHG = `
	INSERT INTO 
		mysql_replication_hostgroups(writer_hostgroup,reader_hostgroup,comment)
	VALUES(%d,%d,%q)`

	// Update a replication hostgroup.
	StmtUpdateOneRHG = `
	UPDATE 
		mysql_replication_hostgroups
	SET 
		writer_hostgroup = %d,
		reader_hostgroup = %d,
		comment = %q
	`

	// Delete a replication hostgroup.
	StmtDeleteOneRHG = `
	DELETE FROM 
		mysql_replication_hostgroups
	WHERE
		writer_hostgroup = %d
	AND
		reader_hostgroup = %d
	`
	// query all RHG informations.

	StmtQueryAllRHGs = `
	SELECT 
		* 
	FROM
		mysql_replication_hostgroups
	LIMIT %d 
	OFFSET %d
	`
)

//NewRHG new replication hostgroup instance
func NewRHG(whg uint64, rhg uint64) (*ReplicationHostgroup, error) {

	newrhg := new(ReplicationHostgroup)

	newrhg.WriterHostgroup = whg
	newrhg.ReaderHostgroup = rhg
	newrhg.Comment = ""

	// return new replication hostgroup instance.
	return newrhg, nil
}

//SetWriterHostGroup set writer hostgroup
func (rhg *ReplicationHostgroup) SetWriterHostGroup(writer uint64) {
	rhg.WriterHostgroup = writer
}

//SetReaderHostGroup set reader hostgroup
func (rhg *ReplicationHostgroup) SetReaderHostGroup(reader uint64) {
	rhg.ReaderHostgroup = reader
}

//SetComment set comment
func (rhg *ReplicationHostgroup) SetComment(comment string) {
	rhg.Comment = comment
}

//AddOneRHG add one new replication hostgroup
func (rhg *ReplicationHostgroup) AddOneRHG(db *sql.DB) error {

	Query := fmt.Sprintf(StmtAddOneRHG, rhg.WriterHostgroup, rhg.ReaderHostgroup, rhg.Comment)
	_, err := db.Exec(Query)
	if err != nil {
		switch {
		case err.(*mysql.MySQLError).Number == 1045:
			return errors.NewAlreadyExists(err, strconv.Itoa(int(rhg.WriterHostgroup))+"-"+strconv.Itoa(int(rhg.ReaderHostgroup)))
		default:
			return errors.Trace(err) //add server failed
		}
	}

	LoadServerToRuntime(db)
	SaveServerToDisk(db)

	return nil
}

//UpdateOneRHG update one replication hostgroup
func (rhg *ReplicationHostgroup) UpdateOneRHG(db *sql.DB) error {

	Query := fmt.Sprintf(StmtUpdateOneRHG, rhg.WriterHostgroup, rhg.ReaderHostgroup, rhg.Comment)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(strconv.Itoa(int(rhg.WriterHostgroup)) + "-" + strconv.Itoa(int(rhg.ReaderHostgroup)))
	}

	LoadServerToRuntime(db)
	SaveServerToDisk(db)

	return nil
}

//DeleteOneRHG delete one replication hostgroup
func (rhg *ReplicationHostgroup) DeleteOneRHG(db *sql.DB) error {

	Query := fmt.Sprintf(StmtDeleteOneRHG, rhg.WriterHostgroup, rhg.ReaderHostgroup)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(strconv.Itoa(int(rhg.WriterHostgroup)) + "-" + strconv.Itoa(int(rhg.ReaderHostgroup)))
	}

	LoadServerToRuntime(db)
	SaveServerToDisk(db)

	return nil

}

// query all replication hostgroups

// QueryAllRHG list all mysql_servers
func QueryAllRHG(db *sql.DB, limit uint64, skip uint64) ([]ReplicationHostgroup, error) {

	var allrhg []ReplicationHostgroup

	Query := fmt.Sprintf(StmtQueryAllRHGs, limit, skip)

	rows, err := db.Query(Query)
	if err != nil {
		return []ReplicationHostgroup{}, errors.Trace(err)
	}
	defer rows.Close()

	for rows.Next() {

		var tmprhg ReplicationHostgroup

		err = rows.Scan(
			&tmprhg.WriterHostgroup,
			&tmprhg.ReaderHostgroup,
			&tmprhg.Comment,
		)

		if err != nil {
			continue
		}

		allrhg = append(allrhg, tmprhg)
	}

	return allrhg, nil
}
