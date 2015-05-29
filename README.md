----------------
freehttp

<a href="https://godoc.org/github.com/nulijiabei/freehttp"><img src="https://godoc.org/github.com/nulijiabei/freehttp?status.svg" alt="GoDoc"></a>

一个快速将类和子方法转换成HTTP接口

主要是对net/http的一个反射封装，便于使用

----------------

安装

	go get github.com/nulijiabei/freehttp
	
----------------

	核心：
		
		freehttp.Request        封装于  http.Request
		freehttp.ResponseWriter 封装于  http.ResponseWriter

	衍生帮助方法:
	
		freehttp.Request.*         基于对 http.Request        的自定义帮助方法
		freehttp.ResponseWriter.*  基于对 http.ResponseWriter 的自定义帮助方法
		
	继承输入类型:
	
		// 不限制使用 http.Request 内所有方法, 只在基础上进行扩展
		freehttp.Request.Request = *http.Request

	继承输出类型:
	
		// 不限制使用 http.ResponseWriter 内所有方法, 只在基础上进行扩展
		freehttp.ResponseWriter.ResponseWriter = http.ResponseWriter 
	
	衍生输入类型:
	
		// Body
		freehttp.Body		对应方法 -> freehttp.Request.ReadAllBody()
		
		// Json Body
		freehttp.BodyJson	对应方法 -> freehttp.Request.ReadAllBodyJson()
		
	衍生输出类型:
		
		// Json 普通格式
		freehttp.Json			对应方法 -> freehttp.ResponseWriter.WriterJson()
		
		// Json 排版格式
		freehttp.JsonIndent		对应方法 -> freehttp.ResponseWriter.WriterJsonIndent()
			
		// HTTP Status
		freehttp.Status			对应方法 -> freehttp.ResponseWriter.WriteHeader()


----------------

	使用方法:
	
		不限制Struct的类型及名称
		不限制Struct所属函数类型及名称(函数首字母大写)
		输出参数和输出只能使用继承及衍生的类型作为参数
		
		可以将任意的输入类型作为输入参数使用，任意组合
		可以将任意的输出类型作为输出参数使用，任意组合
		框架会通过反射机制识别类型并调用对应方法执行.
	
	例如:
	
		func (this *MyStruct) MyFunc(
			// 这里的传入参数只能使用继承或衍生输入类型
		)   // 这里的返回参数只能使用继承或衍生输出类型
		{}

----------------

	// 例
	package main

	import (
		"fmt"
		"freehttp" // 导入 freehttp 包
	)

	// 随便定义一个类
	type Web struct {
	}

	// 输入参数，输出参数任意组合
	func (this *Web) ReadWrite(w *freehttp.ResponseWriter, r *freehttp.Request) error {
		// r.Request.PostForm
		w.ResponseWriter.Write([]byte("print"))
		return fmt.Errorf("...")
	}
	
	// 输入参数，输出参数任意组合
	func (this *Web) WriteJson(r *freehttp.Request) freehttp.Json {
		m := make(map[string]interface{})
		m["baidu"] = "www.baidu.com"
		return m
	}
	
	// 输入参数，输出参数任意组合
	func (this *Web) Hello(w *freehttp.ResponseWriter, r *freehttp.Request, 
		body freehttp.Body, bodyJson freehttp.BodyJson) (freehttp.Status, error) {
		fmt.Println(body)
		fmt.Println(bodyJson)
		return 404, fmt.Errorf("...")
	}

	// 主要
	func main() {

		// 创建 Server
		s := server.NewServer()

		// 传入 Web 类
		if err := s.Register(new(Web)); err != nil {
			fmt.Println(err)
		}
	
		// 启动监听端口
		s.Start(":8080")

	}

