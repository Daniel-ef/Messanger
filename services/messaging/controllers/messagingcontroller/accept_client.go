package messagingcontroller

import (
	"context"
	"github.com/google/uuid"
	"github.com/messanger/services/messaging/logger"
	"github.com/messanger/services/messaging/transport/websocket"
	"go.uber.org/zap"
)

func (s *Service) AcceptClient(ctx context.Context, userID uuid.UUID, hub *websocket.Hub) {
	s.servingClients[userID] = hub

	go hub.Run(context.Background(), s.SendMessageToKafka)

	logger.Info(ctx, "got new client", zap.String("user_id", userID.String()))
}
