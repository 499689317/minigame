package controllers

import (
	"minigame/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"minigame/common"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title Login
// @Description register and login
// @Param	body		body 	models.UserMsgPkg	true		"body for user content"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /login [post]
func (u *UserController) Login() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "玩家登陆成功"

	var pkg models.UserMsgPkg
	json.Unmarshal(u.Ctx.Input.RequestBody, &pkg)

	var userid,uuserid,pluginid,inviteuid,username,userhead string = pkg.UserId,pkg.UUserId,pkg.PluginId,pkg.InviteUid,pkg.UserName,pkg.UserHead
	beego.Info("userid: ", userid)
	beego.Info("uuserid: ", uuserid)
	beego.Info("pluginid: ", pluginid)
	beego.Info("inviteuid: ", inviteuid)
	beego.Info("username: ", username)
	beego.Info("userhead: ", userhead)

	// 检查玩家登陆参数
	if userid == "" || userid == "undefined" || userid == "string" || uuserid == "" || uuserid == "undefined" || uuserid == "string" {
		beego.Error("Login 接口参数错误", userid, uuserid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	
	// 1. 检查玩家是否是新玩家
	userdata, err := models.GetUser(userid)
	if err != nil {
		beego.Error("Login 查询玩家信息错误")
		msg["errcode"] = common.ErrorCode.Error_User_Login_Illegal
		msg["msgdesc"] = "玩家登陆失败"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	var gamecount,freecount,gamecash,looktvcount,sharecount int64 = 0,0,0,0,0
	if userdata != nil {
		beego.Warn("Login 玩家已注册", userid, uuserid)

		// 2. 如果玩家已经注册，更新昵称与头像
		userinfo := map[string]interface{} {"username": username, "userhead": userhead}
		models.UpdateUser(userid, userinfo)

		// 2.1 返回老玩家必要信息
		gamecount,freecount,gamecash,looktvcount,sharecount = userdata.GameCount,userdata.FreeCount,userdata.GameCash,userdata.LookTvCount,userdata.ShareCount
	} else {

		// 3. 创建新玩家
		beego.Warn("Login 注册新玩家", userid, uuserid)
		b := models.AddUser(userid, uuserid, pluginid, inviteuid, username, userhead)
		if b {
			beego.Warn("Login 新玩家注册成功", userid, uuserid)
		} else {
			beego.Warn("Login 新玩家注册失败", userid, uuserid)
		}
	}

	// 4. 返回玩家基本信息
	info := map[string]interface{} {
		"userid": userid,
		"uuserid": uuserid,
		"pluginid": pluginid,
		"username": username,
		"userhead": userhead,
		"gamecount": gamecount,
		"freecount": freecount,
		"gamecash": gamecash,
		"looktvcount": looktvcount,
		"sharecount": sharecount,
	}
	msg["data"] = info
	u.Data["json"] = msg
	u.ServeJSON()
}


// @Title LookTv
// @Description LookTv get gamecount
// @Param	userid		query 	string	true		"user id"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /look_tv [get]
func (u *UserController) LookTv() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "记录玩家观看视频成功"

	userid := u.GetString("userid")
	beego.Info("userid: ", userid)

	if userid == "" || userid == "undefined" {
		beego.Error("LookTv 接口参数错误", userid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	// 观看视频是否已达上限
	isCountTop := models.IsLookTvTop(userid)
	if isCountTop {
		beego.Error("LookTv 观看视频次数已达上限")
		msg["errcode"] = common.ErrorCode.Error_User_TvCount_Illegal
		msg["msgdesc"] = "观看视频次数已达上限"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	// 累加观看视频次数
	looktvcount, e := models.AddLookTvCount(userid)
	if e != nil {
		beego.Error("LookTv 累计观看视频次数失败")
		msg["errcode"] = common.ErrorCode.Error_User_Add_TvCount_Illegal
		msg["msgdesc"] = "累计观看视频次数失败"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	// 返还玩家游戏次数
	gamecount, ee := models.AddGameCount(userid, 1)
	if ee != nil {
		beego.Error("LookTv 累计玩家游戏次数失败")
		msg["errcode"] = common.ErrorCode.Error_User_Add_GameCount_Illegal
		msg["msgdesc"] = "累计玩家游戏次数失败"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	beego.Info("looktvcount: ", looktvcount, "gamecount: ", gamecount)
	msg["data"] = map[string]int64{"looktvcount": looktvcount, "gamecount": gamecount}
	u.Data["json"] = msg
	u.ServeJSON()
}

// @Title ShareGame
// @Description share game other
// @Param	userid		query 	string	true		"user id"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /share_game [get]
func (u *UserController) ShareGame() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "玩家分享游戏成功"

	userid := u.GetString("userid")
	beego.Info("userid: ", userid)

	if userid == "" || userid == "undefined" {
		beego.Error("ShareGame 接口参数错误", userid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	isCountTop := models.IsShareCountEnough(userid)
	if !isCountTop {
		beego.Warn("ShareGame 玩家分享次数已上限", userid)
		msg["errcode"] = common.ErrorCode.Error_User_ShareCount_Illegal
		msg["msgdesc"] = "玩家分享次数已上限"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	// 分享成功，累计玩家分享次数
	sharecount, err := models.AddShareCount(userid)
	if err != nil {
		beego.Error("ShareGame 玩家分享次数累计失败", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Share_Incr_Illegal
		msg["msgdesc"] = "玩家分享次数累计失败"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	// 次数累计成功，发放分享奖励
	gamecount, ee := models.AddGameCount(userid, 1)
	if ee != nil {
		beego.Error("ShareGame 累计玩家游戏次数失败")
		msg["errcode"] = common.ErrorCode.Error_User_Add_GameCount_Illegal
		msg["msgdesc"] = "累加玩家游戏次数失败"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	
	info := map[string]int64 {"sharecount": sharecount, "gamecount": gamecount}
	msg["data"] = info
	u.Data["json"] = msg
	u.ServeJSON()
}

// @Title GetInviteList
// @Description get invite list
// @Param	userid		query 	string	true		"user id"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /get_invite_list [get]
func (u *UserController) GetInviteList() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "获取玩家邀请用户列表成功"

	userid := u.GetString("userid")
	beego.Info("userid: ", userid)

	if userid == "" || userid == "undefined" {
		beego.Error("GetInviteList 接口参数错误", userid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	ret, err := models.GetInviteList(userid)
	if err != nil {
		beego.Error("GetInviteList 获取邀请用户列表失败", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Invite_Illegal
		msg["msgdesc"] = "获取邀请用户列表失败"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	beego.Info(ret)
	info := map[string]interface{} {"list": ret}
	msg["data"] = info
	u.Data["json"] = msg
	u.ServeJSON()
}

// @Title GetCashBonus
// @Description GetCashBonus get gamecash
// @Param	userid		query 	string	true		"user id"
// @Param	pluginid 	query 	string	true		"pluginid id"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /get_cash_bonus [get]
func (u *UserController) GetCashBonus() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "玩家领取现金红包成功"

	userid := u.GetString("userid")
	pluginid := u.GetString("pluginid")
	beego.Info("userid: ", userid, "pluginid: ", pluginid)

	if userid == "" || userid == "undefined" || pluginid == "" || pluginid == "undefined" {
		beego.Error("GetCashBonus 接口参数错误", userid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	// 是否已经领取过奖励
	if models.IsReceiveCash(userid) {
		beego.Error("GetCashBonus 已领取过现金", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Receiv_Illegal
		msg["msgdesc"] = "已领取过现金"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	// 是否已经达到领取条件
	if !models.IsGetCashCondition(userid) {
		beego.Error("GetCashBonus 未达到领取现金条件", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Cash_Con_Illegal
		msg["msgdesc"] = "未达到领取现金条件"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	
	// 计算当前可以获得多少现金
	cash := models.CalcCashBonus()
	beego.Info("cash: ", cash)
	
	// 同步现金到提现系统
	allcash, err := models.AddGameCash(userid, pluginid, cash)
	if err != nil {
		beego.Error("GetCashBonus 累计保存玩家现金失败", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Cash_Rec_Illegal
		msg["msgdesc"] = "累计保存玩家现金失败"
		u.Data["json"] = msg
		return
	}
	
	isCash := models.UpdateGameCash(userid, cash)
	if !isCash {
		beego.Error("GetCashBonus 玩家当天现金记录失败", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Cash_Rec_Illegal
		msg["msgdesc"] = "玩家当天现金记录失败"
		u.Data["json"] = msg
		return
	}

	b := models.UpdateCashStatus(userid)
	if !b {
		beego.Error("GetCashBonus 玩家领取现金状态未记录上，不记录现金到玩家账上", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Cash_Sta_Illegal
		msg["msgdesc"] = "玩家领取现金状态未记录上，不记录现金到玩家账上"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	
	info := map[string]int64 {"cash": cash, "allcash": allcash}
	msg["data"] = info
	u.Data["json"] = msg
	u.ServeJSON()
}

// @Title GetCashInfo
// @Description GetCashInfo get allcash
// @Param	userid		query 	string	true		"user id"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /get_cash_info [get]
func (u *UserController) GetCashInfo() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "玩家领取现金红包成功"

	userid := u.GetString("userid")
	beego.Info("userid: ", userid)

	if userid == "" || userid == "undefined" {
		beego.Error("GetCashInfo 接口参数错误", userid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	cash, err := models.GetGameCash(userid)
	if err != nil {
		beego.Error("GetCashInfo 获取玩家现金奖励失败", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Cash_Get_Illegal
		msg["msgdesc"] = "获取玩家现金奖励失败"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	cashcount := models.GetCashCount(userid)

	msg["data"] = map[string]interface{} {"cash": cash, "cashcount": cashcount}
	u.Data["json"] = msg
	u.ServeJSON()
}

// @Title DelCashBonus
// @Description DelCashBonus del cash
// @Param	userid		query 	string	true		"user id"
// @Param	uuserid		query 	string	true		"uuserid id"
// @Param	pluginid	query 	string	true		"pluginid id"
// @Param	appid		query 	string	true		"ddz id"
// @Param	cash		query 	int 	true		"cash"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /del_cash_bonus [get]
func (u *UserController) DelCashBonus() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "玩家提现成功"

	userid := u.GetString("userid")
	uuserid := u.GetString("uuserid")
	pluginid := u.GetString("pluginid")
	appid := u.GetString("appid")
	cash, _ := u.GetInt64("cash")
	beego.Info("userid: ", userid, "uuserid: ", uuserid, "pluginid: ", pluginid, "appid: ", appid, "cash: ", cash)

	if userid == "" || userid == "undefined" || uuserid == "" || uuserid == "undefined" || appid == "" || appid == "undefined" || pluginid == "" || pluginid == "undefined" {
		beego.Error("DelCashBonus 接口参数错误", userid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	if models.IsLowCash(cash) {
		beego.Error("DelCashBonus 提现金额少于1元")
		msg["errcode"] = common.ErrorCode.Error_User_Cash_Low_Illegal
		msg["msgdesc"] = "提现金额少于1元"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}
	if models.IsMaxCash(cash) {
		beego.Error("DelCashBonus 提现金额大于5元")
		msg["errcode"] = common.ErrorCode.Error_User_Cash_Top_Illegal
		msg["msgdesc"] = "提现金额大于5元"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	if models.IsTopCashCount(userid) {
		beego.Error("DelCashBonus 当天提现次数已达上限")
		msg["errcode"] = common.ErrorCode.Error_User_Cash_Count_Illegal
		msg["msgdesc"] = "当天提现次数已达上限"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	allcash, err := models.DelGameCash(userid, uuserid, pluginid, appid, cash)
	if err != nil {

		beego.Error("DelCashBonus 提现失败", userid)

		if allcash == -1 {
			msg["errcode"] = -1
			msg["msgdesc"] = "未绑定公众号"
			allcash = 0
		} else {
			msg["errcode"] = common.ErrorCode.Error_User_Cash_Del_Illegal
			msg["msgdesc"] = "提现失败"
		}

		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	// 提现成功后记录当天提现次数
	cashcount, _ := models.UpdateCashCount(userid)

	msg["data"] = map[string]interface{} {"cashcount": cashcount, "allcash": allcash}
	u.Data["json"] = msg
	u.ServeJSON()
}

// @Title GetCashList
// @Description GetCashList get cash list
// @Param	userid		query 	string	true		"user id"
// @Param	pluginid	query 	string	true		"pluginid id"
// @Success 200 {map[string]interface{}}
// @Failure 403 body is empty
// @router /get_cash_list [get]
func (u *UserController) GetCashList() {

	msg := make(map[string]interface{})
	msg["errcode"] = common.ErrorCode.SUCCESS
	msg["msgdesc"] = "获取现金流水成功"

	userid := u.GetString("userid")
	pluginid := u.GetString("pluginid")

	beego.Info("userid: ", userid, "pluginid: ", pluginid)

	if userid == "" || userid == "undefined" || pluginid == "" || pluginid == "undefined" {
		beego.Error("GetCashList 接口参数错误", userid)
		msg["errcode"] = common.ErrorCode.Error_Param_Illegal
		msg["msgdesc"] = "接口参数错误"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	list, err := models.GetCashList(userid, pluginid, 20001)
	if err != nil {
		beego.Error("GetCashList 现金流水获取失败", userid)
		msg["errcode"] = common.ErrorCode.Error_User_Cash_List_Illegal
		msg["msgdesc"] = "现金流水获取失败"
		u.Data["json"] = msg
		u.ServeJSON()
		return
	}

	msg["data"] = map[string]interface{} {"list": list}
	u.Data["json"] = msg
	u.ServeJSON()
}
