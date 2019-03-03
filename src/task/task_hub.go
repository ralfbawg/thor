package task

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type TaskHub struct {
	appName string

	appId string

	appKey string

	// Registered clients.
	clients map[string]*TaskClient

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *TaskClient

	// Unregister requests from clients.
	unregister chan *TaskClient

	mainChan chan []byte
}

func newTaskHub(appId string, appKey string) *TaskHub {
	return &TaskHub{
		appId:      appId,
		appKey:     appKey,
		broadcast:  make(chan []byte),
		register:   make(chan *TaskClient),
		unregister: make(chan *TaskClient),
		clients:    make(map[string]*TaskClient),
	}
}

func (h *TaskHub) run() {
	for {
		select {
		case client := <-h.register:

			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
