package main

import (
	"log"
	"net"

	server "github.com/cloudwego/kitex/server"
	"go-playground/internal/microdemo/config"
	"go-playground/internal/microdemo/recommend"
	"go-playground/kitex_gen/microdemo/recommendservice"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", config.RecommendServiceAddr)
	if err != nil {
		log.Fatalf("resolve recommend address: %v", err)
	}

	svr := recommendservice.NewServer(
		recommend.NewService(),
		server.WithServiceAddr(addr),
	)

	log.Printf("recommend-service listening on %s", config.RecommendServiceAddr)
	if err := svr.Run(); err != nil {
		log.Fatalf("recommend-service exited: %v", err)
	}
}
