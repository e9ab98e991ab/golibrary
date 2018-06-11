/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
package socket

import (
	"errors"
	"fmt"
)

//异步并发数量
const (
	MAX_SERVER_CONN_NUM     = 1000
	MAX_CLIENT_CONN_NUM     = 10
	DEFAULT_SERVER_READ_OUT = 10
	DEFAULT_CLIENT_READ_OUT = 10
	USER                    = "unphp"
	PASSWORD                = "123456"
)

var logsHander []func(s string)

var RouteMap map[string]RouterInterface

func init() {
	RouteMap = make(map[string]RouterInterface)
	logsHander = make([]func(s string), 0)
	logsHander = append(logsHander, func(s string) {
		fmt.Println(s)
	})
}

func AddModulesRoute(k string, r RouterInterface) {
	RouteMap[k] = r
}

func getModulesRoute(module string) (RouterInterface, error) {
	r, ok := RouteMap[module]
	if ok {
		return r, nil
	}
	return r, errors.New("Error::not have this modules![" + module + "]")

}

func AddlogHander(hander func(s string)) {
	logsHander = append(logsHander, hander)
}

func log(s string) {
	for _, hander := range logsHander {
		hander(s)
	}
}
