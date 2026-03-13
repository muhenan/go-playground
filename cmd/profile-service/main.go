package main

import (
	"log"
	"net"

	server "github.com/cloudwego/kitex/server"
	"go-playground/internal/microdemo/config"
	"go-playground/internal/microdemo/profile"
	"go-playground/kitex_gen/microdemo/profileservice"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", config.ProfileServiceAddr)
	if err != nil {
		log.Fatalf("resolve profile address: %v", err)
	}

	svr := profileservice.NewServer(
		profile.NewService(),
		server.WithServiceAddr(addr),
	)

	log.Printf("profile-service listening on %s", config.ProfileServiceAddr)
	if err := svr.Run(); err != nil {
		log.Fatalf("profile-service exited: %v", err)
	}
}
