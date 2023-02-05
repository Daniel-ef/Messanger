package httppresenter

import (
	"context"
	"github.com/messanger/services/messaging/transport/messagingapi"
)

func (p *MessagingPresenter) InitConnection(
	ctx context.Context,
	params messagingapi.InitConnectionParams,
) (messagingapi.InitConnectionRes, error) {
	return &messagingapi.InitConnectionOK{}, nil
}
