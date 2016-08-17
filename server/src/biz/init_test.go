package biz

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
)

func printJSON(t *testing.T, data interface{}) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Error("JSON parse error: ", err)
		return
	} 
	t.Log(string(content))
}

var configPath = getConfigPath()

func getConfigPath() string {
	_, file, _, _ := runtime.Caller(1)
	dir := filepath.Dir(file)
	for {
		d := filepath.Join(dir, "conf", "app.test.conf")
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
	beego.TestBeegoInit(filepath.Dir(filepath.Dir(configPath)))
}
