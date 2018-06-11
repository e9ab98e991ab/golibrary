/**
 * @author : godfeer@aliyun.com 
 * @date : 2018/6/11/011 
 **/


package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	CONTEXT_RENDERED = "context_rendered"
	CONTEXT_END      = "context_end"
	CONTEXT_SEND     = "context_send"
)

type Context struct {
	Request     *http.Request       // raw *http.Request
	Base        string              // Base url, as http://domain/
	Url         string              // Path url, as http://domain/path
	RequestUrl  string              // Request url, as http://domain/path?queryString#fragment
	Method      string              // Request method, GET,POST, etc
	Ip          string              // Client Ip
	UserAgent   string              // Client user agent
	Referer     string              // Last visit refer url
	Host        string              // Request host
	Ext         string              // Request url suffix
	IsSSL       bool                // Is https
	IsAjax      bool                // Is ajax
	Response    http.ResponseWriter // native http.ResponseWriter
	Status      int                 // Response status
	Header      map[string]string   // Response header map
	Body        []byte              // Response body bytes
	routeParams map[string]string
	flashData   map[string]interface{}
	eventsFunc  map[string][]reflect.Value
	IsSend      bool // Response is sent or not
	IsEnd       bool // Response is end or not
	app         *App
	layout      string
}

//
func NewContext(app *App, res http.ResponseWriter, req *http.Request) *Context {

	// 初始化Context及其字段属性
	context := new(Context)
	context.flashData = make(map[string]interface{})
	context.eventsFunc = make(map[string][]reflect.Value)
	context.app = app
	context.IsSend = false
	context.IsEnd = false

	// 初始化请求request的属性
	context.Request = req
	context.Url = req.URL.Path
	context.RequestUrl = req.RequestURI
	context.Method = req.Method
	context.Ext = path.Ext(req.URL.Path)
	context.Host = req.Host
	context.Ip = strings.Split(req.RemoteAddr, ":")[0]
	context.IsAjax = req.Header.Get("X-Requested-With") == "XMLHttpRequest"
	context.IsSSL = req.TLS != nil
	context.Referer = req.Referer()
	context.UserAgent = req.UserAgent()
	context.Base = "://" + context.Host + "/"
	if context.IsSSL {
		context.Base = "https" + context.Base
	} else {
		context.Base = "http" + context.Base
	}

	// 初始化响应response的属性
	context.Response = res
	context.Status = 200
	context.Header = make(map[string]string)
	context.Header["Content-Type"] = "text/html;charset=UTF-8"

	// parse form automatically
	req.ParseForm()

	return context
}

// 返回路由匹配的参数值
// Param returns route param by key string which is defined in router pattern string.
func (ctx *Context) Param(key string) string {
	return ctx.routeParams[key]
}

// Flash sets values to this context or gets by key string.
// The flash items are alive in this context only.
func (ctx *Context) Flash(key string, v ...interface{}) interface{} {
	if len(v) > 0 {
		return ctx.flashData[key]
	}
	ctx.flashData[key] = v[0]
	return nil
}

// 注册绑定事件
// 事件名称：
// CONTEXT_RENDERED---模板渲染事件
// STATUS---异常抛出事件（与异常状态挂钩）
// CONTEXT_END---请求处理结束事件
// CONTEXT_SEND---请求响应结束事件
//
func (ctx *Context) On(e string, fn interface{}) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		println("only support function type for Context.On method")
		return
	}
	if ctx.eventsFunc[e] == nil {
		ctx.eventsFunc[e] = make([]reflect.Value, 0)
	}
	ctx.eventsFunc[e] = append(ctx.eventsFunc[e], reflect.ValueOf(fn))
}

// 触发执行绑定事件
// Do invokes event functions of name string in order of that they are be on.
// If args are less than function args, print error and return nil.
// If args are more than function args, ignore extra args.
// It returns [][]interface{} after invoked event function.
func (ctx *Context) Do(e string, args ...interface{}) [][]interface{} {
	_, ok := ctx.eventsFunc[e]
	if !ok {
		return nil
	}
	if len(ctx.eventsFunc[e]) < 1 {
		return nil
	}
	fns := ctx.eventsFunc[e]
	resSlice := make([][]interface{}, 0)
	for _, fn := range fns {
		if !fn.IsValid() {
			println("invalid event function caller for " + e)
		}
		numIn := fn.Type().NumIn()
		if numIn > len(args) {
			println("not enough parameters for Context.Do(" + e + ")")
			return nil
		}
		rArgs := make([]reflect.Value, numIn)
		for i := 0; i < numIn; i++ {
			rArgs[i] = reflect.ValueOf(args[i])
		}
		resValue := fn.Call(rArgs)
		if len(resValue) < 1 {
			resSlice = append(resSlice, []interface{}{})
			continue
		}
		res := make([]interface{}, len(resValue))
		for i, v := range resValue {
			res[i] = v.Interface()
		}
		resSlice = append(resSlice, res)
	}
	return resSlice
}

// 返回所有的input数据，以Map形式
// Input returns all input data map.
func (ctx *Context) Input() map[string]string {
	data := make(map[string]string)
	for key, v := range ctx.Request.Form {
		data[key] = v[0]
	}
	return data
}

// 返回字符串切片，根据给定的键名
// Strings returns string slice of given key.
func (ctx *Context) Strings(key string) []string {
	return ctx.Request.Form[key]
}

// 获取请求参数
func (ctx *Context) String(key string) string {
	return ctx.Request.FormValue(key)
}

// 过滤数据，强制转换为string,并设置默认值
func (ctx *Context) StringOr(key string, def string) string {
	value := ctx.String(key)
	if value == "" {
		return def
	}
	return value
}

// 过滤数据，强制转换为整型
func (ctx *Context) Int(key string) int {
	str := ctx.String(key)
	i, _ := strconv.Atoi(str)
	return i
}

// 过滤数据，强制转换为整型，并设置默认值
func (ctx *Context) IntOr(key string, def int) int {
	i := ctx.Int(key)
	if i == 0 {
		return def
	}
	return i
}

// 过滤数据，强制转换为浮点值，不设置默认值
func (ctx *Context) Float(key string) float64 {
	str := ctx.String(key)
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

// 过滤数据，强制转换为浮点值，并且设置默认值（当key位空时）
func (ctx *Context) FloatOr(key string, def float64) float64 {
	f := ctx.Float(key)
	if f == 0.0 {
		return def
	}
	return f
}

// 过滤数据，强行转换为bool值
func (ctx *Context) Bool(key string) bool {
	str := ctx.String(key)
	b, _ := strconv.ParseBool(str)
	return b
}

// 获取Cookie值
// 或：设置Cookie值，当第二个参数存在时（存在时，设置两个才有效，第一个为值，第二个为过期时间）
func (ctx *Context) Cookie(key string, value ...string) string {
	if len(value) < 1 {
		c, e := ctx.Request.Cookie(key)
		if e != nil {
			return ""
		}
		return c.Value
	}
	if len(value) == 2 {
		t := time.Now()
		expire, _ := strconv.Atoi(value[1])
		t = t.Add(time.Duration(expire) * time.Second)
		cookie := &http.Cookie{
			Name:    key,
			Value:   value[0],
			Path:    "/",
			MaxAge:  expire,
			Expires: t,
		}
		http.SetCookie(ctx.Response, cookie)
		return ""
	}
	return ""
}

// 获取请求的header信息
func (ctx *Context) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

// 302跳转
func (ctx *Context) Redirect(url string, status ...int) {
	ctx.Header["Location"] = url
	if len(status) > 0 {
		ctx.Status = status[0]
		return
	}
	ctx.Status = 302
}

// 设置content-type值
func (ctx *Context) ContentType(contentType string) {
	ctx.Header["Content-Type"] = contentType
}

// 设置json格式的响应，并设置json的响应头信息
func (ctx *Context) Json(data interface{}) {
	bytes, e := json.MarshalIndent(data, "", "    ")
	if e != nil {
		panic(e)
	}
	ctx.ContentType("application/json;charset=UTF-8")
	ctx.Body = bytes
}

// 发送请求响应（如果响应已经发送，不会重复发送）
func (ctx *Context) Send() {
	if ctx.IsSend {
		return
	}
	for name, value := range ctx.Header {
		ctx.Response.Header().Set(name, value)
	}
	ctx.Response.WriteHeader(ctx.Status)
	ctx.Response.Write(ctx.Body)
	ctx.IsSend = true
	ctx.Do(CONTEXT_SEND) //响应发送结束事件
}

// 终止请求处理句柄，并返回请求响应
func (ctx *Context) End() {
	if ctx.IsEnd {
		return
	}
	if !ctx.IsSend {
		ctx.Send()
	}
	ctx.IsEnd = true
	ctx.Do(CONTEXT_END) //触发事件处理---请求处理时间结束事件
}

// 抛出异常，并终止响应
func (ctx *Context) Throw(status int, message ...interface{}) {
	e := strconv.Itoa(status)
	ctx.Status = status
	ctx.Do(e, message...) //触发事件处理---事件名称为状态名称，例如“404”
	ctx.End()
}

// Layout sets layout string.
func (ctx *Context) Layout(str string) {
	ctx.layout = str
}

// 渲染模板并呈现数据，并以字符方式返回
// 如果有错误，将抛出致命异常错误
func (ctx *Context) Tpl(tpl string, data map[string]interface{}) string {
	b, e := ctx.app.view.Render(tpl+".html", data)
	if e != nil {
		panic(e)
	}
	return string(b)
}

// 渲染模板并呈现数据（用以组合模板不同布局中的子模板）
// 结果将以byte格式传递给context.Body.
// 如果有错误，将抛出致命异常错误
func (ctx *Context) Render(tpl string, data map[string]interface{}) {
	b, e := ctx.app.view.Render(tpl+".html", data)
	if e != nil {
		panic(e)
	}
	if ctx.layout != "" {
		l, e := ctx.app.view.Render(ctx.layout+".layout", data)
		if e != nil {
			panic(e)
		}
		b = bytes.Replace(l, []byte("{@Content}"), b, -1)
	}
	ctx.Body = b
	ctx.Do(CONTEXT_RENDERED)
}

// 添加视图函数，它将影响视图实例的全局
func (ctx *Context) Func(name string, fn interface{}) {
	ctx.app.view.FuncMap[name] = fn
}

// 返回App实例引用
func (ctx *Context) App() *App {
	return ctx.app
}

// 下载文件
func (ctx *Context) Download(file string) {
	f, e := os.Stat(file)
	if e != nil {
		ctx.Status = 404
		return
	}
	if f.IsDir() {
		ctx.Status = 403
		return
	}
	output := ctx.Response.Header()
	output.Set("Content-Type", "application/octet-stream")
	output.Set("Content-Disposition", "attachment; filename="+path.Base(file))
	output.Set("Content-Transfer-Encoding", "binary")
	output.Set("Expires", "0")
	output.Set("Cache-Control", "must-revalidate")
	output.Set("Pragma", "public")
	http.ServeFile(ctx.Response, ctx.Request, file)
	ctx.IsSend = true
}
