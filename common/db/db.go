package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dialInfo *mgo.DialInfo
var mongoSession *mgo.Session

func init() {
	beego.AddAPPStartHook(initDBInfo)
}

func initDBInfo() error {
	c, err := beego.AppConfig.GetSection("serverdb")
	if err != nil {
		return err
	}
	dbName := c["name"]
	host := c["host"]
	userName := c["user"]
	pwd := c["password"]

	if dbName == "" {
		return errors.New("尚未配置指定数据库名称,配置信息[serverdb::name]")
	}
	if host == "" {
		return errors.New("尚未配置数据库连接地址,配置信息[serverdb:host]")
	}
	if userName == "" {
		return errors.New("尚未配置数据库连接用户，配置信息[serverdb:user]")
	}
	if pwd == "" {
		return errors.New("尚未配置数据库连接用户密码，配置信息[serverdb:pwd]")
	}

	dialInfo = &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  60 * time.Second,
		Database: dbName,
		Username: userName,
		Password: pwd,
	}
	beego.Info("现尝试连接数据库")
	mongoSession, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return fmt.Errorf("连接数据库失败，%s", err.Error())
	}
	mongoSession.SetMode(mgo.Monotonic, true)
	beego.Info("登陆用户", userName, "连接数据库", host, "/", dbName, "成功")
	return nil
}

// NewSession 新
func NewSession() *Session {
	session := mongoSession.Copy()
	// session.SetMode(mgo.Monotonic, false)
	// session.DB(dialInfo.Database)
	return &Session{session}
}

// Session 封装
type Session struct {
	*mgo.Session
}

// DefaultDB 使用默认DB
func (s *Session) DefaultDB() *mgo.Database {
	return s.DB(dialInfo.Database)
}

// Do 执行数据库操作
func Do(fn func(db *mgo.Database) error) error {
	session := NewSession()
	defer session.Close()
	db := session.DefaultDB()
	return fn(db)
}

// RecordIsExist 检查记录是否存在
func RecordIsExist(collection string, query interface{}) bool {
	if collection == "" {
		return false
	}
	session := NewSession()
	defer session.Close()
	db := session.DefaultDB()
	count, err := db.C(collection).Find(query).Count()
	if err != nil {
		return false
	}
	return count > 0
}

// RecordIsExistByID 通过ID查询记录是否存在
func RecordIsExistByID(collection string, id interface{}) bool {
	return RecordIsExist(collection, bson.M{"_id": id})
}
