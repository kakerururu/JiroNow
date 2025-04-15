package main

import (
	"JiroNow/env"
	"JiroNow/handler"
	"JiroNow/infra"
	"JiroNow/infra/x"
	"JiroNow/oas"
	"log"
	"net/http"
)

func main() {

	envs, err := env.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	xClient := x.NewXClient(envs.XBearerToken, http.DefaultClient)
	geminiClient := infra.NewGeminiClient(envs.GeminiToken, http.DefaultClient)

	h := handler.NewHandler(xClient, geminiClient)

	server, err := oas.NewServer(h)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s ...", envs.Port)
	if err := http.ListenAndServe(":"+envs.Port, server); err != nil {
		log.Fatal(err)
	}
}
