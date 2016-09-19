package task

import (
	"testing"

	"gopkg.in/mgo.v2/bson"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/ysqi/atop/common"
	"github.com/ysqi/atop/common/models"
)

func TestUpdateStatus(t *testing.T) {
	Convey("更新任务状态", t, func() {
		Convey("不合理参数", func() {
			cases := []struct {
				taskID    bson.ObjectId
				newStatus models.TaskStatus
			}{
				{},
				{taskID: "", newStatus: ""},
				{taskID: bson.NewObjectId(), newStatus: ""},
				{taskID: "", newStatus: ""},
				{taskID: bson.NewObjectId(), newStatus: models.TaskStatusNew},
				{taskID: bson.NewObjectId(), newStatus: models.TaskStatusProcessing},
				{taskID: bson.NewObjectId(), newStatus: models.TaskStatusCompleted},
			}
			for _, c := range cases {
				got := UpdateTaskStatus(c.taskID, c.newStatus)
				So(got.(common.ErrBadBody), ShouldNotBeNil)
			}
		})

		Convey("正常更新", func() {
			//新创建任务
		})
	})

}
