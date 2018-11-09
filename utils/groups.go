package utils

import (
	"github.com/xitehip/obo/define"
)

func GetUserNameByNickName(s *define.Session, nickName string) string  {
	for _, group := range s.Groups.List {
		if group.NickName == nickName {
			return group.UserName
		}
	}
	return ""
}
