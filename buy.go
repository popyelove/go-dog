package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/garyburd/redigo/redis"
	"github.com/json-iterator/go"
	"go-dog/st"
	"gopkg.in/gomail.v2"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func http_post(url string, jsonStr []byte, configuration st.Configuration, ch chan string) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", configuration.COOKIE[account_index])
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		ch <- string(body)
	} else {
		ch <- ""
	}
	return
}
func http_open_post(url string, jsonStr []byte, configuration st.Configuration, ch chan string) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		ch <- string(body)
	} else {
		ch <- ""
	}
	return
}

//获取狗的列表
func dog_list(configuration st.Configuration) string {
	if (index_dog >= len(dog_filter)) {
		index_dog = 0;
	}
	//Durl := "https://pet-chain.duxiaoman.com/data/market/queryPetsOnSale"
	url := "https://pet-chain.duxiaoman.com/data/market/queryPetsOnSale"
	var jsonStr = []byte(`
	{
		"pageNo":` + strconv.Itoa(index_page) + `,
		"pageSize":` + strconv.Itoa(configuration.PAGE_SIZE) + `,
		"querySortType":"` + configuration.SORT_TYPE + `",
		"petIds":[],
		"lastAmount":"",
		"lastRareDegree":"",
		"filterCondition":"{` + dog_filter[index_dog] + `}",
		"appId":1,
		"tpl":"",
		"type":null,
		"requestId":1522231859931,
		"timeStamp":null,
		"nounce":null,
		"token":null
    }
	`)
	ch_run := make(chan string)
	go http_open_post(url, jsonStr, configuration, ch_run)
	select {
	case res := <-ch_run:
		if (res != "") {
			fmt.Print("抢狗进行中...", time.Now())
			fmt.Print("\n")
		}
		return res
	case <-time.After(dog_list_timeout * time.Second):
		fmt.Println("拉取狗狗列表接口超时,请检查刷狗频率参数是否过小。。。。\n")
		return ""
	}
	return ""
}

//下单买狗
func bug_dog(petId string, amount string, seed string, code string, validCode string, configuration st.Configuration) string {
	type tiaojian struct {
		Petid     string `json:"petId"`
		Amount    string `json:"amount"`
		Seed      string `json:"seed"`
		Captcha   string `json:"captcha"`
		ValidCode string `json:"validCode"`
		RequestId string `json:"requestId"`
		Appid     string `json:"appId"`
		Tpl       string `json:"tpl"`
	}
	json_tiaojian := tiaojian{Petid: petId, Amount: amount, Seed: seed, Captcha: code, ValidCode: validCode, RequestId: "1520241678619", Appid: "1", Tpl: ""}
	url := "https://pet-chain.duxiaoman.com/data/txn/sale/create"
	jsonStr, _ := json.Marshal(json_tiaojian)
	ch_run := make(chan string)
	go http_post(url, jsonStr, configuration, ch_run)
	select {
	case res := <-ch_run:
		go switch_account(res, petId, amount, configuration)
		return res
	case <-time.After(buy_dog_timeout * time.Second):
		fmt.Println("交易火爆中，请稍后再试。。。！")
		return ""
	}
	return ""
}
func bugbaby_dog(petId string, amount string, seed string, code string, validCode string, configuration st.Configuration) string {
	type tiaojian struct {
		Petid       string `json:"petId"`
		Amount      string `json:"amount"`
		Seed        string `json:"seed"`
		Captcha     string `json:"captcha"`
		ValidCode   string `json:"validCode"`
		RequestId   string `json:"requestId"`
		Appid       string `json:"appId"`
		Tpl         string `json:"tpl"`
		SenderPetId string `json:"senderPetId"`
	}
	json_tiaojian := tiaojian{Petid: petId, Amount: amount, Seed: seed, Captcha: code, ValidCode: validCode, RequestId: "1520241678619", Appid: "1", Tpl: "", SenderPetId: configuration.MAKE_BABY_PETID}
	url := "https://pet-chain.duxiaoman.com/data/txn/breed/create"
	jsonStr, _ := json.Marshal(json_tiaojian)
	ch_run := make(chan string)
	go http_post(url, jsonStr, configuration, ch_run)
	select {
	case res := <-ch_run:
		go switch_account(res, petId, amount, configuration)
		return res
	case <-time.After(buy_dog_timeout * time.Second):
		fmt.Println("交易火爆中，请稍后再试。。。！")
		return ""
	}
	return ""
}

//获取狗的稀有属性
func get_dog_rareDegree(petid string, configuration st.Configuration) (int, int) {
	url := "https://pet-chain.duxiaoman.com/data/pet/queryPetById"
	var jsonStr = []byte(`{
        "petId":` + petid + `,
        "requestId":1520241678619,
        "appId":1,
        "tpl":"",
        "timeStamp":"",
        "nounce":"",
        "token":""
    }`)

	ch_run := make(chan string)
	go http_post(url, jsonStr, configuration, ch_run)
	select {
	case res := <-ch_run:
		js, err := simplejson.NewJson([]byte(res))
		if err != nil {
			fmt.Print(err)
		}
		if js == nil {
			return 0, 0
		}
		count_rareDegree := 0
		dogtype := 0
		for i := 0; i < 8; i++ {
			s := js.Get("data").Get("attributes").GetIndex(i).MustMap()
			if s["rareDegree"] == "稀有" {
				count_rareDegree = count_rareDegree + 1
			}
			if (contain(s["value"], configuration.BODY_TYPE)) {
				dogtype += 1
			}
			if (contain(s["value"], configuration.EYES_TYPE)) {
				dogtype += 1
			}
			if (contain(s["value"], configuration.MOUTH_TYPE)) {
				dogtype += 1
			}
			if (contain(s["value"], configuration.BODY_COLOR)) {
				dogtype += 1
			}
		}
		return count_rareDegree, dogtype
	case <-time.After(get_dog_rare_timeout * time.Second):
		fmt.Println("获取狗的稀有属性超时。。。。。\n")
		return 0, 0
	}
	return 0, 0
}

//传说狗
func chuanshuo_dog(dog map[string]interface{}, configuration st.Configuration) bool {
	rareDegree, _ := jsoniter.MarshalToString(dog["rareDegree"]) //稀有度
	amount := jsoniter.Wrap(dog["amount"]).ToFloat32()           //价额
	timeLeft := jsoniter.Wrap(dog["coolingInterval"]).ToString() //休息时间
	generation, _ := jsoniter.MarshalToString(dog["generation"]) //代数
	if (rareDegree == "5") {
		if (generation == "0") {
			if (amount <= configuration.CHUANSHUO0_8DOG_0_PRICE && timeLeft == "0分钟") {
				return true
			}
			if (amount <= configuration.CHUANSHUO0_8DOG_24_PRICE && timeLeft == "24小时") {
				return true
			}
			if (amount <= configuration.CHUANSHUO0_8DOG_2_PRICE && timeLeft == "2天") {
				return true
			}
		} else {
			if (amount <= configuration.CHUANSHUO_8DOG_OLD_PRICE) {
				return true
			}
		}
	}
	return false
}

//神话狗
func shenhua_dog(dog map[string]interface{}, configuration st.Configuration) bool {
	rareDegree, _ := jsoniter.MarshalToString(dog["rareDegree"])                     //稀有度
	amount := jsoniter.Wrap(dog["amount"]).ToFloat32()                               //价额
	timeLeft := jsoniter.Wrap(dog["coolingInterval"]).ToString()                     //休息时间
	generation, _ := jsoniter.MarshalToString(dog["generation"])                     //代数
	rareDegrees, dogtype := get_dog_rareDegree(dog["petId"].(string), configuration) //属性稀有个数
	if (rareDegrees == 6 && rareDegree == "4" && configuration.GOD_6DOG_SWITCH == 1) {
		//六稀神话
		if (generation == "0") {
			//0代神话0分钟满足特殊属性价格
			if (amount <= configuration.GOD0_6_0SPECIAL_PRICE && timeLeft == "0分钟" && dogtype == count_raredegree) {
				return true
			}
			//0代神话24满足特殊属性价格
			if (amount <= configuration.GOD0_6_24SPECIAL_PRICE && timeLeft == "24小时" && dogtype == count_raredegree) {
				return true
			}
			//0代神话2天满足特殊属性价格
			if (amount <= configuration.GOD0_6_2SPECIAL_PRICE && timeLeft == "2天" && dogtype == count_raredegree) {
				return true
			}
			//0代神话价格
			if (amount <= configuration.GOD0_6DOG_0_PRICE && timeLeft == "0分钟") {
				return true
			}
			if (amount <= configuration.GOD0_6DOG_24_PRICE && timeLeft == "24小时") {
				return true
			}
			if (amount <= configuration.GOD0_6DOG_2_PRICE && timeLeft == "2天") {
				return true
			}
		}
		if (generation == "1") {
			//1代神话0分钟满足特殊属性价格
			if (amount <= configuration.GOD1_6_0SPECIAL_PRICE && timeLeft == "0分钟" && dogtype == count_raredegree) {
				return true
			}
			//1代神话2满足特殊属性价格
			if (amount <= configuration.GOD1_6_2SPECIAL_PRICE && timeLeft == "2天" && dogtype == count_raredegree) {
				return true
			}
			//1代神话4天满足特殊属性价格
			if (amount <= configuration.GOD1_6_4SPECIAL_PRICE && timeLeft == "4天" && dogtype == count_raredegree) {
				return true
			}
			//一代神话价格配置
			if (amount <= configuration.GOD1_6DOG_0_PRICE && timeLeft == "0分钟") {
				return true
			}
			if (amount <= configuration.GOD1_6DOG_2_PRICE && timeLeft == "2天") {
				return true
			}
			if (amount <= configuration.GOD1_6DOG_4_PRICE && timeLeft == "4天") {
				return true
			}
		}
		//大于一代6神话
		if (generation != "0" && generation != "1") {
			///大于一代6神话满足特殊属性的
			if (amount <= configuration.GOD_6DOG_OLD1_SPECIAL_PRICE && dogtype == count_raredegree) {
				return true
			}
			//大于一代6神话价格
			if (amount <= configuration.GOD_6DOG_OLD1_PRICE) {
				return true
			}
		}
	}
	if (rareDegrees == 7 && rareDegree == "4" && configuration.GOD_7DOG_SWITCH == 1) {
		//七夕神话
		if (generation == "0") {
			//0代7神话0分钟满足特殊属性价格
			if (amount <= configuration.GOD0_7_0SPECIAL_PRICE && timeLeft == "0分钟" && dogtype == count_raredegree) {
				return true
			}
			//0代7神话24小时满足特殊属性价格
			if (amount <= configuration.GOD0_7_24SPECIAL_PRICE && timeLeft == "24小时" && dogtype == count_raredegree) {
				return true
			}
			//0代7神话2天满足特殊属性价格
			if (amount <= configuration.GOD0_7_2SPECIAL_PRICE && timeLeft == "2天" && dogtype == count_raredegree) {
				return true
			}
			//0代神话价格
			if (amount <= configuration.GOD0_7DOG_0_PRICE && timeLeft == "0分钟") {
				return true
			}
			if (amount <= configuration.GOD0_7DOG_24_PRICE && timeLeft == "24小时") {
				return true
			}
			if (amount <= configuration.GOD0_7DOG_2_PRICE && timeLeft == "2天") {
				return true
			}
		} else {
			if (amount <= configuration.GOD_7DOG_OLD_SPECIAL_PRICE && dogtype == count_raredegree) {
				return true
			}
			if (amount <= configuration.GOD_7DOG_OLD_PRICE) {
				return true
			}
		}

	}

	return false
}

//神话狗
func shenhuababy_dog(dog map[string]interface{}, configuration st.Configuration) bool {
	rareDegree, _ := jsoniter.MarshalToString(dog["rareDegree"]) //稀有度
	amount := jsoniter.Wrap(dog["amount"]).ToFloat32()           //价额
	//timeLeft := jsoniter.Wrap(dog["coolingInterval"]).ToString()                     //休息时间
	//generation, _ := jsoniter.MarshalToString(dog["generation"])                     //代数
	rareDegrees, _ := get_dog_rareDegree(dog["petId"].(string), configuration) //属性稀有个数
	if (rareDegrees == 6 && rareDegree == "4" && configuration.GOD_6DOG_SWITCH == 1) {
		//六稀神话
		if amount <= configuration.GOD_6DOG_BABY_PRICE {
			return true
		}
	}
	if (rareDegrees == 7 && rareDegree == "4" && configuration.GOD_7DOG_SWITCH == 1) {
		//七夕神话
		if amount <= configuration.GOD_7DOG_BABY_PRICE {
			return true
		}
	}

	return false
}

//史诗狗
func shishi_dog(dog map[string]interface{}, configuration st.Configuration) bool {
	rareDegree, _ := jsoniter.MarshalToString(dog["rareDegree"])
	amount := jsoniter.Wrap(dog["amount"]).ToFloat32()
	timeLeft := jsoniter.Wrap(dog["coolingInterval"]).ToString()
	rareDegrees, dogtype := get_dog_rareDegree(dog["petId"].(string), configuration)
	generation, _ := jsoniter.MarshalToString(dog["generation"])
	id, _ := jsoniter.MarshalToString(dog["id"])
	//五稀史诗
	if (rareDegrees == 5 && rareDegree == "3" && configuration.SHISHI_5_SWITCH == 1) {
		if (generation == "0") {
			//0代0分钟满足特殊属性
			if (amount <= configuration.SHISHI0_5_0SPECIAL_PRICE && timeLeft == "0分钟" && dogtype == count_raredegree) {
				return true
			}
			//0代24满足特殊属性
			if (amount <= configuration.SHISHI0_5_24SPECIAL_PRICE && timeLeft == "24小时" && dogtype == count_raredegree) {
				return true
			}
			//0代2满足特殊属性
			if (amount <= configuration.SHISHI0_5_2SPECIAL_PRICE && timeLeft == "2天" && dogtype == count_raredegree) {
				return true
			}
			if (amount <= configuration.SHISHI0_5DOG_0_PRICE && timeLeft == "0分钟") {
				return true
			}
			if (amount <= configuration.SHISHI0_5DOG_24_PRICE && timeLeft == "24小时") {
				return true
			}
			if (amount <= configuration.SHISHI_5BIRTHDAY_PRICE && validate(id)) {
				return true
			}
		} else {

			if (amount <= configuration.SHISHI_5DOG_OLD_SPECIAL_PRICE && dogtype == count_raredegree) {
				return true
			}
			if (amount <= configuration.SHISHI_5DOG_OLD_PRICE) {
				return true
			}
		}

	}
	//4稀有史诗
	if (rareDegrees == 4 && rareDegree == "3" && configuration.SHISHI_4_SWITCH == 1) {
		if (generation == "0") {
			//0代0分钟满足特殊属性
			if (amount <= configuration.SHISHI0_4_0SPECIAL_PRICE && timeLeft == "0分钟" && dogtype == count_raredegree) {
				return true
			}
			//0代24满足特殊属性
			if (amount <= configuration.SHISHI0_4_24SPECIAL_PRICE && timeLeft == "24小时" && dogtype == count_raredegree) {
				return true
			}
			//0代2天满足特殊属性
			if (amount <= configuration.SHISHI0_4_2SPECIAL_PRICE && timeLeft == "2天" && dogtype == count_raredegree) {
				return true
			}
			if (amount <= configuration.SHISHI0_4DOG_0_PRICE && timeLeft == "0分钟") {
				return true
			}
			if (amount <= configuration.SHISHI0_4DOG_24_PRICE && timeLeft == "24小时") {
				return true
			}
			if (amount <= configuration.SHISHI_4BIRTHDAY_PRICE && validate(id)) {
				return true
			}
		} else {
			if (amount <= configuration.SHISHI_4DOG_OLD_SPECIAL_PRICE && dogtype == count_raredegree) {
				return true
			}
			if (amount <= configuration.SHISHI_4DOG_OLD_PRICE) {
				return true
			}
		}

	}
	return false
}

//卓越狗
func zhuoyue_dog(dog map[string]interface{}, configuration st.Configuration) bool {
	rareDegree, _ := jsoniter.MarshalToString(dog["rareDegree"])
	amount := jsoniter.Wrap(dog["amount"]).ToFloat32()
	generation, _ := jsoniter.MarshalToString(dog["generation"])
	id, _ := jsoniter.MarshalToString(dog["id"])
	_, dogtype := get_dog_rareDegree(dog["petId"].(string), configuration)
	timeLeft := jsoniter.Wrap(dog["coolingInterval"]).ToString()

	if rareDegree == "2" && generation == "0" && amount <= configuration.ZHUOYUE0_0SPECIAL_PRICE && timeLeft == "0分钟" && dogtype == count_raredegree {

		return true
	}
	if rareDegree == "2" && generation == "0" && amount <= configuration.ZHUEYUE0_2DOG_0_PRICE && timeLeft == "0分钟" {

		return true
	}
	if rareDegree == "2" && amount <= configuration.ZHUEYUE_BIRTHDAY_PRICE && validate(id) {

		return true
	}
	if rareDegree == "2" && amount <= configuration.ZHUEYUE_GOOD_NUMBER_PRICE && good_num(id) {

		return true
	}
	//大于0代卓越狗
	if rareDegree == "2" && generation != "0" && amount <= configuration.ZHUOYUE_OLDER0_PRICE {
		return true
	}
	return false
}

//稀有狗
func xiyou_dog(dog map[string]interface{}, configuration st.Configuration) bool {
	rareDegree, _ := jsoniter.MarshalToString(dog["rareDegree"])
	amount := jsoniter.Wrap(dog["amount"]).ToFloat32()
	generation, _ := jsoniter.MarshalToString(dog["generation"])
	id, _ := jsoniter.MarshalToString(dog["id"])
	_, dogtype := get_dog_rareDegree(dog["petId"].(string), configuration)
	timeLeft := jsoniter.Wrap(dog["coolingInterval"]).ToString()
	if rareDegree == "1" && generation == "0" && amount <= configuration.XIYOU0_1DOG_0_PRICE && timeLeft == "0分钟" {

		return true
	}
	if rareDegree == "1" && amount <= configuration.XIYOU_BIRTHDAY_PRICE && validate(id) {

		return true
	}
	if rareDegree == "1" && amount <= configuration.XIYOU_GOOD_NUMBER_PRICE && good_num(id) {

		return true
	}
	if rareDegree == "1" && generation == "0" && amount <= configuration.XIYOU0_0SPECIAL_PRICE && timeLeft == "0分钟" && dogtype == count_raredegree {

		return true
	}
	//大于0代稀有狗
	if rareDegree == "1" && generation != "0" && amount <= configuration.XIYOU_OLDER0_DOG_PRICE {
		return true
	}
	return false
}

//普通狗
func putong_dog(dog map[string]interface{}, configuration st.Configuration) bool {
	rareDegree, _ := jsoniter.MarshalToString(dog["rareDegree"])
	id, _ := jsoniter.MarshalToString(dog["id"])
	amount := jsoniter.Wrap(dog["amount"]).ToFloat32()
	generation, _ := jsoniter.MarshalToString(dog["generation"])
	timeLeft := jsoniter.Wrap(dog["coolingInterval"]).ToString()
	if rareDegree == "0" && generation == "0" && amount <= configuration.PUTONG0_0DOG_0_PRICE && timeLeft == "0分钟" {
		return true
	}
	if rareDegree == "0" && amount <= configuration.PUTONG_BIRTHDAY_PRICE && validate(id) {

		return true
	}
	if rareDegree == "0" && amount <= configuration.PUTONG_GOOD_NUMBER_PRICE && good_num(id) {

		return true
	}
	if rareDegree == "0" && generation != "0" && amount <= configuration.PUTONG_OLDER_DOG_PRICE {
		return true
	}
	return false
}

//获取验证吗
func get_code() string {
	code := code_list.Back()
	if code != nil {
		return code.Value.(string)
	}
	return "{}"
}

//获取火爆验证吗
func get_huobao_code() string {
	code := huobao_codes.Back()
	if code != nil {
		return code.Value.(string)
	}
	return "{}"
}

//传说
func dog_chuanshuo(dogs string, configuration st.Configuration) {
	js, _ := simplejson.NewJson([]byte(dogs))
	if js != nil {
		for i := 0; i < configuration.PAGE_SIZE; i++ {
			s := js.Get("data").Get("petsOnSale").GetIndex(i).MustMap()
			if s != nil {
				if chuanshuo_dog(s, configuration) {
					real_buy(s, configuration)
				}

			}

		}

	}
}

//神话
func dog_shenhua(dogs string, configuration st.Configuration) {
	js, _ := simplejson.NewJson([]byte(dogs))
	if js != nil {
		for i := 0; i < configuration.PAGE_SIZE; i++ {
			s := js.Get("data").Get("petsOnSale").GetIndex(i).MustMap()
			if s != nil {
				if shenhua_dog(s, configuration) {
					real_buy(s, configuration)
				}

			}

		}

	}
}

//
func dogbaby_shenhua(dogs string, configuration st.Configuration) {
	js, _ := simplejson.NewJson([]byte(dogs))
	if js != nil {
		for i := 0; i < configuration.PAGE_SIZE; i++ {
			s := js.Get("data").Get("pets4Breed").GetIndex(i).MustMap()
			if s != nil {
				if shenhuababy_dog(s, configuration) {
					realbaby_buy(s, configuration)
				}

			}

		}

	}
}

//史诗
func dog_shishi(dogs string, configuration st.Configuration) {
	js, _ := simplejson.NewJson([]byte(dogs))
	if js != nil {
		for i := 0; i < configuration.PAGE_SIZE; i++ {
			s := js.Get("data").Get("petsOnSale").GetIndex(i).MustMap()
			if s != nil {
				if shishi_dog(s, configuration) {
					real_buy(s, configuration)
				}

			}

		}

	}
}

//卓越
func dog_zhuoyue(dogs string, configuration st.Configuration) {
	js, _ := simplejson.NewJson([]byte(dogs))
	if js != nil {
		for i := 0; i < configuration.PAGE_SIZE; i++ {
			s := js.Get("data").Get("petsOnSale").GetIndex(i).MustMap()
			if s != nil {
				if zhuoyue_dog(s, configuration) {
					real_buy(s, configuration)
				}

			}

		}

	}
}

//稀有
func dog_xiyou(dogs string, configuration st.Configuration) {
	js, _ := simplejson.NewJson([]byte(dogs))
	if js != nil {
		for i := 0; i < configuration.PAGE_SIZE; i++ {
			s := js.Get("data").Get("petsOnSale").GetIndex(i).MustMap()
			if s != nil {
				if xiyou_dog(s, configuration) {
					real_buy(s, configuration)
				}

			}

		}

	}
}
func real_buy(s map[string]interface{}, configuration st.Configuration) {
	codes := get_code()
	json, _ := simplejson.NewJson([]byte(codes))
	if json != nil {
		seed := json.Get("seed").MustString()
		code := json.Get("code").MustString()
		bres := bug_dog(s["petId"].(string), s["amount"].(string), seed, code, s["validCode"].(string), configuration)
		res, _ := simplejson.NewJson([]byte(bres))
		if res != nil {
			errorNo := res.Get("errorNo").MustString()
			//errorMsg := res.Get("errorMsg").MustString()
			if (errorNo == "100" || errorNo == "101") {
				codes := get_code()
				json, _ := simplejson.NewJson([]byte(codes))
				if json != nil {
					seed := json.Get("seed").MustString()
					code := json.Get("code").MustString()
					bug_dog(s["petId"].(string), s["amount"].(string), seed, code, s["validCode"].(string), configuration)
				}
			}
			if errorNo == "08" {
				//交易火爆，区块链处理繁忙，请稍后再试
				codes := get_code()
				json, _ := simplejson.NewJson([]byte(codes))
				seed := json.Get("seed").MustString()
				code := json.Get("code").MustString()
				bug_dog(s["petId"].(string), s["amount"].(string), seed, code, s["validCode"].(string), configuration)
			}
			if errorNo == "10002" {
				//有人抢先下单啦
				fmt.Print("***********************被别人抢走了*****************")
			}
			if errorNo == "00" {
				//success
				fmt.Print("************************抢到狗狗啦！！！！！！！！！！", "\n", s)
			}
		}
	}
}

func realbaby_buy(s map[string]interface{}, configuration st.Configuration) {
	codes := get_code()
	json, _ := simplejson.NewJson([]byte(codes))
	if json != nil {
		seed := json.Get("seed").MustString()
		code := json.Get("code").MustString()
		bres := bugbaby_dog(s["petId"].(string), s["amount"].(string), seed, code, s["validCode"].(string), configuration)
		res, _ := simplejson.NewJson([]byte(bres))
		if res != nil {
			errorNo := res.Get("errorNo").MustString()
			//errorMsg := res.Get("errorMsg").MustString()
			if (errorNo == "100" || errorNo == "101") {
				codes := get_code()
				json, _ := simplejson.NewJson([]byte(codes))
				if json != nil {
					seed := json.Get("seed").MustString()
					code := json.Get("code").MustString()
					bug_dog(s["petId"].(string), s["amount"].(string), seed, code, s["validCode"].(string), configuration)
				}
			}
			if errorNo == "08" {
				//交易火爆，区块链处理繁忙，请稍后再试
				codes := get_code()
				json, _ := simplejson.NewJson([]byte(codes))
				seed := json.Get("seed").MustString()
				code := json.Get("code").MustString()
				bug_dog(s["petId"].(string), s["amount"].(string), seed, code, s["validCode"].(string), configuration)
			}
			if errorNo == "10002" {
				//有人抢先下单啦
				fmt.Print("***********************被别人抢走了*****************")
			}
			if errorNo == "00" {
				make_baby = "0"
				//success
				fmt.Print("************************抢到狗狗啦！！！！！！！！！！", "\n", s)
			}
		}
	}
}

//普通
func dog_putong(dogs string, configuration st.Configuration) {
	js, _ := simplejson.NewJson([]byte(dogs))
	if js != nil {
		for i := 0; i < configuration.PAGE_SIZE; i++ {
			s := js.Get("data").Get("petsOnSale").GetIndex(i).MustMap()
			if s != nil {
				if putong_dog(s, configuration) {
					fmt.Print(s)
					real_buy(s, configuration)
				}

			}

		}

	}
}

//循环刷狗函数
func do_always(configuration st.Configuration) {
	dogs := dog_list(configuration)
	if dogs != "" {
		if (index_dog >= len(dog_filter)) {
			index_dog = 0
		}
		flag := index_dog
		if (index_page >= configuration.PAGE) {
			index_page = 1
			index_dog += 1
		} else {
			index_page += 1
		}

		switch dog_filter[flag] {
		case "1:5":
			fmt.Print(dog_filter[flag])
			dog_chuanshuo(dogs, configuration)
		case "1:4":
			fmt.Print(dog_filter[flag])
			dog_shenhua(dogs, configuration)
		case "1:3":
			fmt.Print(dog_filter[flag])
			dog_shishi(dogs, configuration)
		case "1:2":
			fmt.Print(dog_filter[flag])
			dog_zhuoyue(dogs, configuration)
		case "1:1":
			fmt.Print(dog_filter[flag])
			dog_xiyou(dogs, configuration)
		case "1:0":
			fmt.Print(dog_filter[flag])
			dog_putong(dogs, configuration)
		default:

		}

	}

}

func print_code(configuration st.Configuration) {
	url := "https://pet-chain.duxiaoman.com/data/newCaptcha/getCaptchaUrl"
	var jsonStr = []byte(`{
							"appId":1,"tpl":"",
							"nounce":null,
							"phoneType":"ios"
							"requestId":1523433103485,
							"timeStamp":null,
							"token":null,
							"tpl":""
							}
						`)

	ch_run := make(chan string)
	go http_post(url, jsonStr, configuration, ch_run)
	select {
	case res := <-ch_run:
		js, err := simplejson.NewJson([]byte(res))
		if err != nil {
			fmt.Print("百度服务器繁忙！！！", "\n")
			return
		}
		var seed string
		var imgurl string
		seed, err = js.Get("data").String()
		if err != nil {
			fmt.Print("百度服务器繁忙。。。。。。。。。。", "\n")
			return
		}
		imgurl = seed
		var seeds []string
		seeds = strings.Split(seed, "&")
		seed = seeds[2][8:len(seeds[2])]
		key := configuration.KEY
		code_res := lujun_api(key, imgurl)
		if (code_res == "") {
			return
		}
		code := code_res
		fmt.Print("验证码=" + code)
		if code != "" {
			jsonstr := `{"code":"` + code + `","seed":"` + seed + `"}`
			len := code_list.Len()
			if (len >= code_num) {
				code_list.Init()
			}
			code_list.PushBack(jsonstr)
		}
	case <-time.After(dama_timeout * time.Second):
		fmt.Println("百度验证码接口超时......")
		return
	}

}

func print_huobao_code(configuration st.Configuration) {
	url := "https://pet-chain.duxiaoman.com/data/captcha/gen"
	var jsonStr = []byte(`{
							"requestId":1523433103485,
							"appId":1,"tpl":"",
							"timeStamp":null,
							"nounce":null,
							"token":null}
						`)

	ch_run := make(chan string)
	go http_post(url, jsonStr, configuration, ch_run)
	select {
	case res := <-ch_run:
		js, err := simplejson.NewJson([]byte(res))
		if err != nil {
			fmt.Print("百度服务器繁忙！！！", "\n")
			return
		}
		var seed string
		seed, err = js.Get("data").Get("seed").String()
		if err != nil {
			fmt.Print("百度服务器繁忙。。。。。。。。。。", "\n")
			return
		}
		imgbase64, err := js.Get("data").Get("img").String()
		if err != nil {
			fmt.Print(err)
			return
		}
		key := configuration.KEY
		code_res := lujun_api(key, imgbase64)
		if (code_res == "") {
			return
		}
		js_code, err := simplejson.NewJson([]byte(code_res))
		if err != nil {
			fmt.Print(err)
			return
		}
		status := js_code.Get("status").MustString()
		msg := js_code.Get("msg").MustString()
		if status == "error" {
			fmt.Print(msg, "\n")
			return
		}
		code := js_code.Get("captcha").MustString()
		fmt.Print("[验证码]="+code, "====>[seed]="+seed, "\n")
		if code != "" {
			jsonstr := `{"code":"` + code + `","seed":"` + seed + `"}`
			len := huobao_codes.Len()
			if (len >= code_num) {
				huobao_codes.Init()
			}
			huobao_codes.PushBack(jsonstr)
		}
	case <-time.After(dama_timeout * time.Second):
		fmt.Println("百度验证码接口超时......")
		return
	}

}
func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func remove_file(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Print(err)
	}
}

//验证码识别接口
func lujun_api(key string, imgurl string) string {
	url := imgurl
	name := getRandomString(10) + ".gif"
	out, _ := os.Create(name)
	resp, _ := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		_, err := io.Copy(out, bytes.NewReader(body))
		if err != nil {
			fmt.Print(err)
		}
	}
	out.Close()
	byte, err := ioutil.ReadFile(name)
	resopne, err := http.Post("http://127.0.0.1:9900/", "multipart/form-data", bytes.NewReader(byte))
	if err != nil {
		fmt.Println("err=", err)
	}
	defer func() {
		resopne.Body.Close()
		out.Close()
	}()
	body, err := ioutil.ReadAll(resopne.Body)
	if err != nil {
		fmt.Println(" post err=", err)
	}
	remove_file(name)
	return string(body)
}

//自动打码服务
func dama_code(configuration st.Configuration) {
	ticker := time.NewTicker(configuration.TIMECODE * time.Millisecond)
	for _ = range ticker.C {
		print_code(configuration)
	}

}
func huobao_dama(configuration st.Configuration) {
	ticker := time.NewTicker(dama_time * time.Millisecond)
	for _ = range ticker.C {
		print_huobao_code(configuration)
	}
}

//是否是生日号
func validate(no string) bool {
	reg := regexp.MustCompile(regular1)
	return reg.MatchString(no)
}

//靓号
func good_num(no string) bool {
	reg := regexp.MustCompile(regular2)
	return reg.MatchString(no)
}

const (
	regular1 = "^(19[8-9]{1}[0-9]{1}|20[0-2]{1}[0-9]{1})(1[0-2]|0[1-9])(0[1-9]|[1-2][0-9]|3[0-1])$"
	regular2 = `1{5}|2{5}|3{5}|4{5}|5{5}|6{5}|7{5}|8{5}|9{5}|0{5}`
)

//获得当前软件版本
func get_version() float64 {
	c, err := redis.Dial("tcp", redis_host)
	if (err != nil) {
		fmt.Print(err)
	}
	c.Do("AUTH", redis_pwd)
	version, err := redis.String(c.Do("GET", "version"))
	if (err != nil) {
		return 0
	}

	if version != "" {
		defer c.Close()
		version, _ := strconv.ParseFloat(version, 64)
		return version
	}
	return 0
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

//筛选刷狗条件
func dogfilter(chuanshuo_switch int, god_switch int, shishi_switch int, zhuoyue_switch int, xiyou_switch int, putong_switch int, dog_filter []string) []string {
	if (chuanshuo_switch == 1) {
		dog_filter = append(dog_filter, "1:5")
	}
	if (god_switch == 1) {
		dog_filter = append(dog_filter, "1:4")
	}
	if (shishi_switch == 1) {
		dog_filter = append(dog_filter, "1:3")
	}
	if (zhuoyue_switch == 1) {
		dog_filter = append(dog_filter, "1:2")
	}
	if (xiyou_switch == 1) {
		dog_filter = append(dog_filter, "1:1")
	}
	if (putong_switch == 1) {
		dog_filter = append(dog_filter, "1:0")
	}
	return dog_filter
}

//获取设置特有属性的条件
func get_raredegree_count(body_type []string, eyes_type []string, mouth_type []string, body_color []string) int {
	count := 0
	if (len(body_type) > 0) {
		count += 1
	}
	if (len(eyes_type) > 0) {
		count += 1
	}
	if (len(mouth_type) > 0) {
		count += 1
	}
	if (len(body_color) > 0) {
		count += 1
	}
	return count
}

//是否切换账号
func switch_account(json string, petid string, amount string, configuration st.Configuration) {
	res, _ := simplejson.NewJson([]byte(json))
	if res != nil {
		errorNo := res.Get("errorNo").MustString()
		if (errorNo == "00" || errorNo == "10001" || errorNo == "10003") {
			account_index += 1
			if (account_index >= len(configuration.COOKIE)) {
				account_index = 0
			}
		}
		//购买成功
		if (errorNo == "00") {
			m := gomail.NewMessage()
			m.SetHeader("From", configuration.QQ_EMAIL)
			m.SetHeader("To", configuration.QQ_EMAIL)
			m.SetAddressHeader("Cc", configuration.QQ_EMAIL, "莱茨狗")
			m.SetHeader("Subject", "莱茨狗订单通知")
			html := `<a href=https://pet-chain.duxiaoman.com/chain/detail?channel=market&petId=` + petid + `>详情地址</a><br>狗狗价格：` + amount + "微"
			m.SetBody("text/html", html)
			d := gomail.NewDialer("smtp.qq.com", 587, configuration.QQ_EMAIL, configuration.QQ_AUTH_PWD)
			d.DialAndSend(m);
		}
		//被别人购买
		if (errorNo == "10002") {
			m := gomail.NewMessage()
			m.SetHeader("From", configuration.QQ_EMAIL)
			m.SetHeader("To", configuration.QQ_EMAIL)
			m.SetAddressHeader("Cc", configuration.QQ_EMAIL, "莱茨狗")
			m.SetHeader("Subject", "被别人抢购成功")
			html := `<a href=https://pet-chain.duxiaoman.com/chain/detail?channel=market&petId=` + petid + `>详情地址</a><br>狗狗价格：` + amount + "微"
			m.SetBody("text/html", html)
			d := gomail.NewDialer("smtp.qq.com", 587, configuration.QQ_EMAIL, configuration.QQ_AUTH_PWD)
			d.DialAndSend(m);
		}

	}
}

var config string
var code_list *list.List
var huobao_codes *list.List
var dog_filter = []string{}

//从索引为0的狗扫描
var index_dog = 0

//初始索引
var index_page = 1

//打码间隔 毫秒
var dama_time time.Duration = 5000

//拉取狗列表超时时间秒
var dog_list_timeout time.Duration = 15

//下单超时时间
var buy_dog_timeout time.Duration = 15

//获取狗狗属性超时
var get_dog_rare_timeout time.Duration = 15

//打码超时时间
var dama_timeout time.Duration = 15

//当前版本
var version float64 = 1.2
var redis_host string = "127.0.0.1:6379"
var redis_pwd string = ""
var dama_host string = "http://127.0.0.1:8888/"
var code_num int = 50

//满足稀有属性的个数
var count_raredegree int = 0

//设置初始账号
var account_index int = 0

var make_baby string = "1"

func is_passed() {
	time1 := time.Now().Unix()
	time2 := time.Date(2019, 11, 29, 0, 0, 0, 0, time.Local).Unix()
	if ((time1 - time2) > 30*86400) {
		panic("软件已过期")
	}
}

//自动切换账号每隔一个半小时
func auto_switch_account(configuration st.Configuration) {
	ticker := time.NewTicker(30 * time.Minute)
	for _ = range ticker.C {
		account_index += 1
		if (account_index >= len(configuration.COOKIE)) {
			account_index = 0
		}
		fmt.Print("切换账号success，当前账号:", account_index)
	}

}

//繁殖狗狗列表
func dogbaby_list(configuration st.Configuration) string {
	url := "https://pet-chain.duxiaoman.com/data/market/breed/pets"
	var jsonStr = []byte(`
	{
		"pageNo":` + strconv.Itoa(index_page) + `,
		"pageSize":` + strconv.Itoa(configuration.PAGE_SIZE) + `,
		"querySortType":"` + configuration.SORT_TYPE + `",
		"petIds":[],
		"lastAmount":"",
		"lastRareDegree":"",
		"filterCondition":"{"1:4"}",
		"appId":1,
		"tpl":"",
		"type":null,
		"requestId":1522231859931,
		"timeStamp":null,
		"nounce":null,
		"token":null
    }
	`)
	ch_run := make(chan string)
	go http_open_post(url, jsonStr, configuration, ch_run)
	select {
	case res := <-ch_run:
		if (res != "") {
			fmt.Print("\n")
			fmt.Print("繁殖进行中...", time.Now())
			fmt.Print("\n")
		}
		return res
	case <-time.After(dog_list_timeout * time.Second):
		fmt.Println("拉取狗狗列表接口超时,请检查刷狗频率参数是否过小。。。。\n")
		return ""
	}
	return ""
}

//循环刷狗函数
func do_makebabys(configuration st.Configuration) {
	dogs := dogbaby_list(configuration)
	if dogs != "" {
		if (index_page >= configuration.PAGE) {
			index_page = 1
		} else {
			index_page += 1
		}
		dogbaby_shenhua(dogs, configuration)
	}

}
func main() {
	//初始化一个用来装验证码的容器
	code_list = list.New()
	huobao_codes = list.New()
	f := flag.String("f", "", "配置文件路径")
	flag.Parse() //解析输入的参数
	if (*f == "") {
		fmt.Printf("请输入你的配置文件的绝对路径(例如：D:/file/conf.yaml)：")
		fmt.Scanln(&config)
	} else {
		config = *f
	}
	var configuration st.Configuration
	configuration.GetConf(config)
	//初始化刷狗类型
	dog_filter = dogfilter(configuration.CHUANSHUO_SWITCH, configuration.GOD_SWITCH, configuration.SHISHI_SWITCH, configuration.ZHUOYUE_SWITCH, configuration.XIYOU_SWITCH, configuration.PUTONG_SWITCH, dog_filter)
	//初始化属性条件
	count_raredegree = get_raredegree_count(configuration.BODY_TYPE, configuration.EYES_TYPE, configuration.MOUTH_TYPE, configuration.BODY_COLOR)
	//是否过期
	//is_passed()
	//预先打码

	go dama_code(configuration)
	//每半小时切换一次账号
	go auto_switch_account(configuration)
	ticker := time.NewTicker(configuration.TIME * time.Millisecond)
	for _ = range ticker.C {
		go do_always(configuration)
		automakebaby := configuration.AUTO_MAKE_BABY
		if automakebaby == "1" && make_baby == "1" && account_index == 0 {
			go do_makebabys(configuration)
		}
	}
}
