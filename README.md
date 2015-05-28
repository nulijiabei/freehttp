----------------
freehttp

一个快速将类和子方法转换成HTTP接口

主要是对net/http的一个反射封装，便于使用

	核心：
		
		server.Request        封装于  http.Request
		server.ResponseWriter 封装于  http.ResponseWriter

	衍生帮助方法:
	
		server.Request.*         基于对 http.Request        的自定义帮助方法
		server.ResponseWriter.*  基于对 http.ResponseWriter 的自定义帮助方法
		
	继承动态输出类型
	
		// 错误类型
		error
		
		例如: func (this *MyStruct) MyFunc(...) error
		
		介绍: 返回 error 会作为日志输出
	
	衍生动态输出类型:
	
		// Json 普通格式
		Json 原型 map[string]interface{}
		
		// Json 排版格式
		JsonIndent 原型 map[string]interface{}
	
		例如: func (this *MyStruct) MyFunc() server.Json {}
		
		例如: func (this *MyStruct) MyFunc() server.JsonIndent {}
		
		介绍：返回类型为 server.Json 或 server.JsonIndent 数据会以 Json 方法回写
			
	继承动态输入类型:
	
		// 包装 http.ResponseWriter
		type ResponseWriter struct { ResponseWriter http.ResponseWriter }
		
		// 包装 *http.Request
		type Request struct { Request *http.Request }
	
		例如: func (this *MyStruct) MyFunc(w server.ResponseWriter, r server.Request) {}
		
		介绍: 输入类型为 server.ResponseWriter 或 server.Request 
		
		则: 可以使用 *http.Request 和 http.ResponseWriter 集继承与自定义方法
			
	衍生动态输入类型:
	
		// Body
		Body 原型 []byte
		
		// Json Body
		BodyJson 原型 map[string]interface{}
		
		例如: func (this *MyStruct) MyFunc(body server.Body, bodyJson server.BodyJson) {}
		
		介绍: 输入类型为 server.Body 或 server.BodyJson 时，自动传入 Body 全部数据 或 转为Json传入

----------------

	package main

	import (
		"fmt"
		"server" // 导入 server 包
	)

	// 随便定义一个类
	type Web struct {
	}

	// 任意传入了 server.ResponseWriter 和 server.Request
	func (this *Web) ReadWrite(w server.ResponseWriter, r server.Request) {
		// r.Request.PostForm
		w.ResponseWriter.Write([]byte("print"))
	}
	
	// 返回了一个 server.Json
	func (this *Web) WriteJson() server.Json {
		m := make(map[string]interface{})
		m["baidu"] = "www.baidu.com"
		return m
	}
	
	// 任意传入了 server.Body 和 server.BodyJson 返回了 error
	func (this *Web) Hello(body server.Body, bodyJson server.BodyJson) error {
		fmt.Println(body)
		fmt.Println(bodyJson)
		return fmt.Errorf("...")
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

