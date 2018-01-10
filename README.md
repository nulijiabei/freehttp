----------------
freehttp

<a href="https://godoc.org/github.com/nulijiabei/freehttp"><img src="https://godoc.org/github.com/nulijiabei/freehttp?status.svg" alt="GoDoc"></a>

一个快速将类和子方法转换成HTTP接口

http://127.0.0.1:8080/MyStructName/MyFuncName

主要是对net/http的一个反射封装，通过几行简短代码即可实现复杂功能，再也不需要每个功能点一一实现 ...

(标记类型）通过反射获取自定义类型来判断用户目的，一行代码完成反锁操作。

比如：return "/a/b/c.jpg" -> string 	你永远不知道用户目的。

这样：return "/a/b/c.jpg" -> File 	这样用户目的一目了然。

----------------

安装

	go get github.com/nulijiabei/freehttp
	

----------------

  输入/输出 | 标记类型  | 类型原型 | 类型还原 | 类型说明
  ------------- | ------------- | ------------- | ------------- | ------------- 
 输出 | Json | interface{}  | Func JsonType(v) | 
 输出 | JsonIndent | interface{}  | Func JsonIndentType(v) | a
 输出 | HttpStatus | int  | Func HttpStatusType(v) | a
 输出 | ContentType | string  | Func ContentTypeType(v)  | a
 输入/输出 | Stream | *bufio.Reader  | Func StreamType(v) | a
 输出 | File | string  | Func FileType(v)  | a
 输出 | Redirect | string  | Func RedirectTyoe(v)  | a

----------------

使用方法

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
		
	核心配置类型
	
		// 类型
		freehttp.INI 
		
		// 初始化配置文件（INI格式）
		service := freehttp.NewService()
		service.Config("/profile")
		...
		
		// 例如
		func (this *Struct) MyFunc(conf *freehttp.INI) {
			conf.Show()
			conf.Set("default", "freehttp", "initalize")
			conf.GetString("default.freehttp", "default value")
			conf.Del("default", "freehttp")
			conf.Save()
		}
		
		// 说明：未初始化配置文件的情况下使用*freehttp.INI参数会造成错误
		// Use a non-initialized type: *freehttp.INI
		
	核心输入类型:
	
		// 类型
		freehttp.FreeHttp
	
		// 例如
		func (this *Struct) MyFunc(rw *freehttp.FreeHttp) {
			rw.SuperResponseWriter.ResponseWriter.Write([]byte("print"))
		}
		
	继承输入类型:
	
		// 类型
		freehttp.Request = *http.Request
		
		// 例如
		func (this *Struct) MyFunc(r *freehttp.Request) {
			// r.Request.Body
		}

	继承输出类型:
	
		// 类型
		freehttp.ResponseWriter = http.ResponseWriter 
		
		// 例如
		func (this *Struct) MyFunc(w *freehttp.ResponseWriter) {
			// w.ResponseWriter.Write()
		}
	
	衍生输入类型:
		
		// Bufio.Reader
		freehttp.Stream		原型 -> *bufio.Reader		还原 -> StreamType(v)
		
		// 例如
		func (this *Struct) MyFunc(stream freehttp.Stream) {
			f, err := os.OpenFile("E:\\a.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			...
			if _, err := io.Copy(bufio.NewWriter(f), freehttp.StreamType(stream))
			...
		}
		
	衍生输出类型:
		
		// Json 普通格式
		freehttp.Json			原型 -> map[string]interface{}	 
		
		// Json 排版格式
		freehttp.JsonIndent		原型 -> map[string]interface{}

		// 例如
		func (this *Struct) MyFunc() (freehttp.Json, freehttp.JsonIndent) {
			m := make(map[string]interface{})
			m["baidu"] = "www.baidu.com"
			return m, m
		}
			
		// HTTP Status
		freehttp.HttpStatus		原型 -> int		还原 -> HttpStatusType(v)
		
		// 例如
		func (this *Struct) MyFunc() freehttp.HttpStatus {
			return 404
		}
		
		// Content-Type
		freehttp.ContentType		原型 -> string		还原 -> ContentTypeType(v)
		
		// 例如
		func (this *Struct) MyFunc() freehttp.ContentType {
			return "image/gif"
		}
		
		// Bufio.Reader
		freehttp.Stream			原型 -> *bufio.Reader		还原 -> StreamType(v)
		
		// 例如
		func (this *Struct) MyFunc() freehttp.Stream {
			return bufio.NewReader(strings.NewReader("..."))
		}
		
		// File
		freehttp.File			原型 -> string		还原 -> FileType(v)
		
		// 例如
		func (this *Struct) MyFunc() freehttp.File {
			return ".../freehttp/README.md"
		}
		
		// Redirect	
		freehttp.Redirect		原型 -> strings		还原 -> RedirectType(v)
		
		// 例如
		func (this *Struct) MyFunc() freehttp.Redirect {
			return "http://www.baidu.com"
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
		// http://127.0.0.1:8080/StructName/FuncName
		func (this *Web) MyFunc(输入类型 + 输入类型 + ...) 输出类型 + 输出类型 + ... {
			// 使用输入类型 ...
			// 返回输出类型 ...
		}
		
		// 启动
		func main() {
	
			// 创建一个 service
			service := freehttp.NewService(new(Web))
			// service.Config("/profile")
		
			// 启动服务器
			if err := service.Start(":8080"); err != nil {
				panic(err)
			}
	
		}

-----------------

		// 案例不会随着代码的更新而更新，具体使用方法请根据使用方法

		案例：
		
		package main
		
		import (
			"fmt"
			"freehttp"
			"log"
			"strings"
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
			return freehttp.JsonIndent(up)
		}
		
		func main() {
			service := freehttp.NewService(new(API))
			// service.Config("/profile")
			if err := service.Start(":8080"); err != nil {
				panic(err)
			}
		}

