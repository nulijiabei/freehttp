----------------
freehttp

<a href="https://godoc.org/github.com/nulijiabei/freehttp"><img src="https://godoc.org/github.com/nulijiabei/freehttp?status.svg" alt="GoDoc"></a>

快速将http服务集成到自定义类及子类中，并且通过反射为http接口封装更多人性化的帮助方法 ... 

----------------

为什么要做这个项目

	现在实现一个HTTP-API往往是下面这个方式 ... 
		func sayhelloName(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello!")
		}
		func main() {
			http.HandleFunc("/", sayhelloName) 
			http.ListenAndServe(":9090", nil)
		}
	问题来了，用GO用久了，面向对象，总是这样做 ...
		type Demo struct {
			Data1 ...
			Data2 ...
			Data3 ...
		}
		func (this *Demo) sayhelloName(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, this.Data1 ...)
		}
	但是上面这样用面向对象的形式可行吗 ... 默认结构下是不可以的 ... 你可能会说我可以自定义
		type Demo struct {
			Data1 ...
			Data2 ...
		}
		func (this *Demo) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
	但是这样真的方便吗 ... 自定义路由？ ... 等等 ... 可能你还会有办法比如这样：
		var DEMO *Demo
		type Demo struct {
			Data1 ...
			Data2 ...
		}
		func sayhelloName(w http.ResponseWriter, r *http.Request) {	
			fmt.Fprintf(w, DEMO.Data1 ....)
		}
		func main() {
			DEMO = new(Demo)
			http.HandleFunc("/", sayhelloName) 
			http.ListenAndServe(":9090", nil)
		}
	这样是可行，全局变量 ... 但是又回来了 ... 面向对象的结构呢 ... 
	以上的方法我都曾经在不同的项目里面使用过 ... 但是总觉得差点什么 ...
	现在做了这个项目 ...  现在你只要 ...
		type Demo struct {
			Data1 ...
			Data2 ...
		}
		// 其中 freehttp.JsonIndent 为内置帮助方法
		// 自动把返回的对象(map、struct)转换为 Json 输出 ...
		func (this *Demo) Hello() freehttp.JsonIndent {
			m := make(map[string]interface{})
			m["nljb"] = "www.nljb.net"
			return m
		}
		demo := new(Demo)
		service := freehttp.NewService(demo)
		service.Router("/hello", demo.Hello)
		service.Start(":8080")
	直接访问 http://127.0.0.1:8080/hello 吧 ...
	这样就可以了 ... 不仅这些 ... freehttp 还提供了很多内置帮助方法 ...
----------------

* 一、可以直接将自定义结构类注册到freehttp服务 ...
* 二、可以为自定义结构类中的子类指定路由地址（访问地址URL）
* 三、结构类中子类均可成为HTTP-API并且可以使用结构类中的资源
* 四、提供方法输入与输出参数通过反射封装标记类型提供人性化帮助方法 ...

----------------

安装

	go get github.com/nulijiabei/freehttp
	
----------------

输入类型

 输入 | 标记类型  | 类型原型 | 类型说明
  ------------- | ------------- | ------------- | -------------
 输入 | *freehttp.FreeHttp | -  |  - 
 输入 | *freehttp.Request | -  |  - 
 输入 | *freehttp.ResponseWriter | -  |  - 
 输入 | *freehttp.INI | -  |  - 
 输入 | freehttp.Json | -  |  - 
 输入 | freehttp.ContentType | -  |  - 
 输入 | freehttp.Json | -  |  - 
 输入 | freehttp.Stream | -  |  - 

----------------

输出类型

 输出 | 标记类型  | 类型原型 | 类型说明
  ------------- | ------------- | ------------- | ------------- 
 输出 | freehttp.Json | map[string]interface{}  |  -
 输出 | freehttp.JsonIndent | map[string]interface{}  |  -
 输出 | freehttp.HttpStatus | int  | -
 输出 | freehttp.ContentType | string  | -
 输出 | freehttp.Stream | *bufio.Reader  | -
 输出 | freehttp.File | string  | -
 输出 | freehttp.Redirect | string  | -

----------------

帮助方法

帮助 | 方法继承 | 帮助方法 | 帮助说明
------------- | ------------- | ------------- | -------------
帮助 | freehttp.FreeHttp | NewWebSokcet | HTTP -> WebSocket

----------------

核心结构

	核心默认服务：
	
		// FreeHttp
		freehttp.FreeHttp
		
		SuperResponseWriter 		*freehttp.ResponseWriter
		SuperRequest        		*freehttp.Request

		freehttp.Request		封装于	http.Request
		freehttp.Reader			封装于	http.Request.Body -> bufio.NewReader(Body)

		freehttp.ResponseWriter		封装于	http.ResponseWriter
		freehttp.Writer			封装于	http.ResponseWriter.Body -> bufio.NewWriter(Body)

	衍生帮助方法:
	
		freehttp.Request.*			基于对 http.Request        的自定义帮助方法
		freehttp.ResponseWriter.*		基于对 http.ResponseWriter 的自定义帮助方法
		
----------------
		
	核心配置类型

		// 类型
		freehttp.INI 

		// 初始化配置文件（INI格式）
		service := freehttp.NewService()
		service.Config("/profile")
		...

		// 例如
		func (this *MyStruct) MyFunc(conf *freehttp.INI) {
			conf.Show()
			conf.Set("default", "freehttp", "initalize")
			conf.GetString("default.freehttp", "default value")
			conf.Del("default", "freehttp")
			conf.Save()
		}

		// 说明：未初始化配置文件的情况下使用*freehttp.INI参数会造成错误
		// Use a non-initialized type: *freehttp.INI
		
----------------
		
	核心输入类型:

		// 类型
		freehttp.FreeHttp

		// 例如
		func (this *MyStruct) MyFunc(rw *freehttp.FreeHttp) {
			rw.SuperResponseWriter.ResponseWriter.Write([]byte("print"))
		}

	继承输入类型:

		// 类型
		freehttp.Request = *http.Request

		// 例如
		func (this *MyStruct) MyFunc(r *freehttp.Request) {
			// r.Request.Body
		}

	继承输出类型:

		// 类型
		freehttp.ResponseWriter = http.ResponseWriter 

		// 例如
		func (this *MyStruct) MyFunc(w *freehttp.ResponseWriter) {
			// w.ResponseWriter.Write()
		}
		
----------------
	
	衍生输入类型:

		// Bufio.Reader
		freehttp.Stream		原型 -> *bufio.Reader		还原 -> StreamType(v)

		// 例如
		func (this *MyStruct) MyFunc(stream freehttp.Stream) {
			f, err := os.OpenFile("E:\\a.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			...
			if _, err := io.Copy(bufio.NewWriter(f), freehttp.StreamType(stream))
			...
		}
		
----------------
		
	衍生输出类型:

		// Json 普通格式
		freehttp.Json			原型 -> map[string]interface{}	 

		// Json 排版格式
		freehttp.JsonIndent		原型 -> map[string]interface{}

		// 例如
		func (this *MyStruct) MyFunc() (freehttp.Json, freehttp.JsonIndent) {
			m := make(map[string]interface{})
			m["baidu"] = "www.baidu.com"
			return m, m
		}

		// HTTP Status
		freehttp.HttpStatus		原型 -> int		还原 -> HttpStatusType(v)

		// 例如
		func (this *MyStruct) MyFunc() freehttp.HttpStatus {
			return 404
		}

		// Content-Type
		freehttp.ContentType		原型 -> string		还原 -> ContentTypeType(v)

		// 例如
		func (this *MyStruct) MyFunc() freehttp.ContentType {
			return "image/gif"
		}

		// Bufio.Reader
		freehttp.Stream			原型 -> *bufio.Reader		还原 -> StreamType(v)

		// 例如
		func (this *MyStruct) MyFunc() freehttp.Stream {
			return bufio.NewReader(strings.NewReader("..."))
		}

		// File
		freehttp.File			原型 -> string		还原 -> FileType(v)

		// 例如
		func (this *MyStruct) MyFunc() freehttp.File {
			return ".../freehttp/README.md"
		}

		// Redirect	
		freehttp.Redirect		原型 -> strings		还原 -> RedirectType(v)

		// 例如
		func (this *MyStruct) MyFunc() freehttp.Redirect {
			return "http://www.baidu.com"
		}

		
----------------

	衍生帮助方法:
	
		// HTTP -> WebSocket
		func (this *MyStruct) MyFunc(rw *freehttp.FreeHttp) {
			// HTTP -> WebSocket
			rw.NewWebSokcet(func(conn *freehttp.WSConn) {
				// WSConn = websocket.Conn
				conn.Write([]byte("Hello WebSokcet !!!"))
				r := bufio.NewReader(conn)
				for {
					v, err := r.ReadBytes('\n')
					if err != nil {
						if err != io.EOF {
							panic(err)
						}
						break
					}
					conn.Write(v)
					conn.Write([]byte("\n"))
				}
			})
		}
	
----------------

	// 本例会随着程序版本的更新而更新

	// 例
	package main

	import (
		"fmt"
		"freehttp" // 导入 freehttp 包
	)

	// 随便定义一个类
	type Web struct {}


	// 随便定义一些方法
	func (this *Web) MyFunc(输入类型 + 输入类型 + ...) 输出类型 + 输出类型 + ... {
		// 使用输入类型 ...
		// 使用帮助方法 ...
		// 返回输出类型 ...
	}

	// 启动
	func main() {
		
		// New 自定义结构类
		web := new(Web)

		// 创建一个 service
		service := freehttp.NewService(web)
		// service.Config("/profile")
		
		// 路由 ...
		service.Router("/baidu", web.MyFunc)

		// 启动服务器
		if err := service.Start(":8080"); err != nil {
			panic(err)
		}

	}

-----------------

	// 案例：

	package main
	
	import (
		"log"
	
		"../../../freehttp"
	)
	
	type API struct {
	}
	
	// ...
	type Update struct {
		Version string `json:"version"`
		Build   string `json:"build"`
		Url     string `json:"url"`
	}
	
	func (this *API) Update(r *freehttp.Request) freehttp.JsonIndent {
		// 解析参数
		r.Request.ParseForm()
		// 获取参数
		version := r.Request.FormValue("version")
		build := r.Request.FormValue("build")
		// 输出
		log.Printf("update request form version(%s) build(%s)", version, build)
		// update
		up := new(Update)
		up.Version = "1.0.10"
		up.Build = "10"
		up.Url = "www.baidu.com"
		// 返回
		return up
	}
	
	func main() {
		api := new(API)
		service := freehttp.NewService(api)
		// service.Config("/profile")
		service.Router("/update", api.Update)
		if err := service.Start(":8080"); err != nil {
			panic(err)
		}
	}


