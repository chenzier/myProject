package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"product/common"
	"product/datamodels"
	"product/encrypt"
	"product/rabbitmq"
	"strconv"
	"sync"
)

// 设置集群地址，最好内网IP
var hostArray = []string{"127.0.0.1", "127.0.0.2"}

// 设置本机IP
var localHost = ""

// 数量控制接口服务器内网IP，或者getone的SLB内网IP
var GetOneIp = "127.0.0.1"

var GetOnePort = "8086"

var port = "8085"

var hashConsistent *common.Consistent

// rabbitmq
var rabbitMqValidate *rabbitmq.RabbitMQ

// 用来存放控制信息
type AccessControl struct {
	//用来存放用户想要存放的信息
	sourcesArray map[int]interface{}
	sync.RWMutex //引入锁是为了保证sourcesArray在高并发情景下的安全
}

var accessControl = &AccessControl{sourcesArray: make(map[int]interface{})}

// 获取指定的数据
func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	m.RWMutex.Unlock()
	data := m.sourcesArray[uid]
	return data
}

// 设置记录
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourcesArray[uid] = "hello imooc"
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	//获取用户UID
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	//采用一致性hash算法，根据用户ID，判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	//判断是否为本机
	if hostRequest == localHost {
		//执行本机数据读取和校
		return m.GetDataFromMap(uid.Value)
	} else {
		//不是本机，本机充当代理访问数据返回结果
		return m.GetDataFromOtherMap(hostRequest, req)
	}
}

// 获取本机map，并且处理业务逻辑，返回的结果为bool类型
func (m *AccessControl) GetDataFromMap(uid string) (isok bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)

	//执行判断逻辑
	if data != nil {
		return false
	}
	return
}

func (m *AccessControl) GetDataFromOtherMap(host string, request *http.Request) bool {
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, request)
	if err != nil {
		return false
	}
	//判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

// 模拟请求
func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	//获取uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	//获取sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}
	//模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	//添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	//获取返回的结果
	response, err = client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	//判断状态
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

func Check(w http.ResponseWriter, r *http.Request) {
	//执行正常业务逻辑
	fmt.Println("执行check! ")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println(productString)

	//获取用户cookie
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	//1.分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
	}
	//2.获取数量控制权限，防止秒杀出现超卖现象
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	//判断数量控制接口请求状态
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			//整合下单
			//1.获取商品ID
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			//2.获取用户ID
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			//3.创建消息体
			message := datamodels.NewMessage(userID, productID)
			//类型转化
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			//4.生产消息
			err = rabbitMqValidate.PublishSimple(string((byteMessage)))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}

	}
	w.Write([]byte("false"))
}

// 统一验证拦截器，每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证")
	uidCookie, err1 := r.Cookie("uid")
	if err1 != nil {
		fmt.Println("获取失败")
	}
	fmt.Println(uidCookie)
	//添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return errors.New("验证成功")
}

// 身份校验函数
// 从Cookie中 获取 uid 和 加密的uid——sign
// 然后检测是否匹配
func CheckUserInfo(r *http.Request) error {
	//获取Uid，cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户UID Cookie获取失败")
	}
	//fmt.Println(uidCookie)
	//获取用户加密串
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 Cookie获取失败")
	}
	//fmt.Println(signCookie.Value)

	//对信息进行解密
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改")
	}

	//signString1 := "root"
	//signByte := []byte(signString1)
	//fmt.Println(string(signByte))

	//fmt.Println("结果比对")
	//fmt.Println("用户ID:" + uidCookie.Value)
	//fmt.Println("解密后用户ID:" + string((signByte)))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	//return errors.New("身份验证失败")
	return nil
}

// 自定义逻辑判断
func checkInfo(checkStr string, signStr string) bool {
	fmt.Println(checkStr, signStr)
	if checkStr == signStr {
		return true
	}
	return false

}

func main2() {
	//负载均衡器设置
	//采用一致性哈希算法
	hashConsistent = common.NewConsistent()
	//采用一致性hash算法，添加节点
	for _, v := range hostArray {
		hashConsistent.Add(v)

	}

	localIP, err := common.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIP
	fmt.Println(localHost)
	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("imoocProduct")
	defer rabbitMqValidate.Destory()

	//1.过滤器
	filter := common.NewFilter()
	//注册拦截器

	filter.RegisterFilterUrl("/check", Auth)
	filter.RegisterFilterUrl("/checkRight", Auth)
	//2.启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))

	//启动服务
	fmt.Println("启动服务成功")
	http.ListenAndServe(":8084", nil)
	//fmt.Println("启动服务成功")
}
