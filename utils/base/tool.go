/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
//

package base

import (
	//"fmt"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	//"strconv"
	"os"
	"runtime"
	"strings"
)

//获取根目录的绝对路径
func RootPath() string {
	var tempSlice []string

	if runtime.GOOS == "windows" {
		//windows系统
		root, _ := exec.LookPath(os.Args[0])
		rootPath, _ := filepath.Abs(root)
		tempSlice = strings.Split(rootPath, `\`)
	} else {
		//其他系统
		root, _ := exec.LookPath(os.Args[0])
		rootPath, _ := filepath.Abs(root)
		tempSlice = strings.Split(rootPath, "/")
	}
	return strings.Join(tempSlice[0:len(tempSlice)-1], "/")
}

func PanicIf(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

//获得程序执行的根目录
func GetPwd() string {
	pwd, _ := os.Getwd()
	selfdir := path.Dir(pwd)
	return selfdir
}
