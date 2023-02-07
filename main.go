package main

import "github.com/gin-gonic/gin"

func main() {
	engine := gin.Default()
	engine.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{"msg": "Hello World"})
	})
	engine.Run()
}
