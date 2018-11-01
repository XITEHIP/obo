//test  auto send a img

package filehelper

import (
	"fmt"
	"github.com/xitehip/obo"
	"github.com/xitehip/obo/define"
	"github.com/xitehip/obo/utils"
)

type AutoSendService struct {
}

func (o *AutoSendService) Register(session *define.Session) {
	session.PluginsManager.Handles["autosend_service"] = define.Handle(service)
}

func service(session *define.Session, respMsg *define.ReceiveMessage) {
	if respMsg.MsgFrom == define.MSG_FROM_FILEHELPER {
		resp, _ := obo.SendImg(session, "/Users/xitehip/Desktop/zhye.png", respMsg.FromUserName, respMsg.ToUserName, session.Bc.Lpr, session.Bc.Br, utils.GetCookies(session))
		fmt.Println(resp)
	}
}
