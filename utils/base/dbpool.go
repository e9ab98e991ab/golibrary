/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
// 数据库链接池

package base

import (
	"fmt"
	"net/url"
	"strings"

	//_ "github.com/go-pg/pq"
	_ "github.com/go-sql-driver/mysql"
	xorm "github.com/go-xorm/xorm"
)

// 数据库链接池
type Dbpool struct {
	DbType string
	DbHost string
	DbName string
	DbUser string
	DbPawd string
	Log    Log
	engine *xorm.Engine
}

// 数据库初始化
func (this *Dbpool) Inits() {
	var err error
	switch this.DbType {
	case "mysql":
		dbtype := this.DbType
		connectstr := this.DbUser + ":" + this.DbPawd
		connectstr = connectstr + "@" + this.DbHost + "/" + this.DbName
		this.engine, err = xorm.NewEngine(dbtype, connectstr+"?charset=utf8")
		if err != nil {
			this.Log.Println("DB connetion error:" + err.Error())
		}
	case "postgres":
		dbtype := this.DbType
		connectstr := dbtype + "://" + this.DbUser + ":" + url.QueryEscape(this.DbPawd)
		connectstr = connectstr + "@" + this.DbHost + "/" + this.DbName
		connectstr = connectstr + "?sslmode=disable"
		this.engine, err = xorm.NewEngine(dbtype, connectstr)
		if err != nil {
			this.Log.Println("DB connetion error:" + err.Error())
		}
	}
	//cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	//this.engine.SetDefaultCacher(cacher)

	this.engine.SetMaxIdleConns(30)
	this.engine.SetMaxOpenConns(50)
	//this.engine.SetMapper(SnakeMapper{})
	//this.engine.SetMapper(SameMapper{})
	this.engine.ShowErr = true
	this.engine.ShowSQL = true
	fmt.Println("Connetion DB...")
}

// 安全获取数据库链接句柄
func (this *Dbpool) Conn() *xorm.Engine {
	if this.engine == nil {
		this.Inits()
	}
	if err := this.engine.Ping(); err != nil {
		//this.Log.Println("DB ping error:" + err.Error())
		this.Inits()
	}
	return this.engine
}

// 安全关闭数据库链接句柄
func (this *Dbpool) Close() {
	this.engine.Close()
}

// SameMapper implements IMapper and provides same name between struct and
// database table
type SameMapper struct {
}

func (m SameMapper) Obj2Table(o string) string {
	return o
}

func (m SameMapper) Table2Obj(t string) string {
	return t
}

// SnakeMapper implements IMapper and provides name transaltion between
// struct and database table
type SnakeMapper struct {
}

func snakeCasedName(name string) string {
	newstr := make([]rune, 0)
	for idx, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if idx > 0 {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

/*func pascal2Sql(s string) (d string) {
    d = ""
    lastIdx := 0
    for i := 0; i < len(s); i++ {
        if s[i] >= 'A' && s[i] <= 'Z' {
            if lastIdx < i {
                d += s[lastIdx+1 : i]
            }
            if i != 0 {
                d += "_"
            }
            d += string(s[i] + 32)
            lastIdx = i
        }
    }
    d += s[lastIdx+1:]
    return
}*/

func (mapper SnakeMapper) Obj2Table(name string) string {
	return snakeCasedName(name)
}

func titleCasedName(name string) string {
	newstr := make([]rune, 0)
	upNextChar := true

	name = strings.ToLower(name)

	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= chr && chr <= 'z' {
				chr -= ('a' - 'A')
			}
		case chr == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}

func (mapper SnakeMapper) Table2Obj(name string) string {
	return titleCasedName(name)
}
