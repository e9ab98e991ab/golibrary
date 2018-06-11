/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/

package base

import (
	"html/template"
	"os"
)

type View struct {
	template string
	data     interface{}
	out      string
}

func (v *View) CreatView(f string) error {
	file := f
	fin, err := os.Open(file)
	defer fin.Close()
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	v.template = ""
	for {
		n, _ := fin.Read(buf)
		if 0 == n {
			break
		}
		v.template = v.template + string(buf[:n])
	}

	return nil
}

func (v *View) SetData(d interface{}) *View {
	v.data = d
	return v
}

func (v *View) Disply(n string) (string, error) {
	v.out = ""
	t := template.New(n)
	if t, err := t.Parse(v.template); err != nil {
		return "", err
	} else {
		if err := t.Execute(v, v.data); err != nil {
			return "", err
		} else {
			return v.out, nil
		}
	}
	return "", nil
}

func (v *View) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 {
		v.out = v.out + string(p)
	}
	return
}
