package freehttp

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Server Json HTTP
// http://127.0.0.1:8080/MyStructName/MyFuncName
type Service struct {
	conf    *INI
	name    string
	rcvr    reflect.Value
	typ     reflect.Type
	methods map[string]reflect.Method
}

// 创建 Server 其中 def 为路径处理函数 nil 则使用默认
func NewService(rcvr interface{}) *Service {
	service := new(Service)
	service.methods = make(map[string]reflect.Method)
	if err := service.Register(rcvr); err != nil {
		panic(err)
	}
	return service
}

// 启动服务
func (this *Service) Start(port string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// ListenAndServe(addr string, handler Handler)
	return http.ListenAndServe(port, this)
}

// 将类及方法注册到FreeHttp
func (this *Service) Register(rcvr interface{}) error {
	// 获取类的反射类型
	this.typ = reflect.TypeOf(rcvr)
	// 获取类的反射值
	this.rcvr = reflect.ValueOf(rcvr)
	// 获取类名 ...
	this.name = this.rcvr.Elem().Type().Name()
	// 遍历函数列表 ...
	for m := 0; m < this.typ.NumMethod(); m++ {
		// 获取函数
		method := this.typ.Method(m)
		// mtype := method.Type
		mname := method.Name
		// 以函数名为KEY存储函数 ...
		this.methods[mname] = method
	}
	return nil
}

// 初始化配置文件
func (this *Service) Config(path string) {
	this.conf = NewINI(path)
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
	// 统一小写
	r.URL.Path = strings.ToLower(r.URL.Path)
	// 遍历注册函数
	for name, method := range this.methods {
		// 注册类与函数名转换为URL路径
		url := strings.ToLower(fmt.Sprintf("/%s/%s", this.name, name))
		// 匹配 URL
		if url == r.URL.Path {
			// 标记匹配
			status = true
			// 创建参数集
			value := make([]reflect.Value, method.Type.NumIn())
			// 第一个值为类反射值
			value[0] = this.rcvr
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
				case "*freehttp.INI":
					// panic(fmt.Sprintf("Use a non-initialized type: %s", inType))
					value[n] = reflect.ValueOf(this.conf)
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
					fmt.Printf("%s out value is null -> %s\n", name, reType)
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
					this.CheckError(err, name)
				case "freehttp.JsonIndent":
					err := freehttp.SuperResponseWriter.WriterJsonIndent(content)
					this.CheckError(err, name)
				case "freehttp.Stream":
					err := freehttp.SuperResponseWriter.WriterStream(content)
					this.CheckError(err, name)
				case "freehttp.File":
					freehttp.ServeFiles(content)
				case "freehttp.Redirect":
					freehttp.Redirect(content)
				case "error":
					this.CheckError(content, name)
				default:
					fmt.Printf("unsupported out type: %s\n", reType)
				}
			}
		}
	}
	// 不匹配则 NotFound ...
	if !status {
		http.NotFound(w, r)
	}
}
