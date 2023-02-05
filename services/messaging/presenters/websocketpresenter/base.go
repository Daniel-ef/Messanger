package websocketpresenter

import (
	"github.com/messanger/services/messaging/controllers/messagingcontroller"
)

type WebsocketPresenter struct {
	messagingCtrl *messagingcontroller.Service
}

func NewWebsocketPresenter(messagingCtrl *messagingcontroller.Service) *WebsocketPresenter {
	return &WebsocketPresenter{
		messagingCtrl: messagingCtrl,
	}
}
