package main
import (
	"net/http"
	"io/ioutil"
	"bytes"
	"github.com/bitly/go-simplejson"
	"encoding/json"
	"github.com/json-iterator/go"
	"fmt"
	"time"
	"go-dog/st"
	"strconv"
	"container/list"
	"os"
)

func http_post(url string,jsonStr []byte,configuration st.Configuration,ch chan string	)  {
	req,_:= http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie",configuration.COOKIE)
	client := &http.Client{}
	resp,err:= client.Do(req)
	if err!=nil{
		return
	}
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		ch <- string(body)
	}else {
		ch <- ""
	}
	return
}
//获取狗的列表
func dog_list(configuration st.Configuration) string {
	url := "https://pet-chain.baidu.com/data/market/queryPetsOnSale"
	var jsonStr = []byte(`{
		"pageNo":1,
		"pageSize":`+strconv.Itoa(configuration.PAGE_SIZE)+`,
		"querySortType":"`+configuration.SORT_TYPE+`",
		"petIds":[],
		"lastAmount":"",
		"lastRareDegree":"",
		"filterCondition":"{`+dog_filter[index_dog]+`}",
		"appId":1,
		"tpl":"",
		"type":null,
		"requestId":1522231859931,
		"timeStamp":null,
		"nounce":null,
		"token":null
		}`)
	ch_run := make(chan string)
	go http_post(url,jsonStr,configuration,ch_run)
	select {
	case res := <-ch_run:
		if(res!=""){
			fmt.Print("抢狗进行中...",time.Now())
			fmt.Print("\n")
		}
		return res
	case <-time.After(dog_list_timeout * time.Second):
		fmt.Println("拉取狗狗列表接口超时,请检查刷狗频率参数是否过小。。。。\n")
		return ""
	}
	return ""
}
//下单买狗借口
func bug_dog(petId string,amount string,seed string,code string ,validCode string,configuration st.Configuration) string{
	type tiaojian struct {
		Petid string `json:"petId"`
		Amount string `json:"amount"`
		Seed string `json:"seed"`
		Captcha string `json:"captcha"`
		ValidCode string `json:"validCode"`
		RequestId string `json:"requestId"`
		Appid string `json:"appId"`
		Tpl string `json:"tpl"`
	}
	json_tiaojian :=tiaojian{Petid:petId,Amount:amount,Seed:seed,Captcha:code,ValidCode:validCode,RequestId:"1520241678619",Appid:"1",Tpl:""}
	url := "https://pet-chain.baidu.com/data/txn/create"
	jsonStr,_:=json.Marshal(json_tiaojian)
	ch_run := make(chan string)
	go http_post(url,jsonStr,configuration,ch_run)
	select {
	case res := <-ch_run:
		return res
	case <-time.After(buy_dog_timeout * time.Second):
		fmt.Println("交易火爆中，请稍后再试。。。！")
		return ""
	}
	return ""
}
//获取狗的稀有属性
func get_dog_rareDegree(petid string,configuration st.Configuration)(int,int){
	url := "https://pet-chain.baidu.com/data/pet/queryPetById"
	var jsonStr = []byte(`{
        "petId":`+petid+`,
        "requestId":1520241678619,
        "appId":1,
        "tpl":"",
        "timeStamp":"",
        "nounce":"",
        "token":""
    }`)

	ch_run := make(chan string)
	go http_post(url,jsonStr,configuration,ch_run)
	select {
	case res := <-ch_run:
		js,err:=simplejson.NewJson([]byte(res))
		if err!=nil {
			fmt.Print(err)
		}
		if js==nil {
			return 0,0
		}
		count_rareDegree :=0
		dogtype :=0
		for i:=0;i<8 ;i++  {
			s:= js.Get("data").Get("attributes").GetIndex(i).MustMap()
			if s["rareDegree"]=="稀有" {
				count_rareDegree=count_rareDegree+1
			}
			if (s["value"]=="天使"){
				dogtype+=1
			}
			if (s["value"]=="白眉斗眼"){
				dogtype+=1
			}
		}
		return count_rareDegree,dogtype
	case <-time.After(get_dog_rare_timeout * time.Second):
		fmt.Println("获取狗的稀有属性超时。。。。。\n")
		return 0,0
	}
	return 0,0
}
//神话狗
func shenhua_dog(dog map[string]interface{},configuration st.Configuration)bool  {
	rareDegree,_:=jsoniter.MarshalToString(dog["rareDegree"])//稀有度
	amount:=jsoniter.Wrap(dog["amount"]).ToFloat32()//价额
	timeLeft :=jsoniter.Wrap(dog["coolingInterval"]).ToString()//休息时间
	generation,_:=jsoniter.MarshalToString(dog["generation"])//代数
	rareDegrees,_:=get_dog_rareDegree(dog["petId"].(string),configuration) //属性稀有个数
	if(rareDegrees==6&&rareDegree=="4"){
		//六稀神话
		if(generation=="0"&&configuration.GOD0_6_SWITCH==1){
			//0代神话价格
			if (amount<=configuration.GOD0_6DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
			if (amount<=configuration.GOD0_6DOG_24_PRICE&&timeLeft=="24小时"){
				return true
			}
			if (amount<=configuration.GOD0_6DOG_2_PRICE&&timeLeft=="2天"){
				return true
			}
		}
		if(generation=="1"&&configuration.GOD1_6_SWITCH==1){
			//一代神话价格配置
			if (amount<=configuration.GOD1_6DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
			if (amount<=configuration.GOD1_6DOG_2_PRICE&&timeLeft=="2天"){
				return true
			}
			if (amount<=configuration.GOD1_6DOG_4_PRICE&&timeLeft=="4天"){
				return true
			}
		}
		if(generation=="2"&&configuration.GOD2_6_SWITCH==1){
			//二代神话价格配置
			if (amount<=configuration.GOD2_6DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
			if (amount<=configuration.GOD2_6DOG_4_PRICE&&timeLeft=="4天"){
				return true
			}
			if (amount<=configuration.GOD2_6DOG_6_PRICE&&timeLeft=="6天"){
				return true
			}
		}
		if(generation=="3"&&configuration.GOD3_6_SWITCH==1){
			//三代神话价格配置
			if (amount<=configuration.GOD3_6DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
			if (amount<=configuration.GOD3_6DOG_6_PRICE&&timeLeft=="6天"){
				return true
			}
			if (amount<=configuration.GOD3_6DOG_8_PRICE&&timeLeft=="8天"){
				return true
			}
		}
	}
	if(rareDegrees==7&&rareDegree=="4"){
		//七夕神话
		if(generation=="0"&&configuration.GOD0_7_SWITCH==1){
			//0代神话价格
			if (amount<=configuration.GOD0_7DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
			if (amount<=configuration.GOD0_7DOG_24_PRICE&&timeLeft=="24小时"){
				return true
			}
			if (amount<=configuration.GOD0_7DOG_2_PRICE&&timeLeft=="2天"){
				return true
			}
		}
		if(generation=="1"&&configuration.GOD1_7_SWITCH==1){
			//1代神话价格
			if (amount<=configuration.GOD1_7DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
			if (amount<=configuration.GOD1_7DOG_2_PRICE&&timeLeft=="2天"){
				return true
			}
			if (amount<=configuration.GOD1_7DOG_4_PRICE&&timeLeft=="4天"){
				return true
			}
		}

	}

	return false
}
//获取验证吗
func get_code()string{
	code :=code_list.Back()
	if code!=nil{
		return code.Value.(string)
	}
	return "{}"
}

//循环刷狗函数
func do_always(configuration st.Configuration)  {
	dogs :=dog_list(configuration)
	if dogs !=""{
		js,_:= simplejson.NewJson([]byte(dogs))
		if js !=nil{
			for i :=0;i<configuration.PAGE_SIZE ;i++  {
				s:= js.Get("data").Get("petsOnSale").GetIndex(i).MustMap()
				if s !=nil{
					if shenhua_dog(s,configuration){
						codes :=get_code()
						json,_ :=simplejson.NewJson([]byte(codes))
						if json !=nil{
							seed :=json.Get("seed").MustString()
							code :=json.Get("code").MustString()
							bres :=bug_dog(s["petId"].(string),s["amount"].(string),seed,code,s["validCode"].(string),configuration)
							fmt.Println(bres)
							res,_ :=simplejson.NewJson([]byte(bres))
							if res!=nil {
								errorNo :=res.Get("errorNo").MustString()
								errorMsg :=res.Get("errorMsg").MustString()
								if errorNo=="08"{
									//交易火爆，区块链处理繁忙，请稍后再试
									fmt.Print(errorMsg)
								}
								if errorNo=="10002" {
									//有人抢先下单啦
									fmt.Print(errorMsg)
								}
								if errorNo =="00"{
									//success
									fmt.Print("抢到狗狗啦！！！！！！","\n",s)
								}
							}
						}

					}

				}

			}

		}

	}

}

func print_code(configuration st.Configuration){
	url := "https://pet-chain.baidu.com/data/captcha/gen"
	var jsonStr = []byte(`{
							"requestId":1523433103485,
							"appId":1,"tpl":"",
							"timeStamp":null,
							"nounce":null,
							"token":null}
						`)

	ch_run := make(chan string)
	go http_post(url,jsonStr,configuration,ch_run)
	select {
	case res := <-ch_run:
		js,err:=simplejson.NewJson([]byte(res))
		if err!=nil {
			fmt.Print("百度服务器繁忙！！！","\n")
			return
		}
		var seed string
		seed,err=js.Get("data").Get("seed").String()
		if err!=nil {
			fmt.Print("百度服务器繁忙。。。。。。。。。。","\n")
			return
		}
		imgbase64,err:=js.Get("data").Get("img").String()
		if err!=nil {
			fmt.Print(err)
			return
		}
		key:=configuration.KEY
		code_res:=lujun_api(key,imgbase64)
		js_code,err:=simplejson.NewJson([]byte(code_res))
		if err!=nil {
			fmt.Print(err)
			return
		}
		status :=js_code.Get("status").MustString()
		msg :=js_code.Get("msg").MustString()
		if status=="error" {
			fmt.Print(msg,"\n")
			os.Exit(0)
		}
		code :=js_code.Get("captcha").MustString()
		fmt.Print("验证码="+code,"====>seed="+seed,"\n")
		if code!="" {
			jsonstr:=`{"code":"`+code+`","seed":"`+seed+`"}`
			len :=code_list.Len()
			if(len>=500){
				code_list.Init()
			}
			code_list.PushBack(jsonstr)
		}
	case <-time.After(dama_timeout * time.Second):
		fmt.Println("百度验证码接口超时......")
		return
	}

}
//验证码识别接口
func lujun_api(key string,img64 string) string {
	url := "http://www.popyelove.com:8888/?key="+key+"&img="+img64
	resp,_ := http.Get(url)
	if resp !=nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
		js,err:=simplejson.NewJson([]byte(string(body)))
		if err!=nil {
			fmt.Print(err)
		}
		captcha,_:=js.Get("captcha").String()
		return  captcha
	}
	return ""
}
//自动打码服务
func Timer2(configuration st.Configuration)  {
	ticker := time.NewTicker(dama_time* time.Millisecond)
	for _ = range ticker.C {
		print_code(configuration)
	}
}
var config string
var code_list *list.List
var dog_filter = [1]string{"1:4"}
var index_dog =0
//打码间隔 毫秒
var dama_time time.Duration=10000
//拉取狗列表超时时间秒
var dog_list_timeout time.Duration=15
//下单超时时间
var buy_dog_timeout time.Duration=15
//获取狗狗属性超时
var get_dog_rare_timeout time.Duration=15
//打码超时时间
var dama_timeout time.Duration=15
func main(){
	code_list = list.New()
	fmt.Printf("请输入你的配置文件的绝对路径(例如：D:/file/conf.yaml)：")
	fmt.Scanln(&config)
	var  configuration st.Configuration
	configuration.GetConf(config)
	go Timer2(configuration)
	ticker := time.NewTicker(configuration.TIME* time.Millisecond)
	for _ = range ticker.C {
		go do_always(configuration)
	}

}
