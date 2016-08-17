package test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	_ "git.coding.net/ysqi/atop/server/src/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func getConfigPath() string {
	_, file, _, _ := runtime.Caller(1)
	dir := filepath.Dir(file)
	for {
		d := filepath.Join(dir, "conf", "app.conf")
		if dir == filepath.VolumeName(dir) {
			return d
		}
		if utils.FileExists(d) {
			return d
		}
		// Parent dir.
		dir = filepath.Dir(dir)
	}
}

func init() {
	// _, file, _, _ := runtime.Caller(1)
	// apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	appPath := getConfigPath()
	beego.TestBeegoInit(filepath.Dir(filepath.Dir(appPath))) 
}

// TestMain is a sample to run an endpoint test
func TestMain(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestMain", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}
