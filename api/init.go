package api

import (
	"encoding/json"
	"github.com/xitehip/obo/support"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"time"
	"github.com/xitehip/obo/define"
	"github.com/xitehip/obo/utils"
)

func WebWxInit(lpr *define.LoginPageResp, br *define.BaseRequest) map[string]interface{} {
	km := url.Values{}
	km.Add("pass_ticket", lpr.PassTicket)
	km.Add("skey", lpr.Skey)
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))

	url := SERVER_URI_BASE + "webwxinit?" + km.Encode()
	params := make(map[string]interface{})

	br.Uin = lpr.Wxuin
	br.Sid = lpr.Wxsid
	br.Skey = lpr.Skey
	br.DeviceID = "e" + utils.GetRandomStringFromNum(15)
	params["BaseRequest"] = br
	result := support.GetHttp().PostJson(url, params)

	return result
}

func WebWxStatusNotify(lpr *define.LoginPageResp, br *define.BaseRequest, my *define.Myself) map[string]interface{} {
	km := url.Values{}
	km.Add("pass_ticket", lpr.PassTicket)
	km.Add("lang", "zh_CN")

	url := SERVER_URI_BASE + "webwxstatusnotify?" + km.Encode()

	params := make(map[string]interface{})
	params["BaseRequest"] = br
	params["Code"] = 3
	params["ClientMsgId"] = int(time.Now().Unix())
	params["FromUserName"] = my.UserName
	params["ToUserName"] = my.UserName

	result := support.GetHttp().PostJson(url, params)

	return result
}

func SyncKey(syncKey map[string]interface{}) *define.SyncKeyList {
	skl := &define.SyncKeyList{}
	skl.Count = int(syncKey["Count"].(float64))
	list := syncKey["List"].([]interface{})
	sks := make([]define.SyncKey, 0)
	for _, val := range list {
		tmp := define.SyncKey{Key: int(val.(map[string]interface{})["Key"].(float64)), Val: int(val.(map[string]interface{})["Val"].(float64))}
		sks = append(sks, tmp)
	}
	skl.List = sks
	return skl
}

func SyncCheck(lpr *define.LoginPageResp, br *define.BaseRequest, skl *define.SyncKeyList, cookies []*http.Cookie) (string, string) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix()*1000, 10))
	km.Add("skey", lpr.Skey)
	km.Add("sid", lpr.Wxsid)
	km.Add("uin", lpr.Wxuin)
	km.Add("deviceid", br.DeviceID)
	km.Add("synckey", utils.SyncKeyStr(skl.List))
	km.Add("_", strconv.FormatInt(time.Now().Unix()*1000, 10))
	uri := SERVER_URI_WEBPUSH + "synccheck?" + km.Encode()
	params := make(map[string]interface{})
	params["BaseRequest"] = br
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	support.GetHttp().SetJar(jar)

	resp := support.GetHttp().GetBodyStr(uri, nil)

	reg := regexp.MustCompile("window.synccheck={retcode:\"(\\d+)\",selector:\"(\\d+)\"}")
	sub := reg.FindStringSubmatch(resp)
	retcode := "0"
	selector := "0"
	if len(sub) >= 2 {
		retcode = sub[1]
		selector = sub[2]
	}

	return retcode, selector
}

func WebWxSync(lpr *define.LoginPageResp, br *define.BaseRequest, skl *define.SyncKeyList, cookies []*http.Cookie, msg chan []byte) []*http.Cookie {
	km := url.Values{}
	km.Add("skey", lpr.Skey)
	km.Add("sid", lpr.Wxsid)
	km.Add("pass_ticket", lpr.PassTicket)
	uri := SERVER_URI_BASE + "webwxsync?" + km.Encode()
	params := make(map[string]interface{})
	params["BaseRequest"] = br
	params["SyncKey"] = skl
	params["rr"] = ^int(time.Now().Unix()) + 1

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	support.GetHttp().SetJar(jar)

	resp := support.GetHttp().PostJsonResp(uri, params)
	body, _ := ioutil.ReadAll(resp.Body)
	bodyMap := make(map[string]interface{})
	json.Unmarshal(body, &bodyMap)
	if bodyMap["BaseResponse"].(map[string]interface{})["Ret"] == float64(0){
		msg <- body
		skl.List = skl.List[:0]
		tmp := SyncKey(bodyMap["SyncKey"].(map[string]interface{}))
		skl.List = append(skl.List, tmp.List...)
		skl.Count = tmp.Count

		return resp.Cookies()
	}

	return nil
}
