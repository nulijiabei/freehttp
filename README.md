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
	
		// FreeHttp
		freehttp.FreeHttp
		
			SuperResponseWriter 	*freehttp.ResponseWriter
			SuperRequest        	*freehttp.Request
		
			freehttp.Request		封装于	http.Request
			freehttp.Reader			封装于	http.Request.Body -> bufio.NewReader(Body)
				
			freehttp.ResponseWriter	封装于	http.ResponseWriter
			freehttp.Writer			封装于	http.ResponseWriter.Body -> bufio.NewWriter(Body)

	衍生帮助方法:
	
		freehttp.Request.*			基于对 http.Request        的自定义帮助方法
		freehttp.ResponseWriter.*	基于对 http.ResponseWriter 的自定义帮助方法
		
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
		freehttp.Stream			原型 ->	*bufio.Reader
		
		// 例如
		func (this *Struct) MyFunc(stream freehttp.Stream) {
			f, _ := os.Open("../freehttp/README.md")
			io.Copy(rw.SuperResponseWriter.Writer, bufio.NewReader(f))
		}
		
	衍生输出类型:
		
		// Json 普通格式
		freehttp.Json			原型 ->	map[string]interface{}	 
		
		// Json 排版格式
		freehttp.JsonIndent		原型 ->	map[string]interface{}

		// 例如
		func (this *Struct) MyFunc() (freehttp.Json, freehttp.JsonIndent) {
			m := make(map[string]interface{})
			m["baidu"] = "www.baidu.com"
			return m, m
		}
			
		// HTTP Status
		freehttp.HttpStatus		原型 ->	int
		
		// 例如
		func (this *Struct) MyFunc() freehttp.HttpStatus {
			return 404
		}
		
		// Content-Type
		freehttp.ContentType	原型 ->	string
		
		// 例如
		func (this *Struct) MyFunc() freehttp.ContentType {
			return "image/gif"
		}
		
		// Bufio.Reader
		freehttp.Stream			原型 ->	*bufio.Reader
		
		// 例如
		func (this *Struct) MyFunc() freehttp.Stream {
			return bufio.NewReader(strings.NewReader("..."))
		}
		
		// File
		freehttp.File			原型 -> string
		
		// 例如
		func (this *Struct) MyFunc() freehttp.File {
			return ".../freehttp/README.md"
		}

		
----------------

		自定义:
		
		freehttp.NewServer(
			// mname = StructName
			// name  = FuncName
			func(mname, name string) string {
				// return string == r.URL.Path
				return strings.ToLower(fmt.Sprintf("/%s/%s", mname, name))
			}
		)

----------------

		// 例
		package main
	
		import (
			"fmt"
			"freehttp" // 导入 freehttp 包
		)
	
		// 随便定义一个类
		type Web struct {}
	
	
		// 随便定义一些方法
		// http://127.0.0.1:8080/StructName.FuncName
		func (this *Web) MyFunc(输入类型 + 输入类型 + ...) 输出类型 + 输出类型 + ... {
			// 使用输入类型 ...
			// 返回输出类型 ...
		}
		
		// 启动
		func main() {
	
			// 创建 Server
			s := server.NewServer(nil)
	
			// 传入 Web 类
			if err := s.Register(new(Web)); err != nil {
				panic(err)
			}
		
			// 启动监听端口
			s.Start(":8080")
	
		}

-----------------

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
			s := freehttp.NewServer(
				func(mname, name string) string {
					return strings.ToLower(fmt.Sprintf("/%s/%s", mname, name))
				})
			if err := s.Register(new(API)); err != nil {
				fmt.Println(err)
			}
			s.Start(":9090")
		}

