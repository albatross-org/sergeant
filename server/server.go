package server

import (
	"github.com/albatross-org/sergeant"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var store *sergeant.Store

// Run starts the server.
func Run(s *sergeant.Store) {
	router = gin.Default()
	store = s

	initMiddleware()
	initRoutes()

	router.Run()
}
