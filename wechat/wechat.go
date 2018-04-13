package wechat
import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"  // 导入协议包
)
// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	// 将插件注册到session
	// 第一个参数: 指定消息类型, 所有该类型的消息都会被转发到此插件
	// 第二个参数: 指定消息处理函数, 消息会进入此函数
	// 第三个参数: 自定义插件名，不能重名，switcher插件会用到此名称
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(buy_msg), "dogwechat")

	// 开启插件
	if err := session.HandlerRegister.EnableByName("dogwechat"); err != nil {
		logs.Error(err)
	}
}

// 消息处理函数
func buy_msg(session *wxweb.Session, msg *wxweb.ReceivedMessage) {

	// 可选: 可以用contact manager来过滤, 比如过滤掉没有保存到通讯录的群
	// 注意，contact manager只存储了已保存到通讯录的群组
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("ignore the messages from", msg.FromUserName)
		return
	}

	// 可选: 过滤和自己无关的群组消息
	if msg.IsGroup && msg.Who != session.Bot.UserName {
		return
	}

	// 取出收到的内容
	// 取text
	logs.Info(msg.Content)
	// anything
	// 回复消息
	// 第一个参数: 回复的内容
	// 第二个参数: 机器人ID
	// 第三个参数: 联系人/群组/特殊账号ID
	session.SendText("plugin demo", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
	// 回复图片和gif 参见wxweb/session.go

}
