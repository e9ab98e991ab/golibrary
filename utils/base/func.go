/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
// 公共函数

package base

import (
	//"fmt"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

// 创建日志文件，打开日志句柄
func CreateAppLog(path string, filesName string, mode ...int) (logger *log.Logger, osfile *os.File) {
	// 创建日志目录
	os.MkdirAll(path, 0755)
	files := path + "/" + filesName
	osfile, _ = os.OpenFile(files, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if len(mode) > 0 {
		switch mode[0] {
		case 1:
			logger = log.New(osfile, "", log.Ldate|log.Ltime)
		case 2:
			logger = log.New(osfile, "log: ", log.Ldate|log.Ltime)
		default:
			logger = log.New(osfile, "log: ", log.Ldate|log.Ltime|log.Lshortfile)
		}
	} else {
		logger = log.New(osfile, "log: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	return logger, osfile
}

// 创建日志文件，打开日志句柄
func CreateLog(path string, filesName string) (logger *log.Logger, osfile *os.File) {
	// 创建日志目录
	os.MkdirAll(path, 0755)
	files := path + "/" + filesName
	osfile, _ = os.OpenFile(files, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	logger = log.New(osfile, "log: ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger, osfile
}

func JsonEncode(d interface{}) (rs string) {
	t, _ := json.Marshal(d)
	rs = string(t)
	return
}

// URL链接参数拼接
func UrlParams(url string, params map[string]string) (u string) {
	u = url
	if len(params) > 0 {
		i := strings.Index(url, "?")
		if i == -1 {
			u = u + "?"
			for k, v := range params {
				u = u + k + "=" + v + "&"
			}
			u = strings.Join(strings.Split(u, "")[0:len(u)-1], "")
		} else {
			if i == len(u)-1 {
				for k, v := range params {
					u = u + k + "=" + v + "&"
				}
				u = strings.Join(strings.Split(u, "")[0:len(u)-1], "")

			} else {
				temp := strings.Split(u, "")
				if temp[len(temp)-1] == "&" {
					for k, v := range params {
						u = u + k + "=" + v + "&"
					}
					u = strings.Join(strings.Split(u, "")[0:len(u)-1], "")

				} else {
					for k, v := range params {
						u = u + "&" + k + "=" + v
					}
				}
			}
		}
	}
	return
}

// 获取url里的域名
func UrlDomain(url string) (domain string) {
	reg := regexp.MustCompile(`(?:http://)([^/]+)/.*`)
	domain = reg.ReplaceAllString(url, "${1}")
	return
}

// 检验正则表达式是否正确。
func CheckReg(reg string) (b bool) {
	b = true
	defer func() {
		if err := recover(); err != nil {
			b = false
		}
	}()
	regexp.MustCompile(reg)
	return
}

func RandInt(min, max int) (rs int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := max - min
	return r.Intn(t) + min
}
