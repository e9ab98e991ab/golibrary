// socket服务
/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/

package socket

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const CLIENT_TIME_OUT_DEFAULT = 10

func NewConnClient(mode, host, port string) *ConnClientPool {
	pool := &ConnClientPool{
		mode:         mode,
		host:         host,
		port:         port,
		poolNum:      MAX_CLIENT_CONN_NUM,
		readDeadline: DEFAULT_CLIENT_READ_OUT * time.Second,
	}
	pool.run()
	return pool
}

func NewPconnectClient(mode, host, port string) *PconnectClientPool {
	pool := &PconnectClientPool{
		mode:         mode,
		host:         host,
		port:         port,
		poolNum:      MAX_CLIENT_CONN_NUM,
		readDeadline: DEFAULT_CLIENT_READ_OUT * time.Second,
	}
	pool.run()
	return pool
}

type ClientPooler interface {
	GetMode() string
	GetHost() string
	GetPort() string
	Send(address string, data interface{}, outTime ...int) ([]byte, int, error)
	Ping() bool
	Conn() (ClientPooler, error)
	recyc(client *SocketClient)
	push()
}

//-----------------------
// 客户端连接池
//-----------------------
type ConnClientPool struct {
	mode         string             //协议：tcp/udp
	host         string             //主机ip
	port         string             //主机端口
	poolNum      int                //初始启动链接数
	readDeadline time.Duration      //链接超时时间
	clientChan   chan *SocketClient //连接池（信道）
	//selfClientInstance ClientPooler
}

func (this *ConnClientPool) GetMode() string {
	return this.mode
}

func (this *ConnClientPool) GetHost() string {
	return this.host
}

func (this *ConnClientPool) GetPort() string {
	return this.port
}

//发送请求
func (this *ConnClientPool) Send(address string, data interface{}, outTime ...int) ([]byte, int, error) {
	client := this.client()
	defer this.closes(client)
	return client.Send(address, data, outTime...)
}

//检测链路是否连通
func (this *ConnClientPool) Ping() bool {
	return this.client().Ping()
}

//链接服务端
func (this *ConnClientPool) Conn() (ClientPooler, error) {
	client := this.client()
	defer this.recyc(client)
	if _, err := client.Conn(); err != nil {
		return this, err
	}
	return this, nil
}

func (this *ConnClientPool) client() *SocketClient {
	client := <-this.clientChan
	fmt.Println("get client ...")
	return client
}

func (this *ConnClientPool) run() {
	this.clientChan = make(chan *SocketClient, this.poolNum)
	//预开启子协程
	for i := 0; i < this.poolNum; i++ {
		this.clientChan <- &SocketClient{
			ClientPool: this,
			isPconnect: false,
			isClost:    true,
		}
	}
}

//回收连接资源到连接池
func (this *ConnClientPool) recyc(client *SocketClient) {
	if client.isClost {
		client = nil
		this.push()
	} else {
		this.clientChan <- client
	}
}

func (this *ConnClientPool) closes(client *SocketClient) {
	if client.isClost {
		client = nil
	} else {
		client.Close()
		client = nil
	}
	this.push()
}

func (this *ConnClientPool) push() {
	this.clientChan <- &SocketClient{
		ClientPool: this,
		isPconnect: false,
		isClost:    true,
	}
}

//-----------------------
// 客户端连接池
//-----------------------
type PconnectClientPool struct {
	mode         string             //协议：tcp/udp
	host         string             //主机ip
	port         string             //主机端口
	poolNum      int                //初始启动链接数
	readDeadline time.Duration      //链接超时时间
	clientChan   chan *SocketClient //连接池（信道）
}

func (this *PconnectClientPool) GetMode() string {
	return this.mode
}

func (this *PconnectClientPool) GetHost() string {
	return this.host
}

func (this *PconnectClientPool) GetPort() string {
	return this.port
}

//发送请求
func (this *PconnectClientPool) Send(address string, data interface{}, outTime ...int) ([]byte, int, error) {
	return this.client().Send(address, data, outTime...)
}

//检测链路是否连通
func (this *PconnectClientPool) Ping() bool {
	return this.client().Ping()
}

//长连接
func (this *PconnectClientPool) Conn() (ClientPooler, error) {
	client := this.client()
	defer this.recyc(client)
	if _, err := client.Conn(); err != nil {
		return this, err
	}
	return this, nil
}

func (this *PconnectClientPool) client() *SocketClient {
	client := <-this.clientChan
	fmt.Println("get client ...")
	return client
}

func (this *PconnectClientPool) run() {
	this.clientChan = make(chan *SocketClient, this.poolNum)
	//预开启子协程
	for i := 0; i < this.poolNum; i++ {
		this.clientChan <- &SocketClient{
			ClientPool: this,
			isPconnect: true,
			isClost:    true,
		}
	}
}

//回收连接资源到连接池
func (this *PconnectClientPool) recyc(client *SocketClient) {
	if client.isClost {
		client = nil
		this.push()
	} else {
		this.clientChan <- client
	}
}

func (this *PconnectClientPool) closes(client *SocketClient) {
	if client.isClost {
		client = nil
	} else {
		client.Close()
		client = nil
	}
	this.push()
}

func (this *PconnectClientPool) push() {
	this.clientChan <- &SocketClient{
		ClientPool: this,
		isPconnect: true,
		isClost:    true,
	}
}

//-----------------------
//SocketClient结构体
//-----------------------
type SocketClient struct {
	ClientPool ClientPooler
	conn       net.Conn
	isPconnect bool
	isClost    bool
}

func (this *SocketClient) Ping() bool {
	//defer this.Recover()
	return true
}

//连接
func (this *SocketClient) Conn() (*SocketClient, error) {
	if _, e := this.connect(); e != nil {
		return this, e
	}
	return this, nil
}

//长连接
func (this *SocketClient) Pconnect() (*SocketClient, error) {
	//defer this.Recover()
	if conn, e := this.connect(); e != nil {
		return this, e
	} else {
		if this.isPconnect == false {
			send := this.fiexLeninfo("pconnect", 15)
			var rsCode int
			if _, err := conn.Write(send); err != nil {
				rsCode = 600
			} else {
				if responseData, err := this.readResponse(); err != nil {
					log(err.Error())
					rsCode = 601
				} else {
					if len(responseData) > 5 {
						codeByte := responseData[0:5]
						rs := responseData[5:len(responseData)]
						rsCode, _ = strconv.Atoi(this.trim(codeByte))
						log("[info]" + string(rs))
					} else {
						rsCode = 602
					}
				}
			}
			if rsCode != 200 {
				log("[error]" + "code " + strconv.Itoa(rsCode) + "the request of client for socket to pconnect error") //记录错误日志
				this.Close()
			} else {
				this.isPconnect = true
			}
		}
	}
	return this, nil
}

//连接
func (this *SocketClient) connect() (net.Conn, error) {
	var e error
	if this.conn == nil {
		conn, err := net.Dial(this.ClientPool.GetMode(), this.ClientPool.GetHost()+":"+this.ClientPool.GetPort())
		if err != nil {
			fmt.Println(err.Error())
			e = err
			this.conn = nil
			return this.conn, e
		} else {
			e = nil
			this.isClost = false
			this.conn = conn
		}
	}
	return this.conn, e
}

//发送请求数据
func (this *SocketClient) Send(address string, data interface{}, outTime ...int) (rs []byte, rsCode int, e error) {
	defer func() {
		//this.Recover()
	}()
	if conn, err := this.connect(); err != nil {
		e = err
		rsCode = 603
		return
	} else {
		//conn.SetDeadline(time.Now().Add(this.socketClientPool.deadline)) //设置超时时间（对读、写都起作用，超过时间，自动关闭）
		var sendData []byte
		switch d := data.(type) {
		case []byte:
			sendData = d
		case string:
			sendData = []byte(d)
		}
		route := this.fiexLeninfo(address, 60)
		sendData = append(route, sendData...)
		var sendHeadFirst, sendHeadLast []byte
		sendHeadFirst = this.fiexLeninfo(strconv.Itoa(len(sendData)), 10)
		if len(outTime) == 0 {
			sendHeadLast = this.fiexLeninfo(strconv.Itoa(CLIENT_TIME_OUT_DEFAULT), 5)
		} else {
			sendHeadLast = this.fiexLeninfo(strconv.Itoa(outTime[0]), 5)
		}
		send := append(sendHeadFirst, sendHeadLast...)
		send = append(send, sendData...)
		if _, err := conn.Write(send); err != nil {
			e = err
			log("[error]" + "the request of client for socket to write error:" + err.Error()) //记录错误日志
			rsCode = 600
			return
		} else {

			if responseData, err := this.readResponse(); err != nil {
				e = err
				log(err.Error())
				rsCode = 601
				return
			} else {
				if len(responseData) > 5 {
					codeByte := responseData[0:5]
					rs = responseData[5:len(responseData)]
					rsCode, _ = strconv.Atoi(this.trim(codeByte))
					return
				} else {
					e = errors.New("response data len must more than 5")
					rsCode = 602
					return
				}
			}
		}
	}
	return
}

//关闭释放连接
func (this *SocketClient) Close() {
	this.conn.Close()
	this.isClost = true
}

//读取请求响应返回的数据
func (this *SocketClient) readResponse() (readData []byte, readErr error) {
	var readLen int                //需要读取的数据大小
	var readBufferLen int = 1024   //每次读取的分片单元大小
	headerinfo := make([]byte, 10) //“请求头信息”大小

	var conn net.Conn
	if conn, readErr = this.connect(); readErr != nil {
		return
	}
	//----------------------------------
	//读取“请求头信息”，
	//大小为10个字节，
	//包含着“本次请求”总共将接受到多少个字节。
	//----------------------------------
	responseHead := ""
	_, err := conn.Read(headerinfo)
	if err != nil {
		readErr = err
		return
	}
	for _, b := range headerinfo {
		if b == 0 {
			continue
		}
		responseHead = responseHead + string(b)
	}
	log(responseHead)
	readLen, _ = strconv.Atoi(responseHead)
	//----------------------------------
	//读取“本次请求”剩下的数据---
	//“请求头信息”之后的所有数据。
	//当数据超过1024（上面设置的“分片单元”）
	//时，将分片读取。
	//----------------------------------
	var datainfo []byte        //每次read的缓冲大小
	readTimes := 0             //读取的次数
	readTimesMaxLimit := 10000 //读取次数的限制（主要是为了保护，防止死循环的出现）
	readData = make([]byte, 0) //读取到是数据（累加的）
	if readLen <= 0 {
		return
	}
	for readTimes < readTimesMaxLimit {
		readTimes++
		lastLen := readLen - readBufferLen
		if lastLen <= 0 {
			//----------------------------------
			//当剩下的数据字节长度小于1024，
			//则读完后返回结果。。。。
			//----------------------------------
			datainfo = make([]byte, readLen)
			n, err := conn.Read(datainfo)
			if err != nil {
				readErr = err
				break
			}
			lostLen := readLen - n                        // ++
			readData = append(readData, datainfo[0:n]...) // ++
			if lostLen > 0 {
				if readData, err = this.readLastData(readData, lostLen); err != nil {
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
			n, err := conn.Read(datainfo)
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
		log("[error]" + "api of request to socket's reading fail!")
	}
	return
}

//递归读取剩下的丢失的数据，直至读到为止。
func (this *SocketClient) readLastData(readData []byte, lastLen int) ([]byte, error) {
	datainfo := make([]byte, lastLen)
	if conn, err := this.connect(); err != nil {
		return readData, err
	} else {
		n, err := conn.Read(datainfo)
		if err != nil {
			return readData, err
		}
		lostLen := lastLen - n                        // ++
		readData = append(readData, datainfo[0:n]...) // ++
		if lostLen > 0 {
			return this.readLastData(readData, lostLen)
		} else {
			return readData, nil
		}
	}
}

//返回的“信息头”---包含将要返回多少字节给客户端
func (this *SocketClient) fiexLeninfo(str string, l int) []byte {
	minsend := make([]byte, l)
	first := append([]byte(str), minsend...)
	first = first[0:l]
	return first
}

//过滤空字节
func (this *SocketClient) trim(b []byte) string {
	emptyByte := make([]byte, 1)
	return strings.TrimSpace(strings.Trim(string(b), string(emptyByte)))
}
