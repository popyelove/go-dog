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
)
//获取狗的列表
//"filterCondition":"{\"1\":\"4\"}",
func dog_list(configuration st.Configuration) string {
	url := "https://pet-chain.baidu.com/data/market/queryPetsOnSale"
	var jsonStr = []byte(`{
		"pageNo":1,
		"pageSize":`+strconv.Itoa(configuration.PAGE_SIZE)+`,
		"querySortType":"CREATETIME_DESC",
		"petIds":[],
		"lastAmount":"",
		"lastRareDegree":"",
		"filterCondition":"{}",
		"appId":1,
		"tpl":"",
		"type":null,
		"requestId":1522231859931,
		"timeStamp":null,
		"nounce":null,
		"token":null
		}`)
	req,_:= http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie",configuration.COOKIE)
	client := &http.Client{}
	resp,_:= client.Do(req)
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Print("抢狗进行中...",time.Now())
		fmt.Print("\n")
		return string(body)
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
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie",configuration.COOKIE)
	client := &http.Client{}
	resp,_:= client.Do(req)
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
	return "{}"
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
	req,_:= http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie",configuration.COOKIE)
	client := &http.Client{}
	resp,_:= client.Do(req)
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		js,_:=simplejson.NewJson([]byte(string(body)))
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
				dogtype+=2
			}
		}
		return count_rareDegree,dogtype
	}
	return 0,0
}
//传说狗
func chuanshuo_dog(dog map[string]interface{},configuration st.Configuration)bool {
	rareDegree, _ := jsoniter.MarshalToString(dog["rareDegree"])               //稀有度
	amount := jsoniter.Wrap(dog["amount"]).ToFloat32()                         //价额
	timeLeft := jsoniter.Wrap(dog["coolingInterval"]).ToString()               //休息时间
	generation, _ := jsoniter.MarshalToString(dog["generation"])               //代数
	if(rareDegree=="5"){
		if(generation=="0"){
			if (amount<=configuration.CHUANSHUO0_8DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
		}
		if(generation=="1"){
			if (amount<=configuration.CHUANSHUO1_8DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
		}
		if(generation=="2"){
			if (amount<=configuration.CHUANSHUO2_8DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
		}
		if(generation=="3"){
			if (amount<=configuration.CHUANSHUO3_8DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
		}
	}
	return false
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
		if (generation=="0")||(generation=="1")||(generation=="2")||(generation=="3"){
			//0代神话价格
			if (amount<=configuration.GOD0_6DOG_0_PRICE&&timeLeft=="0分钟"&&generation=="0"){
				return true
			}
			if (amount<=configuration.GOD0_6DOG_24_PRICE&&timeLeft=="24小时"&&generation=="0"){
				return true
			}
			if (amount<=configuration.GOD0_6DOG_2_PRICE&&timeLeft=="2天"&&generation=="0"){
				return true
			}
			//一代神话价格配置
			if (amount<=configuration.GOD1_6DOG_0_PRICE&&timeLeft=="0分钟"&&generation=="1"){
				return true
			}
			if (amount<=configuration.GOD1_6DOG_2_PRICE&&timeLeft=="2天"&&generation=="1"){
				return true
			}
			if (amount<=configuration.GOD1_6DOG_4_PRICE&&timeLeft=="4天"&&generation=="1"){
				return true
			}
			//二代神话价格配置
			if (amount<=configuration.GOD2_6DOG_0_PRICE&&timeLeft=="0分钟"&&generation=="2"){
				return true
			}
			if (amount<=configuration.GOD2_6DOG_4_PRICE&&timeLeft=="4天"&&generation=="2"){
				return true
			}
			if (amount<=configuration.GOD2_6DOG_6_PRICE&&timeLeft=="6天"&&generation=="2"){
				return true
			}

			//三代神话价格配置
			if (amount<=configuration.GOD3_6DOG_0_PRICE&&timeLeft=="0分钟"&&generation=="3"){
				return true
			}
			if (amount<=configuration.GOD3_6DOG_6_PRICE&&timeLeft=="6天"&&generation=="3"){
				return true
			}
			if (amount<=configuration.GOD3_6DOG_8_PRICE&&timeLeft=="8天"&&generation=="3"){
				return true
			}

		}
	}
	if(rareDegrees==7&&rareDegree=="4"){
		//七夕神话
		if (generation=="0")||(generation=="1")||(generation=="2")||(generation=="3"){
			//0代神话价格
			if (amount<=configuration.GOD0_7DOG_0_PRICE&&timeLeft=="0分钟"&&generation=="0"){
				return true
			}
			if (amount<=configuration.GOD0_7DOG_24_PRICE&&timeLeft=="24小时"&&generation=="0"){
				return true
			}
			if (amount<=configuration.GOD0_7DOG_2_PRICE&&timeLeft=="2天"&&generation=="0"){
				return true
			}
			//1代神话价格
			if (amount<=configuration.GOD1_7DOG_0_PRICE&&timeLeft=="0分钟"&&generation=="1"){
				return true
			}
			if (amount<=configuration.GOD1_7DOG_2_PRICE&&timeLeft=="2天"&&generation=="1"){
				return true
			}
			if (amount<=configuration.GOD1_7DOG_4_PRICE&&timeLeft=="4天"&&generation=="1"){
				return true
			}
		}

	}

	return false
}
//史诗狗
func shishi_dog(dog map[string]interface{},configuration st.Configuration)bool{
	rareDegree,_:=jsoniter.MarshalToString(dog["rareDegree"])
	amount:=jsoniter.Wrap(dog["amount"]).ToFloat32()
	timeLeft :=jsoniter.Wrap(dog["coolingInterval"]).ToString()

	rareDegrees,_:=get_dog_rareDegree(dog["petId"].(string),configuration)

	generation,_:=jsoniter.MarshalToString(dog["generation"])
	//五稀史诗
	if(rareDegrees==5&&rareDegree=="3"){

		if (generation=="0"){
			if (amount<=configuration.SHISHI0_5DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
			if (amount<=configuration.SHISHI0_5DOG_24_PRICE&&timeLeft=="24小时"){
				return true
			}
		}

	}
	//4稀有史诗
	if(rareDegrees==4&&rareDegree=="3"){
		if(generation=="0"){

			if (amount<=configuration.SHISHI0_4DOG_0_PRICE&&timeLeft=="0分钟"){
				return true
			}
			if (amount<=configuration.SHISHI0_4DOG_24_PRICE&&timeLeft=="24小时"){
				return true
			}
		}

	}
	return false
}
//卓越狗
func zhuoyue_dog(dog map[string]interface{},configuration st.Configuration)bool  {

	rareDegree,_:=jsoniter.MarshalToString(dog["rareDegree"])
	amount:=jsoniter.Wrap(dog["amount"]).ToFloat32()
	generation,_:=jsoniter.MarshalToString(dog["generation"])
	if rareDegree=="2"&&generation=="0"&&amount<=configuration.ZHUEYUE0_2DOG_0_PRICE{

		return true
	}
	return false
}
//稀有狗
func xiyou_dog(dog map[string]interface{},configuration st.Configuration)bool  {

	rareDegree,_:=jsoniter.MarshalToString(dog["rareDegree"])
	amount:=jsoniter.Wrap(dog["amount"]).ToFloat32()
	generation,_:=jsoniter.MarshalToString(dog["generation"])
	if rareDegree=="1"&&generation=="0"&&amount<=configuration.XIYOU0_1DOG_0_PRICE{

		return true
	}
	return false
}
//普通狗
func putong_dog(dog map[string]interface{},configuration st.Configuration)bool  {

	rareDegree,_:=jsoniter.MarshalToString(dog["rareDegree"])
	amount:=jsoniter.Wrap(dog["amount"]).ToFloat32()
	generation,_:=jsoniter.MarshalToString(dog["generation"])
	if rareDegree=="0"&&generation=="0"&&amount<=configuration.PUTONG0_1DOG_0_PRICE{

		return true
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
							if shenhua_dog(s,configuration)||putong_dog(s,configuration)||shishi_dog(s,configuration)||zhuoyue_dog(s,configuration)||xiyou_dog(s,configuration)||chuanshuo_dog(s,configuration){
								codes :=get_code()
								json,_ :=simplejson.NewJson([]byte(codes))
								if json !=nil{
									seed :=json.Get("seed").MustString()
									code :=json.Get("code").MustString()
									bres :=bug_dog(s["petId"].(string),s["amount"].(string),seed,code,s["validCode"].(string),configuration)
									res,_ :=simplejson.NewJson([]byte(bres))
									if res!=nil {
										errorNo :=json.Get("errorNo").MustString()
										errorMsg :=json.Get("errorMsg").MustString()
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
											fmt.Print("抢到狗狗啦！！！！！！")
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

		seed,err:=js.Get("data").Get("seed").String()
		if err!=nil {
			fmt.Print(err)
		}
		if seed=="" {
			fmt.Print("百度莱茨狗验证码接口出现问题","\n")
			return
		}
		imgbase64,err:=js.Get("data").Get("img").String()
		if err!=nil {
			fmt.Print(err)
		}
		key:=configuration.KEY
		code:=lujun_api(key,imgbase64)
		fmt.Print("验证码="+code,"====>seed="+seed,"\n")
		if code!="" {
			jsonstr:=`{"code":"`+code+`","seed":"`+seed+`"}`
			len :=code_list.Len()
			if(len>=500){
				code_list.Init()
			}
			code_list.PushBack(jsonstr)
		}

	}

}
//验证码识别接口
func lujun_api(key string,img64 string) string {
	url := "http://www.popyelove.com:8888/?key="+key+"&img="+img64
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
//自动打码服务
func Timer2(configuration st.Configuration)  {
	ticker := time.NewTicker(5000* time.Millisecond)
	for _ = range ticker.C {
		print_code(configuration)
	}
}
var config string
var code_list *list.List
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
