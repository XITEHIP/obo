package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xitehip/obo/support"
	"github.com/xitehip/obo/support/tools"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/xitehip/obo/define"
)



func SendMsg(lpr *define.LoginPageResp, br *define.BaseRequest, msg string,
	from, to string, cookies []*http.Cookie) map[string]interface{} {
	km := url.Values{}
	km.Add("pass_ticket", lpr.PassTicket)

	uri := SERVER_URI_BASE + "webwxsendmsg?" + km.Encode()
	params := make(map[string]interface{})
	params["BaseRequest"] = br
	params["Msg"] = &define.TextMessage{
		Type:         1,
		Content:      msg,
		FromUserName: from,
		ToUserName:   to,
		LocalID:      int(time.Now().Unix() * 1e4),
		ClientMsgId:  int(time.Now().Unix() * 1e4),
	}
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	support.GetHttp().SetJar(jar)
	body := support.GetHttp().PostJson(uri, params)

	return body
}

func SendImg(lpr *define.LoginPageResp, br *define.BaseRequest, mediaId string,
	from, to string, cookies []*http.Cookie) map[string]interface{} {

	km := url.Values{}
	km.Add("pass_ticket", lpr.PassTicket)
	km.Add("fun", "async")
	km.Add("f", "json")
	km.Add("lang", "zh_CN")
	uri := SERVER_URI_BASE + "webwxsendmsgimg?" + km.Encode()

	params := make(map[string]interface{})
	params["BaseRequest"] = br
	params["Msg"] = &define.MediaMessage{
		Type:         3,
		Content:      "",
		FromUserName: from,
		ToUserName:   to,
		LocalID:      int(time.Now().Unix() * 1e4),
		ClientMsgId:  int(time.Now().Unix() * 1e4),
		MediaId:      mediaId,
	}
	params["Scene"] = 0
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	support.GetHttp().SetJar(jar)
	body := support.GetHttp().PostJson(uri, params)

	return body
}

func UploadMedia(filename string, from, to string,
	lpr *define.LoginPageResp, br *define.BaseRequest, cookies []*http.Cookie) (string, error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("filename", filename)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, bytes.NewReader(content)); err != nil {
		return "", err
	}

	ss := strings.Split(filename, ".")
	suffix := ss[len(ss)-1]

	fw, _ = w.CreateFormField("id")
	fw.Write([]byte("WU_FILE_0"))

	fw, _ = w.CreateFormField("name")
	fw.Write([]byte(filepath.Base(filename)))

	fw, _ = w.CreateFormField("type")
	if suffix == "gif" {
		fw.Write([]byte("image/gif"))
	} else if suffix == "png" {
		fw.Write([]byte("image/png"))
	} else if suffix == "jpeg" {
		fw.Write([]byte("image/jpeg"))
	}

	fw, _ = w.CreateFormField("lastModifieDate")
	fw.Write([]byte("Mon Feb 13 2018 14:52:29 GMT+0800 (CST)"))

	fw, _ = w.CreateFormField("size")
	fw.Write([]byte(strconv.Itoa(len(content))))

	fw, _ = w.CreateFormField("mediatype")
	if suffix == "gif" {
		fw.Write([]byte("doc"))
	} else {
		fw.Write([]byte("pic"))
	}

	umr := make(map[string]interface{})
	umr["UploadType"] = 2
	umr["BaseRequest"] = br
	umr["ClientMediaId"] = int(time.Now().Unix() * 1e4)
	umr["TotalLen"] = len(content)
	umr["StartPos"] = 0
	umr["DataLen"] = len(content)
	umr["MediaType"] = 4
	umr["FromUserName"] = from
	umr["ToUserName"] = to
	umr["FileMd5"] = tools.Md5File(filename)

	jb, _ := json.Marshal(umr)
	fw, _ = w.CreateFormField("uploadmediarequest")
	fw.Write(jb)

	fw, _ = w.CreateFormField("webwx_data_ticket")
	for _, v := range cookies {
		if strings.Contains(v.String(), "webwx_data_ticket") {
			fw.Write([]byte(strings.Split(v.String(), "=")[1]))
			break
		}
	}
	fw, _ = w.CreateFormField("pass_ticket")
	fw.Write([]byte(lpr.PassTicket))

	w.Close()

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(SERVER_URI_FILE)
	jar.SetCookies(u, cookies)
	support.GetHttp().SetJar(jar)
	body := support.GetHttp().Upload(SERVER_URI_FILE + "webwxuploadmedia?f=json", &b, map[string]string{"Content-Type": w.FormDataContentType()})
	bodyMap := make(map[string]interface{})
	json.Unmarshal(body, &bodyMap)
	ret := int(bodyMap["BaseResponse"].(map[string]interface{})["Ret"].(float64))
	if ret != 0 {
		return "", fmt.Errorf("BaseResponse.Ret=%d", ret)
	}

	return bodyMap["MediaId"].(string), nil
}
