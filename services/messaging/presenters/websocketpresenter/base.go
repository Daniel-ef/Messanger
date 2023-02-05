package websocketpresenter

import (
	"context"
	"github.com/google/uuid"
	kafkacontroller "github.com/messanger/services/messaging/controllers/kafka_controller"
	"github.com/messanger/services/messaging/transport/websocket"
)

type HubsMap map[uuid.UUID]*websocket.Hub

type WebsocketPresenter struct {
	servingClients HubsMap
	kafkaCtrl      *kafkacontroller.KafkaController
}

func NewWebsocketPresenter(ctx context.Context, kafkaConn *kafkacontroller.KafkaController) *WebsocketPresenter {
	wp := &WebsocketPresenter{
		kafkaCtrl:      kafkaConn,
		servingClients: make(HubsMap),
	}

	go wp.HandleMessageFromKafka(ctx)

	return wp
}
