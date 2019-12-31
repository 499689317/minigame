package logger

import (
	"encoding/json"
	"github.com/astaxie/beego"
	// beego的日志处理是基于logs模块搭建的，默认可以输出到控制台
	// "github.com/astaxie/beego/logs"
)

// beego.Emergency("this is emergency")
// beego.Alert("this is alert")
// beego.Critical("this is critical")
// beego.Error("this is error")
// beego.Warning("this is warning")
// beego.Notice("this is notice")
// beego.Informational("this is informational")
// beego.Debug("this is debug")

type Opt struct {
	FileName string `json:"filename"`
	Level int `json:"level"`
	MaxLines int `json:"maxlines"`
	MaxSize int `json:"maxsize"`
	Daily bool `json:"daily"`
	MaxDays int `json:"maxdays"`
}


func InitLogger() {

	var env = beego.AppConfig.String("runmode")

	/**
	 * beego log级别
	 * LevelEmergency
	 * LevelAlert
	 * LevelCritical
	 * LevelError
	 * LevelWarning
	 * LevelNotice
	 * LevelInformational
	 * LevelDebug
	 */
	// beego.SetLevel(beego.LevelInformational)

	filename := "logs/" + env + ".log"
	level := beego.LevelInformational // 1,2,3
	maxlines := 100000// 最大行数
	maxsize := 100000000// 最大字节（大概100M）
	daily := true// 每天做切分
	maxdays := 7// 保留7天

	opt := Opt{filename, level, maxlines, maxsize, daily, maxdays}
	optbuf, err := json.Marshal(opt)
	if err != nil {
		beego.Error("logger 配置文件解析错误")
		return
	}
	optjson := string(optbuf)
	beego.Info(optjson)
	beego.SetLogger("file", optjson)

	// 关闭控制台输出
	if env == "prod" {
		beego.BeeLogger.DelLogger("console")
	}

	// 日志默认不输出调用文件名和行号，如果期望输出文件名与行号，进行如下设置
	beego.SetLogFuncCall(true)

	beego.Debug("服务log系统配置完成")
}

