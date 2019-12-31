package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"minigame/logger"
	"minigame/dbrds"
	"minigame/dbmgo"
	"minigame/common/jsondata"
	"minigame/common/timer"
	"minigame/common/schedule"
	_ "minigame/routers"
)

func main() {
	if beego.BConfig.RunMode != "prod" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// 初始化服务log配置
	logger.InitLogger()
	// 初始化静态数据表
	jsons_is_ok := jsondata.InitJsonData()
	if !jsons_is_ok {
		beego.Error("jsondata error 【启动失败】")
		return
	}
	jsondata.InitCashInfo()

	// 连接redis
	var rdsip string = beego.AppConfig.String("rdsip")
	var rdsport string = beego.AppConfig.String("rdsport")
	hostname := fmt.Sprintf("%s:%s", rdsip, rdsport)
	var rdspass string = beego.AppConfig.String("rdspass")
	rdsdb, e := beego.AppConfig.Int("rdsdb")
	if e != nil {
		beego.Error("获取redis配置信息错误")
		return
	}
	beego.Info("hostname: ", hostname)
	beego.Info("rdspass: ", rdspass)
	beego.Info("rdsdb: ", rdsdb)
	rds_is_ok := dbrds.InitRedis(hostname, rdspass, rdsdb)
	if !rds_is_ok {
		beego.Error("redis error 【启动失败】")
		return
	}
	// 连接mongodb
	var mgoip string = beego.AppConfig.String("mgoip")
	var mgoport string = beego.AppConfig.String("mgoport")
	addrs := []string{fmt.Sprintf("%s:%s", mgoip, mgoport)}
	var dbname string = beego.AppConfig.String("mgodb")
	var username string = beego.AppConfig.String("mgouser")
	var password string = beego.AppConfig.String("mgopass")
	beego.Info("addrs: ", addrs)
	beego.Info("dbname: ", dbname)
	beego.Info("username: ", username)
	beego.Info("password: ", password)
	mgo_is_ok := dbmgo.InitMongoDB(addrs, dbname, username, password)
	if !mgo_is_ok {
		beego.Error("mongodb error 【启动失败】")
		return
	}

	// 服务定时任务系统
	sys_schedule := schedule.NewSchedule()
	sys_schedule.Init()

	// schedule.ResetUserInfo()
	// schedule.ResetGameInfo()

	// test_mongodb_is_ok()
	// test_redis_is_ok()
	// test_timer_is_ok()

	/**
	 * test go channel
	 * go是阻塞性语言，与nodejs编程思想不一样
	 */
	// c := make(chan int)
	// r := make(chan string)
	// go func() {
	// 	defer beego.Info("channel执行完毕")
	// 	// 如果没有消费掉channel中的值，for循环会被一直阻塞
	// 	for i := 0; i < 5; i++ {
	// 		beego.Info("推入channel消息", i)
	// 		c <- i
	// 	}
	// 	// 使用完channel后，如果不执行close(c)关闭channel，则value, ok := <- c这条语句将一直被阻塞
	// 	defer close(c)
	// }()

	// for {
	// 	// 取channel值，取完退出
	// 	if value, ok := <- c; ok {
			
	// 		beego.Info("取值: ", value)
	// 	} else {
	// 		beego.Info("退出循环")
	// 		break
	// 	}
	// }
	
	// 使用range来迭代操作channel
	// for value := range c {
	// 	beego.Info("range value", value)
	// }

	// goroutine的select多路复用
	// select case语句在匹配channel时会依次从上到下进行，找到匹配就会执行对应语句块(多个匹配会随机执行一个)，如果都不匹配会执行default(有default部份的话)
	// for {
		// select {
		// 	case <- c:
		// 		beego.Info("获取channel c消息")
		// 		break
		// 	case <- r:
		// 		beego.Info("获取channel r消息")
		// 		break
		// 	default:
		// 		break
		// }
	// }

	beego.Run()
}



type TestStruct struct {
	UserId 	string
	Name 	string
	Head 	string
	Gold 	string
	Jewel 	string
}
func test_mongodb_is_ok() {
	
	/**
	 * 在测试的过程中发现有些坑，一些很奇怪的问题
	 * struct内的属性字段首字母必须大写，否则mongodb写不进数据库，读取时也无法识别字段（怀疑属性大写是可读态的）
	 * 虽然struct内的属性字段首字母大写，但是写入mongodb后会自动变为小写(所有的字母都变为小写，真的坑)，但是不会影响读取数据时struct的分片接收数据
	 */
	test_struct := TestStruct{"test_id", "机器人", "xxx.xxxx.xxxx.com", "1000", "2000"}

	table := "test_coll"

	beego.Info(test_struct.UserId)
	beego.Info(test_struct.Name)
	beego.Info(test_struct.Head)
	beego.Info(test_struct.Gold)
	beego.Info(test_struct.Jewel)

	// test_map := make(map[string]string)
	// test_map["userId"] = "test_id"
	// test_map["name"] = "机器人"
	// test_map["head"] = "xxx.xxxx.xxxx.com"
	// test_map["gold"] = "1000"
	// test_map["jewel"] = "2000"

	// dbmgo.Insert(table, &test_struct)

	// key := "userid"
	// value := "test_id"
	// dbmgo.Remove(table, key, value)
	// dbmgo.RemoveAll(table, key, value)

	selector := make(map[string]string)
	selector["userid"] = "test_id"

	// test1_map := make(map[string]interface{})
	// test2_map := make(map[string]interface{})
	// test2_map["xxx"] = test_map
	// test1_map["$set"] = test2_map
	// dbmgo.Update(table, selector, test1_map)
	// dbmgo.Upsert(table, selector, test1_map)
	// dbmgo.UpdateAll(table, selector, test1_map)
	
	var tests []TestStruct
	var test TestStruct
	dbmgo.Find(table, selector, &tests, 100)
	dbmgo.FindOne(table, selector, &test)

	beego.Info(tests)
	beego.Info(test)

}
func test_redis_is_ok() {
	

	/**
	 *
	 * 测试字符串操作
	 */
	test_key := "test_key"
	test_str := "hello world"
	b := dbrds.Set(test_key, test_str, 0)
	if b {
		beego.Info("字符串插入ok")
	}
	str, err := dbrds.Get(test_key)
	if err != nil {
		beego.Error("redis Get错误")
	}
	beego.Info(str)

	/**
	 * 测试超时时间
	 * @type {[type]}
	 */
	bb := dbrds.Set("kkk", "this is string", 100)
	if bb {
		beego.Info("insert kkk ok")
	}
}

func test_timer_is_ok() {
	
	/**
	 * timer api
	 */
	
	// 今天日期26号，判断今天是否是26号
	test_day := 26
	b := timer.IsToday(test_day)
	beego.Info(b)

	// 当前周是本年度第几周
	w := timer.WeekInYear()
	beego.Info(w)

	// 今天开始时间/结束时间的时间戳（unix秒）
	bs := timer.GetTodayBeginSec()
	beego.Info(bs)
	es := timer.GetTodayEndSec()
	beego.Info(es)

	// 今天已经过了/剩余多少秒
	rs := timer.GetTodayRunSec()
	beego.Info(rs)
	ls := timer.GetTodayLeftSec()
	beego.Info(ls)

	s2t := timer.Str2Time("2018/11/26 14:35:00")
	beego.Info(s2t)
	t2s := timer.Time2Str(s2t)
	beego.Info(t2s)



}
