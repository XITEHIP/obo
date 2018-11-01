package define

import (
	"encoding/xml"
	"net/http"
	"sync"
)

type LoginConfig struct {
	Uuid     string
	Redirect string
}

type LoginPageResp struct {
	XMLName     xml.Name `xml:"error"`
	Ret         int      `xml:"ret"`
	Message     string   `xml:"message"`
	Skey        string   `xml:"skey"`
	Wxsid       string   `xml:"wxsid"`
	Wxuin       string   `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	IsGrayscale int      `xml:"isgrayscale"`
}

type BotClient struct {
	//DeviceID    string
	UserAgent string
}

type BaseRequest struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
}

type SyncKey struct {
	Key int
	Val int
}

type SyncKeyList struct {
	Count int
	List  []SyncKey
}

type Myself struct {
	AppAccountFlag    int
	ContactFlag       int
	HeadImgFlag       int
	HeadImgUrl        string
	HideInputBarFlag  int
	NickName          string
	PYInitial         string
	PYQuanPin         string
	RemarkName        string
	RemarkPYInitial   string
	RemarkPYQuanPin   string
	Sex               int
	Signature         string
	SnsFlag           int
	StarFriend        int
	Uin               int
	UserName          string
	VerifyFlag        int
	WebWxPluginSwitch int
}

type Contact struct {
	Alias            string
	AppAccountFlag   int
	AttrStatus       int
	ChatRoomId       int
	City             string
	ContactFlag      int
	DisplayName      string
	EncryChatRoomId  string
	HeadImgUrl       string
	HideInputBarFlag int
	IsOwner          int
	KeyWord          string
	MemberCount      int
	MemberList       []*Contact
	NickName         string
	OwnerUin         int
	PYInitial        string
	PYQuanPin        string
	Province         string
	RemarkName       string
	RemarkPYInitial  string
	RemarkPYQuanPin  string
	Sex              int
	Signature        string
	SnsFlag          int
	StarFriend       int
	Statues          int
	Uin              int
	UniFriend        int
	UserName         string
	VerifyFlag       int
}

type ContactWrap struct {
	List  map[string]*Contact
	Count int
}

type BotConfig struct {
	Lc  *LoginConfig
	Lpr *LoginPageResp
	Bc  *BotClient
	Br  *BaseRequest
}

//微信专用
type Specials struct {
	*ContactWrap
}

//官方
type Officials struct {
	*ContactWrap
}

//群
type Groups struct {
	*ContactWrap
}

//通讯录
type Friends struct {
	*ContactWrap
}

type TextMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
}

type MediaMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
	MediaId      string
}

type ReceiveMessage struct {
	MsgId        string
	MsgType      int
	MsgFrom      int
	FromUserName string
	ToUserName   string
	Content      string
	OriContent   string
	Url          string
}

type Handle func(*Session, *ReceiveMessage)

type PluginsManager struct {
	Handles map[string]Handle
}

type Session struct {
	Cookies   []*http.Cookie
	Bc        *BotConfig
	Skl       *SyncKeyList
	Myself    *Myself
	Specials  *Specials
	Officials *Officials
	Groups    *Groups
	Friends   *Friends
	*PluginsManager
	MuCookie  sync.Mutex
}
