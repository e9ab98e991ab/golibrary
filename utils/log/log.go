/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
// 系统日志。

package log

import (
	"fmt"

	base "golibrary/utils/base"
	//"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var SysLog *Log

func New() *Log {
	if SysLog == nil {
		SysLog = &Log{
			ChanLogMsg:    NewMyLogChan(),
			ServerSafeMap: base.NewSafeMapRun(),
			AppSafeMap:    base.NewSafeMapRun(),
			Mode:          0,
			AppLogMap:     make(map[string]interface{}),
		}
		SysLog.Init()
	}
	return SysLog
}

// 无限叠加器：实现多协程之间无阻塞写入。
type myLogChan chan []string

// 实现无阻塞的信道“叠加器”方法
func (this myLogChan) Stack(value string) {
	newdata := make([]string, 0)
	newdata = append(newdata, value)
	for {
		select {
		case this <- newdata:
			return
		case old := <-this:
			old = append(old, newdata[0])
			newdata = old
		}
	}
}

func NewMyLogChan() myLogChan {
	return make(chan []string, 1)
}

// 系统服务日志
type Log struct {
	ChanLogMsg      myLogChan
	Mode            int
	ServerSafeMap   base.SafeMap
	AppSafeMap      base.SafeMap
	AppLogMap       map[string]interface{}
	selfSafeMapCode string
}

func (this *Log) Init() {
	this.selfSafeMapCode = "aecmp_ana"
	temp := make([]string, 0)
	this.ServerSafeMap.Insert(this.selfSafeMapCode, temp)
}

// 记录日志
func (this *Log) Println(s string) {
	//funcName, file, line, ok := runtime.Caller(1)
	msg := "[Shop-yun.go] " + time.Now().Format("2006-01-02 15:04:05")
	if _, file, line, ok := runtime.Caller(1); ok {
		fileArr := strings.Split(file, "/")
		msg = msg + " " + fileArr[len(fileArr)-1] + ":" + strconv.Itoa(line) + " "
	}
	msg = msg + ":" + strings.TrimSpace(s)
	switch this.Mode {
	case 0:
		this.add(msg)
	case 1:
		fmt.Println(msg)
	case 2:
		this.add(msg)
		fmt.Println(msg)
	default:
		this.add(msg)
		fmt.Println(msg)
	}
}

func (this *Log) PrintlnTemp(s string) {
	//funcName, file, line, ok := runtime.Caller(1)
	msg := "[Shop-yun.go] " + time.Now().Format("2006-01-02 15:04:05")
	//if _, file, line, ok := runtime.Caller(1); ok {
	//	fileArr := strings.Split(file, "/")
	//	msg = msg + " " + fileArr[len(fileArr)-1] + ":" + strconv.Itoa(line)
	//}
	msg = msg + ":" + s
	switch this.Mode {
	case 0:
		this.addTemp(msg)
	case 1:
		fmt.Println(msg)
	case 2:
		this.addTemp(msg)
		fmt.Println(msg)
	default:
		this.addTemp(msg)
		fmt.Println(msg)
	}
}

func (this *Log) Get() (value []string, found bool) {
	if data, ok := this.ServerSafeMap.Find(this.selfSafeMapCode); ok {
		switch d := data.(type) {
		case []string:
			value = d
		}
		found = ok
	} else {
		found = false
	}
	return
}

func (this *Log) add(s string) {
	this.ServerSafeMap.Update(this.selfSafeMapCode, this.getLogTempSafeUpdateFun(s, 500, 1))
}

func (this *Log) addTemp(s string) {
	this.ServerSafeMap.Update(this.selfSafeMapCode, this.getLogTempSafeUpdateFun(s, 500, 0))
}

func (this *Log) getLogTempSafeUpdateFun(s string, l int, writefile int) base.UpdateFunc {
	return func(data interface{}, has bool) (rs interface{}) {
		if has {
			switch d := data.(type) {
			case []string:
				if writefile == 1 {
					this.ChanLogMsg.Stack(s)
				}
				d = append(d, s)
				count := len(d)
				if count > l {
					rs = d[count-l : count]
				} else {
					rs = d
				}
			}
		}
		return
	}
}

//func (this *Log) AppLog(code string, safelog interface{}) interface{} {
//	if log, ok := this.AppLogMap[code]; ok {
//		return log
//	} else {
//		params := &base.SafeLogParams{
//			Code:        code,
//			SafeMap:     this.ServerSafeMap,
//			LogDataPath: CoreConf.Get("common", "logpath").String() + "/applog",
//		}
//		log := safelog.(base.InterfaceSafeLog)
//		log.Init(params)
//		this.AppLogMap[code] = log
//		return log
//	}
//}
