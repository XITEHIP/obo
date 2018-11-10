package group

import (
	"github.com/xitehip/obo/define"
	"github.com/xitehip/obo/support"
)

type CustomerService struct {
}

func (o *CustomerService) Register(session *define.Session) {
	session.PluginsManager.Handles["customer_service"] = define.Handle(service)
}

func service(session *define.Session, respMsg *define.ReceiveMessage) func(fun define.TransmitFun) {
	if respMsg.MsgFrom == define.MSG_FROM_GROUP {
		support.Cl().Message(respMsg.Content)
	}
	return nil
}
