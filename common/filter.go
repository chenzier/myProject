package common

import (
	"net/http"
	"strings"
)

// 声明一个新的数据类型 这是一个函数类型
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// 拦截器结构体
type Filter struct {
	//用来存储需要拦截的URI
	filterMap map[string]FilterHandle //用于存储需要拦截的 URI 和对应的拦截器函数。
}

// Filter初始化函数
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

// 注册拦截器
func (f *Filter) RegisterFilterUrl(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// 根据Uri获取对应的handle
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

// 声明新的函数类型
type WebHandle func(rw http.ResponseWriter, req *http.Request)

// 执行拦截器，返回函数类型
// 输入WebHandle类型的函数，返回一个 函数类型
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) { //自定义一个函数
		for path, handle := range f.filterMap {
			if strings.Contains(r.RequestURI, path) { //它首先检查请求的 URI 是否存在于 filterMap 中
				//执行拦截业务逻辑
				err := handle(rw, r) // 如果存在则执行对应的拦截器函数，然后再执行正常的 Web 处理函数。
				//这种实现方式允许你在处理特定 URI 的 HTTP 请求之前，执行一些额外的逻辑，例如身份验证、权限控制等。
				//通过注册拦截器函数和 URI 的关联，你可以灵活地控制哪些请求需要经过拦截器处理
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				//跳出循环
				break
			}
		}
		//执行正常注册的函数
		webHandle(rw, r)
	}
}
