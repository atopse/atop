package biz

import (
	"testing"

	"github.com/ysqi/atop/common/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFindAgent(t *testing.T) {
	agent := models.AgentInfo{Name: "test", URL: "http://127.0.0.1:6060", IP: "127.0.0.1"}

	Convey("Find New Agent", t, func() {
		AgentMgt.FindAgent(agent)
		got, err := AgentMgt.GetAgentInfo(agent.IP, false)
		So(err, ShouldBeNil)
		So(got, ShouldNotBeNil)
		So(got.Status, ShouldEqual, models.AgentStatusOnline)

		agent.Name = "Test2"
		agent.URL = "http://192.168.1.1:6060"
		agent.IP = "192.168.1.1"

		AgentMgt.FindAgent(agent)
		got, err = AgentMgt.GetAgentInfo(agent.IP, true)

		So(err, ShouldBeNil)
		So(got, ShouldNotBeNil)
		So(got.Status, ShouldEqual, models.AgentStatusOnline)
		So(got.Name, ShouldEqual, agent.Name)
		So(got.URL, ShouldEqual, agent.URL)
		So(got.IP, ShouldEqual, agent.IP)
	})

}

func TestUpdateAgentStauts(t *testing.T) {
	agent := models.AgentInfo{Name: "test", URL: "http://127.0.0.2:6060", IP: "127.0.0.2"}
	Convey("更新Agent状态", t, func() {
		want := models.AgentStatusUnknown
		AgentMgt.saveAgent(&agent)
		err := AgentMgt.UpdateAgentStatus(agent.IP, want)
		So(err, ShouldBeNil)
		got, err := AgentMgt.GetAgentInfo(agent.IP, true)
		So(err, ShouldBeNil)
		So(got.Status, ShouldEqual, want)
	})
}
