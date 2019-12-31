package models

import (
	"fmt"
	"errors"
	"strconv"
	"github.com/astaxie/beego"
	"minigame/dbrds"
	"minigame/dbmgo"
	"minigame/common/jsondata"
)

var (
	mini_games string
	mini_status string
)

func init() {
	mini_games = "mini_games"
	mini_status = "mini_status"
}

/**
 * gameid游戏id
 * score历史最高分
 * status是否创关成功
 */
type GameInfo struct {
	UserId   string `json:"userid"`
	GameId   int64 	`json:"gameid"`
	Score 	 int64  `json:"score"`
	Status   bool   `json:"status"`
}

func GetGameInfoById(gameid int64) interface{} {
	for _, value := range jsondata.PluginsData {
		if gameid == value.Id {
			return value
		}
	}
	return nil
}
func IsGameIdLegal(gameid int64) bool {
	
	// t := reflect.TypeOf(jsondata.PluginsData)
	// v := reflect.ValueOf(jsondata.PluginsData)
	// for i := 0; i < v.NumField(); i++ {
	// 	if gameid == v.Field(i).Interface()["Id"] {
	// 		return true
	// 	}
	// 	beego.Info("%s----%v\n", i, v.Field(i).Interface()["Id"])
	// }
	// return false
	
	for _, value := range jsondata.PluginsData {
		// beego.Info("%v", value.Id)
		if gameid == value.Id {
			return true
		}
	}
	return false
}
func AddGameInfo(uid string, gameid int64) bool {
	
	// 初始化游戏数据
	gameInfo := make(map[string]interface{})
	gameInfo["userid"] = uid
	gameInfo["gameid"] = gameid
	gameInfo["score"] = 0
	// gameInfo["status"] = false

	beego.Info(gameInfo)

	key := fmt.Sprintf("%s:%d:%s", mini_games, gameid, uid)
	beego.Info("AddGameInfo HMSet key", key)

	b := dbrds.HMSet(key, gameInfo)
	if !b {
		beego.Error("redis AddGameInfo插入玩家游戏数据失败", uid, gameid)
		return false
	}
	isMgoOk := dbmgo.Insert(mini_games, gameInfo)
	if !isMgoOk {
		beego.Error("mongodb AddGameInfo插入玩家游戏数据失败", uid, gameid)
		return false
	}
	return true
}
// 判断玩家是否存在游戏数据
func IsGameInfoExist(uid string, gameid int64) bool {

	key := fmt.Sprintf("%s:%d:%s", mini_games, gameid, uid)
	beego.Info("IsGameInfoExist Exists key", key)

	b := dbrds.Exists(key)
	if b {
		beego.Warn("玩家游戏数据已存在")
		return true
	}
	return false
}
func GetRedisGameInfo(uid string, gameid int64) (g *GameInfo, err error) {
	
	key := fmt.Sprintf("%s:%d:%s", mini_games, gameid, uid)
	beego.Info("GetRedisGameInfo HGetAll key", key)

	ret, err := dbrds.HGetAll(key)
	if err != nil {
		beego.Error("redis GetRedisGameInfo查找玩家游戏数据失败", uid, gameid)
		return nil, err
	}
	if len(ret) == 0 {
		beego.Warn("redis GetRedisGameInfo玩家游戏数据为空", uid, gameid)
		return nil, nil
	}
	beego.Info(ret)


	var score int64
	// score, e := strconv.Atoi(ret["score"])
	score, _ = strconv.ParseInt(ret["score"], 10, 64)


	var status bool
	statusKey := fmt.Sprintf("%s:%d", mini_status, gameid)
	gameStatus, _ := dbrds.HGet(statusKey, uid)
	if gameStatus == "" {
		status = false
	}
	status, _ = strconv.ParseBool(gameStatus)
	

	gameInfo := GameInfo{uid, gameid, score, status}
	return &gameInfo, nil
}
func GetMongodbGameInfo(uid string) (*[]GameInfo, error) {
	
	selector := map[string]string{"userid": uid}
	var gameinfos []GameInfo
	b := dbmgo.Find(mini_games, selector, &gameinfos, 100)
	if !b {
		beego.Error("mongodb GetMongodbGameInfo查询玩家游戏数据失败")
		return nil, errors.New("mongodb GetMongodbGameInfo查询玩家游戏数据失败")
	}
	return &gameinfos, nil
}

// 玩家游戏状态
func IsStatusOk(uid string, gameid int64) bool {
	key := fmt.Sprintf("%s:%d", mini_status, gameid)
	beego.Info("IsStatusOk HGet key", key)
	ret, err := dbrds.HGet(key, uid)
	if err != nil {
		beego.Error("redis IsStatusOk查找玩家status数据失败", uid, gameid)
		return false
	}
	beego.Info(ret)
	if ret == "" {
		beego.Warn("redis IsStatusOk玩家状态status未完成")
		return false
	}
	status, e := strconv.ParseBool(ret)// 除了可以将特定的string转为bool外，其它值都会返回一个错误
	if e != nil {
		beego.Error("status 字符串转bool错误", e)
	}
	return status
}
func GetGameScoreById(uid string, gameid int64) int64 {
	key := fmt.Sprintf("%s:%d:%s", mini_games, gameid, uid)
	beego.Info("IsScoreOk HGet key", key)
	ret, err := dbrds.HGet(key, "score")
	if err != nil {
		beego.Error("redis IsScoreOk查找玩家score数据失败", uid, gameid)
		return 0
	}
	score, _ := strconv.ParseInt(ret, 10, 64)
	beego.Info(ret, score)
	return score
}
func UpdateRdsStatus(uid string, gameid int64) bool {
	key := fmt.Sprintf("%s:%d", mini_status, gameid)
	beego.Info("UpdateRdsStatus HSet key", key)
	b := dbrds.HSet(key, uid, true)// 为什么更新成功后还是返回false呢
	if !b {
		beego.Error("redis UpdateRdsStatus更新玩家status失败", uid, gameid)
		return false
	}
	return true
}
func UpdateRdsGameInfo(uid string, gameid int64, value map[string]interface{}) bool {
	key := fmt.Sprintf("%s:%d:%s", mini_games, gameid, uid)
	beego.Info("UpdateRdsGameInfo HMSet key", key)
	return dbrds.HMSet(key, value)
}
func UpdateMgoGameInfo(uid string, gameid int64, value map[string]interface{}) bool {
	selector := map[string]interface{} {"userid": uid, "gameid": gameid}
	updatevalue := map[string]interface{} {"$set": value}
	isMgoOk := dbmgo.Update(mini_games, selector, updatevalue)
	if !isMgoOk {
		beego.Error("mongodb UpdateGameInfo更新玩家游戏数据失败")
		return false
	}
	return true
}
func UpdateGameInfo(uid string, gameid int64, value map[string]interface{}) bool {
	
	key := fmt.Sprintf("%s:%d:%s", mini_games, gameid, uid)
	beego.Info("UpdateGameInfo HMSet key", key)

	b := dbrds.HMSet(key, value)
	if !b {
		beego.Error("redis UpdateGameInfo更新玩家游戏数据失败")
		return false
	}
	// selector := map[string]interface{}{"userid": uid, "gameid": gameid}
	// updatevalue := map[string]interface{}{"$set": value}
	// isMgoOk := dbmgo.Update(mini_games, selector, updatevalue)
	// if !isMgoOk {
	// 	beego.Error("mongodb UpdateGameInfo更新玩家游戏数据失败")
	// 	return false
	// }
	// return true
	return UpdateMgoGameInfo(uid, gameid, value)
}


