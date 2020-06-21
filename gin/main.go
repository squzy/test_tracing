package main

import (
	"github.com/gin-gonic/gin"
	"github.com/squzy/squzy_go/core"
	sg "github.com/squzy/squzy_go/integrations/gin"
	sGrpc "github.com/squzy/squzy_go/integrations/grpc"
	"google.golang.org/grpc"
	"log"
	service "test_app/generated"
)

func main() {
	squzy, err := core.CreateApplication(nil, &core.Options{
		ApiHost:         "http://localhost:8080",
		ApplicationName: "Go app test gin",
	})

	if err != nil {
		log.Fatal(err)
	}

	engine := gin.New()
	engine.Use(gin.Recovery(), sg.New(squzy))
	conn, err := grpc.Dial("localhost:7879", grpc.WithInsecure(), grpc.WithUnaryInterceptor(sGrpc.NewClientUnaryInterceptor(squzy)))
	if err != nil {
		log.Fatal(err)
	}
	clint := service.NewEchoServiceClient(conn)
	engine.GET("hello", func(context *gin.Context) {
		res, err := clint.Echo(context, &service.EchoMsg{})
		if err != nil {
			context.AbortWithError(200, err)
			return
		}
		context.JSON(200, res)
	})

	log.Fatal(engine.Run(":7878"))
}
