package api

import (
	"encoding/xml"
	"fmt"
	"github.com/mdp/qrterminal"
	"github.com/xitehip/obo/define"
	"github.com/xitehip/obo/support"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getUuid() string {
	query := make(map[string]string, 0)
	query["appid"] = "wx782c26e4c19acffb"
	query["fun"] = "new"
	query["lang"] = "zh_CN"
	query["_"] = strconv.FormatInt(time.Now().Unix(), 10)
	response := support.GetHttp().GetBodyStr("https://login.weixin.qq.com/jslogin", query)
	rs := strings.Split(response, "\"")
	return rs[1]
}

func ShowQr(lc *define.LoginConfig) {
	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}
	lc.Uuid = getUuid()
	fmt.Println("https://login.weixin.qq.com/l/"+lc.Uuid)
	qrterminal.GenerateWithConfig("https://login.weixin.qq.com/l/"+lc.Uuid, config)
}

func ListenScan(tip int64, lc *define.LoginConfig) (string, string) {
	uv := url.Values{}
	uv.Add("loginicon", "true")
	uv.Add("uuid", lc.Uuid)
	uv.Add("tip", strconv.FormatInt(tip, 10))
	uv.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	uv.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	url := SERVER_URI_BASE + "login?" + uv.Encode()
	response := support.GetHttp().GetBodyStr(url, nil)
	reg := regexp.MustCompile(`(\d){3}`)
	matches := reg.FindAllString(response, -1)

	return matches[0], response
}

func WebWxNewLoginPage(lc *define.LoginConfig, lpr *define.LoginPageResp) ([]*http.Cookie, error) {
	u, _ := url.Parse(lc.Redirect)
	km := u.Query()
	km.Add("fun", "new")
	uri := SERVER_URI_BASE + "webwxnewloginpage?" + km.Encode()

	resp := support.GetHttp().Get(uri, nil)
	body, _ := ioutil.ReadAll(resp.Body)
	if err := xml.Unmarshal(body, lpr); err != nil {
		return nil, err
	}
	if lpr.Ret != 0 {
		return nil, fmt.Errorf("xc.Ret != 0: %s", string(body))
	}

	return resp.Cookies(), nil
}
