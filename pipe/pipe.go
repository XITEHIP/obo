package pipe

import (
	"encoding/json"
	"fmt"
	"github.com/xitehip/obo/api"
	"github.com/xitehip/obo/define"
	"github.com/xitehip/obo/plugins"
	"github.com/xitehip/obo/support/tools"
	"github.com/xitehip/obo/utils"
	"strings"
	"github.com/xitehip/obo/support"
	"log"
)

const (
	PIPE_NODE_ShowQr = iota
	PIPE_NODE_Login
	PIPE_NODE_WebWxNewLoginPage
	PIPE_NODE_WebWxInit
	PIPE_NODE_Listen
	PIPE_NODE_Customer
	PIPE_NODE_Exit
)

var isInit = false

type Pipe struct {
	flowChan chan int
	Traces   []int
	session  *define.Session
	errChan  chan error
	msgChan  chan []byte
	receiveMsgChan chan string
	transmit define.TransmitFun
}

func (o *Pipe)Session() *define.Session  {
	return o.session
}

func initCfg() *define.BotConfig {
	bc := &define.BotConfig{}
	bc.Lc = &define.LoginConfig{}
	bc.Lpr = &define.LoginPageResp{}
	bc.Bc = &define.BotClient{}
	bc.Br = &define.BaseRequest{}
	return bc
}

func New() *Pipe {
	o := &Pipe{}
	o.flowChan = make(chan int, 1)
	o.errChan = make(chan error)
	o.msgChan = make(chan []byte, 1024)
	o.session = &define.Session{}
	o.session.PluginsManager = &define.PluginsManager{Handles: make(map[string]define.Handle)}
	o.receiveMsgChan = make(chan string)

	return o
}

func (o *Pipe) AttachPlugins(plugins []plugins.PluginProviderInterface) *Pipe {
	if len(plugins) > 0 {
		for _, plugin := range plugins {
			plugin.Register(o.session)
		}
	}

	return o
}

//绑定转发消息
func (o *Pipe) AttachTransmit(f define.TransmitFun) *Pipe {

	if f != nil {
		o.transmit = f
	}
	return o
}


func (o *Pipe) Run() {

	check()
	bc := initCfg()

	o.flowChan <- PIPE_NODE_ShowQr

	for {
		switch <-o.flowChan {
		case PIPE_NODE_ShowQr:
			isInit = false
			api.ShowQr(bc.Lc)
			o.flowChan <- PIPE_NODE_Login
		case PIPE_NODE_Login:
			o.login(bc.Lc, o.flowChan)
		case PIPE_NODE_WebWxNewLoginPage:
			o.webWxNewLoginPage(bc)
		case PIPE_NODE_WebWxInit:
			o.webWxInit()
		case PIPE_NODE_Listen:
			isInit = true
			//go o.receiveListen()
			go o.listen()
		case PIPE_NODE_Customer:
			o.customer()
		case PIPE_NODE_Exit:
			break
		}
	}
	log.Fatal("obo is exit!")
}

func (o *Pipe)IsInited() bool {
	return isInit
}

func (o *Pipe) login(lc *define.LoginConfig, ch chan int) {
	tip := int64(1)
	support.Cl().Message("Please scan the qrCode with wechat.")
	for i := 0; i < 10; i++ {
		code, response := api.ListenScan(tip, lc)
		switch code {
		case "201":
			support.Cl().Message("Please confirm login in wechat.")
			tip = 0
		case "200":
			rs := strings.Split(response, "\"")
			lc.Redirect = rs[1] + "&fun=new"
			o.flowChan <- PIPE_NODE_WebWxNewLoginPage
			o.loginAfter()
			return
		case "408":
			tip = 1
			support.Cl().Message("Login timeout, response code 408.")
		default:
			tip = 1
			support.Cl().Message("Login error " + code)
		}
	}
	ch <- PIPE_NODE_Exit
}

func (o *Pipe) listen() {
	o.flowChan <- PIPE_NODE_Customer
	support.Cl().Message("obo begin listen...")
	quitCurrClientCode := []string{"1100", "1101", "1102", "1205"}

	for {
		ret, sel := api.SyncCheck(o.session.Bc.Lpr, o.session.Bc.Br, o.session.Skl, utils.GetCookies(o.session))
		fmt.Printf("ret:" + ret + " " + "sel:" + sel + "\n")
		if tools.FindArr(ret, quitCurrClientCode) {
			o.errChan <- fmt.Errorf("api blocked, ret:%s", ret)
			break
		}
		if sel == "0" {
			continue
		}
		cookies := api.WebWxSync(o.session.Bc.Lpr, o.session.Bc.Br, o.session.Skl, utils.GetCookies(o.session), o.msgChan)
		if cookies == nil {
			continue
		}
		utils.SetCookies(o.session, cookies)
	}
}

func (o *Pipe) customer() {
	for {
		select {
		case err := <-o.errChan:
			support.Cl().Message(err.Error())
			//try reload login
			o.flowChan <- PIPE_NODE_Login
			return
		case msg := <-o.msgChan:
			msgMap := make(map[string]interface{})
			json.Unmarshal(msg, &msgMap)
			if int(msgMap["AddMsgCount"].(float64)) > 0 {
				utils.ParsingAddMsgList(msgMap["AddMsgList"].([]interface{}), func(message *define.ReceiveMessage) {
					plugins.Fire(o.session, message, o.transmit)
				})
			}
		}
	}
}

func (o *Pipe) webWxNewLoginPage(bc *define.BotConfig) {
	cookies, _ := api.WebWxNewLoginPage(bc.Lc, bc.Lpr)
	utils.SetCookies(o.session, cookies)
	o.session.Bc = bc
	o.flowChan <- PIPE_NODE_WebWxInit
}

func (o *Pipe) webWxInit() {
	support.Cl().Message("wx init begin.")
	resp := api.WebWxInit(o.session.Bc.Lpr, o.session.Bc.Br)
	utils.GenerateSyncKey(o.session, resp)
	go o.myself(resp)
	go o.initContactList("ContactList", resp)
	o.initFriends()
	o.initGroups()
	o.flowChan <- PIPE_NODE_Listen
	support.Cl().Message("wx init end.")
}

func (o *Pipe) initContactList(key string, result map[string]interface{}) {
	cl := result[key].([]interface{})
	c := utils.InitContactList(cl)
	o.session.Specials = c[0].(*define.Specials)
	o.session.Officials = c[1].(*define.Officials)
	o.session.Groups = c[2].(*define.Groups)
	o.session.Friends = c[3].(*define.Friends)
}

func (o *Pipe) initFriends() {
	friendsByte := api.WebWxGetContact(o.session.Bc.Lpr, utils.GetCookies(o.session))
	contactsMap := make(map[string]interface{})
	json.Unmarshal(friendsByte, &contactsMap)
	o.initContactList("MemberList", contactsMap)
}

func (o *Pipe) initGroups() {
	resp := api.WebWxBatchGetContact(o.session.Bc.Br, utils.GetCookies(o.session), o.session.Groups)
	o.initContactList("ContactList", resp)
}

func (o *Pipe) myself(result map[string]interface{}) {
	user := result["User"].(map[string]interface{})
	o.session.Myself = utils.GetMyself(user)
}

func (o *Pipe) webWxStatusNotify() {
	api.WebWxStatusNotify(o.session.Bc.Lpr, o.session.Bc.Br, o.session.Myself)
}

func check() {
}

func (o *Pipe) loginAfter() {
	support.Cl().Message("Login success!")
}

//其他系统主动发送给obo text msg
func (o *Pipe)receiveListen()  {
	support.TcpServer(o.receiveMsgChan)
	for msg := range o.receiveMsgChan {
		api.SendMsg(o.session.Bc.Lpr, o.session.Bc.Br, msg, o.session.Myself.UserName, "filehelper", o.session.Cookies)
	}
}

func (o *Pipe)SendMsg(msg string, to string )  (map[string]interface{}, error)  {
	to = utils.GetUserNameByNickName(o.Session(), to)
	return api.SendMsg(o.session.Bc.Lpr, o.session.Bc.Br, msg, o.session.Myself.UserName, to, o.session.Cookies), nil
}

func (o *Pipe)SendImg(filename string, to string) (map[string]interface{}, error) {
	to = utils.GetUserNameByNickName(o.Session(), to)
	mediaId, err := api.UploadMedia(filename, o.session.Myself.UserName, to, o.session.Bc.Lpr, o.session.Bc.Br, o.session.Cookies)
	if err != nil {
		support.Cl().Error(err.Error())
		return nil, err
	}
	return api.SendImg(o.session.Bc.Lpr, o.session.Bc.Br, mediaId, o.session.Myself.UserName, to, o.session.Cookies), nil
}
