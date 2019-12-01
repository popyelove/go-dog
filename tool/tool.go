package tool
import (
	"bytes"
	"io/ioutil"
	"net/http"
	"fmt"
	"time"
)
//获取个人详情
func GetInfo(cookie string) string {
	url := "https://pet-chain.duxiaoman.com/data/user/get"
	var jsonStr = []byte(`{
	"appId":1,
	"nounce":null,
	"phoneType":"ios",
	"requestId":1575190291281,
	"timeStamp":null,
	"token":null,
	"tpl":""
		}`)
	req,_:= http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie",cookie)
	client := &http.Client{}
	resp,_:= client.Do(req)
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
	fmt.Println(time.Now())
	fmt.Print("\n")
	return ""
}
//查询列表 获取已上架的狗狗id amount
func GetList(cookie string,pageNo string) string {
	url := "https://pet-chain.duxiaoman.com/data/user/pet/list"
	var jsonStr = []byte(`{
		"pageNo":`+pageNo+`,
		"pageSize":10,
		"pageTotal":-1,
		"totalCount":0,
		"requestId":1524651063100,
		"appId":1,
		"tpl":"",
		"timeStamp":null,
		"nounce":null,
		"token":null
		}`)
	req,_:= http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie",cookie)
	client := &http.Client{}
	resp,_:= client.Do(req)
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
	fmt.Println(time.Now())
	fmt.Print("\n")
	return ""
}
//上架
func Sale(petid string,amount string,cookie string) string {
	url := "https://pet-chain.duxiaoman.com/data/market/sale/shelf/create"
	var jsonStr = []byte(`
	{
		"petId":`+petid+`,
		"amount":`+amount+`,
		"requestId":1521014055739,
		"appId":1,
		"tpl":"",
		"timeStamp":null,
		"nounce":null,
		"token":null
	}`)
	req,_:= http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie",cookie)

	client := &http.Client{}
	resp,_:= client.Do(req)
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
	return ""
}


//下架狗狗
func Unsale(petid string,cookie string) string {
	url := "https://pet-chain.duxiaoman.com/data/market/unsalePet"
	var jsonStr = []byte(`
	{
		"petId":`+petid+`,
		"requestId":1521014647899,
		"appId":1,
		"tpl":"",
		"timeStamp":null,
		"nounce":null,
		"token":null
	}
`)
	req,_:= http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie",cookie)

	client := &http.Client{}
	resp,_:= client.Do(req)
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
	return ""
}
