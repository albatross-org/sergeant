package server

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

// Run starts the server.
func Run() {
	router = gin.Default()

	initMiddleware()
	initRoutes()

	router.Run()
}
