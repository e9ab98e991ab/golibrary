/**
 * @author : godfeer@aliyun.com 
 * @date : 2018/6/11/011 
 **/


package web

import (
	"bytes"
	"html/template"
	"os"
	"path"
	"strings"
)

// 实例视图
type View struct {
	Dir           string                        // 视图模板目录
	FuncMap       template.FuncMap              // 视图函数
	IsCache       bool                          // 是否开启缓存
	templateCache map[string]*template.Template // 模板缓存映射
}

func (v *View) getTemplateInstance(tpl []string) (*template.Template, error) {
	key := strings.Join(tpl, "-")
	// if IsCache, get cached template if exist
	if v.IsCache {
		if v.templateCache[key] != nil {
			return v.templateCache[key], nil
		}
	}
	var (
		t    *template.Template
		e    error
		file []string = make([]string, len(tpl))
	)
	for i, tp := range tpl {
		file[i] = path.Join(v.Dir, tp)
	}
	t = template.New(path.Base(tpl[0]))
	t.Funcs(v.FuncMap)
	t, e = t.ParseFiles(file...)
	if e != nil {
		return nil, e
	}
	if v.IsCache {
		v.templateCache[key] = t
	}
	return t, nil

}

// Render renders template with data.
// Tpl is the file names under template directory, like tpl1,tpl2,tpl3.
func (v *View) Render(tpl string, data map[string]interface{}) ([]byte, error) {
	t, e := v.getTemplateInstance(strings.Split(tpl, ","))
	if e != nil {
		return nil, e
	}
	var buf bytes.Buffer
	e = t.Execute(&buf, data)
	if e != nil {
		return nil, e
	}
	return buf.Bytes(), nil
}

// Has checks the template file existing.
func (v *View) Has(tpl string) bool {
	f := path.Join(v.Dir, tpl)
	_, e := os.Stat(f)
	return e == nil
}

// NoCache sets view cache off and clean cached data.
func (v *View) NoCache() {
	v.IsCache = false
	v.templateCache = make(map[string]*template.Template)
}

// NewView returns view instance with directory.
// It contains bundle template function HTML(convert string to template.HTML).
func NewView(dir string) *View {
	v := new(View)
	v.Dir = dir
	v.FuncMap = make(template.FuncMap)
	v.FuncMap["Html"] = func(str string) template.HTML {
		return template.HTML(str)
	}
	v.IsCache = false
	v.templateCache = make(map[string]*template.Template)
	return v
}
