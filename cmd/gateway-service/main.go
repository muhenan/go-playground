package main

import (
	"log"

	"github.com/cloudwego/kitex/client"
	"github.com/gin-gonic/gin"
	"go-playground/internal/microdemo/config"
	"go-playground/internal/microdemo/gateway"
	"go-playground/kitex_gen/microdemo/profileservice"
	"go-playground/kitex_gen/microdemo/recommendservice"
)

func main() {
	profileClient, err := profileservice.NewClient(
		config.ProfileServiceName,
		client.WithHostPorts(config.ProfileServiceAddr),
	)
	if err != nil {
		log.Fatalf("create profile client: %v", err)
	}

	recommendClient, err := recommendservice.NewClient(
		config.RecommendServiceName,
		client.WithHostPorts(config.RecommendServiceAddr),
	)
	if err != nil {
		log.Fatalf("create recommend client: %v", err)
	}

	r := gin.Default()
	handler := gateway.NewHandler(profileClient, recommendClient)
	handler.RegisterRoutes(r)

	log.Printf("gateway-service listening on http://%s", config.GatewayHTTPAddr)
	if err := r.Run(config.GatewayHTTPAddr); err != nil {
		log.Fatalf("gateway-service exited: %v", err)
	}
}
