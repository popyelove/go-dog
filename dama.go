package main

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"time"
	"go-dog/st"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"flag"
)

func print_code(configuration st.Configuration){
	url := "https://pet-chain.baidu.com/data/captcha/gen"
	var jsonStr = []byte(`{
							"requestId":1523433103485,
							"appId":1,"tpl":"",
							"timeStamp":null,
							"nounce":null,
							"token":null}
						`)
	req,_:= http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie",configuration.COOKIE)
	client := &http.Client{}
	resp,_:= client.Do(req)
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		js,err:=simplejson.NewJson([]byte(string(body)))
		if err!=nil {
			fmt.Print(err)
		}
		seed,_:=js.Get("data").Get("seed").String()
		imgbase64,_:=js.Get("data").Get("img").String()
		key:=configuration.KEY
		code:=lujun_api(key,imgbase64)
		fmt.Print("验证码="+code,"====>seed="+seed,"\n")
		if code!="" {
			c,_:= redis.Dial("tcp", "127.0.0.1:6379")
			jsonstr:=`{"code":"`+code+`","seed":"`+seed+`"}`
			c.Do("rpush", "code_list",jsonstr)
		}

	}

}
//验证码识别接口
func lujun_api(key string,img64 string) string {
	url := "http://api.lujun.co:8888/?key="+key+"&img="+img64
	resp,_ := http.Get(url)
	if resp !=nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		js,err:=simplejson.NewJson([]byte(string(body)))
		if err!=nil {
			fmt.Print(err)
		}
		captcha,_:=js.Get("captcha").String()
		return  captcha
	}
	return ""
}

func main(){
	config := flag.String("f", "", "配置文件")
	flag.Parse()
	configfile:=*config
	var  configuration st.Configuration
	configuration.GetConf(configfile)
	ticker := time.NewTicker(configuration.TIMECODE* time.Millisecond)
	for _ = range ticker.C {
		print_code(configuration)
	}

}