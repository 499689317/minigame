package models

import (
	"reflect"
	"github.com/astaxie/beego"
)

var (
	MSG_TYPE MSGTYPE
)
type MSGTYPE struct {
	SCORE 	int32
	STATUS 	int32
}

type UserMsgPkg struct {
	UserId   	string 		`json:"userid"`
	UUserId     string      `json:"uuserid"`
	PluginId 	string 		`json:"pluginid"`
	InviteUid 	string 		`json:"inviteuid"`
	UserName 	string 		`json:"username"`
	UserHead 	string 		`json:"userhead"`
}
type GameMsgPkg struct {
	UserId   	string 		`json:"userid"`
	GameId   	int64 		`json:"gameid"`
	MsgType  	int32  		`json:"msgtype"`
	Score 	 	int64  		`json:"score"`
	Status   	bool   		`json:"status"`
}

func init() {
	MSG_TYPE.SCORE = 1
	MSG_TYPE.STATUS = 2
}
func IsMsgTypeIllegal(tp int32) bool {

	v := reflect.ValueOf(MSG_TYPE)
	for i := 0; i < v.NumField(); i++ {
		// if tp == v.Field(i) {
		// 	return true
		// }
		beego.Info("%s----%v\n", i, v.Field(i).Interface())
	}
	return false
}
func GetUpdateDataByType(tp int32) map[string]interface{} {
	// if tp == MSG_SCORE_TYPE {
	// 	return map[string]interface{}{"score": score}
	// } else if tp == MSG_STATUS_TYPE {
	// 	return map[string]interface{}{"score": score, "status": status}
	// }
	return nil
}



