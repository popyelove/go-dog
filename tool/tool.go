package tool
import (
	"bytes"
	"io/ioutil"
	"net/http"
	"go-dog/st"
)
//上架
func Sale(petid string,amount string,configuration st.Configuration) string {
	url := "https://pet-chain.baidu.com/data/market/salePet"
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
	req.Header.Set("Cookie",configuration.COOKIE[0])

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
func Unsale(petid string,configuration st.Configuration) string {
	url := "https://pet-chain.baidu.com/data/market/unsalePet"
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
	req.Header.Set("Cookie",configuration.COOKIE[0])

	client := &http.Client{}
	resp,_:= client.Do(req)
	if resp !=nil{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
	return ""
}
