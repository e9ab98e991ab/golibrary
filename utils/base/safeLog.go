/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
//

package base

import (
	"fmt"
	"os"
)

type SafeLogParams struct {
	Code        string
	SafeMap     SafeMap
	LogDataPath string
}

type SafeLog struct {
	SafeMap      SafeMap
	LogCode      string
	LogDataPath  string
	LogSliceChan MySliceChan // 任务日志
}

func (this *SafeLog) Init(params *SafeLogParams) {
	this.LogCode = params.Code
	this.LogDataPath = params.LogDataPath
	this.SafeMap = params.SafeMap

}

// 新增 “计划任务日志”
func (this *SafeLog) Create() {
	temp := make([]string, 0)
	this.SafeMap.Insert(this.LogCode, temp)
	this.LogSliceChan = NewMySliceChan()
	go func() {
		applogPath := this.LogDataPath
		os.MkdirAll(applogPath, 0777)
		applog, osfile := CreateAppLog(applogPath, this.LogCode, 1)
		defer func() {
			osfile.Close()
		}()
		for {
			data := <-this.LogSliceChan
			if data == nil {
				break
			}
			for _, msg := range data {
				applog.Println(msg)
			}
		}
	}()
}

// 新增 “计划任务日志”
func (this *SafeLog) Add(s string) {
	fmt.Println(s)
	this.SafeMap.Update(this.LogCode, this.getLogTempSafeUpdateFun(s, 500))
}

// 获得 “计划任务日志”
func (this *SafeLog) Find() (value []string, found bool) {
	if data, ok := this.SafeMap.Find(this.LogCode); ok {
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

// 删除 “计划任务日志”
func (this *SafeLog) Delete() {
	//s, b := this.AppSafeMap.Find(crontabCode)
	//if b {
	//	close(s.(myTempChan))
	//}

	this.SafeMap.Delete(this.LogCode)
	close(this.LogSliceChan)
}

func (this *SafeLog) getLogTempSafeUpdateFun(s string, l int) UpdateFunc {
	return func(data interface{}, has bool) (rs interface{}) {
		if has {
			switch d := data.(type) {
			case []string:
				this.LogSliceChan.Stack(s)
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
