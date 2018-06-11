// socket服务
/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/

package socket

import (
	"errors"
	"net"
	"os"
	"time"
)

func NewSocketServer(mode, host, port string) *SocketServer {
	return &SocketServer{
		mode:         mode,
		host:         host,
		port:         port,
		poolNum:      MAX_SERVER_CONN_NUM,
		readDeadline: DEFAULT_SERVER_READ_OUT * time.Second,
		RouteMap:     make(map[string]RouterInterface),
	}
}

//SocketServer结构体
type SocketServer struct {
	mode         string        //协议：tcp/udp
	host         string        //主机ip
	port         string        //端口
	poolNum      int           //
	readDeadline time.Duration //
	RouteMap     map[string]RouterInterface

	//stoptag string   //发送数据的终止符
}

func (this *SocketServer) AddModulesRoute(k string, r RouterInterface) {
	this.RouteMap[k] = r
}

func (this *SocketServer) getModulesRoute(module string) (RouterInterface, error) {
	r, ok := this.RouteMap[module]
	if ok {
		return r, nil
	}
	return r, errors.New("Error::not have this modules![" + module + "]")

}

func (this *SocketServer) SetPoolNum(i int) *SocketServer {
	this.poolNum = i
	return this
}

func (this *SocketServer) GetPoolNum() int {
	return this.poolNum
}

func (this *SocketServer) SetDefaultDeadline(t time.Duration) *SocketServer {
	this.readDeadline = t
	return this
}

//开始启动socket服务
func (this *SocketServer) Run() {
	listener, err := net.Listen(this.mode, this.host+":"+this.port)
	if err != nil {
		//log("socket error : " + err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	log("socket running ...")
	conn_chan := make(chan net.Conn)
	//预开启子协程
	for i := 0; i < this.poolNum; i++ {
		go func() {
			for conn := range conn_chan {
				log("get request ... ")
				connSocket := new(ConnSocket)
				connSocket.Conn = conn
				connSocket.SocketServer = this
				connSocket.Doing()
				connSocket = nil
				log("over request ... ")
			}
		}()
	}
	//开始监听
	for {
		conn, err := listener.Accept()
		if err != nil {
			log("Error accept:" + err.Error())
			return
		}
		//通过信道，转交给预开启的子协程处理，达到非阻塞监听处理请求
		conn_chan <- conn
	}

}
