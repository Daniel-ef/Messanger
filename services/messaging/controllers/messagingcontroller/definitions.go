package messagingcontroller

import "github.com/google/uuid"

type InboundMessage struct {
	Type     string    `json:"type"`
	ChatID   uuid.UUID `json:"chat_id"`
	SendTo   uuid.UUID `json:"send_to"`
	SendFrom uuid.UUID `json:"send_from"`
	Text     string    `json:"text"`
}
