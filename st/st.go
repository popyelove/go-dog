package st

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Configuration struct {
	COOKIE          []string      `yaml:"COOKIE"`
	KEY             string        `yaml:"KEY"`
	TIME            time.Duration `yaml:"TIME"`
	TIMECODE        time.Duration `yaml:"TIMECODE"`
	SELL_TIME       time.Duration `yaml:"SELL_TIME"`
	PAGE            int           `yaml:"PAGE"`
	PAGE_SIZE       int           `yaml:"PAGE_SIZE"`
	SORT_TYPE       string        `yaml:"SORT_TYPE"`
	BODY_TYPE       []string      `yaml:"BODY_TYPE"`
	EYES_TYPE       []string      `yaml:"EYES_TYPE"`
	MOUTH_TYPE      []string      `yaml:"MOUTH_TYPE"`
	BODY_COLOR      []string      `yaml:"BODY_COLOR"`
	QQ_EMAIL        string        `yaml:"QQ_EMAIL"`
	QQ_AUTH_PWD     string        `yaml:"QQ_AUTH_PWD"`
	AUTO_MAKE_BABY  string        `yaml:"AUTO_MAKE_BABY"`
	MAKE_BABY_PETID string        `yaml:"MAKE_BABY_PETID"`

	CHUANSHUO_SWITCH         int     `yaml:"CHUANSHUO_SWITCH"`         //传说总开关
	CHUANSHUO0_8DOG_0_PRICE  float32 `yaml:"CHUANSHUO0_8DOG_0_PRICE"`  //传说00价格
	CHUANSHUO0_8DOG_24_PRICE float32 `yaml:"CHUANSHUO0_8DOG_24_PRICE"` //传说024价格
	CHUANSHUO0_8DOG_2_PRICE  float32 `yaml:"CHUANSHUO0_8DOG_2_PRICE"`  //传说02价格
	CHUANSHUO_8DOG_OLD_PRICE float32 `yaml:"CHUANSHUO_8DOG_OLD_PRICE"` //大于0代传说价格

	GOD_SWITCH             int     `yaml:"GOD_SWITCH"`
	GOD_6DOG_BABY_PRICE    float32 `yaml:"GOD_6DOG_BABY_PRICE"`
	GOD_6DOG_SWITCH        int     `yaml:"GOD_6DOG_SWITCH"`    //6稀有神话开关
	GOD0_6DOG_0_PRICE      float32 `yaml:"GOD0_6DOG_0_PRICE"`  //0代神话0分钟价格
	GOD0_6DOG_24_PRICE     float32 `yaml:"GOD0_6DOG_24_PRICE"` //0代神话24小时价格
	GOD0_6DOG_2_PRICE      float32 `yaml:"GOD0_6DOG_2_PRICE"`  //0代神话2天价格
	GOD0_6_0SPECIAL_PRICE  float32 `yaml:"GOD0_6_0SPECIAL_PRICE"`
	GOD0_6_24SPECIAL_PRICE float32 `yaml:"GOD0_6_24SPECIAL_PRICE"`
	GOD0_6_2SPECIAL_PRICE  float32 `yaml:"GOD0_6_2SPECIAL_PRICE"`

	GOD1_6DOG_0_PRICE           float32 `yaml:"GOD1_6DOG_0_PRICE"` //1代神话0天价格
	GOD1_6DOG_2_PRICE           float32 `yaml:"GOD1_6DOG_2_PRICE"` //1代神话2天价格
	GOD1_6DOG_4_PRICE           float32 `yaml:"GOD1_6DOG_4_PRICE"` //1代神话4天价格
	GOD1_6_0SPECIAL_PRICE       float32 `yaml:"GOD1_6_0SPECIAL_PRICE"`
	GOD1_6_2SPECIAL_PRICE       float32 `yaml:"GOD1_6_2SPECIAL_PRICE"`
	GOD1_6_4SPECIAL_PRICE       float32 `yaml:"GOD1_6_4SPECIAL_PRICE"`
	GOD_6DOG_OLD1_PRICE         float32 `yaml:"GOD_6DOG_OLD1_PRICE"`
	GOD_6DOG_OLD1_SPECIAL_PRICE float32 `yaml:"GOD_6DOG_OLD1_SPECIAL_PRICE"`

	GOD_7DOG_SWITCH        int     `yaml:"GOD_7DOG_SWITCH"` //7稀有神话开关
	GOD_7DOG_BABY_PRICE    float32 `yaml:"GOD_7DOG_BABY_PRICE"`
	GOD0_7DOG_0_PRICE      float32 `yaml:"GOD0_7DOG_0_PRICE"`  //0代七稀神话0分钟
	GOD0_7DOG_24_PRICE     float32 `yaml:"GOD0_7DOG_24_PRICE"` //0代七稀神话24小时
	GOD0_7DOG_2_PRICE      float32 `yaml:"GOD0_7DOG_2_PRICE"`  //0代七稀神话2天
	GOD0_7_0SPECIAL_PRICE  float32 `yaml:"GOD0_7_0SPECIAL_PRICE"`
	GOD0_7_24SPECIAL_PRICE float32 `yaml:"GOD0_7_24SPECIAL_PRICE"`
	GOD0_7_2SPECIAL_PRICE  float32 `yaml:"GOD0_7_2SPECIAL_PRICE"`

	GOD_7DOG_OLD_PRICE         float32 `yaml:"GOD_7DOG_OLD_PRICE"`
	GOD_7DOG_OLD_SPECIAL_PRICE float32 `yaml:"GOD_7DOG_OLD_SPECIAL_PRICE"`

	SHISHI_SWITCH                 int     `yaml:"SHISHI_SWITCH"`
	SHISHI_5_SWITCH               int     `yaml:"SHISHI_5_SWITCH"`
	SHISHI0_5DOG_0_PRICE          float32 `yaml:"SHISHI0_5DOG_0_PRICE"`  //0代五稀史诗0天
	SHISHI0_5DOG_24_PRICE         float32 `yaml:"SHISHI0_5DOG_24_PRICE"` //0代五稀史诗24
	SHISHI_5BIRTHDAY_PRICE        float32 `yaml:"SHISHI_5BIRTHDAY_PRICE"`
	SHISHI0_5_0SPECIAL_PRICE      float32 `yaml:"SHISHI0_5_0SPECIAL_PRICE"`
	SHISHI0_5_24SPECIAL_PRICE     float32 `yaml:"SHISHI0_5_24SPECIAL_PRICE"`
	SHISHI0_5_2SPECIAL_PRICE      float32 `yaml:"SHISHI0_5_2SPECIAL_PRICE"`
	SHISHI_5DOG_OLD_PRICE         float32 `yaml:"SHISHI_5DOG_OLD_PRICE"`
	SHISHI_5DOG_OLD_SPECIAL_PRICE float32 `yaml:"SHISHI_5DOG_OLD_SPECIAL_PRICE"`

	SHISHI_4_SWITCH               int     `yaml:"SHISHI0_4_SWITCH"`
	SHISHI0_4DOG_0_PRICE          float32 `yaml:"SHISHI0_4DOG_0_PRICE"`
	SHISHI0_4DOG_24_PRICE         float32 `yaml:"SHISHI0_4DOG_24_PRICE"`
	SHISHI_4BIRTHDAY_PRICE        float32 `yaml:"SHISHI_4BIRTHDAY_PRICE"`
	SHISHI0_4_0SPECIAL_PRICE      float32 `yaml:"SHISHI0_4_0SPECIAL_PRICE"`
	SHISHI0_4_24SPECIAL_PRICE     float32 `yaml:"SHISHI0_4_24SPECIAL_PRICE"`
	SHISHI0_4_2SPECIAL_PRICE      float32 `yaml:"SHISHI0_4_2SPECIAL_PRICE"`
	SHISHI_4DOG_OLD_PRICE         float32 `yaml:"SHISHI_4DOG_OLD_PRICE"`
	SHISHI_4DOG_OLD_SPECIAL_PRICE float32 `yaml:"SHISHI_4DOG_OLD_SPECIAL_PRICE"`

	ZHUOYUE_SWITCH            int     `yaml:"ZHUOYUE_SWITCH"`
	ZHUEYUE0_2DOG_0_PRICE     float32 `yaml:"ZHUEYUE0_2DOG_0_PRICE"` //0,0卓越
	ZHUEYUE_BIRTHDAY_PRICE    float32 `yaml:"ZHUEYUE_BIRTHDAY_PRICE"`
	ZHUEYUE_GOOD_NUMBER_PRICE float32 `yaml:"ZHUEYUE_GOOD_NUMBER_PRICE"`
	ZHUOYUE0_0SPECIAL_PRICE   float32 `yaml:"ZHUOYUE0_0SPECIAL_PRICE"`
	ZHUOYUE_OLDER0_PRICE      float32 `yaml:"ZHUOYUE_OLDER0_PRICE"`

	XIYOU_SWITCH            int     `yaml:"XIYOU_SWITCH"`
	XIYOU0_1DOG_0_PRICE     float32 `yaml:"XIYOU0_1DOG_0_PRICE"`  //00 稀有
	XIYOU_BIRTHDAY_PRICE    float32 `yaml:"XIYOU_BIRTHDAY_PRICE"` //00 稀有
	XIYOU_GOOD_NUMBER_PRICE float32 `yaml:"XIYOU_GOOD_NUMBER_PRICE"`
	XIYOU0_0SPECIAL_PRICE   float32 `yaml:"XIYOU0_0SPECIAL_PRICE"`
	XIYOU_OLDER0_DOG_PRICE  float32 `yaml:"XIYOU_OLDER0_DOG_PRICE"`

	PUTONG_SWITCH            int     `yaml:"PUTONG_SWITCH"`
	PUTONG0_0DOG_0_PRICE     float32 `yaml:"PUTONG0_0DOG_0_PRICE"` //00 普通
	PUTONG_BIRTHDAY_PRICE    float32 `yaml:"PUTONG_BIRTHDAY_PRICE"`
	PUTONG_GOOD_NUMBER_PRICE float32 `yaml:"PUTONG_GOOD_NUMBER_PRICE"`
	PUTONG_OLDER_DOG_PRICE   float32 `yaml:"PUTONG_OLDER_DOG_PRICE"`
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
		fmt.Print("配置文件不合法，请检查配置文件内容", err, "\n")
		fmt.Print("\n")
	}

	return configuration
}
