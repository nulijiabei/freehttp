package freehttp

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
)

type Domain struct {
	name    string
	rcvr    reflect.Value
	typ     reflect.Type
	methods map[string]reflect.Method
}

type R struct {
	name   string
	method string
}

type Service struct {
	domain map[string]*Domain
	router map[string]R
}

// 创建 Service
func NewService() *Service {
	service := new(Service)
	service.domain = make(map[string]*Domain)
	service.router = make(map[string]R)
	return service
}

// 添加路由
func (this *Service) Router(path string, method interface{}) {
	name := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
	a1 := regexp.MustCompile(`\(\*.*\)`).FindAllString(name, -1)[0]
	a2 := regexp.MustCompile(`\).*-`).FindAllString(name, -1)[0]
	this.router[path] = R{a1[2 : len(a1)-1], a2[2 : len(a2)-1]}
}

// 启动服务
func (this *Service) Start(port string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// ListenAndServe(addr string, handler Handler)
	return http.ListenAndServe(port, this)
}

// 将类及方法注册到 FreeHttp
func (this *Service) Register(rcvr interface{}) {
	// Domain 初始化
	do := new(Domain)
	do.methods = make(map[string]reflect.Method)
	// 获取类的反射类型
	do.typ = reflect.TypeOf(rcvr)
	// 获取类的反射值
	do.rcvr = reflect.ValueOf(rcvr)
	// 获取类名 ...
	do.name = do.rcvr.Elem().Type().Name()
	// 遍历函数列表 ...
	for m := 0; m < do.typ.NumMethod(); m++ {
		// 获取函数
		method := do.typ.Method(m)
		// mtype := method.Type
		mname := method.Name
		// 以函数名为KEY存储函数 ...
		do.methods[mname] = method
	}
	// Reg to Service
	this.domain[do.name] = do
}

// 检查错误返回
func (this *Service) CheckError(err interface{}, name string) {
	if err != nil {
		fmt.Errorf("%s exception -> %s", name, err.(error).Error())
	}
}

// 包含 ServeHTTP(ResponseWriter, *Request) 函数符合 Handler 接口
func (this *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 是否匹配 ...
	status := false
	// 重要：创建FreeHttp并初始化 ...
	freehttp := NewFreeHttp(w, r)
	// 匹配路径
	if r, ok := this.router[r.URL.Path]; ok {
		// 匹配域
		if do, ok := this.domain[r.name]; ok {
			// 匹配函数
			if method, ok := do.methods[r.method]; ok {
				// 标记匹配
				status = true
				// 创建参数集
				value := make([]reflect.Value, method.Type.NumIn())
				// 第一个值为类反射值
				value[0] = do.rcvr
				// 遍历注册函数参数类型
				for n := 1; n < method.Type.NumIn(); n++ {
					// 获取参数类型 -> String
					inType := method.Type.In(n).String()
					/*
						传入输入参数时如果类型为指针时可以为空
						传入输入参数时如果类型非指针不可以为空
					*/
					switch inType {
					case "*freehttp.FreeHttp":
						value[n] = reflect.ValueOf(freehttp)
					case "*freehttp.Request":
						value[n] = reflect.ValueOf(freehttp.SuperRequest)
					case "*freehttp.ResponseWriter":
						value[n] = reflect.ValueOf(freehttp.SuperResponseWriter)
					case "freehttp.Json":
						value[n] = reflect.ValueOf(freehttp.SuperRequest.ReadJson())
					case "freehttp.ContentType":
						value[n] = reflect.ValueOf(freehttp.SuperRequest.ReadContentType())
					case "freehttp.Stream":
						value[n] = reflect.ValueOf(freehttp.SuperRequest.ReadStream())
					default:
						fmt.Printf("unsupported in type: %s\n", inType)
					}
				}
				// 调用函数 ... 传参 ... 并获取返回值 ...
				returnValues := method.Func.Call(value)
				// 遍历返回值 ...
				for t := 0; t < method.Type.NumOut(); t++ {
					// 获取返回值类型
					reType := method.Type.Out(t).String()
					// 获取返回值内容
					content := returnValues[t].Interface()
					// 返回 nil 触发 ...
					if reType != "error" && content == nil {
						fmt.Printf("%s out value is null -> %s\n", r.name, reType)
						continue
					}
					/*
						输出参数其返回值只能为Error且需要手动捕捉输出
					*/
					switch reType {
					case "freehttp.HttpStatus":
						freehttp.SuperResponseWriter.WriteHeader(content)
					case "freehttp.ContentType":
						freehttp.SuperResponseWriter.SetContentType(content)
					case "freehttp.Json":
						err := freehttp.SuperResponseWriter.WriterJson(content)
						this.CheckError(err, r.name)
					case "freehttp.JsonIndent":
						err := freehttp.SuperResponseWriter.WriterJsonIndent(content)
						this.CheckError(err, r.name)
					case "freehttp.Stream":
						err := freehttp.SuperResponseWriter.WriterStream(content)
						this.CheckError(err, r.name)
					case "freehttp.File":
						freehttp.SuperResponseWriter.WriterFile(content)
					case "freehttp.Redirect":
						freehttp.Redirect(content)
					case "error":
						this.CheckError(content, r.name)
					default:
						fmt.Printf("unsupported out type: %s\n", reType)
					}
				}
				// Auto Close
				freehttp.SuperRequest.Close()
			}
		}
	}
	// 不匹配则 NotFound ...
	if !status {
		http.NotFound(w, r)
	}
}
