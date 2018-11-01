package obo


import (
	"net/http"
	"github.com/xitehip/obo/api"
	"github.com/xitehip/obo/define"
)

func SendText(s *define.Session, msg, from, to string) {
	api.SendMsg(s.Bc.Lpr, s.Bc.Br, msg, from, to, s.Cookies)
}

func SendImg(s *define.Session, filename string, from, to string,
	lpr *define.LoginPageResp, br *define.BaseRequest, cookies []*http.Cookie) (map[string]interface{}, error) {

	mediaId, err := api.UploadMedia(filename, from, to, s.Bc.Lpr, s.Bc.Br, cookies)
	if err != nil {
		return nil, err
	}
	return api.SendImg(lpr, br, mediaId, from, to, cookies), nil
}

