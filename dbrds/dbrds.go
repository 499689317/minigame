package dbrds

/**
 * go版本redis驱动
 * go get github.com/go-redis/redis
 */

import (
	"errors"
	"time"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)

var (
	client *redis.Client
)

/**
 * [连接redis]
 */
func InitRedis(hostname, password string, rdsdb int) bool {

	opt := redis.Options{
		Addr: hostname,
		Password: password,
		DB: rdsdb,
	}
	client = redis.NewClient(&opt);
	err := client.Ping().Err()
	if err != nil {
		beego.Error("redis连接失败", err)
		return false
	}
	beego.Debug("redis连接成功")
	return true
}


/**
 * 项目中需要用到的redis方法，先在这里注册一下
 * redis调用接口挂到dbrds上，方便管理
 */

/**
 * String method
 */
func Get(key string) (string, error) {
	cmd := client.Get(key)
	ret, err := cmd.Result()
	if err != nil {
		beego.Error("redis Get返回字符串错误")
		return "", errors.New("redis Get返回字符串错误")
	}
	return ret, nil
}
/**
 * expiration单位为ns
 * 1s = 1000ms = 1000000000ns
 * ex单位为秒
 */
func Set(key, value string, ex time.Duration) bool {
	err := client.Set(key, value, ex * time.Second).Err()
	if err != nil {
		beego.Error("redis Set插入字符串错误")
		return false
	}
	return true
}
/**
 * 可变参数keys ...string
 * 变量v类型为[]T，中间使用...分隔，v ...T
 * T 为可变参数的类型，当 T 为 interface{} 时，传入的可以是任意类型。
 */
func Exists(key string) bool {
	ret, err := client.Exists(key).Result()
	// beego.Info(ret)
	// beego.Info(err)
	if err != nil {
		beego.Error("redis Exists判断key是否存在错误")
		return false
	}
	if ret == 1 {
		return true
	}
	return false
}

func Del(key string) bool {
	ret, err := client.Del(key).Result()
	// beego.Info(ret)
	// beego.Info(err)
	if err != nil {
		beego.Error("redis Del删除键错误")
		return false
	}
	if ret == 1 {
		return true
	}
	return false
}


/**
 * map method
 */
func HGet(key, field string) (string, error) {
	// cmd := client.HGet(key, field)
	// ret, err := cmd.Result()
	// beego.Info(ret)
	// if err != nil {
	// 	beego.Error("redis HGet返回value错误", err)
	// 	return "", err
	// }
	ret := client.HGet(key, field).Val()
	return ret, nil
}
func HGetAll(key string) (map[string]string, error) {
	cmd := client.HGetAll(key)
	ret, err := cmd.Result()
	if err != nil {
		beego.Error("redis HGetAll返回value错误")
		return nil, errors.New("redis HGetAll返回value错误")
	}
	return ret, nil
}
func HSet(key, field string, value interface{}) bool {
	err := client.HSet(key, field, value).Err()
	if err != nil {
		beego.Error("redis HSet插入value错误")
		return false
	}
	return true
}
func HMSet(key string, value map[string]interface{}) bool {
	err := client.HMSet(key, value).Err()
	if err != nil {
		beego.Error("redis HMSet插入map value错误")
		return false
	}
	return true
}
func HIncrBy(key, field string, incr int64) (int64, error) {
	cmd := client.HIncrBy(key, field, incr)
	ret, err := cmd.Result()
	if err != nil {
		beego.Error("redis HIncrBy累计数据错误")
		return 0, errors.New("redis HIncrBy累计数据错误")
	}
	return ret, nil
}

