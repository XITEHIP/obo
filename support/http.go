package support

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
	"errors"
	"fmt"
)

type HttpClient struct {
	client    *http.Client
	userAgent string
}

var httpObj *HttpClient

var err error

func GetHttp() *HttpClient {
	if httpObj == nil {

		var netTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		}
		cookieJar, _ := cookiejar.New(nil)
		httpClient := &http.Client{
			Timeout:   time.Second * 100,
			Transport: netTransport,
			Jar:       cookieJar,
		}
		httpObj = &HttpClient{
			client:    httpClient,
			userAgent: "ApiV2 Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36 ",
		}
	}
	return httpObj
}

//Get body byte
func (o *HttpClient) GetBodyByte(url string, query map[string]string) []byte {

	defer func() {
		if r := recover(); r != nil {
			Cl().Error(fmt.Sprintf("GetBodyByte:%s", r))
		}
	}()

	resp := o.Get(url, query)
	if resp == nil {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		checkErr(err)
	}
	return body
}

//Get body byte
func (o *HttpClient) GetBodyMap(url string, query map[string]string) map[string]interface{} {
	resp := o.Get(url, query)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		checkErr(err)
	}
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)

	return result
}

// Get body string
func (o *HttpClient) GetBodyStr(url string, query map[string]string) string {
	return string(o.GetBodyByte(url, query))
}

//Get Resp
func (o *HttpClient) Get(uri string, query map[string]string) *http.Response {

	if query != nil && len(query) > 0 {
		uv := url.Values{}
		for key, val := range query {
			uv.Add(key, val)
		}
		if strings.Contains(uri, "?") {
			uri += "&" + uv.Encode()
		} else {
			uri += "?" + uv.Encode()
		}
	}
	return o.httpDo("GET", uri, nil)
}

//Post resp byte
func (o *HttpClient) Post(url string, params map[string]interface{}) []byte {
	resp := o.httpDo("POST", url, params)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		checkErr(err)
	}
	return body
}

//Post json resp byte
func (o *HttpClient) PostJson(url string, params map[string]interface{}) map[string]interface{} {
	resp := o.httpDo("POST", url, params)
	result := make(map[string]interface{})
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		checkErr(err)
	}
	json.Unmarshal(body, &result)
	return result
}

func (o *HttpClient) PostJsonResp(url string, params map[string]interface{}) *http.Response {
	resp := o.httpDo("POST", url, params)
	return resp
}

//base http
func (o *HttpClient) httpDo(method string, url string, params map[string]interface{}) *http.Response {

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknow panic")
			}
			Cl().Error(err.Error())
		}
	}()

	if params == nil {
		params = make(map[string]interface{})
	}
	p, err := json.Marshal(params)
	if err != nil {
		checkErr(err)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(p))
	if err != nil {
		checkErr(err)
	}
	//req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		checkErr(err)
	}
	resp, err := o.client.Do(req)
	if err != nil {
		checkErr(err)
	}

	return resp
}

func (o *HttpClient) Upload(uri string, body io.Reader, headers map[string]string) []byte {

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		checkErr(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := o.client.Do(req)
	if err != nil {
		checkErr(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return respBody
}

func (o *HttpClient) SetJar(jar http.CookieJar) {
	o.client.Jar = jar
}

func checkErr(err error) {
	if err != nil {
		Cl().Error("http error:" + err.Error())
	}
}


func(o *HttpClient) PostForm(u string, values url.Values) map[string]interface{} {
	resp, err := http.PostForm(u, values)

	if err != nil {
		checkErr(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		checkErr(err)
	}
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)

	return result
}

