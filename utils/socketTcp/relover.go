// 接受请求，分发中转到模块
/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/

package socketTcp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type ConnSocket struct {
	Conn         *net.TCPConn
	SocketServer *SocketServer
	isPconnect   bool
	isClose      bool
}

//接受请求，分发处理
func (this *ConnSocket) Doing() {
	defer func() {
		if err := recover(); err != nil {
			log(fmt.Sprint(err))
			//状态码500，服务内部错误
			this.ConnWrite(500, fmt.Sprint(err))
			this.close()
			return
		}
		this.close()
	}()
	this.isPconnect = false
	this.isClose = false
	//this.Conn.SetReadDeadline(time.Now().Add(this.SocketServer.readDeadline)) //设置超时时间（读超时，自动关闭）
	//-----------------------------------------------------------------------------------------------------------
	var readLen int              //需要读取的数据大小
	var readBufferLen int = 1024 //每次读取的分片单元大小
	newRequest := true           //是否是一次新的请求
	requestTime := 0
	for {
		if this.isClose {
			return
		}
		if newRequest {
			//----------------------------------
			//读取“请求头信息”，
			//大小为10个字节，
			//包含着“本次请求”总共将接受到多少个字节。
			//----------------------------------
			headerinfo := make([]byte, 0)
			var e error
			if headerinfo, e = this.readLastData(headerinfo, 15, 1); e != nil { //如果请求头错误，执行完一次直接断开
				if this.isPconnect && requestTime > 0 {
					//如果是长连接，执行完一次请求后，则继续监听
					fmt.Println(e.Error(), "====================")
					break
				} else {
					//如果不是长连接，执行完一次请求后，则直接断开
					break
				}
			} else { //
				requestTime++
				headerFirst := this.trim(headerinfo[0:10])
				//"长连接"请求
				if this.isPconnect == false {
					if headerFirst == "pconnect" {
						keepAliveTime, _ := strconv.Atoi(this.trim(headerinfo[10:15]))
						if keepAliveTime < 10 {
							keepAliveTime = 10
						}
						//开启长连接的心跳机制
						if err := this.Conn.SetKeepAlive(true); err != nil {
							this.ConnWrite(507, err.Error())
							return
						}
						//设置心跳周期
						if err := this.Conn.SetKeepAlivePeriod(time.Duration(keepAliveTime) * time.Second); err != nil {
							this.ConnWrite(508, err.Error())
							return
						}
						this.ConnWrite(200, "pconnect successfull!")
						this.isPconnect = true
						continue
					}
				}
				//非"长连接"
				readLen, _ = strconv.Atoi(headerFirst)
				newRequest = false
				//
				if this.isPconnect == false { //
					timeOutStr := this.trim(headerinfo[10:15])
					//非"长连接"，设置超时时间，以保护线程能及时释放回收
					timeOut, _ := strconv.Atoi(timeOutStr)
					if timeOut <= 0 {
						this.Conn.SetDeadline(time.Now().Add(this.SocketServer.readDeadline)) //默认超时时间
					} else if timeOut > 2000 {
						this.Conn.SetDeadline(time.Now().Add(time.Duration(2000) * time.Second)) //最大超时时间
					} else {
						this.Conn.SetDeadline(time.Now().Add(time.Duration(timeOut) * time.Second)) //自定义超时时间
					}
				}
			}
		} else {
			if readLen <= 0 {
				//状态码501，请求错误，数据为空
				log("[warn]" + "api of request to socket's reading get no data!")
				this.ConnWrite(502, "api of request to socket's reading get no data!")
				return
			}
			//----------------------------------
			//读取“本次请求”剩下的数据---
			//“请求头信息”之后的所有数据。
			//当数据超过1024（上面设置的“分片单元”）
			//时，将分片读取。
			//----------------------------------
			var readErr error
			var datainfo []byte         //每次read的缓冲大小
			readTimes := 0              //读取的次数
			readTimesMaxLimit := 10000  //读取次数的限制（主要是为了保护，防止死循环的出现）
			readData := make([]byte, 0) //读取到是数据（累加的）
			for readTimes < readTimesMaxLimit {
				readTimes++
				lastLen := readLen - readBufferLen
				if lastLen <= 0 {
					//----------------------------------
					//当剩下的数据字节长度小于1024，
					//则读完后返回结果。。。。
					//----------------------------------
					datainfo = make([]byte, readLen)
					n, err := this.Conn.Read(datainfo)
					if err != nil {
						readErr = err
						break
					}
					lostLen := readLen - n                        // ++
					readData = append(readData, datainfo[0:n]...) // ++
					if lostLen > 0 {
						if readData, err = this.readLastData(readData, lostLen, 1); err != nil {
							readErr = err
							break
						}
					}
					break
				} else {
					//----------------------------------
					//当剩下的数据字节长度大于1024，
					//则读取1024个字节。。。。
					//----------------------------------
					datainfo = make([]byte, readBufferLen)
					n, err := this.Conn.Read(datainfo)
					if err != nil {
						readErr = err
						break
					}
					lostLen := readBufferLen - n                  // ++
					readData = append(readData, datainfo[0:n]...) // ++
					readLen = lastLen + lostLen
				}
			}
			if readErr != nil {
				//状态码501，请求错误，数据为空
				log("[error]" + "api of request to socket's reading fail!")
				this.ConnWrite(503, "api of request to socket's reading fail!")
				return
			} else {
				this.sendsData(readData)
				newRequest = true
			}
		}
	}
}

//递归读取剩下的丢失的数据，直至读到为止。
func (this *ConnSocket) readLastData(readData []byte, lastLen int, readTime int) ([]byte, error) {
	datainfo := make([]byte, lastLen)
	if readTime > 1000 {
		return readData, errors.New("read lost data more than 1000 times!")
	}
	if n, err := this.Conn.Read(datainfo); err != nil {
		return readData, err
	} else {
		lostLen := lastLen - n
		readData = append(readData, datainfo[0:n]...)
		if lostLen > 0 {
			readTime = readTime + 1
			return this.readLastData(readData, lostLen, readTime)
		} else {
			return readData, nil
		}
	}

}

//处理接受到的数据，并最终返回结果。
func (this *ConnSocket) sendsData(get []byte) {
	defer func() {
		if err := recover(); err != nil {
			log(fmt.Sprint(err))
			//状态码500，服务内部错误
			log("[error]" + fmt.Sprint(err))
			this.ConnWrite(504, fmt.Sprint(err))
			return
		}
	}()
	//专门为ping做一个状态码200，响应链接是存在的
	getStr := this.trim(get[0:60])
	if getStr == "ping" {
		this.ConnWrite(200, "ok")
		return
	}
	//解析请求
	route := strings.Split(getStr, ".")
	if len(route) != 3 {
		//状态码502，请求接口地址的格式不符合规范
		log("[warn]" + getStr + " ,the api to address's format is error,parsing fail!") //记录错误日志
		this.ConnWrite(505, "the api to address's format is error,parsing fail!")
		return
	}
	module := route[0]
	controller := route[1]
	action := route[2]
	//解析成功
	if modelsRoute, err := this.SocketServer.getModulesRoute(module); err == nil {
		//交由模块进行处理
		d := get[60:len(get)]
		code, rsData := modelsRoute.Route(controller, action, d)
		this.ConnWrite(code, rsData)
		return
	} else {
		//状态码503，请求的接口不存在
		log("[warn]" + err.Error()) //记录错误日志
		this.ConnWrite(506, err.Error())
		return
	}
}

//返回结果
func (this *ConnSocket) ConnWrite(code int, rs string) {
	rs = string(this.fiexLeninfo(strconv.Itoa(code), 5)) + rs
	first := this.fiexLeninfo(strconv.Itoa(len(rs)), 10)
	sendbyte := []byte(rs)
	send := append(first, sendbyte...)
	if _, err := this.Conn.Write(send); err != nil {
		errlog := "[error]" + "the response of api for socket to write error:" + err.Error()
		log(errlog) //记录错误日志
		this.isClose = true
	}
}

//返回的“信息头”---包含将要返回多少字节给客户端
func (this *ConnSocket) fiexLeninfo(str string, l int) []byte {
	minsend := make([]byte, l)
	first := append([]byte(str), minsend...)
	first = first[0:l]
	return first
}

func (this *ConnSocket) trim(b []byte) string {
	emptyByte := make([]byte, 1)
	return strings.TrimSpace(strings.Trim(string(b), string(emptyByte)))
}

func (this *ConnSocket) close() {
	this.Conn.Close()
	this = nil
}
