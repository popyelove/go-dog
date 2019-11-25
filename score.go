package main

import (
	"go-dog/tool"
	"go-dog/st"
	"fmt"
	"github.com/bitly/go-simplejson"
	"flag"
	"time"
	"math"
	"strconv"
)
func dog_lists(cookie string,account int)  {
	body := tool.GetList(cookie,"1")
	js,_:= simplejson.NewJson([]byte(body))
	if js!=nil{
		count:=js.Get("data").Get("totalCount").MustFloat64()
		dealMap := make(map[interface{}]interface{})
		if count > 10 {
			size := math.Ceil(count/10)
			for p :=1; p <= int(size); p++ {
				if p > 1 {
					body := tool.GetList(cookie,strconv.Itoa(p))
					js,_:= simplejson.NewJson([]byte(body))
					for i :=0;i < PAGE_SIZE;i++{
						info:=js.Get("data").Get("dataList").GetIndex(i).MustMap()
						if info!=nil{
							dealMap[info["petId"]] = info["rareDegree"]
						}
					}
				} else {
					for i :=0;i < PAGE_SIZE;i++{
						info:=js.Get("data").Get("dataList").GetIndex(i).MustMap()
						if info!=nil{
							dealMap[info["petId"]] = info["rareDegree"]
						}
					}
				}
			}
		} else {
			for i :=0;i < PAGE_SIZE;i++{
				info:=js.Get("data").Get("dataList").GetIndex(i).MustMap()
				if info!=nil{
					dealMap[info["petId"]] = info["rareDegree"]
				}
			}
		}

		for petid,raredgee:= range dealMap {
			fmt.Print("petid====>",petid,"      rareDegree==>",raredgee)
			fmt.Print("\n\n")
		}

	} else {
		fmt.Print("List==>null",time.Now())
		fmt.Print("\n\n")
	}
}
var PAGE_SIZE=10
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
	cookies:=configuration.COOKIE
	for i:=0;i<len(cookies);i++{
		dog_lists(cookies[i],i)
	}
}