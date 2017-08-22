package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chanxuehong/wechat.v2/mp/core"
	//"github.com/chanxuehong/wechat.v2/mp/menu"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/request"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/response"
)

var (
	WECHAT_ADMIN string

	//WECHAT_TOKEN   = "12345678"
	//WECHAT_APPID   = "wxa38733da0e716f43"
	//WECHAT_AES_KEY = "PUqU7yIh4l0d1uye2nD1OEeM7n4lVYaE1OAehim2BP8"
	//WECHAT_ADMIN = "owKlz0UrhB0cqu9hQal2cwup4AyM"

	wechatServer *core.Server
)

func init() {
	WECHAT_TOKEN := os.Getenv("WECHAT_TOKEN")
	WECHAT_APPID := os.Getenv("WECHAT_APPID")
	WECHAT_AES_KEY := os.Getenv("WECHAT_AES_KEY")
	WECHAT_ADMIN = os.Getenv("WE_ADMIN")

	wemux := core.NewServeMux()
	wemux.MsgHandleFunc(request.MsgTypeText, coreHandleMessage)
	wemux.MsgHandleFunc(request.MsgTypeImage, coreHandleImage)
	wemux.DefaultEventHandleFunc(coreHandleEvent)

	wechatServer = core.NewServer("", WECHAT_APPID, WECHAT_TOKEN, WECHAT_AES_KEY, wemux, nil)
}

func WechatCallback(w http.ResponseWriter, r *http.Request) {
	wechatServer.ServeHTTP(w, r, nil)
}

func coreHandleMessage(ctx *core.Context) {
	log.Println("\nMESSAGE HANDLER INVOKED")

	msg := request.GetText(ctx.MixedMsg)
	data, err := parseContent(msg.Content)
	if err != nil {
		respText := "received text!"
		resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, respText)
		ctx.RawResponse(resp)

		pub := fmt.Sprintf("Type: \"%s\"  Content: \"%s\"  Response: \"%s\"  FromID: \"%s\"", msg.MsgType, msg.Content, respText, msg.FromUserName)
		subManager.Publish("wxstream", pub)
		//ctx.AESResponse(resp, 0, "", nil)
	}

	log.Println(data)
}

func coreHandleImage(ctx *core.Context) {
	log.Println("\nIMAGE HANDLER INVOKED")

	msg := request.GetImage(ctx.MixedMsg)
	respText := "received image!"
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, respText)
	ctx.RawResponse(resp)

	imageLink := `<a target="_blank" href="%s">image</a>`
	pub := fmt.Sprintf("Type: \"%s\"  Content: "+imageLink+"  Response: \"%s\"  FromID: \"%s\"", msg.MsgType, msg.PicURL, respText, msg.FromUserName)
	subManager.Publish("wxstream", pub)
}

func coreHandleEvent(ctx *core.Context) {
	log.Println("\nEVENT HANDLER INVOKED")

	switch ctx.MixedMsg.EventType {
	case request.EventTypeSubscribe:
		log.Println("Subscribe: " + ctx.MixedMsg.FromUserName)
		msg := request.GetSubscribeEvent(ctx.MixedMsg)
		respText := "welcome subscriber!"
		resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, respText)
		ctx.RawResponse(resp)

		pub := fmt.Sprintf("Type: \"%s\"  Response: \"%s\"  FromID: \"%s\"", msg.EventType, respText, msg.FromUserName)
		subManager.Publish("wxstream", pub)
		break
	case request.EventTypeUnsubscribe:
		log.Println("Unsubscribe: " + ctx.MixedMsg.FromUserName)
		msg := request.GetUnsubscribeEvent(ctx.MixedMsg)

		pub := fmt.Sprintf("Type: \"%s\"  FromID: \"%s\"", msg.EventType, msg.FromUserName)
		subManager.Publish("wxstream", pub)
		break
	}
}
