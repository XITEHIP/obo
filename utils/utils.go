package utils

import (
	"github.com/xitehip/obo/define"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetRandomStringFromNum(length int) string {
	bytes := []byte("0123456789")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func ParsingAddMsgList(data []interface{}, callback func(*define.ReceiveMessage)) {

	if len(data) > 0 {
		for _, item := range data {
			tmp := item.(map[string]interface{})
			msg := &define.ReceiveMessage{}
			msg.MsgId = tmp["MsgId"].(string)
			msg.MsgType = int(tmp["MsgType"].(float64))
			msg.FromUserName = tmp["FromUserName"].(string)
			msg.ToUserName = tmp["ToUserName"].(string)
			msg.Content = tmp["Content"].(string)
			msg.OriContent = tmp["OriContent"].(string)
			msg.Url = tmp["Url"].(string)

			if strings.Contains(msg.FromUserName, "@@") || strings.Contains(msg.ToUserName, "@@") {
				msg.MsgFrom = define.MSG_FROM_GROUP
			} else if msg.ToUserName == "filehelper" {
				msg.MsgFrom = define.MSG_FROM_FILEHELPER
			}
			callback(msg)
		}
	}
}

func SyncKeyStr(sks []define.SyncKey) string {
	sksSlice := make([]string, 0)
	for _, v := range sks {
		sksSlice = append(sksSlice, strconv.Itoa(v.Key) + "_" + strconv.Itoa(v.Val))
	}
	return strings.Join(sksSlice, "|")
}

func GetMyself(myOrig map[string]interface{}) *define.Myself {
	myself := &define.Myself{}
	fields := reflect.ValueOf(myself).Elem()
	for k, v := range myOrig {
		field := fields.FieldByName(k)
		if ftv, ok := v.(float64); ok {
			field.Set(reflect.ValueOf(int(ftv)))
		} else {
			field.Set(reflect.ValueOf(v))
		}
	}
	return myself
}

func SetCookies(s *define.Session, cookies []*http.Cookie) {
	s.MuCookie.Lock()
	defer s.MuCookie.Unlock()

	s.Cookies = cookies
}

func GetCookies(s *define.Session) []*http.Cookie {
	s.MuCookie.Lock()
	defer s.MuCookie.Unlock()
	return s.Cookies
}
