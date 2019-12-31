package jsondata

/**
 * 服务启动时加载必要的json静态数据表
 * Go在操作文件时，提供了一系列的内置方法
 * 1. io/ioutil操作文件
 * 读文件
 * func ReadAll(r io.Reader) ([]byte, error)
 * func ReadFile(filename string) ([]byte, error)
 * 写文件
 * func WriteFile(filename string, data []byte, perm os.FileMode) error
 * 2. os操作文件
 * 创建文件
 * func Create(name string) (file *File, err error)
 * func NewFile(fd uintptr, name string) *File
 * 打开文件
 * func Open(name string) (*File, error)
 * func OpenFile(name string, flag int, perm os.FileMode) (*File, error)
 * 读文件
 * func (f *File) Read(b []byte) (n int, err error)
 * 写文件
 * func (f *File) Write(b []byte) {n int, err error}
 * 3. bufio操作文件
 * 读文件
 * func NewReader(rd io.Reader) *Reader
 * func (b *Reader) Read(p []byte) (n int, err error)
 * func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
 * 写文件
 * func NewWriter(w io.Writer) *Writer
 * func (b *Writer) Write(p []byte) (n int, err error)
 */

import (
 	// "fmt"
 	"strings"
	"strconv"
 	"encoding/json"
 	"io/ioutil"
 	"github.com/astaxie/beego"
)

var (
	// map声明且定义时必须要明确type为Plugins类型，否则无法用.形式来取属性值，好坑
	// Data map[string]Plugins = make(map[string]Plugins)
	// PluginsData Plugins
	PluginsData []PluginInfo = make([]PluginInfo, 0)
	GameInfoData GameInfo
	CashInfoData map[int]map[string]int64 = make(map[int]map[string]int64)
)

/**
 * 休闲游戏信息json表数据
 * data/plugins.json
 * {
 * 		plugins: {
 * 			plugininfo: value
 * 		}
 * }
 */
type PluginInfo struct {
	Id 		int64 	`json:"id"`
	Name 	string 	`json:"name"`
	Scene 	string 	`json:"scene"`
	Path 	string 	`json:"path"`
	Icon 	string 	`json:"icon"`
}
type Plugins struct {
	Box 	PluginInfo `json:"box"`
	Archery PluginInfo `json:"archery"`
	Jump 	PluginInfo `json:"jump"`
	Eyes 	PluginInfo `json:"eyes"`
	Avoid 	PluginInfo `json:"avoid"`
	Parkour PluginInfo `json:"parkour"`
}



/**
 * gameinfo json表数据
 * looktv: 观看视频配置数据
 * cashbonus: 获取现金奖励配置数据
 */
type CashBonusInfo struct {
	MaxCount 	int64 	`json:"maxcount"`
	LowCash     int64   `json:"lowcash"`
	MaxCash 	int64 	`json:"maxcash"`
	CashPro  	string 	`json:"cashpro"`
	CashNum     string  `json:"cashnum"`
}
type GameInfo struct {
	TvMaxCount 		int64 		`json:"tvmaxcount"`
	FreeMaxCount   	int64   	`json:"freemaxcount"`
	ShareMaxCount   int64       `json:"sharemaxcount"`
	CashBonus 	CashBonusInfo 	`json:"cashbonus"`
}




// =====================================
// =====================================
// =====================================
// ======    初始化相关启动游戏要加载的必要数据
// =====================================
// =====================================
// =====================================


// 初始化玩家领取红包概率及红包获取数量区间数据
func InitCashInfo() {

	pros := strings.Split(GameInfoData.CashBonus.CashPro, ",")
	nums := strings.Split(GameInfoData.CashBonus.CashNum, ",")

	var tmp_val int64
	var tmp_pros []int64
	for i := 0; i < len(pros); i++ {
		val1, _ := strconv.ParseInt(pros[i], 10, 64)
		tmp_val += val1
		tmp_pros = append(tmp_pros, tmp_val)
	}
	// beego.Info(tmp_pros)
	for i := 0; i < len(tmp_pros); i++ {
		pro := make(map[string]int64)
		var p int64
		if i - 1 >= 0 {
			p = tmp_pros[i - 1]
		}
		pro["minpro"] = p
		pro["maxpro"] = tmp_pros[i]
		
		tmp_nums := strings.Split(nums[i], "-")
		// beego.Info(tmp_nums)
		pro["minnum"], _ = strconv.ParseInt(tmp_nums[0], 10, 64)
		pro["maxnum"], _ = strconv.ParseInt(tmp_nums[1], 10, 64)

		CashInfoData[i] = pro
	}
	// beego.Info(CashInfoData)
}






func InitJsonData() bool {
	// Data必须要在声明阶段初始化，否则外部包引用Data变量时获取不到Data内的值
	// Data := make(map[string]interface{})

	jsons := make(map[string]string)
	LoadJsonData(jsons)
	// beego.Info(jsons)
	for key, value := range jsons {
		if key == "plugins.json" {
			err := json.Unmarshal([]byte(value), &PluginsData)
			if err != nil {
				beego.Error("plugins.json文件数据解码失败")
				return false
			}
		}
		if key == "gameinfo.json" {
			err := json.Unmarshal([]byte(value), &GameInfoData)
			if err != nil {
				beego.Error("gameinfo.json文件数据解码失败")
				return false
			}
		}
	}
	// beego.Info(PluginsData)
	// beego.Info(GameInfoData.LookTv.MaxCount)
	// beego.Info(GameInfoData.CashBonus.MaxCount)
	beego.Debug("服务静态数据表初始化完成")
	return true
}
func LoadJsonData(jsonmap map[string]string) {

	// 读取文件夹下所有文件
	filedir := "data/"
	fileinfo, err := ioutil.ReadDir(filedir)
	if err != nil {
		beego.Error("读取静态数据表目录错误")
		return
	}
	for _, info := range fileinfo {
		
		if info.IsDir() {
			beego.Warn("目录: ", info.Name())
		} else {
			// 将文件一个一个加载到内存
			name := info.Name()
			filename := filedir + name
			beego.Info("读取静态数据表文件: ", filename)
			filestr := ReadFileByIoutil(filename)
			if filestr == "" {
				beego.Warn("未读取到静态表数据", filename)
				continue
			}
			// beego.Info(filestr)
			jsonmap[name] = filestr
		}
	}
}
// 使用io/ioutil方法读取文件
func ReadFileByIoutil(filename string) string {
	if contents, err := ioutil.ReadFile(filename); err == nil {
		return string(contents)
	}
	return ""
}

