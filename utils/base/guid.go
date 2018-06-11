/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
//
package base

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const SYS_GUID_BINGE int64 = 1e2

type Guid struct {
	Recodefile string
	PrefixNum  int64
	autoid     int64
	idsChan    chan int64
}

func (gid *Guid) Run() {
	gid.idsChan = make(chan int64, SYS_GUID_BINGE)
	go gid.creatId()
}

func (gid *Guid) creatId() {
	recodePath, _ := filepath.Abs(gid.Recodefile)
	var tempSlice []string
	if runtime.GOOS == "windows" {
		tempSlice = strings.Split(recodePath, `\`)

	} else {
		tempSlice = strings.Split(recodePath, "/")
	}
	path := strings.Join(tempSlice[0:len(tempSlice)-1], "/")
	a := tempSlice[len(tempSlice)-1]
	h := md5.New()
	h.Write([]byte(a))
	b := hex.EncodeToString(h.Sum(nil))
	fmt.Println(path)
	// 创建日志目录
	os.MkdirAll(path, 0755)
	for {
		aFile := path + "/" + a
		bFile := path + "/" + b
		var id int64
		if file_a, err := os.Open(aFile); err != nil && os.IsNotExist(err) { //当a文件不存在
			if file_b, err := os.Open(bFile); err != nil && os.IsNotExist(err) { //当b文件不存在
				id = gid.initId()
				gid.initAfiles(aFile, id)
			} else { //当b文件存在
				r := gid.readData(file_a)
				if r == "" { //当b文件没有数据（可能上次写入失败）
					id = gid.initId()
					gid.initAfiles(aFile, id)
				} else { //当b文件有数据
					if temp, err := strconv.Atoi(r); err != nil {
						id = gid.initId()
						gid.initAfiles(aFile, id)
					} else {
						id = int64(temp) + 1
						gid.initAfiles(aFile, id)
					}
				}
				file_b.Close()
				os.Remove(bFile)
			}
		} else {
			r := gid.readData(file_a)
			if r == "" { //当b文件没有数据（可能上次写入失败）
				id = gid.initId()
				gid.initAfiles(bFile, id)
			} else { //当b文件有数据
				if temp, err := strconv.Atoi(r); err != nil {
					id = gid.initId()
					gid.initAfiles(bFile, id)
				} else {
					id = int64(temp) + 1
					gid.initAfiles(bFile, id)
				}
			}
			file_a.Close()
			os.Remove(aFile)
			os.Rename(bFile, aFile)
		}
		fmt.Println(id)
		gid.idsChan <- id
	}
}

func (gid *Guid) Getid(prefix int64) int64 {
	rs := <-gid.idsChan
	if prefix > 0 {
		s := strconv.Itoa(int(prefix)) + strconv.Itoa(int(rs))
		rs, _ = strconv.ParseInt(s, 10, 64)
	}
	return rs
}

func (gid *Guid) initAfiles(path string, id int64) {
	if temp_a, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644); err != nil {
		panic(err.Error())
	} else {
		defer temp_a.Close()
		_, err := temp_a.WriteString(strconv.Itoa(int(id)))
		if err != nil {
			panic(err.Error())
		}
	}
}

func (gid *Guid) readData(fin *os.File) string {
	buf := make([]byte, 1024)
	var str string
	str = ""
	for {
		n, _ := fin.Read(buf)
		if 0 == n {
			break
		}
		str = str + string(buf[:n])
	}
	return str
}

func (gid *Guid) initId() int64 {

	return (time.Now().Unix() + gid.PrefixNum*10000000000) * 100
}
