package group

import (
	"fmt"
	"github.com/xitehip/obo/define"
)


type CustomerService struct {
}

func (o *CustomerService) Register(session *define.Session) {
	session.PluginsManager.Handles["customer_service"] = define.Handle(service)
}

func service(session *define.Session, respMsg *define.ReceiveMessage) {
	if respMsg.MsgFrom == define.MSG_FROM_GROUP {
		fmt.Println(respMsg.Content)
	}
}


