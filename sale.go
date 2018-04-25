package main

import (
	"go-dog/tool"
	"go-dog/st"
	"fmt"

	"github.com/bitly/go-simplejson"
	"flag"
	"time"
)

func auto_do_big(configuration st.Configuration)  {
	for petid :=range petids{
		res := tool.Sale(petid, petids[petid],configuration)
		res = tool.Unsale(petid,configuration)
		res = tool.Sale(petid, petids[petid],configuration)
		js,_:= simplejson.NewJson([]byte(res))
		if js!=nil{
			errno:=js.Get("errorNo").MustString()
			if errno =="30007"||errno=="30003"{
				delete(petids,petid);
			}
			fmt.Println(res)
		}else{
			fmt.Println("请到配置文件中配置你的需要卖出的狗狗")
		}
	}
}
var petids = map[string]string{}
func main() {
	var config string
	f := flag.String("f", "", "配置文件路径")
	flag.Parse() //解析输入的参数
	if(*f==""){
		fmt.Printf("请输入你的配置文件的绝对路径(例如：D:/file/conf.yaml)：")
		fmt.Scanln(&config)
	}else{
		config=*f
	}
	var  configuration st.Configuration
	configuration.GetConf(config)
	petids=configuration.SALE_LIST
	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		go auto_do_big(configuration)
	}
}