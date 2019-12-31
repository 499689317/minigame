package controllers

import (
	"minigame/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"minigame/common"
)

// Operations about Games
type GameController struct {
	beego.Controller
}


// @Title StartGame
// @Description start play game
// @Param	userid		query 	string	true		"user id"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /start_game [get]
func (g *GameController) StartGame() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "开始游戏成功"

	userid := g.GetString("userid")
	beego.Info("userid: ", userid)

	if userid == "" || userid == "undefined" {
		beego.Error("StartGame 接口参数错误", userid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}
	
	b := models.IsGameCountEnough(userid)
	if !b {
		beego.Warn("StartGame 玩家游戏次数不足", userid)
		msg["errcode"] = common.ErrorCode.Error_Game_Count_Illegal
		msg["msgdesc"] = "玩家游戏次数不足"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}
	freecount, gamecount, e := models.DelGameCount(userid)
	if e != nil {
		beego.Error("StartGame 扣除本次游戏次数失败", userid)
		msg["errcode"] = common.ErrorCode.Error_Game_DelCount_Illegal
		msg["msgdesc"] = "扣除本次游戏次数失败"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}

	// TODO 有一个方案可以在这个接口做防刷机制，后期有时间优化可以在这加上

	msg["data"] = map[string]int64 {"freecount": freecount, "gamecount": gamecount}
	g.Data["json"] = msg
	g.ServeJSON()
}


// @Title Update
// @Description update minigame info
// @Param	body		body 	models.GameMsgPkg	true		"body for game content"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /update [post]
func (g *GameController) Update() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "游戏数据更新成功"

	var pkg models.GameMsgPkg
	json.Unmarshal(g.Ctx.Input.RequestBody, &pkg)

	userid,gameid := pkg.UserId,pkg.GameId
	beego.Info("userid: ", userid)
	beego.Info("gameid: ", gameid)

	if userid == "" || userid == "undefined" || gameid == 0 {
		beego.Error("Update 接口参数错误", userid, gameid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}
	// 判断玩家在user表是否存在，防止有人单刷这个接口
	isUserOk := models.IsUserExist(userid)
	if !isUserOk {
		beego.Error("玩家未注册，但更新数据，记录可疑情况", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Regist_Illegal
		msg["msgdesc"] = "玩家未注册，记录可疑情况"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}
	// 判断gameid是否非法
	b := models.IsGameIdLegal(gameid)
	if !b {
		beego.Error("Update gameid不合法", gameid)
		msg["errcode"] = common.ErrorCode.Error_Game_Id_Illegal
		msg["msgdesc"] = "非法的gameid"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}
	// 判断玩家游戏数据是否存在
	isExistGameInfo := models.IsGameInfoExist(userid, gameid)
	if !isExistGameInfo {
		beego.Warn("Update 玩家游戏数据为空，新增初始游戏数据", userid, gameid)
		isAddOk := models.AddGameInfo(userid, gameid)
		if !isAddOk {
			beego.Error("Update 新增玩家游戏数据失败", userid, gameid)
			msg["errcode"] = common.ErrorCode.Error_Game_Add_Illegal
			msg["msgdesc"] = "新增玩家数据失败"
			g.Data["json"] = msg
			g.ServeJSON()
			return
		}
	}

	// 更新玩家数据
	gameinfo := make(map[string]interface{})
	score := models.GetGameScoreById(userid, gameid)
	if score >= pkg.Score {
		beego.Warn("Update 玩家未打破最高分", userid, gameid)
	} else {
		gameinfo["score"] = pkg.Score
		models.UpdateRdsGameInfo(userid, gameid, gameinfo)
	}
	// 判断玩家是否已经完成过任务
	isTaskOk := models.IsStatusOk(userid, gameid)
	if isTaskOk {
		beego.Warn("Update 玩家已完成过任务", userid, gameid)
	} else {
		gameinfo["status"] = true
		models.UpdateRdsStatus(userid, gameid)
	}
	beego.Info(gameinfo)
	if len(gameinfo) == 0 {
		beego.Warn("Update 不更新玩家游戏数据", userid, gameid)
		msg["errcode"] = common.ErrorCode.Error_Game_MsgType_Illegal
		msg["msgdesc"] = "不更新玩家游戏数据"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}
	isUpdateOk := models.UpdateMgoGameInfo(userid, gameid, gameinfo)
	if !isUpdateOk {
		beego.Error("Update 更新玩家数据失败", userid, gameid)
		msg["errcode"] = common.ErrorCode.Error_Game_Update_Illegal
		msg["msgdesc"] = "玩家数据更新失败"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}
	msg["data"] = "update ok"
	g.Data["json"] = msg
	g.ServeJSON()
}

// @Title GetGameInfo
// @Description get hall info
// @Param	userid		query 	string	true		"user id"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /get_game_info [get]
func (g *GameController) GetGameInfo() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "获取游戏数据成功"

	userid := g.GetString("userid")
	beego.Info("userid: ", userid)

	if userid == "" || userid == "undefined" {
		beego.Error("GetGameInfo 接口参数错误", userid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}

	// 判断玩家游戏数据是否存在
	ret, err := models.GetMongodbGameInfo(userid)
	if err != nil {
		beego.Error("GetGameInfo 取玩家大厅数据失败", userid)
		msg["errcode"] = common.ErrorCode.Error_Game_GameInfo_Illegal
		msg["msgdesc"] = "取玩家大厅数据失败"
		g.Data["json"] = msg
		g.ServeJSON()
		return
	}
	beego.Info(ret)

	msg["data"] = map[string]interface{}{"gameinfos": ret}
	g.Data["json"] = msg
	g.ServeJSON()
}

