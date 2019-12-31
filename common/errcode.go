package common

/**
 * 接口errcode定义说明
 */

var (
	ErrorCode Code
)
type Code struct {
	SUCCESS int32
	Error_Param_Illegal int32

	Error_User_Login_Illegal int32
	Error_User_Regist_Illegal int32
	Error_User_TvCount_Illegal int32
	Error_User_Add_TvCount_Illegal int32
	Error_User_Add_GameCount_Illegal int32
	Error_User_ShareCount_Illegal int32
	Error_User_Share_Incr_Illegal int32
	Error_User_Invite_Illegal int32
	Error_User_Receiv_Illegal int32
	Error_User_Cash_Con_Illegal int32
	Error_User_Cash_Sta_Illegal int32
	Error_User_Cash_Rec_Illegal int32
	Error_User_Cash_Get_Illegal int32
	Error_User_Cash_Del_Illegal int32
	Error_User_Cash_List_Illegal int32
	Error_User_Cash_Low_Illegal int32
	Error_User_Cash_Top_Illegal int32
	Error_User_Cash_Count_Illegal int32

	Error_Game_Id_Illegal int32
	Error_Game_Update_Illegal int32
	Error_Game_Add_Illegal int32
	Error_Game_Status_Illegal int32
	Error_Game_MsgType_Illegal int32
	Error_Game_GameInfo_Illegal int32
	Error_Game_User_Illegal int32
	Error_Game_Count_Illegal int32
	Error_Game_DelCount_Illegal int32
}
func init() {
	ErrorCode.SUCCESS = 0// 接口调用成功
	ErrorCode.Error_Param_Illegal = 10000// 接口参数错误

	ErrorCode.Error_User_Login_Illegal = 10001// 玩家登陆错误
	ErrorCode.Error_User_Regist_Illegal = 10002 // 玩家未注册
	ErrorCode.Error_User_TvCount_Illegal = 10003// 玩家观看视频次数达上限
	ErrorCode.Error_User_Add_TvCount_Illegal = 10004// 累计看视频次数失败
	ErrorCode.Error_User_Add_GameCount_Illegal = 10005// 累计游戏次数失败
	ErrorCode.Error_User_ShareCount_Illegal = 10006// 玩家分享次数不足
	ErrorCode.Error_User_Share_Incr_Illegal = 10007// 玩家累计分享次数失败
	ErrorCode.Error_User_Invite_Illegal = 10008// 获取玩家邀请列表错误
	ErrorCode.Error_User_Receiv_Illegal = 10009// 玩家已领取过现金
	ErrorCode.Error_User_Cash_Con_Illegal = 10010// 玩家未达到领取现金条件
	ErrorCode.Error_User_Cash_Sta_Illegal = 10011// 玩家领取现金未记录上
	ErrorCode.Error_User_Cash_Rec_Illegal = 10012// 记录玩家当天领取现金失败
	ErrorCode.Error_User_Cash_Get_Illegal = 10013// 获取玩家现金奖励失败
	ErrorCode.Error_User_Cash_Del_Illegal = 10014// 提现失败
	ErrorCode.Error_User_Cash_List_Illegal = 10015// 现金流水获取失败
	ErrorCode.Error_User_Cash_Low_Illegal = 10016// 提现金额小于1块
	ErrorCode.Error_User_Cash_Top_Illegal = 10017// 提现金额大于5块
	ErrorCode.Error_User_Cash_Count_Illegal = 10018// 当天提现次数达上限

	ErrorCode.Error_Game_Id_Illegal = 20000// 更新数据gameid非法
	ErrorCode.Error_Game_Update_Illegal = 20001 // 更新数据失败
	ErrorCode.Error_Game_Add_Illegal = 20002// 新增游戏数据失败
	ErrorCode.Error_Game_Status_Illegal = 20003// 玩家已完成过任务
	ErrorCode.Error_Game_MsgType_Illegal = 20004// 接口消息类型不匹配
	ErrorCode.Error_Game_GameInfo_Illegal = 20005// 取玩家大厅数据失败
	ErrorCode.Error_Game_User_Illegal = 20006// 获取玩家数据异常
	ErrorCode.Error_Game_Count_Illegal = 20007// 玩家游戏次数不足
	ErrorCode.Error_Game_DelCount_Illegal = 20008// 扣除玩家游戏次数失败
}
