
package schedule

import (
	"fmt"
 	"time"
	"github.com/astaxie/beego"
	"minigame/common/jsondata"
	"minigame/common/timer"
	"minigame/dbrds"
	"minigame/dbmgo"
)

var (
	g_timer *timer.Timer
)

type Schedule struct {

}
func NewSchedule() *Schedule {
	self := new(Schedule)
	return self
}
func (self *Schedule) Init() {
	beego.Debug("初始化schedule系统")

	// g_timer = timer.NewHourTimer(24)// 从第二天启动服务器开始监听
	// g_timer = timer.NewSecondTimer(time.Duration(timer.GetTodayLeftSec()))// 从当天00:00:00开始监听
	g_timer = timer.NewSecondTimer(10)// 开服后10秒开始监听定时任务
	
	// 初始化DayShedule
	dayShedule := new(DaySchedule)
	dayShedule.Init()

	hourSchedule := new(HourSchedule)
	hourSchedule.Init()

}



/**
 * 每天00:00:00执行一次定时任务
 */
type DaySchedule struct {

}
func (self *DaySchedule) Init() {
	beego.Info("开启每天定时任务系统")
	g_timer.AddTimeFunc(timer.GetTodayLeftSec(), timer.OneDay_SecCnt, -1, self)
}
func (self *DaySchedule) OnTimerRefresh(now int64) bool {

	ct := time.Now()
	d := fmt.Sprintf("%d-%d-%d %d:%d:%d", ct.Year(), ct.Month(), ct.Day(), ct.Hour(), ct.Minute(), ct.Second())
	beego.Info("执行", d, "任务")

	ResetUserInfo()
	ResetGameInfo()

	return true
}
func (self *DaySchedule) OnTimerRunEnd(now int64) {
	
}



/**
 * 每小时执行一次定时任务
 */
type HourSchedule struct {

}
func (self *HourSchedule) Init() {
	beego.Info("开启每小时定时任务系统")

	const ONE_HOUR_SEC int = 3600

	now := time.Now()
	hour := now.Hour()
	nowSec := now.Unix()
	nextHour := hour + 1
	if hour == 24 {
		nextHour = 0
	}
	nextSec := time.Date(now.Year(), now.Month(), now.Day(), nextHour, 0, 0, 0, now.Location()).Unix()
	beego.Info("nowSec: ", nowSec)
	beego.Info("nextSec: ", nextSec)

	leftSec := nextSec - nowSec
	beego.Info("leftSec: ", leftSec)

	g_timer.AddTimeFunc(leftSec, ONE_HOUR_SEC, -1, self)
}
func (self *HourSchedule) OnTimerRefresh(now int64) bool {

	ct := time.Now()
	d := fmt.Sprintf("%d-%d-%d %d:%d:%d", ct.Year(), ct.Month(), ct.Day(), ct.Hour(), ct.Minute(), ct.Second())
	beego.Info("执行", d, "任务")

	return true
}
func (self *HourSchedule) OnTimerRunEnd(now int64) {
	
}

/**
 *
 * 1. 每天重置任务
 *    玩家每天看视频次数looktvcount清0
 *    玩家每天游戏次数gamecount加满，前一天的免费次数不累加
 *    玩家每天完成任务的状态status清除
 * 
 * 
 */

func ResetUserInfo() {
	
	beego.Info("重置mini_users玩家数据")

	dbrds.Del("mini_tvcount")
	dbrds.Del("mini_freecount")
	dbrds.Del("mini_sharecount")
	dbrds.Del("mini_cashbonus")
	dbrds.Del("mini_cashstatus")
	dbrds.Del("mini_cashcount")

	// tvselect := map[string]interface{} {"looktvcount": map[string]int64 {"$ne": 0}}
	// tvinfo := map[string]interface{} {"$set": map[string]int64 {"looktvcount": 0}}
	// dbmgo.UpdateAll("mini_users", tvselect, tvinfo)
	// gcselect := map[string]interface{} {"gamecount": map[string]int64 {"$ne": 10}}
	// gcinfo := map[string]interface{} {"$set": map[string]int64 {"gamecount": 10}}
	// dbmgo.UpdateAll("mini_users", gcselect, gcinfo)

}
func ResetGameInfo() {
	
	beego.Info("重置mini_games游戏数据")
	
	for _, value := range jsondata.PluginsData {
		key := fmt.Sprintf("%s:%d", "mini_status", value.Id)
		beego.Info("key: ", key)
		dbrds.Del(key)
	}

	selector := map[string]bool {"status": true}
	gameinfo := map[string]interface{} {"$set": map[string]bool {"status": false}}
	dbmgo.UpdateAll("mini_games", selector, gameinfo)

}







