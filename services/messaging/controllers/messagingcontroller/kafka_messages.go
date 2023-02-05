package messagingcontroller

import (
	"context"
	"encoding/json"
	"github.com/messanger/services/messaging/transport/kafkatransport"

	"go.uber.org/zap"

	"github.com/messanger/services/messaging/logger"
)

func (s *Service) SendMessageToKafka(ctx context.Context, message []byte) {
	logger.Info(ctx, "got message", zap.ByteString("message", message))

	// send to Kafka
	err := s.kafkaCtrl.SendMessage(ctx, kafkacontroller.TopicMessages, message)
	if err != nil {
		logger.Error(ctx, "could not write message to kafkatransport", zap.Error(err))
		return
	}
}

func (s *Service) HandleMessageFromKafka(ctx context.Context) {
	for {
		message, err := s.kafkaCtrl.ReadMessage(ctx, kafkacontroller.TopicMessages)
		if err != nil {
			logger.Error(ctx, "could not read message from kafkatransport", zap.Error(err))
			continue
		}

		var inboundMessage InboundMessage
		err = json.Unmarshal(message.Value, &inboundMessage)
		if err != nil {
			logger.Error(ctx, "could not unmarshal message", zap.Error(err))
			continue
		}

		if hub, ok := s.servingClients[inboundMessage.SendTo]; ok {
			hub.SendMessage(message.Value)

			logger.Info(ctx, "message sent to client", zap.String("user_id", inboundMessage.SendTo.String()))
		} else {
			logger.Warn(ctx, "this pod doesn't serve the client", zap.String("user_id", inboundMessage.SendTo.String()))
		}
	}
}
