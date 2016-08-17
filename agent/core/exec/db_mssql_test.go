package exec

import (
	"testing"

	"git.coding.net/ysqi/atop/common/models"

	. "github.com/smartystreets/goconvey/convey"
)

const mssqlConnStr = `server=192.168.230.113\SD;user id=readonly;password=readonly;`

func TestMSSQLSingleQuery(t *testing.T) {
	cmd := MssqlCmd{}
	Convey("MSSQL:单查询结果", t, func() {
		cmdInfo := &models.CmdInfo{
			Command: "select getDate() as NowTime",
			Args:    []string{mssqlConnStr},
		}
		data, err := cmd.Exec(cmdInfo)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)

		var realData []map[string]interface{}
		So(data, ShouldHaveSameTypeAs, realData)
		realData = data.([]map[string]interface{})
		So(realData, ShouldNotBeEmpty)
		So(realData[0], ShouldContainKey, "NowTime")

		t.Logf("Result:%#v", realData[0]["NowTime"])
	})
}

func TestMSSQLNullQuery(t *testing.T) {
	cmd := MssqlCmd{}
	Convey("MSSQL:空查询结果", t, func() {
		cmdInfo := &models.CmdInfo{
			Command: "select * from sys.all_views where 1!=1 ",
			Args:    []string{mssqlConnStr},
		}
		data, err := cmd.Exec(cmdInfo)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		So(data, ShouldBeEmpty)
	})
}

func TestMSSQLTableQuery(t *testing.T) {
	cmd := MssqlCmd{}
	Convey("MSSQL:空查询结果", t, func() {
		cmdInfo := &models.CmdInfo{
			Command: "select name,object_id,create_date,modify_date from test.sys.all_objects where type=?1 ",
			Args:    []string{mssqlConnStr, "V"},
		}
		data, err := cmd.Exec(cmdInfo)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		So(data, ShouldNotBeEmpty)
	})
}
