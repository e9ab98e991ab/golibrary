/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
//

package base

//"fmt"

// 无限叠加器：实现多协程之间无阻塞写入。
type MySliceChan chan []string

// 实现无阻塞的信道“叠加器”方法
func (this MySliceChan) Stack(value string) {
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

func NewMySliceChan() MySliceChan {
	return make(chan []string, 1)
}
