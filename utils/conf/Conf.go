/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
//
// 配置数据处理。

package conf

import (
	"strconv"
	"strings"

	config "github.com/gokyle/goconfig"
)

type Conf struct {
	Initconf config.ConfigMap
}

func (this *Conf) Get(section, key string) (data *confData) {
	d, ok := this.Initconf[section][key]
	if !ok {
		panic("conf not have the '[" + section + "][" + key + "]'")
	} else {
		data = &confData{
			Data: d,
		}
	}
	return
}

func (this *conf) GetOr(key, defaultValue string) string {
	keySlice := strings.Split(key, ".")
	d, ok := this.Initconf[keySlice[0]][keySlice[1]]
	if !ok {
		return defaultValue
	} else {
		return d
	}
}

func (this *Conf) IsGet(section, key string) bool {
	_, ok := this.Initconf[section][key]
	return ok
}

type confData struct {
	Data string
}

func (this *confData) String() string {
	return this.Data
}

func (this *confData) Int() int {
	i, _ := strconv.Atoi(this.Data)
	return i
}

func (this *confData) Int64() int64 {
	i, _ := strconv.Atoi(this.Data)
	return int64(i)
}
