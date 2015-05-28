----------------
freehttp

一个快速将类和子方法转换成HTTP接口

主要是对net/http的一个反射封装，便于使用

	server.ResponseWriter = http.ResponseWriter + 帮助方法
	// 也就是 server.ResponseWriter.ResponseWriter = http.ResponseWriter
	// 也就是 server.ResponseWriter.* = 帮助方法

	server.Request = http.Request + 帮助方法
	// 也就是 server.Request.Request = *http.Request
	// 也就是 server.Request.* = 帮助方法

	// 子方法名任意但是要首字母大写，Go权限设计
	// 任意使用 server.ResponseWriter 或 server.Request 
	// 当然，也可以不使用，直接 func MyFunc() {} 即可
	// 使用任意也可以比如 func MyFunca(w server.ResponseWriter) {}
	// 使用任意也可以比如 func MyFunca(r server.Request) {}

	// Json 设计
	Request // 内绑定了很多读取的帮助方法
	ResponseWriter // 内绑定了很多写入的帮助方法

----------------

	package main

	import (
		"fmt"
		"server" // 导入 server 包
	)

	// 随便定义一个类
	type Web struct {
	}

	// 给类随便定义一个方法，可以任意使用 server.ResponseWriter 或 server.Request 数据
	func (this *Web) Print(w server.ResponseWriter, r server.Request) {
		// 回写了一个数据
		w.ResponseWriter.Write([]byte("print"))
	}

	// ...
	func (this *Web) Hello(w server.ResponseWriter, r server.Request) {
		// ...
		w.ResponseWriter.Write([]byte("hello"))
	}

	// ...
	func (this *Web) Ifconfig(r server.Request) error {
		// 这里只是简单的设计，返回一个错误时，会打印日志
		return fmt.Errorf("error")
	}

	// ResponseWriter 内捆绑着很多帮助方法，比如发送JSON
	func (this *Web) Json(w server.ResponseWriter) {
		m := make(map[string]interface{})
		m["baidu"] = "www.baidu.com"
		w.WriterJsonLine(m)
		w.WriterJsonIndent(m)
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

