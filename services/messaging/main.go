// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	kafkacontroller "github.com/messanger/services/messaging/controllers/kafka_controller"
	"github.com/messanger/services/messaging/logger"
	"github.com/messanger/services/messaging/presenters/websocketpresenter"
)

var addr = flag.String("addr", ":8080", "http_presenter service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.Context(), "serving home page")
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	logger.SetNewGlobalLoggerOnce(logger.Config{LogLevel: "DEBUG", Format: "json", Sink: "stdout"})
	defer logger.Close()

	ctx := context.Background()

	kafkaCtr, closer := kafkacontroller.NewKafkaController(ctx, "localhost:29092")

	websocketPresenter := websocketpresenter.NewWebsocketPresenter(ctx, kafkaCtr)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/init", websocketPresenter.InitConnection)

	logger.Info(ctx, fmt.Sprintf("Starting server on %s", *addr))

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		logger.Fatal(ctx, "could not start server", zap.Error(err))
	}

	shutdownCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	closer(shutdownCtx)
}
