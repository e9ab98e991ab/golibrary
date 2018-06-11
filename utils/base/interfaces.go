/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
// 公共接口

package base

//"fmt"

type StackChaner interface {
	Stack(string)
	Get() ([]string, bool)
}

type Log interface {
	Println(string)
}

type InterfaceSafeLog interface {
	Init(params *SafeLogParams)
	Create()
	Add(s string)
	Find() (value []string, found bool)
	Delete()
}
