package proxysql

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/juju/errors"
)

//Schedulers 定义调度
type Schedulers struct {
	ID         int64  `json:"id" db:"id"`
	Active     int64  `json:"active" db:"active"`
	IntervalMs int64  `json:"interval_ms" db:"interval_ms"`
	FileName   string `json:"filename" db:"filename"`
	Arg1       string `json:"arg1" db:"arg1"`
	Arg2       string `json:"arg2" db:"arg2"`
	Arg3       string `json:"arg3" db:"arg3"`
	Arg4       string `json:"arg4" db:"arg4"`
	Arg5       string `json:"arg5" db:"arg5"`
	Comment    string `json:"comment" db:"comment"`
}

//定义查询条件
const (
	/*add a new scheduler*/
	StmtAddOneScheduler = `
	INSERT 
	INTO 
		scheduler(id,filename,interval_ms,active,arg1,arg2,arg3,arg4,arg5) 
	VALUES(%d,%q,%d,%d,%q,%q,%q,%q,%q)`

	/*delete a scheduler*/
	StmtDeleteOneScheduler = `
	DELETE 
	FROM 
		scheduler 
	WHERE id = %d`

	/*update a scheduler*/
	StmtUpdateOneScheduler = `
	UPDATE 
		scheduler 
	SET 
		active = %d,
		interval_ms=%d,
		filename = %q,
		arg1=%q,
		arg2=%q,
		arg3=%q,
		arg4=%q,
		arg5=%q,
		comment=%q 
	WHERE 
		id = %d`

	/*query all schedulers.*/
	StmtFindAllScheduler = `
	SELECT 
		id,
		active,
		interval_ms,
		filename,
		ifnull(arg1,""),
		ifnull(arg2,""),
		ifnull(arg3,""),
		ifnull(arg4,""),
		ifnull(arg5,""),
		comment 
	FROM 
		scheduler 
	LIMIT %d 
	OFFSET %d`

	/*query Last Insert Id*/
	StmtFindLastInsertID = `
	SELECT 
		max(id)
	FROM scheduler
	WHERE
		filename = %q
	AND 
		interval_ms = %d`
)

//FindAllSchedulerInfo query all schedulers
func FindAllSchedulerInfo(db *sql.DB, limit uint64, skip uint64) ([]Schedulers, error) {

	var allscheduler []Schedulers

	Query := fmt.Sprintf(StmtFindAllScheduler, limit, skip)

	rows, err := db.Query(Query)
	if err != nil {
		return []Schedulers{}, errors.Trace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tmpscheduler Schedulers
		err = rows.Scan(
			&tmpscheduler.ID,
			&tmpscheduler.Active,
			&tmpscheduler.IntervalMs,
			&tmpscheduler.FileName,
			&tmpscheduler.Arg1,
			&tmpscheduler.Arg2,
			&tmpscheduler.Arg3,
			&tmpscheduler.Arg4,
			&tmpscheduler.Arg5,
			&tmpscheduler.Comment,
		)

		if err != nil {
			continue
		}

		allscheduler = append(allscheduler, tmpscheduler)
	}

	return allscheduler, nil
}

//NewSch new NewSch
func NewSch(filename string, intervalMS int64) (*Schedulers, error) {

	sched := new(Schedulers)

	sched.FileName = filename
	switch {
	case intervalMS < 100:
		sched.IntervalMs = 100
	case intervalMS > 100000000:
		sched.IntervalMs = 100000000
	default:
		sched.IntervalMs = intervalMS
	}

	sched.Active = 0
	sched.Arg1 = "NULL"
	sched.Arg2 = "NULL"
	sched.Arg3 = "NULL"
	sched.Arg4 = "NULL"
	sched.Arg5 = "NULL"

	return sched, nil

}

//SetSchedulerID Set Scheduler id
func (sched *Schedulers) SetSchedulerID(id int64) {
	sched.ID = id
}

//SetSchedulerActive Set Scheduler Active status
func (sched *Schedulers) SetSchedulerActive(active int64) {
	if active >= 1 {
		sched.Active = 1
	} else {
		sched.Active = 0
	}
}

//SetSchedulerArg1 Set Scheduler arg1
func (sched *Schedulers) SetSchedulerArg1(arg1 string) {
	sched.Arg1 = arg1
}

//SetSchedulerArg2 Set Scheduler arg2
func (sched *Schedulers) SetSchedulerArg2(arg2 string) {
	sched.Arg2 = arg2
}

//SetSchedulerArg3 Set Scheduler arg3
func (sched *Schedulers) SetSchedulerArg3(arg3 string) {
	sched.Arg3 = arg3
}

//SetSchedulerArg4 Set Scheduler arg4
func (sched *Schedulers) SetSchedulerArg4(arg4 string) {
	sched.Arg4 = arg4
}

//SetSchedulerArg5 Set Scheduler arg5
func (sched *Schedulers) SetSchedulerArg5(arg5 string) {
	sched.Arg5 = arg5
}

//SetSchedulerIntervalMs Set scheduler interval_ms
func (sched *Schedulers) SetSchedulerIntervalMs(intervalMS int64) {
	switch {
	case intervalMS < 100:
		sched.IntervalMs = 100
	case intervalMS > 100000000:
		sched.IntervalMs = 100000000
	default:
		sched.IntervalMs = intervalMS
	}
}

//AddOneScheduler add a new scheduler
func (sched *Schedulers) AddOneScheduler(db *sql.DB) error {

	Query := fmt.Sprintf(StmtAddOneScheduler, sched.ID, sched.FileName, sched.IntervalMs, sched.Active, sched.Arg1, sched.Arg2, sched.Arg3, sched.Arg4, sched.Arg5)

	_, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	Query = fmt.Sprintf(StmtFindLastInsertID, sched.FileName, sched.IntervalMs)
	rows := db.QueryRow(Query)

	/*
		FIX:
		It will always return 0 when you use sql.Result.LastInsertId() function to get last inserted row id.
		And go-sql-driver/mysql not support transaction.
		So,I Query a max(id) after insert a row.
	*/
	err = rows.Scan(&sched.ID)
	if err != nil {
		return errors.Trace(err)
	}

	LoadSchedulerToRuntime(db)
	SaveSchedulerToDisk(db)

	return nil
}

//DeleteOneScheduler delete a scheduler
func (sched *Schedulers) DeleteOneScheduler(db *sql.DB) error {

	Query := fmt.Sprintf(StmtDeleteOneScheduler, sched.ID)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(strconv.Itoa(int(sched.ID)))
	}

	LoadSchedulerToRuntime(db)
	SaveSchedulerToDisk(db)
	return nil
}

//UpdateOneSchedulerInfo update a scheduler.
func (sched  *Schedulers) UpdateOneSchedulerInfo(db *sql.DB) error {

	Query := fmt.Sprintf(StmtUpdateOneScheduler, sched .Active, sched .IntervalMs, sched .FileName, sched .Arg1, sched .Arg2, sched .Arg3, sched .Arg4, sched .Arg5, sched .Comment, sched .ID)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(strconv.Itoa(int(sched .ID)))
	}

	LoadSchedulerToRuntime(db)
	SaveSchedulerToDisk(db)

	return nil
}
