package websocketpresenter

import (
	"github.com/google/uuid"
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

	userIDRaw := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		logger.Error(r.Context(), "could not parse user_id", zap.Error(err))
		return
	}

	p.messagingCtrl.AcceptClient(ctx, userID, hub)
}
