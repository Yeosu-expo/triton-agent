package handler

import (
	"net/http"

	"github.com/ahr-i/triton-agent/tritonController"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
}

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler: mux,
	}

	downloaded := make(chan string)
	channel := &downloaded

	go tritonController.Seeding("filePath", channel)

	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")                                                                                                // Ping check
	mux.HandleFunc("/model/{model:[a-z-_]+}/{version:[0-9]+}/infer", handler.inferHandler).Methods("POST")                                                     // Inference version 1.0
	mux.HandleFunc("/model/{model:[a-z-_]+}/{version:[0-9]+}/ready", handler.readyHandler).Methods("GET")                                                      // Model check
	mux.HandleFunc("/repository/index", handler.repositoryIndexHandler).Methods("POST")                                                                        // Get Triton repository index
	mux.HandleFunc("/provider/{provider:[0-9a-zA-Z-_]+}/model/{model:[0-9a-zA-Z-_]+}/{version:[0-9]+}/infer", handler.inferV2Handler).Methods("POST")          // Inference version 2.0
	mux.HandleFunc("/provider/{provider:[0-9a-zA-Z-_]+}/model/{model:[0-9a-zA-Z-_]+}/{version:[0-9]+}/infer/test", handler.testInferV2Handler).Methods("POST") // Inference version 2.0
	mux.HandleFunc("/serving", handler.servingHandler).Methods("POST")                                                                                         // Model serving API

	return handler
}
