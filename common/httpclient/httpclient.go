package httpclient

/**
 * httplib 库主要用来模拟客户端发送 HTTP 请求，类似于 Curl 工具，支持 JQuery 类似的链式操作。
 * go get github.com/astaxie/beego/httplib
 * 支持生产环境必要的超时及日志输出
 */

import (
 	"github.com/astaxie/beego"
 	"github.com/astaxie/beego/httplib"
 	"time"
 	"encoding/json"
)
var env = beego.AppConfig.String("runmode")

func Get(url string) ([]byte, error) {

	req := httplib.Get(url)
	
	if env != "prod" {
		req.Debug(true)
	}
	// connectTimeout, readWriteTimeout
	req.SetTimeout(5 * time.Second, 5 * time.Second)

	// 以流形式接收数据
	return req.Bytes()
}

// 支持param k/v参数
func Post(url string, param map[string]interface{}) ([]byte, error) {
	
	req := httplib.Post(url)

	if env != "prod" {
		req.Debug(true)
	}
	req.SetTimeout(5 * time.Second, 5 * time.Second)

	/**
	 * map[key string]type
	 * param = {
	 * 		key: value
	 * }
	 * struct = {
	 * 		p string
	 * 		p2 int
	 * }
	 */
	// for key, value := range param {
	// 	beego.Info("%s -> %s\n", key, value)
	// 	req.Param(key, value)
	// }
	xxx, _ := json.Marshal(param)
	var data []byte = []byte(xxx)
	req.Body(data)

	return req.Bytes()
}

