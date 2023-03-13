package main

import (
	"Gin_framework/controller"
	"Gin_framework/service"
	"github.com/gin-gonic/gin"
)

var (
	videoService    service.VideoService       = service.New()
	videoController controller.VideoController = controller.New(videoService)
)

func main() {
	server := gin.Default()

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
