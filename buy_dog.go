package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"go-dog/st"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/http"
	"time"
)
func http_post1(url string, jsonStr []byte, configuration st.Configuration, ch chan string) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", configuration.COOKIE[account_index1])
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
//是否切换账号
func switch_account1(json string, petid string, amount string, configuration st.Configuration) {
	res, _ := simplejson.NewJson([]byte(json))
	if res != nil {
		errorNo := res.Get("errorNo").MustString()
		if (errorNo == "00" || errorNo == "10001" || errorNo == "10003") {
			account_index1 += 1
			if (account_index1 >= len(configuration.COOKIE)) {
				account_index1 = 0
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
		//if(errorNo=="10002"){
		//	m := gomail.NewMessage()
		//	m.SetHeader("From",configuration.QQ_EMAIL)
		//	m.SetHeader("To",configuration.QQ_EMAIL)
		//	m.SetAddressHeader("Cc", configuration.QQ_EMAIL, "莱茨狗")
		//	m.SetHeader("Subject", "被别人抢购成功")
		//	html:=`<a href=https://pet-chain.duxiaoman.com/chain/detail?channel=market&petId=`+petid+`>详情地址</a><br>狗狗价格：`+amount+"微"
		//	m.SetBody("text/html", html)
		//	d:=gomail.NewDialer("smtp.qq.com", 587, configuration.QQ_EMAIL,configuration.QQ_AUTH_PWD)
		//	d.DialAndSend(m);
		//}

	}
}
//下单买狗
func bugdog(petId string, amount string, seed string, code string, validCode string, configuration st.Configuration) string {
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
	go http_post1(url, jsonStr, configuration, ch_run)
	select {
	case res := <-ch_run:
		go switch_account1(res, petId, amount, configuration)
		return res
	case <-time.After(buy_dog_timeout1 * time.Second):
		fmt.Println("交易火爆中，请稍后再试。。。！")
		return ""
	}
	return ""
}
//下单超时时间
var buy_dog_timeout1 time.Duration = 15
var config1 string
//设置初始账号
var account_index1 int = 0
func main() {
	f := flag.String("f", "", "配置文件路径")
	flag.Parse() //解析输入的参数
	if (*f == "") {
		fmt.Printf("请输入你的配置文件的绝对路径(例如：D:/file/conf.yaml)：")
		fmt.Scanln(&config1)
	} else {
		config1 = *f
	}
	var configuration st.Configuration
	configuration.GetConf(config1)
	res := bugdog(
		"2287954367406252036",
		"1",
		"captchaservice30386131385273482b684c7a386f33782b345071765756325765356161794d344c626a75566d6251506d376f5047516c2f525a6875387648476878716b4c75393331534c5a6c5051794378573877734a693561363634303277614842742f52314f4a4a465271454548426a62686555506b5a473251742b627766474a676b50774a526567547179736f524a726c5632517755596c6553666a47794c314173625a614434782b56694335633869415151614c434147323874397732592f5733475451586e794779593562693846657631465179432f433863536e3042456639316771583856667037484848324264672b59644f31665462597a502b5243515a3351716f575348725458625a68504844443856614f326157305a3539395845672b54506b6f573562686f6833685763384e6866394670",
		"zydv",
		"",
		configuration)
	fmt.Print(res)
}
