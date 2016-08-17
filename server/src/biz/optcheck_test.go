package biz

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadCheckitem(t *testing.T) {

	Convey("加载检查项", t, func() {
		err := initCheckItemGroup()
		So(err, ShouldBeNil)
		printJSON(t, checkItemList)
	})
}
