package db

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/ysqi/atop/common/config"
	"github.com/ysqi/atop/common/log2"
	"gopkg.in/mgo.v2/bson"
)

var dialInfo *mgo.DialInfo
var mainSession *mgo.Session

func init() {
	log2.Info("初始化数据库连接...")
	c, err := config.AppCfg.GetSection("serverdb")
	if err != nil {
		log2.Fatalf("获取数据库配置[serverdb]失败,%s", err)
	}
	dbName := c["name"]
	host := c["host"]
	userName := c["user"]
	pwd := c["password"]

	if dbName == "" {
		log2.Fatalln("尚未配置指定数据库名称,配置信息[serverdb::name]")
	}
	if host == "" {
		log2.Fatalln("尚未配置数据库连接地址,配置信息[serverdb:host]")
	}
	if userName == "" {
		log2.Fatalln("尚未配置数据库连接用户，配置信息[serverdb:user]")
	}
	if pwd == "" {
		log2.Fatalln("尚未配置数据库连接用户密码，配置信息[serverdb:pwd]")
	}

	dialInfo = &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  60 * time.Second,
		Database: dbName,
		Username: userName,
		Password: pwd,
	}
	log2.Infof("测试数据库MongoDB='%s@%s/%s连接'", userName, dbName, host)
	mainSession, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log2.Fatalf("连接数据库失败,%s", err)
	}
	log2.Infof("初始化数据库连接成功,登录用户<%s>,数据库<%s@%s>", userName, dbName, host)
}

// Session 封装
type Session struct {
	*mgo.Session
}

// DefaultDB 使用默认DB
func (s *Session) DefaultDB() *mgo.Database {
	return s.DB(dialInfo.Database)
}

// GetSession 获取个连接Session
func GetSession() (session *Session, err error) {
	if mainSession == nil {
		mainSession, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			return
		}
		mainSession.SetMode(mgo.Monotonic, false)
	}
	session = &Session{mainSession.Copy()}
	return
}

// Do 执行数据库操作
func Do(fn func(db *mgo.Database) error) error {
	session, err := GetSession()
	if err != nil {
		return err
	}
	defer session.Close()
	db := session.DefaultDB()
	return fn(db)
}

// QueryCount 根据条件查询记录数
func QueryCount(collection string, query interface{}) (int, error) {
	if collection == "" {
		return 0, errors.New("集合名称为空")
	}
	session, err := GetSession()
	if err != nil {
		return 0, err
	}
	db := session.DefaultDB()
	count, err := db.C(collection).Find(query).Count()
	session.Close()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindByID 按条件查找单个记录.
func FindByID(collection string, id interface{}, result interface{}) error {
	return FindOne(collection, bson.M{"_id": id}, result)
}

// FindOne 按条件查找单个记录.
func FindOne(collection string, query interface{}, result interface{}) error {
	if collection == "" {
		return errors.New("集合名称为空")
	}
	session, err := GetSession()
	if err != nil {
		return err
	}
	db := session.DefaultDB()
	err = db.C(collection).Find(query).One(result)
	session.Close()
	return err
}

// FindAll 按条件查找所有记录.
func FindAll(collection string, query interface{}, result interface{}) error {
	if collection == "" {
		return errors.New("集合名称为空")
	}
	session, err := GetSession()
	if err != nil {
		return err
	}
	db := session.DefaultDB()
	err = db.C(collection).Find(query).All(result)
	session.Close()
	return err
}

// Insert 按条件查找所有记录.
func Insert(collection string, docs ...interface{}) error {
	if collection == "" {
		return errors.New("集合名称为空")
	}
	session, err := GetSession()
	if err != nil {
		return err
	}
	db := session.DefaultDB()
	err = db.C(collection).Insert(docs...)
	session.Close()
	return err
}
