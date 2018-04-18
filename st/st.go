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
	PAGE int `yaml:"PAGE"`
	PAGE_SIZE int `yaml:"PAGE_SIZE"`
	SORT_TYPE string `yaml:"SORT_TYPE"`

	CHUANSHUO_SWITCH int `yaml:"CHUANSHUO_SWITCH"`
	CHUANSHU1_SWITCH int `yaml:"CHUANSHU1_SWITCH"`
	CHUANSHU2_SWITCH int `yaml:"CHUANSHU2_SWITCH"`
	CHUANSHU3_SWITCH int `yaml:"CHUANSHU3_SWITCH"`
	CHUANSHUO0_8DOG_0_PRICE float32 `yaml:"CHUANSHUO0_8DOG_0_PRICE"`
	CHUANSHUO1_8DOG_0_PRICE float32 `yaml:"CHUANSHUO1_8DOG_0_PRICE"`
	CHUANSHUO2_8DOG_0_PRICE float32 `yaml:"CHUANSHUO2_8DOG_0_PRICE"`
	CHUANSHUO3_8DOG_0_PRICE float32 `yaml:"CHUANSHUO3_8DOG_0_PRICE"`

	GOD0_6_SWITCH int `yaml:"GOD0_6_SWITCH"`
	GOD0_6DOG_0_PRICE float32 `yaml:"GOD0_6DOG_0_PRICE"`  	//0代神话0分钟价格
	GOD0_6DOG_24_PRICE float32 `yaml:"GOD0_6DOG_24_PRICE"` 	//0代神话24小时价格
	GOD0_6DOG_2_PRICE float32  `yaml:"GOD0_6DOG_2_PRICE"`	//0代神话2天价格

	GOD1_6_SWITCH int `yaml:"GOD1_6_SWITCH"`
	GOD1_6DOG_0_PRICE float32 `yaml:"GOD1_6DOG_0_PRICE"`	//1代神话0天价格
	GOD1_6DOG_2_PRICE float32 `yaml:"GOD1_6DOG_2_PRICE"`	//1代神话2天价格
	GOD1_6DOG_4_PRICE float32 `yaml:"GOD1_6DOG_4_PRICE"`	//1代神话4天价格

	GOD2_6_SWITCH int `yaml:"GOD2_6_SWITCH"`
	GOD2_6DOG_0_PRICE float32 `yaml:"GOD2_6DOG_0_PRICE"`	//2代神话0天价格
	GOD2_6DOG_4_PRICE float32 `yaml:"GOD2_6DOG_4_PRICE"`	//2代神话4天价格
	GOD2_6DOG_6_PRICE float32 `yaml:"GOD2_6DOG_6_PRICE"`	//2代神话6天价格

	GOD3_6_SWITCH int `yaml:"GOD3_6_SWITCH"`
	GOD3_6DOG_0_PRICE float32 `yaml:"GOD3_6DOG_0_PRICE"`	//3代神话0天价格
	GOD3_6DOG_6_PRICE float32 `yaml:"GOD3_6DOG_6_PRICE"`	//3代神话6天价格
	GOD3_6DOG_8_PRICE float32 `yaml:"GOD3_6DOG_8_PRICE"`	//3代神话8天价格

	GOD0_7_SWITCH int `yaml:"GOD0_7_SWITCH"`
	GOD0_7DOG_0_PRICE float32 `yaml:"GOD0_7DOG_0_PRICE"`	//0代七稀神话0分钟
	GOD0_7DOG_24_PRICE float32 `yaml:"GOD0_7DOG_24_PRICE"`	//0代七稀神话24小时
	GOD0_7DOG_2_PRICE float32 `yaml:"GOD0_7DOG_2_PRICE"`	//0代七稀神话2天

	GOD1_7_SWITCH int `yaml:"GOD1_7_SWITCH"`
	GOD1_7DOG_0_PRICE float32 `yaml:"GOD1_7DOG_0_PRICE"`	//1代七稀神话0分钟
	GOD1_7DOG_2_PRICE float32 `yaml:"GOD1_7DOG_2_PRICE"`	//1代七稀神话2
	GOD1_7DOG_4_PRICE float32 `yaml:"GOD1_7DOG_4_PRICE"`	//1代七稀神话4天

	SHISHI0_5_SWITCH int `yaml:"SHISHI0_5_SWITCH"`
	SHISHI0_5DOG_0_PRICE float32 `yaml:"SHISHI0_5DOG_0_PRICE"`	//0代五稀史诗0天
	SHISHI0_5DOG_24_PRICE float32 `yaml:"SHISHI0_5DOG_24_PRICE"`	//0代五稀史诗24
	SHISHI_5BIRTHDAY_PRICE float32 `yaml:"SHISHI_5BIRTHDAY_PRICE"`

	SHISHI0_4_SWITCH int `yaml:"SHISHI0_4_SWITCH"`
	SHISHI0_4DOG_0_PRICE float32 `yaml:"SHISHI0_4DOG_0_PRICE"`
	SHISHI0_4DOG_24_PRICE float32 `yaml:"SHISHI0_4DOG_24_PRICE"`
	SHISHI_4BIRTHDAY_PRICE float32 `yaml:"SHISHI_4BIRTHDAY_PRICE"`

	ZHUEYUE0_2_SWITCH int `yaml:"ZHUEYUE0_2_SWITCH"`
	ZHUEYUE0_2DOG_0_PRICE float32 `yaml:"ZHUEYUE0_2DOG_0_PRICE"` 	//0,0卓越
	ZHUEYUE_BIRTHDAY_PRICE float32 `yaml:"ZHUEYUE_BIRTHDAY_PRICE"`

	XIYOU0_1_SWITCH int `yaml:"XIYOU0_1_SWITCH"`
	XIYOU0_1DOG_0_PRICE float32 `yaml:"XIYOU0_1DOG_0_PRICE"` 	//00 稀有
	XIYOU_BIRTHDAY_PRICE float32 `yaml:"XIYOU_BIRTHDAY_PRICE"` 	//00 稀有

	PUTONG0_0_SWITCH int `yaml:"PUTONG0_0_SWITCH"`
	PUTONG0_0DOG_0_PRICE float32 `yaml:"PUTONG0_0DOG_0_PRICE"` 	//00 普通
	PUTONG_BIRTHDAY_PRICE float32 `yaml:"PUTONG_BIRTHDAY_PRICE"`
}

func (configuration *Configuration) GetConf(config string) *Configuration {

	yamlFile, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Print("配置文件路径不对请核查", err)
		fmt.Print("\n")
		fmt.Printf("请输入你的配置文件的绝对路径(例如：D:/file/conf.yaml)：")
		fmt.Scanln(&config)
		configuration.GetConf(config)
	}
	err = yaml.Unmarshal(yamlFile, configuration)
	if err != nil {
		fmt.Print("配置文件不合法，请检查配置文件内容", err,"\n")
		fmt.Print("\n")
	}

	return configuration
}

