package proxysql

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/juju/errors"
)

type (
	//Users 定义user info
	Users struct {
		Username              string `db:"username" json:"username"`
		Password              string `db:"password" json:"password"`
		Active                uint64 `db:"active" json:"active"`
		UseSsl                uint64 `db:"use_ssl" json:"use_ssl"`
		DefaultHostgroup      uint64 `db:"default_hostgroup" json:"default_hostgroup"`
		DefaultSchema         string `db:"default_schema" json:"default_schema"`
		SchemaLocked          uint64 `db:"schema_locked" json:"schema_locked"`
		TransactionPersistent uint64 `db:"transaction_persistent" json:"transaction_persistent"`
		FastForward           uint64 `db:"fast_forward" json:"fast_forward"`
		Backend               uint64 `db:"backend" json:"backend"`
		Frontend              uint64 `db:"frontend" json:"frontend"`
		MaxConnections        uint64 `db:"max_connections" json:"max_connections"`
	}
)

//定义查询条件
const (
	/*add a new users*/
	StmtAddOneUser = `
	INSERT INTO 
		mysql_users(username,password,default_hostgroup,default_schema)
	VALUES(%q,%q,%d,%q)`

	/*delete a user*/
	StmtDeleteOneUser = `
	DELETE FROM 
		mysql_users 
	WHERE 
		username = %q
	AND
		backend = %d
	AND
		frontend = %d
	`

	/*list all users*/
	StmtFindAllUserInfo = `
	SELECT 
		ifnull(username,""),
		ifnull(password,""),
		ifnull(active,0),
		ifnull(use_ssl,0),
		ifnull(default_hostgroup,0),
		ifnull(default_schema,""),
		ifnull(schema_locked,0),
		ifnull(transaction_persistent,0),
		ifnull(fast_forward,0),
		ifnull(backend,0),
		ifnull(frontend,0),
		ifnull(max_connections,0) 
	FROM mysql_users 
	LIMIT %d 
	OFFSET %d`

	/*update a users*/
	StmtUpdateOneUser = `
	UPDATE 
		mysql_users 
	SET 
		password=%q,
		active=%d,
		use_ssl=%d,
		default_hostgroup=%d,
		default_schema=%q,
		schema_locked=%d,
		transaction_persistent=%d,
		fast_forward=%d,
		backend=%d,
		frontend=%d,
		max_connections=%d 
	WHERE 
		username = %q
	AND
		backend = %d
	AND
		frontend = %d
		`
)

//FindAllUserInfo list all users.
func FindAllUserInfo(db *sql.DB, limit uint64, skip uint64) ([]Users, error) {
	var alluser []Users

	Query := fmt.Sprintf(StmtFindAllUserInfo, limit, skip)

	rows, err := db.Query(Query)
	if err != nil {
		return []Users{}, errors.Trace(err)
	}
	defer rows.Close()

	for rows.Next() {

		var tmpusr Users

		err = rows.Scan(
			&tmpusr.Username,
			&tmpusr.Password,
			&tmpusr.Active,
			&tmpusr.UseSsl,
			&tmpusr.DefaultHostgroup,
			&tmpusr.DefaultSchema,
			&tmpusr.SchemaLocked,
			&tmpusr.TransactionPersistent,
			&tmpusr.FastForward,
			&tmpusr.Backend,
			&tmpusr.Frontend,
			&tmpusr.MaxConnections,
		)

		if err != nil {
			continue
		}

		alluser = append(alluser, tmpusr)
	}
	return alluser, nil
}

//NewUser 新建user
func NewUser(username string, password string, defaultHostgroup uint64, defaultSchema string) (*Users, error) {
	newuser := new(Users)

	newuser.Username = username
	newuser.Password = password
	newuser.DefaultHostgroup = defaultHostgroup
	newuser.DefaultSchema = defaultSchema

	newuser.Active = 1
	newuser.UseSsl = 0
	newuser.SchemaLocked = 0
	newuser.TransactionPersistent = 0
	newuser.FastForward = 0
	newuser.Backend = 1
	newuser.Frontend = 1
	newuser.MaxConnections = 10000

	return newuser, nil
}

//SetFastForward set fast_forward
func (users *Users) SetFastForward(fastForward uint64) {
	if fastForward >= 1 {
		users.FastForward = 1
	} else {
		users.FastForward = 0
	}
}

//SetMaxConnections set max_connections
func (users *Users) SetMaxConnections(maxConnections uint64) {
	switch {
	case maxConnections >= 10000:
		users.MaxConnections = 10000
	case maxConnections <= 1:
		users.MaxConnections = 1
	default:
		users.MaxConnections = maxConnections
	}
}

//SetBackend set backend
func (users *Users) SetBackend(backend uint64) {
	if backend >= 1 {
		users.Backend = 1
	} else {
		users.Backend = 0
	}
}

//SetFrontend set fronted
func (users *Users) SetFrontend(frontend uint64) {
	if frontend >= 1 {
		users.Frontend = 1
	} else {
		users.Frontend = 0
	}
}

//SetUserActive set user active/disactive
func (users *Users) SetUserActive(active uint64) {
	if active >= 1 {
		users.Active = 1
	} else {
		users.Active = 0
	}
}

//SetUseSSL Set users UseSSL
func (users *Users) SetUseSSL(useSSL uint64) {
	if useSSL >= 1 {
		users.UseSsl = 1
	} else {
		users.UseSsl = 0
	}
}

//SetSchemaLocked set users SchemaLocked
func (users *Users) SetSchemaLocked(schemaLocked uint64) {
	if schemaLocked >= 1 {
		users.SchemaLocked = 1
	} else {
		users.SchemaLocked = 0
	}
}

//SetTransactionPersistent set users transaction_persistent
func (users *Users) SetTransactionPersistent(transactionPersistent uint64) {
	if transactionPersistent >= 1 {
		users.TransactionPersistent = 1
	} else {
		users.TransactionPersistent = 0
	}
}

//AddOneUser add a new user.
func (users *Users) AddOneUser(db *sql.DB) error {

	Query := fmt.Sprintf(StmtAddOneUser, users.Username, users.Password, users.DefaultHostgroup, users.DefaultSchema)

	_, err := db.Exec(Query)
	if err != nil {
		switch {
		case err.(*mysql.MySQLError).Number == 1045:
			return errors.NewAlreadyExists(err, users.Username)
		default:
			return errors.Trace(err) //add user failed
		}
	}

	LoadUserToRuntime(db)
	SaveUserToDisk(db)

	return nil
}

//DeleteOneUser delete a user.
func (users *Users) DeleteOneUser(db *sql.DB) error {

	Query := fmt.Sprintf(StmtDeleteOneUser, users.Username, users.Backend, users.Frontend)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err) //delte failed
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(users.Username)
	}

	LoadUserToRuntime(db)
	SaveUserToDisk(db)

	return nil //delete success
}

//UpdateOneUserInfo update a user.
func (users *Users) UpdateOneUserInfo(db *sql.DB) error {

	Query := fmt.Sprintf(StmtUpdateOneUser,
		users.Password,
		users.Active,
		users.UseSsl,
		users.DefaultHostgroup,
		users.DefaultSchema,
		users.SchemaLocked,
		users.TransactionPersistent,
		users.FastForward,
		users.Backend,
		users.Frontend,
		users.MaxConnections,
		users.Username,
		users.Backend,
		users.Frontend)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(users.Username)
	}

	LoadUserToRuntime(db)
	SaveUserToDisk(db)

	return nil
}
