package api

import (
	"github.com/xitehip/obo/define"
	"github.com/xitehip/obo/support"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"
)

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
