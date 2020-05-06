package proxysql

import (
	"database/sql"
	"fmt"

	"github.com/juju/errors"
	//"fmt"
)

type (
	//Variables define variable
	Variables struct {
		VariablesName string `db:"Variable_name" json:"variable_name"`
		Value         string `db:"Value" json:"variable_value"`
	}
)

//定义查询条件
const (
	StmtGlobalVariables   = `SHOW GLOBAL VARIABLES`
	StmtUpdateOneVariable = `
	UPDATE 
		global_variables 
	SET 
		variable_value=%q 
	WHERE variable_name = %q`
)

//UpdateOneConfig 更新配置
func UpdateOneConfig(db *sql.DB, varName string, varValue string) error {
	st := fmt.Sprintf(StmtUpdateOneVariable, varValue, varName)

	result, err := db.Exec(st)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotImplementedf(varName)
	}

	LoadMySQLVariablesToRuntime(db)
	LoadAdminVariablesToRuntime(db)
	SaveMySQLVariablesToDisk(db)
	SaveAdminVariablesToDisk(db)
	return nil
}

//GetConfig 获取配置信息
func GetConfig(db *sql.DB) ([]Variables, error) {
	var tmparray []Variables
	var tmp Variables

	rows, err := db.Query(StmtGlobalVariables)
	if err != nil {
		return []Variables{}, errors.Trace(err)
	}

	for rows.Next() {
		tmp = Variables{}
		err = rows.Scan(&tmp.VariablesName, &tmp.Value)
		tmparray = append(tmparray, tmp)
	}

	return tmparray, nil
}
