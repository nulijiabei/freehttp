package freehttp

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// INI
type INI struct {
	path string
	conf map[string]string
	lock *sync.RWMutex
}

// New INI
func NewINI(path string) *INI {
	ini := new(INI)
	ini.path = path
	ini.conf = make(map[string]string)
	ini.lock = new(sync.RWMutex)
	err := ini.load()
	if err != nil {
		panic(err)
	}
	return ini
}

// 初始化时载入配置文件
func (this *INI) load() error {
	f, err := os.Open(this.path)
	if err != nil {
		return err
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	var group string
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		line = Trim(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group = line[1:(len(line) - 1)]
		}
		if strings.Contains(line, "=") {
			ss := strings.Split(line, "=")
			if IsBlank(group) {
				this.Set("default", Trim(ss[0]), Trim(ss[1]))
			} else {
				this.Set(group, Trim(ss[0]), Trim(ss[1]))
			}
		}
	}
	return nil
}

// 设置配置项目 (组,对象,值)
func (this *INI) Set(group, key, value string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.conf[fmt.Sprintf("%s.%s", group, key)] = value
}

// 删除配置项目
func (this *INI) Del(group, key string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.conf, fmt.Sprintf("%s.%s", group, key))
}

// 查看配置内容
func (this *INI) Show() {
	this.lock.RLock()
	defer this.lock.RUnlock()
	data := make(map[string]map[string]string)
	for k, v := range this.conf {
		arr := strings.Split(k, ".")
		if val, ok := data[arr[0]]; ok {
			val[arr[1]] = v
		} else {
			data[arr[0]] = make(map[string]string)
			data[arr[0]][arr[1]] = v
		}
	}
	var content string
	for k, v := range data {
		content += fmt.Sprintf("[%s]\n", k)
		for kk, vv := range v {
			content += fmt.Sprintf("%s=%s\n", kk, vv)
		}
	}
	fmt.Printf("\n", this.path, content)
}

// 将当前配置保存到配置文件内
func (this *INI) Save() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	data := make(map[string]map[string]string)
	for k, v := range this.conf {
		arr := strings.Split(k, ".")
		if val, ok := data[arr[0]]; ok {
			val[arr[1]] = v
		} else {
			data[arr[0]] = make(map[string]string)
			data[arr[0]][arr[1]] = v
		}
	}
	var content string
	for k, v := range data {
		content += fmt.Sprintf("[%s]\n", k)
		for kk, vv := range v {
			content += fmt.Sprintf("%s=%s\n", kk, vv)
		}
	}
	return ioutil.WriteFile(this.path, []byte(content), 0644)
}

// 返回字符串值
// key(root.system)
// def = default 失败则返回
func (this *INI) GetString(key, def string) string {
	this.lock.RLock()
	defer this.lock.RUnlock()
	val, ok := this.conf[key]
	if ok {
		return val
	}
	return def
}
