package main
import (
	"net/http"
	"io/ioutil"
	"bytes"
	"github.com/bitly/go-simplejson"
	"encoding/json"
	"github.com/json-iterator/go"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"time"
	"go-dog/st"
	"strconv"
	"flag"
)
//获取狗的列表
func dog_list(configuration st.Configuration) string {
	url := "https://pet-chain.baidu.com/data/market/queryPetsOnSale"
	var jsonStr = []byte(`{
		"pageNo":1,
		"pageSize":`+strconv.Itoa(configuration.PAGE_SIZE)+`,
		"querySortType":"CREATETIME_DESC",
		"petIds":[],
		"lastAmount":"",
		"lastRareDegree":"",
		"filterCondition":"{\"1\":\"4\"}",
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
		fmt.Print("*")
		return string(body)
	}
	fmt.Println(time.Now())
	fmt.Print("\n")
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
		fmt.Print(string(body))
		fmt.Print("\n")
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
func shenhua_dog(dog map[string]interface{},configuration st.Configuration)bool  {
	rareDegree,_:=jsoniter.MarshalToString(dog["rareDegree"])//稀有度
	amount:=jsoniter.Wrap(dog["amount"]).ToFloat32()//价额
	timeLeft :=jsoniter.Wrap(dog["coolingInterval"]).ToString()//休息时间
	generation,_:=jsoniter.MarshalToString(dog["generation"])//代数
	rareDegrees,_:=get_dog_rareDegree(dog["petId"].(string),configuration) //属性稀有个数
	if(rareDegree=="4"){
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
	c,_:= redis.Dial("tcp", "127.0.0.1:6379")
	values, _ := redis.String(c.Do("rpop", "code_list"))
	if values!=""{
		defer c.Close()
		return values
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
							if shenhua_dog(s,configuration)||shishi_dog(s,configuration){
								codes :=get_code()
								json,_ :=simplejson.NewJson([]byte(codes))
								if json !=nil{
									seed :=json.Get("seed").MustString()
									code :=json.Get("code").MustString()
									bug_dog(s["petId"].(string),s["amount"].(string),seed,code,s["validCode"].(string),configuration)
								}

							}

					}

				}

			}

		}

}

func main(){
	config := flag.String("f", "", "配置文件")
	flag.Parse()
	configfile:=*config
	var  configuration st.Configuration
	configuration.GetConf(configfile)
	ticker := time.NewTicker(configuration.TIME* time.Millisecond)
	for _ = range ticker.C {
		go do_always(configuration)
	}

}
