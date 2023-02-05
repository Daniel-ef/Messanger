package websocketpresenter

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	kafkacontroller "github.com/messanger/services/messaging/controllers/kafka_controller"
	"github.com/messanger/services/messaging/logger"
	"github.com/messanger/services/messaging/transport/websocket"
	"go.uber.org/zap"
	"net/http"
)

// InitConnection handles websocket requests from the peer.
func (p *WebsocketPresenter) InitConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hub := websocket.NewHub()

	_, err := websocket.NewClient(w, r, hub)
	if err != nil {
		logger.Error(r.Context(), "could not create websocket client", zap.Error(err))
		return
	}
	go hub.Run(context.Background(), p.SendMessageToKafka)

	userIDRaw := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		logger.Error(r.Context(), "could not parse user_id", zap.Error(err))
		return
	}

	p.servingClients[userID] = hub

	logger.Info(ctx, "got new client", zap.String("user_id", userID.String()))
}

type InboundMessage struct {
	Type     string    `json:"type"`
	ChatID   uuid.UUID `json:"chat_id"`
	SendTo   uuid.UUID `json:"send_to"`
	SendFrom uuid.UUID `json:"send_from"`
	Text     string    `json:"text"`
}

func (p *WebsocketPresenter) SendMessageToKafka(ctx context.Context, message []byte) {
	logger.Info(ctx, "got message", zap.ByteString("message", message))

	// send to Kafka
	err := p.kafkaCtrl.SendMessage(ctx, kafkacontroller.TopicMessages, message)
	if err != nil {
		logger.Error(ctx, "could not write message to kafka_controller", zap.Error(err))
		return
	}
}

func (p *WebsocketPresenter) HandleMessageFromKafka(ctx context.Context) {
	for {
		message, err := p.kafkaCtrl.ReadMessage(ctx, kafkacontroller.TopicMessages)
		if err != nil {
			logger.Error(ctx, "could not read message from kafka_controller", zap.Error(err))
			continue
		}
		logger.Info(ctx, "got message from kafka_controller", zap.ByteString("message", message.Value))

		var inboundMessage InboundMessage
		err = json.Unmarshal(message.Value, &inboundMessage)
		if err != nil {
			logger.Error(ctx, "could not unmarshal message", zap.Error(err))
			continue
		}

		if hub, ok := p.servingClients[inboundMessage.SendTo]; ok {
			hub.SendMessage(message.Value)

			logger.Info(ctx, "message sent to client", zap.String("user_id", inboundMessage.SendTo.String()))
		} else {
			logger.Warn(ctx, "this pod doesn't serve the client", zap.String("user_id", inboundMessage.SendTo.String()))
		}
	}
}
