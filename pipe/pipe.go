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
	"time"
)

const (
	PIPE_NODE_ShowQr = iota
	PIPE_NODE_Login
	PIPE_NODE_WebWxNewLoginpage
	PIPE_NODE_WebWxInit
	PIPE_NODE_Listen
	PIPE_NODE_Customer
	PIPE_NODE_Exit
)

type Pipe struct {
	flowChan chan int
	Traces   []int
	session  *define.Session
	errChan  chan error
	msgChan  chan []byte
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

func (o *Pipe) Run() {

	check()

	bc := initCfg()

	o.flowChan <- PIPE_NODE_ShowQr

	for {
		switch <-o.flowChan {
		case PIPE_NODE_ShowQr:
			api.ShowQr(bc.Lc)
			o.flowChan <- PIPE_NODE_Login
		case PIPE_NODE_Login:
			o.login(bc.Lc)
		case PIPE_NODE_WebWxNewLoginpage:
			o.webWxNewLoginPage(bc)
		case PIPE_NODE_WebWxInit:
			o.webWxInit()
		case PIPE_NODE_Listen:
			go o.listen()
		case PIPE_NODE_Customer:
			o.customer()
		case PIPE_NODE_Exit:
			break
		}
	}

	fmt.Println("obo is exit!")
}

func (o *Pipe) login(lc *define.LoginConfig) {
	tip := int64(1)
	retryTime := 10
	fmt.Println("Please scan the qrCode with wechat.")
	for i := 0; i < retryTime; i++ {
		code, response := api.ListenScan(tip, lc)
		switch code {
		case "201":
			fmt.Println("Please confirm login in wechat.")
			tip = 0
		case "200":
			rs := strings.Split(response, "\"")
			lc.Redirect = rs[1] + "&fun=new"
			o.flowChan <- PIPE_NODE_WebWxNewLoginpage
			o.after()
			return

		case "408":
			tip = 1
			retryTime -= 1
			time.Sleep(time.Second * 1)
		default:
			tip = 1
			retryTime -= 1
			time.Sleep(time.Second * 1)
		}
	}
}

func (o *Pipe) listen() {
	o.flowChan <- PIPE_NODE_Customer

	quitCurrClientCode := []string{"1100", "1101", "1102", "1205"}

	for {
		ret, sel := api.SyncCheck(o.session.Bc.Lpr, o.session.Bc.Br, o.session.Skl, utils.GetCookies(o.session))
		fmt.Printf("ret:" + ret + " " + "sel:" + sel + "\n")
		if tools.Find(ret, quitCurrClientCode) {
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
			fmt.Println(err.Error())
			//try reload login
			o.flowChan <- PIPE_NODE_Login

			return
		case msg := <-o.msgChan:
			msgMap := make(map[string]interface{})
			json.Unmarshal(msg, &msgMap)
			if int(msgMap["AddMsgCount"].(float64)) > 0 {
				utils.ParsingAddMsgList(msgMap["AddMsgList"].([]interface{}), func(message *define.ReceiveMessage) {
					plugins.Fire(o.session, message)
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
	fmt.Println("wx init begin...")
	resp := api.WebWxInit(o.session.Bc.Lpr, o.session.Bc.Br)
	 o.generateSyncKey(resp)
	go o.myself(resp)
	go o.initContactList("ContactList", resp)
	o.initFriends()
	o.initGroups()
	o.flowChan <- PIPE_NODE_Listen
	fmt.Println("wx init end...")

}

func (o *Pipe) generateSyncKey(result map[string]interface{}) {
	syncKey := result["SyncKey"].(map[string]interface{})
	o.session.Skl = api.SyncKey(syncKey)
}

func (o *Pipe) initContactList(key string, result map[string]interface{}) {

	cl := result[key].([]interface{})
	c := api.InitContactList(cl)
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

func (o *Pipe) after() {
	fmt.Println("Init  after")
}
