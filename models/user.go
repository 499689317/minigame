package models

import (
	"fmt"
	"time"
	"errors"
	"math/rand"
	"encoding/json"
	"strconv"
	"github.com/astaxie/beego"
	"minigame/dbrds"
	"minigame/dbmgo"
	"minigame/common/jsondata"
	"minigame/common/httpclient"
)

var (
	mini_users 		string
	mini_tvcount 	string
	mini_freecount 	string
	mini_sharecount string
	mini_cashbonus 	string
	mini_cashstatus string
	mini_cashcount  string
)

/**
 * {
 * 		UserId: 用户id
 * 		PluginId: 插件id
 * 		UserName: 用户昵称
 * 		UserHead: 用户头像
 * 		GameCount: 玩家获得的游戏次数，这个数据是不能清理掉的
 * 		FreeCount: 游戏给玩家分配免费游戏次数，每天都要做清理
 * 		GameCash: 玩家获得的红包现金币
 * 		LookTvCount: 玩家每天观看视频次数
 * }
 */
type User struct {
	UserId   		string 		`json:"userid"`
	UUserId			string 		`json:"uuserid"`
	PluginId 		string 		`json:"pluginid"`
	InviteUid 		string		`json:"inviteuid"`
	UserName 		string 		`json:"username"`
	UserHead 		string 		`json:"userhead"`
	RegTime 		string 		`json:"regtime"`
	GameCount 		int64  		`json:"gamecount"`
	FreeCount   	int64  		`json:"freecount"`
	ShareCount 		int64 		`json:"sharecount"`
	GameCash  		int64  		`json:"gamecash"`
	LookTvCount 	int64  		`json:"looktvcount"`
}
type InviteUser struct {
	UserId    		string  	`json:"userid"`
	UserName 		string		`json:"username"`
	UserHead 		string 		`json:"userhead"`
	RegTime 		string 		`json:"regtime"`
}
type CashStatus struct {
	Status 			bool 		`json:"status"`
}

func init() {
	mini_users = "mini_users"
	mini_tvcount = "mini_tvcount"
	mini_freecount = "mini_freecount"
	mini_sharecount = "mini_sharecount"
	mini_cashbonus = "mini_cashbonus"
	mini_cashstatus = "mini_cashstatus"
	mini_cashcount = "mini_cashcount"
	rand.Seed(time.Now().UnixNano()) //利用当前时间的UNIX时间戳初始化rand包
}

// 新增玩家信息
func AddUser(userid, uuserid, pluginid, inviteuid, username, userhead string) bool {

	// 先查找玩家是否已经注册
	freshUser := make(map[string]interface{})
	freshUser["userid"] 	= 	userid
	freshUser["uuserid"]    =   uuserid
	freshUser["pluginid"] 	= 	pluginid
	freshUser["inviteuid"] 	= 	inviteuid
	freshUser["username"] 	= 	username
	freshUser["userhead"] 	= 	userhead
	freshUser["regtime"]	= 	strconv.Itoa(int(time.Now().Unix()))
	freshUser["gamecount"] 	= 	0
	
	key := fmt.Sprintf("%s:%s", mini_users, userid)
	beego.Info("AddUser hmset key", key)

	b := dbrds.HMSet(key, freshUser)
	if !b {
		beego.Error("redis AddUser插入玩家信息失败")
		return false
	}

	// 数据入库
	isMgoOk := dbmgo.Insert(mini_users, freshUser)
	if !isMgoOk {
		beego.Error("mongodb AddUser插入玩家信息失败")
		return false
	}
	return true
}

// 判断用户是否已经存在
func IsUserExist(uid string) bool {

	key := fmt.Sprintf("%s:%s", mini_users, uid)
	beego.Info("IsUserExist Exists key", key)

	// keys := []string{key}
	b := dbrds.Exists(key)
	if b {
		return true
	}
	return false
}

// 获取玩家信息
func GetUser(uid string) (u *User, err error) {

	// mini_user表数据
	key := fmt.Sprintf("%s:%s", mini_users, uid)
	beego.Info("GetUser hgetall key", key)
	userinfo, err := dbrds.HGetAll(key)
	if err != nil {
		beego.Error("GetUser redis mini_users查找玩家信息错误", uid)
		return nil, err
	}
	if len(userinfo) == 0 {
		beego.Warn("redis mini_users玩家未注册信息", uid)
		return nil, nil
	}
	beego.Info(userinfo)
	// redis数据映射到User struct上，读取userinfo时不能用点的形式读取，只能用[]读取形式，否则编译过不了
	userid 		:= 		userinfo["userid"]
	uuserid     :=      userinfo["uuserid"]
	pluginid 	:= 		userinfo["pluginid"]
	inviteuid   := 		userinfo["inviteuid"]
	username 	:= 		userinfo["username"]
	userhead 	:= 		userinfo["userhead"]
	regtime    	:= 		userinfo["regtime"]
	gamecount, _ 	:= 	strconv.ParseInt(userinfo["gamecount"], 10, 64)
	
	beego.Info("userid: ", userid)
	beego.Info("uuserid: ", uuserid)
	beego.Info("pluginid: ", pluginid)
	beego.Info("inviteuid: ", inviteuid)
	beego.Info("username: ", username)
	beego.Info("userhead: ", userhead)
	beego.Info("gamecount: ", gamecount)

	// mini_cashbonus
	cashval, _ := dbrds.HGet(mini_cashbonus, uid)
	gamecash, _ := 	strconv.ParseInt(cashval, 10, 64)

	// mini_freecount表数据
	freeval, _ := dbrds.HGet(mini_freecount, uid)
	freecount, _ := strconv.ParseInt(freeval, 10, 64)
	
	// mini_tv表数据
	looktvval, _ := dbrds.HGet(mini_tvcount, uid)
	looktvcount, _  :=  strconv.ParseInt(looktvval, 10, 64)

	// mini_sharecount表数据
	shareval, _ := dbrds.HGet(mini_sharecount, uid)
	sharecount, _ := strconv.ParseInt(shareval, 10, 64)
	
	// 将数据汇总到User struct
	user := User{userid,uuserid,pluginid,inviteuid,username,userhead,regtime,gamecount,freecount,sharecount,gamecash,looktvcount}
	return &user, nil
}

func IsLookTvTop(uid string) bool {
	ret, err := dbrds.HGet("mini_tvcount", uid)
	if err != nil {
		beego.Error("redis IsLookTvTop查找玩家looktvcount数据失败", uid)
		return true
	}
	beego.Info("IsLookTvTop looktvcount: ", ret)
	looktvcount, _  :=  strconv.ParseInt(ret, 10, 32)// 获取到的还是int64
	beego.Info("looktvcount: ", looktvcount, "jsondata.GameInfoData.TvMaxCount: ", jsondata.GameInfoData.TvMaxCount)
	if looktvcount >= jsondata.GameInfoData.TvMaxCount {
		return true
	}
	return false
}
func AddLookTvCount(uid string) (int64, error) {
	ret, err := dbrds.HIncrBy(mini_tvcount, uid, 1)
	if err != nil {
		beego.Error("redis AddLookTvCount累加looktvcount失败", uid)
		return 0, err
	}
	beego.Info("AddLookTvCount looktvcount: ", ret)
	// value := map[string]interface{} {"looktvcount": ret}
	// isMgoOk := UpdateMgoUser(uid, value)
	// if !isMgoOk {
	// 	beego.Error("mongodb AddLookTvCount累加looktvcount失败", uid)
	// 	return 0, errors.New("mongodb AddLookTvCount累加looktvcount失败")
	// }
	return ret, nil
}


/**
 * 玩家增加游戏次数
 * 判断玩家游戏次数是否足够，包括免费次数在内
 * 扣除玩家游戏次数，包括免费次数
 */
func AddGameCount(uid string, count int64) (int64, error) {
	key := fmt.Sprintf("%s:%s", mini_users, uid)
	beego.Info("AddGameCount HIncrBy key", key)
	field := "gamecount"
	ret, err := dbrds.HIncrBy(key, field, count)
	if err != nil {
		beego.Error("redis AddGameCount累加gamecount失败", uid)
		return 0, err
	}
	beego.Info("AddGameCount gamecount: ", ret)
	value := map[string]interface{} {"gamecount": ret}
	isMgoOk := UpdateMgoUser(uid, value)
	if !isMgoOk {
		beego.Error("mongodb AddGameCount累加gamecount失败", uid)
		return 0, errors.New("mongodb AddGameCount累加gamecount失败")
	}
	return ret, nil
}
func GetFreeCount(uid string) (int64, error) {
	ret, _ := dbrds.HGet(mini_freecount, uid)
	beego.Info("GetFreeCount: ", ret)
	if ret == "" || ret == "0" {
		beego.Warn("redis GetFreeCount玩家没有免费次数", uid)
		return 0, nil
	}
	return strconv.ParseInt(ret, 10, 64)
}
func GetGameCount(uid string) (int64, error) {
	key := fmt.Sprintf("%s:%s", mini_users, uid)
	beego.Info("GetGameCount HGet key", key)
	ret, err := dbrds.HGet(key, "gamecount")
	if err != nil {
		beego.Error("redis GetGameCount获取gamecount失败", uid)
		return 0, err
	}
	return strconv.ParseInt(ret, 10, 64)
}
func IsGameCountEnough(uid string) bool {
	freecount, _ := GetFreeCount(uid)
	gamecount, _ := GetGameCount(uid)
	beego.Info("IsGameCountExist freecount: ", freecount, "gamecount: ", gamecount)
	if freecount >= jsondata.GameInfoData.FreeMaxCount && gamecount < 1 {
		return false
	}
	return true
}
func DelGameCount(uid string) (int64, int64, error) {
	// 优先扣除免费次数，免费次数不足则扣除获取次数
	freecount, _ := GetFreeCount(uid)
	gamecount, _ := GetGameCount(uid)
	if freecount < jsondata.GameInfoData.FreeMaxCount {
		beego.Warn("使用免费游戏次数游戏", freecount)
		freecount, e := dbrds.HIncrBy(mini_freecount, uid, 1)
		return freecount, gamecount, e
	}
	beego.Warn("免费游戏次数已用完", freecount)
	gamecount, ee := AddGameCount(uid, -1)
	return freecount, gamecount, ee
}




/**
 * 分享游戏链接次数
 */
func AddShareCount(uid string) (int64, error) {
	ret, err := dbrds.HIncrBy(mini_sharecount, uid, 1)
	if err != nil {
		beego.Error("redis AddShareCount累加sharecount失败", uid)
		return 0, err
	}
	return ret, nil
}
func GetShareCount(uid string) (int64, error) {
	ret, _ := dbrds.HGet(mini_sharecount, uid)
	beego.Info("GetShareCount: ", ret)
	if ret == "" || ret == "0" {
		beego.Warn("redis GetShareCount玩家没有分享次数", uid)
		return 0, nil
	}
	return strconv.ParseInt(ret, 10, 64)
}
func IsShareCountEnough(uid string) bool {
	sharecount, err := GetShareCount(uid)
	if err != nil {
		return false
	}
	beego.Info("sharecount: ", sharecount, "jsondata.GameInfoData.ShareMaxCount: ", jsondata.GameInfoData.ShareMaxCount)
	if sharecount >= jsondata.GameInfoData.ShareMaxCount {
		return false
	}
	return true
}

/**
 * 获取玩家邀请其它玩家列表
 */
func GetInviteList(uid string) (*[]InviteUser, error) {
	selector := map[string]string{"inviteuid": uid}
	var inviteuser []InviteUser
	b := dbmgo.Find(mini_users, selector, &inviteuser, 100)
	if !b {
		beego.Error("mongodb GetInviteList查询玩家邀请用户列表数据错误")
		return nil, errors.New("mongodb GetInviteList查询玩家邀请用户列表数据错误")
	}
	return &inviteuser, nil
}


/**
 * 玩家红包现金计算
 */
func IsReceiveCash(uid string) bool {
	ret, err := dbrds.HGet(mini_cashstatus, uid)
	if err != nil {
		beego.Error("redis IsReceiveCash查找玩家领取红包状态失败", uid)
		return true
	}
	beego.Info("IsReceiveCash: ", ret)
	if ret == "" {
		beego.Warn("redis IsReceiveCash玩家领取红包暂无记录")
		return false
	}
	isReceive, _  :=  strconv.ParseBool(ret)
	// beego.Info(isReceive)
	return isReceive
}
func UpdateCashStatus(uid string) bool {
	b := dbrds.HSet(mini_cashstatus, uid, true)
	if !b {
		beego.Error("redis UpdateCashStatus更新玩家领取红包状态失败", uid)
		return false
	}
	return true
}

func IsLowCash(cash int64) bool {
	beego.Info("cash: ", cash, "jsondata.GameInfoData.CashBonus.LowCash: ", jsondata.GameInfoData.CashBonus.LowCash)
	if cash < jsondata.GameInfoData.CashBonus.LowCash {
		return true
	}
	return false
}
func IsMaxCash(cash int64) bool {
	beego.Info("cash: ", cash, "jsondata.GameInfoData.CashBonus.MaxCash: ", jsondata.GameInfoData.CashBonus.MaxCash)
	if cash > jsondata.GameInfoData.CashBonus.MaxCash {
		return true
	}
	return false
}
func IsTopCashCount(uid string) bool {
	cashcount := GetCashCount(uid)
	beego.Info("IsTopGetCashCount cashcount: ", cashcount)
	if cashcount >= jsondata.GameInfoData.CashBonus.MaxCount {
		return true
	}
	return false
}
func UpdateCashCount(uid string) (int64, error) {
	ret, err := dbrds.HIncrBy(mini_cashcount, uid, 1)
	if err != nil {
		beego.Error("redis UpdateGetCashCount累加cashcount失败", uid)
		return 0, err
	}
	return ret, nil
}
func GetCashCount(uid string) int64 {
	ret, err := dbrds.HGet(mini_cashcount, uid)
	if err != nil {
		beego.Error("redis GetCashCount查询玩家提现次数失败", uid)
		return 0
	}
	if ret == "" {
		beego.Warn("redis GetCashCount玩家当天还未提现")
		return 0
	}
	cashcount, _  :=  strconv.ParseInt(ret, 10, 64)
	return cashcount
}


func IsGetCashCondition(uid string) bool {
	selector := map[string]string {"userid": uid}
	var cashstatus []CashStatus
	b := dbmgo.Find("mini_games", selector, &cashstatus, 10)
	// beego.Info(cashstatus)
	if !b {
		beego.Error("mongodb IsGetCashCondition查询玩家领取现金状态错误", uid)
		return false
	}
	count := 0
	statuslen := len(cashstatus)

	if statuslen < 6 {
		beego.Warn("IsGetCashCondition玩家游戏个数少于全部个数", uid, statuslen)
		return false
	}
	for i := 0; i < statuslen; i++ {
		if cashstatus[i].Status {
			count++
		}
	}
	beego.Info("count: ", count, "statuslen: ", statuslen)
	if count >= statuslen {
		return true
	}
	return false
}
func CalcCashBonus() int64 {
	randnum := int64(rand.Intn(98)) //返回[0,98)的随机整数
	for _, value := range jsondata.CashInfoData {
		if randnum >= value["minpro"] && randnum < value["maxpro"] {
			beego.Info("当前红包奖励等级为：", randnum, value["minpro"], value["maxpro"])
			minnum := int(value["minnum"])
			maxnum := int(value["maxnum"])
			randnum2 := int64(rand.Intn(maxnum - minnum) + minnum)
			beego.Info("当前红包奖励数目为：", randnum2, minnum, maxnum)
			return randnum2
		}
	}
	return 0
}
func UpdateGameCash(uid string, cash int64) bool {
	b := dbrds.HSet(mini_cashbonus, uid, cash)
	if !b {
		beego.Error("redis UpdateGameCash记录玩家领取现金失败", uid)
		return false
	}
	return true
}

func GetLocalUrl() string {
	if beego.BConfig.RunMode == "dev" {
		return "http://192.168.10.167:9005"
	}
	if beego.BConfig.RunMode == "test" {
		return "http://127.0.0.1:9005"
	}
	return "http://120.79.210.232:9005"
}
func GetDelCashUrl() string {
	return fmt.Sprintf("%s%s", GetLocalUrl(), "/minigame/user/walletout")
}
func GetCashInfoUrl() string {
	return fmt.Sprintf("%s%s", GetLocalUrl(), "/minigame/user/iteminfo")
}
func GetCashListUrl() string {
	return fmt.Sprintf("%s%s", GetLocalUrl(), "/minigame/user/itemcharge")
}
type ItemInfo struct {
	ItemId int64 `json:"itemId"`
	ItemCount int64 `json:"itemCount"`
}
type CashMsg struct {
	ErrCode int64 `json:"errCode"`
	ErrDesc string `json:"errDesc"`
	Data []ItemInfo `json:"data"`
}
func GetGameCash(userid string) (int64, error) {

	gid, _ := strconv.ParseInt(userid, 10, 64)
	url := fmt.Sprintf("%s?gid=%d", GetCashInfoUrl(), gid)
	beego.Info("GetGameCash: ", url)
	ret, err := httpclient.Get(url)
	if err != nil {
		beego.Error("GetGameCash 请求玩家道具错误，超时等其它原因")
		return 0, err
	}
	var msg CashMsg
	json.Unmarshal(ret, &msg)
	// beego.Info(msg)

	if msg.ErrCode != 0 {
		beego.Info("GetGameCash 请求返回玩家道具失败")
		return 0, errors.New("GetGameCash errCode不为0错误")
	}
	for _, v := range msg.Data {
		if v.ItemId == 20001 {
			return v.ItemCount, nil
		}
	}
	beego.Warn("GetGameCash 玩家未获得现金奖励")
	return 0, nil
}

type AddInfo struct {
	ItemId int64 `json:"itemId"`
	Balance int64 `json:"balance"`
	ChangeCount int64 `json:"changeCount"`
}
type AddMsg struct {
	ErrCode int64 `json:"errCode"`
	ErrDesc string `json:"errDesc"`
	Data AddInfo `json:"data"`
}
func AddGameCash(userid, gameid string, cash int64) (int64, error) {
	gid, _ := strconv.ParseInt(userid, 10, 64)
	appId, _ := strconv.ParseInt(gameid, 10, 64)
	param := make(map[string]interface{})
	param["gid"] = gid
	param["itemId"] = 20001
	param["changeType"] = 45
	param["changeCount"] = cash
	param["appId"] = appId
	param["remark"] = "休闲小游戏添加现金奖励"
	param["operate"] = ""
	beego.Info(param)

	ret, err := httpclient.Post(GetCashInfoUrl(), param)
	if err != nil {
		beego.Error("AddGameCash 请求添加现金奖励错误")
		return 0, err
	}
	var msg AddMsg
	json.Unmarshal(ret, &msg)
	beego.Info(msg)

	if msg.ErrCode != 0 {
		beego.Error("AddGameCash 请求添加现金奖励失败")
		return 0, errors.New("AddGameCash errCode不为0错误")
	}
	return msg.Data.Balance, nil
}

type BalInfo struct {
	Balance int64 `json:"balance"`
	AwardAmount int64 `json:"awardAmount"`
}
type DelMsg struct {
	ErrCode int64 `json:"errCode"`
	ErrDesc string `json:"errDesc"`
	Data BalInfo `json:"data"`
}
func DelGameCash(userid, uuserid, gameid, appid string, cash int64) (int64, error) {

	gid, _ := strconv.ParseInt(userid, 10, 64)
	uid, _ := strconv.ParseInt(uuserid, 10, 64)
	appId, _ := strconv.ParseInt(appid, 10, 64)
	param := make(map[string]interface{})
	param["gid"] = gid
	param["userId"] = uid
	param["redpId"] = fmt.Sprintf("%s_%d", userid, time.Now().UnixNano())
	param["gameId"] = gameid
	param["appId"] = appId
	param["activityName"] = "minigame"
	param["awardAmount"] = cash
	beego.Info(param)

	ret, err := httpclient.Post(GetDelCashUrl(), param)
	if err != nil {
		beego.Error("DelGameCash 请求提现错误")
		return 0, err
	}
	var msg DelMsg
	json.Unmarshal(ret, &msg)
	beego.Info(msg)

	if msg.ErrCode != 0 {
		beego.Error("DelGameCash 请求提现失败")
		if msg.ErrCode == -1 {
			return -1, errors.New("DelGameCash 未绑定公众号")
		}
		return 0, errors.New("DelGameCash errCode不为0错误")
	}
	return msg.Data.Balance, nil
}
// 查询玩家流水
type ListInfo struct {
	ItemId int64 `json:"itemId"`
	AppId int64 `json:"appId"`
	ChangeType int64 `json:"changeType"`
	ChangeCount int64 `json:"changeCount"`
	Balance int64 `json:"balance"`
	Remark string `json:"remark"`
	CreateTime string `json:"createTime"`
	Operate string `json:"operate"`
}
type ListMsg struct {
	ErrCode int64 `json:"errCode"`
	ErrDesc string `json:"errDesc"`
	TotalCount int64 `json:"totalCount"`
	BeginNum int64 `json:"beginNum"`
	PageNum int64 `json:"pageNum"`
	Data []ListInfo `json:"data"`
}
func GetCashList(userid, gameid string, itemid int64) (*ListMsg, error) {
	gid, _ := strconv.ParseInt(userid, 10, 64)
	appId, _ := strconv.ParseInt(gameid, 10, 64)
	url := fmt.Sprintf("%s?gid=%d&appId=%d&itemId=%d", GetCashListUrl(), gid, appId, itemid)
	beego.Info("GetCashList: ", url)
	ret, err := httpclient.Get(url)
	if err != nil {
		beego.Error("GetCashList 请求现金流水错误")
		return nil, err
	}
	var msg ListMsg
	json.Unmarshal(ret, &msg)
	beego.Info(msg)

	if msg.ErrCode != 0 {
		beego.Error("GetCashList 请求现金流水失败")
		return nil, errors.New("GetCashList errCode不为0错误")
	}
	return &msg, nil
}


// 更新玩家信息
func UpdateRdsUser() {
	
}
func UpdateMgoUser(uid string, value map[string]interface{}) bool {
	selector := map[string]string {"userid": uid}
	updatevalue := map[string]interface{} {"$set": value}
	isMgoOk := dbmgo.Update(mini_users, selector, updatevalue)
	if !isMgoOk {
		beego.Error("mongodb UpdateMgoUser更新玩家信息失败")
		return false
	}
	return true
}
func UpdateUser(uid string, value map[string]interface{}) bool {
	
	key := fmt.Sprintf("%s:%s", mini_users, uid)
	beego.Info("UpdateUser HMSet key", key)

	b := dbrds.HMSet(key, value)
	if !b {
		beego.Error("redis UpdateUser更新玩家信息失败")
		return false
	}

	// selector := map[string]string{"userid": uid}
	// updatevalue := map[string]interface{}{"$set": value}
	// isMgoOk := dbmgo.Update(mini_users, selector, updatevalue)
	// if !isMgoOk {
	// 	beego.Error("mongodb UpdateUser更新玩家信息失败")
	// 	return false
	// }
	// return true
	
	return UpdateMgoUser(uid, value)
}
