// 最顶级路由---type
/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/

package socketTcp

//socket接口
type RouterInterface interface {
	Route(controller, action string, b []byte) (code int, data string)
}
