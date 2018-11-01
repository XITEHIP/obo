package api

import (
	"encoding/json"
	"github.com/xitehip/obo/support"
	"github.com/xitehip/obo/support/tools"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
	"github.com/xitehip/obo/define"
)

var SPECIAL_USERS = []string{
	"newsapp", "fmessage", "filehelper", "weibo", "qqmail",
	"fmessage", "tmessage", "qmessage", "qqsync", "floatbottle",
	"lbsapp", "shakeapp", "medianote", "qqfriend", "readerapp",
	"blogapp", "facebookapp", "masssendapp", "meishiapp",
	"feedsapp", "voip", "blogappweixin", "weixin", "brandsessionholder",
	"weixinreminder", "wxid_novlwrv3lqwv11", "gh_22b87fa7cb3c",
	"officialaccounts", "notification_messages", "wxid_novlwrv3lqwv11",
	"gh_22b87fa7cb3c", "wxitil", "userexperience_alarm", "notification_messages",
}

var (
	specials  *define.Specials
	officials *define.Officials
	groups    *define.Groups
	friends   *define.Friends
)

func InitContactList(contactsOrig []interface{}) []interface{} {

	if specials == nil {
		specials = &define.Specials{&define.ContactWrap{List: make(map[string]*define.Contact)}}
	}
	if officials == nil {
		officials = &define.Officials{&define.ContactWrap{List: make(map[string]*define.Contact)}}
	}
	if groups == nil {
		groups = &define.Groups{&define.ContactWrap{List: make(map[string]*define.Contact)}}
	}
	if friends == nil {
		friends = &define.Friends{&define.ContactWrap{List: make(map[string]*define.Contact)}}
	}
	s := false
	o := false
	g := false
	f := false
	for _, contact := range contactsOrig {
		contactObj := &define.Contact{}
		contactByte, _ := json.Marshal(contact)
		json.Unmarshal(contactByte, contactObj)
		userName := contactObj.UserName
		if tools.Find(userName, SPECIAL_USERS) {
			if s == false {
				specials.Count = 0
				s = true
			}
			specials.List[userName] = contactObj
			specials.Count++
		} else if (contactObj.VerifyFlag & 8) != 0 {
			if o == false {
				officials.Count = 0
				o = true
			}
			officials.List[userName] = contactObj
			officials.Count++
		} else if strings.Contains(userName, "@@") {
			if g == false {
				groups.Count = 0
				g = true
			}
			groups.List[userName] = contactObj
			groups.Count++
		} else {
			if f == false {
				friends.Count = 0
				f = true
			}
			friends.List[userName] = contactObj
			friends.Count++
		}
	}
	return []interface{}{specials, officials, groups, friends}
}

func WebWxGetContact(lpr *define.LoginPageResp, cookies []*http.Cookie) []byte {
	km := url.Values{}
	km.Add("lang", "zh_CN")
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("seq", "0")
	km.Add("skey", lpr.Skey)
	km.Add("pass_ticket", lpr.PassTicket)
	urlStr := SERVER_URI_BASE + "webwxgetcontact?" + km.Encode()

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(urlStr)
	jar.SetCookies(u, cookies)
	support.GetHttp().SetJar(jar)

	result := support.GetHttp().GetBodyByte(urlStr, nil)

	return result
}

func WebWxBatchGetContact(br *define.BaseRequest, cookies []*http.Cookie, gr *define.Groups) map[string]interface{} {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("type", "ex")
	uri := SERVER_URI_BASE + "webwxbatchgetcontact?" + km.Encode()

	params := make(map[string]interface{})
	params["BaseRequest"] = br
	params["Count"] = len(gr.List)

	plist := make([]map[string]string, 0)
	for _, v := range gr.List {
		one := make(map[string]string)
		one["EncryChatRoomId"] = ""
		one["UserName"] = v.UserName
		plist = append(plist, one)
	}
	params["List"] = plist

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	support.GetHttp().SetJar(jar)

	result := support.GetHttp().PostJson(uri, params)
	return result
}
