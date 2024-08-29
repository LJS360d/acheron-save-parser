package main

import (
	"acheron-save-parser/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/decode", handlers.HandleDecode)
	r.Run(":8080")
}
