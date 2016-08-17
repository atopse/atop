package exec

import (
	"testing"

	"git.coding.net/ysqi/atop/common/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOSCommand(t *testing.T) {
	cmd := OSCommand{}
	Convey("OSCommand:命令执行测试", t, func() {
		Convey("Ping 127.0.0.1", func() {
			cmdInfo := &models.CmdInfo{
				Command: "ping",
				Args:    []string{"127.0.0.1"},
			}
			data, err := cmd.Exec(cmdInfo)
			So(err, ShouldBeNil)
			So(data, ShouldNotBeNil)

			var realData string
			So(data, ShouldHaveSameTypeAs, realData)
			realData = data.(string)
			So(realData, ShouldNotBeEmpty)
			So(realData, ShouldContainSubstring, "0%")
		})
		Convey("Dir .", func() {
			cmdInfo := &models.CmdInfo{
				Command: "dir",
			}
			data, err := cmd.Exec(cmdInfo)
			So(err, ShouldBeNil)
			So(data, ShouldNotBeNil)
			var realData string
			So(data, ShouldHaveSameTypeAs, realData)
			realData = data.(string)
			So(realData, ShouldNotBeEmpty)
			So(realData, ShouldContainSubstring, "os_command_test.go")
		})

		Convey("Dir . with params", func() {
			cmdInfo := &models.CmdInfo{
				Command: "dir",
				Args:    []string{".", "/B"},
			}
			data, err := cmd.Exec(cmdInfo)
			So(err, ShouldBeNil)
			So(data, ShouldNotBeNil)
			var realData string
			So(data, ShouldHaveSameTypeAs, realData)
			realData = data.(string)
			So(realData, ShouldNotBeEmpty)
			So(realData, ShouldContainSubstring, "os_command_test.go")
			t.Log(realData)
		})
	})

}
