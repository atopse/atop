package exec

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/ysqi/com"

	"git.coding.net/ysqi/atop/common/models"
	// 加载MSSQL驱动
	_ "github.com/denisenkom/go-mssqldb"
)

// MssqlCmd 执行MSSQL脚本
type MssqlCmd struct {
}

// Exec 执行命令
func (c *MssqlCmd) Exec(cmd *models.CmdInfo) (interface{}, error) {
	//必须存在参数，第一个参数为 数据库连接信息
	if len(cmd.Args) == 0 || cmd.Args[0] == "" {
		return nil, errors.New("执行MSSQL命令，第一次参数必须是数据库连接信息，目前缺失")
	}

	dataSourceName := cmd.Args[0]
	conn, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return c.query(conn, cmd.Command, cmd.Args[1:]...)
}

func (c *MssqlCmd) query(db *sql.DB, cmd string, args ...string) ([]map[string]interface{}, error) {

	argItems := []interface{}{}
	for _, v := range args {
		argItems = append(argItems, v)
	}
	rows, err := db.Query(cmd, argItems...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	vals := []map[string]interface{}{}

	//列
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	//说明没有数据
	if cols == nil {
		return vals, nil
	}
	//提取数据
	for rows.Next() {
		rowValues := make([]interface{}, len(cols))
		for i := 0; i < len(rowValues); i++ {
			rowValues[i] = new(interface{})
		}
		err = rows.Scan(rowValues...)
		if err != nil {
			return nil, err
		}
		m := make(map[string]interface{}, len(cols))
		for i := 0; i < len(rowValues); i++ {
			v := rowValues[i]
			//获取值，转换为字符串
			m[cols[i]] = com.ToStr(reflect.Indirect(reflect.ValueOf(v)))
		}
		vals = append(vals, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return vals, nil
}
