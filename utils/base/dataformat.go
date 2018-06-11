/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/

package base

//返回数据的结构
type SendData struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func RsError(code int, err string) (rs SendData) {
	rs.Code = 1
	rs.Data = map[string]interface{}{"rs": code, "error": err}
	return
}

func RsData(data interface{}) (rs SendData) {
	rs.Code = 1
	rs.Data = map[string]interface{}{"rs": 1, "data": data}
	return
}
