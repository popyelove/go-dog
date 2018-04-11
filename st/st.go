package st

import (
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
	"time"
)
type Configuration struct {
	COOKIE string `yaml:"COOKIE"`
	KEY string `yaml:"KEY"`
	TIME time.Duration `yaml:"TIME"`
	TIMECODE time.Duration `yaml:"TIMECODE"`
	GOD0_6DOG_0_PRICE float32 `yaml:"GOD0_6DOG_0_PRICE"`  	//0代神话0分钟价格
	GOD0_6DOG_24_PRICE float32 `yaml:"GOD0_6DOG_24_PRICE"` 				//0代神话24小时价格
	GOD0_6DOG_2_PRICE float32  `yaml:"GOD0_6DOG_2_PRICE"`				//0代神话2天价格

	GOD1_6DOG_0_PRICE float32 `yaml:"GOD1_6DOG_0_PRICE"`	//1代神话0天价格
	GOD1_6DOG_2_PRICE float32 `yaml:"GOD1_6DOG_2_PRICE"`	//1代神话2天价格
	GOD1_6DOG_4_PRICE float32 `yaml:"GOD1_6DOG_4_PRICE"`	//1代神话4天价格

	GOD2_6DOG_0_PRICE float32 `yaml:"GOD2_6DOG_0_PRICE"`	//2代神话0天价格
	GOD2_6DOG_4_PRICE float32 `yaml:"GOD2_6DOG_4_PRICE"`	//2代神话4天价格
	GOD2_6DOG_6_PRICE float32 `yaml:"GOD2_6DOG_6_PRICE"`	//2代神话6天价格


	GOD3_6DOG_0_PRICE float32 `yaml:"GOD3_6DOG_0_PRICE"`	//3代神话0天价格
	GOD3_6DOG_6_PRICE float32 `yaml:"GOD3_6DOG_6_PRICE"`	//3代神话6天价格
	GOD3_6DOG_8_PRICE float32 `yaml:"GOD3_6DOG_8_PRICE"`	//3代神话8天价格

	GOD0_7DOG_0_PRICE float32 `yaml:"GOD0_7DOG_0_PRICE"`	//0代七稀神话0分钟
	GOD0_7DOG_24_PRICE float32 `yaml:"GOD0_7DOG_24_PRICE"`	//0代七稀神话24小时
	GOD0_7DOG_2_PRICE float32 `yaml:"GOD0_7DOG_2_PRICE"`	//0代七稀神话2天

	GOD1_7DOG_0_PRICE float32 `yaml:"GOD1_7DOG_0_PRICE"`	//1代七稀神话0分钟
	GOD1_7DOG_2_PRICE float32 `yaml:"GOD1_7DOG_2_PRICE"`	//1代七稀神话2
	GOD1_7DOG_4_PRICE float32 `yaml:"GOD1_7DOG_4_PRICE"`	//1代七稀神话4天

	SHISHI0_5DOG_0_PRICE float32 `yaml:"SHISHI0_5DOG_0_PRICE"`	//0代五稀史诗0天
	SHISHI0_5DOG_24_PRICE float32 `yaml:"SHISHI0_5DOG_24_PRICE"`	//0代五稀史诗24

	ZHUEYUE0_2DOG_0_PRICE float32 `yaml:"ZHUEYUE0_2DOG_0_PRICE"` 	//0,0卓越

	XIYOU0_1DOG_0_PRICE float32 `yaml:"XIYOU0_1DOG_0_PRICE"` 	//00 稀有

	PUTONG0_1DOG_0_PRICE float32 `yaml:"PUTONG0_1DOG_0_PRICE"` 	//00 普通
}

func (configuration *Configuration) GetConf() *Configuration {

	yamlFile, err := ioutil.ReadFile("/Users/admin/go/src/go-dog/conf.yaml")
	if err != nil {
		fmt.Print("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, configuration)
	if err != nil {
		fmt.Print("Unmarshal: %v", err)
	}

	return configuration
}

