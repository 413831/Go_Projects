package main

import (
	"Gin_framework/controller"
	"Gin_framework/middlewares"
	"Gin_framework/service"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
	"io"
	"os"
)

var (
	videoService    service.VideoService       = service.New()
	videoController controller.VideoController = controller.New(videoService)
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	setupLogOutput()

	server := gin.New()

	server.Use(gin.Recovery(), middlewares.Logger(), middlewares.BasicAuth(), gindump.Dump())

	server.GET("/test", check)
	server.GET("/", findAll)
	server.POST("/save", save)

	server.Run(":8080")
}

func save(ctx *gin.Context) {
	ctx.JSON(200, videoController.Save(ctx))
}

func findAll(ctx *gin.Context) {
	ctx.JSON(200, videoController.FindAll())
}

func check(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "OK!!",
	})
}
