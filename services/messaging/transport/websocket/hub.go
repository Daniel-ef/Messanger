// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import "context"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
//
// Current hub should have only one client
type Hub struct {
	// Inbound messages from the clients.
	broadcast chan []byte

	client *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast: make(chan []byte),
	}
}

func (h *Hub) Register(client *Client) {
	h.client = client
}

func (h *Hub) Unregister() {
	close(h.client.send)
	h.client = nil
}

func (h *Hub) Run(ctx context.Context, messageHandler func(context.Context, []byte)) {
	for {
		select {
		case message := <-h.broadcast:
			messageHandler(ctx, message)
		}
	}
}

func (h *Hub) SendMessage(message []byte) {
	select {
	case h.client.send <- message:
	default:
		close(h.client.send)
		h.client = nil
	}
}
