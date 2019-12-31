package dbmgo

/**
 * go版本mongodb驱动
 * go get gopkg.in/mgo.v2
 */

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	database *mgo.Database
)

/**
 * [连接数据库]
 */
func InitMongoDB(addrs []string, dbname, username, password string) bool {

	opt := mgo.DialInfo{
		Addrs:    addrs,
		Timeout:  10 * time.Second,
		Database: dbname,
		Username: username,
		Password: password,
	}
	session, err := mgo.DialWithInfo(&opt)
	if err != nil {
		beego.Error("mognodb连接失败", err)
		return false
	}
	database = session.DB(dbname)
	beego.Debug("mongodb连接成功")
	return true
}

/**
 * 
 */
func Insert(table string, data interface{}) bool {
	coll := database.C(table)
	err := coll.Insert(data)
	if err != nil {
		beego.Error("插入数据错误: %v \r\ntable: %s \r\ndata: %v \r\n", err.Error(), table, data)
		return false
	}
	return true
}
/**
 * data bson.M
 */
func Remove(table, key string, value interface{}) bool {
	coll := database.C(table)
	err := coll.Remove(bson.M{key: value})
	if err != nil {
		beego.Error("删除数据错误: %v \r\ntable: %s \r\nkey value: %s:%v \r\n", err.Error(), table, key, value)
		return false
	}
	return true
}
/**
 * TODO remove操作是key value参数
 */
func RemoveAll(table, key string, value interface{}) bool {
	coll := database.C(table)
	_, err := coll.RemoveAll(bson.M{key: value})
	if err != nil {
		beego.Error("删除选择的全部数据错误: %v \r\ntable: %s \r\nkey value: %s:%v \r\n", err.Error(), table, key, value)
		return false
	}
	return true
}
/**
 * 
 */
func Upsert(table string, selector, data interface{}) bool {
	coll := database.C(table)
	_, err := coll.Upsert(selector, data)
	if err != nil {
		beego.Error("插入更新数据错误: %v \r\ntable: %s \r\nselector: %v \r\ndata: %v \r\n", err.Error(), table, selector, data)
		return false
	}
	return true
}
/**
 * 
 */
func Update(table string, selector, data interface{}) bool {
	coll := database.C(table)
	err := coll.Update(selector, data)
	if err != nil {
		beego.Error("更新数据错误: %v \r\ntable: %s \r\nselector: %v \r\ndata: %v \r\n", err.Error(), table, selector, data)
		return false
	}
	return true
}
/**
 * 
 */
func UpdateAll(table string, selector, data interface{}) bool {
	coll := database.C(table)
	_, err := coll.UpdateAll(selector, data)
	if err != nil {
		beego.Error("更新选择的全部数据错误: %v \r\ntable: %s \r\nselector: %v \r\ndata: %v \r\n", err.Error(), table, selector, data)
		return false
	}
	return true
}
/**
 * TODO 要兼容selector为nil的条件
 * limit为了防止一次查询结果过大，内存消耗完导致系统崩溃
 * data 分片地址
 */
func Find(table string, selector, data interface{}, limit int) bool {
	coll := database.C(table)
	iter := coll.Find(selector).Limit(limit).Iter()
	err := iter.All(data)
	if err != nil {
		beego.Error("查询数据错误: table: %s \r\nselector: %v \r\n", table, selector)
		return false
	}
	return true
}
/**
 * data 非分片地址
 */
func FindOne(table string, selector, data interface{}) bool {
	coll := database.C(table)
	err := coll.Find(selector).One(data)
	if err != nil {
		if err == mgo.ErrNotFound {
			beego.Error("未查询到数据: %s \r\nselector: %v \r\n", table, selector)
		} else {
			beego.Error("查询单个数据错误: %v \r\ntable: %s \r\nselector: %v\r\n", err.Error(), table, selector)
		}
		return false
	}
	return true
}
