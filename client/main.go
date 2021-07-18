package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main(){
	for i:=1;i<1000;i++{
		go func() {
			readName()
		}()
		time.Sleep(1*time.Second)
	}
	time.Sleep(100*time.Second)
}

func readName(){
	//使用Get方法获取服务器响应包数据
	resp, err := http.Get("http://localhost:8888/name")

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	//读取body内的内容
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(content))
}