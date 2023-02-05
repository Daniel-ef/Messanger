package httppresenter

import (
	"github.com/google/uuid"
	"github.com/messanger/services/messaging/transport/messagingapi"
	"github.com/messanger/services/messaging/transport/websocket"
)

type MessagingPresenter struct {
	servingClients *map[uuid.UUID]*websocket.Client
}

func NewMessagingPresenter(servingClients *map[uuid.UUID]*websocket.Client) *MessagingPresenter {
	return &MessagingPresenter{servingClients: servingClients}
}

var _ messagingapi.Handler = &MessagingPresenter{}
