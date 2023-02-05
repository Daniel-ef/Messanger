package messagingcontroller

import (
	"context"
	"github.com/google/uuid"
	"github.com/messanger/services/messaging/transport/kafkatransport"
	"github.com/messanger/services/messaging/transport/websocket"
)

type HubsMap map[uuid.UUID]*websocket.Hub

type Service struct {
	servingClients HubsMap
	kafkaCtrl      *kafkacontroller.KafkaController
}

func NewService(ctx context.Context, kafkaConn *kafkacontroller.KafkaController) *Service {
	s := &Service{
		kafkaCtrl:      kafkaConn,
		servingClients: make(HubsMap),
	}

	go s.HandleMessageFromKafka(ctx)

	return s
}
