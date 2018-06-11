/**
 * @author : godfeer@aliyun.com 
 * @date : 2018/6/11/011 
 **/


package web

import (
	goUrl "net/url"
	"path"
	"regexp"
	"strings"
)

const (
	ROUTER_METHOD_GET    = "GET"
	ROUTER_METHOD_POST   = "POST"
	ROUTER_METHOD_PUT    = "PUT"
	ROUTER_METHOD_DELETE = "DELETE"
)

//
type Router struct {
	routeSlice []*Route
}

// 注册并返回一个路由实例
func NewRouter() *Router {
	rt := new(Router)
	rt.routeSlice = make([]*Route, 0)
	return rt
}

func newRoute() *Route {
	route := new(Route)
	route.params = make([]string, 0)
	return route
}

// 注册路由协议：Get方式的，满足路由规则（正则表达式）的处理句柄函数
func (rt *Router) Get(pattern string, fn ...Handler) {
	route := newRoute()
	route.regex, route.params = rt.parsePattern(pattern)
	route.method = ROUTER_METHOD_GET
	route.fn = fn
	rt.routeSlice = append(rt.routeSlice, route)
}

// 注册路由协议：Post方式的，满足路由规则（正则表达式）的处理句柄函数
func (rt *Router) Post(pattern string, fn ...Handler) {
	route := newRoute()
	route.regex, route.params = rt.parsePattern(pattern)
	route.method = ROUTER_METHOD_POST
	route.fn = fn
	rt.routeSlice = append(rt.routeSlice, route)
}

// 注册路由协议：Put方式的，满足路由规则（正则表达式）的处理句柄函数
func (rt *Router) Put(pattern string, fn ...Handler) {
	route := newRoute()
	route.regex, route.params = rt.parsePattern(pattern)
	route.method = ROUTER_METHOD_PUT
	route.fn = fn
	rt.routeSlice = append(rt.routeSlice, route)
}

// 注册路由协议：Delete方式的，满足路由规则（正则表达式）的处理句柄函数
func (rt *Router) Delete(pattern string, fn ...Handler) {
	route := newRoute()
	route.regex, route.params = rt.parsePattern(pattern)
	route.method = ROUTER_METHOD_DELETE
	route.fn = fn
	rt.routeSlice = append(rt.routeSlice, route)
}

func (rt *Router) parsePattern(pattern string) (regex *regexp.Regexp, params []string) {
	params = make([]string, 0)
	segments := strings.Split(goUrl.QueryEscape(pattern), "%2F")
	for i, v := range segments {
		if strings.HasPrefix(v, "%3A") {
			segments[i] = `([\w-%]+)`
			params = append(params, strings.TrimPrefix(v, "%3A"))
		}
	}
	regex, _ = regexp.Compile("^" + strings.Join(segments, "/") + "$")
	return
}

// Find does find matched rule and parse route url, returns route params and matched handlers.
func (rt *Router) Find(url string, method string) (params map[string]string, fn []Handler) {
	sfx := path.Ext(url)
	url = strings.Replace(url, sfx, "", -1)
	// fix path end slash
	url = goUrl.QueryEscape(url)
	if !strings.HasSuffix(url, "%2F") && sfx == "" {
		url += "%2F"
	}
	url = strings.Replace(url, "%2F", "/", -1)
	for _, r := range rt.routeSlice {
		if r.regex.MatchString(url) && r.method == method {
			p := r.regex.FindStringSubmatch(url)
			if len(p) != len(r.params)+1 {
				continue
			}
			params = make(map[string]string)
			for i, n := range r.params {
				params[n] = p[i+1]
			}
			fn = r.fn
			return
		}
	}
	return nil, nil
}

// 路由协议结构体
type Route struct {
	regex  *regexp.Regexp
	method string
	params []string
	fn     []Handler
}

// Handler defines route handler, middleware handler type.
type Handler func(context *Context)

// 路由缓存
type routerCache struct {
	param map[string]string
	fn    []Handler
}
