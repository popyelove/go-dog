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
		"captchaservice376464663633553951756f6b626b4659466a576634572b6178455376454c30315a4362582f327761626c49732f58417a7134544b4b44422b4355716d395131303463746777465a4e66704831436c6c56304931466f626f456e67486b2b5a387577374644757074706f727a383370756b45576441462f63335348595975444b517355374633726c4664503565704e4d4f674f4e6f71534e485a523847543771736142557371395a6b45793747635542354a75454542354d797734713432674d77386e67636c41625979516f736a4e59633264554673523649786a6f30657248376344614d496f6a4146333832486664595835624d796b33785574415772733158674d773263557250742b4359307453776a634f7870475434747456544c75396f4552696d736d6555414d5566646f6d466351",
		"4lsb",
		"",
		configuration)
	fmt.Print(res)
}
